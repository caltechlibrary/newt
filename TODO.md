
# Action Items

## Bugs

## Next

- [ ] Present at Code4Lib meetup, July 14, 2023
- [ ] Demo for DLD staff
- [ ] Decide if I want a "redirect on success", "redirect on fail" columns in the routes CSV file. This could be used to make webform integration smoother for simple webform handling (it would save using a dummy Pandoc template to has redirect handler)
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
    - [ ] Does the environment need to merge with objects sent to data API via POST, PUT, PATCH?
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

- [ ] Have Newt read a YAML file expression a declarative presentation of routes and form validation
	- this could function like the backend of forms processing service
- [ ] Form validation could be  used to generate default HTML markup for the form, might standardize that code a bit in the process (e.g. allow a vanilla JavaScript integration for complex form fields)
- [ ] A tool that would Convert CSV, Excel, Open Office Spreadsheets into YAML defining our routes table and forms with their validation rules leveraging the router DSL types
- [ ] Support Excel and Open Office Spreadsheets
	- [ ] If form validation is supported in Newt then you could have sheet "routes" that defines routes and other sheets defining forms (pairs of input element name attributes and types from route dsl)
- [ ] Add a declaritive way to validate form input in Newt without sending it to a JSON data API for validation, this would improve POST, PUT, PATCH handling in the same way URL vetting is implemented
	- [ ] forms could be validated from a simple YAML file using the router DSL types paired with a input element's name attribute
- [ ] Add file upload support to Newt via an S3 like service call

