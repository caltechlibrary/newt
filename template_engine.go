package newt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
	//"path"
	"strings"
	//"time"

	// 3rd Party Templates
	"github.com/cbroglie/mustache"
	//"gopkg.in/yaml.v3"
)


// ResolvePath reviews the `.Request` attribute and updates the Vars using PatternKeys()
func (mt *Template) ResolvePath() error {
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
func (mt *Template) ResolveTemplate() error {
	if mt.Template != "" {
		//FIXME: Load partials from `{base_dir}/{partials_dir}`
		/*
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
		*/
		if mt.Debug {
			log.Printf("attempting to parse single template %q", mt.Template)
		}
		tmpl, err := mustache.ParseFile(mt.Template)
		if err != nil {
			return err
		}
		//FIXME: Need to attach the resulting template to the template object
		mt.Tmpl = tmpl
		if mt.Debug {
			log.Printf("templates parsed successfully")
		}
		return nil
	}
	return fmt.Errorf("no template found")
}

// Handler decodes a the request body and then processes that as a template engine.
func (mt *Template) Handler(w http.ResponseWriter, r *http.Request) {
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
		if mt.Debug {
			log.Printf(".body -> %s", src)
		}
		dec := json.NewDecoder(bytes.NewBuffer(src))
		dec.UseNumber()
		if err := dec.Decode(&body); err != nil && err != io.EOF {
			if mt.Debug {
				log.Printf("failed to decode JSON Response body, %s", err)
			}
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}
		if body == nil {
			if mt.Debug {
				log.Printf("no data for template processing")
			}
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}
	}
	document := map[string]string{}
	vars := map[string]string{}
	// Copy in the options into page objcet's options.
	if mt.Document != nil {
		if mt.Debug {
			log.Printf("Document -> %+v\n", mt.Document)
		}
		for k, v := range mt.Document {
			document[k] = v
		}
		if mt.Debug {
			log.Printf("obj after processing Document -> %+v", document)
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
		"body":       body,
		"options":    options,
		"vocabulary": vocabulary,
		"vars":       vars,
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

func (nm *TemplateEngine) ListenAndServe() error {
	mux := http.NewServeMux()
	// Setup our handlers, POST for process data with the template and GET to retreive the template
	// ast.
	for _, mt := range nm.Templates {
		//FIXME: Need to map in the vocabularies into .Document 
		//FIXME: Need to map in app options into Document object passed to template
		/*
		if len(nm.Options) > 0 {
			mt.Options = map[string]string{}
			for k, v := range nm.Options {
				mt.Options[k] = v
			}
		}
		*/
		// FIXME: Process the data displaying template if a GET is used.
		// Process the data with template if a POST.
		mux.HandleFunc("POST "+mt.Pattern, func(w http.ResponseWriter, r *http.Request) {
			if mt.Debug {
				log.Printf("mux.HandleFunc(%q, mt.Handler)", "POST "+mt.Pattern)
				log.Printf(".vars -> %+v", mt.Vars)
			}
			mt.Handler(w, r)
		})
	}
	// Now create my http server
	svr := new(http.Server)
	svr.Addr = fmt.Sprintf(":%d", nm.Port)
	svr.Handler = mux
	if err := svr.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	return nil
}
