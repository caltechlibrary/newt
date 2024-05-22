
# Newt YAML syntax

The Newt programs are configured in a YAML file. Each Newt program may focus on some properties and ignore others. What is described below is the complete YAML syntax use in a Newt project that uses all of the Newt programs.

## Top level properties

These are the top level properties in YAML files.

`applications`
: (optional) holds the run time configuration used by the Newt applications.

`models`
: (required by newtgenerator) This holds the description of the data models in your application. Each model describes a series of HTML 5 elements that make up the model. 

`routes`
: (required by newtrouter) This holds the routes for the data pipeline (e.g. JSON API and template engine sequence)

`templates`
: (required by newtmustache)

## The applications property

The applications properties are optional. Some maybe set via command line. See Newt application's manual page for specific ones. These properties lets you override the default settings of Newt programs.

`router`
: this contains configuration for the Newt Router, i.e. port and htdocs

`mustache`
: this contains configuration for Newt Mustache, i.e. port

`postgres`
: this contains configuration information for the running Postgres database

`postgrest`
: this contains configuration information managing PostgREST application with `newt`

`options`
: holds key value pairs of which can be referenced in the values of models, routes and templates.

`environment`
: (optional) this is a list of operating system environment variables that will be available to models, routes and templates. You can use this to pass in secrets (e.g. credentials) to your pipelined services. The pairing of environment variable and value are merged into the options when Newt services start up.

### Configuring Newt programs

`app_path`
: (postgrest) the path to the program(s) managed by `newt`

`conf_path`
: (postgrest) path to PostgREST configuration file

`namespace`
: (postgres) uses this in the SQL generated for setting up Postgres+PostgREST

`port`
: (all) Port number to used for Newt web service running on localhost

`htdocs`
: (newt, newtrouter) Directory that holds your application's static content


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

Models holds a list of individual models used by our data pipelines. The models are used by Newt code generator and the Newt router. Models were inspired by a subset of GitHub YAML issue template syntax (abbr: GHYITS). Models contain lists of elements, elements are descriptions of HTML5 input elements and their attributes.

### a model object

The model object is based largely on GitHub YAML issue template syntax with a couple extra properties that are Newt enhancements.

`id`
: (required) this is the name identifying the model. It must conform to variable name rules[^21]

The following properties are inspired by GitHub YAML issue template syntax[^22] (abbr: GHYITS)

`description`
: (required: GHYITS, optional: newt) A human description of the model, It will appear in the web form and SQL components generated from the model

`elements`
: (required) A a list of input types. Each input type maps to columns in SQL, input element in web forms and or HTML elements in read only pages

#### a model's input types

This is based on HTML5 form input types. They were inspired by GitHub YAML issue template (abbr: GHYITS) input types[^23]. 

`id`
: (required) an identifier for the element. Must conform to variable name rules[^21]. It is used to SQL as a column name and in web forms for the input property.

`type`
: (required) Identifies the type of elements (input, textarea, markdown, checkbox, dropdown).

`attributes`
: (optional) A key-value list that define properties of the element. These used in rendering the element in SQL or HTML.

`primary_key`
: (boolean, defaults to false) if true this element is treat as the object identifier or primary key. It used in routing and in generate SQL schema and queries.

## input types

The model types supported in Newt are based on the HTML5 form input types. They are inspired by GHYITS scheme[^23]. Example
types include

`text`
: A single line text field. This conforms to value input types in HTML 5 and can be expressed using the CSS selector notation as `input[type=date]`.

`textarea`
: (models only) A multi-line text field

`dropdown`
: A dropdown menu. In SQL this could render as an enumerated type. In HTML it would render as a drop down list

`checkboxes`
: A checkbox element. In SQL if the checkbox is exclusive (e.g. a radio button) then the result is stored in a single column, if multiple checks are allowed it is stored as a JSON Array column.

`email`
: This conforms to a text string with validation for email format

See <https://developer.mozilla.org/en-US/docs/Learn/Forms/HTML5_input_types> for more possibilities.

IMPORTANT: Newt does NOT support the file input type because Newt applications do not support file uploads. 

Newt extends support for additional identifiers used by galleries, libraries, archives and museum by generating the necessary
validation at the SQL level and use regular expressions in the pattern attribute of HTML5 input elements.  When expressed as
HTML these can be identified by the additional attribute "extended-type".

orcid
: This identifies the element as an ORCID. Server side validation will be generated at the SQL level and client side validation via HTML5 input element with pattern attribute.

~~~html
<input type="text" extended-type="orcid" pattern="[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9A-Z]">
~~~


## Example Newt YAML for router and code generator

~~~yaml
applications:
  router:
    port: 8010
    htdocs: htdocs
  template_engine:
    port: 8011
  postgres:
    namespace: people
    port: 5432
    dsn: postgres://{PGUSER}:{PGPASSWORD}@localhost:5432/people.yaml
  postgrest:
    app_path: postgrest
    conf_path: postgrest.conf
    port: 3000
  enviroment:
    - PGUSER
    - PGPASSWORD
models:
  - id: people
    description: |
      This models a curated set of profiles of colleagues
