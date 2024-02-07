
# Newt YAML syntax

Newt's configuration and modeling is based on a YAML file. That YAML file has a
specific syntax.  The top level of that syntax are three properties.

application
: holds the run time configuration used by the Newt web service and metadata about the application you're creating.

models
: This holds the description of the data models in your application. Each model uses the [GitHub YAML issue template syntax](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/configuring-issue-templates-for-your-repository#creating-issue-forms)

routes
: This holds the routes for the data pipe line (e.g. JSON API and template engine sequence)


## application property

The application property itself has four properties.

port
: (optional, default is This port number the Newt web services uses to listen for request on localhost

htdocs
: (optional) Directory that holds your application's static content

metadata
: (optional) This holds metadata about your application using the [CITATION.cff](https://citation-file-format.github.io/) YAML syntax under metadata.

templates
: (optional) The directory where you put your application templates (e.g. Pandoc templates or Mustache template depending which engine you're using)

sql
: (optional) Teh directory where the SQL code is generated


## models property

## routes


