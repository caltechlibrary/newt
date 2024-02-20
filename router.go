package newt

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// NewtRouter is used to implement the Newt Router
type NewtRouter struct {
	// Port is the port the router will listen on
	Port string 

	// Routes holds a list of route
	Routes []*NewtRoute

	// Htdocs holds the location of a the static files if used
	Htdocs string
}

// This holds the route definitions, e.g. request, description, pipeline, debug
type NewtRoute struct {
	// Pattern holds the HTTP Method and URL path, may include Go 1.22 patterns
	Pattern string `json:"request,required" yaml:"request"`

	// Description holds a human describe of the purpose of this route.
	Description string `json:"description,omitempty" yaml:"description"`

	// Pipeline holds the series of http services context with the output of one sent to another.
	Pipeline []*Service `json:"pipeline,omitempty" yaml:"pipeline"`

	// Debug if true log verbosely
	Debug bool `json:"debug,omitempty" yaml:"debug"`

	// Env holds a map of defaults that are available from the environment and from path values in the url
	Options map[string]string `json:"options,omitempty" yaml:"options"`

	// Vars holds the variables defined in the route
	Vars []string
}

// Service holds the necessary information to contact a data source and retrieve the results for use in a pipeline.
type Service struct {
	// Service holds the http Request Pattern to request a resource from a service
	Service string `json:"service,required" yaml:"service"`

	// Description describes the service and purpose of contact. Human readable.
	Description string `json:"description,omitempty" yaml:"description"`

	// Timeout sets a timeout value to recieve a response from the service.
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout"`
}

// ResolveRoute reviews the `.Request` attribute and updates the Vars using PatternKeys()
func (nr *NewtRoute) ResolveRoute() error {
	// Does the `.Request` hold a pattern or a fixed string?
	if strings.Contains(nr.Pattern, "{") {
		if !strings.Contains(nr.Pattern, "}") {
			return fmt.Errorf("%q is malformed", nr.Pattern)
		}
		// Record our list of var names so handler can override the object being constructed from a path.
		nr.Vars = PatternKeys(nr.Pattern)
	}
	if nr.Debug {
		log.Printf("assigning mb.Pattern -> %q\n", nr.Pattern)
		log.Printf("vars -> %+v\n", nr.Vars)
	}
	return nil
}

// Request a service, sending any import if provided.
// Returns a byte splice of results and an error value
func (s *Service) MakeRequest(env map[string]string, input []byte, debug bool) ([]byte, int, string, error) {
	// The default method for a service is POST, it can be overwritten by what is supplied in the .Service's pattern.
	method := http.MethodPost
	uri := s.Service
	if strings.Contains(s.Service, " ") {
		parts := strings.SplitN(uri, " ", 2)
		method, uri = parts[0], parts[1]
	}
	if len(env) > 0 && strings.Contains(uri, "{") && strings.Contains(uri, "}") {
		for k, v := range env {
			varhandle := fmt.Sprintf("{%s}", k)
			if strings.Contains(uri, varhandle) {
				uri = strings.ReplaceAll(uri, varhandle, v)
			}
		}
	}
	if debug {
		log.Printf("making a %s request to %s", method, uri)
	}
	var timeout time.Duration
	if s.Timeout > 0 {
		timeout = s.Timeout * time.Second
	} else {
		timeout = 30 * time.Second
	}
	// We should now have enough information to make our request.
	req, err := http.NewRequest(method, uri, bytes.NewReader(input))
	if err != nil {
		return nil, http.StatusInternalServerError, "", err
	}
	client := http.Client{
		Timeout: timeout,
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, response.StatusCode, "", err
	}
	statusText := http.StatusText(response.StatusCode)
	if response.StatusCode < http.StatusOK || response.StatusCode > http.StatusNoContent {
		err := fmt.Errorf("request failed for %s %s, %d %s", method, uri, response.StatusCode, statusText)
		if debug {
			log.Printf("%s", err)
		}
		return nil, response.StatusCode, "", err
	}
	contentType := ""
	if response.Header != nil {
		contentType = response.Header.Get("Content-Type")
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		if debug {
			log.Printf("failed to read response body %s %s %d %s, %s", method, uri, response.StatusCode, statusText, err)
		}
		return nil, response.StatusCode, "", err
	}
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	if debug {
		l := len(data)
		log.Printf("%s %s returning %s, %d byte(s)", method, uri, contentType, l)
	}
	return data, response.StatusCode, contentType, nil
}

