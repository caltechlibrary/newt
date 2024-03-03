/**
 * cli.go an implements runners for the cli of the Newt Project.
 *
 * @author R. S. Doiel
 */
package newt
	

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"time"

	// 3rd Party Templates
	"github.com/cbroglie/mustache"
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
	ROUTER_PORT = ":8010"
	MUSTACHE_PORT    = ":8011"
	MUSTACHE_TIMEOUT = 3 * time.Second
	SWS_PORT   = 8000
	SWS_HTDOCS = "."
)

// RunNewtGenerator is a runner for generating SQL and templates from our Newt YAML file.
func RunNewtGenerator(in io.Reader, out io.Writer, eout io.Writer, args []string) int {
	//appName := "Newt Generator"
	fName, target, codeType := "", "", ""
	if len(args) > 0 {
		fName = args[0]
	} else {
		fmt.Fprintf(eout, "missing YAML configuration file\n")
		return CONFIG
	}
	if len(args) > 1 {
		target = args[1]
	} else {
		fmt.Fprintf(eout, "missing generator name, e.g. postgres, postgrest, mustache\n")
		return CONFIG
	}
	if len(args) > 2 {
		codeType = args[2]
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
	if err := generator.Generate(target, codeType); err != nil {
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
	if mt.Port == "" {
		mt.Port = MUSTACHE_PORT
	}
	// See if we have a command line option for port to process.
	if port != 0 {
		mt.Port = fmt.Sprintf(":%d", port)
	}
	if verbose {
		fmt.Fprintf(out, "port set to %q\n", mt.Port)
	}
	// Onelast check to make sure the port number as the colon prefix
	if !strings.HasPrefix(mt.Port, ":") {
		mt.Port = fmt.Sprintf(":%s", mt.Port)
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
	fmt.Printf("starting %s listening on port %s (press Ctrl-c to exit)\n", appName, mt.Port)
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
	if router.Port == "" {
		router.Port = ROUTER_PORT
	}
	if port != 0 {
		router.Port = fmt.Sprintf(":%d", port)
	}
	// Prefix the port number with a colon
	if !strings.HasPrefix(router.Port, ":") {
		router.Port = fmt.Sprintf(":%s", router.Port)
	}
	if htdocs != "" {
		router.Htdocs = htdocs
	}

	// Are we ready to run service?
	if router.Routes == nil && router.Htdocs == "" {
		fmt.Fprintf(eout, "nether routes or htdocs are set.")
		return CONFIG
	}

	if router.Port == "" || router.Port == ":" {
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

	// Launch web service
	fmt.Fprintf(out, "starting %s listening on port %s (use Ctr-c to exit)\n", appName, router.Port)
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
				port = fmt.Sprintf(":%d", cfg.Applications.NewtRouter.Port)
			}
			fmt.Fprintf(out, "Newt Router configured, port set to %s\n", port)
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
			port := fmt.Sprintf("%d", cfg.Applications.NewtMustache.Port)
			if port == "" {
				port = fmt.Sprintf(":%d", MUSTACHE_PORT)
			}
			fmt.Fprintf(out, "Newt Mustache configured, port set to %s\n", port)
			fmt.Fprintf(out, "%d Mustache Templates are defined\n", len(cfg.Templates))
			for _, mt := range cfg.Templates {
				tList := []string{
					mt.Template,
				}
				if len(mt.Partials) > 0 {
					tList = append(tList, mt.Partials...)
				}
				fmt.Fprintf(out, "http://localhost%s%s points at %s\n", port, mt.Pattern, strings.Join(tList, ", "))
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
			fmt.Fprintf(eout, "%s init action is not implemented\n", appName)
			return INIT_FAIL
		case "check":
			return RunNewtCheckYAML(in, out, eout, args, verbose)
		case "generate":
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
		func() {
			RunNewtRouter(in, out, eout, args, false, 0, "", verbose)
		}()
	}
	// NOTE: we need to wait for a signal so that our external process and Go routines aan run.

	// Set up channel on which to send signal notifications.
	c := make(chan os.Signal, 1)

	// We're listening for all signals, probably should narrow this down.
	signal.Notify(c)

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
	if pageElements == nil || len(pageElements) ==10 {
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
