
# Action Items

## Bugs

- [ ] In the YAML route generation the identifier isn't getting passed for updates
- [ ] Newt command runner doesn't seem to be mapping ports correctly generated code and Newt YAML
  - [ ] postgres.conf generation needs to have password set/generated for the authenticator account generate in setup in setup.sql
  - [ ] DSN string generated for access PostgREST in `.applications.postgres.dsn` needs to have port match the port Postgres is listening on
  - [ ] The DSN used in the postgres.conf file needs to match port set in `.applications.postgres.port`
- Interactive UI is not consistent, menus aren't obvious how to use
  - [ ] Attribute mapping is inconsistent with the rest of the interactive UI, e.g. add attribute should accept the attribute name to add, it needs to prompt for the value
  -  [ ] Interactive UI needs to be consistent in what enter without command means or become "q" everywhere
  - [ ] "object key" and "primary key" should be one or the other, confusing when editing elements
    - "identifier" seems like the best choice, I think it is something that the something library community understands
  - [X] The address becomes lower case in conversation UI (`.app_metadata`)
- [ ] Check should summarize each service not just template engine or only show errors and report OK otherwise
  - [ ] Missing information about router, should indicate port as well as if htdocs is set
- [X] In generated Newt YAML file the bang line is incorrect (missing bang)
- [ ] When creating a new object, the resulting JSON should return the identifier then a template can trigger a redirect in the browser to a read end point using the identifier

## Next for the third prototype

NOTE: X is completed, P is partial completion, D (dropped) from implementation plans

- [ ] Evaluate how to organize root YAML to support a selection of backends, e.g. Postgres+PostgREST or Dataset+datasetd. 
  - It is important to be able to evolve the back end without causing a major code refactor each time
  - [ ] Look at `just` tool written in rust and see what might need to be adjusted in the context of Newt application pipelines
  - [ ] Decide what fields the runnable applications need for both running and configuration generation (or drop generating postgrest.conf)
  - [ ] The list of applications for the backend stack are ordered and reflected in the startup sequence
  - [ ] Is there a specific set of roles the backend components play?
- [ ] The `*_response.hbs` templates seem unnecessary, the "read" template given the data can is the same thing
  - [ ] The "response" pages could be implemented as a redirect to the read end point
  - [ ] It can be done as a template using that generates a page redirect via head, meta elements
  - [ ] It'd be nice to be able to have the pipeline handle the redirect as a psuedo service
- [ ] First step in process should be modeling the data
- [ ] Drop "config" verb, this should be automatic if the ports need to be changed then that is a whole other level of knowledge needed
  - [ ] Presume the stack is Postgres + PostgREST
  - [ ] Presume Deno is available and the ports are free for the validator service
    - [ ] You can have a single validator service generated so it takes only a single port, models would be included via imports
  - [ ] Presume nte is used
  - [ ] Generate systemd service file(s) with Bash wrappers
- [ ] Rename "generate" to "build" since that is actually what is happening
- [ ] Need to Generate TypeScript validator service, tsgen.go??
- [ ] Need to generate a Bash (bashgen.go???) for running the database setup and configuration or have this done "automagically" as a step done by `newt` tool
- [ ] Need to generate a `deno.json` file suitable for managing the project and for running via Newt command

- [ ] Review [go-webui](https://github.com/webui-dev/go-webui) for implementing a GUI for newt
- [ ] Adding modeler.go shows I need to cleanup code and normalize an API info working with models, updating routes and templates based on updated models list.
- [ ] Postgres configuration in yamlgen.go needs to include a DSN to make the connection for loading our SQL
- [ ] I need to implement the third prototype code generator once I've debugged the Newt YAML syntax
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

## Someday, maybe

- [ ] Update `newt config` to support Dataset+datasetd
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
