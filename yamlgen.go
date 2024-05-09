package newt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

// setupRouter prompt to configure the Router
func setupRouter(ast *AST, buf *bufio.Reader, out io.Writer, appFName string, skipPrompts bool) {
	var answer string
	if skipPrompts {
		answer = "y"
	} else {
		fmt.Fprintf(out, "Will %s use Newt Router (Y/n)? ", appFName)
		answer = getAnswer(buf, "y", true)
	}
	if answer == "y" {
		if ast.Applications == nil {
			ast.Applications = &Applications{}
		}
		if ast.Applications.Router == nil {
			ast.Applications.Router = &Application{}
		}
		// NOTE: If port is zero, we haven't configure the router.
		if ast.Applications.Router.Port == 0 {
			ast.Applications.Router.Port = ROUTER_PORT
			if info, err := os.Stat("htdocs"); err == nil && info.IsDir() {
				ast.Applications.Router.Htdocs = "htdocs"
			} else {
				ast.Applications.Router.Htdocs = ""
			}
		}

		if ast.Routes == nil {
			ast.Routes = []*Route{}
		}
	} else {
		if ast.Applications != nil {
			ast.Applications.Router = nil
		}
	}
}

// setupPostgREST prompt to configure PostgREST
func setupPostgREST(ast *AST, buf *bufio.Reader, out io.Writer, appFName string, skipPrompts bool) {
	var answer string
	if skipPrompts {
		answer = "y"
	} else {
		fmt.Fprintf(out, "Will %s use PostgREST (Y/n)? ", appFName)
		answer = getAnswer(buf, "y", true)
	}
	if answer == "y" {
		if ast.Applications == nil {
			ast.Applications = &Applications{}
		}
		if ast.Applications.PostgREST == nil {
			ast.Applications.PostgREST = &Application{}
		}
		if ast.Applications.PostgREST.Port == 0 {
			ast.Applications.PostgREST.Port = POSTGREST_PORT
			ast.Applications.PostgREST.AppPath = "postgrest"
			ast.Applications.PostgREST.ConfPath = "postgrest.conf"
		}
	} else {
		if ast.Applications != nil {
			ast.Applications.PostgREST = nil
		}
	}
}

// setupPostgres prompt to configure Postgres
func setupPostgres(ast *AST, buf *bufio.Reader, out io.Writer, appFName string, skipPrompts bool) {
	var answer string
	if skipPrompts {
		answer = "y"
	} else {
		fmt.Fprintf(out, "Will %s use Postgres (Y/n)? ", appFName)
		answer = getAnswer(buf, "y", true)
	}
	if answer == "y" {
		if ast.Applications == nil {
			ast.Applications = &Applications{}
		}
		if ast.Applications.Postgres == nil {
			ast.Applications.Postgres = &Application{}
		}
		if ast.Applications.Postgres.Port == 0 {
			ast.Applications.Postgres.Port = POSTGRES_PORT
		}
		if ast.Applications.Postgres.DSN == "" {
			ast.Applications.Postgres.DSN = fmt.Sprintf("postgres://{PGUSER}:{PGPASSWORD}@localhost:%d/%s", ast.Applications.Postgres.Port, appFName)
			// Now we need to make sure we allow PGUSER and PGPASSWORD to pass through in the environment
			if ast.Applications.Environment == nil {
				ast.Applications.Environment = []string{}
			}
			ast.Applications.Environment = append(ast.Applications.Environment, "PGUSER", "PGPASSWORD")
		}
		if ast.Applications.NewtGenerator == nil {
			ast.Applications.NewtGenerator = &Application{
				Namespace: appFName,
			}
		}
	} else {
		if ast.Applications != nil {
			ast.Applications.Postgres = nil
		}
		if ast.Applications.NewtGenerator != nil {
			ast.Applications.NewtGenerator.Namespace = ""
		}
	}
}

