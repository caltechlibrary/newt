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

{app_name} [OPTIONS] YAML_CONFIG_FILE

# DESCRIPTION

**{app_name}** is a web service that provides a Mustache template rendering inspired by Pandoc server.

Unlike Pandoc web server, `+"`"+`{app_name}`+"`"+` expects a YAML_CONFIG_FILE. The format is
described below. That file specifics the request to template mapping along with any ancillary information
to merge into the submitted object for the template to process. In additional to options expressed in
the configuration the pattern describing the request can also be merged into the JSON objects the template
will process.

`+"`"+`{app_name}`+"`"+` expects a JSON object to process. This means you should normally use a POST
in defining your request pattern to be match.  If no method is specified then a POST method will be
assumed.

The request expression is based on the Go 1.22 patterns used in it's `+"`"+`http`+"`"+` package.
<https://pkg.go.dev/net/http@master#hdr-Patterns>. The only modification is if not method is included
a POST instead of a GET will be assumed.

Like Pandoc web service `+"`"+`{app_name}`+"`"+` does not normally log requests. It's a quick transaction.
If you want to debug your templates use the verbose option or turn on debug for specific requests.

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

-timeout SECONDS
: Timeout in seconds, after which a template rendering is aborted.  Default: 3.

-verbose
: If set provide verbose debugging output for requests

# The templates

Mustache templates are documented at <https://mustache.github.io>. The template engine
used is based on Go package <https://github.com/cbroglie/mustache>.

## Features

- Newt Mustache only runs on localhost at the designated port (default is 3032).
- Templates are read in at startup and are retained in memory bound to the request path.
- No addition reads are performed once the web service starts listening.
- Patterns expressed in the request definitions are available in the object passed to the template

# YAML_CONFIG_FILE

This is a list of the Newt YAML syntax relevant to **{app_name}**.

## Top level properties

These are the top level properties in YAML files.

applications
: (optional) holds the run time configuration used by the Newt applications.

templates
: (required by newtmustache)

## The applications property

The applications properties are optional. Some maybe set via command line. See Newt application's manual page for specific ones. These properties lets you override the default settings of Newt programs.

newtmustache
: this contains configuration for Newt Mustache, i.e. port

options
: holds key value pairs of which can be referenced in the values of models, routes and templates.

### newtmustache settings

port
: (all) Port number to used for Newt web service running on localhost

### the "routes" property

Routes hosts a list of request descriptions and their data pipelines. This property is only used by Newt router and Newt code generator.

## templates property

This property is used by Newt Mustache. It is ignore by Newt router and code generator.

templates
: (optional: newtmustache) this holds a list of template objects

### template object model

The template objects are used by Newt Mustache template engine. If you're not using it you can skip these.

`+"`"+`request [METHOD ][PATH]`+"`"+`
: (required) This holds the request HTTP method and path. If the HTTP method is missing a POST is assumed

`+"`"+`template`+"`"+`
: (required: newtmustache only) This is the path to the template associated with request. NOTE: Pandoc web service does not support partial templates. Mustache does support partial templates

`+"`"+`partials`+"`"+`
: (optional, newtmustache only) A list of paths to partial Mustache templates used by `+"`"+`.template`+"`"+`.

`+"`"+`options`+"`"+`
: (optional, newtmustache only) An object that can be merged in with JSON received for processing by your Mustache template

`+"`"+`debug`+"`"+`
: (optional) this turns on debugging output for this template

# EXAMPLES

Example of newtmustache YAML:

~~~yaml
applications:
  newtmustache:
    port: 8012
templates:
  - request: GET /hello/{name}
    template: testdata/simple.mustache
  - request: GET /hello
    template: testdata/simple.mustache
    options:
      name: Universe
  - request: GET /hi/{name}
    template: testdata/hithere.html
    partials:
      - testdata/name.mustache
    debug: true
  - request: GET /hi
    template: testdata/hithere.html
    partials:
      - testdata/name.mustache
    options:
      name: Universe
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
	port, timeout, verbose := 0, 0, false
	flag.IntVar(&port, "port", port, fmt.Sprintf("set the port %s will listen on", appName))
	flag.IntVar(&timeout, "timeout", timeout, "timeout for requests")
	flag.BoolVar(&verbose, "verbose", verbose, "display template debugging and request information")

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
	os.Exit(newt.RunNewtMustache(in, out, eout, args, port, timeout, verbose))
}
