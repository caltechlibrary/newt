
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

`newtrouter`
: this contains configuration for the Newt Router, i.e. port and htdocs

`newtgenerator`
: this contains configuration for the Newt Generator, i.e. namespace

`newtmustache`
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
: (newt) the path to the program(s) managed by `newt`

`namespace`
: (newtgenerator) uses this in the SQL generated for setting up Postgres+PostgREST

`port`
: (all) Port number to used for Newt web service running on localhost

`htdocs`
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
  newtgenerator:
    namespace: people # E.g. "people" Namespace to use generating Postgres SQL
  newtrouter:
    port: 8011 # Port number for Newt Router
    htdocs: htdocs # Path to static content directory if required
  newtmustache:
    port: 8012 # Port number for Newt Mustache 
  postgres:
    port: 5432 # Where to find the Postgres port
  postgrest: # Configuration for running PostgREST via newt
    port: 3000
    app_path: /usr/local/bin/postgrest
    conf: people.conf
  environment:
    - PGUSER
    - PGPASSWORD
    - PGHOST
models:
  - id: people
    name: People Profiles
    description: |
      This models a curated set of profiles of colleagues
    elements:
      - id: people_id
        type: input
        attributes:
          label: A unique person id, no spaces, alpha numeric
          placeholder: ex. jane-do-007
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
          required: true
      - id: given_name
        type: input
        attributes:
          label: (optional, encouraged) A person's given name
          placeholder: ex. Jane
      - id: orcid
        type: orcid
        attributes:
          label: (optional) A person's ORCID identifier
          placeholder: ex. 0000-0000-0000-0000
      - id: email
        type: email
        attributes:
          label: (optional) A person public email address
      - id: website
        type: url
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
    desciption: Update a person's profile
    request: "GET /person/{{people.people_id}}"
    pipeline:
      - description: Retrieve a person's profile
        service: POST "https://{{DB_USER}}@{{DB_HOST}}:3000/person/{{people.people_id}}"
        content_type: application/json
      - description: |
          Render a person's profile
        service: POST http:localhost:3032/profile.tmpl
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
      - description: Displace the result of what happed in the removal
        service: POST http://localhost:3032/removed_person.tmpl
  - id: list_people
    description: List people profiles available
    request: GET /people
    pipeline:
      - description: Retrieve a list of all people profiles available
        service: GET https://{{DB_HOST}}@{{DB_HOST}}:3000/people
        content_type: application/json
      - discription: format a browsable people list linking to individual profiles
        service: POST http://localhost:3030/list_people.tmpl
        content_type: applicatin/json
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
  newtmustache:
    port: 8012
templates:
  - id: hello
    request: /hello/{name}
    template: testdata/simple.mustache
  - id: hello
    request: /hello
    template: testdata/simple.mustache
    options:
      name: Universe
  - id: hi
    request: /hi/{name}
    template: testdata/hithere.html
    partials:
      - testdata/name.mustache
    debug: true
  - id: hi
    request: /hi
    template: testdata/hithere.html
    partials:
      - testdata/name.mustache
    options:
      name: Universe
~~~

