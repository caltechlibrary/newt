package newt

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// setupPostgRESTService creates a Service object for interacting with PostgREST
func setupPostgRESTService(cfg *Config, action string, objName string) *Service {
	var (
		serviceURL string
		description string
		method string
		port int
	)
	switch action {
	case "create":
		method = http.MethodPost
	case "update": 
		method = http.MethodPut
	case "delete":
		method = http.MethodDelete
	default:
		method = http.MethodGet
	}
	description = fmt.Sprintf("Access PostgREST API for %s %s", action, objName)
	if cfg.Applications != nil && cfg.Applications.PostgREST != nil {
		port = cfg.Applications.PostgREST.Port
	} else {
		port = 3000
	}
	serviceURL = fmt.Sprintf("%s http://localhost:%d/rpc/%s_%s", method, port, objName, action)
	return &Service {
		Service: serviceURL,
		Description: description,
	}
}

// setupTemplService creates a Service object to process with a template
func setupTmplService(cfg *Config, tmplPattern string, description string) *Service {
	var port int
	if cfg.Applications != nil && cfg.Applications.NewtMustache != nil {
		port = cfg.Applications.NewtMustache.Port
	} else {
		port = 8011
	}
	serviceURL := fmt.Sprintf("POST http://localhost:%d%s", port, tmplPattern)
	return &Service {
		Service: serviceURL,
		Description: description,
	}
}

// setupWebFormHandling generates the routes and template handling for retrieving and submitting
// webforms for "create", "update" or "delete".
func setupWebFormHandling(cfg *Config, action string, objName string) {
	var (
		pathSuffix string
		service *Service
	)
	if action == "update" || action == "delete" {
		pathSuffix = "/{oid}"
	}
	// Setup templates, webform
	tmplPattern := fmt.Sprintf("/%s_%s", action, objName)
	tmplName := "page.tmpl" // fmt.Sprintf("page_%s.tmpl", objName)
	partialName := fmt.Sprintf("%s_%s_form.tmpl", action, objName)
	tmplDescription := fmt.Sprintf("Display a %s for %s", action, objName)
	cfg.Templates = append(cfg.Templates, &MustacheTemplate{
		Pattern:     tmplPattern,
		Template:    tmplName,
		Partials: []string{
			partialName,
		},
		Description: tmplDescription,
	})
	// Handle web form request
	routeId := fmt.Sprintf("%s_%s", action, objName)
	request := fmt.Sprintf("%s /%s_%s%s", http.MethodGet, action, objName, pathSuffix)
	routeDescription := fmt.Sprintf("Handle retrieving the webform for %s %s", action, objName)
	route := &NewtRoute{
		Id:          routeId,
		Pattern:     request,
		Description: routeDescription,
		Pipeline: []*Service{},
	}
	if pathSuffix != "" {
		// NOTE: If we have an update or delete we want to retrieve the record before calling the template
		service = setupPostgRESTService(cfg, "read", objName)
		route.Pipeline = append(route.Pipeline, service)
	}
	service = setupTmplService(cfg, tmplPattern, tmplDescription)
	route.Pipeline = append(route.Pipeline, service)
	cfg.Routes = append(cfg.Routes, route)

	// Setup template submit result
	tmplPattern = fmt.Sprintf("/%s_%s_response", action, objName)
	tmplName = "page.tmpl" // fmt.Sprintf("page_%s.tmpl", objName)
	partialName = fmt.Sprintf("%s_%s_response.tmpl", action, objName)
	tmplDescription = fmt.Sprintf("This is an result template for %s %s", action, objName)
	cfg.Templates = append(cfg.Templates, &MustacheTemplate{
		Pattern:     tmplPattern,
		Template:    tmplName,
		Partials: []string{
			partialName,
		},
		Description: tmplDescription,
	})

	// Handle submission routing 
	routeId = fmt.Sprintf("%s_%s", action, objName)
	routeDescription = fmt.Sprintf("Handle form submission for %s %s", action, objName)
	request = fmt.Sprintf("%s /%s_%s_response", http.MethodPost, action, objName)
	route = &NewtRoute{
		Id:          routeId,
		Pattern:     request,
		Description: routeDescription,
		Pipeline: []*Service{},
	}
	service = setupPostgRESTService(cfg, action, objName)
	route.Pipeline = append(route.Pipeline, service)
	service = setupTmplService(cfg, tmplPattern, tmplDescription)
	route.Pipeline = append(route.Pipeline, service)
	cfg.Routes = append(cfg.Routes, route)
}

