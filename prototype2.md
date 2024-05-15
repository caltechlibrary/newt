
# Prototype 2

## Features

- newt developer tool
  - config (configure the stack used with your Newt application)
  - model (model your data)
  - generate (generate configuration file, SQL and templates)
  - run (Newt Router, Newt Mustache and PostgREST)
- newt router (a data router web service)
  - validator (a service that read the models and tests the submissions for valid data)
  - data pipelines
  - static file service
  - related routes linked by an id
- newt mustache (a Mustache template engine inspired by Pandoc's server mode)
  - stateless template engine
  - loads templates at startup
  - support partial templates
  - templates referenced through unique path
  - related templates linked by an id
- Extensible to support for other JSON data sources through editing Newt YAML

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


