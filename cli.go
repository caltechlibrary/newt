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

// RunNewtGenerator is a runner for generating SQL and templates from our Newt YAML file.
func RunNewtGenerator(in io.Reader, out io.Writer, eout io.Writer, args []string) int {
	const (
		// These constants are used for exit code. FIXME: look up the POSIX recommendation on exit
		// codes and adopt those.
		OK = iota
		CONFIG
		GENFAIL
	)
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
		return GENFAIL
	}
	return OK
}

// RunNewtMustache is a runner for a Mustache redner engine service based on the Pandoc server API.
func RunNewtMustache(in io.Reader, out io.Writer, eout io.Writer, args []string, port int, timeout int, verbose bool) int {
	const (
		// These constants are used for exit code. FIXME: look up the POSIX recommendation on exit
		// codes and adopt those.
		OK = iota
		CONFIG
		RESOLVE
		HANDLER
		WEBSERVICE

		// Default port number for tmplbnld
		PORT    = ":8011"
		TIMEOUT = 3 * time.Second
	)
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
		mt.Port = PORT
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
		mt.Timeout = TIMEOUT
	}
	if len(mt.Templates) == 0 {
		fmt.Fprintf(eout, "no templates in configuration\n")
		return CONFIG
	}
	fmt.Printf("starting %s\n", path.Base(os.Args[0]))
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
	fmt.Printf("listening on port %s\n", mt.Port)
	if err := mt.ListenAndServe(); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return WEBSERVICE
	}
	return OK
}

// RunNewtRouter is a runner for Newt data router and static file service
func RunNewtRouter(in io.Reader, out io.Writer, eout io.Writer, args []string, dryRun bool, port int, htdocs string, verbose bool) int {
	const (
		// These constants are used for exit code. FIXME: look up the POSIX recommendation on exit
		// codes and adopt those.
		OK = iota
		CONFIG
		RESOLVE
		HANDLER
		WEBSERVICE

		// Default port number for tmplbnld
		PORT = ":8010"
	)
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
		router.Port = PORT
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
		return WEBSERVICE
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
	fmt.Fprintf(out, "%s listening on port %s\n", path.Base(os.Args[0]), router.Port)
	if err := router.ListenAndServe(); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return WEBSERVICE
	}
	return OK
}

// RunStaticWebServer this provides a localhost for static file content.
func RunStaticWebServer(in io.Reader, out io.Writer, eout io.Writer, args []string, port int, verbose bool) int {
	const (
		OK = iota
		SERVER_ERROR

		// Defaults
		PORT   = 8000
		HTDOCS = "."
	)
	if port == 0 {
		port = PORT
	}
	htdocs := HTDOCS
	if len(args) > 0 {
		htdocs = args[0]
	}
	if err := NewtStaticFileServer(port, htdocs, verbose); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return SERVER_ERROR
	}
	return OK
}