func setupReadHandling(cfg *Config, action string, objName string) {
	// Setup template for results of read request
	tmplPattern := fmt.Sprintf("/%s_%s", action, objName)
	tmplName := "page.tmpl" // fmt.Sprintf("page_%s.tmpl", objName)
	partialName := fmt.Sprintf("%s_%s.tmpl", action, objName)
	tmplDescription := fmt.Sprintf("This template handles %s %s", action, objName)
	cfg.Templates = append(cfg.Templates, &MustacheTemplate{
		Pattern:     tmplPattern,
		Template:    tmplName,
		Partials: []string{
			partialName,
		},
		Description: tmplDescription,
	})
	// Handle requesting object or list of objects
	routeId := fmt.Sprintf("%s_%s", action, objName)
	routeDescription := fmt.Sprintf("Retrieve object(s) for %s %s", action, objName)
	request := fmt.Sprintf("%s /%s_%s", http.MethodPost, action, objName)
	route := &NewtRoute{
		Id:          routeId,
		Pattern:     request,
		Description: routeDescription,
		Pipeline: []*Service{},
	}
	service := setupPostgRESTService(cfg, action, objName)
	route.Pipeline = append(route.Pipeline, service)
	service = setupTmplService(cfg, tmplPattern, tmplDescription)
	route.Pipeline = append(route.Pipeline, service)
	cfg.Routes = append(cfg.Routes, route)
}

func setupNewtMustache(cfg *Config, buf *bufio.Reader, out io.Writer, appFName string, objName string) {
	fmt.Fprintf(out, "Will %s use Newt Mustache (Y/n)? ", appFName)
	answer := getAnswer(buf, true)
	if answer != "n" {
		if cfg.Applications.NewtMustache == nil {
			cfg.Applications.NewtMustache = &Application{
				Port: MUSTACHE_PORT,
			}
		}
		if cfg.Templates == nil {
			cfg.Templates = []*MustacheTemplate{}
			// Handle the special cases of routes for retrieving forms for create, update and delete.
			// E.g. retrieve the web form, handle the submit of the web form as two actions.
			setupWebFormHandling(cfg, "create", objName)
			setupWebFormHandling(cfg, "update", objName)
			setupWebFormHandling(cfg, "delete", objName)
			// Now add the mappings for read and list
			setupReadHandling(cfg, "read", objName)
			setupReadHandling(cfg, "list", objName)
		}
	}
}

func setupNewtGenerator(cfg *Config, buf *bufio.Reader, out io.Writer, appFName string, objName string) {
	fmt.Fprintf(out, "Will %s use Newt Generator (Y/n)? ", appFName)
	answer := getAnswer(buf, true)
	if answer != "n" {
		if cfg.Applications.NewtGenerator == nil {
			cfg.Applications.NewtGenerator = &Application{
				Namespace: objName,
			}
		}
		if cfg.Models == nil {
			cfg.Models = []*NewtModel{}
			cfg.Models = append(cfg.Models, &NewtModel{
				Id:          objName,
				Description: "This is where you would model your application data",
				Body: []*Element{
					&Element{
						Id:   "data_attribute",
						Type: "input",
						Attributes: map[string]string{
							"name":            "data_attribute",
							"description":     "This is an example input element",
							"placeholdertext": "ex. of placeholder text",
							"title":           "this is an example element in your model",
						},
						Validations: map[string]interface{}{
							"required": true,
						},
					},
				},
			})
		}
	}
}

func setupEnvironment(cfg *Config, buf *bufio.Reader, out io.Writer, appFName string, objName string) {
	fmt.Fprintf(out, "Will %s need to import environment variables (y/N)? ", appFName)
	answer := getAnswer(buf, true)
	if answer == "y" {
		if cfg.Applications.Environment == nil {
			cfg.Applications.Environment = []string{}
		}
		if len(cfg.Applications.Environment) > 0 {
			fmt.Fprintf(out, "You currently have the following environment defined:\n\t%s\n",
				strings.Join(cfg.Applications.Environment, "\n\t"))
		}
		fmt.Fprintf(out, "Enter the environment variable name (one per line, enter empty line when complete):\n")
		answer = " "
		for answer != "" {
			answer = getAnswer(buf, false)
			if answer != "" {
				cfg.Applications.Environment = append(cfg.Applications.Environment, answer)
			}
		}
	}
}

func setupOptions(cfg *Config, buf *bufio.Reader, out io.Writer, appFName string, objFName string) {
	fmt.Fprintf(out, "Will %s provide options to the services (y/N)? ", appFName)
	answer := getAnswer(buf, true)
	if answer == "y" {
		if cfg.Applications.Options == nil {
			cfg.Applications.Options = map[string]string{}
		}
		if len(cfg.Applications.Options) > 0 {
			fmt.Fprintf(out, "You currently have the following options defined:\n")
			for k, v := range cfg.Applications.Options {
				fmt.Fprintf(out, "\t%s -> %q\n", k, v)
			}
		}
		fmt.Fprintf(out, "Enter the options (separated key/value by colon, enter empty line when complete):\n")
		answer = " "
		for answer != "" {
			answer = getAnswer(buf, false)
			if strings.Contains(answer, ":") {
				parts := strings.SplitN(answer, ":", 2)
				k := strings.ReplaceAll(strings.TrimSpace(parts[0]), " ", "_")
				v := strings.TrimSpace(parts[1])
				cfg.Applications.Options[k] = v
			} else if answer != "" {
				fmt.Fprintf(out, "%q is missing a colon, can't tell key from value, try again\n", answer)
			}
		}
	}
}

