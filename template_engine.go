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
	"path/filepath"
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
		// NOTE: Need use MustParse and also handle regestering partials, etc.
		// This is done at startup and should fail if unsuccesful.
		tName := path.Join(t.BaseDir, t.Template + t.ExtName)
		src, err := os.ReadFile(tName)
		if err != nil {
			return err
		}
		tmpl := raymond.MustParse(fmt.Sprintf("%s", src))
		partials := map[string]string{}
		// FIXME: Need to attach the resulting template to the template object
		// NOTE: Load partials from `{base_dir}/{partials_dir}`
		if t.PartialsDir != ""{
			if t.Debug {
				log.Printf("handling primary and partial templates")
			}
			pattern := path.Join(t.BaseDir, t.PartialsDir, "*" + t.ExtName)
			names, err := filepath.Glob(pattern)
			if err != nil {
				log.Printf("failed to read %q, %s", pattern, err)
				return err
			}
			for _, fName := range names {
				// Get the basename without path and file extension, then register it as a partial template.
				name := path.Base(fName)
				if strings.HasSuffix(name, t.ExtName) {
					name = strings.TrimSuffix(name, t.ExtName)
				}
				if _, hasName := partials[name]; ! hasName {
					src, err := os.ReadFile(fName)
					if err != nil {
						return err
					}
					partials[name] = fmt.Sprintf("%s", src)
				}
			}
		}
		if len(partials) > 0 {
			tmpl.RegisterPartials(partials)
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
	document := make(map[string]interface{})
	vars := map[string]string{}
	// Copy in the options into page object's options.
	log.Printf("Debug Document -> %+v\n", t.Document)
	if t.Document != nil && len(t.Document) > 0 {
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
	log.Printf("Debug Vars -> %+v\n", t.Vars)
	if len(t.Vars) > 0 {
		if t.Debug {
			log.Printf("varnames -> %+v\n", t.Vars)
		}
		for _, varname := range t.Vars {
			val := r.PathValue(varname)
			if t.Debug {
				log.Printf("varname: %q -> %+v", varname, val)
			}
			if val != "" {
				// val presidence over t.Options
				vars[varname] = val
			}
		}
		if t.Debug {
			log.Printf("obj after processing varnames -> %+v", vars)
		}
	}
	log.Printf("Debug vars -> %+v\n", vars)
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
		// FIXME: Add the process to handle the GET request to display template source.
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
