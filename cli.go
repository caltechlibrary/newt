/**
 * cli.go an implements runners for the cli of the Newt Project.
 *
 * @author R. S. Doiel
 */
package newt

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"time"

	// 3rd Party Templates
	"github.com/cbroglie/mustache"
	"gopkg.in/yaml.v3"
)

const (
	// These constants are used for exit code. FIXME: look up the POSIX recommendation on exit
	// codes and adopt those.
	OK = iota
	CONFIG

	// General failure of a command or service
	INIT_FAIL
	CHECK_FAIL
	GENERATOR_FAIL
	ROUTER_FAIL
	MUSTACHE_FAIL
	NEWT_FAIL
	SWS_FAIL
	POSTGREST_FAIL

	// Internal service failures
	RESOLVE
	HANDLER
	SERVER_ERROR
	UNSUPPORTED_ACTION
	DATA_ERROR
	READ_ERROR
	DECODE_ERROR
	TEMPLATE_ERROR

	// Default service settings
	ROUTER_PORT      = 8010
	MUSTACHE_PORT    = 8011
	MUSTACHE_TIMEOUT = 3 * time.Second
	SWS_PORT         = 8000
	SWS_HTDOCS       = "."
	POSTGREST_PORT   = 3000
)

// RunNewtGenerator is a runner for generating SQL and templates from our Newt YAML file.
func RunNewtGenerator(in io.Reader, out io.Writer, eout io.Writer, args []string) int {
	//appName := "Newt Generator"
	fName, generatorName, action, modelName := "", "", "", ""
	if len(args) > 0 {
		fName = args[0]
	} else {
		fmt.Fprintf(eout, "missing YAML configuration file\n")
		return CONFIG
	}
	if len(args) > 1 {
		generatorName = args[1]
	} else {
		fmt.Fprintf(eout, "missing generator name, e.g. postgres, postgrest, mustache\n")
		return CONFIG
	}
	if len(args) > 2 {
		action = args[2]
	}
	if len(args) > 3 {
		modelName = args[3]
	}
	cfg, err := LoadConfig(fName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
	}
	if cfg.Applications == nil || cfg.Applications.NewtGenerator == nil {
		fmt.Fprintf(eout, "missing newtgenerator configuration, aborting\n")
		return CONFIG
	}
	generator, err := NewGenerator(cfg)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
	}
	generator.out = out
	generator.eout = eout
	if err := generator.Generate(generatorName, action, modelName); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return GENERATOR_FAIL
	}
	return OK
}

// RunNewtMustache is a runner for a Mustache redner engine service based on the Pandoc server API.
func RunNewtMustache(in io.Reader, out io.Writer, eout io.Writer, args []string, port int, timeout int, verbose bool) int {
	appName := "Newt Mustache"
	// Configure the template bundler webservice
	fName := ""
	if len(args) > 0 {
		fName = args[0]
	} else {
		fmt.Fprintf(eout, "missing Newt YAML configuration\n")
		return CONFIG
	}
	// Load the Newt YAML syntax file holding the configuration
	// and make sure it conforms.
	cfg, err := LoadConfig(fName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
	}
	// Instantiate the specific application with the filename and Config object
	mt, err := NewNewtMustache(cfg)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
	}
	// If port is not set in the config, set it to the default port.
	if mt.Port == 0 {
		mt.Port = MUSTACHE_PORT
	}
	// See if we have a command line option for port to process.
	if port != 0 {
		mt.Port = port
	}
	if verbose {
		fmt.Fprintf(out, "port set to %d\n", mt.Port)
	}

	if timeout != 0 {
		mt.Timeout = time.Duration(timeout) * time.Second
	}
	if mt.Timeout == 0 {
		mt.Timeout = MUSTACHE_TIMEOUT
	}
	if len(mt.Templates) == 0 {
		fmt.Fprintf(eout, "no templates in configuration\n")
		return CONFIG
	}
	// Create mux for http service
	// Resolve partial templates and build handlers
	for _, tmpl := range mt.Templates {
		if verbose {
			tmpl.Debug = true
		}
		if err := tmpl.ResolveTemplate(); err != nil {
			fmt.Fprintf(eout, "%s failed to resolve, %s\n", tmpl.Template, err)
			return RESOLVE
		}
		if err := tmpl.ResolvePath(); err != nil {
			fmt.Fprintf(eout, "failed to build handler for %q, %s\n", tmpl.Pattern, err)
			return HANDLER
		}
	}
	// Launch web service
	fmt.Printf("starting %s listening on port :%d (press Ctrl-c to exit)\n", appName, mt.Port)
	if err := mt.ListenAndServe(); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return MUSTACHE_FAIL
	}
	return OK
}

