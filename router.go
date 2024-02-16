package newt

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// 3rd Party packages
	"gopkg.in/yaml.v3"
)

// NewtRouter is used to implement the Newt Router
type NewtRouter struct {
	// Port is the port the router will listen on
	Port string `json:"port,omitempty" yaml:"port"`

	// Routes holds a list of route
	Routes []*NewtRoute `json:"routes,omitempty" yaml:"routes"`

	// Htdocs holds the location of a the static files if used
	Htdocs string `json:"htdocs,omitempty" yaml:"htdocs"`

	// Environment holds those environment values propogates at startup.
	// these are mapped into each route's "options"
	Environment []string `json:"environment,omitempty" yaml:"environment"`
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
func (s *Service) MakeRequest(env map[string]string, input []byte, debug bool) ([]byte, error) {
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
		return nil, err
	}
	client := http.Client{
		Timeout: timeout,
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	statusText := http.StatusText(response.StatusCode)
	if response.StatusCode < http.StatusOK || response.StatusCode > http.StatusNoContent {
		err := fmt.Errorf("request failed for %s %s, %d %s", method, uri, response.StatusCode, statusText)
		if debug {
			log.Printf("%s", err)
		}
		return nil, err
	}
	data, err := io.ReadAll(response.Body)
	if err != nil {
		if debug {
			log.Printf("failed to read response body %s %s %d %s, %s", method, uri, response.StatusCode, statusText, err)
		}
		return nil, err
	}
	if debug {
		contentType := http.DetectContentType(data)
		l := len(data)
		log.Printf("%s %s returning %s, %d byte(s)", method, uri, contentType, l)
	}
	return data, nil
}

func (nr *NewtRoute) RunPipeline(w http.ResponseWriter, r *http.Request, env map[string]string) ([]byte, error) {
	var (
		input  []byte
		output []byte
		err    error
	)
	if nr.Debug && len(env) > 0 {
		log.Printf("pipeline environment %+v", env)
	}
	for i, service := range nr.Pipeline {
		output, err = service.MakeRequest(env, input, nr.Debug)
		if err != nil {
			if nr.Debug {
				log.Printf("service %d, %s error %s", i, service.Service, err)
			}
			// We stop processing there was a problem.
			return nil, err
		}
		input = output
	}
	return output, nil
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
	data, err := nr.RunPipeline(w, r, env)
	if err != nil {
		// FIXME: should return Error page
		if nr.Debug {
			log.Printf("pipeline failed for route %s, %s", nr.Pattern, err)
		}
		// Need to decide which http status is right
		// 500 Internel Server Error
		// 501 Not implemented
		// 503 Service Unavailable
		// 504 Gateway Timeout
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// NewNewtRouter creates a newt router suprisingly
func NewNewtRouter(fName string, cfg *Config) (*NewtRouter, error) {
	rtr := &NewtRouter{}
	if fName == "" {
		return rtr, fmt.Errorf("missing configuration file")
	}
	src, err := os.ReadFile(fName)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(src, &rtr)
	if err != nil {
		return nil, err
	}
	env := map[string]string{}
	if cfg.Application != nil {
		if rtr.Port == "" && cfg.Application.Port != 0 {
			rtr.Port = fmt.Sprintf(":%d", cfg.Application.Port)
		}
		// Populate an env from our allowed environment variables
		if len(cfg.Application.Environment) > 0 {
			for _, envar := range cfg.Application.Environment {
				val := os.Getenv(envar)
				if val != "" {
					env[envar] = val
				}
			}
		}
		if len(env) > 0 {
			for _, route := range rtr.Routes {
				// Copy in the environment values without overwriting defined options.
				for k, v := range env {
					if _, conflict := route.Options[k]; !conflict {
						route.Options[k] = v
					}
				}
			}
		}
	}
	// Prefix the port number with a colon
	if !strings.HasPrefix(rtr.Port, ":") {
		rtr.Port = fmt.Sprintf(":%s", rtr.Port)
	}
	for _, route := range rtr.Routes {
		// Map the environment into options if not set
		for k, v := range env {
			if _, conflict := route.Options[k]; !conflict {
				route.Options[k] = v
			}
		}
	}
	return rtr, nil
}

// ListenAndServe() runs the router web service
func (rtr *NewtRouter) ListenAndServe() error {
	mux := http.NewServeMux()
	for _, nr := range rtr.Routes {
		mux.HandleFunc(nr.Pattern, func(w http.ResponseWriter, r *http.Request) {
			if nr.Debug {
				log.Printf("mux.HandleFunc(%q, nr.Handler)", nr.Pattern)
				log.Printf(".vars -> %+v", nr.Vars)
			}
			nr.Handler(w, r)
		})
	}
	// Now create my http server
	svr := new(http.Server)
	svr.Addr = rtr.Port
	svr.Handler = mux
	if err := svr.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
	return nil
}
