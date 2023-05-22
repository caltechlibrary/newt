package newt

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	pandocPort
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
		"pandoc_port",
		"pandoc_options",
		"pandoc_template",
		"res_headers",
	}
)

type Route struct {
	ReqPath        *RouteDSL
	ReqMethod      string
	ReqContentType string
	ApiURL         string
	ApiMethod      string
	ApiContentType string
	PandocPort     string
	PandocOptions  string
	PandocTemplate string
	ResHeaders     map[string]string
}

/*FIXME:
func (route *Route) ResolveApiURL() (string , error) {
	return route.ReqPath.Resolve(route.ApiURL)
}
*/


type Router struct {
	Env    map[string]string
	Routes []*Route
}

// ReadCSV filename
func (router *Router) ReadCSV(fName string) error {
	return fmt.Errorf("ReadCSV not limited")
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

func rowToRoute(rowNo int, row []string) (*Route, error) {
	route := new(Route)
	for coNo, val := range row {
		val := strings.TrimSpace(val)
		switch coNo {
		case reqPath:
			rdsl, err := NewRouteDSL(val)
			if err != nil {
				return nil, fmt.Errorf("row %d req_path %q, %s", rowNo, val, err)
			}
			route.ReqPath = rdsl
		case reqMethod:
			if !isHTTPMethod(val) {
				return nil, fmt.Errorf("row %d req_method %q, %s", rowNo, val, "unsupport HTTP method")
			}
			route.ReqMethod = strings.ToUpper(val)
		case reqContentType:
			route.ReqContentType = val
		case apiURL:
			_, err := url.Parse(strings.ReplaceAll(strings.ReplaceAll(val, "{", ""), "}", ""))
			if err != nil {
				return nil, fmt.Errorf("row %d api_url %q, %s", rowNo, val, err)
			}
			route.ApiURL = val
		case apiMethod:
			if !isHTTPMethod(val) {
				return nil, fmt.Errorf("row %d api_method %q, %s", rowNo, val, "unsupport HTTP method")
			}
			route.ApiMethod = val
		case apiContentType:
			route.ApiContentType = val
		case pandocPort:
			route.PandocPort = strings.TrimPrefix(val, ":")
		case pandocOptions:
			route.PandocOptions = val
		case pandocTemplate:
			if val != "" {
				data, err := os.ReadFile(val)
				if err != nil {
					return nil, fmt.Errorf("row %d, template %q, %s", rowNo, val, err)
				}
				route.PandocTemplate = fmt.Sprintf("%s", data)
			}
		case resHeaders:
			if val != "" {
				m := map[string]string{}
				if err := json.Unmarshal([]byte(val), &m); err != nil {
					return nil, fmt.Errorf("row %d, res_headers %q, %s", rowNo, val, err)
				}
				route.ResHeaders = m
			}
		}
	}
	return route, nil
}

func RouterFromCSV(fName string) (*Router, error) {
	in, err := os.Open(fName)
	if err != nil {
		return nil, err
	}
	defer in.Close()

	router := new(Router)
	r := csv.NewReader(in)
	rowNo := 0
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("row %d: %s", rowNo+1, err)
		}
		if rowNo == 0 {
			// Verify column names
			if len(row) < len(routerColumns) {
				return nil, fmt.Errorf("missing column names")
			}
			for i, colName := range routerColumns {
				if strings.Compare(colName, row[i]) != 0 {
					return nil, fmt.Errorf("expected column %q, got %q", colName, row[i])
				}
			}
		} else {
			route, err := rowToRoute(rowNo, row)
			if err != nil {
				return nil, err
			}
			if route.ReqPath != nil {
				router.Routes = append(router.Routes, route)
			}
		}
		rowNo++
	}
	return router, nil
}

