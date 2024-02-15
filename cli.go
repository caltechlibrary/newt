/**
 * cli.go an implements runners for the cli of the Newt Project.
 *
 * @author R. S. Doiel
 */
package newt

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

// RunPostgresSQL is a runner for generating SQL from a Newt YAML file.
func RunPostgresSQL(in io.Reader, out io.Writer, eout io.Writer, args []string, pgSetupSQL bool, pgModelsSQL bool, pgModelsTestSQL bool) int {
	configFName := ""
	if len(args) > 0 {
		configFName = args[0]
	}
	cfg, err := NewtLoad(configFName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	if cfg.Models == nil {
		fmt.Fprintln(eout, "-- No modules defined.")
		return 1
	}
	// For each module we generate a table create statement,
	// default view.
	if configFName == "" {
		configFName = "standard input"
	}
	exitCode := 0
	if pgSetupSQL {
		src, err := PgSetupSQL(configFName, cfg.Namespace, "")
		if err != nil {
			fmt.Fprintf(eout, "-- could not create setup for %q based on %q, %s\n", cfg.Namespace, configFName, err)
			exitCode = 1
		} else {
			fmt.Fprintf(out, "%s\n", src)
		}
	}
	if pgModelsSQL {
		modelNames := []string{}
		for i, model := range cfg.Models {
			src, err := PgModelSQL(configFName, model)
			if err != nil {
				fmt.Fprintf(eout, "-- could not create model %q (%d), %s\n", model.Name, i, err)
				exitCode = 1
			} else {
				fmt.Fprintf(out, "%s\n", src)
			}
			modelNames = append(modelNames, model.Name)
		}

		src, err := PgModelPermissions(configFName, cfg.Namespace, modelNames)
		if err != nil {
			fmt.Fprintf(eout, "-- could not permissions for models in %q, %s\n", configFName, err)
			exitCode = 1
		} else {
			fmt.Fprintf(out, "%s\n", src)
		}
	}
	if pgModelsTestSQL {
		for name, model := range cfg.Models {
			src, err := PgModelTestSQL(configFName, model)
			if err != nil {
				fmt.Fprintf(eout, "-- could not create model test %q, %s\n", name, err)
				exitCode = 1
			} else {
				fmt.Fprintf(out, "%s\n", src)
			}
		}
	}
	return exitCode
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
		PORT = ":3029"
	)
	// Configure the template bundler webservice
	fName := ""
	if len(args) > 0 {
		fName = args[0]
	}
	mb, err := NewNewtMustache(fName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
	}
	if port != 0 {
		mb.Port = fmt.Sprintf("%d", port)
	}
	if mb.Port == "" {
		mb.Port = PORT
	}
	// Prefix the port number with a colon
	if ! strings.HasPrefix(mb.Port, ":") {
		mb.Port = fmt.Sprintf(":%s", mb.Port)
	}

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
	if err := mb.ListenAndServe(); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return WEBSERVICE
	}
	return OK
}

// RunTemplateBundler is a runner for tmplbndl a service that perpares a JSON object
// for submission to a service like the Pandoc web service.
func RunTemplateBundler(in io.Reader, out io.Writer, eout io.Writer, args []string, port int, timeout int, verbose bool) int {
	const (
		// These constants are used for exit code. FIXME: look up the POSIX recommendation on exit
		// codes and adopt those.
		OK = iota
		CONFIG
		RESOLVE
		HANDLER
		WEBSERVICE

		// Default port number for tmplbnld
		PORT = ":3029"
	)
	// Configure the template bundler webservice
	fName := ""
	if len(args) > 0 {
		fName = args[0]
	}
	tb, err := NewTemplateBundler(fName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
	}
	if port != 0 {
		tb.Port = fmt.Sprintf("%d", port)
	}
	if tb.Port == "" {
		tb.Port = PORT
	}
	// Prefix the port number with a colon
	if ! strings.HasPrefix(tb.Port, ":") {
		tb.Port = fmt.Sprintf(":%s", tb.Port)
	}

	// Create mux for http service
	// Resolve partial templates and build handlers
	for _, bndl := range tb.Templates {
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
	if err := tb.ListenAndServe(); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return WEBSERVICE
	}
	return OK
}

// RunNewt is a runner for Newt data router and static file server
func RunNewt(in io.Reader, out io.Writer, eout io.Writer, args []string, dryRun bool) int {
	configFName := ""
	if len(args) > 0 {
		configFName = args[0]
	}
	cfg, err := NewtLoad(configFName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	// Finally make sure we have cfg.Htdocs or cfg.Routes 
	// set to run service.
	if cfg.Routes == nil && cfg.Htdocs == "" {
		if cfg.FName != "" {
			fmt.Fprintf(eout, "routes and htdocs are not set.")
			return 1
		} 
		fmt.Fprintf(eout, "NEWT_ROUTES and NEWT_HTDOCS are undefined.")
		return 1
	}

	router := new(Router)
	if cfg != nil {
		if err := router.Configure(cfg); err != nil {
			fmt.Fprintf(eout, "error configuring router from %q, %s\n", cfg.FName, err)
			return 1
		}
	}
	if err != nil {
		fmt.Fprintf(eout, "error reading routes from %q, %s\n", cfg.FName, err)
		return 1
	}
	if dryRun {
		fmt.Fprintf(out, "configuration and routes successfully processed\n")
		return 0
	}

	appName := path.Base(os.Args[0])
	mux := http.NewServeMux()
	switch {
	case cfg.Htdocs != "" && cfg.Routes != nil:
		log.Printf("%s using %s for static content and %s for router", appName, cfg.Htdocs, cfg.Routes)
		mux.Handle("/", NewLogger(router.Newt(http.FileServer(http.Dir(cfg.Htdocs)))))
	case cfg.Htdocs == "" && cfg.Routes != nil:
		log.Printf("%s using %s for router only", appName, cfg.Routes)
		mux.Handle("/", NewLogger(router.Newt(http.NotFoundHandler())))
	case cfg.Htdocs != "" && cfg.Routes == nil:
		log.Printf("%s using %s for static content only", appName, cfg.Htdocs)
		mux.Handle("/", NewLogger(http.FileServer(http.Dir(cfg.Htdocs))))
	default:
		log.Printf("Not configured, aborting")
		return 1
	}

	log.Printf("%s listening on port %s", appName, cfg.Port)
	http.ListenAndServe(":"+cfg.Port, mux)
	return 0
}

