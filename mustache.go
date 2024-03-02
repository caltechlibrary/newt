package newt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	// 3rd Party Templates
	"github.com/cbroglie/mustache"
	"gopkg.in/yaml.v3"
)

// NewtMustache defines the `newtmustache` application configuration YAML
type NewtMustache struct {
	Port      string
	Templates []*MustacheTemplate
	Timeout   time.Duration
}

// MustacheTemplate hold the request to template mapping for NewtMustache struct
type MustacheTemplate struct {
	// Pattern holds a request path, e.g. `/blog_post`. NOTE: the method is ignored. A POST
	// is presumed to hold data that will be processed by the template engine. A GET retrieves the
	// unresolved template.
	Pattern string `json:"request,required" yaml:"request"`

	// Template holds a path to the primary template file for this route. Path can be relative
	// to the current working directory.
	Template string `json:"template,required" yaml:"template"`

	// Partials holds a list of path to partial templates used by the primary template. `newtmustache` will
	// attempt to replace references in the primary template with the content of the partials. Recursive
	// partials are not supported. The goal is to facilate including sub templates.
	Partials []string `json:"partials,omitempty" yaml:"partials"`

	// Options hold the JSON object that will be resolve by `newtmustache`. The values `.text` and `.template`
	// will be replaced by the contents specified in Bundles and received in the request.
	Options map[string]interface{} `json:"options,omitempty" yaml:"options"`

	// Debug logs more verbosely if true
	Debug bool `json:"debug,omitempty" yaml:"debug"`

	// Tmpl holds the parsed template
	Tmpl *mustache.Template

	// Vocabulary holds the path to a YAML file used to populate Vocabulary at startup.
	Vocabulary string `json:"vocabulary,omitempty" yaml:"vocabulary"`

	// Voc holds a map of variable names to values. It is read in when NewtMustache starts from a separate YAML
	// file.
	Voc map[string]interface{}

    // Vars holds the names of any variables expressed in the pattern, these an be used to replace elements of
    // the output object.
    Vars []string
}

// NewNewtMustache create a new NewtMustache struct. If a filename
// is provided it reads the file and sets things up accordingly.
func NewNewtMustache(cfg *Config) (*NewtMustache, error) {
	nm := &NewtMustache{
		Templates: cfg.Templates,
	}
	if cfg.Applications.NewtMustache.Port != 0 {
		nm.Port = fmt.Sprintf(":%d", cfg.Applications.NewtMustache.Port)
	}
	if cfg.Applications.NewtMustache.Timeout != 0 {
		nm.Timeout = cfg.Applications.NewtMustache.Timeout * time.Second
	}
	return nm, nil
}

// LoadVocabary retrieves the YAML file contents found in .VocabularFName and builds the map[string]interface{} that
// holds .Vocabulary
func (mt *MustacheTemplate) LoadVocabulary() error {
	voc := map[string]interface{}{}
	if mt.Vocabulary != "" {
		src, err := os.ReadFile(mt.Vocabulary)
		if err != nil {
			return err
		}
		if err := yaml.Unmarshal(src, &voc); err != nil {
			return err
		}
	}
	if len(voc) > 0 {
		mt.Voc = voc
	}
	return nil
}

// ResolvePath reviews the `.Request` attribute and updates the Vars using PatternKeys()
func (mt *MustacheTemplate) ResolvePath() error {
	// Does the `.Request` hold a pattern or a fixed string?
	if strings.Contains(mt.Pattern, "{") {
		if !strings.Contains(mt.Pattern, "}") {
			return fmt.Errorf("%q is malformed", mt.Pattern)
		}
		// Record our list of var names so handler can override the object being constructed from a path.
		mt.Vars = PatternKeys(mt.Pattern)
	}
	if mt.Debug {
		log.Printf("assigning mt.Pattern -> %q\n", mt.Pattern)
		log.Printf("vars -> %+v\n", mt.Vars)
	}
	return nil
}

// ResolvesTemplate is responsible for reading and parse the template and partials associated with a mapped request.
// If an error is encountered a error value is returned.
func (mt *MustacheTemplate) ResolveTemplate() error {
	if mt.Template != "" {
		if len(mt.Partials) > 0 {
			if mt.Debug {
				log.Printf("handling primary and partial templates")
			}
			sp := mustache.StaticProvider{}
			sp.Partials = map[string]string{}
			for _, fName := range mt.Partials {
				if mt.Debug {
					log.Printf("attempting to read %q", fName)
				}
				src, err := os.ReadFile(fName)
				if err != nil {
					return err
				}
				name := strings.TrimSuffix(path.Base(fName), path.Ext(fName))
				sp.Partials[name] = fmt.Sprintf("%s", src)
			}
			if mt.Debug {
				log.Printf("attempting to parse template with partials %q", mt.Template)
			}
			src, err := os.ReadFile(mt.Template)
			if err != nil {
				return err
			}
			txt := fmt.Sprintf("%s", src)
			tmpl, err := mustache.ParseStringPartials(txt, &sp)
			if err != nil {
				return err
			}
			mt.Tmpl = tmpl
			if mt.Debug {
				log.Printf("templates and partials parsed successfully")
			}
			return nil
		}
		if mt.Debug {
			log.Printf("attempting to parse single template %q", mt.Template)
		}
		tmpl, err := mustache.ParseFile(mt.Template)
		if err != nil {
			return err
		}
		mt.Tmpl = tmpl
		if mt.Debug {
			log.Printf("templates parsed successfully")
		}
		return nil
	}
	return fmt.Errorf("no template found")
}