// RunNewtRouter is a runner for Newt data router and static file service
func RunNewtRouter(in io.Reader, out io.Writer, eout io.Writer, args []string, dryRun bool, port int, htdocs string, verbose bool) int {
	appName := "Newt Router"
	// You can run Newt Router with just an htdocs directory. If so you don't require a config file.
	var err error
	cfg := &Config{
		Applications: &Applications{
			NewtRouter: &Application{},
		},
	}
	router := &NewtRouter{}
	fName := ""
	if htdocs == "" || len(args) > 0 {
		if len(args) > 0 {
			fName = args[0]
		}
		cfg, err = LoadConfig(fName)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return CONFIG
		}
		// Finally Instantiate the router from fName and Config object
		router, err = NewNewtRouter(cfg)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return CONFIG
		}
	}
	if router.Port == 0 {
		router.Port = ROUTER_PORT
	}
	if port != 0 {
		router.Port = port
	}
	if htdocs != "" {
		router.Htdocs = htdocs
	}

	// Are we ready to run service?
	if router.Routes == nil && router.Htdocs == "" {
		fmt.Fprintf(eout, "nether routes or htdocs are set.")
		return CONFIG
	}

	if router.Port == 0 {
		fmt.Fprintf(eout, "port is not set, default is not available\n")
		return CONFIG
	}

	if verbose && router.Routes != nil {
		for _, route := range router.Routes {
			route.Debug = true
		}
	}

	if dryRun {
		fmt.Fprintf(out, "configuration and routes successfully processed\n")
		return OK
	}

	// Launch web services
	fmt.Fprintf(out, "starting %s listening on port :%d (use Ctr-c to exit)\n", appName, router.Port)
	if router.Htdocs != "" {
		dir, err := filepath.Abs(router.Htdocs)
		if err == nil {
			fmt.Fprintf(out, "\tstatic content %s\n", dir)
		} else {
			fmt.Fprintf(out, "\tstatic content %s (warning: %s)\n", router.Htdocs, err)
		}
	}
	if err := router.ListenAndServe(); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return ROUTER_FAIL
	}
	return OK
}

// RunStaticWebServer this provides a localhost for static file content.
func RunStaticWebServer(in io.Reader, out io.Writer, eout io.Writer, args []string, port int, verbose bool) int {
	appName := "Newt Static Web Server"
	if port == 0 {
		port = SWS_PORT
	}
	htdocs := SWS_HTDOCS
	if len(args) > 0 {
		htdocs = args[0]
	}
	fmt.Fprintf(out, "starting %s listening on port :%d (use Ctr-c to exit)\n", appName, port)
	if err := NewtStaticFileServer(port, htdocs, verbose); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return SWS_FAIL
	}
	return OK
}

