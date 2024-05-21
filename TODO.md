
# Action Items

## Bugs

- [ ] Newt modeler needs to indicate at each sub menu where you're at. E.g. when editing a model, it should clearly identify you're on the model editing menu, when editing elements of a model it should show the model name and element name, when editing the attributes of an element then the model id, element id as well as the attribute id needs to be clearly displayed
- [ ] When display a list of key/value pairs I need to show the values too, this will make it easier to understand when you need to modify a model, element, attribute, options, etc.
- [ ] Newt Router currently should restrict POST to urlencoded or application/JSON encoded but explicitly disallow multipart form encoding since Newt Router doesn't support file uploads yet.
- [X] When using Newt to run services, each services needs to display the port it is running on.
- [X] Need to wrap the web forms in the appropriate form element and attributes
- [X] Select partials in page.tmpl isn't viable the way I thought I could implemented it, need to figure out a better way, 
- [X] Newt Router needs to bubble up the HTTP error code from the last retrieved  HTTP response in the pipeline
- [X] decide if it is object name then action or action then object name, I've seem to have flipped flopped around on this in the code.
- [X] Newt Config should not add routes or templates if they are already defined in the previously read YAML

## Next for the second prototype 

NOTE: X is completed, P is partial completion, D (dropped) is skipping implementation

- [ ] In selectMenuItem there should be a way to turn off numbered items when they don't make sense (e.g. at the modify model menu or modify element menu)
    - [ ] could do this by having two different types of menu, selectMenuItem, selectMenuById or some such (not happy with these names)
- [ ] If I am going to support numbered item selection then I should also allow the order of the number and action to be reverse, e.g. "1 modify" would indicate modify item 1. This would give the numbers more obvious meaning and useful ness, otherwise it's just noise in the display.
- [ ] Rename Newt verb "init" to "config" since that better describes what it does
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
- [ ] If you run Newt Generator from the `newt` command it should assign predictable filenames for SQL files, PostgREST configuration and Mustache templates. This could insure the turn key operation of a bare prototype. It needs to align with what `newt init` generators.
  - SQL generated should organize files by function, e.g. `{{model_name}}_{{action_name}}.sql`, `*_setup.sql` would create database, schema, roles, tables, `*_access_control.sql` would be regenerated and map roles and permissions for any functions found created by `{{model_name}}_{{action_name}}.sql`. I need to figure out which tables to query to identify the functions that are available in the metadata schema so that the mapping can be complete. 
- [ ] Verify we have a simple developer workflow
  - [ ] `newt init` generate a default YAML for project
    - [ ] `newt init` allow automatic generation of the project code base?
  - [ ] `newt generator` generated the code for the project
  - [ ] `newt run PROJECT_YAML` should be able to run the rendered project.
- [ ] Newt Check should detect mis-aligned form names when paired with a Postgres function.
  - [ ] Need to document how the web form input "names" attributes need to match the SQL functions, it is obvious if you understand PostgREST but not so obvious if you are unfamiliar with either Postgres functions or PostgREST
- [ ] I still think the string representation of UUID is problematic in the model generation, I need to decide how to deal with this
  - [X] Include Python functions in Postgres to support shorter unique identifiers
    - [ ] Add Python short object id via Newt Generator SQL for models.
  - [D] Implement an "object manager" that will sit in front of PostgREST that will wrap the objects with their object id, created and updated timestamps, this could enforce REST method behaviors but I would need to tie models to routes to do so
- [P] Nail down the second prototype YAML syntax
- [ ] I need to implement the second prototype code generator once I've debugged the Newt YAML syntax
  - [X] setup.sql
  - [X] models.sql
  - [X] models\_test.sql
  - [X] postgrest.conf
  - [ ] create mustache template
  - [ ] read mustache template
  - [ ] update mustache template
  - [ ] delete mustache template
  - [ ] list mustache template
- [ ] `newt` runner should be able to manage a postgrest instance. This will simplify using Newt in the development setting
  - [X] `newt` need to track the pid of the external process, then folder that into the signal handlers (using a go routine to start it don't think I need to track this, though it would be nice to log it.
- [ ] `newt init`, an interactive Newt YAML file generator, need to decide on approach
  - [X] could be done as a cli interactive tool
  - [ ] could be done as a GUI form based application
- [ ] A "newt object manager" takes a model, validates it and if OK sends the result to the next stage in the pipeline for storage. It should deal with things like unpacking shortened UUID using a base encoding scheme
- [ ] Present to DLD and interested staff
- [ ] Present/announce via SoCal Code4Lib (recorded or in person presentation)
- [ ] Create birds demo for 2nd Prototype (using Postgres+PostgREST, Newt Router and Newt Mustache)
- [ ] Create postcards demo, armchair archive example
- [ ] Implement Thesis Management System core in Newt (not email features)
- [ ] Implement COLD core
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
