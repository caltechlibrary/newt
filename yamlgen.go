package newt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// setupAppMetadata to configure the project's metadata. Project metadata is used in
// generating the template partials.
func setupAppMetadata(ast *AST, buf *bufio.Reader, out io.Writer, appFName string, skipPrompts bool) {
	var answer string
	if skipPrompts {
		answer = "y"
	} else {
		fmt.Fprintf(out, "Update project metadata for %s (Y/n)? ", appFName)
		answer = getAnswer(buf, "y", true)
	}
	if answer != "y" {
		return
	}
	appName := strings.TrimSuffix(path.Base(appFName), ".yaml")
	year := time.Now().Format("2006")
	if ast.AppMetadata == nil {
		ast.AppMetadata = new(AppMetadata)
	}
	if ast.AppMetadata.AppName == "" {
		ast.AppMetadata.AppName = appName
	}
	if ast.AppMetadata.AppTitle == "" {
		ast.AppMetadata.AppTitle = appName
	}
	if ast.AppMetadata.CopyrightYear == "" {
		ast.AppMetadata.CopyrightYear = year
	}
	if ast.AppMetadata.CSSPath == "" {
		ast.AppMetadata.CSSPath = "/css/site.css"
	}
	fmt.Fprint(out, "Set the project metadata for %s\n", appFName)
	for quit := false; !quit; {
		menuList := []string{
			fmt.Sprintf("Set [an] App Name: %s", ast.AppMetadata.AppName),
			fmt.Sprintf("Set [at] App Title: %s", ast.AppMetadata.AppTitle),
			fmt.Sprintf("Set [cy] Copyright Year: %s", ast.AppMetadata.CopyrightYear),
			fmt.Sprintf("Set [cl] Copyright Link: %s", ast.AppMetadata.CopyrightLink),
			fmt.Sprintf("Set [ct] Copyright Text: %s", ast.AppMetadata.CopyrightText),
			fmt.Sprintf("Set [ll] License Link: %s", ast.AppMetadata.LicenseLink),
			fmt.Sprintf("Set [lt] License Text: %s", ast.AppMetadata.LicenseText),
			fmt.Sprintf("Set [hl] Header Link: %s", ast.AppMetadata.HeaderLink),
			fmt.Sprintf("Set [ht] Header Text: %s", ast.AppMetadata.HeaderText),
			fmt.Sprintf("SEt [gl] Logo Link: %s", ast.AppMetadata.LogoLink),
			fmt.Sprintf("SEt [gt] Logo Link: %s", ast.AppMetadata.LogoText),
			fmt.Sprintf("Set [ca] Contact Address: %s", ast.AppMetadata.ContactAddress),
			fmt.Sprintf("Set [cp] Contact Phone: %s", ast.AppMetadata.ContactPhone),
			fmt.Sprintf("Set [ce] Contact EMail: %s", ast.AppMetadata.ContactEMail),
			fmt.Sprintf("Set [c]  CSS path/url: %s", ast.AppMetadata.CSSPath),
		}
		menu, opt := selectMenuItem(buf, out,
			"Manage Project Metadata",
			"Type menu letter(s) and press enter to modify or press enter when done",
			menuList, false, "", "", true)
		if len(menu) > 0 {
			menu = strings.TrimSpace(menu)
		}
		switch menu {
			case "an":
				if opt == "" {
					fmt.Fprintf(out, "Enter App Name: ")
					opt = getAnswer(buf, appName, true)
				}
				if opt != ast.AppMetadata.AppName {
					ast.AppMetadata.AppName = opt
					ast.isChanged = true
				}
			case "at":
				if opt == "" {
					fmt.Fprintf(out, "Enter App Title: ")
					opt = getAnswer(buf, "", true)
				}
				if opt != ast.AppMetadata.AppTitle {
					ast.AppMetadata.AppTitle = opt
					ast.isChanged = true
				}
			case "cy":
				if opt == "" {
					fmt.Fprintf(out, "Enter Copyright Year: ")
					opt = getAnswer(buf, year, true)
				}
				if opt != ast.AppMetadata.CopyrightYear {
					ast.AppMetadata.CopyrightYear = opt
					ast.isChanged = true
				}
			case "cl":
				if opt == "" {
					fmt.Fprintf(out, "Enter Copyright Link: ")
					opt = getAnswer(buf, "", true)
				}
				if opt != ast.AppMetadata.CopyrightLink {
					ast.AppMetadata.CopyrightLink = opt
					ast.isChanged = true
				}
			case "ct":
				if opt == "" {
					fmt.Fprintf(out, "Enter Copyright Text: ")
					opt = getAnswer(buf, "", true)
				}
				if opt != ast.AppMetadata.CopyrightText {
					ast.AppMetadata.CopyrightText = opt
					ast.isChanged = true
				}
			case "ll":
				if opt == "" {
					fmt.Fprintf(out, "Enter License Link: ")
					opt = getAnswer(buf, "", true)
				}
				if opt != ast.AppMetadata.LicenseLink {
					ast.AppMetadata.LicenseLink = opt
					ast.isChanged = true
				}
			case "lt":
				if opt == "" {
					fmt.Fprintf(out, "Enter License Text: ")
					opt = getAnswer(buf, "", true)
				}
				if opt != ast.AppMetadata.LicenseText {
					ast.AppMetadata.LicenseText = opt
					ast.isChanged = true
				}
			case "gl":
				if opt == "" {
					fmt.Fprintf(out, "Enter Logo Link: ")
					opt = getAnswer(buf, "", true)
				}
				if opt != ast.AppMetadata.LogoLink {
					ast.AppMetadata.LogoLink = opt
					ast.isChanged = true
				}
			case "gt":
				if opt == "" {
					fmt.Fprintf(out, "Enter Logo Text: ")
					opt = getAnswer(buf, "", true)
				}
				if opt != ast.AppMetadata.LogoText {
					ast.AppMetadata.LogoText = opt
					ast.isChanged = true
				}
			case "hl":
				if opt == "" {
					fmt.Fprintf(out, "Enter Header Link: ")
					opt = getAnswer(buf, "", true)
				}
				if opt != ast.AppMetadata.HeaderLink {
					ast.AppMetadata.HeaderLink = opt
					ast.isChanged = true
				}
			case "ht":
				if opt == "" {
					fmt.Fprintf(out, "Enter Header Text: ")
					opt = getAnswer(buf, "", true)
				}
				if opt != ast.AppMetadata.HeaderText {
					ast.AppMetadata.HeaderText = opt
					ast.isChanged = true
				}
			case "ca":
				if opt == "" {
					fmt.Fprintf(out, "Enter Contact Address: ")
					opt = getAnswer(buf, "", true)
				}
				if opt != ast.AppMetadata.ContactAddress {
					ast.AppMetadata.ContactAddress = opt
					ast.isChanged = true
				}
			case "cp":
				if opt == "" {
					fmt.Fprintf(out, "Enter Contact Phone: ")
					opt = getAnswer(buf, "", true)
				}
				if opt != ast.AppMetadata.ContactPhone {
					ast.AppMetadata.ContactPhone = opt
					ast.isChanged = true
				}
			case "ce":
				if opt == "" {
					fmt.Fprintf(out, "Enter Contact EMail: ")
					opt = getAnswer(buf, "", true)
				}
				if opt != ast.AppMetadata.ContactEMail {
					ast.AppMetadata.ContactEMail = opt
					ast.isChanged = true
				}
			case "c":
				if opt == "" {
					fmt.Fprintf(out, "Enter CSS path or url: ")
					opt = getAnswer(buf, "", true)
				}
				if opt != ast.AppMetadata.CSSPath {
					ast.AppMetadata.CSSPath = opt
					ast.isChanged = true
				}
			case "q":
				quit = true
			case "":
				quit = true
			default:
				fmt.Fprintf(out, "failed to understand %q\n", menu)
		}

	}
}

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
			ast.Applications = NewApplications()
		}
		if ast.Applications.Router == nil {
			ast.Applications.Router = NewApplication()
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
		for quit := false; !quit; {
			menuList := []string{
				fmt.Sprintf("Set [p]ort: %d", ast.Applications.Router.Port),
				fmt.Sprintf("Set [h]tdocs: %s", ast.Applications.Router.Htdocs),
			}
			menu, opt := selectMenuItem(buf, out,
				"Manage Newt Router Settings",
				"Type menu letter and press enter to modify or press enter when done",
				menuList, false, "", "", true)
			if len(menu) > 0 {
				menu = menu[0:1]
			}
			switch menu {
			case "p":
				if opt == "" {
					fmt.Fprintf(out, "Enter Port number: ")
					opt = getAnswer(buf, strconv.Itoa(ROUTER_PORT), true)
				}
				port, err := strconv.Atoi(opt)
				if err != nil {
					fmt.Fprintf(out, "ERROR: port number post be an integer, got %q\n", opt)
				} else {
					ast.Applications.Router.Port = port
					ast.isChanged = true
				}
			case "h":
				if opt == "" {
					fmt.Fprintf(out, "Enter htdocs value: ")
					opt = getAnswer(buf, "", false)
				}
				if opt != ast.Applications.Router.Htdocs {
					ast.Applications.Router.Htdocs = strings.TrimSpace(opt)
					ast.isChanged = true
				}
			case "q":
				quit = true
			case "":
				quit = true
			default:
				fmt.Fprintf(out, "failed to understand %q\n", menu)
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
	fmt.Fprintf(out, "Configuring PostgREST\n")
	if ast.Applications == nil {
		ast.Applications = NewApplications()
	}
	if ast.Applications.PostgREST == nil {
		ast.Applications.PostgREST = NewApplication()
	}
	if ast.Applications.PostgREST.Port == 0 {
		ast.Applications.PostgREST.Port = POSTGREST_PORT
		ast.Applications.PostgREST.AppPath = "postgrest"
		ast.Applications.PostgREST.ConfPath = "postgrest.conf"
	}
	for quit := false; !quit; {
		menuList := []string{
			fmt.Sprintf("Set [p]ort: %d", ast.Applications.PostgREST.Port),
			fmt.Sprintf("Set [a]pp path: %q", ast.Applications.PostgREST.AppPath),
			fmt.Sprintf("Set [c]onf path: %q", ast.Applications.PostgREST.ConfPath),
		}
		menu, opt := selectMenuItem(buf, out,
			"Manage PostgREST Settings",
			"Type menu letter and press enter to modify or press enter when done",
			menuList, false, "", "", true)
		if len(menu) > 0 {
			menu = menu[0:1]
		}
		switch menu {
		case "p":
			if opt == "" {
				fmt.Fprintf(out, "Enter Postgres port number: ")
				opt = getAnswer(buf, strconv.Itoa(POSTGREST_PORT), true)
			}
			port, err := strconv.Atoi(opt)
			if err != nil {
				fmt.Fprintf(out, "ERROR: port number must be an intereger, got %q\n", opt)
			} else {
				ast.Applications.PostgREST.Port = port
				ast.isChanged = true
			}
		case "a":
			if opt == "" {
				fmt.Fprintf(out, "Enter the path to PostgREST application (an empty path is OK): ")
				opt = getAnswer(buf, "", false)
			}
			if ast.Applications.PostgREST.AppPath != opt {
				ast.Applications.PostgREST.AppPath = strings.TrimSpace(opt)
				ast.isChanged = true
			}
		case "c":
			if opt == "" {
				fmt.Fprintf(out, "Enter the path to PostgREST configuration (an empty path is OK): ")
				opt = getAnswer(buf, "", false)
			}
			if ast.Applications.PostgREST.ConfPath != opt {
				ast.Applications.PostgREST.ConfPath = strings.TrimSpace(opt)
				ast.isChanged = true
			}
		case "q":
			quit = true
		case "":
			quit = true
		default:
			fmt.Fprintf(out, "failed to understand request, %q\n", menu)
		}
	}
}

// setupPostgresAndPostgREST prompt to configure Postgres
func setupPostgresAndPostgREST(ast *AST, buf *bufio.Reader, out io.Writer, appFName string, skipPrompts bool) {
	var answer string
	if skipPrompts {
		answer = "y"
	} else {
		fmt.Fprintf(out, "Will %s use Postgres and PostgREST (Y/n)? ", appFName)
		answer = getAnswer(buf, "y", true)
	}
	if answer == "y" {
		if ast.Applications == nil {
			ast.Applications = NewApplications()
		}
		if ast.Applications.Postgres == nil {
			ast.Applications.Postgres = NewApplication()
		}
		if ast.Applications.Postgres.Port == 0 {
			ast.Applications.Postgres.Port = POSTGRES_PORT
		}
		if ast.Applications.Postgres.Namespace == "" {
			ast.Applications.Postgres.Namespace = fNameToNamespace(appFName)
		}
		if ast.Applications.Postgres.DSN == "" {
			ast.Applications.Postgres.DSN = fmt.Sprintf("postgres://{PGUSER}:{PGPASSWORD}@localhost:%d/%s", ast.Applications.Postgres.Port, appFName)
			// Now we need to make sure we allow PGUSER and PGPASSWORD to pass through in the environment
			if ast.Applications.Environment == nil {
				ast.Applications.Environment = []string{}
			}
			ast.Applications.Environment = append(ast.Applications.Environment, "PGUSER", "PGPASSWORD")
		}
		for quit := false; !quit; {
			menuList := []string{
				fmt.Sprintf("Set [p]ort: %d", ast.Applications.Postgres.Port),
				fmt.Sprintf("Set [d]sn (data source name): %s", ast.Applications.Postgres.DSN),
				fmt.Sprintf("Set [n]amespace: %s", ast.Applications.Postgres.Namespace),
			}
			menu, opt := selectMenuItem(buf, out,
				"Manage Postgres Settings",
				"Type menu letter and press enter to modify or press enter when done",
				menuList, false, "", "", true)
			if len(menu) > 0 {
				menu = menu[0:1]
			}
			switch menu {
			case "n":
				if opt == "" {
					fmt.Fprintf(out, "Enter namespace value: ")
					opt = getAnswer(buf, appFName, true)
				}
				if opt != ast.Applications.Postgres.Namespace {
					ast.Applications.Postgres.Namespace = strings.TrimSpace(opt)
					ast.isChanged = true
				}
			case "p":
				if opt == "" {
					fmt.Fprintf(out, "Enter Postgres port: ")
					opt = getAnswer(buf, strconv.Itoa(POSTGRES_PORT), true)
				}
				port, err := strconv.Atoi(opt)
				if err != nil {
					fmt.Fprintf(out, "ERROR: port number must be an intereger, got %q\n", opt)
				} else {
					ast.Applications.Postgres.Port = port
					ast.isChanged = true
				}
			case "d":
				if opt == "" {
					fmt.Fprintf(out, "Enter DSN in uri form: ")
					opt = getAnswer(buf, "", false)
				}
				if ast.Applications.Postgres.DSN != opt {
					ast.Applications.Postgres.DSN = strings.TrimSpace(opt)
					ast.isChanged = true
				}
			case "q":
				quit = true
			case "":
				quit = true
			default:
				fmt.Fprintf(out, "failed to understand request, %q\n", menu)
			}
		}
		setupPostgREST(ast, buf, out, appFName, skipPrompts)
	} else {
		if ast.Applications != nil {
			ast.Applications.Postgres = nil
		}
	}
}

func setupTemplateEngine(ast *AST, buf *bufio.Reader, out io.Writer, appFName string, skipPrompts bool) {
	var answer string
	if skipPrompts {
		answer = "y"
	} else {
		fmt.Fprintf(out, "Will %s use Newt's template engine (Y/n)? ", appFName)
		answer = getAnswer(buf, "y", true)
	}
	if answer == "y" {
		if ast.Applications == nil {
			ast.Applications = NewApplications()
		}
		if ast.Applications.TemplateEngine == nil {
			ast.Applications.TemplateEngine = NewApplication()
		}
		if ast.Applications.TemplateEngine.Port == 0 {
			ast.Applications.TemplateEngine.Port = TEMPLATE_ENGINE_PORT
		}
		if ast.Applications.TemplateEngine.BaseDir == "" {
			ast.Applications.TemplateEngine.BaseDir = TEMPLATE_ENGINE_BASE_DIR
		}
		if ast.Applications.TemplateEngine.ExtName == "" {
			ast.Applications.TemplateEngine.ExtName = TEMPLATE_ENGINE_EXT_NAME
		}
		if ast.Applications.TemplateEngine.PartialsDir == "" {
			ast.Applications.TemplateEngine.PartialsDir = TEMPLATE_ENGINE_PARTIALS_DIR
		}
		//FIXME: If there are models then templates will need to be updates even when it is NOT nil.
		// When the model list changes then the related templates should change to.
		// A scan of the template routes for removed models needs to happen when the model is "removed" by the modeler.
		if ast.Templates == nil {
			ast.Templates = []*Template{}
			if err := setupTemplateHandlers(ast); err != nil {
				fmt.Fprintf(out, "WARNINGS: %s\n", err)
			}
		}
		for quit := false; !quit; {
			menuList := []string{
				fmt.Sprintf("Set [p]ort: %d", ast.Applications.TemplateEngine.Port),
				fmt.Sprintf("Set [b]ase directory: %s", ast.Applications.TemplateEngine.BaseDir),
				fmt.Sprintf("Set [f]ile extention: %s", ast.Applications.TemplateEngine.ExtName),
				fmt.Sprintf("Set [P]artials sub-directory: %s", ast.Applications.TemplateEngine.PartialsDir),
			}
			// FIXME: You show the current template list here..
			menu, opt := selectMenuItem(buf, out,
				"Manage Newt's template engine Settings",
				"Type menu letter and press enter to modify or press enter when done",
				menuList, false, "", "", true)
			if len(menu) > 0 {
				menu = menu[0:1]
			}
			switch menu {
			case "p":
				if opt == "" {
					fmt.Fprintf(out, "Enter port number: ")
					opt = getAnswer(buf, strconv.Itoa(TEMPLATE_ENGINE_PORT), true)
				}
				port, err := strconv.Atoi(opt)
				if err != nil {
					fmt.Fprintf(out, "ERROR: port number must be an intereger, got %q\n", opt)
				} else {
					ast.Applications.TemplateEngine.Port = port
					ast.isChanged = true
				}
			case "b":
				if opt == "" {
					fmt.Fprintf(out, "Enter base directory: ")
					opt = getAnswer(buf, TEMPLATE_ENGINE_BASE_DIR, true)
				}
				if opt != "" {
					ast.Applications.TemplateEngine.BaseDir = opt
					ast.isChanged = true
				}
			case "f":
				if opt == "" {
					fmt.Fprintf(out, "Enter file extention: ")
					opt = getAnswer(buf, TEMPLATE_ENGINE_EXT_NAME, true)
				}
				if opt != "" {
					ast.Applications.TemplateEngine.ExtName = opt
					ast.isChanged = true
				}
			case "P":
				if opt == "" {
					fmt.Fprintf(out, "Enter partials sub directory: ")
					opt = getAnswer(buf, TEMPLATE_ENGINE_PARTIALS_DIR, true)
				}
				if opt != "" {
					ast.Applications.TemplateEngine.PartialsDir = opt
					ast.isChanged = true
				}
			case "q":
				quit = true
			case "":
				quit = true
			default:
				fmt.Fprintf(out, "failed to understand request, %q\n", menu)
			}
		}
	} else {
		if ast.Applications != nil {
			ast.Applications.TemplateEngine = nil
		}
	}
}

func setupTemplateHandlers(ast *AST) error {
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
	if ast.Applications != nil && ast.Applications.TemplateEngine != nil {
		port = ast.Applications.TemplateEngine.Port
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
	if inList(routeId, routeList) {
		if err := ast.RemoveRouteById(routeId); err != nil {
			return err
		}
	}
	if inList(templateId, templateList) {
		if err := ast.RemoveTemplateById(templateId); err != nil {
			return err
		}
	}

	tSuffix := "_form.hbs"
	tmplName := mkName(objName, action, tSuffix)
	tmplPattern := fmt.Sprintf("/%s_%s", objName, action)
	tmplDescription := fmt.Sprintf("Display a %s for %s", objName, action)
	ast.Templates = append(ast.Templates, &Template{
		Id:          templateId,
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
	tmplName = mkName(objName, action, "_response.hbs")
	tmplPattern = fmt.Sprintf("/%s_%s_response", objName, action)
	tmplDescription = fmt.Sprintf("This is an result template for %s %s", objName, action)
	ast.Templates = append(ast.Templates, &Template{
		Id:          templateId,
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
	if inList(routeId, routeList) {
		if err := ast.RemoveRouteById(routeId); err != nil {
			return err
		}
	}
	if inList(templateId, templateList) {
		if err := ast.RemoveTemplateById(templateId); err != nil {
			return err
		}
	}
	// Setup template for results of read request
	tmplName := mkName(model.Id, action, ".hbs")
	tmplPattern := fmt.Sprintf("/%s_%s", model.Id, action)
	tmplDescription := fmt.Sprintf("This template handles %s %s", model.Id, action)
	ast.Templates = append(ast.Templates, &Template{
		Id:          templateId,
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
				fmt.Sprintf("Manage Environment availability in %s", appFName),
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
				quit = true
			default:
				fmt.Fprintf(out, "do not understand %q\n", menu)
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
			ast.Applications.Options = make(map[string]interface{})
		}
		for quit := false; !quit; {
			optionsList := []string{}
			for k, v := range ast.Applications.Options {
				optionsList = append(optionsList, fmt.Sprintf("%s -> %q", k, v))
			}
			menu, opt := selectMenuItem(buf, out,
				"Enter menu command and option name",
				"Menu [a]dd, [r]emove, [q]uit (making changes)",
				optionsList, false, "", "", true)
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
				if opt != "" && val != "" {
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
				quit = true
			default:
				fmt.Fprintf(out, "do not understand %q\n", menu)
			}
		}
	}
}
