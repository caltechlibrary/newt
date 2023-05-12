package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	// Caltech Library Packages
	"github.com/rsdoiel/newt"
)

var (
	helpText = `---
title: "{app_name}(1) user manual | Version {version}"
pubDate: 2023-05-12
author: "R. S. Doiel"
---

# NAME

{app_name}

# SYNOPSIS

{app_name} [CONFIG_FILE]

# DESCRIPTION

*{app_name}* is a microservice designed to work along side Postgres,
PostgREST, and Pandoc server. It provides URL routing and data flow
between the microserves based on a simple configuration table holding
routing information. It is part of the Newt Project which is exploring
building web services, applications and sites using SQL for data modeling
and define back-end service behaviors along with Pandoc templates used to
generate HTML consumed by the web browser.  {app_name} supports data
hosted in Postgres databases via PostgREST JSON API as well as static
files contained in an "htdocs" directory (e.g. HTML, CSS, JavaScript,
and image assets). 

This goel of Newt project is to be able to assemble an entire backend
from off the self services only requiring data modeling and route
definitions using SQL and a Postgres database. Reducing the back-end
to SQL may simplify application management (it reduces it to a
database administrator activity) and free up developer time to focus
more on front end development and human interaction. It is also
hoped that focusing the back-end on a declaritive model will allow for
a more consistent and reliable back-end.

# CONFIGURATION

{app_name} can be configured through the environment or through
a PostgREST is configuration file. It adds only an optional
uri for the Pandoc server (if used) when it runs on a non-standard
port.

~~~
db-uri = "postgres://birds:my_secret_password@localhost:5432/birds"
db-schemas = "birds"
db-anon-role = "birds_anonymous"
# This is used by Newt to know where to find the Pandoc server
# on localhost.
pandoc-server-port = "3030"
~~~

# EXAMPLE

~~~
	{app_name} postgrest.conf
~~~

`

	showHelp    bool
	showLicense bool
	showVersion bool
)

func fmtTxt(txt string, appName string, version string) string {
	return strings.ReplaceAll(strings.ReplaceAll(txt, `{app_name}`, appName), `{version}`, version)
}

func main() {
	appName := path.Base(os.Args[0])

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
		fmt.Fprintf(out, "%s\n", fmtTxt(helpText, appName, newt.Version))
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintf(out, "%s\n", newt.LicenseText)
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintf(out, "%s %s\n", appName, newt.Version)
		os.Exit(0)
	}
	os.Exit(newt.Run(in, out, eout, args))
}
