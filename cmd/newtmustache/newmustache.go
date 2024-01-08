package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	// Caltech Library Packages
	"github.com/rsdoiel/newt"
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

{app_name}

# DESCRIPTION

**{app_name}** is a microservice that provides a Mustache template rendering inspired by Pandoc server.

There is no configuration file. There are command line options to specify port and timeout.

# API

## Root endpoint

The root (/) endpoint accepts only POST method requests. Other methods will yield a appropriate http status code.

##  Request

The body of the POST request should be a JSON object, with the following fields.  Template is always required. It must
also include at least one of either metadata or variables. If both are set then the maps are merged before processing
with the template.

template (string)
: String contents of a document template, see Mustache <https://mustache.github.io/mustache.5.html>.

data (JSON map)
: String-valued metadata.

content_type (optional, default is text/plain)
: Set the rendered content type, default is text/plain. You probably want text/html for web pages.

## Response

It returns a response body and headers set by the default http writer provided in the http go package.
If there request can't be fullfilled then an http status code and text will be returned.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

--port NUM
: HTTP port on which to run the server.  Default: 3030.

--timeout SECONDS
: Timeout in seconds, after which a conversion is killed.  Default: 2.

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
	port, timeout, verbose := "3030", 2, false
	flag.StringVar(&port, "port", port, "port number to listen on")
	flag.IntVar(&timeout, "timeout", timeout, "timeout for requests")
	flag.BoolVar(&verbose, "verbose", verbose, "display template debugging and request information to standard out")

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
	os.Exit(newt.RunMustache(in, out, eout, port, timeout, verbose))
}
