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
		ast.isChanged = true
	}
	if ast.AppMetadata.AppName == "" {
		ast.AppMetadata.AppName = appName
		ast.isChanged = true
	}
	if ast.AppMetadata.AppTitle == "" {
		ast.AppMetadata.AppTitle = appName
		ast.isChanged = true
	}
	if ast.AppMetadata.CopyrightYear == "" {
		ast.AppMetadata.CopyrightYear = year
		ast.isChanged = true
	}
	if ast.AppMetadata.CSSPath == "" {
		ast.AppMetadata.CSSPath = "/css/site.css"
		ast.isChanged = true
	}
	if skipPrompts {
		return
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
		if ast.Services == nil {
			ast.Services = NewServices()
			ast.isChanged = true
		}
		router := ast.GetService("router")
		if router == nil {
			router = &Service{
				AppName: "router",
				Port: 0,
				Htdocs: "",
			}
			ast.Services = append(ast.Services, router)
			ast.isChanged = true
		}
		// NOTE: If port is zero, we haven't configure the router.
		if router.Port == 0 {
			router.Port = ROUTER_PORT
			if info, err := os.Stat("htdocs"); err == nil && info.IsDir() {
				router.Htdocs = "htdocs"
				ast.isChanged = true
			} else {
				router.Htdocs = ""
				fmt.Fprintf(out, "WARNING: Skipping the htdocs directory for router")
			}
		}
		if skipPrompts {
			return
		}
		for quit := false; !quit; {
			menuList := []string{
				fmt.Sprintf("Set [p]ort: %d", router.Port),
				fmt.Sprintf("Set [h]tdocs: %s", router.Htdocs),
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
					router.Port = port
					ast.isChanged = true
				}
			case "h":
				if opt == "" {
					fmt.Fprintf(out, "Enter htdocs value: ")
					opt = getAnswer(buf, "", false)
				}
				if opt != router.Htdocs {
					router.Htdocs = strings.TrimSpace(opt)
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
		ast.RemoveService("router")
	}
}

// setupPostgREST prompt to configure PostgREST
func setupPostgREST(ast *AST, buf *bufio.Reader, out io.Writer, appFName string, skipPrompts bool) {
	fmt.Fprintf(out, "Configuring PostgREST\n")
	if ast.Services == nil {
		ast.Services = NewServices()
		ast.isChanged = true
	}
	postgREST := ast.GetService("postgrest")
	if postgREST == nil {
		postgREST = &Service{
			AppName: "postgrest",
			Port: 0,
		}
		ast.Services = append(ast.Services, postgREST)
		ast.isChanged = true
	}
	if postgREST.Port == 0 {
		postgREST.Port = POSTGREST_PORT
		postgREST.AppPath = "postgrest"
		postgREST.ConfPath = "postgrest.conf"
		ast.isChanged = true
	}
	if skipPrompts {
		return
	}
	for quit := false; !quit; {
		menuList := []string{
			fmt.Sprintf("Set [p]ort: %d", postgREST.Port),
			fmt.Sprintf("Set [a]pp path: %q", postgREST.AppPath),
			fmt.Sprintf("Set [c]onf path: %q", postgREST.ConfPath),
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
				postgREST.Port = port
				ast.isChanged = true
			}
		case "a":
			if opt == "" {
				fmt.Fprintf(out, "Enter the path to PostgREST application (an empty path is OK): ")
				opt = getAnswer(buf, "", false)
			}
			if postgREST.AppPath != opt {
				postgREST.AppPath = strings.TrimSpace(opt)
				ast.isChanged = true
			}
		case "c":
			if opt == "" {
				fmt.Fprintf(out, "Enter the path to PostgREST configuration (an empty path is OK): ")
				opt = getAnswer(buf, "", false)
			}
			if postgREST.ConfPath != opt {
				postgREST.ConfPath = strings.TrimSpace(opt)
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
		if ast.Services == nil {
			ast.Services = NewServices()
			ast.isChanged = true
		}
		postgres := ast.GetService("postgres")
		if postgres == nil {
			postgres = &Service{
				AppName: "postgres",
				Port: 9,
				Namespace: fNameToNamespace(appFName),
			}
			ast.Services = append(ast.Services, postgres)
			ast.isChanged = true
		}
		if postgres.Port == 0 {
			postgres.Port = POSTGRES_PORT
			ast.isChanged = true
		}
		if postgres.Namespace == "" {
			postgres.Namespace = fNameToNamespace(appFName)
			ast.isChanged = true
		}
		if postgres.DSN == "" {
			postgres.DSN = fmt.Sprintf("postgres://{PGUSER}:{PGPASSWORD}@localhost:%d/%s", postgres.Port, appFName)
			// Now we need to make sure we allow PGUSER and PGPASSWORD to pass through in the environment
			if postgres.Environment == nil {
				postgres.Environment = []string{}
			}
			postgres.Environment = append(postgres.Environment, "PGUSER", "PGPASSWORD")
			ast.isChanged = true
		}
		if skipPrompts {
			return
		}
		for quit := false; !quit; {
			menuList := []string{
				fmt.Sprintf("Set [p]ort: %d", postgres.Port),
				fmt.Sprintf("Set [d]sn (data source name): %s", postgres.DSN),
				fmt.Sprintf("Set [n]amespace: %s", postgres.Namespace),
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
				if opt != postgres.Namespace {
					postgres.Namespace = strings.TrimSpace(opt)
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
					postgres.Port = port
					ast.isChanged = true
				}
			case "d":
				if opt == "" {
					fmt.Fprintf(out, "Enter DSN in uri form: ")
					opt = getAnswer(buf, "", false)
				}
				if postgres.DSN != opt {
					postgres.DSN = strings.TrimSpace(opt)
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
		ast.RemoveService("postgres")
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
		if ast.Services == nil {
			ast.Services = NewServices()
			ast.isChanged = true
		}
		templateEngine := ast.GetService("template_engine")
		if templateEngine == nil {
			templateEngine = &Service{
				AppName: "template_engine",
				Port: 0,
			}
			ast.Services = append(ast.Services, templateEngine)
			ast.isChanged = true
		}
		if templateEngine.Port == 0 {
			templateEngine.Port = TEMPLATE_ENGINE_PORT
			ast.isChanged = true
		}
		if templateEngine.BaseDir == "" {
			templateEngine.BaseDir = TEMPLATE_ENGINE_BASE_DIR
			ast.isChanged = true
		}
		if templateEngine.ExtName == "" {
			templateEngine.ExtName = TEMPLATE_ENGINE_EXT_NAME
			ast.isChanged = true
		}
		if templateEngine.PartialsDir == "" {
			templateEngine.PartialsDir = TEMPLATE_ENGINE_PARTIALS_DIR
			ast.isChanged = true
		}
		//FIXME: If there are models then templates will need to be updates even when it is NOT nil.
		// When the model list changes then the related templates should change to.
		// A scan of the template routes for removed models needs to happen when the model is "removed" by the modeler.
		if ast.Templates == nil {
			ast.Templates = []*Template{}
			if err := setupTemplateHandlers(ast); err != nil {
				fmt.Fprintf(out, "WARNINGS: %s\n", err)
			}
			ast.isChanged = true
		}
		if skipPrompts {
			return
		}
		for quit := false; !quit; {
			menuList := []string{
				fmt.Sprintf("Set [p]ort: %d", templateEngine.Port),
				fmt.Sprintf("Set [b]ase directory: %s", templateEngine.BaseDir),
				fmt.Sprintf("Set [f]ile extention: %s", templateEngine.ExtName),
				fmt.Sprintf("Set [P]artials sub-directory: %s", templateEngine.PartialsDir),
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
					templateEngine.Port = port
					ast.isChanged = true
				}
			case "b":
				if opt == "" {
					fmt.Fprintf(out, "Enter base directory: ")
					opt = getAnswer(buf, TEMPLATE_ENGINE_BASE_DIR, true)
				}
				if opt != "" {
					templateEngine.BaseDir = opt
					ast.isChanged = true
				}
			case "f":
				if opt == "" {
					fmt.Fprintf(out, "Enter file extention: ")
					opt = getAnswer(buf, TEMPLATE_ENGINE_EXT_NAME, true)
				}
				if opt != "" {
					templateEngine.ExtName = opt
					ast.isChanged = true
				}
			case "P":
				if opt == "" {
					fmt.Fprintf(out, "Enter partials sub directory: ")
					opt = getAnswer(buf, TEMPLATE_ENGINE_PARTIALS_DIR, true)
				}
				if opt != "" {
					templateEngine.PartialsDir = opt
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
		ast.RemoveService("template_engine")
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

// setupPostgRESTService creates a RouteService object for interacting with PostgREST
func setupPostgRESTService(ast *AST, model *Model, action string) *RouteService {
	var (
		identifier         string
		identifierSuffix   string
		description string
		method      string
		port        int
	)
	objName := model.Id
	element, ok := model.GetModelIdentifier()
	if ok {
		identifier = fmt.Sprintf("{%s}", element.Id)
	} else {
		identifier = "{identifier}"
	}
	description = fmt.Sprintf("Access PostgREST API for %s %s", objName, action)
	postgREST := ast.GetService("postgrest")
	if ast.Services != nil && postgREST != nil {
		port = postgREST.Port
	} else {
		port = 3000
	}
	switch action {
	case "create":
		// create action doesn't take an identifier
		method = http.MethodPost
	case "read":
		method = http.MethodGet
		identifierSuffix = "/" + identifier
	case "update":
		method = http.MethodPut
		identifierSuffix = "/" + identifier
	case "delete":
		method = http.MethodDelete
		identifierSuffix = "/" + identifier
	default:
		// list action doesn't take an identifier
		method = http.MethodGet
	}
	return &RouteService{
		Service:     fmt.Sprintf("%s http://localhost:%d/rpc/%s_%s%s", method, port, objName, action, identifierSuffix),
		Description: description,
	}
}

// setupTemplService creates a RouteService object to process with a template
func setupTmplService(ast *AST, tmplPattern string, description string) *RouteService {
	var port int
	templateEngine := ast.GetService("template_engine")
	if ast.Services != nil && templateEngine != nil {
		port = templateEngine.Port
	} else {
		port = 8011
	}
	serviceURL := fmt.Sprintf("POST http://localhost:%d%s", port, tmplPattern)
	return &RouteService{
		Service:     serviceURL,
		Description: description,
	}
}

// setupWebFormHandling generates the routes and template handling for retrieving and submitting
// webforms for "create", "update" or "delete".
func setupWebFormHandling(ast *AST, model *Model, action string) error {
	var (
		identifier        string
		pathSuffix string
		service    *RouteService
	)
	objName := model.Id
	element, ok := model.GetModelIdentifier()
	if ok {
		identifier = fmt.Sprintf("{%s}", element.Id)
	} else {
		identifier = "{identifier}"
	}
	if action == "update" || action == "delete" {
		pathSuffix = "/" + identifier
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
		Pipeline:    []*RouteService{},
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

	//FIXME: somewhere here the update route isn't getting the {indentifier} set for it's path

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
		Pipeline:    []*RouteService{},
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
		Pipeline:    []*RouteService{},
	}
	service := setupPostgRESTService(ast, model, action)
	route.Pipeline = append(route.Pipeline, service)
	service = setupTmplService(ast, tmplPattern, tmplDescription)
	route.Pipeline = append(route.Pipeline, service)
	ast.Routes = append(ast.Routes, route)
	ast.isChanged = true
	return nil
}

//FIXME: setupEnvironment is now set at the runtime app level ...
func setupEnvironment(ast *AST, buf *bufio.Reader, out io.Writer, appName string, appFName string, skipPrompts bool) {
	var answer string
	if skipPrompts {
		answer = "n"
	} else {
		fmt.Fprintf(out, "Will %s need to import environment variables (y/N)? ", appFName)
		answer = getAnswer(buf, "n", true)
	}
	if answer == "y" {
		app := ast.GetService(appName)
		if app == nil {
			app := &Service{
				AppName: appName,
			}
			ast.Services = append(ast.Services, app)
			ast.isChanged = true
		}
		if app.Environment == nil {
			app.Environment = []string{}
		}
		for quit := false; !quit; {
			menu, opt := selectMenuItem(buf, out,
				fmt.Sprintf("Manage Environment availability in %s", appFName),
				"Menu [a]dd, [r]emove or press enter when done",
				app.Environment, true, "", "", true)
			if val, ok := getIdFromList(app.Environment, opt); ok {
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
					app.Environment = append(app.Environment, opt)
					ast.isChanged = true
				}
			case "r":
				if opt == "" {
					fmt.Fprintf(out, "Enter environment name to remove: ")
					opt = getAnswer(buf, "", false)
				}
				if opt != "" {
					pos, ok := getItemNoFromList(app.Environment, opt)
					if ok {
						app.Environment = append(app.Environment[:pos], app.Environment[(pos+1):]...)
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

//FIXME: Options are now set per application
func setupOptions(ast *AST, buf *bufio.Reader, out io.Writer, appName string, appFName string, skipPrompts bool) {
	var answer string
	if skipPrompts {
		answer = "n"
	} else {
		fmt.Fprintf(out, "Will %s provide options to the services (y/N)? ", appFName)
		answer = getAnswer(buf, "n", true)
	}
	if answer == "y" {
		app := ast.GetService(appName)
		if app == nil {
			app = &Service{
				AppName: appName,
			}
			ast.Services = append(ast.Services, app)
			ast.isChanged = true
		}
		if app.Options == nil {
			app.Options = make(map[string]interface{})
		}
		for quit := false; !quit; {
			optionsList := []string{}
			for k, v := range app.Options {
				optionsList = append(optionsList, fmt.Sprintf("%s -> %q", k, v))
			}
			menu, opt := selectMenuItem(buf, out,
				"Enter menu command and option name",
				"Menu [a]dd, [r]emove or press enter when done",
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
					app.Options[opt] = val
					ast.isChanged = true
				}
			case "r":
				if opt == "" {
					fmt.Fprintf(out, "Enter option name to remove: ")
					opt = getAnswer(buf, "", false)
				}
				if opt != "" {
					delete(app.Options, opt)
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
