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

{app_name} [OPTIONS] ACTION [YAML_FILE]

# DESCRIPTION

**{app_name}** is a developer tool. It can set up and run your Newt application
during development.  **{app_name}** supports the "config", "check", "model", "generate",
and "run" actions. The "config" command is used when you are starting a
Newt Project. It is an interactive command prompting for various choices regarding
the application you want to create.  

The "check" will analyze the the YAML file and report results of the analyze and as
well as validate the YAML syntax.

The "model" is used to manage your data models. It is interactive like "config".

The "generate" will run the Newt Generator to create the SQL, PostgREST configuration
and Mustache templates. It'll also update the Newt YAML file's routes and templates
properties as needed for the generated content. 

The "run" action will run Newt Router, Newt Mustache and PostgREST
for testing and development.  This allows you to quick run and stop the services
from one command. 

If a Newt YAML filename isn't supplied when you invoke **{app_name}** the application
will look in the current directory for a file with the directory's name and ".yaml" 
extension and if that is not found it'll look for "app.yaml".

If you are working in the root of a Git repository when you run "config" you will
get a warning about which files should be added to your `+"`"+`.gitignore`+"`"+`.

# OPTIONS

The following options are supported by **{app_name}**.

-help
: display this help message

-license
: display the software license

-version
: display version information

-verbose
: If set provide verbose debugging output for requests

# ACTION

config [YAML_FILE]
: this will create or refresh your Newt YAML file based on a set of interactive
questions. It will suggest updates to your `+"`"+`.gitignore`+"`"+`.

model [YAML_FILE]
: this will read the current Newt YAML file and run the interactive modeler updating the Newt YAML file.

check [YAML_FILE]
: analyze the Newt YAML file and report problems if found.

generate [YAML_FILE]
: this is used to generate your SQL, PostgREST configuration and Mustache templates
based on the contents of your Newt YAML file.

run [YAML_FILE]
: this will run the defined services in the application attribute of the Newt YAML file.
This is intended for use in development. In a production setting you'd setup the individual
services to run from systemd or config as services.

# YAML_FILE 

**{app_name}** is configured in a YAML file. What is described below is a summary of 
YAML syntax use in a Newt project that uses all of the Newt programs.

## Top level properties

These are the top level properties in YAML files.

applications
: (optional) holds the run time configuration used by the Newt applications.

models
: (required by newt generator) This holds the description of the data models in your application.  Each model uses HTML 5 element descriptions which can be set using the interactive `+"`"+`{app_name} model`+"`"+` command. 

routes
: (required by newt router) This holds the routes for the data pipeline (e.g. JSON API and template engine sequence)

templates
: (required by newt mustache)

## The applications property

The applications properties are optional. Some maybe set via command line. See Newt application's manual page for specific ones. These properties lets you override the default settings of Newt programs.

router
: this contains configuration for the Newt Router, i.e. port and htdocs

generator
: this contains configuration for code generation. e.g. namespace used by Postgres and PostgREST.

mustache
: this contains configuration for Newt Mustache, i.e. port

options
: holds key value pairs of which can be referenced in the values of models, routes and templates.

environment
: (optional: newtrouter, newtmustache) this is a list of operating system environment that will be available to routes. This is used to pass in secrets (e.g. credentials) need in the pipeline

### Each application may have one or more of these properies.

namespace
: (postgres) uses this in the SQL generated for setting up Postgres+PostgREST

port
: (all) Port number to used for Newt web service running on localhost

htdocs
: (router) Directory that holds your application's static content

timeout
: (router, mustache) The time in seconds to timeout a HTTP transaction

dsn
: (postgres) The data source name (URI for database connection string). Using to connect to Postgres

app_path
: (postgrest) The full path to the application (if it is not in the PATH environment variable)

conf_app
: (postgres) Configuration path for the application.

## the "routes" property

Routes hosts a list of request descriptions and their data pipelines. This property is only used by Newt router and Newt code generator.

### a route object

`+"`"+`id`+"`"+`
: (required) This identifies the pipeline. It maybe used in code generation. It must conform to variable name rules[^21]. You may have more than one route associated with an identifier but they may not share the same method. This allows us to group routes by actions. E.g. creating an object would include a GET to retrieve the web form and a POST to handle the submission. Both should use the same id.

`+"`"+`description`+"`"+`
: (optional, encouraged) This is a human readable description of what you're trying to accomplish in this specific route. It may be used in comments or by documentation generators.

`+"`"+`request [METHOD ][PATH]`+"`"+`
: (required) This is a string that expresses the HTTP method and URL path to used to trigger running the data pipeline. If METHOD is not provided it will match using just the path. This is probably NOT what you want. You can express embedded variables in the PATH element. This is done by using single curl braces around a variable name. E.g. `+"`"+`GET /items/{item_id}`+"`"+` would make `+"`"+`item_id`+"`"+` available in building your service paths in the pipeline. The pattern takes up a whole path segment so `+"`"+`/blog/{year}-{month}-{day}`+"`"+` would not work but `+"`"+`/blog/{year}/{month}/{day}`+"`"+` would capture the individual elements. The Newt router sits closely on top of the Go 1.22 HTTP package route handling. For the details on how Go 1.22 and above request handlers and patterns form see See <https://tip.golang.org/doc/go1.22#enhanced_routing_patterns> and <https://pkg.go.dev/net/http#hdr-Patterns> for explanations.

