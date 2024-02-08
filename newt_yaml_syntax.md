
# Newt YAML syntax

Newt's configuration and modeling is based on a YAML file. That YAML file has a
specific syntax.  The top level of that syntax are three properties.

application
: holds the run time configuration used by the Newt web service and metadata about the application you're creating.

models
: This holds the description of the data models in your application. Each model uses the [GitHub YAML issue template syntax](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/configuring-issue-templates-for-your-repository#creating-issue-forms)

routes
: This holds the routes for the data pipe line (e.g. JSON API and template engine sequence)


## the "application" property

The application property itself has four properties.

port
: (optional, default is This port number the Newt web services uses to listen for request on localhost

htdocs
: (optional) Directory that holds your application's static content

metadata
: (optional) This holds metadata about your application using the [CITATION.cff](https://citation-file-format.github.io/) YAML syntax under metadata.

environment
: (optional) this is a list of operating system environment that will be available to routes. This is used to pass in secrets (e.g. credentials) need in the pipeline

## the "routes" property

Routes hosts a list of request descriptions and their data pipelines

### a route object

name
: (required) This identifies the pipe line so that it can be re-used or included in other pipelines. It must conform to variable name rules[^1]

request
: (required) This is a string that express the paith to listen to for running the data pipeline

pipeline
: (required) this is a list of data request used to form the pipeline needed for this request. the last element in the pipeline is returned as a response to the request

description
: (optional) a description of the pipeline and what it is trying to accomplihs

error
: (optional) this points to a static page that can be displayed when the pipeline fails (e.g. like a 404 page used by web servers)

variables
: (optional) this is a key-value list of input types that may be used in route mapping. These are base on input types in models (see below).

debug
: (optonal) if set to true the pipeline process will be logged to standard out


## the "models" property

Models holds a list of individual models used by our data pipelines. The models are used when generating the SQL code defining our models and data manage. It also is used to render templates that maybe used by a template engine in the data pipeline. It can be used to generate static HTML suitable for embedding in Markdown or HTML documents (e.g. web forms)

### a model object

id
: (required, newt specific) this is the name identifying the model. It must conform to variable name rules[^1]

The following following properties are based on the GitHub YAML issue template syntax[^2] (abbr: GHYT)

name
: (required by GHYT) Must be unique to use with GitHub YAML issue templates[^2]

description
: (required) A human description of the model, It will appear in the web form and SQL commonents  generated from the model

body
: (required) A a list of input types. Each input type maps to columns in SQL, input element in web forms and or HTML elements in read only pages

#### a model's input types

This is based on GitHub YAML issue template (abbr: GHYT) input types[^3]. 

id
: (required) an identifier for the element. Must conform to variable name rules[^1]. It is used to SQL as a column name and in web forms for the input property.

type
: (required) Identifies the type of elements (input, textarea, markdown, checkbox, dropdown).

attributes
: (optional) A key-value list that define properties of the element. These used in rendering the element in SQL or HTML

validations
: (optional, encouraged) A set of key-value pairs setting constraints of the element content. E.g. required, regexp properties, validation rule provided with certain identifiers (e.g. DOI, ROR, ORCID).


## input types

Both the routes and models may continue input types. The types supported in Newt are based on the types found in GitHub YAML issue types[^3]. They include

markdown
: (models only) markdown request displayed to the user but not submitted to the user but not submitted by forms. 

textarea
: (models only) A multi-line text field

input
: A single line text field. This conforms to value input types in HTML 5 and can be expressed using the CSS selector notation. E.g. `input[type=data]` would be a date type. This would result in a date column type in SQL, a date input type in HTML forms and in formatting other HTML elements for display.

dropdown
: A dropdown menu. In SQL this could render as an enumerated type. In HTML it would render as a drop down list

checkboxes
: A checkbox element. In SQL if the checkbox is exclusive (e.g. a radio button) then the result is stored in a single column, if multiple checkes are allowed it is stored as a JSON Array column.


[^1]: variable numbers must start with a letter, may contain numbes but not spaces or punctation except the underscore

[^2]: See <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-issue-forms>, 

[^3]: See <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema>

