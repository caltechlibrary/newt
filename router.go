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
			// FIXME: validate method
			route.ReqMethod  = record[reqMethod]
			// FIXME: validate content type
			route.ReqContentType = record[reqContentType]
			// FIXME: need to make sure this make sense
			route.ApiURL = record[apiURL]
			// FIXME: validate method
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
			fmt.Fprintf(os.Stderr, "DEBUG adding route %s\n", route.String())
			router.Routes = append(router.Routes, route)
		}
		rowNo++
	}
	return nil
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


