package newt

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	//"net/url"
	"os"
	"strings"

	// 3rd Party
	"gopkg.in/yaml.v3"
)

const (
	varDef = iota
	reqPath
	reqMethod
	apiURL
	apiMethod
	apiContentType
	pandocTemplate
	resHeaders
)

var (
	routerColumns = []string{
		"var",
		"req_path",
		"req_method",
		"api_url",
		"api_method",
		"api_content_type",
		"pandoc_template",
		"res_headers",
	}
)

// Route describes an individual route, maps to the attributes in our
// YAML file for a specific your. 
type Route struct {
	// Var holds a map of variable names and types
	Var map[string]string `json:"var,omitempty" yaml:"var,omitempty"`

	// ReqPath, a path described by DSL from a browser or front end web server
	ReqPath *DSL `json:"req_path,omitempty" yaml:"req_path,omitempty"`

	// ReqMethod, the request HTTP method (e.g. GET, POST, PUT, DELETE, PATCH, HEAD)
	ReqMethod string `json:"req_method,omitempty" yaml:"req_method,omitempty"`
	// ApiURL is the URL used to contact the data source (e.g. PostgREST, Solr, Elasticsearch)
	ApiURL string `json:"api_url,omitempty" yaml:"api_url,omitempty"`

	// ApiMethod indicates the HTTP method to be used to contact the data source api
	ApiMethod string `json:"api_method,omitempty" yaml:"api_method,omitempty"`

	// ApiContentType is the content type to be used to contact the data source api
	ApiContentType string `json:"api_content_type,omitempty" yaml:"api_content_type,omitempty"`

	// PandocTemplate holds the source to a Pandoc template
	PandocTemplate string `json:"pandoc_template,omitempty" yaml:"pandoc_template,omitempty"`

	// ResHeaders holds any additional response headers to send back to the
	// browser or front end web server
	ResHeaders map[string]string `json:"res_headers,omitempty" yaml:"res_headers,omitempty"`
}

// String shows a JSON Representation of an individual route
func (route *Route) String() string {
	src, _ := json.MarshalIndent(route, "", "    ")
	return fmt.Sprintf("%s", src)
}

// Router holds the actual router structure. This 
type Router struct {
	Env    map[string]string
	Vars   map[string]string
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
	for _, method := range []string{ http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodHead, "OPTION"} {
		if strings.Compare(method, s) == 0 {
			return true
		}
	}
	return false
}

// ReadYAML filename reads router configuration from a file.
func (router *Router) ReadYAML(fName string) error {
	in, err := os.Open(fName)
	if err != nil {
		return err
	}
	defer in.Close()
	src, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	tmpRouter := &Router{}
	if err = yaml.Unmarshal(src, tmpRouter); err != nil {
		return err
	}
	if tmpRouter.Env != nil {
		for k, v := range tmpRouter.Env {
			router.Env[k] = v
		}
	}
	if tmpRouter.Vars != nil {
		for k, v := range tmpRouter.Vars {
			router.Vars[k] = v
		}
	}
	if tmpRouter.Routes != nil {
		router.Routes = tmpRouter.Routes
	}
 	return nil	
}

// Configure takes a Config object and maps in the relatavent values to
// the router.
func (router *Router) Configure(cfg *Config) {
	// Copy the environment into the router
	for _, envar := range cfg.Env {
		router.Env[envar] = os.Getenv(envar)
	}
	// Copy cfg.Routes into router.Routes
	for _, route := range cfg.Routes {
		router.Routes = append(router.Routes, route)
	}
}

// ReadCSV filename reads router configuration from a file.
func (router *Router) ReadCSV(fName string) error {
	in, err := os.Open(fName)
	if err != nil {
		return err
	}
	defer in.Close()
	r := csv.NewReader(in)
	r.Comment = '#'
	r.LazyQuotes = true
	r.TrimLeadingSpace = true
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
				// Figure out where the columns are off, and return error
				for i, k := range routerColumns {
					if i >= len(record) || record[i] != k {
						return fmt.Errorf("missing columns %q, line %d, starting with %s", fName, rowNo, k)
					}
				}
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

			if record[varDef] != "" {
				route.Var = map[string]string{}
				if err = json.Unmarshal([]byte(record[varDef]), &route.Var); err != nil {
					return fmt.Errorf("var error, line %d in %s, %s", rowNo, fName, err)
				}
			}
			route.ReqPath, err = NewDSL(record[reqPath], route.Var)
			if err != nil {
				return fmt.Errorf("req_path error, line %d in %q, %s", rowNo, fName, err)
			}
			if !isHTTPMethod(record[reqMethod]) {
				return fmt.Errorf("rep_method error, line %d in %q, unsupported method %q", rowNo, fName, record[reqMethod])
			}
			route.ReqMethod = record[reqMethod]
			// FIXME: need to make sure this make sense
			route.ApiURL = record[apiURL]
			if !isHTTPMethod(record[reqMethod]) {
				return fmt.Errorf("api_method error, line %d in %q, unsupported method %q", rowNo, fName, record[reqMethod])
			}
			route.ApiMethod = record[apiMethod]
			// FIXME: validate content type
			route.ApiContentType = record[apiContentType]
			templateName := record[pandocTemplate]
			if templateName != "" {
				tSrc, err := os.ReadFile(templateName)
				if err != nil {
					return fmt.Errorf("failed to load template %q, line %d in %q, %s", templateName, rowNo, fName, err)
				}
				route.PandocTemplate = fmt.Sprintf("%s", tSrc)
			}
			if record[resHeaders] != "" {
				headers := map[string]string{}
				if err := json.Unmarshal([]byte(record[resHeaders]), &headers); err != nil {
					return fmt.Errorf("failed to parse JSON expression for headers, line %d in %q, %s", rowNo, fName, err)
				}
				route.ResHeaders = headers
			}
			router.Routes = append(router.Routes, route)
		}
		rowNo++
	}
	return nil
}