`+"`"+`pipeline`+"`"+`
: (required) this is a list of URLs to one or more web services visible on localhost. The first stage to fail will abort the pipeline returning an HTTP error status. If done fail then the result of the last stage it returned to the requesting browser.

`+"`"+`debug`+"`"+`
: (optional) if set to true the `+"`"+`newt`+"`"+` service will log verbose results to standard out for this specific pipeline

#### a pipeline object

A pipeline is a list of web services containing a type, URL, method and content types

`+"`"+`service [METHOD ][URL]`+"`"+`
: (required) The HTTP method is included in the URL The URL to be used to contact the web service, may contain embedded variable references drawn from the request path as well as those passed in through `+"`"+`.application.environment`+"`"+`.  All the elements extracted from the elements derived from the request path are passed through strings. These are then used to construct a simple key-value object of variable names and objects which are then passed through the Mustache template representing the target service URL. 

`+"`"+`description`+"`"+`
: (optional, encouraged) This is a description of what this stage of the pipe does. It is used when debug is true in the log output and in program documentation.

`+"`"+`timeout`+"`"+`
: (optional) Set the timeout in seconds for receiving a response from the web server. Remember the time spent at each stage is the cumulative time your browser is waiting for a response. For this reason you may want to set the timeout to a small number.


## the "models" property

Models holds a list of individual models used by our data pipelines. The models are describe by a list of HTML5 form elements which intern are described by their attributes.

### a model object

The model object was inspired by GitHub YAML issue template syntax. It is much simpler and more focus on
representing HTML5 webform elements and their attributes directly. HTML5 form elements imply a mapping to SQL
data types.

id
: (required) this is the name identifying the model. It must conform to variable name rules[^21]

description
: (optional) A human description of the model, It will appear in the web form and SQL components generated from the model

elements
: (required) A a list of input types. Each input type maps to columns in SQL, input element in web forms and or HTML elements in read only pages

#### a model's input types

id
: (required) an identifier for the element. Must conform to variable name rules[^21]. It is used to SQL as a column name and in web forms for the input property.

type
: (required) Identifies the type of elements (input, textarea, markdown, checkbox, dropdown).

attributes
: (optional) A key-value list that define properties of the element. These used in rendering the element in SQL or HTML. These correspond with the attributess you'd find in HTML 5 input elements.

primary_key
: (boolean, defaults to false) if true then Newt will use this element as the identifier in the SQL code and routes. You can only have one element used as the object identifier and it must be unique.

## input types

Inputs must map to their HTML 5 form element types. Here is an example of HTML 5 input types.

text
: A single line text field. This conforms to value input types in HTML 5 and can be expressed using the CSS selector notation. E.g. `+"`"+`input[type=data]`+"`"+` would be a date type. This would result in a date column type in SQL, a date input type in HTML forms and in formatting other HTML elements for display.

textarea
: A multi-line text field

checkboxes
: A checkbox element. In SQL if the checkbox is exclusive (e.g. a radio button) then the result is stored in a single column, if multiple checks are allowed it is stored as a JSON Array column.

More comple types include 

radio
: For generating a radio button list

select
: A dropdown menu. In SQL this could render as an enumerated type. In HTML it would render as a drop down list

auto_complete
: A input plus data list, useful for large vocabularies


NOTE: The file input type is explicitly NOT supported by Newt at this time. Newt may add additional types in
the future through aliasing input with regular expression validation or web components.

## Example Newt YAML file, "app.yaml" generated with the `+"`"+`{app_name} config app.yaml`+"`"+` and
`+"`"+`{app_name} model app.yaml`+` commands.

~~~yaml
applications:
  newtrouter:
    port: 8010
  newtmustache:
    port: 8011
  newtgenerator:
    namespace: app
  postgres:
    port: 5432
    dsn: postgres://{PGUSER}:{PGPASSWORD}@localhost:5432/app
  postgrest:
    app_path: postgrest
    conf_path: postgrest.conf
    port: 3000
  environment:
    - PGUSER
    - PGPASSWORD
models:
  - id: app
    description: This is where you would model your application data
    body:
      - type: input
        id: data_attribute
        attributes:
          description: This is an example input element
          name: data_attribute
          placeholdertext: ex. of placeholder text
          title: this is an example element in your model
        validations:
          required: true