func setupNewtMustache(ast *AST, buf *bufio.Reader, out io.Writer, appFName string, skipPrompts bool) {
	var answer string
	if skipPrompts {
		answer = "y"
	} else {
		fmt.Fprintf(out, "Will %s use Newt Mustache (Y/n)? ", appFName)
		answer = getAnswer(buf, "y", true)
	}
	if answer == "y" {
		if ast.Applications == nil {
			ast.Applications = &Applications{}
		}
		if ast.Applications.NewtMustache == nil {
			ast.Applications.NewtMustache = &Application{}
		}
		if ast.Applications.NewtMustache.Port == 0 {
			ast.Applications.NewtMustache.Port = MUSTACHE_PORT
		}
		//FIXME: If there are models then templates will need to be updates even when it is NOT nil.
		// When the model list changes then the related templates should change to.
		// A scan of the template routes for removed models needs to happen when the model is "removed" by the modeler.
		if ast.Templates == nil {
			ast.Templates = []*MustacheTemplate{}
			if err := setupMustacheTemplateHandlers(ast); err != nil {
				fmt.Fprintf(out, "WARNINGS: %s\n", err)
			}
		}
	} else {
		if ast.Applications != nil {
			ast.Applications.NewtMustache = nil
		}
	}
}

func setupMustacheTemplateHandlers(ast *AST) error {
	eBuf := bytes.NewBuffer([]byte{})
	hasError := false
	for _, m := range ast.Models {
		// Handle the special cases of routes for retrieving forms for create, update and delete.
		// E.g. retrieve the web form, handle the submit of the web form as two actions.
		for _, action := range []string{"create", "update", "delete"} {
			if err := setupWebFormHandling(ast, m, action); err != nil {
				fmt.Fprintf(eBuf, "%s\n", err)
				hasError = true
			}
		}
		// Now add the mappings for read and list
		for _, action := range []string{"read", "list"} {
			if err := setupReadHandling(ast, m, action); err != nil {
				fmt.Fprintf(eBuf, "%s\n", err)
				hasError = true
			}
		}
	}
	if hasError {
		return fmt.Errorf("%s", eBuf.Bytes())
	}
	return nil
}

func setupNewtGenerator(ast *AST, buf *bufio.Reader, out io.Writer, appFName string, skipPrompts bool) {
	var answer string
	if skipPrompts {
		answer = "y"
	} else {
		fmt.Fprintf(out, "Will %s use Newt Generator (Y/n)? ", appFName)
		answer = getAnswer(buf, "y", true)
	}
	if answer == "y" {
		if ast.Applications == nil {
			ast.Applications = &Applications{}
		}
		if ast.Applications.NewtGenerator == nil {
			ast.Applications.NewtGenerator = &Application{}
		}
		if ast.Applications.NewtGenerator.Namespace == "" {
			ast.Applications.NewtGenerator.Namespace = appFName
		}
		if ast.Models == nil {
			// FIXME: This is the same add adding a model in modeler.go so this code needs to be unified.
			ast.Models = []*Model{}
		}
	} else {
		if ast.Applications != nil {
			ast.Applications.NewtGenerator = nil
		}
	}
}

// setupPostgRESTService creates a Service object for interacting with PostgREST
func setupPostgRESTService(ast *AST, model *Model, action string) *Service {
	var (
		oid         string
		oidSuffix   string
		description string
		method      string
		port        int
	)
	objName := model.Id
	element, ok := model.GetModelIdentifier()
	if ok {
		oid = fmt.Sprintf("{%s}", element.Id)
	} else {
		oid = "{oid}"
	}
	description = fmt.Sprintf("Access PostgREST API for %s %s", objName, action)
	if ast.Applications != nil && ast.Applications.PostgREST != nil {
		port = ast.Applications.PostgREST.Port
	} else {
		port = 3000
	}
	switch action {
	case "create":
		// create action doesn't take an oid
		method = http.MethodPost
	case "read":
		method = http.MethodGet
		oidSuffix = "/" + oid
	case "update":
		method = http.MethodPut
		oidSuffix = "/" + oid
	case "delete":
		method = http.MethodDelete
		oidSuffix = "/" + oid
	default:
		// list action doesn't take an oid
		method = http.MethodGet
	}
	return &Service{
		Service:     fmt.Sprintf("%s http://localhost:%d/rpc/%s_%s%s", method, port, objName, action, oidSuffix),
		Description: description,
	}
}

// setupTemplService creates a Service object to process with a template
func setupTmplService(ast *AST, tmplPattern string, description string) *Service {
	var port int
	if ast.Applications != nil && ast.Applications.NewtMustache != nil {
		port = ast.Applications.NewtMustache.Port
	} else {
		port = 8011
	}
	serviceURL := fmt.Sprintf("POST http://localhost:%d%s", port, tmplPattern)
	return &Service{
		Service:     serviceURL,
		Description: description,
	}
}


