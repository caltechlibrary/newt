package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	// Caltech Library Packages
	"github.com/caltechlibrary/newt"
)

var (
	helpText = `---
title: "{app_name}(1) user manual | {version} {release_hash}"
pubDate: {release_date}
author: "R. S. Doiel"
---

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS]

# DESCRIPTION

**{app_name}** is a web service that provides a Mustache template rendering inspired by Pandoc server.

Like Pandoc web server there is no configuration file. There are a few command line options, e.g.
port, template directory and timeout.

# Template engine API

**{app_name}** accepts POST of JSON content and maps them to a template name expressed in the 
request URL. If I had a template called `+"`"+`template/list_objects.tmpl`+"`"+`
then the url would be formed like `+"`"+`http://localhost:3032/list_objects.tmpl`+"`"+`. The JSON
encoded post would then be sent through the "list_objects.tmpl" template and returned to the browser.
The JSON object POSTed does not need a wrapping object like required by Pandoc server.  **{app_name}**
reads in the templates found in the template directory at start up. It creates a map between the 
basename of the template and the URL handler built from the that name. As such the templates are fixed
at startup and do not need to be passed along with your data object.

This improves of Pandoc web servers in a few ways. 

- Now wrapping object
- Template parse errors are known earlier own and are visible from the output log
- Parsing of the templates happens once at startup
- Partial templates can be supported because the startup phase of **{app_name}** handles resolving partials


# OPTIONS

The following options are supported by **{app_name}**.

-h, -help
: display this help message

-license
: display the software license

-version
: display version information

-port NUMBER
: (default is port is 3032) Set the port number to listen on

-templates
: (default directory name is "templates") Pick an alternative location for finding templates

--timeout SECONDS
: Timeout in seconds, after which a template rendering is aborted.  Default: 3.

# The templates

Mustache templates are documented at <https://mustache.github.io>. The template engined
used is based on Go package <https://github.com/cbroglie/mustache>.

`

	showHelp    bool
	showLicense bool
	showVersion bool
)

func main() {
	appName := path.Base(os.Args[0])
	// NOTE: The following variables are set when version.go is generated
	version := newt.Version
	releaseDate := newt.ReleaseDate
	releaseHash := newt.ReleaseHash
	fmtHelp := newt.FmtHelp

	// Standard Options
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "version", false, "display version")
	
	// App option(s)
	port, timeout, verbose := "3032", 3, false
	templates := "templates"
	flag.StringVar(&port, "port", port, "port number to listen on")
	flag.IntVar(&timeout, "timeout", timeout, "timeout for requests")
	flag.BoolVar(&verbose, "verbose", verbose, "display template debugging and request information to standard out")
	flag.StringVar(&templates, "templates", templates, "set the templates directory to scan")

	// We're ready to process args
	flag.Parse()
	flag.Args()

	in := os.Stdin
	out := os.Stdout
	eout := os.Stderr

	if showHelp {
		fmt.Fprintf(out, "%s\n", fmtHelp(helpText, appName, version, releaseDate, releaseHash))
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintf(out, "%s\n", newt.LicenseText)
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintf(out, "%s %s %s\n", appName, version, releaseHash)
		os.Exit(0)
	}
	os.Exit(newt.RunMustache(in, out, eout, port, timeout, verbose, templates))
}