func (route *Route) HasReqMethod(method string) bool {
	if route.ReqMethod == method {
		return true
	}
	return false
}

func (router *Router) ResolveRoute(u string, method string) (int, map[string]string, bool) {
	for i, r := range router.Routes {
		// Compare our HTTP method then route to determine response.
		if router.Routes[i].HasReqMethod(method) {
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
		return router.Routes[no].ReqPath.Resolve(m, router.Routes[no].ApiURL), true
	}
	return "", false
}

// RequestDataAPI
func (router *Router) RequestDataAPI(rNo int, apiURL string, body []byte) ([]byte, string, int) {
	if rNo > 0 && rNo >= len(router.Routes) {
		return nil, http.StatusText(502), 502
	}
	// FIXME: Make an http call to the JSON data API
	method, contentType := strings.ToUpper(router.Routes[rNo].ApiMethod), router.Routes[rNo].ApiContentType
	var (
		res *http.Response
		err error
	)
	buf := bytes.NewBuffer(body)
	switch method {
	case http.MethodGet:
		res, err = http.Get(apiURL)
	case http.MethodPost:
		res, err = http.Post(apiURL, contentType, buf)
	default:
		log.Printf("http method not supported %s", method)
		return nil, http.StatusText(502), 502
	}
	if err != nil {
		log.Printf("%s", err)
		// Bad Gateway
		return nil, http.StatusText(502), 502
	}
	defer res.Body.Close()
	src, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("%s", err)
		// Bad Gateway
		return nil, http.StatusText(502), 502
	}
	return src, http.StatusText(200), 200
}

// toFrontMatter converts JSON source to front matter
// for a Pandoc Markdown document.
func JSONSrcToFrontMatter(src []byte) (string, error) {
	if len(src) == 0 {
		return "", nil
	}
	var data interface{}
	if err := json.Unmarshal(src, &data); err != nil {
		return "", err
	}
	metadata := map[string]interface{}{
		"data": data,
	}
	mSrc, err := yaml.Marshal(metadata)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("---\n%s\n---\n\n", mSrc), nil
}

// RequestPandoc
func (router *Router) RequestPandoc(rNo int, src []byte) ([]byte, string, int) {
	template := router.Routes[rNo].PandocTemplate
	u := "http://localhost:3030"
	// NOTE: Need to unpack our JSON data and repack it as front matter with
	// the name "data".
	txt, err := JSONSrcToFrontMatter(src)
	if err != nil {
		log.Printf("RequestPandoc(%d, %s), failed to build text from JSON resource, %s", rNo, src, err)
		return nil, http.StatusText(502), 502
	}
	m := map[string]interface{}{}
	m["standalone"] = true
	m["text"] = txt
	if len(template) > 0 {
		m["template"] = template
	}
	// FIXME: Make an http call to Pandoc server
	buf, err := json.Marshal(m)
	if err != nil {
		log.Printf("RequestPandoc(%d, src), %s", rNo, err)
		return nil, http.StatusText(502), 502
	}
	res, err := http.Post(u, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		log.Printf("RequestPandoc(%d, src), %s", rNo, err)
		return nil, http.StatusText(502), 502
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("RequestPandoc(%d, src), %s", rNo, err)
		return nil, http.StatusText(res.StatusCode), res.StatusCode
	}
	return body, http.StatusText(res.StatusCode), res.StatusCode
}

// WriteResponse
func (router *Router) WriteResponse(w http.ResponseWriter, rNo int, src []byte) {
	for k, v := range router.Routes[rNo].ResHeaders {
		w.Header().Set(k, v)
	}
	w.WriteHeader(http.StatusOK)
	// Write back the content
	fmt.Fprintf(w, "%s", src)
}

// Newt implements the router as an http middleware. This allows our
// router to be used more generally and gives us options like having
// the Newt application be both a data router and static file server.
func (router *Router) Newt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rNo, m, ok := router.ResolveRoute(req.URL.Path, req.Method)
		if !ok {
			// NOTE: We don't match a route, pass this on to the next handler
			next.ServeHTTP(w, req)
			return
		}
		var (
			body []byte
			err  error
		)
		// FIXME: Handle http.MethodOption
		// FIXME: Handle http.MethodHead
		// FIXME: Handle http.MethodGet
		// FIXME: Handle http.MethodDelete
		if (req.Method == http.MethodPost) ||
			(req.Method == http.MethodPut) ||
			(req.Method == http.MethodPatch) {
			// FIXME: I only want to use ParseForm is
			// the client is NOT posting JSON content.
			if err := req.ParseForm(); err != nil {
				http.Error(w, fmt.Sprintf("%s", err), 400)
				return
			}
			m := map[string]interface{}{}
			for k, _ := range req.Form {
				m[k] = req.FormValue(k)
			}
			body, err = json.Marshal(m)
			if err != nil {
				http.Error(w, fmt.Sprintf("%s", err), 400)
				return
			}

			//fmt.Fprintf(os.Stderr, "DEBUG body (%q) \n-->%s<--\n", req.Method, body)
		} 

		// Get content from the JSON Data API
		apiURL, ok := router.ResolveApiURL(rNo, m)
		if !ok {
			// Handle 502, Bad Gateway
			http.Error(w, http.StatusText(502), 502)
			return
		}
		src, statusText, statusCode := router.RequestDataAPI(rNo, apiURL, body)
		if statusCode < 200 || statusCode >= 300 {
			// echo back the data request status code and text
			http.Error(w, statusText, statusCode)
			return
		}
		// NOTE: if Pandoc transform data request
		if router.Routes[rNo].PandocTemplate != "" {
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
