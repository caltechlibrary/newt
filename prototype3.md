
# Prototype 3

## Questions and Explorations

- Generated TypeScript middleware targetting Deno (e.g. validation, data enhancement)
- Explore Dataset/datasetd as a simpler JSON data source to Postgres+PostgREST

## Desirable Features

- [ ] newt (developer tool)
  - [X] implemeted a conversational UI
  - [X] configuration manager (e.g. target stack, services required by application)
  - [X] model (model your data)
    - [X] Add app metadata to AST and update newt YAML syntax docs
  - [ ] Workflow testing and verification
    - [ ] Model
    - [ ] Generator
    - [ ] Run
    - [ ] Config (optional) available if needed
    - [ ] Check (optional) available if needed
- [ ] generate the following
    - [ ] SQL for Postgres+PostgREST
    - [ ] Handlebar template targetting Newt's template engine
    - [ ] Automatic route management based on models and simple pipelines
    - [ ] Automatic template mapping based on models and simple pipelines
    - (time permitting) [ ] TypeScript for model validator service
- run (target stack based on new implementation AST and data models)
  - [ ] Test building a new model from scratch and run
  - [ ] Test updating and existing model and run
  - [ ] Document the run command more fully
- data router  (ndr)
  - data pipelines
    - [ ] validator integration via generated code (e.g. TypeScript run by Deno?)
    - [ ] Postgres+PostgREST integration via generated code and configuration
    - [ ] Templates render data acurately (needs automated testing)
  - static file service
    - [ ] Does it make sense to auto generate some static content like the about page and dashboard page?
  - [ ] related routes linked by an id
- [ ] Generator application code for two stacks (newt)
  - [X] SQL and config Postgres+PostgREST
  - [X] Handlebar Templates and partials for the Template Engine
  - [ ] routes to enabled services using their port settings in the YAML file
  - [ ] template mappings
  - [ ] Generate a TypeScript web service for implementing serverside validation of submitted content for all models
  - [ ] Generate JavaScript for use with web form validation in the browser
- [X] Newt Template engine (nte)
  - [X] A template engine that works as the last stage of a JSON data pipeline

## Constraints

- Simple data models only
- No file upload support
- Models targets Postgres+PostgREST
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

