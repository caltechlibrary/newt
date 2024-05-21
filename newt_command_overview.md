
# Newt's `newt` command

The core of Newt is the command `newt`. It handles the following for your Newt based project. The command first parameter is an "action" The actions
let you progress from configuration, to data modeling and ending with generating code.

1. Create a Newt YAML file through an interactive dialog using the "config" action.
2. Generate and manage your models through an interactive dialog using the "model" action.
3. Verify the YAML using the "check" action.
4. Generate code (e.g. postgrest.conf, setup.sql, models.sql, various Mustache templates) using "generate" action.
5. After setting up Postgres via the SQL files you can run your application using the "run" action.

The `newt` command is a development tool. It simplifies standing up a Newt application built around the Newt Router, Newt Mustache template engine and Postgres+PostgREST.  A `newt` generated application will provide the basic create, read, update, delete and list operations you need for managing the metadata described by your models.

## Newt YAML file structure

Newt command is a tool for managing your Newt YAML file. The top level structure of a Newt application is as follows.

applications
: Holds configuration information for PostgREST, Newt Router and Newt Mustache template engine

models
: This holds the data models your application will implement, this property is inspired by the [GitHub YAML issue template syntax]()

routes
: These are the request Newt Router will manage they are descriptions of HTTP method and path along with any data pipeline processing needed to respond to HTTP request

templates
: This holds the configuration for the Newt Mustache template engine. The template engine accepts POST requests and associates Mustache templates with a given request path.  Mustache templates support the concept of "partial" templates, these are also specified in relationship to the request path and primary template association. Templates are read in with Newt Mustache starts and are not re-read until you restart Newt Mustache.

The responsibility for managing these properties is split between two Newt actions, "config" and "model".  In the next section we will cover the `newt config` which is a tool for setting the first property, applications.

Next see [newt config](newt_config_explained.md).

