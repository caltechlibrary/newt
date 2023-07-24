
# Action Items

## Bugs

## Next

- [ ] Update generate SQL to confirm to the style suggestion in "The Art of Postgres"
- [ ] It would be nice to be able to define some global constants to minimize repeated strings like a common PostgREST url prefix, `http://localhost:3000` expressed as `{api}` 
- [ ] Prune required boiler plate in Newt's YAML file
    - [ ] Make sure `*_method:` defaults to "GET"
    - [ ] Make sure `api_content_type:` defaults to "application/json"
- [ ] Present at Code4Lib meetup, July 14, 2023
- [ ] Demo for DLD staff
- [x] Decide if I want a "redirect on success", "redirect on fail" columns in the routes CSV file. This could be used to make webform integration smoother for simple webform handling (it would save using a dummy Pandoc template to has redirect handler
- [ ] Create a model attribute in YAML, see if that can be used to generate a basic SQL file for integration with Postgres/PostgREST, model could also be used for validating forms but maybe "global" to the routes
    - [ ] generate PostgREST setup SQL
    - [x] generate SQL table create statement
    - [ ] generate a default view
    - [ ] generate a default create, update, delete, functions for table
    - [ ] generate Pandoc template web for using the create/update functions 
- [ ] Finish Newt prototype
    - [ ] Newt router handling for OPTIONS
    - [ ] Newt router handling for HEAD
    - [x] Newt router handling for GET
    - [x] Newt router handling for POST
    - [ ] Newt router handling for PUT
    - [ ] Newt router handling for PATCH
    - [ ] Newt router handling for DELETE
    - [ ] Test newt router handling OPTIONS
    - [ ] Test newt router handling HEAD
    - [ ] Test newt router handling GET
    - [ ] Test newt router handling POST
    - [ ] Test newt router handling PUT
    - [ ] Test newt router handling PATCH
    - [ ] Test newt router handling DELETE
    - [ ] Test newt birds3 demo and debug
    - [x] Does the environment need to merge with objects sent to data API via POST, PUT, PATCH?
    - [x] Debug interaction between static file service and Router
    - [x] Test assembled Newt static file handling
    - [x] Test newt router only handling
    - [x] implement Path DSL
    - [x] Figure out where resoling the types, values and url request should go
    - [x] Nail down where table of routes comes from
    - [x] Nail down and implement configuration file
    - [x] Add load configuration to Runner
    - [x] Add load routes to Runner
    - [x] Implement web service for Newt router
    - [x] Dry run and start up web service to Runner
- [x] Create birds demo 3 (using Postgres+PostgREST, Pandoc, Newt)
- [x] Create birds demo 1 (static site with Pandoc)
- [x] Create birds demo 2 (dynamic site with Postgres+PostgREST and JavaScript using xhr calls)
- [x] Draft presentation
- [x] Documenting ideas

## Someday, maybe

- [ ] Newt could delegate file uploads to an S3 like service, Minio is written in Go and I think supports a file streaming model
- [ ] Integrate SQLite3 support as a "JSON Data Source", might be a separate service or synthesize the results from direct access to SQLite3 database file(s)
- [ ] Consider porting Newt from Golang to Haskell for integration opportunities with Pandoc
- [ ] Create a community portal integrated with GitHub for sharing project YAML, SQL and Pandoc templates

