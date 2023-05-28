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
	ReqPath *RouteDSL
	// ReqMethod, the request HTTP method (e.g. GET, POST, PUT, DELETE, PATCH, HEAD)
	ReqMethod string
	// ReqContentType, content of request, e.g. text/html, application/json
	ReqContentType string
	// ApiURL is the URL used to contact the data source (e.g. PostgREST, Solr, Elasticsearch)
	ApiURL string
	// ApiMethod indicates the HTTP method to be used to contact the data source api
	ApiMethod string
	// ApiContentType is the content type to be used to contact the data source api
	ApiContentType string
	// Pandoc if true route the data source content through Pandoc server
	Pandoc bool
	// PandocOptions is a URL query string to pass to the Pandoc server
	PandocOptions string
	// PandocTemplate holds the source to a Pandoc template
	PandocTemplate []byte
	// ResHeaders holds any additional response headers to send back to the
	// browser or front end web server
	ResHeaders map[string]string
}

func (route *Route) String() string {
	src, _ := json.MarshalIndent(route, "", "    ")
	return fmt.Sprintf("%s", src)
}

type Router struct {
	Debug bool
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

// OverlayEnv merges the env with m (a map[string]string).
// For keys in common the value in env replaces that in m.
func (router *Router) OverlayEnv(m map[string]string) map[string]string {
	for k, v := range router.Env {
		m[k] = v
	}
	return m
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
			if !isHTTPMethod(record[reqMethod]) {
				return fmt.Errorf("rep_method error, line %d in %q, unsupported method %q", rowNo, fName, record[reqMethod])
			}
			route.ReqMethod = record[reqMethod]
			// FIXME: validate content type
			route.ReqContentType = record[reqContentType]
			// FIXME: need to make sure this make sense
			route.ApiURL = record[apiURL]
			if !isHTTPMethod(record[reqMethod]) {
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
					return fmt.Errorf("failed to load template %q, line %d in %q, %s", templateName, rowNo, fName, err)
				}
			}
			if record[resHeaders] != "" {
				headers := map[string]string{}
				if err := json.Unmarshal([]byte(record[resHeaders]), &headers); err != nil {
					return fmt.Errorf("failed to parse JSON expression for headers, line %d in %q, %s", rowNo, fName, err)
				}
				route.ResHeaders = headers
			}
			if router.Debug {
				fmt.Fprintf(os.Stderr, "DEBUG adding %s\n", route.String())
			}
			router.Routes = append(router.Routes, route)
		}
		rowNo++
	}
	return nil
}

func (route *Route) HasContentType(contentTypes string) bool {
	if strings.Contains(contentTypes, route.ReqContentType) {
		return true
	}
	return false
}

func (route *Route) HasReqMethod(method string) bool {
	if route.ReqMethod == method {
		return true
	}
	return false
}

func (router *Router) ResolveRoute(u string, method string, contentType string) (int, map[string]string, bool) {
	for i, r := range router.Routes {
		// Compare our HTTP method and make sure Content Types
		// are supported by request route.
		if router.Routes[i].HasReqMethod(method) && router.Routes[i].HasContentType(contentType) {
			// Now see if the path matches
			if m, ok := r.ReqPath.Eval(u); ok {
				// Merge the environent map of router.Env with the
				// returned map from Eval.
				//
				// NOTE: environment overwrites the returned map!
				// This is deliberate.
				m = router.OverlayEnv(m)
				return i, m, true
			}
		}
	}
	return -1, nil, false
}

// ResolveApiURL takes an in bound Request URL, matches it to a route
// and returns a resolved JSON data API URL based the related api_url.
// It included an error if something went wrong.
func (router *Router) ResolveApiURL(no int, m map[string]string) (string, bool) {
	if no >= 0 && no < len(router.Routes) {
		res, ok := router.Routes[no].ReqPath.Resolve(m, router.Routes[no].ApiURL), true
		if router.Debug {
			fmt.Fprintf(os.Stderr, "DEBUG resolved api URL %q\n", res)
		}
		return res, ok
	}
	return "", false
}

// RequestDataAPI
func (router *Router) RequestDataAPI(rNo int, apiURL string) ([]byte, string, int) {
	// Make an http call to the JSON data API
	return nil, http.StatusText(501), 501
}

// RequestPandoc
func (router *Router) RequestPandoc(rNo int, src []byte) ([]byte, string, int) {
	// Make an http call to Pandoc server
	return nil, http.StatusText(501), 501
}

// WriteResponse
func (router *Router) WriteResponse(w http.ResponseWriter, rNo int, src []byte) {
	for k, v := range router.Routes[rNo].ResHeaders {
		w.Header().Set(k, v)
	}
	w.WriteHeader(http.StatusOK)
	// Write back the content
	if router.Debug {
		fmt.Fprintf(os.Stderr, "DEBUG response content\n%s\n", src)
	}
	fmt.Fprintf(w, "%s", src)
}

// Newt implements the router as an http middleware. This allows our
// router to be used more generally and gives us options like having
// the Newt application be both a data router and static file server.
func (router *Router) Newt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		contentTypes := req.Header.Get("Content-Type")
		rNo, m, ok := router.ResolveRoute(req.URL.Path, req.Method, contentTypes)
		if !ok {
			/*
				// Handle 404
				http.Error(w, http.StatusText(404), 404)
				return
			*/
			// We don't match a route, pass this on to the next handler
			next.ServeHTTP(w, req)
			return
		}
		// Get content from the JSON Data API
		apiURL, ok := router.ResolveApiURL(rNo, m)
		if !ok {
			// Handle 502, Bad Gateway
			http.Error(w, http.StatusText(502), 502)
			return
		}
		src, statusText, statusCode := router.RequestDataAPI(rNo, apiURL)
		if statusCode < 200 || statusCode >= 300 {
			// echo back the data request status code and text
			http.Error(w, statusText, statusCode)
			return
		}
		// NOTE: if Pandoc transform data request
		if router.Routes[rNo].Pandoc {
			src, statusText, statusCode = router.RequestPandoc(rNo, src)
			if statusCode < 200 || statusCode >= 300 {
				// echo back the pandoc request status code and text
				http.Error(w, statusText, statusCode)
				return
			}
		}
		// NOTE: We have content from either JSON data API or what
		// we sent through Pandoc, write out.
		router.WriteResponse(w, rNo, src)
		return
	})
}
