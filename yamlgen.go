package newt

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// getAnswer get a Y/N response from buffer
func getAnswer(buf *bufio.Reader, defaultAnswer string, lower bool) string {
	answer, err := buf.ReadString('\n')
	if err != nil {
		return ""
	}
	answer = strings.TrimSpace(answer)
	if answer == "" {
		answer = defaultAnswer
	}
	if lower {
		return strings.ToLower(answer)
	}
	return answer
}

// setupNewtRouter prompt to configure the NewtRouter
func setupNewtRouter(cfg *Config, buf *bufio.Reader, out io.Writer, appFName string, objName string) {
	fmt.Fprintf(out, "Will %s use Newt Router (Y/n)? ", appFName)
	answer := getAnswer(buf, "y", true)
	if answer == "y" {
		if cfg.Applications == nil {
			cfg.Applications = &Applications{ }
		}
		if cfg.Applications.NewtRouter == nil {
			cfg.Applications.NewtRouter = &Application{}
		}
		// NOTE: If port is zero, we haven't configure the router.
		if cfg.Applications.NewtRouter.Port == 0 {
			cfg.Applications.NewtRouter.Port = ROUTER_PORT
			if info, err := os.Stat("htdocs"); err == nil && info.IsDir() {
				cfg.Applications.NewtRouter.Htdocs = "htdocs"
			} else {
				cfg.Applications.NewtRouter.Htdocs = ""
			}
		}

		if cfg.Routes == nil {
			cfg.Routes = []*NewtRoute{}
		}
	} else {
		if cfg.Applications != nil {
			cfg.Applications.NewtRouter = nil
		}
	}
}

// setupPostgREST prompt to configure PostgREST
func setupPostgREST(cfg *Config, buf *bufio.Reader, out io.Writer, appFName string, objName string) {
	fmt.Fprintf(out, "Will %s use PostgREST (Y/n)? ", appFName)
	answer := getAnswer(buf, "y", true)
	if answer == "y" {
		if cfg.Applications == nil {
			cfg.Applications = &Applications{}
		}
		if cfg.Applications.PostgREST == nil {
			cfg.Applications.PostgREST = &Application{}
		}
		if cfg.Applications.PostgREST.Port == 0 {
			cfg.Applications.PostgREST.Port = POSTGREST_PORT
			cfg.Applications.PostgREST.AppPath = "postgrest"
			cfg.Applications.PostgREST.ConfPath = "postgrest.conf"
		}
	} else {
		if cfg.Applications != nil {
			cfg.Applications.PostgREST = nil
		}
	}
}

// setupPostgres prompt to configure Postgres
func setupPostgres(cfg *Config, buf *bufio.Reader, out io.Writer, appFName string, objName string) {
	fmt.Fprintf(out, "Will %s use Postgres (Y/n)? ", appFName)
	answer := getAnswer(buf, "y", true)
	if answer == "y" {
		if cfg.Applications == nil {
			cfg.Applications = &Applications{}
		}
		if cfg.Applications.Postgres == nil {
			cfg.Applications.Postgres = &Application{}
		}
		if cfg.Applications.Postgres.Port == 0 {
			cfg.Applications.Postgres.Port = POSTGRES_PORT
		}
		if cfg.Applications.Postgres.DSN == "" {
			cfg.Applications.Postgres.DSN = fmt.Sprintf("postgres://{PGUSER}:{PGPASSWORD}@localhost:%d/%s", cfg.Applications.Postgres.Port, objName)
			// Now we need to make sure we allow PGUSER and PGPASSWORD to pass through in the environment
			if cfg.Applications.Environment == nil {
				cfg.Applications.Environment = []string{}
			}
			cfg.Applications.Environment = append(cfg.Applications.Environment, "PGUSER", "PGPASSWORD")
		}
	} else {
		if cfg.Applications != nil {
			cfg.Applications.Postgres = nil
		}
	}
}


