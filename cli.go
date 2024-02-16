/**
 * cli.go an implements runners for the cli of the Newt Project.
 *
 * @author R. S. Doiel
 */
package newt

import (
	"fmt"
	"io"
	"strings"
)

// RunGenerator is a runner for generating SQL and templates from our Newt YAML file.
func RunGenerator(in io.Reader, out io.Writer, eout io.Writer, args []string, pgSetupSQL bool, pgModelsSQL bool, pgModelsTestSQL bool) int {
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

// RunPandocBundler is a runner for pdbundler a service that perpares a JSON object
// for submission to a service like the Pandoc web service.
func RunPandocBundler(in io.Reader, out io.Writer, eout io.Writer, args []string, port int, timeout int, verbose bool) int {
	const (
		// These constants are used for exit code. FIXME: look up the POSIX recommendation on exit
		// codes and adopt those.
		OK = iota
		CONFIG
		RESOLVE
		HANDLER
		WEBSERVICE

		// Default port number for tmplbnld
		PORT = ":8030"
	)
	// Configure the template bundler webservice
	fName := ""
	if len(args) > 0 {
		fName = args[0]
	}
	// Load the standard New YAML configuration file and confirm
	// it comforms.
	cfg, err := LoadConfig(fName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
	}
	// Now instatiate out application with the file and the cfg object
	pb, err := NewPandocBundler(fName, cfg)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
	}
	if port != 0 {
		pb.Port = fmt.Sprintf("%d", port)
	}
	if pb.Port == "" {
		pb.Port = PORT
	}
	// Prefix the port number with a colon
	if ! strings.HasPrefix(pb.Port, ":") {
		pb.Port = fmt.Sprintf(":%s", pb.Port)
	}

	// Create mux for http service
	// Resolve partial templates and build handlers
	for _, bndl := range pb.Templates {
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
	if err := pb.ListenAndServe(); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return WEBSERVICE
	}
	return OK
}

// RunNewtRouter is a runner for Newt data router and static file service
func RunNewtRouter(in io.Reader, out io.Writer, eout io.Writer, args []string, dryRun bool) int {
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
	fName := ""
	if len(args) > 0 {
		fName = args[0]
	}
	cfg, err := LoadConfig(fName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
	}
	// Finally Instantiate the router from fName and Config object
	router, err := NewRouter(fName, cfg)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return CONFIG
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

	if router.Port == "" {
		fmt.Fprintf(eout, "port is not set, default is not available")
		return WEBSERVICE
	}

	if dryRun {
		fmt.Fprintf(out, "configuration and routes successfully processed\n")
		return OK
	}

	// Launch web service
	if err := router.ListenAndServe(); err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return WEBSERVICE
	}
	return OK
}