// setupWebFormHandling generates the routes and template handling for retrieving and submitting
// webforms for "create", "update" or "delete".
func setupWebFormHandling(ast *AST, model *Model, action string) error {
	var (
		oid        string
		pathSuffix string
		service    *Service
	)
	objName := model.Id
	element, ok := model.GetModelIdentifier()
	if ok {
		oid = fmt.Sprintf("{%s}", element.Id)
	} else {
		oid = "{oid}"
	}
	if action == "update" || action == "delete" {
		pathSuffix = "/" + oid
	}
	// Setup templates and webforms. Names are formed by objName combined with action.
	templateList := ast.GetTemplateIds()
	templateId := mkName(model.Id, action, "")
	routeList := ast.GetRouteIds()
	routeId := mkName(objName, action, "")
	if inList(routeId, routeList) || inList(templateId, templateList) {
		return fmt.Errorf("routes and templates exist for %s %s", model.Id, action)
	}
	
	tSuffix := "_form.tmpl"
	tmplName := mkName(objName, action, tSuffix)
	tmplPattern := fmt.Sprintf("/%s_%s", objName, action)
	tmplDescription := fmt.Sprintf("Display a %s for %s", objName, action)
	ast.Templates = append(ast.Templates, &MustacheTemplate{
		Id: templateId,
		Pattern:     tmplPattern,
		Template:    tmplName,
		Description: tmplDescription,
	})

	// Handle web form request
	request := fmt.Sprintf("%s /%s_%s%s", http.MethodGet, objName, action, pathSuffix)
	routeDescription := fmt.Sprintf("Handle retrieving the webform for %s %s", objName, action)
	route := &Route{
		Id:          routeId,
		Pattern:     request,
		Description: routeDescription,
		Pipeline:    []*Service{},
	}
	// NOTE: If we have an update or delete we want to retrieve the record before calling the template
	if action == "update" || action == "delete" {
		service = setupPostgRESTService(ast, model, "read")
		service.Description = fmt.Sprintf("Retrieve %s from PostgREST API before %s", objName, action)
		route.Pipeline = append(route.Pipeline, service)
	}
	service = setupTmplService(ast, tmplPattern, tmplDescription)
	route.Pipeline = append(route.Pipeline, service)
	ast.Routes = append(ast.Routes, route)

	// Setup template submit result
	tmplName = mkName(objName, action, "_response.tmpl")
	tmplPattern = fmt.Sprintf("/%s_%s_response", objName, action)
	tmplDescription = fmt.Sprintf("This is an result template for %s %s", objName, action)
	ast.Templates = append(ast.Templates, &MustacheTemplate{
		Id: templateId,
		Pattern:     tmplPattern,
		Template:    tmplName,
		Description: tmplDescription,
	})
	// Handle submission routing
	routeId = mkName(objName, action, "")
	routeDescription = fmt.Sprintf("Handle form submission for %s %s", objName, action)
	request = fmt.Sprintf("%s /%s_%s", http.MethodPost, objName, action)
	route = &Route{
		Id:          routeId,
		Pattern:     request,
		Description: routeDescription,
		Pipeline:    []*Service{},
	}
	service = setupPostgRESTService(ast, model, action)
	route.Pipeline = append(route.Pipeline, service)
	service = setupTmplService(ast, tmplPattern, tmplDescription)
	route.Pipeline = append(route.Pipeline, service)
	ast.Routes = append(ast.Routes, route)
	ast.isChanged = true
	return nil
}

func setupReadHandling(ast *AST, model *Model, action string) error {
	templateList := ast.GetTemplateIds()
	templateId := fmt.Sprintf("%s_%s", model.Id, action)
	routeList := ast.GetRouteIds()
	routeId := mkName(model.Id, action, "")
	if inList(routeId, routeList) || inList(templateId, templateList) {
		return fmt.Errorf("route or template exists for %s %s", model.Id, action)
	}
	// Setup template for results of read request
	tmplName := mkName(model.Id, action, ".tmpl")
	tmplPattern := fmt.Sprintf("/%s_%s", model.Id, action)
	tmplDescription := fmt.Sprintf("This template handles %s %s", model.Id, action)
	ast.Templates = append(ast.Templates, &MustacheTemplate{
		Id: templateId,
		Pattern:     tmplPattern,
		Template:    tmplName,
		Description: tmplDescription,
	})
	
	// Handle requesting object or list of objects
	routeDescription := fmt.Sprintf("Retrieve object(s) for %s %s", model.Id, action)
	request := fmt.Sprintf("%s /%s_%s", http.MethodPost, model.Id, action)
	route := &Route{
		Id:          routeId,
		Pattern:     request,
		Description: routeDescription,
		Pipeline:    []*Service{},
	}
	service := setupPostgRESTService(ast, model, action)
	route.Pipeline = append(route.Pipeline, service)
	service = setupTmplService(ast, tmplPattern, tmplDescription)
	route.Pipeline = append(route.Pipeline, service)
	ast.Routes = append(ast.Routes, route)
	ast.isChanged = true
	return nil
}