// NewtRunCheckYAML will load a Newt YAML fiel and make sure it can parse the configuration.
func RunNewtCheckYAML(in io.Reader, out io.Writer, eout io.Writer, args []string, verbose bool) int {
	if len(args) == 0 {
		fmt.Fprintf(eout, "missing Newt YAML filename to check\n")
		return CHECK_FAIL
	}
	fName := args[0]
	cfg, err := LoadConfig(fName)
	if err != nil {
		fmt.Fprintf(eout, "%s error: %s\n", fName, err)
		return CHECK_FAIL
	}
	if cfg.Applications == nil {
		fmt.Fprintf(eout, "%s has no applications defined\n", fName)
		return CHECK_FAIL
	}
	if cfg.Models == nil || len(cfg.Models) == 0 {
		if cfg.Applications.PostgREST != nil {
			fmt.Fprintf(eout, "WARNING: %s has no models defined\n", fName)
		} else if cfg.Applications.NewtMustache != nil {
			fmt.Fprintf(eout, "WARNING: %s has no models defined\n", fName)
		}
	}
	if cfg.Routes == nil || len(cfg.Routes) == 0 {
		if cfg.Applications.NewtRouter != nil {
			fmt.Fprintf(eout, "%s has no routes defined for Newt Router\n", fName)
			return CHECK_FAIL
		}
	}
	if cfg.Templates == nil || len(cfg.Templates) == 0 {
		if cfg.Applications.NewtMustache != nil {
			fmt.Fprintf(eout, "%s has no templates defined for Newt Mustache\n", fName)
			return CHECK_FAIL
		}
	}
	if verbose {
		if cfg.Applications.PostgREST != nil {
			fmt.Fprintf(out, "PostgREST configuration is %s\n", cfg.Applications.PostgREST.ConfPath)
			fmt.Fprintf(out, "PostgREST will be run with the command %q\n", strings.Join([]string{
				cfg.Applications.PostgREST.AppPath,
				cfg.Applications.PostgREST.ConfPath,
			}, " "))
		}
		if cfg.Models != nil {
			for _, m := range cfg.Models {
				fmt.Printf("models %s defined, %d elements\n", m.Id, len(m.Body))
				if m.Description != "" {
					fmt.Fprintf(out, "\t%s\n\n", m.Description)
				}
			}
		}
		if cfg.Applications.NewtRouter != nil {
			port := ROUTER_PORT
			if cfg.Applications.NewtRouter.Port != 0 {
				port = cfg.Applications.NewtRouter.Port
			}
			fmt.Fprintf(out, "Newt Router configured, port set to :%d\n", port)
			if cfg.Applications.NewtRouter.Htdocs != "" {
				fmt.Fprintf(out, "Static content will be served from %s\n", cfg.Applications.NewtRouter.Htdocs)
			}
			if cfg.Routes != nil {
				for _, r := range cfg.Routes {
					fmt.Fprintf(out, "route %s defined, request path %s, pipeline size %d\n", r.Id, r.Pattern, len(r.Pipeline))
					if r.Description != "" {
						fmt.Fprintf(out, "\t%s\n\n", r.Description)
					}
				}
			}
		}
		if cfg.Applications.NewtMustache != nil {
			port := cfg.Applications.NewtMustache.Port
			if port == 0 {
				port = MUSTACHE_PORT
			}
			fmt.Fprintf(out, "Newt Mustache configured, port set to :%d\n", port)
			fmt.Fprintf(out, "%d Mustache Templates are defined\n", len(cfg.Templates))
			for _, mt := range cfg.Templates {
				tList := []string{
					mt.Template,
				}
				if len(mt.Partials) > 0 {
					tList = append(tList, mt.Partials...)
				}
				fmt.Fprintf(out, "http://localhost:%s%s points at %s\n", port, mt.Pattern, strings.Join(tList, ", "))
				if mt.Description != "" {
					fmt.Fprintf(out, "\t%s\n\n", mt.Description)
				}
			}
		}
	}
	return OK
}