func setupNewtMustache(cfg *Config, buf *bufio.Reader, out io.Writer, appFName string, objName string) {
	fmt.Fprintf(out, "Will %s use Newt Mustache (Y/n)? ", appFName)
	answer := getAnswer(buf, "y", true)
	if answer == "y" {
		if cfg.Applications == nil {
			cfg.Applications = &Applications{}
		}
		if cfg.Applications.NewtMustache == nil {
			cfg.Applications.NewtMustache = &Application{}
		}
		if cfg.Applications.NewtMustache.Port == 0 {
			cfg.Applications.NewtMustache.Port = MUSTACHE_PORT
		}
		//FIXME: If there are models then templates will need to be updates even when it is NOT nil.
		// When the model list changes then the related templates should change to.
		// A scan of the template routes for removed models needs to happen when the model is "removed" by the modeler.
		if cfg.Templates == nil {
			cfg.Templates = []*MustacheTemplate{}
			// Handle the special cases of routes for retrieving forms for create, update and delete.
			// E.g. retrieve the web form, handle the submit of the web form as two actions.
			setupWebFormHandling(cfg, objName, "create")
			setupWebFormHandling(cfg, objName, "update")
			setupWebFormHandling(cfg, objName, "delete")
			// Now add the mappings for read and list
			setupReadHandling(cfg, objName, "read")
			setupReadHandling(cfg, objName, "list")
		}
	} else {
		if cfg.Applications != nil {
			cfg.Applications.NewtMustache = nil
		}
	}
}

