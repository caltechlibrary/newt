package newt

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	//"net/url"
	"os"
	"strings"
)

const (
	reqPath = iota
	reqMethod
	reqContentType
	apiURL
	apiMethod
	apiContentType
	pandoc
	pandocOptions
	pandocTemplate
	resHeaders
)

var (
	routerColumns = []string{
		"req_path",
		"req_method",
		"req_content_type",
		"api_url",
		"api_method",
		"api_content_type",
		"pandoc",
		"pandoc_options",
		"pandoc_template",
		"res_headers",
	}
)

// Route describes an individual route, maps to the columns of
// a route CSV file.
type Route struct {
	// ReqPath, a path described by RouteDSL from a browser or front end web server
	ReqPath        *RouteDSL
	// ReqMethod, the request HTTP method (e.g. GET, POST, PUT, DELETE, PATCH, HEAD)
	ReqMethod      string
	// ReqContentType, content of request, e.g. text/html, application/json
	ReqContentType string
	// ApiURL is the URL used to contact the data source (e.g. PostgREST, Solr, Elasticsearch)
	ApiURL         string
	// ApiMethod indicates the HTTP method to be used to contact the data source api
	ApiMethod      string
	// ApiContentType is the content type to be used to contact the data source api
	ApiContentType string
	// Pandoc if true route the data source content through Pandoc server
	Pandoc         bool
	// PandocOptions is a URL query string to pass to the Pandoc server
	PandocOptions  string
	// PandocTemplate holds the source to a Pandoc template
	PandocTemplate []byte
	// ResHeaders holds any additional response headers to send back to the 
	// browser or front end web server
	ResHeaders     map[string]string
}

func (route *Route) String() string {
	src, _ := json.MarshalIndent(route, "", "    ")
	return fmt.Sprintf("%s", src)
}


type Router struct {
	Env    map[string]string
	Routes []*Route
}

// Setenv adds an explicit environment name and value to Router.Env
func (router *Router) Setenv(key string, val string) {
	if router.Env == nil {
		router.Env = map[string]string{}
	}
	router.Env[key] = val
}

// Getenv retrieves an environment value using a name
func (router *Router) Getenv(key string) string {
	if val, ok := router.Env[key]; ok {
		return val
	}
	return ""
}


// ReadCSV filename
func (router *Router) ReadCSV(fName string) error {
	in, err := os.Open(fName) 
	if err != nil {
		return err
	}
	defer in.Close()
	r := csv.NewReader(in)
	rowNo := 1
	for {
		record, err := r.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("%q line %d: %s", fName, rowNo, err)
		}
		if rowNo == 1 {
			if len(record) < len(routerColumns) {
				return fmt.Errorf("missing columns %q, line %d", fName, rowNo)
			}
			for i, k := range routerColumns {
				if record[i] != k {
					return fmt.Errorf("invalid column name (%d), expected %q got %q", i, k, record[i])
				}
			}
		} else {
			if len(record) < len(routerColumns) {
				return fmt.Errorf("missing columns for routes line %d in %q", rowNo, fName)
			} 
			var err error
			route := new(Route)
			route.ReqPath, err = NewRouteDSL(record[reqPath])
			if err != nil {
				return fmt.Errorf("req_path error, line %d in %q, %s", rowNo, fName, err)
			}
			if ! isHTTPMethod(record[reqMethod]) {
				return fmt.Errorf("rep_method error, line %d in %q, unsupported method %q", rowNo, fName, record[reqMethod])
			}
			route.ReqMethod  = record[reqMethod]
			// FIXME: validate content type
			route.ReqContentType = record[reqContentType]
			// FIXME: need to make sure this make sense
			route.ApiURL = record[apiURL]
			if ! isHTTPMethod(record[reqMethod]) {
				return fmt.Errorf("api_method error, line %d in %q, unsupported method %q", rowNo, fName, record[reqMethod])
			}
			route.ApiMethod = record[apiMethod]
			// FIXME: validate content type
			route.ApiContentType = record[apiContentType]
			// FIXME: make sure this is a port number
			if (strings.ToLower(record[pandoc]) == "true") || (record[pandoc] == "1") {
				route.Pandoc = true
			} else {
				route.Pandoc = false
			}
			// FIXME: make sure this is formed as a query string
			route.PandocOptions = record[pandocOptions]
			// FIXME: validate template name, save template source
			templateName := record[pandocTemplate]
			if templateName != "" {
				route.PandocTemplate, err = os.ReadFile(templateName)
				if err != nil {
					return fmt.Errorf("failed to load template %q, line %d in %q, %s",  templateName, rowNo, fName, err)
				}
			}
			if record[resHeaders] != "" {
				headers := map[string]string{}
				if err := json.Unmarshal([]byte(record[resHeaders]), &headers); err != nil {
					return fmt.Errorf("failed to parse JSON expression for headers, line %d in %q, %s", rowNo, fName, err)
				}
				route.ResHeaders = headers
			}
			//fmt.Fprintf(os.Stderr, "DEBUG adding route %s\n", route.String())
			router.Routes = append(router.Routes, route)
		}
		rowNo++
	}
	return nil
}

// ResolveApiURL takes an in bound Request URL, matches it to a route
// and returns a resolved JSON data API URL based the related api_url.
// It included an error if something went wrong.
func (router *Router) ResolveApiURL(no int, m map[string]string, u string) (string, bool) {
	if no >= 0 && no < len(router.Routes) {
		return router.Routes[no].ReqPath.Resolve(m, router.Routes[no].ApiURL), true
	}
	return "", false
}

func (router *Router) ResolveRoute(u string, method string, contentType string) (int, map[string]string, bool) {
	for i, r := range router.Routes {
		// Compare our HTTP method and Content Type specrified
		if router.Routes[i].ReqMethod == method && 
				router.Routes[i].ReqContentType == contentType {
			// Now see if the path matches
			if m, ok := r.ReqPath.Eval(u); ok {
				// Merge the environent map of router.Env with the 
				// returned map from Eval.
				//
				// NOTE: environment overwrites the returned map!
				// This is deliberate.
				for k, v := range router.Env {
					m[k] = v
				}
				return i, m, true
			}
		}
	}
	return -1, nil, false
}


// Handler implements the http handler used for Routing requests.
func (router *Router) Handler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Handler not implemented yet\n\n")
}

func isHTTPMethod(s string) bool {
	s = strings.ToUpper(s)
	for _, method := range []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD"} {
		if strings.Compare(method, s) == 0 {
			return true
		}
	}
	return false
}


