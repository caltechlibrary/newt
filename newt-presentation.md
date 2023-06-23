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
updateDate: 2023-06-23
pubDate: 2023-07-14
place: UCLA
date: July 14, 2023
section-titles: false
toc: true
keywords: [ "code4lib", "meetup", "Postgres", "PostgREST", "Pandoc" ]
url: "https://caltechlibrary.github.io/newt/presentation"
---

# Scaling

scale
: a euphemism for **scaling big**, usage "google scale", "amazon scale"

- Scaling big is hard
- Scaling big can make things really complex

# Scaling, lessons 

- distributed application design
- concept of microservices
- programmable infrastructure

# Scaling

- Microservices illuminates a path to scaling small
- **We pack only what is needed!**

# Scaling small

- Limit the moving parts
- Limit the cognitive shifts
- Try to **Write less code**

# Three abstractions

- JSON source to manage data => [Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org)
- A powerful template engine => [Pandoc](https://pandoc.org)
- data router, form validator => [Newt](https://github.com/caltechlibrary/newt/)

# Cognitive shifts

- Write YAML, orchestrate our microservice conversation
- Write SQL (Postgres), yields JSON source (PostgREST)
- Write Pandoc templates, use JSON, yields HTML

# Writing less code

1. Use "off the shelf" microservices
2. Standardize templating with Pandoc
3. Take advantage of SQL

> **Got middle-ware?** We use the Newt stack, 
> Postgres+PostgREST+Pandoc+Newt

# Why SQL?

- SQL is good at describing structured data
- SQL can do queries, views, functions, triggers, ...
- SQL allows us to model our data once
- PostgREST gives us a JSON api based on the SQL we write

> **Minimize the source Luke!**
> PostgreSQL+PostgREST is a code savings plan.

# A scale small experiment

Let's compare three implementations of a bird sighting website

# [birds 1](../demos/make-birds1.bash "this is a static website")

CSV file, Pandoc, 2 directories, 5 files, **72 total line count**, static site hosting. Included is a typical setup I would use for a static site project.

Lines   Files
------  ---------------
    21  [README.md](../demos/birds1/README.html)
     4  [birds.csv](../demos/birds1/birds.csv)
     6  [build.sh](../demos/birds1/build.sh)
     7  [page.tmpl](../demos/birds1/page.tmpl)
    32  htdocs/index.html

See <https://caltechlibrary.github.io/newt/demos/birds1/htdocs/index.html>

# [birds 2](../demos/make-birds2.bash "this website requires a machine")

SQL (Postgres + PostgREST), Browser JavaScript, 2 directories, 8 files, **232 total line count**, dynamic site requires hosting

Lines    Files
------   --------------
    29   [README.md](../demos/birds2/README.html)
     4   [birds.csv](../demos/birds2/birds.csv)
    34   [setup.sql](../demos/birds2/setup.sql)
    60   [models.sql](../demos/birds2/models.sql)
    15   [models_test.sql](../demos/birds2/models_test.sql)
     3   [postgrest.conf](../demos/birds2/postgrest.conf)
    24   htdocs/[index.html](../demos/birds2/htdocs/index.html)
    63   htdocs/[sightings.js](../demos/birds2/htdocs/sightings.js)

# [birds 3](../demos/make-birds3.bash "this website requires a machine")

SQL (Postgres + PostgREST), Pandoc, Newt, 1 directory, 7 files, 277 total line count, dynamic site requires hosting, **no JavaScript required**

Lines    Files
------   ---------------
    43   [README.md](../demos/birds3/README.html)
     4   [birds.csv](../demos/birds2/birds.csv)
    34   [setup.sql](../demos/birds2/setup.sql)
    60   [models.sql](../demos/birds2/models.sql)
    15   [models_test.sql](../demos/birds2/models_test.sql)
     3   [postgrest.conf](../demos/birds2/postgrest.conf)
    25   [birds.yaml](../demos/birds3/birds.yaml)
    36   [page.tmpl](../demos/birds3/page.tmpl)
     7   [post_result.tmpl](../demos/birds3/post_result.tmpl)


# A microservice conversation

1. web browser => Newt
2. Newt => PostgREST
3. Newt => Pandoc
4. Newt => web browser

# Newt's YAML file

- htdocs directory for static concetnt (option)
- environment variable used to access JSON sources (optional)
- route definitions
  - (optional) variable definitions (path and form data)
  - request routing details (e.g. path, method)
  - JSON source details (e.g. API URL, method, content type)
  - (optional) Pandoc template filename

See <https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/birds3/birds.yaml>

# Developer workflow

1. Model data in Postgres
2. Create/update Pandoc templates
3. Create/update routes/form data validation in YAML file
4. (Re)start PostgREST and Newt to (re)load models and routes
5. Test with our web browser

**Repeat as needed**

# Recap

version   site type   pros                     cons
-------   ---------   -----------------------  ----------------------------
birds 1   static      easy to conceptualize,   read only
                      no JavaScript required
birds 2   dynamic     read/write data          requires SQL knowledge
                                               requires browser JavaScript
                                               JavaScript is complex
birds 3   dynamic     read/write data          requires SQL knowledge
                      easy to conceptualize    requires knowledge of Pandoc
                      no JavaScript required   requires knowledge of YAML

# Newt stack weakness

- Newt is **an experimental prototype** (June/July 2023)
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
- SQL turns some people off, extending Newt YAML could mitigate that 
- Pandoc templates are simple to learn
- Newt stack does not enhibit all the HTML5 goodness available for front-end- Moving from JavaScript on the front end to Pandoc templates and birds.yaml for the back in was a wash interms of lines of code

# An unexpected result

> Newt stack should scale big!

- Newt can be scaled wide
- Pandoc server can be scaled wide
- PostgREST can be scale wide
- Postgres can be clustered

# Next steps for Newt?

1. Testing Solr/Elasticsearch as alternate JSON sources
2. I am building staff facing applications this Summer (2023)
3. I hope to move beyond my proof of concept in Winter (2023), Spring (2024)

# Someday, maybe ...

- Have Newt delegate file uploads to an S3 like service
- Enhance the Newt's YAML to generate SQL models
- Explore integrating SQLite3 support

# If Newt had a community ...

- share SQL code
- share Pandoc templates
- share YAML files
- improve Newt

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