// Handler decodes a the request body and then processes that as a Mustache template.
func (mt *MustacheTemplate) Handler(w http.ResponseWriter, r *http.Request) {
	if mt.Debug {
		log.Printf(".Handler(w, %s %s)", r.Method, r.URL.Path)
	}

	src, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if mt.Debug {
			log.Printf("failed to read request body, %s\n", err)
		}
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// If we have src to decode, let's decode it.
	var (
		body *interface{}
	)
	if len(src) > 0 {
		dec := json.NewDecoder(bytes.NewBuffer(src))
		dec.UseNumber()
		if err := dec.Decode(&body); err != nil && err != io.EOF {
			if mt.Debug {
				log.Printf("failed to decode JSON Response body, %s", err)
			}
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}
		if mt.Debug {
			log.Printf("obj populated from request body, %+v", body)
		}
		if body == nil {
			if mt.Debug {
				log.Printf("no data for template processing")
			}
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}
	}
	options := map[string]interface{}{}
	vocabulary := map[string]interface{}{}
	vars := map[string]string{}
	// Copy in the options into page objcet's options.
	if mt.Options != nil {
		if mt.Debug {
			log.Printf("options -> %+v\n", mt.Options)
		}
		for k, v := range mt.Options {
			options[k] = v
		}
		if mt.Debug {
			log.Printf("obj after processing options -> %+v", options)
		}
	}
	if mt.Voc != nil {
		if mt.Debug {
			log.Printf("vocabulary -> %+v\n", mt.Voc)
		}
		for k, v := range mt.Voc {
			vocabulary[k] = v
		}
		if mt.Debug {
			log.Printf("obj after processing vocabulary -> %+v", vocabulary)
		}
	}
	// Merge in path values into .vars
	if len(mt.Vars) > 0 {
		if mt.Debug {
			log.Printf("varnames -> %+v\n", mt.Vars)
		}
		for _, varname := range mt.Vars {
			val := r.PathValue(varname)
			if val != "" {
				// val presidence over mt.Options
				vars[varname] = val
			}
		}
		if mt.Debug {
			log.Printf("obj after processing varnames -> %+v", vars)
		}
	}
	obj := map[string]interface{}{
		"body": body,
		"options": options,
		"vocabulary": vocabulary,
		"vars": vars,
	}
	if mt.Debug {
		log.Printf("obj after processing options -> %+v", obj)
	}
	// Handle case where some how the service was started before setting up template processing
	if mt.Tmpl == nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	if mt.Debug {
		log.Printf("mt.Tmpl -> %+v", mt.Tmpl)
	}
	// We want to write to a buffer so we can do content detection and set the headers correctly.
	buf := bytes.NewBuffer([]byte{})
	mt.Tmpl.FRender(buf, obj)
	src = buf.Bytes()
	contentType := http.DetectContentType(src)
	if bytes.HasPrefix(src, []byte("<!DOC")) {
		contentType = "text/html; charset=utf-8"
	}
	if mt.Debug {
		log.Printf("content type: %q -> %q", contentType, src)
	}
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.Write(src)
}

func (nm *NewtMustache) ListenAndServe() error {
	mux := http.NewServeMux()
	// Setup our handlers, POST for process data with the template and GET to retreive the template
	// source.
	for _, mt := range nm.Templates {
		if err := mt.LoadVocabulary(); err != nil {
			log.Fatal(err)
		}
		// Process the data with template if a POST.
		mux.HandleFunc("POST " + mt.Pattern, func(w http.ResponseWriter, r *http.Request) {
			if mt.Debug {
				log.Printf("mux.HandleFunc(%q, mt.Handler)", "POST " + mt.Pattern)
				log.Printf(".vars -> %+v", mt.Vars)
			}
			mt.Handler(w, r)
		})
	}
	// Now create my http server
	svr := new(http.Server)
	svr.Addr = nm.Port
	svr.Handler = mux
	if err := svr.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	return nil
}
