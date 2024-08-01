
# Prototype 3

## Questions and Explorations

- Generated middleware using TypeScript running in Deno (e.g. validation, data enhancement)
- Explore Dataset/datasetd as a simpler JSON data source to Postgres+PostgREST

## Desirable Features

- newt (developer tool)
  - implemeted with web UI
  - configuration manager (e.g. target stack, services required by application)
    - searchengine configurations and indexing models
    - metadata storage (e.g. Postgres+PostgREST or SQLite3+PL/Lua)
    - object storage (e.g. S3 protocol service integration)
  - model (model your data)
    - continue to support Prototype 2 modeling
    - file type elements integrating an OCFL S3 service
  - generate (generate code for Postgres+PostgREST Stack or Python+Flask stack)
  - run (target stack based on AST and data models)
- data router 
  - validator integration via generated code (e.g. TypeScript run by Deno)
  - data pipelines
  - static file service
  - related routes linked by an id
- Generator application code for two stacks
  - (default) Dataset collection with Dataset YAML for datasetd
  - (optional) Postgres+PostgREST
- a template engine to easily turn a JSON object into HTML or other formats

## Constraints

- Simple data models only
- No file upload support
- Modeler targets Dataset+Datasetd or Postgres+PostgREST as option
- Newt services are restricted to localhost
- Access control defered to front end web service
  - Apache 2 or NginX plus Shibboleth and acls for URL paths

## Data Types

Newt presents HTML 5 form elements as the native (base) set of data types. These imply
predictable SQL column types when persisting the data in Postgres.

Additional data types implemented via HTML 5 input element's pattern attribute (e.g.
similar to validation found in Python's idutils package).

Complex data types implented as textarea containing YAML (for fallback) and using
web components to provide more intitive user experience if JavaScript is available
in the browser.


## Proposed Roadmap for Prototype 3

- Targeted use cases (at Caltech Library)
  - [cold](https://github.com/caltechlibrary/cold)
  - [tms](https://github.com/caltechlibrary/tms)
- TBD Demonstration to DLD group
- TBD Demonstration to Caltech Library staff
- TBD Demonstration for SoCal Code4Lib Meetup
- TBD Promote via recorded presentation
