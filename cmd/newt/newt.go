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
title: "{app_name}(1) user manual | Version 0.0.1 f75250d"
pubDate: 2023-06-03
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
YAML file.  {app_name} first sends requests to Postgres+PostgREST
then processes the results via Pandoc running as a web service.
While {app_name} was created to work specifically with PostgREST
it can actually talk to any JSON data source that can be specified
in URL (e.g. Solr, Elasticsearch, Opensearch). 

Before contacting the JSON data source the request URL are validated
including and extracting any variable contents embedded in the request.
The variables can then be used in forming the URL to make the request
to the JSON data source.  If a web form is part of the request the
defined variables and their types are used to vet and validate the
request before submitting it to the JSON data source.

When the data source replies the results can be fed through Pandoc server 
or returned directly to the browser depending on how the route is
configured.

Newt also can provide static content hosting so that related HTML,
CSS, JavaScript and other page assets can be integrated into the 
response to the request.

Newt's configuration uses a declaritive model expressed in YAML.
It can also allow environment variables read at start up to be
part of the data for mapping JSON data source requests or used
referenced by a Pandoc template.

This goal of Newt Project is to be able to assemble an entire back-end
from off the self services only specify data modeling and end point
definitions using SQL and a Postgres database. Reducing the back-end
to SQL may simplify application management (it reduces it to a
database administrator activity) and free up developer time to focus
more on front end development and human interaction. It is also
hoped that focusing the back-end on a declarative model will allow for
a more consistent and reliable back-end.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-dry-run
: Load YAML configuration and report any errors found


# CONFIGURATION

The three things {app_name} needs to know to run are port number,
where to find the "routes" YAML file and a list of any POSIX environment
variables to import and make available inside the router.

{app_name} can be configured via a POSIX environment.

~~~
NEWT_PORT="8000"
NEWT_ROUTES="routes.yaml"
NEWT_ENV="DB_NAME;DB_USER;DB_PASSWORD"
export NEWT_PORT
export NEWT_ROUTES
export NEWT_ENV
~~~

The environment variables can also be configured 
in the Newt YAML file.

~~~
{app_name}_port: "8000"
{app_name}_env: [ "DB_NAME", "DB_USER", "DB_PASSWORD" ]
{app_name}_htdocs: "htdocs"
~~~

The environment will load first then the configuration file if
provided. The configuration file takes precedence over the 
environment.

{app_name} does not have secrets but could use secrets passed
in via the environment. This allows your routes in the YAML file
to be safely saved along side your SQL source code for
your Newt Project (e.g. your YAML and SQL files can be checked
safely into your source code repository without saving secrets)

# Routing data

For {app_name} to function as a data router it needs information
about which requests will be serviced and how to map them to a
JSON data API source before (optionally) sending to Pandoc.

The routes are held in YAML file under the "routes" attribute.
The following attributes are supported for a route.

var
: One or more variable name and type pairs used in request path or form data. The types allow for data validation.

req_path
: This is the URL path to watch for incoming requests, it may be a literal path or one containing variable declarations used in forming a API URL call.

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

# Defining variables

Each route can have an associated variables available to form a JSON data API request. The defined variables also are used to validate webform input for the request. Only the defined variables will be passed on to the JSON data API and the Pandoc template if specified. Here's an example "var" definitiondefining three form variables for a route, they are "bird", "place" and "sighted" with the types "String", "String" and "Date".

~~~
var: { "bird": "String", "place": "String", "sighted": "Date" }
~~~

If a web browser injected additional form values they would not get passed along via the JSON data API request, they would be ignored. This is part of the declaritive approach for defining Newt's behavior.

# Variable types

# Route DSL

Variables that are defined for a route can be extracted from a request path where they represent a path element (e.g. directory, filename or extension). The extracted variables can then be used in forming the JSON data API request and also inside a Pandoc template as part of the "data" values passed to it.  The variable is referenced by a dollar sign and open curly bracket, the variable name and a closing closing curly brackets. This is similar to how
Pandoc template variables are represented.


~~~
/blog/${yr}/${mo}/${dy}/${title-slug}
~~~

The "var" attribute in the route would look like 

~~~
var: { "yr": "Year", "mo": "Month", "dy": "Day", "title-slug": "String" }
~~~

In the above example the values "yr", "mo", "dy" and "title-slug" would
be extracted from a request path. These might then be used to form an
API request path along with environment variables imported by the YAML
file.

~~~
https://localhost:3000/blog?date=${yr}-${mo}-{dy}&title-slug=${title-slug}
~~~

The resulting data would be bound to the variable "data" and passed to
 Pandoc to be processed along with the appropriate template.

There are times you might need to treat the "filename" part of a path as a file's basename and extension.  Two types data types handle that. So for a "var" defined like

~~~
var: { "yr": "Year", "mo": "Month", "dy": "Day", "title-slug": "BaseName", "ext": "Extname" }
~~~

The request URL pattern could like like

~~~
/blog/${yr}/${mo}/${dy}/${title-slug}${ext}
~~~

The related JSON data source URL might look something like

~~~
https://localhost:3000/blog?date=${yr}-${mo}-{dy}&title-slug=${title-slug}&format=${ext}
~~~


NOTE: that "Basename" and "Extname" only make sense in the context of a path. If those same values are used in a form they will be validated as a string only.

In this prototype phase there are a very limited number of variables types supported. This is likely to grow and to change overtime if the prototype is successful.

## variable types

String
: Any sequence of characters. If the variabe is embedded in a path then "/" will be used to delimited path parts and would not be passed into the variables value.

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
	export NEWT_ROUTES="{app_name}.yaml"
	export NEWT_ENV="DB_USER;DB_PASSWORD"
~~~

Configuration from a YAML file called "{app_name}.yaml"

~~~
{app_name} {app_name}.yaml
~~~

An example of a YAML file describing blog display routes.

~~~
htdocs: htdocs
routes:
	- var: [ "yr": "Year", "mo": "Month", "dy": "Day" }
	  req_path: "/blog/${yr}/${mo}//${dy}"
	  req_method: GET
	  api_url: "http://localhost:3000/posts?year=${yr}&month=${mo}&day=${dy},posts.tmpl"
	  api_method: GET
	  api_content_type: "application/json"
	  pandoc_template: article_list.tmpl
	  res_headers: { "content-type": "text/html" }
	- var: [ "yr": "Year", "mo": "Month", "dy": "Day" }
	  req_path: "/blog/${yr}/${mo}//${dy}/${title-slug}"
	  req_method: GET
	  api_url": "http://localhost:3000/posts?year=${yr}&month=${mo}&day=${dy}&title-slug=${title-slug}"
	  pandoc_template: "article.tmpl"
	  res_headers: { "content-type": "text/html" }
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