routes:
  - id: people_create
    request: GET /people_create
    description: Handle retrieving the webform for people create
    pipeline:
      - service: POST http://localhost:8011/people_create
        description: Display a people for create
  - id: people_create
    request: POST /people_create
    description: Handle form submission for people create
    pipeline:
      - service: POST http://localhost:3000/rpc/people_create
        description: Access PostgREST API for people create
      - service: POST http://localhost:8011/people_create_response
        description: This is an result template for people create
  - id: people_update
    request: GET /people_update/{oid}
    description: Handle retrieving the webform for people update
    pipeline:
      - service: GET http://localhost:3000/rpc/people_read/{oid}
        description: Retrieve people from PostgREST API before update
      - service: POST http://localhost:8011/people_update
        description: Display a people for update
  - id: people_update
    request: POST /people_update
    description: Handle form submission for people update
    pipeline:
      - service: PUT http://localhost:3000/rpc/people_update/{oid}
        description: Access PostgREST API for people update
      - service: POST http://localhost:8011/people_update_response
        description: This is an result template for people update
  - id: people_delete
    request: GET /people_delete/{oid}
    description: Handle retrieving the webform for people delete
    pipeline:
      - service: GET http://localhost:3000/rpc/people_read/{oid}
        description: Retrieve people from PostgREST API before delete
      - service: POST http://localhost:8011/people_delete
        description: Display a people for delete
  - id: people_delete
    request: POST /people_delete
    description: Handle form submission for people delete
    pipeline:
      - service: DELETE http://localhost:3000/rpc/people_delete/{oid}
        description: Access PostgREST API for people delete
      - service: POST http://localhost:8011/people_delete_response
        description: This is an result template for people delete
  - id: people_read
    request: POST /people_read
    description: Retrieve object(s) for people read
    pipeline:
      - service: GET http://localhost:3000/rpc/people_read/{oid}
        description: Access PostgREST API for people read
      - service: POST http://localhost:8011/people_read
        description: This template handles people read
  - id: people_list
    request: POST /people_list
    description: Retrieve object(s) for people list
    pipeline:
      - service: GET http://localhost:3000/rpc/people_list
        description: Access PostgREST API for people list
      - service: POST http://localhost:8011/people_list
        description: This template handles people list
templates:
  - id: people_create
    request: /people_create
    template: people_create_form.tmpl
    description: Display a people for create
  - id: people_create
    request: /people_create_response
    template: people_create_response.tmpl
    description: This is an result template for people create
  - id: people_update
    request: /people_update
    template: people_update_form.tmpl
    description: Display a people for update
  - id: people_update
    request: /people_update_response
    template: people_update_response.tmpl
    description: This is an result template for people update
  - id: people_delete
    request: /people_delete
    template: people_delete_form.tmpl
    description: Display a people for delete
  - id: people_delete
    request: /people_delete_response
    template: people_delete_response.tmpl
    description: This is an result template for people delete
  - id: people_read
    request: /people_read
    template: people_read.tmpl
    description: This template handles people read
  - id: people_list
    request: /people_list
    template: people_list.tmpl
    description: This template handles people list
~~~


[^21]: variable numbers must start with a letter, may contain numbers but not spaces or punctuation except the underscore

[^22]: See <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-issue-forms>, 

[^23]: See <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema>

## templates property

This property is used by Newt Mustache. It allows Newt Mustache to map a POST to a specific template or
template and set of partials.

templates
: (optional: newtmustache) this holds a list of template objects.

### template object model

The template objects are used by Newt Mustache template engine. If you're not using it you can skip these.
Newt Mustache only responds to POST methods which have a defined path in the templates section of the 
Newt YAML template file.

`id`
: (required) the identifier for related templates. E.g. a web form and the response form are associated with a common identifier. Newt's modeler uses this field to manage template updates and removals.

`request [PATH]`
: (required) This holds the request HTTP method and path. If the HTTP method is missing a POST is assumed. Requests in templates like with request in routes can contain inline string variables. Example would be `{name}` used in the YAML below.

`template`
: (required) This is the path to the template associated with request. NOTE: Pandoc web service does not support partial templates. Mustache does support partial templates

`partials`
: (optional) A list of paths to partial Mustache templates used by `.template`.

`vocabularies`
: (optional) This is a YAML file containing "vocabulary" items that can be incorporated into the output of your Mustache template. Example would be auto complete form elements. It is read at Newt Mustache startup along with the template.

`options`
: (optional) An object that can be merged in with JSON received for processing by your Mustache template. Options can be set per template but are also inherited from the resolved names and values in the applications' attribute of the Newt YAML file.

`debug`
: (optional) this turns on debugging output for this template.

Example of newtmustache YAML:

~~~yaml
applications:
  template_engine:
    port: 8012
templates:
  - id: hello
    request: /hello/{name}
    template: testdata/simple.tmpl
  - id: hello
    request: /hello
    template: testdata/simple.tmpl
    options:
      name: Universe
  - id: hi
    request: /hi/{name}
    template: testdata/hithere.tmpl
    partials:
      - testdata/name.tmpl
    debug: true
  - id: hi
    request: /hi
    template: testdata/hithere.tmpl
    partials:
      - testdata/name.tmpl
    options:
      name: Universe
~~~

