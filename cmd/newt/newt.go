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
title: "{app_name}(1) user manual | Version {version} {release_hash}"
pubDate: {release_date}
author: "R. S. Doiel"
---

# NAME

{app_name}

# SYNOPSIS

{app_name} [CONFIG_FILE]

# DESCRIPTION

*{app_name}* is a microservice designed to work along side Postgres,
PostgREST, and Pandoc server. It provides URL routing and data flow
between the microservices based on a list of "routes" described in a
CSV file.  {app_name} is part of the Newt Project which is exploring
building web services, applications and sites using SQL for data modeling
and define back-end service behaviors along with Pandoc templates used to
generate HTML consumed by the web browser.  {app_name} supports data
hosted in Postgres databases via PostgREST JSON API as well as static
files contained in an "htdocs" directory (e.g. HTML, CSS, JavaScript,
and image assets). 

This goal of Newt Project is to be able to assemble an entire backend
from off the self services only requiring data modeling and end point
definitions using SQL and a Postgres database. Reducing the back-end
to SQL may simplify application management (it reduces it to a
database administrator activity) and free up developer time to focus
more on front end development and human interaction. It is also
hoped that focusing the back-end on a declaritive model will allow for
a more consistent and reliable back-end.

# CONFIGURATION

The three things {app_name} needs to know to run are port number,
where to find the "routes" CSV file and a list of any POSIX environment
variables to import and make available inside the router.

{app_mame} can be configured via a POSIX environment.

~~~
NEWT_PORT="3030"
NEWT_ROUTES="routes.csv"
NEWT_ENV="DB_USER;DB_PASSWORD"
export NEWT_PORT
export NEWT_ROUTES
export NEWT_ENV
~~~

It can also be configured using a configuration file.


~~~
newt_port = "3030"
newt_routes = "routes.csv"
newt_env = [ "DB_USER", "DB_PASSWORD" ]
~~~


The environment will load first then the configuration file if
provided. The configuration file takes presidence to the environment.

{app_name} does not have secrets but could use secrets passed
in via the environment. This allows your routes CSV file to be safely
saved along side your SQL source code for your Newt Project.

# EXAMPLES

Configuration from the environment

~~~
	export NEWT_PORT="3030"
	export NEWT_ROUTES="routes.csv"
	export NEWT_ENV="DB_USER;DB_PASSWORD"
	{app_name}
~~~

Configuration from a YAML file called "newt.yaml"

~~~
{app_name} newt.yaml
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
	os.Exit(newt.Run(in, out, eout, args))
}
