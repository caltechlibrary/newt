---
title: "Newt, a small system experiment"
author: "R. S. Doiel, <rsdoiel@caltech.edu>"
institute: |
  Caltech Library,
  Digital Library Development
description: Code4Lib Meet up, Los Angeles
urlcolor: blue
linkstyle: bold
aspectratio: 169
createDate: 2023-05-16
updateDate: 2023-06-27
pubDate: 2023-07-14
place: UCLA
date: July 14, 2023
section-titles: false
toc: true
keywords: [ "code4lib", "microservice", "Postgres", "PostgREST", "Pandoc" ]
url: "https://caltechlibrary.github.io/newt/presentation"
---

# The experiment

> Build a web application using a scaling small approach

# My Scaling Small Approach

- Minimize the moving parts
- Minimize the cognative shifts
- Leverage "off the shelf" software

# Three abstractions

- A Template engine => [Pandoc](https://pandoc.org)
- A JSON source to manage data => [Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org)
- A data router, form validator => [Newt](https://github.com/caltechlibrary/newt/)

# Limiting our cognitive shifts

- Use templates to generate HTML (Pandoc)
- Model data in SQL (Postgres) get JSON (PostgREST)
- Orchestrate our microservice conversation with YAML (Newt)

# Why SQL?

> PostgreSQL+PostgREST, my code savings plan

- Good at describing structured data
- Supports queries, views, functions, triggers, events, ...
- Allows us to model our data once

> **Minimize the source Luke!**

# The scale small experiment

Let's compare three implementations of a bird sighting website

# [birds 1](../demos/make-birds1.bash "this is a static website")

CSV file, Pandoc, 2 directories, 5 files, **72 total line count**, static site hosting. Included is a typical setup I would use for a static site project.

Lines   Files
------  ---------------
    21  [README.md](../demos/birds1/README.html)
     4  birds.csv
     6  build.sh
     7  page.tmpl
    32  htdocs/index.html

See <https://github.com/caltechlibrary/newt/tree/main/demos/birds1>

# [birds 2](../demos/make-birds2.bash "this website requires a machine")

SQL (Postgres + PostgREST), Browser JavaScript, 2 directories, 8 files, **232 total line count**, dynamic site requires hosting

Lines    Files
------   --------------
    29   [README.md](../demos/birds2/README.html)
     4   birds.csv
    34   setup.sql
    60   models.sql
    15   models_test.sql
     3   postgrest.conf
    24   htdocs/index.html
    63   htdocs/sightings.js

See <https://github.com/caltechlibrary/newt/tree/main/demos/birds2>

# [birds 3](../demos/make-birds3.bash "this website requires a machine")

SQL (Postgres + PostgREST), Pandoc, Newt, 1 directory, 7 files, 277 total line count, dynamic site requires hosting, **no JavaScript required**

Lines    Files
------   ---------------
    43   [README.md](../demos/birds3/README.html)
     4   birds.csv
    34   setup.sql
    60   models.sql
    15   models_test.sql
     3   postgrest.conf
    25   birds.yaml
    36   page.tmpl
     7   post_result.tmpl

See <https://github.com/caltechlibrary/newt/tree/main/demos/birds3>

# Newt's microservice conversation

1. web browser => Newt
2. Newt => PostgREST
3. Newt => Pandoc
4. Newt => web browser

# Newt's YAML file

- htdocs directory for static content (option)
- environment variable used to access JSON sources (optional)
- models (optional, used to bootstrap SQL code)
- route definitions
  - (optional) variable definitions (path and form data)
  - request routing details (e.g. path, method)
  - JSON source details (e.g. API URL, method, content type)
  - (optional) Pandoc template filename

See <https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/birds3/birds.yaml>

# Developer workflow

1. Define models/routes/form validation in Newt YAML file
2. Enhance models in SQL and Postgres
3. Create/update Pandoc templates
4. (Re)start PostgREST and Newt to (re)load models and routes
5. Test with our web browser

**Repeat as needed**

# Newt stack weakness

- Newt is **an experimental prototype** (June/July 2023)
- Newt's generated SQL presumes customization by developer
- Newt doesn't support file uploads
- Newt doesn't eliminate learning curves, e.g. Postgres and SQL; Pandoc; using HTTP methods; YAML

# Newt stack strengths 

> A mature foundation

- 20th Century tech
  - SQL (1974), HTTP (1991), HTML (1993), Postgres (1996)
- 21st Century tech
  - JSON (2001), YAML (2001), Pandoc (2006), PostgREST (2014)

# Experimental insights so far

- "Off the shelf" microservices can make application construction easier
- Orchestrating the data pipeline in YAML seems reasonable
- SQL turns some people off, Newt may help mitigate that
- Pandoc templates are simple to learn
- Newt stack plays well with HTML5 approaches and best practices

# An unexpected result

- Newt maintains very little "state"
- Pandoc server requires no "state"
- PostgREST maintains very little "state
- Postgres can be clustered

> The Newt stack can scale big

# Next steps for Newt?

1. Enhance SQL generation
2. Test with Solr/Elasticsearch as alternate JSON sources
3. Build staff facing applications this Summer (2023)
4. Hopefully move beyond my proof of concept in Winter (2023), Spring (2024)

# Someday, maybe ...

- Have Newt delegate file uploads to an S3 like service
- Explore integrating SQLite3 support in addition to Postgres+PostgREST
- A Newt community to share YAML, SQL and Pandoc templates

# Additional resources 

- Newt <https://github.com/caltechlibrary/newt>
- Postgres <https://postgres.org> + PostgREST <https://postgrest.org>
  - PostgREST Community Tutorials <https://postgrest.org/en/stable/ecosystem.html>
- Compiling Pandoc or PostgREST requires Haskell
  - Install Haskell GHCup <https://https://www.haskell.org/ghcup/>

# Thank you!

- This Presentation <https://caltechlibrary.github.io/newt/presentation/>
- Project: <https://github.com/caltechlibrary/newt>
- Email: rsdoiel@caltech.edu

