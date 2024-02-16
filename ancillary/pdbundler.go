package newt

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	// 3rd party packages
	"gopkg.in/yaml.v3"
)

// PandocBundler models the application `pdbundler`
type PandocBundler struct {
	Port string `json:"port,omitempty" yaml:"port"`
	Templates []*PandocTemplate `json:"templates,omitempty" yaml:"templates"`
}

// PandocTemplate hold the request to template mapping for `pdbundler`
type PandocTemplate struct {
	// Pattern holds a request pattern, e.g. `POST /blog_post`
	// A request is associated with a template to be bundled into
	// an JSON object. The pattern conforms to Go 1.22 or later's
	// HTTP handler function pattern, see <https://tip.golang.org/doc/go1.22#enhanced_routing_patterns>
	Pattern string `json:"request,required" yaml:"request"`

	// Template holds a path to the primary template file for this route. Path can be relative
	// to the current working directory.
	Template string `json:"template,required" yaml:"template"`

	// Options hold the JSON object that will be resolve by `tmplbndl`. The values `.text` and `.template`
	// will be replaced by the contents specified in Bundles and received in the request.
	Options map[string]interface{} `json:"options,omitempty" yaml:"options"`

	// Debug logs more verbosely if true
	Debug bool `json:"debug,omitempty" yaml:"debug"`

	// Src holds the resolved template content
	Src []byte

	// Vars holds the names of any variables expressed in the pattern, these an be used to replace elements of
	// the output object.
	Vars []string
}

// NewPandocBundler create a new PandocBundler struct. If a filename
// is provided it reads the file and sets things up accordingly.
func NewPandocBundler(fName string, cfg *Config) (*PandocBundler, error) {
	pb := &PandocBundler{}
	if fName == "" {
		return pb, fmt.Errorf("missing configuration file")
	}
	src, err := os.ReadFile(fName)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(src, &pb)
	if err != nil {
		return nil, err
	}
	if pb.Port == "" && cfg.Application != nil &&  cfg.Application.Port != 0 {
		pb.Port = fmt.Sprintf(":%d", cfg.Application.Port)
	}
	// Prefix the port number with a colon
	if ! strings.HasPrefix(pb.Port, ":") {
		pb.Port = fmt.Sprintf(":%s", pb.Port)
	}
	return pb, nil
}

// ResolvesTemplate. If not template name is available it is assumed
// you're going to use one of Pandoc default templates. If a name is
// provided then it reads the file saving the results in `.Src`
// An error is returned in a problem is encountered.
func (pt *PandocTemplate) ResolveTemplate() error {
	if pt.Template != "" {
 		src, err := os.ReadFile(pt.Template)
		if err != nil {
			return err
		}
		pt.Src = src
	}
	return nil
}

// Handler provides a HandleFunc for use with an http.ServeMux struct.
func (pt *PandocTemplate) Handler(w http.ResponseWriter, r *http.Request) {
	if pt.Debug {
		log.Printf("DEBUG .Handler(w, %s %s)", r.Method, r.URL.Path)
	}
	obj := map[string]interface{}{}
	// Copy in our options
	if pt.Options != nil {
		for k, v := range pt.Options {
			obj[k] = v
		}
	}
	// Merge in path values into options
	for _, varname := range pt.Vars {
		val := r.PathValue(varname)
		if val != "" {
			obj[varname] = val
		}
	}
	if pt.Debug {
		log.Printf("DEBUG varnames -> %+v\n", pt.Vars)
		log.Printf("DEBUG obj now -> %+v\n", obj)
	}
	// Add our resolved template source
	obj["template"] = fmt.Sprintf("%s", pt.Src)
	// Add our `.text` value from the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// Need to return an HTTP 400 error status. r.Body wasn't welformed
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	body = bytes.TrimSpace(body)
	// Is the body a JSON object or some other binary object?
	switch {
		case bytes.HasPrefix(body, []byte("{")):
			// Handle case of JSON object
			m := map[string]interface{}{}
			if err := json.Unmarshal(body, &m); err != nil {
				// Need to return an HTTP 400 error status. r.Body wasn't welformed
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("%s, %s", http.StatusText(http.StatusBadRequest), err)))
				return
			}
			src, _ := yaml.Marshal(m)
			obj["text"] = fmt.Sprintf("%s", src)
		case bytes.HasPrefix(body, []byte("[")):
			// Handle case of JSON array
			a := []interface{}{}
			if err := json.Unmarshal(body, &a); err != nil {
				// Need to return an HTTP 400 error status. r.Body wasn't welformed
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("%s, %s", http.StatusText(http.StatusBadRequest), err)))
				return
			}
			src, _ := yaml.Marshal(a)
			obj["text"] = fmt.Sprintf("%s", src)
		default:
			// If we a text format then we just attach to `.text` in obj. Otherwise we need to encode it.
			contentType := http.DetectContentType(body)
			if strings.HasPrefix(contentType, "text/") {
				obj["text"] = body
			} else {
				dst := make([]byte, base64.StdEncoding.EncodedLen(len(body)))
				base64.StdEncoding.Encode(dst, body)
				obj["text"] = dst
			}
	}
	// Build a respond with a object and JSON encode it.
	src, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		// Need to return an HTTP 400 error status. r.Body wasn't welformed
		w.WriteHeader(http.StatusUnprocessableEntity)
		// FEAT: Be nice to allow a debug hook to debug content problems before writing error ...
		w.Write([]byte(http.StatusText(http.StatusUnprocessableEntity)))
		return
	}
	w.Write(src)
}


// ResolvePath reviews the `.Request` attribute and updates the Vars using PatternKeys()
func (pt *PandocTemplate) ResolvePath() error {
	// Does the `.Request` hold a pattern or a fixed string?
	if strings.Contains(pt.Pattern, "{") {
		if ! strings.Contains(pt.Pattern, "}") {
			return fmt.Errorf("%q is malformed", pt.Pattern)
		}
		// Record our list of var names so handler can override the object being constructed from a path.
		pt.Vars = PatternKeys(pt.Pattern)
	}
	if pt.Debug {
		log.Printf("DEBUG assigning pt.Pattern -> %q\n", pt.Pattern)
		log.Printf("DEBUG vars -> %+v\n", pt.Vars)
	}
	return nil
}

func (pb *PandocBundler) ListenAndServe() error {
	mux := http.NewServeMux()
	for _, bndl := range pb.Templates {
		mux.HandleFunc(bndl.Pattern, func(w http.ResponseWriter, r *http.Request) {
			if bndl.Debug {
				log.Printf("DEBUG mux.HandleFunc(%q, bndl.Handler)", bndl.Pattern)
				log.Printf("DEBUG .vars -> %+v", bndl.Vars)
			}
			bndl.Handler(w, r)
		})
	}
	// Now create my http server
	svr := new(http.Server)
	svr.Addr = pb.Port
	svr.Handler = mux
	log.Printf("%s listening on %s", path.Base(os.Args[0]), svr.Addr)
	if err := svr.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	return nil
}

