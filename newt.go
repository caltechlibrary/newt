/**
 * newt.go an implementation of the Newt URL router.
 *
 * @author R. S. Doiel
 */
package newt

import (
	"text/scanner"
)

// RouteFunc implements a "type" rule that returns true of the input
// string contains an element defined by the func and the string
// value extracted by the function. If the "type" rule is encounters
// and error then it returns "false" and an empty string value.
type RouteFunc func(string) (bool, string)

// A route holds a definition and an pipeline of route functions
// used to extract the route's parts used to validate a route match
// and return any extracted values.
type Route struct {
	// This holds the source string defining the route and any
	// embedded variable definitions.
	Defn string
	// Pipeline is an array of Route functions used to parse
	// a input string. The route functions are applied in order
	// they are added to the router with `AddRoute()`.
	// If any function returns false the pipeline processing is aborted
	// The pipeline functions are used by Route's `Eval(string)` which
	// returns a success booleana and a map to any defined variable
	// names and their values.
	Pipeline []*RouteFunc
}

// Eval takes an input string and returns true if the string
// matches the route's definition and a map of variable names
// and values if found. If the route's defn does not match it
// returns false and a nil map.
func (r *Route) Eval(string) (bool, map[string]string) {
}

type Router struct {
	Routes []*Route
}

// NewRouter returns a new Router where you can add 
// route definitions.
func NewRouter() *Router {
	return new(Router)
}

// AddRoute takes a route definition and adds it to the router
// or returns an error if the definition fails to parse.
func (router *Router) AddRoute(defn string) error {
	r := new(Route)
	r.Defn = defn
	// FIXME: Parse the defn into a pipe line of RouteFuncs
	return fmt.Errorf("r.AddRoute() not implemented")
}


// GetStringValue evaluates a string and returns a success value and 
// string value. If "l" is negative then true and the whole string.
// If "l" is greater than -1 then and the string is of that length then
// true and the string is returned.  If the string is longer than "l"
// then false and a string of length "l" is returned.
// 
func GetStringValue(s string, l int) (bool, string) {
	if l < 0 || len(s) == l {
		return true, s[:]
	}
	return false, ""
}

// GetLiteral returns the longest prefixed literal string in s. If no
// literal prefix found, stops
// and return a 
func GetLiteral(s string) (bool, string) {
	pos := strings.Index(s, "{")
	if pos > -1 {
		return true, s[0:pos]
	}
	return false, ""
}

// GetVarDefn returns the next variable definition found
func GetVarDefn(s string) (bool, string) {
	start := strings.Index(s, "{")
	end := strings.Index(s, "}")
	if (start > -1) &&  (end > -1) {
		return true, s[start:end]
	}
	return false, ""
}

// EvalRoute takes a list of routes description and a path value and returns
// a boolean match and if match is true a map of names and values 
//
// ```
// routeDef := `/blog/{year year}/{month month}/{day day}/`
// reqPath := `/blog/2023/05/11/`
// if ok, m := newt.HasRoute(routeDef, reqPath); ok {
//    // ... using m assemble a PostgREST request, process with Pandoc ...
// } else {
//    // ... handle 404, requested route is not defined ...
// }
// ```
func EvalRoute(routes []*Routes, reqPath string) (bool, map[string]string) {
	// Parse route returning literals or variable declarations
	for _, routes

	// FIXME: HasRoute not implemented
	return false, nil
}
