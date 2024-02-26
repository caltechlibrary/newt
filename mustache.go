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
)

// NewtMustache defines the `newtmustache` application configuration YAML
type NewtMustache struct {
	Port      string
	Templates []*MustacheTemplate
	Timeout   time.Duration
}

// MustacheTemplate hold the request to template mapping for NewtMustache struct
type MustacheTemplate struct {
	// Pattern holds a request pattern, e.g. `POST /blog_post`. If the METHOD is not specified a POST is assumed.
	// A request is associated with a template to be bundled into
	// an JSON object. The pattern conforms to Go 1.22 or later's
	// HTTP handler function pattern, see <https://tip.golang.org/doc/go1.22#enhanced_routing_patterns>
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
	// FIXME: Think about what it means if a GET, HEAD, PUT, DELETE are to be handled.
	obj := map[string]interface{}{}
	src, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if mt.Debug {
			log.Printf("failed to read request body, %s\n", err)
		}
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	// If we have src to decode, let's decode it.
	if len(src) > 0 {
		dec := json.NewDecoder(bytes.NewBuffer(src))
		dec.UseNumber()
		if err := dec.Decode(&obj); err != nil && err != io.EOF {
			if mt.Debug {
				log.Printf("failed to decode JSON Response body, %s", err)
			}
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}
		if mt.Debug {
			log.Printf("obj populated from request body, %+v", obj)
		}
	}
	params := r.URL.Query()
	if len(params) > 0 {
		// Let's check if data came in as query paramters and add it to our object.
		if mt.Debug {
			log.Printf("URL Query parameters -> %+v", params)
		}
		for k, v := range params {
			if k != "" {
				// Take the first value set (e.g. in POST or GET QUERY parameters)
				if _, conflict := obj[k]; !conflict {
					obj[k] = v
				}
			}
		}
		if mt.Debug {
			log.Printf("obj after processing query parameters, %+v", obj)
		}
	}
	// Merge (without overwriting our POST content) in our options into obj
	if mt.Options != nil {
		if mt.Debug {
			log.Printf("options -> %+v\n", mt.Options)
		}
		for k, v := range mt.Options {
			// Options take presidence over POST or GET QUERY parameters.
			obj[k] = v
		}
		if mt.Debug {
			log.Printf("obj after processing options -> %+v", obj)
		}
	}
	// Merge in path values into obj
	if len(mt.Vars) > 0 {
		if mt.Debug {
			log.Printf("varnames -> %+v\n", mt.Vars)
		}
		for _, varname := range mt.Vars {
			val := r.PathValue(varname)
			if val != "" {
				// val presidence over mt.Options
				obj[varname] = val
			}
		}
		if mt.Debug {
			log.Printf("obj after processing varnames -> %+v", obj)
		}
	}
	if obj == nil {
		if mt.Debug {
			log.Printf("no data attribute defined for template processing")
		}
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	// Handle case where some how the service was started before setting up template processing
	if mt.Tmpl == nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	if mt.Debug {
		log.Printf("mt.Tmpl -> %+v", mt.Tmpl)
	}
	mt.Tmpl.FRender(w, obj)
}

func (nm *NewtMustache) ListenAndServe() error {
	mux := http.NewServeMux()
	for _, mt := range nm.Templates {
		mux.HandleFunc(mt.Pattern, func(w http.ResponseWriter, r *http.Request) {
			if mt.Debug {
				log.Printf("mux.HandleFunc(%q, mt.Handler)", mt.Pattern)
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
