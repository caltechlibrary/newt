---
title: "newtgenerator(1) user manual | 0.0.7-dev f7a98be"
pubDate: 2024-02-20
author: "R. S. Doiel"
---

# NAME

newtgenerator

# SYNOPSIS

newtgenerator CONFIG_FILE GENERATOR [PARAMS]

# DESCRIPTION

**newtgenerator** is a command line tool for generate SQL and configuration files used in bootstrapping a web service. Currently the following generatation is supported.

- Postgres SQL for use with PostgREST
- PostgREST configuration file
- Mustache templates for use with Newt Mustache

**newtgenerator**'s configuration uses a declaritive model expressed in YAML.  It can also allow environment variables read at start up to be part of the data for mapping JSON data source requests. This is particularly helpful for supplying access credentials. You do not express secrets in the **newtgenerator** YAML configuration file. This follows the best practice used when working with container services and Lambda like systems.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

# CONFIG_FILE

**newtgenerator** uses the Newt YAML syntax. What follows are those properties of specific
relevance to **newtgenerator** configuration.

### Top level properties

These are the top level properties in YAML files.

applications
: (optional) holds the run time configuration used by the Newt applications.

models
: (required by newtgenerator) This holds the description of the data models in your application. Each model uses the [GitHub YAML issue template syntax](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/configuring-issue-templates-for-your-repository#creating-issue-forms) (abbr: GHYTS)

### The applications property

newtgenerator
: this contains configuration for the Newt Generator, e.g. port, namespace

options
: holds key value pairs of which can be referenced in the values of models, routes and templates.

### newtgenerator property

namespace
: newtgenerator uses this in the SQL generated for setting up Postgres+PostgREST

### the "models" property

Models holds a list of individual models used by our data pipelines. The models are by Newt code generator and the Newt router. Models defines a superset of the GitHub YAML issue template syntax (abbr: GHYTS).

### a model object

The model object is based largely on GitHub YAML issue template syntax with a couple extra properties that are Newt enhancements.

id
: (required, newt specific) this is the name identifying the model. It must conform to variable name rules[^21]

The following properties are based on the GitHub YAML issue template syntax[^22] (abbr: GHYTS)

name
: (required: GHYTS, optional: newt) Must be unique to use with GitHub YAML issue templates[^22]. In Newt it will be used in populating comments in generated SQL

description
: (required: GHYTS, optional: newt) A human description of the model, It will appear in the web form and SQL components generated from the model

body
: (required) A a list of input types. Each input type maps to columns in SQL, input element in web forms and or HTML elements in read only pages

#### a model's input types

This is based on GitHub YAML issue template (abbr: GHYTS) input types[^23]. 

id
: (required) an identifier for the element. Must conform to variable name rules[^21]. It is used to SQL as a column name and in web forms for the input property.

type
: (required) Identifies the type of elements (input, textarea, markdown, checkbox, dropdown).

attributes
: (optional) A key-value list that define properties of the element. These used in rendering the element in SQL or HTML.

validations
: (optional, encouraged) A set of key-value pairs setting constraints of the element content. E.g. required, regexp properties, validation rule provided with certain identifiers (e.g. DOI, ROR, ORCID).

#### input types

Both the routes and models may contain input types. The types supported in Newt are based on the types found in the GHYTS for scheme[^23]. They include

markdown
: (models only) markdown request displayed to the user but not submitted to the user but not submitted by forms. 

textarea
: (models only) A multi-line text field

input
: A single line text field. This conforms to value input types in HTML 5 and can be expressed using the CSS selector notation. E.g. `input[type=data]` would be a date type. This would result in a date column type in SQL, a date input type in HTML forms and in formatting other HTML elements for display.

dropdown
: A dropdown menu. In SQL this could render as an enumerated type. In HTML it would render as a drop down list

checkboxes
: A checkbox element. In SQL if the checkbox is exclusive (e.g. a radio button) then the result is stored in a single column, if multiple checks are allowed it is stored as a JSON Array column.

Newt may add additional types in the future.

[^21]: variable numbers must start with a letter, may contain numbers but not spaces or punctuation except the underscore

[^22]: See <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-issue-forms>, 

[^23]: See <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema>

# GENERATOR

Currently three types of generators are being implemented in the 2nd Newt Protoype. This parameter
lets you set which one you are using. It is required.

# PARAMS

Some generators may need additional information. This is true of the Mustache template generator.
This is where you would provide those parameters. The Postgres and PostgREST generators do not
currently take any parameters.

# EXAMPLES

~~~yaml
applications:
  newtgenerator:
    namespace: people # E.g. "people" Namespace to use generating Postgres SQL
models:
  - id: people
    name: People Profiles
    description: |
      This models a curated set of profiles of colleagues
    body:
      - id: people_id
        type: input
        attributes:
          label: A unique person id, no spaces, alpha numeric
          placeholder: ex. jane-do-007
        validations:
          required: true
      - id: display_name
        type: input
        attributes:
          label: (optional) A person display name
          placeholder: ex. J. Doe, journalist
      - id: family_name
        type: input
        attributes:
          label: (required) A person's family name or singular when only one name exists
          placeholder: ex. Doe
        validations:
          required: true
      - id: given_name
        type: input
        attributes:
          label: (optional, encouraged) A person's given name
          placeholder: ex. Jane
      - id: orcid
        type: input
        attributes:
          label: (optional) A person's ORCID identifier
          placeholder: ex. 0000-0000-0000-0000
        validations:
          pattern: "[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]"
      - id: ror
        type: input
        attributes:
          label: (optional) A person's ROR identifing their affiliation
      - id: email
        type: "input[type=email]"
        attributes:
          label: (optional) A person public email address
      - id: website
        type: "input[type=url]"
        attributes:
          label: (optional) A person's public website
          placeholder: ex. https://jane.doe.example.org
~~~