routes:
  - id: app_create
    request: GET /app_create
    description: Handle retrieving the webform for app create
    pipeline:
      - service: POST http://localhost:8011/app_create
        description: Display a app for create
  - id: app_create
    request: POST /app_create
    description: Handle form submission for app create
    pipeline:
      - service: POST http://localhost:3000/rpc/app_create
        description: Access PostgREST API for app create
      - service: POST http://localhost:8011/app_create_response
        description: This is an result template for app create
  - id: app_update
    request: GET /app_update/{oid}
    description: Handle retrieving the webform for app update
    pipeline:
      - service: GET http://localhost:3000/rpc/app_read/{oid}
        description: Retrieve app from PostgREST API before update
      - service: POST http://localhost:8011/app_update
        description: Display a app for update
  - id: app_update
    request: POST /app_update
    description: Handle form submission for app update
    pipeline:
      - service: PUT http://localhost:3000/rpc/app_update/{oid}
        description: Access PostgREST API for app update
      - service: POST http://localhost:8011/app_update_response
        description: This is an result template for app update
  - id: app_delete
    request: GET /app_delete/{oid}
    description: Handle retrieving the webform for app delete
    pipeline:
      - service: GET http://localhost:3000/rpc/app_read/{oid}
        description: Retrieve app from PostgREST API before delete
      - service: POST http://localhost:8011/app_delete
        description: Display a app for delete
  - id: app_delete
    request: POST /app_delete
    description: Handle form submission for app delete
    pipeline:
      - service: DELETE http://localhost:3000/rpc/app_delete/{oid}
        description: Access PostgREST API for app delete
      - service: POST http://localhost:8011/app_delete_response
        description: This is an result template for app delete
  - id: app_read
    request: POST /app_read
    description: Retrieve object(s) for app read
    pipeline:
      - service: GET http://localhost:3000/rpc/app_read/{oid}
        description: Access PostgREST API for app read
      - service: POST http://localhost:8011/app_read
        description: This template handles app read
  - id: app_list
    request: POST /app_list
    description: Retrieve object(s) for app list
    pipeline:
      - service: GET http://localhost:3000/rpc/app_list
        description: Access PostgREST API for app list
      - service: POST http://localhost:8011/app_list
        description: This template handles app list
templates:
  - id: app_create
    request: /app_create
    template: app_create_form.tmpl
    description: Display a app for create
  - id: app_create
    request: /app_create_response
    template: app_create_response.tmpl
    description: This is an result template for app create
  - id: app_update
    request: /app_update
    template: app_update_form.tmpl
    description: Display a app for update
  - id: app_update
    request: /app_update_response
    template: app_update_response.tmpl
    description: This is an result template for app update
  - id: app_delete
    request: /app_delete
    template: app_delete_form.tmpl
    description: Display a app for delete
  - id: app_delete
    request: /app_delete_response
    template: app_delete_response.tmpl
    description: This is an result template for app delete
  - id: app_read
    request: /app_read
    template: app_read.tmpl
    description: This template handles app read
  - id: app_list
    request: /app_list
    template: app_list.tmpl
    description: This template handles app list
~~~


[^21]: variable numbers must start with a letter, may contain numbers but not spaces or punctuation except the underscore

## templates property

This property is used by Newt Mustache. It is ignore by Newt router and code generator.

templates
: (optional) this holds a list of template objects

### template object model

The template objects are used by Newt Mustache template engine. If you're not using it you can skip these.

`+"`"+`id`+"`"+`
: (required) is the template id. It can be used to relate one or more templates to an action. E.g. a web form an
submission result.

`+"`"+`request [PATH]`+"`"+`
: (required) This holds the request HTTP method and path. Note the HTTP method is missing as all request to Newt Mustache must be done using POST.

`+"`"+`template`+"`"+`
: (required) This is the path to the template associated with request. NOTE: Pandoc web service does not support partial templates. Mustache does support partial templates

`+"`"+`partials`+"`"+`
: (optional) A list of paths to partial Mustache templates used by `+"`"+`.template`+"`"+`.

`+"`"+`vocabulary`+"`"+`
: (opitonal) A path to a "vocabulary" file which is read in when the Newt Mustache is started and the contents are available to the template for processing.

`+"`"+`options`+"`"+`
: (optional) An object that can be merged in with JSON received for processing by your Mustache template

`+"`"+`debug`+"`"+`
: (optional) this turns on debugging output for this template

Example of newtmustache YAML:

~~~yaml
applications:
  newtmustache:
    port: 8012
templates:
  - id: hello
    request: POST /hello/{name}
    template: testdata/simple.mustache
  - id: hello
    request: POST /hello
    template: testdata/simple.mustache
    options:
      name: Universe
  - id: hi
    request: POST /hi/{name}
    template: testdata/hithere.html
    partials:
      - testdata/name.mustache
    debug: true
  - id: hi
    request: POST /hi
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
	verbose := false
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
	if len(args) > 0 && args[0] == "help" {
		fmt.Fprintf(out, "%s\n", fmtHelp(helpText, appName, version, releaseDate, releaseHash))
		os.Exit(0)
	}
	os.Exit(newt.RunNewt(in, out, eout, args, verbose))
}
