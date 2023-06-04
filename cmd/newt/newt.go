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

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-dry-run
: Load configuration and routes CSV but don't start web service


# CONFIGURATION

The three things {app_name} needs to know to run are port number,
where to find the "routes" CSV file and a list of any POSIX environment
variables to import and make available inside the router.

{app_mame} can be configured via a POSIX environment.

~~~
NEWT_PORT="8000"
NEWT_ROUTES="routes.csv"
NEWT_ENV="DB_NAME;DB_USER;DB_PASSWORD"
export NEWT_PORT
export NEWT_ROUTES
export NEWT_ENV
~~~

It can also be configured using a configuration file.


~~~
newt_port = "8000"
newt_routes = "routes.csv"
newt_env = [ "DB_NAME", "DB_USER", "DB_PASSWORD" ]
~~~

The environment will load first then the configuration file if
provided. The configuration file takes presidence to the environment.

{app_name} does not have secrets but could use secrets passed
in via the environment. This allows your routes CSV file to be safely
saved along side your SQL source code for your Newt Project.

# Routing data

For {app_name} to function as a data router is needs information
about which requests will be serviced and how to map them to a
JSON data API source before (optionally) sending to Pandoc.

The routes are held in CSV file with the following columns

req_path
: This is the URL path to watch for incoming requests, it may be a literal path or one containing variable declarations used in forming a api URL call.

req_method
: This is the HTTP method to listen for. Maybe "GET", "POST", "PUT", "PATCH" or "DELETE"

api_url
: This is the URL used to connect to the JSON data source (e.g. PostgREST, Solr, Elasticsearch). It may contain variables defined in the request path.

api_method
: This is the HTTP method used to access the JSON data source. Maybe "OPTIONS", "GET", "POST", "PUT", "PATCH" or "DELETE"

api_content_type
: This is the HTTP content type string to send with your JSON data source request, typically it is "application/json". 

pandoc_template
: If included Newt will load the Pandoc template file into memory and use it when results are returned from a JSON data source.

res_headers
: This is any additional HTTP headers you want to send back to the client.

# Route DSL

A simple domain specific language (DSL) can be used to define values taken
from a request path and used again to form a JSON data API URL. Variables can be defined in each of the request path's directory name(s), file basename and file extension. The variable is defined by an opening double curly bracket, the variable name, a space, a type and closing double curly brackets.

~~~
/blog/{{yr Year}}/{{mo Month}}/{{dy Day}}/{{title-slug String}}
/blog/{{yr Year}}/{{mo Month}}/{{dy Day}}/{{title-slug Basename}}{{ext Extname}}
~~~

In the first line the variables defined are "yr" of type "Year", "mo" of type "Month", "dy" of type "Day", "title-slug" of type "String". In the second line the "title-slug" is of type "Basename" (i.e. filename without an extension) and "ext" is of type "Extname" (i.e. the file extension).

In this prototype phase there are a very limited number of variables types
supported. This is likely to grow overtime if the prototype is successful.

## variable types

String
: Any sequence of characters except "/" which delimits the directory parts

Year
: A four digit year (e.g. 2023)

Month
: A two digit month (e.g. "01" for January, "10" for October)

Day
: A two digit day (e.g. "01" for the first, "11" for the eleventh)

Basename
: A file's basename (filename without an extension)

Extname
: A file's extension (e.g. ".html", ".txt", ".rss", ".js")

Isbn10
: An ten digit ISBN

Isbn13
: A thirteen digit ISBN

Isbn
: An ISBN (either 10 ro 13 digit)

Issn
: An ISSN

DOI
: A DOI (digital object identifier)

Isni
: An ISNI

ORCID
: An ORCID identifier
 

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
	dryRun      bool
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
	flag.BoolVar(&dryRun, "dry-run", false, "evaluate configuration and routes but don't start web service")

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
	os.Exit(newt.Run(in, out, eout, args, dryRun))
}
