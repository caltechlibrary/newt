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

{app_name} [OPTIONS] [HTDOCS]

# DESCRIPTION

**{app_name}** is a simple static web server running on localhost. The default port
is 8000 and the current working directory is the default document root. There is no
configuration aside from providing a directory path, optionally changing the port and
setting a verbose mode useful for debugging requests.

This is a minimal web server. No http, no embedded programming languages. No remapping
content types, redirect or other fancy stuff. It provides a quick way to allow your
static content to be available to your web browser over http for development purposes.

Dot files are not served out.

# OPTIONS

The following options are supported by **{app_name}**.

-h, -help
: display this help message

-license
: display the software license

-version
: display version information

-port
: set the port the web server listens on

-verbose
: show verbose logging of requests, e.g. contents of a POST

# EXAMPLE

In the example below the web server would listen for `+"`"+`http://localhost:8080`+"`"+`
and respond with the content in htdocs.

~~~shell
ws -port 8080 htdocs
~~~

An example of using the static file server to debug a form submission by showing
what the forms submits in the log output.

~~~shell
ws -verbose -port 8080 htdocs
~~~


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
	port, verbose := 0, false
	flag.IntVar(&port, "port", port, "set the port number for service")
	flag.BoolVar(&verbose, "verbose", verbose, "display detail logging of requests")

	// We're ready to process args
	flag.Parse()
	args := flag.Args()

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
	os.Exit(newt.RunStaticWebServer(in, out, eout, args, port, verbose))
}
