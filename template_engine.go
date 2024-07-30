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
	//"time"

	// 3rd Party Templates
	//"github.com/cbroglie/mustache"
	"github.com/aymerick/raymond"
	//"gopkg.in/yaml.v3"
)


// ResolvePath reviews the `.Request` attribute and updates the Vars using PatternKeys()
func (t *Template) ResolvePath() error {
	// Does the `.Request` hold a pattern or a fixed string?
	if strings.Contains(t.Pattern, "{") {
		if !strings.Contains(t.Pattern, "}") {
			return fmt.Errorf("%q is malformed", t.Pattern)
		}
		// Record our list of var names so handler can override the object being constructed from a path.
		t.Vars = PatternKeys(t.Pattern)
	}
	if t.Debug {
		log.Printf("assigning t.Pattern -> %q\n", t.Pattern)
		log.Printf("vars -> %+v\n", t.Vars)
	}
	return nil
}

// ResolvesTemplate is responsible for reading and parse the template and partials associated with a mapped request.
// If an error is encountered a error value is returned.
func (t *Template) ResolveTemplate() error {
	if t.Template != "" {
		if t.Debug {
			log.Printf("attempting to parse single template %q", t.Template)
		}
		//FIXME: Need to use MustParse and also handle regestering partials, etc.
		src, err := os.ReadFile(path.Join(t.BaseDir, t.Template + t.ExtName))
		if err != nil {
			return err
		}
		tmpl := raymond.MustParse(fmt.Sprintf("%s", src))
		//FIXME: Need to attach the resulting template to the template object
		//NOTE: Load partials from `{base_dir}/{partials_dir}`
		if t.PartialsDir != ""{
			if t.Debug {
				log.Printf("handling primary and partial templates")
			}
			//FIXME: Read in the directory of partial templates and register them.
			//partialsDir := path.Join(t.BaseDir, t.PartialsDir)
			//FIXME: read the partial files found in the partial dir and Register them ...
			log.Printf("DEBUG reading and registering partials not implemented yet")
		}
		//FIXME: Handle Layouts for the template
		if t.LayoutsDir != "" {
			log.Printf("DEBUG reading and handling layouts not implemented yet")
		}
		if t.DefaultLayout != "" {
			log.Printf("DEBUG reading and handling default layouts not implemented yet")
		}
		if t.Helpers != nil && len(t.Helpers) > 0 {
			log.Printf("DEBUG reading and handling helpers not implemented yet")
		}
		if t.CompilerOptions != nil && len(t.CompilerOptions) > 0 {
			log.Printf("DEBUG reading and handling compiler options not implemented yet")
		}

		// Now that we've read in all the template parts we can assign the Handlebars template
		// for use in the web service. 
		t.Tmpl = tmpl
		if t.Debug {
			log.Printf("templates parsed successfully")
		}
		return nil
	}
	return fmt.Errorf("no template found")
}

// Handler decodes a the request body and then processes that as a template engine.
func (t *Template) Handler(w http.ResponseWriter, r *http.Request) {
	if t.Debug {
		log.Printf(".Handler(w, %s %s)", r.Method, r.URL.Path)
	}

	src, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if t.Debug {
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
		if t.Debug {
			log.Printf(".body -> %s", src)
		}
		dec := json.NewDecoder(bytes.NewBuffer(src))
		dec.UseNumber()
		if err := dec.Decode(&body); err != nil && err != io.EOF {
			if t.Debug {
				log.Printf("failed to decode JSON Response body, %s", err)
			}
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}
		if body == nil {
			if t.Debug {
				log.Printf("no data for template processing")
			}
			http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
			return
		}
	}
	document := map[string]string{}
	vars := map[string]string{}
	// Copy in the options into page objcet's options.
	if t.Document != nil {
		if t.Debug {
			log.Printf("Document -> %+v\n", t.Document)
		}
		for k, v := range t.Document {
			document[k] = v
		}
		if t.Debug {
			log.Printf("obj after processing Document -> %+v", document)
		}
	}
	// Merge in path values into .vars
	if len(t.Vars) > 0 {
		if t.Debug {
			log.Printf("varnames -> %+v\n", t.Vars)
		}
		for _, varname := range t.Vars {
			val := r.PathValue(varname)
			if val != "" {
				// val presidence over t.Options
				vars[varname] = val
			}
		}
		if t.Debug {
			log.Printf("obj after processing varnames -> %+v", vars)
		}
	}
	obj := map[string]interface{}{
		"body":       body,
		"document":   document,
		"vars":       vars,
	}
	if t.Debug {
		log.Printf("obj after processing options -> %+v", obj)
	}
	// Handle case where some how the service was started before setting up template processing
	if t.Tmpl == nil {
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	if t.Debug {
		log.Printf("t.Tmpl -> %+v", t.Tmpl)
	}
	// We want to write to a buffer so we can do content detection and set the headers correctly.
	txt, err := t.Tmpl.Exec(obj)
	if err != nil {
		log.Printf("t.Tmpl.Exec(obj) error -> %s: %s", t.Template, err)
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	src = []byte(txt)
	contentType := http.DetectContentType(src)
	if strings.HasPrefix(strings.ToLower(txt), "<!doc") || strings.HasPrefix(strings.ToLower(txt), "<html") {
		contentType = "text/html; charset=utf-8"
	}
	if t.Debug {
		log.Printf("content type: %q -> %q", contentType, src)
	}
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.Write(src)
}

func (te *TemplateEngine) ListenAndServe() error {
	mux := http.NewServeMux()
	// Setup our handlers, POST for process data with the template and GET to retreive the template
	// ast.
	for _, t := range te.Templates {
		//FIXME: Need to map in the vocabularies into .Document 
		//FIXME: Need to map in app options into Document object passed to template
		/*
		if len(te.Options) > 0 {
			t.Options = map[string]string{}
			for k, v := range te.Options {
				t.Options[k] = v
			}
		}
		*/
		// FIXME: Process the data displaying template if a GET is used.
		// Process the data with template if a POST.
		mux.HandleFunc("POST "+t.Pattern, func(w http.ResponseWriter, r *http.Request) {
			if t.Debug {
				log.Printf("mux.HandleFunc(%q, t.Handler)", "POST "+t.Pattern)
				log.Printf(".vars -> %+v", t.Vars)
			}
			t.Handler(w, r)
		})
	}
	// Now create my http server
	svr := new(http.Server)
	svr.Addr = fmt.Sprintf(":%d", te.Port)
	svr.Handler = mux
	if err := svr.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	return nil
}