func (nr *NewtRoute) RunPipeline(w http.ResponseWriter, r *http.Request, env map[string]string) ([]byte, int, string, error) {
	var (
		input  []byte
		output []byte
		err    error
	)
	if nr.Debug && len(env) > 0 {
		log.Printf("pipeline environment %+v", env)
	}
	contentType := ""
	statusCode := http.StatusOK
	for i, service := range nr.Pipeline {
		output, statusCode, contentType, err = service.MakeRequest(env, input, nr.Debug)
		if err != nil {
			if nr.Debug {
				log.Printf("service %d, %s error %s", i, service.Service, err)
			}
			// We stop processing there was a problem.
			return nil, statusCode, contentType, err
		}
		input = output
	}
	return output, statusCode, contentType, nil
}

// ResolvePattern takes the request Pattern and pulls out the values from the actual request.
// returns values in map[string]string and an error value.
func (nr *NewtRoute) ResolvePattern(r *http.Request) map[string]string {
	// Copy in our options, then overwrite with any  vars
	m := map[string]string{}
	for k, v := range nr.Options {
		m[k] = v
	}
	// Now add our extracted Path Values
	for _, varname := range nr.Vars {
		val := r.PathValue(varname)
		if val != "" {
			m[varname] = val
		}
	}
	return m
}

// Handler creates an http handler func for a given route.
func (nr *NewtRoute) Handler(w http.ResponseWriter, r *http.Request) {
	// Resolve the enviroment with options and path values
	env := nr.ResolvePattern(r)

	// Build now we can run our pipeline and get back some data
	data, statusCode, contentType, err := nr.RunPipeline(w, r, env)
	if err != nil {
		// FIXME: should return Error page
		if nr.Debug {
			log.Printf("pipeline failed for route %s, %s", nr.Pattern, err)
		}
		// Bubble up any HTTP error codes
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	if nr.Debug {
		log.Printf("%s return content type: %q", nr.Pattern, contentType)
	}
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}
	w.Write(data)
}

// NewNewtRouter creates a newt router suprisingly
func NewNewtRouter(cfg *Config) (*NewtRouter, error) {
	router := &NewtRouter{
		Routes: cfg.Routes,
	}
	if cfg.Applications.NewtRouter.Htdocs != "" {
		router.Htdocs = cfg.Applications.NewtRouter.Htdocs
	}
	if cfg.Applications.NewtRouter.Port != 0 {
		router.Port = fmt.Sprintf(":%d", cfg.Applications.NewtRouter.Port)
	}
	// Populate an env from our allowed environment variables
	for _, route := range router.Routes {
		// Copy in the applications options without overwriting the route specific
		// defined options.  NOTE: Application options have already been resolved
		// with the environment by cfg.LoadConfig()
		for k, v := range cfg.Applications.Options {
			if _, conflict := route.Options[k]; ! conflict {
				route.Options[k] = v
			}
		}
	}
	// Prefix the port number with a colon
	if !strings.HasPrefix(router.Port, ":") {
		router.Port = fmt.Sprintf(":%s", router.Port)
	}
	return router, nil
}

// ListenAndServe() runs the router web service
func (rtr *NewtRouter) ListenAndServe() error {
	mux := http.NewServeMux()
	for _, nr := range rtr.Routes {
		// FIXME: need to warn if the patter is / when the htdocs 
		// directory is set.
		if (nr.Pattern == "/" || strings.HasSuffix(nr.Pattern, " /")) && rtr.Htdocs != "" {
			log.Println("WARNING: you have a htdocs directory set to service files and you've mapped a route to the name end point, /")
		}
		mux.HandleFunc(nr.Pattern, func(w http.ResponseWriter, r *http.Request) {
			if nr.Debug {
				log.Printf("mux.HandleFunc(%q, nr.Handler)", nr.Pattern)
				log.Printf(".vars -> %+v", nr.Vars)
			}
			nr.Handler(w, r)
		})
	}
	// Do we need to support htdocs static resources?
	if rtr.Htdocs != "" {
		fsys := dotFileHidingFileSystem{http.Dir(rtr.Htdocs)}
		mux.Handle("/", http.FileServer(fsys))
	}
	// Now create my http server
	svr := new(http.Server)
	svr.Addr = rtr.Port
	svr.Handler = NewLogger(mux)
	if err := svr.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	return nil
}