func setupEnvironment(ast *AST, buf *bufio.Reader, out io.Writer, appFName string, skipPrompts bool) {
	var answer string
	if skipPrompts {
		answer = "n"
	} else {
		fmt.Fprintf(out, "Will %s need to import environment variables (y/N)? ", appFName)
		answer = getAnswer(buf, "n", true)
	}
	if answer == "y" {
		if ast.Applications.Environment == nil {
			ast.Applications.Environment = []string{}
		}
		for quit := false; !quit; {
			menu, opt := selectMenuItem(buf, out,
				"Enter menu command and environment name",
				"Menu [a]dd, [r]emove, [q]uit (making changes)",
					ast.Applications.Environment, true, "", "", true)
			if val, ok := getIdFromList(ast.Applications.Environment, opt); ok {
				opt = val
			}
			if len(menu) > 0 {
				menu = menu[0:1]
			}
			switch menu {
				case "a":
				   if opt == "" {
					   fmt.Fprintf(out, "Enter environment name to include: ")
					   opt = getAnswer(buf, "", false)
				   }
				   if opt != "" {
						ast.Applications.Environment = append(ast.Applications.Environment, opt)
						ast.isChanged = true
				   }
				case "r":
					if opt == "" {
						fmt.Fprintf(out, "Enter environment name to remove: ")
					    opt = getAnswer(buf, "", false)
					}
					if opt != "" {
						pos, ok := getItemNoFromList(ast.Applications.Environment, opt)
						if ok {
							ast.Applications.Environment = append(ast.Applications.Environment[:pos], ast.Applications.Environment[(pos+1):]...)
							ast.isChanged = true
						}
					}
				case "q":
					quit = true
				case "":
					// do nothing
				default:
					fmt.Fprint(out, "do not understand %q\n", menu)
			}
		}
	}
}

func setupOptions(ast *AST, buf *bufio.Reader, out io.Writer, appFName string, skipPrompts bool) {
	var answer string
	if skipPrompts {
		answer = "n"
	} else {
		fmt.Fprintf(out, "Will %s provide options to the services (y/N)? ", appFName)
		answer = getAnswer(buf, "n", true)
	}
	if answer == "y" {
		if ast.Applications.Options == nil {
			ast.Applications.Options = map[string]string{}
		}
		for quit := false; !quit; {
			optionsList := getAttributeIds(ast.Applications.Options)
			menu, opt := selectMenuItem(buf, out,
				"Enter menu command and option name",
				"Menu [a]dd, [r]emove, [q]uit (making changes)",
					optionsList, true, "", "", true)
			if val, ok := getIdFromList(optionsList, opt); ok {
				opt = val
			}
			if len(menu) > 0 {
				menu = menu[0:1]
			}
			switch menu {
				case "a":
				   if opt == "" {
					   fmt.Fprintf(out, "Enter option name: ")
					   opt = getAnswer(buf, "", false)
				   }
				   fmt.Fprintf(out, "Enter option value: ")
				   val := getAnswer(buf, "", false)
				   if opt != "" && val != ""{
					    ast.Applications.Options[opt] = val
						ast.isChanged = true
				   }
				case "r":
					if opt == "" {
						fmt.Fprintf(out, "Enter option name to remove: ")
					    opt = getAnswer(buf, "", false)
					}
					if opt != "" {
						delete(ast.Applications.Options, opt)
						ast.isChanged = true
					}
				case "q":
					quit = true
				case "":
					// do nothing
				default:
					fmt.Fprint(out, "do not understand %q\n", menu)
			}
		}
	}
}
