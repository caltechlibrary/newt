package newt

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"io"
	"strings"

	// 3rd party packages
	"gopkg.in/yaml.v3"
)

// PandocBundler models the application `pdbundler`
type PandocBundler struct {
	Port string `json:"port,omitempty" yaml:"port"`
	Templates []*PandocBundle `json:"templates,omitempty" yaml:"templates"`
}

// PandocBundle hold the request to template mapping for `pdbundler`
type PandocBundle struct {
	// Pattern holds a request pattern, e.g. `POST /blog_post`
	// A request is associated with a template to be bundled into
	// an JSON object. The pattern conforms to Go 1.22 or later's
	// HTTP handler function pattern, see <https://tip.golang.org/doc/go1.22#enhanced_routing_patterns>
	Pattern string `json:"request,required" yaml:"request"`
	// Template holds a path to the primary template file for this route. Path can be relative
	// to the current working directory.
	Template string `json:"template,required" yaml:"template"`
	// Partials holds a list of path to partial templates used by the primary template. `tmplbndl` will
	// attempt to replace references in the primary template with the content of the partials. Recursive
	// partials are not supported. The goal is to facilate including sub templates.
	Partials []string `json:"partials,omitempty" yaml:"partials"`
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
	tb := &PandocBundler{}
	if fName == "" {
		return tb, fmt.Errorf("missing configuration file")
	}
	src, err := os.ReadFile(fName)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(src, &tb)
	if err != nil {
		return nil, err
	}
	if tb.Port == "" && cfg.Application != nil &&  cfg.Application.Port != ""{
		tb.Port = cfg.Application.Port	
	}
	// Prefix the port number with a colon
	if ! strings.HasPrefix(tb.Port, ":") {
		tb.Port = fmt.Sprintf(":%s", tb.Port)
	}
	return tb, nil
}

// ResolvesTemplate. If not template name is available it is assumed
// you're going to use one of Pandoc default templates. If a name is
// provided then it reads the file saving the results in `.Src`
// An error is returned in a problem is encountered.
func (b *PandocBundle) ResolveTemplate() error {
	if b.Template != "" {
 		src, err := os.ReadFile(b.Template)
		if err != nil {
			return err
		}
		b.Src = src
	}
	return nil
}

// Handler provides a HandleFunc for use with an http.ServeMux struct.
func (b *PandocBundle) Handler(w http.ResponseWriter, r *http.Request) {
	if b.Debug {
		log.Printf("DEBUG .Handler(w, %s %s)", r.Method, r.URL.Path)
	}
	obj := map[string]interface{}{}
	// Copy in our options
	if b.Options != nil {
		for k, v := range b.Options {
			obj[k] = v
		}
	}
	// Merge in path values into options
	for _, varname := range b.Vars {
		val := r.PathValue(varname)
		if val != "" {
			obj[varname] = val
		}
	}
	if b.Debug {
		log.Printf("DEBUG varnames -> %+v\n", b.Vars)
		log.Printf("DEBUG obj now -> %+v\n", obj)
	}
	// Add our resolved template source
	obj["template"] = fmt.Sprintf("%s", b.Src)
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
			obj["text"] = m
		case bytes.HasPrefix(body, []byte("[")):
			// Handle case of JSON array
			a := []interface{}{}
			if err := json.Unmarshal(body, &a); err != nil {
				// Need to return an HTTP 400 error status. r.Body wasn't welformed
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(fmt.Sprintf("%s, %s", http.StatusText(http.StatusBadRequest), err)))
				return
			}
			obj["text"] = a
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
func (b *PandocBundle) ResolvePath() error {
	// Does the `.Request` hold a pattern or a fixed string?
	if strings.Contains(b.Pattern, "{") {
		if ! strings.Contains(b.Pattern, "}") {
			return fmt.Errorf("%q is malformed", b.Pattern)
		}
		// Record our list of var names so handler can override the object being constructed from a path.
		b.Vars = PatternKeys(b.Pattern)
	}
	if b.Debug {
		log.Printf("DEBUG assigning b.Pattern -> %q\n", b.Pattern)
		log.Printf("DEBUG vars -> %+v\n", b.Vars)
	}
	return nil
}

func (tb *PandocBundler) ListenAndServe() error {
	mux := http.NewServeMux()
	for _, bndl := range tb.Templates {
		mux.HandleFunc(bndl.Pattern, func(w http.ResponseWriter, r *http.Request) {
			if bndl.Debug {
				log.Printf("DEBUG mux.HandleFunc(%q, bndl.Handler)", bndl.Pattern)
				log.Printf("DEBUG .vars -> %+v", bndl.Vars)
			}
			bndl.Handler(w, r)
		})
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			fmt.Fprintf(w, "\nSomething fell through -> %q %q\n", r.Method, r.URL.Path)
			//http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page!")
	})
	// Now create my http server
	svr := new(http.Server)
	svr.Addr = tb.Port
	svr.Handler = mux
	if err := svr.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	return nil
}

