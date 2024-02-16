
# Newt YAML syntax

Newt programs are configured in YAML files. Newt programs may focus on some properties and ignore others. The interpretation is specific to the program.

These are the top level properties in YAML files.

application
: (optional: newtrouter, newtgenerator and mustache) holds the run time configuration used by the Newt web service and metadata about the application you're creating.

models
: (optional: newtrouter and newtgenerator) This holds the description of the data models in your application. Each model uses the [GitHub YAML issue template syntax](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/configuring-issue-templates-for-your-repository#creating-issue-forms) (abbr: GHYTS)

routes
: (optional: newtrouter and newtgenerator) This holds the routes for the data pipeline (e.g. JSON API and template engine sequence)

templates
: (optional: newtmustache, pdbundler)

## the "application" property

The application properties are optional. Some maybe set via command line. The the Newt application manual pages.

namespace
: (optional: newtgenerator) uses this in the SQL generated for setting up Postgres+PostgREST

port
: (optional: newtrouter, newtmustache, pdbundler) default is This port number the Newt web services uses to listen for request on localhost

htdocs
: (optional: newtrouter only) Directory that holds your application's static content

metadata
: (optional) This holds metadata about your application using the [CITATION.cff](https://citation-file-format.github.io/) YAML syntax under metadata.

environment
: (optional: newtrouter, newtmustache) this is a list of operating system environment that will be available to routes. This is used to pass in secrets (e.g. credentials) need in the pipeline

The fifth attribute in application is special. It can be used in instead of `metadata`. If you maintain a CITATION.cff file you can point to it to avoid maintaining it in two places. When the Newt router or code generated is started up it will copy the contents into the `metadata` property.

citation
: (optional) This points at an file (e.g. CITATION.cff). It is used to populate the metadata property at startup

## the "routes" property

Routes hosts a list of request descriptions and their data pipelines. This property is only used by Newt router and Newt code generator.

### a route object

`id`
: (required) This identifies the pipeline. It maybe used in code generation. It must conform to variable name rules[^1]

`description`
: (optional, encouraged) This is a human readable description of what you're trying to accomplish in this specific route. It may be used in comments or by documentation generators.

`request [METHOD ][PATH]`
: (required) This is a string that expresses the HTTP method and URL path to used to trigger running the data pipeline. If METHOD is not provided it will match using just the path. This is probably NOT what you want. You can express embedded variables in the PATH element. This is done by using single curl braces around a variable name. E.g. `GET /items/{item_id}` would make `item_id` available in building your service paths in the pipeline. The pattern takes up a whole path segment so `/blog/{year}-{month}-{day}` would not work but `/blog/{year}/{month}/{day}` would capture the individual elements. The Newt router sits closely on top of the Go 1.22 HTTP package route handling. For the details on how Go 1.22 and above request handlers and patterns form see See <https://tip.golang.org/doc/go1.22#enhanced_routing_patterns> and <https://pkg.go.dev/net/http#hdr-Patterns> for explanations.

`pipeline`
: (required) this is a list of URLs to one or more web services visible on localhost. The first stage to fail will abort the pipeline returning an HTTP error status. If done fail then the result of the last stage it returned to the requesting browser.

`error`
: (optional) this points to a static page that can be displayed when the pipeline fails (e.g. like a 404 page used by web servers)

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

Models holds a list of individual models used by our data pipelines. The models are by Newt code generator and the Newt router. Models defines a superset of the GitHub YAML issue template syntax (abbr: GHYTS).

### a model object

The model object is based largely on GitHub YAML issue template syntax with a couple extra properties at are Newt specific.

id
: (required, newt specific) this is the name identifying the model. It must conform to variable name rules[^1]

routing
: (optional, newt specific) this holds a list of route ids associated with this model. It is use in code generation, e.g. to populate a web form's action and model

The following properties are based on the GitHub YAML issue template syntax[^2] (abbr: GHYTS)

name
: (required by GHYTS, optional in newt) Must be unique to use with GitHub YAML issue templates[^2]. In Newt it will be used in populating comments in generated SQL

description
: (required) A human description of the model, It will appear in the web form and SQL components generated from the model

body
: (required) A a list of input types. Each input type maps to columns in SQL, input element in web forms and or HTML elements in read only pages

#### a model's input types

This is based on GitHub YAML issue template (abbr: GHYTS) input types[^3]. 

id
: (required) an identifier for the element. Must conform to variable name rules[^1]. It is used to SQL as a column name and in web forms for the input property.

type
: (required) Identifies the type of elements (input, textarea, markdown, checkbox, dropdown).

attributes
: (optional) A key-value list that define properties of the element. These used in rendering the element in SQL or HTML.

validations
: (optional, encouraged) A set of key-value pairs setting constraints of the element content. E.g. required, regexp properties, validation rule provided with certain identifiers (e.g. DOI, ROR, ORCID).


## input types

Both the routes and models may contain input types. The types supported in Newt are based on the types found in the GHYTS for scheme[^3]. They include

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


## Example Newt YAML for router and code generator

```yaml
application:
  port: 8011
  htdocs: htdocs
  metadata:
    cff-version: 1.2.0
    message: Demo of Newt YAML file
    type: software
    title: Newt a faster way to build metadata curation applications
    abstract: |
      This is a demonstation of a YAML that can generate a simple
      application to manage people and groups
    version: 0.0.0
    status: concept
    authors:
      - family-names: Doiel
        given-names: R. S.
        orcid: "https://orcid.org/0000-0003-0900-6903"
    keywords:
      - demo
      - newt
      - rapid application development
  environment:
    - DB_USER
    - DB_PASSWORD
    - DB_HOST
models:
  - id: people
    name: People Profiles
    description: |
      This models a curated set of profiles of colleagues
    routing:
      - create_person
      - read_person
      - update_person
      - delete_person
      - list_people
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
          label: (optional) A person's ROR identifing their affiliation
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
```


[^1]: variable numbers must start with a letter, may contain numbers but not spaces or punctuation except the underscore

[^2]: See <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-issue-forms>, 

[^3]: See <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema>

## templates property

This property is used by Newt Mustache and pdbundler. It is ignore by Newt router and code generator.

templates
: (optional: newtmustache, pdbundler) this holds a list of template objects

### template object model

`request [METHOD ][PATH]`
: (required) This holds the request HTTP method and path. If the HTTP method is missing a POST is assumed

`template`
: (required: newtmustache, optional: pdbundler) This is the path to the template associated with request. NOTE: Pandoc web service does not support partial templates. Mustache does support partial templates

`partials`
: (optional, newtmustache only) A list of paths to partial Mustache templates used by `.template`.

`options`
: (optional: pdbundler) This are the default options you want to supply Pandoc web service, e.g. "to", "from", "standalone", "title"

`debug`
: (optional) this turns on debugging output for this template

Example of newtmustache YAML:

```yaml
application:
    port: 3032
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
```

Example of pdbundler YAML:

```yaml
application:
    port: 3029
templates:
  - request: POST /hello
    template: hello.md.tmpl
    options:
      from: markdown
      to: markdown
      standalone: false
    debug: true
  - request: POST /hello/{name}
    template: hello.md.tmpl
    options:
      from: markdown
      to: markdown
    debug: true
  - request: "POST /custom_page"
    template: testdata/pdoc.html.tmpl
    options:
      from: markdown
      to: html5
      standalone: true
      title: This is the custom template with this title
    debug: true
  - request: "POST /custom_page_with_title/{title}"
    template: testdata/pdoc.html.tmpl
    options:
      from: markdown
      to: html5
      standalone: true
      title: This title is overwritten by the one in the request
    debug: true
  - request: "POST /custom_page_include"
    template: testdata/pdoc.html.tmpl
    options:
      from: markdown
      to: html5
      standalone: false
    debug: true
  - request: "POST /default_html5"
    options:
      from: markdown
      to: html5
      standalone: true
      title: A Page using the default template
    debug: true
  - request: "POST /default_html5_with_title/{title}"
    options:
      from: markdown
      to: html5
      standalone: true
      title: This title is replaced by the title in the URL
  - request: "POST /default_html5_include"
    options:
      from: markdown
      to: html5
      standalone: false
    debug: true
```


