---
<<<<<<< HEAD
title: newt(1) user manual | 0.0.8 7af9eec
=======
title: newt(1) user manual | 0.0.8 26115a9
>>>>>>> 6924d3476fef94ed779922846bce6bf411b8a9fa
pubDate: 2024-04-18
author: R. S. Doiel
---

# NAME

newt

# SYNOPSIS

newt [OPTIONS] ACTION YAML_CONFIG_FILE 

# DESCRIPTION

**newt** sets up up and can run your Newt application during development.
**newt** supports the "init", "generate" and "run" actions. The "init"
command is used when you are starting a Newt Project. It will guide you through
creating your initial Newt YAML file and also optionally make additions to your
`.gitignore`. "generate" will run the Newt Generator to create the
SQL, PostgREST configuration and Mustache templates based on the contents of the
Newt YAML file. The third action is "run". It is used to run your Newt based
application. For the applications you have defined in your Newt YAML_CONFIG_FILE
it's start them up. These include the Newt Router, Newt Mustache and PostgREST.
This allows you to quick run and stop the services are you craft your
Newt YAML file for your project.


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

init
: this will create your initial Newt YAML file based on a set of interactive
questions. It can also update your `.gitignore` file too.

generate
: this is used to generate your SQL, PostgREST configuration and Mustache templates
based on the contents of your Newt YAML file.

run
: this will run the defined services in the application attribute of the Newt YAML file.
This is intended for use in development. In a production setting you'd setup the individual
services to run from systemd or init as services.

sws
: this will run Newt's static web server.

# YAML_CONFIG_FILE 

**newt** is configured in a YAML file. What is described below is the complete
YAML syntax use in a Newt project that uses all of the Newt programs.

## Top level properties

These are the top level properties in YAML files.

applications
: (optional) holds the run time configuration used by the Newt applications.

models
: (required by newtgenerator) This holds the description of the data models in your application. Each model uses the [GitHub YAML issue template syntax](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/configuring-issue-templates-for-your-repository#creating-issue-forms) (abbr: GHYITS)

routes
: (required by newtrouter) This holds the routes for the data pipeline (e.g. JSON API and template engine sequence)

templates
: (required by newtmustache)

## The applications property

The applications properties are optional. Some maybe set via command line. See Newt application's manual page for specific ones. These properties lets you override the default settings of Newt programs.

newtrouter
: this contains configuration for the Newt Router, i.e. port and htdocs

newtgenerator
: this contains configuration for the Newt Generator, e.g. port, namespace

newtmustache
: this contains configuration for Newt Mustache, i.e. port

options
: holds key value pairs of which can be referenced in the values of models, routes and templates.

environment
: (optional: newtrouter, newtmustache) this is a list of operating system environment that will be available to routes. This is used to pass in secrets (e.g. credentials) need in the pipeline

### Configuring Newt programs

namespace
: (newtgenerator) uses this in the SQL generated for setting up Postgres+PostgREST

port
: (all) Port number to used for Newt web service running on localhost

htdocs
: (newtrouter) Directory that holds your application's static content

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

## Example Newt YAML for router and code generator

~~~yaml
applications:
  newtgenerator:
    namespace: people # E.g. "people" Namespace to use generating Postgres SQL
  newtrouter:
    port: 8011 # Port number for Newt Router
    htdocs: htdocs # Path to static content directory if required
  newtmustache:
    port: 8012 # Port number for Newt Mustache 
  #options:
    # name value pairs used for aliasing strings in routes, models, and templates 
  environment:
    - DB_USER
    - DB_PASSWORD
    - DB_HOST
models:
  - id: people
    name: People Profiles
    description: |
      This models a curated set of profiles of colleagues
    body:
      - id: people_id
        type: input
        attributes:
          label: A unique person id, no spaces, alpha numeric
          placeholder: ex. jane-do-007
        validations:
          required: true
      - id: display_name
        type: input
        attributes:
          label: (optional) A person display name
          placeholder: ex. J. Doe, journalist
      - id: family_name
        type: input
        attributes:
          label: (required) A person's family name or singular when only one name exists
          placeholder: ex. Doe
        validations:
          required: true
      - id: given_name
        type: input
        attributes:
          label: (optional, encouraged) A person's given name
          placeholder: ex. Jane
      - id: orcid
        type: input
        attributes:
          label: (optional) A person's ORCID identifier
          placeholder: ex. 0000-0000-0000-0000
        validations:
          pattern: "[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]"
      - id: ror
        type: input
        attributes:
          label: (optional) A person's ROR identifying their affiliation
      - id: email
        type: "input[type=email]"
        attributes:
          label: (optional) A person public email address
      - id: website
        type: "input[type=url]"
        attributes:
          label: (optional) A person's public website
          placeholder: ex. https://jane.doe.example.org
routes:
  - id: create_person
    description: Create new person profile
    request: POST /person
    pipeline:
      - description: This will generate a new user in the database
        service: POST "https://{{DB_USER}}@{{DB_HOST}}:3000/rpc/people"
        content_type: application/json
      - description: |
          This sends the results of creating a person to the template engine
        service: POST http://localhost:3032/people_update_result.tmpl
        content_type: application/json
  - id: read_person
    description: Update a person's profile
    request: "GET /person/{{people.people_id}}"
    pipeline:
      - description: Retrieve a person's profile
        service: POST "https://{{DB_USER}}@{{DB_HOST}}:3000/person/{{people.people_id}}"
        content_type: application/json
      - description: |
          Render a person's profile
        service: POST http://localhost:3032/profile.tmpl
        content_type: application/json
  - id: update_person
    description: Update person's profile
    request: "PUT /person/{{people.people_id}}"
    pipeline:
      - description: This will update a person record in the database
        service: PUT "https://{{DB_USER}}@{{DB_HOST}}:3000/rpc/people"
        content_type: application/json
      - description: |
          This sends the results of updating a person to the template engine
        service: POST http://localhost:3032/people_update_result.tmpl
        content_type: application/json
  - id: delete_person
    description: Remove person's profile
    request: "DELETE /person/{{people.people_id}}"
    pipeline:
      - description: Remove the person for the database
        service: DELETE "https://{{DB_USER}}@{{DB_HOST}}:3000/people/{{people.people_id}}"
        content_type: application/json
      - description: Displace the result of what happened in the removal
        service: POST http://localhost:3032/removed_person.tmpl
  - id: list_people
    description: List people profiles available
    request: GET /people
    pipeline:
      - description: Retrieve a list of all people profiles available
        service: GET https://{{DB_HOST}}@{{DB_HOST}}:3000/people
        content_type: application/json
      - description: format a browsable people list linking to individual profiles
        service: POST http://localhost:3030/list_people.tmpl
        content_type: application/json
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


