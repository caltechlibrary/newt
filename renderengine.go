package newt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	// 3rd Party Packages
	"github.com/cbroglie/mustache"
)

// RenderRequest holds a JSON encoded Mustache server request.
type RenderRequest struct {
	// This is the source text of the template.
	Template string `json:"template,omitempty"`
	// Data is a map of variable names and values to be processed with the template
	Data map[string]interface{} `json:"data,omitempty"`
	// ContentType sets the return content-type of the response. If omitted defaults to text/plain.
	ContentType string `json:"content_type,omitempty"`
}

// MustacheHandleRequest takes a POST, decodes a RenderRequest then processes that
// as a Mustache template.
func MustacheHandleRequest(w http.ResponseWriter, r *http.Request, verbose bool, out io.Writer) {
	// Make sure we have a post
	if r.Method != http.MethodPost {
		if verbose {
			fmt.Fprintf(out, "http method %q not supported\n", r.Method)
		}
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	src, err := ioutil.ReadAll(r.Body)
	if err != nil {
		if verbose {
			fmt.Fprintf(out, "failed to read request body, %s\n", err)
		}
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	} 
	renderReq := &RenderRequest{}
	dec := json.NewDecoder(bytes.NewBuffer(src))
	dec.UseNumber()
	if err := dec.Decode(&renderReq); err != nil  && err != io.EOF {
		if verbose {
			fmt.Fprintf(out, "failed to decode JSON Response body, %s\n", err)
		}
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	template := renderReq.Template
	if template == "" {
		if verbose {
			fmt.Fprintf(out, "failed to decode template in JSON Response, %s\n", err)
		}
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	tmpl, err := mustache.ParseString(template)
	if err != nil {
		if verbose {
			fmt.Fprintf(out, "template failed to parse %s\n", err)
		}
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	if renderReq.Data == nil {
		if verbose {
			fmt.Fprintf(out, "no data attribute defined for template processing\n")
		}
		http.Error(w, http.StatusText(http.StatusUnprocessableEntity), http.StatusUnprocessableEntity)
		return
	}
	if renderReq.ContentType != "" {
		w.Header().Set("Content-Type", renderReq.ContentType)
	}
	tmpl.FRender(w, renderReq.Data)
}

// MustacheServer provides a Pandoc server like experience for Mustache templates.
func MustacheServer(out io.Writer, eout io.Writer, port string, timeout int, verbose bool) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		MustacheHandleRequest(w, r, verbose, out)
	})
	if ! strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	if verbose {
		fmt.Fprintf(out, "starting server on %s, timeout set to %d\n", port, timeout)
	}
	srv := &http.Server{
		Addr: port,
		ReadTimeout: time.Duration(timeout) * time.Second,
	}
	err := srv.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Fprintf(out, "server closed\n")
	} else if err != nil {
		fmt.Fprintf(eout, "error starting server: %s\n", err)
	}
	return err
}
