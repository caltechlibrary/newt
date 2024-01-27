---
title: "newtpg(1) user manual | 0.0.6 cef35c8"
pubDate: 2024-01-26
author: "R. S. Doiel"
---

# NAME

newtpg

# SYNOPSIS

newtpg [CONFIG_FILE]

# DESCRIPTION

**newtpg** is a command line tool for generate SQL suitable to bootstrap a microservice implemented with PostgREST and Postgres.  It uses the same YAML file as the Newt web service using the namespace and models attributes to render table structure, views and funcitons to support basic CRUD and list operations in SQL.

**newtpg**'s configuration uses a declaritive model expressed in YAML.  It can also allow environment variables read at start up to be part of the data for mapping JSON data source requests. This is particularly helpful for supplying access credentials. You do not express secrets in the **newtpg** YAML configuration file. This follows the best practice used when working with container services and Lambda like systems.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-dry-run
: Load YAML configuration and report any errors found

Newt has some experimental options to render Postgres dialect of
SQL from a YAML file containing models. These options will render SQL
suitable to bootstrap a Newt+Pandoc+Postgres+PostgREST based project.

-setup
: This renders a SQL document suitable for bootstraping Postgres+PostgREST access

-models
: This renders a SQL file to bootstrap modeling data with Postgres+PostgREST

-models-test
: This renders a SQL file to bootstrap writing SQL tests for Postgres+PostgREST


# CONFIGURATION

**newtpg** looks for two attributes in the Newt YAML file.

namespace
: This indicates the Postgres schema associated with your application

models
: This is a list of models that map to tables in your Postgres schema/database.

The **models** attribute holds a list of models expressed in Newt's data model DSL. Models are optional but can be used to by Newt to generate bootstrap SQL code for use with PostgREST+Postgres. This is very experimental (2024) and likely to change as usage of Newt increases. Each model has a `model` attribute holding the models name (conforming to a variable name found in langauges like JavaScript, Python, or Lua). Each model also contains a `var` attribute which is a list of key/value pairs. The key/value pairs are made from a variable name (key) and type definition (value). The type definitions are mapped to suitable Postgres SQL schema when generating table definitions. Example models used for groups and people metadata. The asterix at the end of a type string indicates this is to be used as a key when looking up the object.

~~~yaml
namespace: groups_and_people
models:
- model: cl_person
  var:
    family_name: String
    given_name: String
    orcid: ORCID
    ror: ROR
    created: Timestamp
    updated: Timestamp
- model: cl_group
  var:
    cl_group_id: String*
    short_name: String
    display_name: String
    description: Text
    contact: EMail
    created: Timestamp
    updated: Timestamp
    founded: Date 2006-01-02
    disbanded: Date 2006-01-02
    approx_founding: Boolean
    active: Boolean
    website: URL
    ror: ROR
    grid: String
    isni: ISNI
    ringold: String
    viaf: String
~~~

# Defining variables

Each model can have its own associated set of variables. Variables are "typed".  The code for type definitions includes validation. When a variable is detected in a request path or form data it is vetted using it's associated type. Only if the variables past validation are they allowed to be used to assemble a request to a JSON data source. 

Variables are defined in the YAML "var" attribute. Here's an example "var" definition defining three form variables for a route. The variable names are "bird", "place" and "sighted" with the types "String", "String" and "Date". The "bird" variable is also a "key" for the table so has its type end in an asterix.

~~~
var:
  bird: String*
  place: String
  sighted: Date 2006-01-02
~~~

# Variable types

String
: Any sequence of characters. If the variabe is embedded in a path then "/" will be used to delimited path parts and would not be passed into the variables value.

Date
: (default) A year, month, day string like 2006-01-02

Date 2006
: A four digit year (e.g. 2023)

Date 01
: A two digit month (e.g. "01" for January, "10" for October)

Date 02
: A two digit day (e.g. "01" for the first, "11" for the eleventh)

Basename
: A file's basename (filename without an extension)

Extname
: A file's extension (e.g. ".html", ".txt", ".rss", ".js")

Isbn10
: An ten digit ISBN

Isbn13
: A thirteen digit ISBN

Isbn
: An ISBN (either 10 ro 13 digit)

Issn
: An ISSN

DOI
: A DOI (digital object identifier)

Isni
: An ISNI

ORCID
: An ORCID identifier
 
NOTE: The current names associated with types will likely change
as the prototype **newtpg** evolves. It is planned for them to be
stable if and when we get to a v1 release (e.g. when we're out of the
prototype phase).


