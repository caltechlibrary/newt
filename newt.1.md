---
title: newt(1) user manual | 0.0.8 bc1bb84
pubDate: 2024-05-01
author: R. S. Doiel
---

# NAME

newt

# SYNOPSIS

newt [OPTIONS] ACTION [YAML_FILE]

# DESCRIPTION

**newt** is a developer tool. It can set up and run your Newt application
during development.  **newt** supports the "init", "check", "model", "generate",
"run" and "ws" actions. The "init" command is used when you are starting a
Newt Project. It is an interactive command prompting for various choices regarding
the application you want to create.  

The "check" will analyze the the YAML file and report results of the analyze and as
well as validate the YAML syntax.

The "model" is used to manage your data models. It is interactive like "init".

The "generate" will run the Newt Generator to create the SQL, PostgREST configuration
and Mustache templates. It'll also update the Newt YAML file's routes and templates
properties as needed for the generated content. 

The "run" action will run Newt Router, Newt Mustache and PostgREST
for testing and development.  This allows you to quick run and stop the services
from one command. 

The "ws" action runs a simple static web service. It is useful when you are debugging
static web assets for your project.


If a Newt YAML filename isn't supplied when you invoke **newt** the application
will look in the current directory for a file with the directory's name and ".yaml" 
extension and if that is not found it'll look for "app.yaml".

If you are working in the root of a Git repository when you run "init" you will
get a warning about which files should be added to your `.gitignore`.

# OPTIONS

The following options are supported by **newt**.

-h, -help
: display this help message

-license
: display the software license

-version
: display version information

-verbose
: If set provide verbose debugging output for requests

# ACTION

init [YAML_FILE]
: this will create your initial Newt YAML file based on a set of interactive
questions. It can also update your `.gitignore` file too.

model [YAML_FILE]
: this will read the current Newt YAML file and run the interactive modeler updating the Newt YAML file if model(s) are accepted.

check [YAML_FILE]
: analyze the Newt YAML file and report problems if found.

generate [YAML_FILE]
: this is used to generate your SQL, PostgREST configuration and Mustache templates
based on the contents of your Newt YAML file.

run [YAML_FILE]
: this will run the defined services in the application attribute of the Newt YAML file.
This is intended for use in development. In a production setting you'd setup the individual
services to run from systemd or init as services.

ws [YAML_FILE]
: run Newt's static web server.

# YAML_FILE 

**newt** is configured in a YAML file. What is described below is a summary of 
YAML syntax use in a Newt project that uses all of the Newt programs.

## Top level properties

These are the top level properties in YAML files.

applications
: (optional) holds the run time configuration used by the Newt applications.

models
: (required by newt generator) This holds the description of the data models in your application.  Each model uses HTML 5 element descriptions which can be set using the interactive `newt model` command. 

routes
: (required by newt router) This holds the routes for the data pipeline (e.g. JSON API and template engine sequence)

templates
: (required by newt mustache)

## The applications property

The applications properties are optional. Some maybe set via command line. See Newt application's manual page for specific ones. These properties lets you override the default settings of Newt programs.

newtrouter
: this contains configuration for the Newt Router, i.e. port and htdocs

newtgenerator
: this contains configuration for code generation. e.g. namespace used by Postgres and PostgREST.

newtmustache
: this contains configuration for Newt Mustache, i.e. port

options
: holds key value pairs of which can be referenced in the values of models, routes and templates.

environment
: (optional: newtrouter, newtmustache) this is a list of operating system environment that will be available to routes. This is used to pass in secrets (e.g. credentials) need in the pipeline

### Each application may have one or more of these properies.

namespace
: (newtgenerator) uses this in the SQL generated for setting up Postgres+PostgREST

port
: (all) Port number to used for Newt web service running on localhost

htdocs
: (newtrouter) Directory that holds your application's static content

timeout
: (newt router, newt mustache) The time in seconds to timeout a HTTP transaction

dsn
: (postgres) The data source name (URI for database connection string). Using to connect to Postgres

app_path
: (postgrest) The full path to the application (if it is not in the PATH environment variable)

conf_app
: (postgres) Configuration path for the application.

## the "routes" property

Routes hosts a list of request descriptions and their data pipelines. This property is only used by Newt router and Newt code generator.

### a route object

`id`
: (required) This identifies the pipeline. It maybe used in code generation. It must conform to variable name rules[^21]

`description`
: (optional, encouraged) This is a human readable description of what you're trying to accomplish in this specific route. It may be used in comments or by documentation generators.

`request [METHOD ][PATH]`
: (required) This is a string that expresses the HTTP method and URL path to used to trigger running the data pipeline. If METHOD is not provided it will match using just the path. This is probably NOT what you want. You can express embedded variables in the PATH element. This is done by using single curl braces around a variable name. E.g. `GET /items/{item_id}` would make `item_id` available in building your service paths in the pipeline. The pattern takes up a whole path segment so `/blog/{year}-{month}-{day}` would not work but `/blog/{year}/{month}/{day}` would capture the individual elements. The Newt router sits closely on top of the Go 1.22 HTTP package route handling. For the details on how Go 1.22 and above request handlers and patterns form see See <https://tip.golang.org/doc/go1.22#enhanced_routing_patterns> and <https://pkg.go.dev/net/http#hdr-Patterns> for explanations.