// RunNewt is a runner that can run Newt Mustache, Newt Router and PostgREST if defined in the Newt YAML file.
func RunNewt(in io.Reader, out io.Writer, eout io.Writer, args []string, verbose bool) int {
	appName := path.Base(os.Args[0])
	action := ""
	// Extract our actions from the args
	switch len(args) {
	case 0:
		return NEWT_FAIL
	case 1:
		action, args = args[0], []string{}
	default:
		action, args = args[0], args[1:]
	}

	switch action {
	case "init":
		return RunNewtInit(in, out, eout, args, verbose)
	case "check":
		return RunNewtCheckYAML(in, out, eout, args, verbose)
	case "generate":
	//FIXME: I need to back up all the expected filenames for project first, then
	// for each generate action option a new output buffer to render each new version of the file.
		fmt.Fprintf(eout, "%s init action is not implemented\n", appName)
		return INIT_FAIL
	case "run":
		return RunNewtApplications(in, out, eout, args, verbose)
	case "sws":
		return RunStaticWebServer(in, out, eout, args, 0, verbose)
	default:
		fmt.Fprintf(eout, "%s does %q is an unsupported action, see %s -help\n", appName, action, appName)
		return UNSUPPORTED_ACTION
	}
	return OK
}

// RunNewtApplications will run the applictions defined in your Newt YAML file.
func RunNewtApplications(in io.Reader, out io.Writer, eout io.Writer, args []string, verbose bool) int {
	var fName string
	appName := path.Base(os.Args[0])
	// get the Newt YAML filename.
	if len(args) > 0 {
		fName = args[0]
	}
	// Newt YAML should be present so the actions available are "generate" and "run"
	if fName == "" {
		fmt.Fprintf(eout, "%s expected a Newt YAML filename\n", appName)
		return CONFIG
	}
	cfg, err := LoadConfig(fName)
	if err != nil {
		fmt.Fprintf(eout, "%s failed to load %q, %s", appName, fName, err)
		return CONFIG
	}
	// Startup PostgREST if configured in the Newt YAML file.
	if cfg.Applications != nil && cfg.Applications.PostgREST != nil &&
		cfg.Applications.PostgREST.ConfPath != "" && cfg.Applications.PostgREST.AppPath != "" {
		postgREST := cfg.Applications.PostgREST
		cwd, err := os.Getwd()
		if err != nil {
			log.Println(err)
			return POSTGREST_FAIL
		}
		cmd := exec.Command(postgREST.AppPath, postgREST.ConfPath)
		cmd.Dir = cwd
		// Setup the stdin and stdout to be visible from Newt
		cmd.Stdout = out
		cmd.Stderr = eout
		log.Printf("starting %s %s listening on :%d (use Ctrl-c to exit)", postgREST.AppPath, postgREST.ConfPath, postgREST.Port)
		err = cmd.Start()
		if err != nil {
			log.Println(err)
			return POSTGREST_FAIL
		}
		log.Printf("%s running with pid %d in the backround", postgREST.AppPath, cmd.Process.Pid)
		cmd.Process.Release()
	}
	// Setup and start Newt Mustache first
	if cfg.Applications != nil && cfg.Applications.NewtMustache != nil {
		go func() {
			RunNewtMustache(in, out, eout, args, 0, 0, verbose)
		}()
	}

	// The router starts up second and is what prevents service from falling through.
	if cfg.Applications != nil && cfg.Applications.NewtRouter != nil {
		go func() {
			RunNewtRouter(in, out, eout, args, false, 0, "", verbose)
		}()
	}
	// NOTE: we need to wait for a signal so that our external process and Go routines aan run.

	// Set up channel on which to send signal notifications.
	c := make(chan os.Signal, 1)

	// We're listening for all signals, probably should narrow this down.
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until any signal is received.
	s := <-c
	fmt.Println("exited with signal:", s)
	return OK
}

