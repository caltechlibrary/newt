
# Prototype 3

## Features

- newt developer (implemented as a native GUI or Web UI)
  - implemeted as web UI and/or native GUI (e.g. via Fyne)
  - configuration manager (e.g. target stack, services required by application)
    - searchengine configurations and indexing models
    - metadata storage (e.g. Postgres+PostgREST or SQLite3+PL/Lua)
    - object storage (e.g. S3 protocol service integration) 
  - model (model your data)
    - continue to support Prototype 2 modeling
    - file type elements integrating an OCFL S3 service
  - generate (generate code for Postgres+PostgREST Stack or Python+Flask stack)
  - run (target stack based on AST and data models)
- newt server
  - validator 
  - data pipelines
    - including non-HTTP routes to an embedded template engine and Lua lamba functions
  - static file service
  - related routes linked by an id
  - embedded Lua lambda functions available (useful for custom validation requirements)
  - embeddd SQLite3 support
    - implemetns PL/Lua for database functions and procedures
- newt object server
  - An OCFL enable S3 service for persistant storage of large objects
- Generator application code for two stacks
  - Postgres+PostgREST and Mustache templates
  - SQLite3+PL/Lua and Mustache templates


## Constraints

- Simple data models
- No file upload support
- Modeler assumes the Postgres+PostgREST stack
- Newt services required to restricted to localhost
- Access control defered to front end web service
  - Apache 2 or NginX plus Shibboleth and acls for URL paths

## Data Types

Newt presents HTML 5 form elements as the native (base) set of data types. These imply
predictable SQL column types when persisting the data in Postgres.

Additional data types implemented via HTML 5 input element's pattern attribute (e.g.
similar to validation found in Python's idutils package)

Complex data types implented as textarea containing YAML (for fallback) and using
web components to provide more intitive user experience if JavaScript is available
in the browser.

## Proposed Roadmap for Prototype 2

- Targeted for use in developing [cold](https://caltechlibrary.github.io/cold) for Caltech Library
- Demonstration to DLD group
- Demonstration to Caltech Library staff
- Demonstration for SoCal Code4Lib Meetup
- Promote via 