// RunNewt is a runner that can run Newt Mustache, Newt Router and PostgREST if defined in the Newt YAML file.
func RunNewt(in io.Reader, out io.Writer, eout io.Writer, args []string, verbose bool) int {
	const (
		// These constants are used for exit code. FIXME: look up the POSIX recommendation on exit
		// codes and adopt those.
		OK = iota
		CONFIG
		RESOLVE
		HANDLER
		WEBSERVICE
		POSTGREST
	)
	appName := path.Base(os.Args[0])
	fName := ""
	if len(args) > 0 {
		fName = args[0]
	} else {
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
			return POSTGREST
		}
		cmd := exec.Command(postgREST.AppPath, postgREST.ConfPath)
		cmd.Dir = cwd
		// Setup the stdin and stdout to be visible from Newt
		cmd.Stdout = out
		cmd.Stderr = eout
		log.Printf("Starting %s %s on port :%d", postgREST.AppPath, postgREST.ConfPath, postgREST.Port)
		err = cmd.Start()
		if err != nil {
			log.Println(err)
			return POSTGREST
		}
		log.Printf("%s running with pid %d in the backround", postgREST.AppPath, cmd.Process.Pid)
		cmd.Process.Release()
	}
	// Setup and start Newt Mustache first
	if cfg.Applications != nil && cfg.Applications.NewtMustache != nil {
		go func() {
			const (
				PORT    = ":8011"
				TIMEOUT = 3 * time.Second
			)
			// Instantiate the specific application with the filename and Config object
			mt, err := NewNewtMustache(cfg)
			if err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return
			}
			// If port is not set in the config, set it to the default port.
			if mt.Port == "" {
				mt.Port = PORT
			}
			if verbose {
				fmt.Fprintf(out, "port set to %q\n", mt.Port)
			}
			// Onelast check to make sure the port number as the colon prefix
			if !strings.HasPrefix(mt.Port, ":") {
				mt.Port = fmt.Sprintf(":%s", mt.Port)
			}
			if mt.Timeout == 0 {
				mt.Timeout = TIMEOUT
			}
			if len(mt.Templates) == 0 {
				fmt.Fprintf(eout, "no templates in configuration\n")
				return
			}
			fmt.Printf("starting %s\n", path.Base(os.Args[0]))
			// Create mux for http service
			// Resolve partial templates and build handlers
			for _, tmpl := range mt.Templates {
				if verbose {
					tmpl.Debug = true
				}
				if err := tmpl.ResolveTemplate(); err != nil {
					fmt.Fprintf(eout, "%s failed to resolve, %s\n", tmpl.Template, err)
					return
				}
				if err := tmpl.ResolvePath(); err != nil {
					fmt.Fprintf(eout, "failed to build handler for %q, %s\n", tmpl.Pattern, err)
					return
				}
			}
			// Launch web service
			fmt.Printf("%s Newt Mustache listening on port %s\n", appName, mt.Port)
			if err := mt.ListenAndServe(); err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return
			}
		}()
	}

	// The router starts up second and is what prevents service from falling through.
	if cfg.Applications != nil && cfg.Applications.NewtRouter != nil {
		func() {
			const (
				// Default port number for tmplbnld
				PORT = ":8010"
			)
			// Finally Instantiate the router from fName and Config object
			router, err := NewNewtRouter(cfg)
			if err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return
			}

			if router.Port == "" {
				router.Port = PORT
			}
			// Prefix the port number with a colon
			if !strings.HasPrefix(router.Port, ":") {
				router.Port = fmt.Sprintf(":%s", router.Port)
			}

			// Are we ready to run service?
			if router.Routes == nil && router.Htdocs == "" {
				fmt.Fprintf(eout, "nether routes or htdocs are set.")
				return
			}

			if router.Port == "" || router.Port == ":" {
				fmt.Fprintf(eout, "port is not set, default is not available\n")
				return
			}

			if verbose && router.Routes != nil {
				for _, route := range router.Routes {
					route.Debug = true
				}
			}

			// Launch web service
			fmt.Fprintf(out, "%s Newt Router listening on port %s\n", appName, router.Port)
			if err := router.ListenAndServe(); err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return
			}
			return
		}()
	}

	// NOTE: we need to wait for a signal so that our external process and Go routines aan run.

	// Set up channel on which to send signal notifications.
	c := make(chan os.Signal, 1)

	// We're listening for all signals, probably should narrow this down.
	signal.Notify(c)

	// Block until any signal is received.
	s := <-c
	fmt.Println("Got signal:", s)
	return OK
}

// RunMustacheCLI this provides a cli for checking your templates using static JSON files and
// displaying results to stdout.
func RunMustacheCLI(in io.Reader, out io.Writer, eout io.Writer, args []string, pageElements map[string]interface{}) int {
	const (
		OK = iota
		ERROR
		READ_ERROR
		DECODE_ERROR
		TEMPLATE_ERROR
	)
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
		return ERROR
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
	if pageElements == nil || len(pageElements) == 0 {
		pageElements = map[string]interface{}{}
	}
	pageElements["body"] = data
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
