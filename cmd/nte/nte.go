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
title: {app_name}(1) user manual | {version} {release_hash}
pubDate: {release_date}
author: R. S. Doiel
---

# NAME

{app_name}

# SYNOPSIS

{app_name} [OPTIONS] YAML_CONFIG_FILE

# DESCRIPTION

**{app_name}** is a web service that provides a template rendering engine inspired
by Pandoc server.

Unlike Pandoc web server, `+"`"+`{app_name}`+"`"+` expects a YAML configuration file.
The format is described below. That file specifies the runtime configuration. It specifies
the request path to template mapping. It can also specify ancillary information made
available to the template associated with the request path and template.

The `+"`"+`{app_name}`+"`"+` template engine listens for a POST requests with JSON encoded data.
It  checks requested path to see if that matches the request path described in the YAML
file. If there is a match it processes the request returning the rendered results matched with
 any data found in the POST. `+"`"+`{app_name}`+"`"+`.

The content of the POST is passed to the template as `+"`"+`.body`+"`"+`, applications options
are merged into a `+"`"+`.document`+"`"+` object along with any specified vocabularies.
Finally if you've defined a variables in the path to the tempalte those are provided via
the `+"`"+`.vars`+"`"+` property.

If you use a GET request then the unprocessed referenced template is returned (minus partials,
layouts, etc).

# OPTIONS

The following options are supported by **{app_name}**.

`+"`"+`-h`+"`"+`
: display this help message

`+"`"+`-license`+"`"+`
: display the software license

`+"`"+`-version`+"`"+`
: display version information

`+"`"+`-port NUMBER`+"`"+`
: (default is port is 3032) Set the port number to listen on

`+"`"+`-base-dir PATH`+"`"+`
: set the base directory path (where you have your templates).

`+"`"+`-timeout SECONDS`+"`"+`
: Timeout in seconds, after which a template rendering is aborted.  Default: 3.

`+"`"+`-verbose`+"`"+`
: If set provide verbose debugging output for requests

# The templates

The template engine supports the [Handlebars](https://handlebarsjs.com) template langauge
which is largely a superset of Mustache templates documented at <https://mustache.github.io>.
The template engine used is based on Go package <github.com/aymerick/raymond>.

## Features

- Newt template engine only runs on localhost at a designated port (default is 8011).
- Templates are read in at startup and are retained in memory bound to the request path.
- JSON data is provided to the template in a `+"`"+`.body`+"`"+` object.
- Vocabulary files are read in at startup and bound to the request path and propogated
  to the template via the `+"+"+`.document`+"`"+` object.
- Variables found expressed in the request path are available in the `+"`"+`.vars`+"`"+`
  passed to the template.
- Except for variables no addition reads are performed once the web service starts listening.

# YAML_CONFIG_FILE

This is a list of the Newt YAML syntax relevant to **{app_name}**.

## Top level properties

These are the top level properties in YAML files.

applications
: (required) holds the run time configuration used by the Newt applications.

templates
: (required) holds a list of template objects

## The applications property

template_engine
: (requred) this contains configuration for Newt template engine, e.g. port, base_dir, ext_name.

### template engine properties

port
: (required) port number to used for to ued for Newt Template Engine

base_dir
: (required) base directory holding the primary templates

partial_dir
: (optional) the sub directory holding the partial templates

layout_Dir
: (optional) the sub directory holding the layouts

default_layout
: (optional) the default layout

`+"`"+`vocabularies`+"`"+`
: (optional) this holds a map of vocabulary name to vocabulary filename. A
vocabulary file is a YAML file that is made available in templates via the
`+"`"+`.document`+"`"+` object. It can be used to provide common document
attributes between a set of templates.

helpers
: (optional) this holds a map of handlebars helpers


## templates property

This property is used by Newt template engine. It provides a list of
template objects.

### template object

The template objects are used by Newt template engine. If you're not using it you can skip these.

`+"`"+`request PATH`+"`"+`
: (required) This holds the request URL's path. `+"`"+`{app_name}`+"`"+` only listens for POST method.

`+"`"+`template`+"`"+`
: (required) This is the name of the primary template (without file extnesion). The primary
template may also include partials and those will be read from the partials sub directory
defined in the template engine property.

`+"`"+`document`+"`"+`
: this will provide template specific data merged with the any vocaluaries
defined in template_engine property.

`+"`"+`debug`+"`"+`
: (optional) this turns on debugging output for this template

# EXAMPLES

Example of newtmustache YAML that only runs the template engine by itself.
The paths are used to provide template content.

~~~yaml
applications:
  template_engine:
    port: 8011
	# this is the path to the primary templates
	base_dir: testadata
templates:
  - id: hello
    request: /hello/{name}
    template: simple
  - id: hello
    request: /hello
    template: simple
    document:
      name: Universe
  - id: hi
    request: /hi/{name}
    template: hithere
    debug: true
  - id: hi
    request: /hi
    template: hithere
    document:
      name: Universe
~~~

NOTE: the template name doesn't require the extension since that is set at the 
template engine level.

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
	os.Exit(newt.RunTemplateEngine(in, out, eout, args, port, timeout, verbose))
}
