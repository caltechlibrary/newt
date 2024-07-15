
# Action Items

## Bugs

- [ ] Newt Run bugs
- [ ] Newt Router bugs
  - [ ] Newt Router should restrict POST that expect support for multipart encoding since Newt does not support file uploads
- [ ] Newt Generator
  - [ ] Newt Generator needs to generate the SQL to validate inputs based on their element types. If this is too complex then the router needs to do the validation before sending through the pipeline (that will require a change to the AST (e.g. `validate_with:<MODEL_NAME>` attribute)
- [ ] Newt Template Engine bugs

## Next for the second prototype 

NOTE: X is completed, P is partial completion, D (dropped) from implementation plans

- [ ] Review [go-webui](https://github.com/webui-dev/go-webui) for implementing a GUI for newt
- [ ] Review [deno](https://deno.land) and [Typescript](https://www.typescriptlang.org/) as a target for generating a validation service for our models
    - deno can compile typescript/javascript into an executable
    - these could be installed in a bin directory next to the newt YAML project file
    - the typescript can can run as a a validation service in the pipeline for the form processing
    - I would need to map the service name sainly
    - Investigate using embedding WASM code into Newt pipeline for validation, I could then generate TypeScript validation code compile to WASM using javy and transpile to JavaScript for browser (review how fresh doing this)
- [ ] The configuration Newt YAML needs to specify PostgreSQL+PostgREST as one target and Deno+datasetd as another back end with/without Newt router
- [ ] The generator for Newt needs to generate TypeScript apps for validation and could also be used to generate standalone Deno+datasetd based webapps
- [ ] See if I can turn handlebarjs into a WASM module, then implement a NewtHandlebars service, mutch easier working with Handlebars then Mustache fo r the types of templates we do.
- [ ] Newt should support dataset based applications since datasetd exists, newt wouldn't model apps in SQL in this case but would use the models to create the templates, validation, adata pipelines and routing
- [ ] Adding modeler.go shows I need to cleanup code and normalize an API info working with models, updating routes and templates based on updated models list.
- [ ] Postgres configuration in yamlgen.go needs to include a DSN to make the connection for loading our SQL
- [ ] Newt Router needs to validate it's inputs for POST, PUT, PATCH against a specific data model
    - [ ] Need to decide on how to vet submitted web forms, does vetting code get generate in the RDBMS, do I include a model in the route definition or should this be done as an independent part of the pipeline? Or all of they above?
- [ ] Newt needs a web hook service that can be placed in the pipeline to trigger a non-pipeline action, sorta of like the Unix tee command
  - use cases
    - the web hook would receive the JSON data from the previous service in the pipeline
    - it could trigger another URL pipeline
    - it could trigger running a program on the local system (e.g. trigger search engine indexing run after record change)
    - it could insert an action into an event/job queue
- [ ] I need to implement the second prototype code generator once I've debugged the Newt YAML syntax
  - [ ] setup.sql
  - [ ] models.sql
  - [ ] postgrest.conf
  - [ ] create related templates
  - [ ] read template
  - [ ] update related templates
  - [ ] delete realted template
  - [ ] list template
- [ ] `newt generate` needs to be brought into line with changes made in the modeler development
  - [ ] I need to have tests for all the generating functions to avoid regressions
- [ ] Create a useful set of demos
    - [ ] Create birds demo for 2nd Prototype (using Postgres+PostgREST, Newt Router and Newt Mustache)
    - [ ] Garden management app
    - [ ] A text only blog
    - [ ] Create postcards demo, armchair archive example
    - [ ] Implement Thesis Management System core in Newt (not email features)
    - [ ] Implement COLD core
- [ ] Present to DLD and interested staff
- [ ] Present/announce via SoCal Code4Lib (recorded or in person presentation)
- [X] `newt config`, an interactive Newt YAML file generator, need to decide on approach
  - [X] could be done as a cli interactive tool
  - [D] could be done as a GUI form based application
- [X] `newt` runner should be able to manage a postgrest instance. This will simplify using Newt in the development setting
  - [X] `newt` need to track the pid of the external process, then folder that into the signal handlers (using a go routine to start it don't think I need to track this, though it would be nice to log it.
- [X] In selectMenuItem there should be a way to turn off numbered items when they don't make sense (e.g. at the modify model menu or modify element menu)
    - [X] could do this by having two different types of menu, selectMenuItem, selectMenuById or some such (not happy with these names)
- [X] Rename Newt verb "init" to "config" since that better describes what it does
- [X] If you run Newt Generator from the `newt` command it should assign predictable filenames for SQL files, PostgREST configuration and Mustache templates. This could insure the turn key operation of a bare prototype. It needs to align with what `newt init` generators.
  - SQL generated should organize files by function, e.g. `{{model_name}}_{{action_name}}.sql`, `*_setup.sql` would create database, schema, roles, tables, `*_access_control.sql` would be regenerated and map roles and permissions for any functions found created by `{{model_name}}_{{action_name}}.sql`. I need to figure out which tables to query to identify the functions that are available in the metadata schema so that the mapping can be complete. 
- [X] Verify we have a simple developer workflow
  - [X] `newt config` generate a default YAML for project
    - [X] `newt config` allow automatic generation of the project code base?
  - [X] `newt generator` generated the code for the project
  - [X] `newt run PROJECT_YAML` should be able to run the rendered project.
- [X] Nail down the second prototype YAML syntax
- [X] Port attributes in the struct need to all be either string or int (probably int), it'll make the code read better to be consistent
- [X] There should be a "newt" command that wraps the router, generator and mustache engine in the way the go command wraps govet, gofmt, etc. This will be convenient for development
- [X] (rethought the application concept in favor of a single YAML file) Should Application metadata really be it's own top level attribute? Wouldn't having a service that reads a codemeta.json or CITATION.cff make more sense in a service oriented architecture?
- [X] (one configuration used by all Newt tools) Should Newt Router, Newt Mustache and Newt Generator use separate YAML files? or a combined file?
  - [X] (future prototype can do OS level suggested conf files) Using a combined file would make it easy to generate startup scripts or systemd configurations
- [X] Decide what is logged by default for Newt Mustache and Newt Router
- [X] Decide what is logged in "debug" or "verbose" model by Newt Mustache and Newt Router
- [X] Generate SQL confirming to the style suggestion in "The Art of Postgres" (link as a reference in Newt documentation)
- [X] (options and environment can solve this) Writing the URL for a localhost service can be tedious and obscure what is happening, create an example where you use a environment variable or application option to express the service name to a variable that can then be reference in the URL pattern
- [X] Adopt GitHub's YAML Issue Syntax for describing models
  - [X] evaluate the DSL that Newt's existing has to see if it still makes sense (probably doesn't)
  - [X] Can the model/type DSL be made compatible with [GitHub YAML issue template schema](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema)? Or should it be replaced by it?
- [X] Present at Code4Lib SoCal meetup, July 14, 2023
- [X] Present at Code4Lib SoCal meetup, April 19, 2024
- [X] Cleanup data models, abandon strict following of GHYITS, I don't need that level of complexity and simplifying data model and AST will be easier to document leaving out the GHYITS connection

## Someday, maybe

- [ ] If I am going to support numbered item selection then I should also allow the order of the number and action to be reverse, e.g. "1 modify" would indicate modify item 1. This would give the numbers more obvious meaning and useful ness, otherwise it's just noise in the display.
- [ ] A "newt object manager" takes a model, validates it and if OK sends the result to the next stage in the pipeline for storage. It should deal with things like unpacking shortened UUID using a base encoding scheme
- [ ] Generate Python code that can be compiled into WASM module that is run directly in Newt Router without need to an external network call. This might be useful for implementing custom validation for the models or if they want to write a proxy service to outside localhost.
- [ ] Explore the idea of "complex models", these would be composed by either combining or aggregating lists of other models. This would let me mimic RDM/ArchivesSpace rich data models from Newt
- [ ] If the Newt Router is going to proxy file uploads it will need to specify the route as multipart encoded, if Newt Router can generate object ids when creating new objects along with the file upload then you could auto magically map the object into an S3 protocol path for file storage and further processing
- [ ] The generate could generating an Apache include conf file or an NGINX configuration file
- [ ] Solr could be specified like PostgREST in the `.applications` property, this could easy integrating Solr
- [ ] If Newt Router supported multipart encoding it could proxy to a web service that provided file upload management
- [ ] `newt build` would generate a stand alone Go based program for the application described by the Newt YAML file (think OBNC/Ofront/Ofront+ generating C then compiling the C to rendering an executable or library).
- [ ] Explore alternatives to a UUID for object identifiers, some sort of short id like RDM would be very nice.
- [ ] seeing the data past thru a pipeline will be an important part in debugging problems. Newt Inspector could address that as a pass thru web service that output debugging info to the log output.
- [ ] An OCFL service would allow Newt to support building repository like applications. U. C. Santa Barbara has some Go packages related to this.
- [ ] Newt should support development with SQLite 3 databases
  - [ ] Look at sqlite2rest (a Python project, automatic REST API for SQLite databases), work up a demo using it with Newt server
  - [ ] Should Newt supply a newtsqlite service?
- [ ] Can I extend Newt to support file uploads?
  - [ ] Should this be a separate service, a stage in the pipeline?
  - [ ] Should I integrate S3 protocol support to allow file upload handling?