func setupNewtGenerator(cfg *Config, buf *bufio.Reader, out io.Writer, appFName string, objName string) {
	fmt.Fprintf(out, "Will %s use Newt Generator (Y/n)? ", appFName)
	answer := getAnswer(buf, "y", true)
	if answer == "y" {
		if cfg.Applications == nil {
			cfg.Applications = &Applications{}
		}
		if cfg.Applications.NewtGenerator == nil {
			cfg.Applications.NewtGenerator = &Application{}
		}
		if cfg.Applications.NewtGenerator.Namespace == "" {
			cfg.Applications.NewtGenerator.Namespace = objName
		}
		if cfg.Models == nil {
			// FIXME: This is the same add adding a model in modeler.go so this code needs to be unified.
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
	} else {
		if cfg.Applications != nil {
			cfg.Applications.NewtGenerator = nil
		}
	}
}


// setupPostgRESTService creates a Service object for interacting with PostgREST
func setupPostgRESTService(cfg *Config, objName string, action string) *Service {
	var (
		oidSuffix string
		description string
		method string
		port int
	)
	description = fmt.Sprintf("Access PostgREST API for %s %s", objName, action)
	if cfg.Applications != nil && cfg.Applications.PostgREST != nil {
		port = cfg.Applications.PostgREST.Port
	} else {
		port = 3000
	}
	switch action {
	case "create":
		// create action doesn't take an oid
		method = http.MethodPost
	case "read":
		method = http.MethodGet
		oidSuffix = "/{oid}"
	case "update": 
		method = http.MethodPut
		oidSuffix = "/{oid}"
	case "delete":
		method = http.MethodDelete
		oidSuffix = "/{oid}"
	default:
		// list action doesn't take an oid
		method = http.MethodGet
	}
	return &Service {
		Service: fmt.Sprintf("%s http://localhost:%d/rpc/%s_%s%s", method, port, objName, action, oidSuffix),
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
func setupWebFormHandling(cfg *Config, objName string, action string) {
	var (
		pathSuffix string
		service *Service
	)
	if action == "update" || action == "delete" {
		pathSuffix = "/{oid}"
	}
	// Setup templates and webforms. Names are formed by objName combined with action.
	tmplPattern := fmt.Sprintf("/%s_%s", objName, action)
	tmplName := fmt.Sprintf("%s_%s_form.tmpl", objName, action)
	tmplDescription := fmt.Sprintf("Display a %s for %s", objName, action)
	cfg.Templates = append(cfg.Templates, &MustacheTemplate{
		Pattern:     tmplPattern,
		Template:    tmplName,
		Description: tmplDescription,
	})
	// Handle web form request
	routeId := fmt.Sprintf("%s_%s", objName, action)
	request := fmt.Sprintf("%s /%s_%s%s", http.MethodGet, objName, action, pathSuffix)
	routeDescription := fmt.Sprintf("Handle retrieving the webform for %s %s", objName, action)
	route := &NewtRoute{
		Id:          routeId,
		Pattern:     request,
		Description: routeDescription,
		Pipeline: []*Service{},
	}
	// NOTE: If we have an update or delete we want to retrieve the record before calling the template
	if action == "update" || action == "delete" {
		service = setupPostgRESTService(cfg, objName, "read")
		service.Description = fmt.Sprintf("Retrieve %s from PostgREST API before %s", objName, action)
		route.Pipeline = append(route.Pipeline, service)
	}
	service = setupTmplService(cfg, tmplPattern, tmplDescription)
	route.Pipeline = append(route.Pipeline, service)
	cfg.Routes = append(cfg.Routes, route)

	// Setup template submit result
	tmplPattern = fmt.Sprintf("/%s_%s_response", objName, action)
	tmplName = fmt.Sprintf("%s_%s_response.tmpl", objName, action)
	tmplDescription = fmt.Sprintf("This is an result template for %s %s", objName, action)
	cfg.Templates = append(cfg.Templates, &MustacheTemplate{
		Pattern:     tmplPattern,
		Template:    tmplName,
		Description: tmplDescription,
	})

	// Handle submission routing 
	routeId = fmt.Sprintf("%s_%s", objName, action)
	routeDescription = fmt.Sprintf("Handle form submission for %s %s", objName, action)
	request = fmt.Sprintf("%s /%s_%s", http.MethodPost, objName, action)
	route = &NewtRoute{
		Id:          routeId,
		Pattern:     request,
		Description: routeDescription,
		Pipeline: []*Service{},
	}
	service = setupPostgRESTService(cfg, objName, action)
	route.Pipeline = append(route.Pipeline, service)
	service = setupTmplService(cfg, tmplPattern, tmplDescription)
	route.Pipeline = append(route.Pipeline, service)
	cfg.Routes = append(cfg.Routes, route)
}

func setupReadHandling(cfg *Config, objName string, action string) {
	// Setup template for results of read request
	tmplPattern := fmt.Sprintf("/%s_%s", objName, action)
	tmplName := fmt.Sprintf("%s_%s.tmpl", objName, action)
	tmplDescription := fmt.Sprintf("This template handles %s %s", objName, action)
	cfg.Templates = append(cfg.Templates, &MustacheTemplate{
		Pattern:     tmplPattern,
		Template:    tmplName,
		Description: tmplDescription,
	})
	// Handle requesting object or list of objects
	routeId := fmt.Sprintf("%s_%s", objName, action)
	routeDescription := fmt.Sprintf("Retrieve object(s) for %s %s", objName, action)
	request := fmt.Sprintf("%s /%s_%s", http.MethodPost, objName, action)
	route := &NewtRoute{
		Id:          routeId,
		Pattern:     request,
		Description: routeDescription,
		Pipeline: []*Service{},
	}
	service := setupPostgRESTService(cfg, objName, action)
	route.Pipeline = append(route.Pipeline, service)
	service = setupTmplService(cfg, tmplPattern, tmplDescription)
	route.Pipeline = append(route.Pipeline, service)
	cfg.Routes = append(cfg.Routes, route)
}

func setupEnvironment(cfg *Config, buf *bufio.Reader, out io.Writer, appFName string, objName string) {
	fmt.Fprintf(out, "Will %s need to import environment variables (y/N)? ", appFName)
	answer := getAnswer(buf, "n", true)
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
			answer = getAnswer(buf, "", false)
			if answer != "" {
				cfg.Applications.Environment = append(cfg.Applications.Environment, answer)
			}
		}
	}
}

func setupOptions(cfg *Config, buf *bufio.Reader, out io.Writer, appFName string, objFName string) {
	fmt.Fprintf(out, "Will %s provide options to the services (y/N)? ", appFName)
	answer := getAnswer(buf, "n", true)
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
			answer = getAnswer(buf, "", false)
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

