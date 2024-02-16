package newt

import (
	"fmt"
)

// Router is used to implement the Newt Router
type Router struct {
	// Port is the port the router will listen on
	Port string `json:"port,omitempty" yaml:"port"`

	// Routes holds a list of route
	Routes []*Route `json:"routes,omitempty" yaml:"routes"`

	// Htdocs holds the location of a the static files if used
	Htdocs string `json:"htdocs,omitempty" yaml:"htdocs"`
}

// NewNewtRouter creates a newt router suprisingly
func NewRouter(fName string, cfg *Config) (*Router, error) {
	return nil, fmt.Errorf("NewRouter() not implemented")
}

// ListenAndServe() runs the router web service
func (rtr *Router) ListenAndServe() error {
	return fmt.Errorf("rtr.ListenAndServe() not implemented")
}

// This holds the route definitions, e.g. request, description, pipeline, debug
type Route struct {
}