`pipeline`
: (required) this is a list of URLs to one or more web services visible on localhost. The first stage to fail will abort the pipeline returning an HTTP error status. If done fail then the result of the last stage it returned to the requesting browser.

`debug`
: (optional) if set to true the `newt` service will log verbose results to standard out for this specific pipeline

#### a pipeline object

A pipeline is a list of web services containing a type, URL, method and content types

`service [METHOD ][URL]`
: (required) The HTTP method is included in the URL The URL to be used to contact the web service, may contain embedded variable references drawn from the request path as well as those passed in through `.application.environment`.  All the elements extracted from the elements derived from the request path are passed through strings. These are then used to construct a simple key-value object of variable names and objects which are then passed through the Mustache template representing the target service URL. 

`description`
: (optional, encouraged) This is a description of what this stage of the pipe does. It is used when debug is true in the log output and in program documentation.

`timeout`
: (optional) Set the timeout in seconds for receiving a response from the web server. Remember the time spent at each stage is the cumulative time your browser is waiting for a response. For this reason you may want to set the timeout to a small number.


## the "models" property

Models holds a list of individual models used by our data pipelines. The models are by Newt code generator and the Newt router. Models defines a superset of the GitHub YAML issue template syntax (abbr: GHYITS).

### a model object

The model object is based largely on GitHub YAML issue template syntax with a couple extra properties that are Newt enhancements.

id
: (required, newt specific) this is the name identifying the model. It must conform to variable name rules[^21]

The following properties are based on the GitHub YAML issue template syntax[^22] (abbr: GHYITS)

name
: (required: GHYITS, optional: newt) Must be unique to use with GitHub YAML issue templates[^22]. In Newt it will be used in populating comments in generated SQL

description
: (required: GHYITS, optional: newt) A human description of the model, It will appear in the web form and SQL components generated from the model

body
: (required) A a list of input types. Each input type maps to columns in SQL, input element in web forms and or HTML elements in read only pages

#### a model's input types

This is based on GitHub YAML issue template (abbr: GHYITS) input types[^23]. 

id
: (required) an identifier for the element. Must conform to variable name rules[^21]. It is used to SQL as a column name and in web forms for the input property.

type
: (required) Identifies the type of elements (input, textarea, markdown, checkbox, dropdown).

attributes
: (optional) A key-value list that define properties of the element. These used in rendering the element in SQL or HTML.

validations
: (optional, encouraged) A set of key-value pairs setting constraints of the element content. E.g. required, regexp properties, validation rule provided with certain identifiers (e.g. DOI, ROR, ORCID).


## input types

Both the routes and models may contain input types. The types supported in Newt are based on the types found in the GHYITS for scheme[^23]. They include

markdown
: (models only) markdown request displayed to the user but not submitted to the user but not submitted by forms. 

textarea
: (models only) A multi-line text field

input
: A single line text field. This conforms to value input types in HTML 5 and can be expressed using the CSS selector notation. E.g. `input[type=data]` would be a date type. This would result in a date column type in SQL, a date input type in HTML forms and in formatting other HTML elements for display.

dropdown
: A dropdown menu. In SQL this could render as an enumerated type. In HTML it would render as a drop down list

checkboxes
: A checkbox element. In SQL if the checkbox is exclusive (e.g. a radio button) then the result is stored in a single column, if multiple checks are allowed it is stored as a JSON Array column.

Newt may add additional types in the future.

## Example Newt YAML file, "app.yaml" 

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
  - request: /app_create
    template: app_create_form.tmpl
    description: Display a app for create
  - request: /app_create_response
    template: app_create_response.tmpl
    description: This is an result template for app create
  - request: /app_update
    template: app_update_form.tmpl
    description: Display a app for update
  - request: /app_update_response
    template: app_update_response.tmpl
    description: This is an result template for app update
  - request: /app_delete
    template: app_delete_form.tmpl
    description: Display a app for delete
  - request: /app_delete_response
    template: app_delete_response.tmpl
    description: This is an result template for app delete
  - request: /app_read
    template: app_read.tmpl
    description: This template handles app read
  - request: /app_list
    template: app_list.tmpl
    description: This template handles app list
~~~


[^21]: variable numbers must start with a letter, may contain numbers but not spaces or punctuation except the underscore

[^22]: See <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-issue-forms>, 

[^23]: See <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema>

## templates property

This property is used by Newt Mustache. It is ignore by Newt router and code generator.

templates
: (optional: newtmustache) this holds a list of template objects

### template object model

The template objects are used by Newt Mustache template engine. If you're not using it you can skip these.

`request [METHOD ][PATH]`
: (required) This holds the request HTTP method and path. If the HTTP method is missing a POST is assumed

`template`
: (required: newtmustache only) This is the path to the template associated with request. NOTE: Pandoc web service does not support partial templates. Mustache does support partial templates

`partials`
: (optional, newtmustache only) A list of paths to partial Mustache templates used by `.template`.

`options`
: (optional, newtmustache only) An object that can be merged in with JSON received for processing by your Mustache template

`debug`
: (optional) this turns on debugging output for this template

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


