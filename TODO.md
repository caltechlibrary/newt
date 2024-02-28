
# Action Items

## Bugs

- [ ] Newt Router needs to bubble up the HTTP error code from the last retrieved  HTTP response in the pipeline

## Next

- [ ] Nail down the second prototype YAML syntax
- [x] There should be a "newt" command that wraps the router, generator and mustache engine in the way the go coommand wraps govet, gofmt, etc. This will be convenient for development
- [ ] I need to implement the second prototype code generator once I've debugged the Newt YAML syntax
- [ ] An "object manager" takes a model and maps it to the SQL JSON storage model of the generated code (id, created, updated, object), it also validates the submitted object against that model
- [X] (rethought the application concept in favor of a single YAML file) Should Application metadata really be it's own top level attribute? Wouldn't having a service that reads a codemeta.json or CITATION.cff make more sense in a service oriented architecture?
- [X] (one configuration used by all Newt tools) Should Newt Router, Newt Mustache and Newt Generator use separate YAML files? or a combined file?
  - [X] (future prototype can do OS level suggested conf files) Using a combined file would make it easy to generate startup scripts or systemd configurations
- [X] Decide what is logged by default for Newt Mustache and Newt Router
- [X] Decide what is logged in "debug" or "verbose" model by Newt Mustache and Newt Router
- [ ] Generate SQL confirming to the style suggestion in "The Art of Postgres" (link as a reference in Newt documentation)
- [X] (options and environment can solve this) Writing the URL for a localhost service can be tedious and obscure what is happening, create an example where you use a environment variable or application option to express the service name to a variable that can then be reference in the URL pattern
- [X] Adopt GitHub's YAML Issue Syntax for describing models
  - [x] evaluate the DSL that Newt's existing has to see if it still makes sense (probably doesn't)
  - [x] Can the model/type DSL be made compatible with [GitHub YAML issue template schema](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema)? Or should it be replaced by it?
- [x] Present at Code4Lib meetup, July 14, 2023
- [ ] Demo second prototype for DLD developer group, newt-presentation2.md
- [ ] Create birds demo for 2nd Prototype (using Postgres+PostgREST, Newt Router and Newt Mustache)
- [ ] Create postcards demo, armchair archive example

## Someday, maybe

- [ ] seeing the data past thru a pipeline will be an important part in debugging problems. Newt Inspector could adress that as a pass thru web service that output debugging info to the log output.
- [ ] An OCFL service would allow Newt to support building repository like applications. U. C. Santa Barbara has some Go packages related to this.
- [ ] Newt should support development with SQLite 3 databases
  - [ ] Look at sqlite2rest (a Python project, automatic REST API for SQLite databases), work up a demo using it with Newt server
  - [ ] Should a write a newtsqlite service?
- [ ] Can I extend Newt to support file uploads?
  - [ ] Should this be a separate service, a stage in the pipeline?
  - [ ] Should I integrate S3 protocol support to allow file upload handling?