// RunMustacheCLI this provides a cli for checking your templates using static JSON files and
// displaying results to stdout.
func RunMustacheCLI(in io.Reader, out io.Writer, eout io.Writer, args []string, pageElements map[string]interface{}) int {
	var (
		tmplFName string
		dataFName string
		txt       []byte
		src       []byte
		data      *interface{}
		err       error
	)
	if len(args) == 1 {
		tmplFName, dataFName = args[0], "-"
	} else if len(args) == 2 {
		tmplFName, dataFName = args[0], args[1]
	} else {
		fmt.Fprintf(eout, "expected a JSON data file and template filename\n")
		return DATA_ERROR
	}
	txt, err = os.ReadFile(tmplFName)
	if err != nil {
		fmt.Fprintf(eout, "failed reading %q, %s\n", tmplFName, err)
		return READ_ERROR
	}

	if dataFName == "" || dataFName == "-" {
		dataFName = "stdin"
		src, err = io.ReadAll(in)
	} else {
		src, err = os.ReadFile(dataFName)
	}
	if err != nil {
		fmt.Fprintf(eout, "failed reading data %q, %s\n", dataFName, err)
		return READ_ERROR
	}
	decoder := json.NewDecoder(bytes.NewBuffer(src))
	decoder.UseNumber()
	if err = decoder.Decode(&data); err != nil {
		fmt.Fprintf(eout, "failed decoding %q, %s\n", dataFName, err)
		return DECODE_ERROR
	}
	if pageElements == nil || len(pageElements) == 10 {
		pageElements = map[string]interface{}{}
	}

	tmpl, err := mustache.ParseString(fmt.Sprintf("%s", txt))
	if err != nil {
		fmt.Fprintf(eout, "failed template parse error %q, %s\n", dataFName, err)
		return TEMPLATE_ERROR
	}
	if err = tmpl.FRender(out, pageElements); err != nil {
		fmt.Fprintf(eout, "failed render error %q, %s\n", dataFName, err)
		return TEMPLATE_ERROR
	}
	return OK
}

func getAnswer(buf *bufio.Reader, lower bool) string {
	answer, err := buf.ReadString('\n')
	if err != nil {
		return ""
	}
	answer = strings.TrimSpace(answer)
	if lower {
		return strings.ToLower(answer)
	}
	return answer
}

// appendToPipeline will append a Newt Mustache to the end of the pipeline for routeID and objName
func appendToPipeline(routes []*NewtRoute, routeId string, objName string, method string, port int, servicePath string) error {
	for _, r := range routes {
		if r.Id == routeId {
			r.Pipeline = append(r.Pipeline, &Service{
				Service:     fmt.Sprintf("%s http://localhost:%d/%s", method, port, servicePath),
				Description: fmt.Sprintf("Example of %s service for %s", routeId, objName),
			})
			return nil
		}
	}
	return fmt.Errorf("failed to find route %q", routeId)
}

func setupNewtRouter(cfg *Config, buf *bufio.Reader, out io.Writer, appFName string, objName string) {
	fmt.Fprintf(out, "Will %s use Newt Router (Y/n)? ", appFName)
	answer := getAnswer(buf, true)
	if answer != "n" {
		if cfg.Applications == nil {
			cfg.Applications = &Applications{}
		}
		if cfg.Applications.NewtRouter == nil {
			cfg.Applications.NewtRouter = &Application{
				Port:   ROUTER_PORT,
				Htdocs: "htdocs",
			}
		}
		if cfg.Routes == nil {
			cfg.Routes = []*NewtRoute{}
		}
	}
}

