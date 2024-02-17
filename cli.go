/**
 * cli.go an implements runners for the cli of the Newt Project.
 *
 * @author R. S. Doiel
 */
package newt

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

// RunNewtGenerator is a runner for generating SQL and templates from our Newt YAML file.
func RunNewtGenerator(in io.Reader, out io.Writer, eout io.Writer, args []string, pgSetupSQL bool, pgModelsSQL bool, pgModelsTestSQL bool) int {
	const (
		// These constants are used for exit code. FIXME: look up the POSIX recommendation on exit
		// codes and adopt those.
		OK = iota
		CONFIG
		GENFAIL
	)
	fName := ""
	if len(args) > 0 {
		fName = args[0]
	}
	cfg, err := LoadConfig(fName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
	}
	generator, err := NewGenerator(fName, cfg)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
	}
	if err := generator.Generate(); err != nil {
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
		PORT = ":8040"
	)
	// Configure the template bundler webservice
	fName := ""
	if len(args) > 0 {
		fName = args[0]
	}
	// Load the Newt YAML syntax file holding the configuration
	// and make sure it conforms.
	cfg, err := LoadConfig(fName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
	}
	// Instantiate the specific application with the filename and Config object
	mb, err := NewNewtMustache(fName, cfg)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
	}
	if cfg.Application != nil  && cfg.Application.Port != 0 {
		mb.Port = fmt.Sprintf("%d", port)
	}
	if mb.Port == "" {
		mb.Port = PORT
	}
	// Prefix the port number with a colon
	if ! strings.HasPrefix(mb.Port, ":") {
		mb.Port = fmt.Sprintf(":%s", mb.Port)
	}

    fmt.Printf("starting %s\n", path.Base(os.Args[0]))
	// Create mux for http service
	// Resolve partial templates and build handlers
	for _, bndl := range mb.Templates {
		if verbose {
			bndl.Debug = true
		}
		if err := bndl.ResolveTemplate(); err != nil {
			fmt.Fprintf(eout, "%s failed to resolve, %s\n", bndl.Template, err)
			return RESOLVE
		}
		if err := bndl.ResolvePath(); err != nil {
			fmt.Fprintf(eout, "failed to build handler for %q, %s\n", bndl.Pattern, err)
			return HANDLER
		}
	}
	// Launch web service
    fmt.Printf("listening on port %s\n", mb.Port)
	if err := mb.ListenAndServe(); err != nil {
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
		PORT = ":8020"
	)
	// You can run Newt Router with just an htdocs directory. If so you don't require a config file.
	var err error
	cfg := &Config{}
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
		router, err = NewNewtRouter(fName, cfg)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return CONFIG
		}
		if cfg.Application != nil {
			if port == 0 && cfg.Application.Port > 0{
				router.Port = fmt.Sprintf(":%d", cfg.Application.Port)
			}
			if htdocs == "" && cfg.Application.Htdocs != "" {
				router.Htdocs = cfg.Application.Htdocs
			}
			if len(cfg.Application.Environment) > 0 {
				router.Environment = append(router.Environment, cfg.Application.Environment...)
			}
		}
	}
	if router.Port == "" {
		router.Port = PORT
	}
	if port != 0 {
		router.Port = fmt.Sprintf(":%d", port)
	}
	// Prefix the port number with a colon
	if ! strings.HasPrefix(router.Port, ":") {
		router.Port = fmt.Sprintf(":%s", router.Port)
	}

	if htdocs != "" {
		router.Htdocs = htdocs
	}

	// Are we ready to run service?
	if router.Routes == nil && router.Htdocs == "" {
		if fName != "" {
			fmt.Fprintf(eout, "nether routes or htdocs are set.")
			return CONFIG
		} 
		fmt.Fprintf(eout, "NEWT_ROUTES and NEWT_HTDOCS are undefined.")
		return CONFIG
	}

	if router.Port == "" || router.Port == ":" {
		fmt.Fprintf(eout, "port is not set, default is not available")
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