func setupPostgREST(cfg *Config, buf *bufio.Reader, out io.Writer, appFName string, objName string) {
	fmt.Fprintf(out, "Will %s use PostgREST (y/N)? ", appFName)
	answer := getAnswer(buf, true)
	if answer == "y" {
		if cfg.Applications.PostgREST == nil {
			cfg.Applications.PostgREST = &Application{
				Port:     3000,
				AppPath:  "postgrest",
				ConfPath: "postgrest.conf",
			}
		}
		// We create our restful actions for interacting with the PostgREST JSON API.
		// NOTE: If Newt Mustache is used then the Mustache setup will add routes for retrieving editor forms and will append
		// result templates the routes defined here.
		for action, method := range map[string]string{"create": http.MethodPost, "read": http.MethodGet, "update": http.MethodPut, "delete": http.MethodDelete, "list": http.MethodGet} {
			routeId := fmt.Sprintf("%s_%s", objName, action)
			request := fmt.Sprintf("%s /%s", method, objName)
			if action == "read" || action == "update" || action == "delete" {
				request = fmt.Sprintf("%s /%s/{oid}", method, routeId)
			}
			servicePath := fmt.Sprintf("rpc/%s_%s", objName, action)
			cfg.Routes = append(cfg.Routes, &NewtRoute{
				Id:      routeId,
				Pattern: request,
			})
			appendToPipeline(cfg.Routes, routeId, objName, method, cfg.Applications.PostgREST.Port, servicePath)
		}
	}
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
			action := "create"
			method := http.MethodGet
			routeId := fmt.Sprintf("%s_%s_form", objName, action)
			request := fmt.Sprintf("%s /%s/{$}", method, objName)
			pattern := fmt.Sprintf("/%s_%s_form", objName, action)
			tName := fmt.Sprintf("%s_%s_form.tmpl", objName, action)
			cfg.Routes = append(cfg.Routes, &NewtRoute{
				Id:          routeId,
				Pattern:     request,
				Description: "Example display a create web form",
			})
			cfg.Templates = append(cfg.Templates, &MustacheTemplate{
				Pattern:     pattern,
				Template:    tName,
				Description: "Example display an create web from",
			})
			appendToPipeline(cfg.Routes, routeId, objName, http.MethodPost, cfg.Applications.NewtMustache.Port, pattern)

			action = "update"
			routeId = fmt.Sprintf("%s_%s_form", objName, action)
			request = fmt.Sprintf("%s /%s/{oid}", method, objName)
			pattern = fmt.Sprintf("/%s_%s_form", objName, action)
			tName = fmt.Sprintf("%s_%s_form.tmpl", objName, action)
			cfg.Routes = append(cfg.Routes, &NewtRoute{
				Id:          routeId,
				Pattern:     request,
				Description: "Example display a update web form",
			})
			cfg.Templates = append(cfg.Templates, &MustacheTemplate{
				Pattern:     pattern,
				Template:    tName,
				Description: "Example display an update web from",
			})
			appendToPipeline(cfg.Routes, routeId, objName, http.MethodPost, cfg.Applications.NewtMustache.Port, pattern)

			action = "delete"
			routeId = fmt.Sprintf("%s_%s_form", objName, action)
			request = fmt.Sprintf("%s /%s/{oid}", method, objName)
			pattern = fmt.Sprintf("/%s_%s_form", objName, action)
			tName = fmt.Sprintf("%s_%s_form.tmpl", objName, action)
			cfg.Routes = append(cfg.Routes, &NewtRoute{
				Id:          routeId,
				Pattern:     request,
				Description: "Example display a delete web form",
			})
			cfg.Templates = append(cfg.Templates, &MustacheTemplate{
				Pattern:     pattern,
				Template:    tName,
				Description: "Example display a delete web from",
			})
			appendToPipeline(cfg.Routes, routeId, objName, http.MethodPost, cfg.Applications.NewtMustache.Port, pattern)

			// Now add the mappings for templates that return results.
			for _, action := range []string{"create", "read", "update", "delete", "list"} {
				routeId := fmt.Sprintf("%s_%s", objName, action)
				pattern = fmt.Sprintf("/%s_%s", objName, action)
				tName := fmt.Sprintf("%s_%s.tmpl", objName, action)
				cfg.Templates = append(cfg.Templates, &MustacheTemplate{
					Pattern:     pattern,
					Template:    tName,
					Description: "This is an example of defining a template handler",
				})
				appendToPipeline(cfg.Routes, routeId, objName, http.MethodPost, cfg.Applications.NewtMustache.Port, pattern)
			}
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

// RunNewtInit will initialize a Newt project by creating a Newt YAML file interactively.
func RunNewtInit(in io.Reader, out io.Writer, eout io.Writer, args []string, verbose bool) int {
	var (
		answer   string
		mkBackUp bool
	)
	cfg := &Config{}
	readBuffer := bufio.NewReader(in)
	// Step 1. Figure out what we're going to call our generated Newt YAML file.
	appFName := ""
	if len(args) > 0 {
		appFName = args[0]
	}
	if appFName == "" {
		s, err := os.Getwd()
		if err == nil {
			s = path.Base(s)
		}
		if s != "" {
			appFName = strings.ToLower(strings.ReplaceAll(s, " ", "_")) + ".yaml"
		} else {
			appFName = "app.yaml"
		}
	}
	if _, err := os.Stat(appFName); err == nil {
		mkBackUp = true
		cfg, err = LoadConfig(appFName)
		fmt.Fprintf(out, "found %q, continue (y/N)? ", appFName)
		answer = getAnswer(readBuffer, true)
		if answer != "y" {
			fmt.Fprintf(eout, "aborting init, %q already exists\n", appFName)
			return INIT_FAIL
		}
	}
	// objName is used to gererate example elements in the Newt YAML file.
	objName := strings.TrimSuffix(appFName, ".yaml")
	// Step 2. Figure out which applications will be running
	if cfg.Applications == nil {
		cfg.Applications = &Applications{}
	}
	for {
		setupNewtRouter(cfg, readBuffer, out, appFName, objName)
		setupPostgREST(cfg, readBuffer, out, appFName, objName)
		setupNewtMustache(cfg, readBuffer, out, appFName, objName)
		setupNewtGenerator(cfg, readBuffer, out, appFName, objName)
		setupEnvironment(cfg, readBuffer, out, appFName, objName)
		setupOptions(cfg, readBuffer, out, appFName, objName)

		// Now output the YAML
		comment := []byte(fmt.Sprintf(`#!/bin/env newt check
#
# This was generated by %s, version %s %s, on %s.
#
`, path.Base(os.Args[0]), Version, ReleaseHash, ReleaseDate))
		data := bytes.NewBuffer(comment)
		encoder := yaml.NewEncoder(data)
		encoder.SetIndent(2)
		if err := encoder.Encode(cfg); err != nil {
			fmt.Fprintf(eout, "Failed to generate %s, %s\n", appFName, err)
			return INIT_FAIL
		}
		src := data.Bytes()

		fmt.Fprintf(out, "%s\n", src)
		fmt.Fprintf(out, "Save and exit (y/N)? ")
		answer = getAnswer(readBuffer, true)
		if answer == "y" {
			// We're ready to write out result.
			// If file exists make a back up copy
			if mkBackUp {
				buf, err := os.ReadFile(appFName)
				if err != nil {
					fmt.Fprintf(eout, "failed to back up %s, aborting write\n", appFName)
					return INIT_FAIL
				}
				if err := os.WriteFile(appFName+".bak", buf, 0666); err != nil {
					fmt.Fprintf(eout, "failed to write back up %s, aborting write, %s\n", appFName, err)
					return INIT_FAIL
				}
			}
			if err := os.WriteFile(appFName, src, 0666); err != nil {
				fmt.Fprintf(eout, "failed to write %s, %s\n", appFName, err)
				return INIT_FAIL
			}
			break
		} else {
			fmt.Fprintf(out, "Exit without saving (y/N)? ")
			answer = getAnswer(readBuffer, true)
			if answer == "y" {
				fmt.Fprintf(out, "aborting write of %q\n", appFName)
				return INIT_FAIL
			}
		}
	}
	if _, err := os.Stat(".git"); err == nil {
		fmt.Fprintf(out, `It appears your are using the git revision control system.
You should make sure that generated code containing secrets is NOT included in your
repository. It is recommented that add the following to your .gitignore file.

    # Newt Project ignore list.
    *setup*.sql
    postgrest.conf
`, objName)
	}
	return OK
}
