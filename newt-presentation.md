---
title: "Newt, assemble web applications with Pandoc, Postgres and PostgREST"
author: "R. S. Doiel, <rsdoiel@caltech.edu>"
institute: |
  Caltech Library,
  Digital Library Development
description: Code4Lib Meet up, Los Angeles
urlcolor: blue
linkstyle: bold
aspectratio: 169
createDate: 2023-05-16
updateDate: 2023-06-14
pubDate: 2023-07-14
place: UCLA
date: July 14, 2023
section-titles: false
toc: true
keywords: [ "code4lib", "meetup", "Postgres", "webstack", "PostgREST", "Pandoc" ]
url: "https://caltechlibrary.github.io/newt/presentation"
---

# Newt

## a small system experiment

### ... but first some context

# LAMP and stacks of complexity

Software           Stack            
-----------------  ---------------------
EPrints            Clasic LAMP + Xapian
Islandora          Apache+MySQL+PHP 
                   Fedora and Solr
ArchivesSpace      Apache, MySQL, 
                   Tomcat,
                   Java+jRuby, Solr
Invenio RDM        NginX, Python/pip,
                   Postgres, Redis,
                   Elasticsearch,
                   Docker, NodeJS/NPM

# Our Legacy of complexity

1. Applications are built on a stack
2. Stacks are complex and divergent
3. Sustaining them requires many coping strategies

# The problem

- Our coping strategies are not sustainable
- Complexity is part of the problem

# Why are these things so complex?

1. We want more from our application, more code gets written
2. We want "enhancements", complexity accrues overtime
3. We build "systems designed to scale"

# Scaling

scale
: a euphemism for **scaling big**, usage "google scale", "amazon scale"

# Scaling

- Scaling big is hard
- Scaling big can make things really complex
- Scaling big, a path to scaling small
  - **Can we pack only what is needed?**

# Scaling strategies

- distributed application design
- containers
- programmable infrastructure

# Scaling small

- Limit the moving parts
- Limit the cognitive shifts
- Try to **Write less code**

# Limit the moving parts, three abstractions

- [Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org) => JSON source to manage data
- [Pandoc](https://pandoc.org) =>  a powerful template engine
- [Newt](https://github.com/caltechlibrary/newt/) => data router, form data validator

# A microservice conversation

1. web browser => Newt
2. Newt => PostgREST
3. Newt => Pandoc
4. Newt => web browser

# Limit the cognitive shifts

- Write SQL (Postgres), get JSON source (PostgREST)
- Write Pandoc templates, use JSON, get HTML
- Write YAML, orchestrate our microservice conversation

# Writing less code

1. Use "off the shelf" microservices
2. Take advantage of SQL
3. Standardize templating with Pandoc

> Postgres+PostgREST+Pandoc+Newt =>
> No need to write any middle-ware

# Why SQL?

- SQL is good at describing structured data
- SQL provides views, functions, triggers, ...
- SQL allows us to model our data once

> **Minimize the source Luke!**
> PostgreSQL+PostgREST is a code savings plan.

# A scale small experiment

Let's compare three implementations of a bird sighting website

# [birds 1](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/make-birds1.bash "this is a static website")

- Built with Pandoc from a CSV file and a Pandoc template
- 2 directories, 5 files, 63 total line count

Lines   Files
------  ---------------
     5  [README.md](https://caltechlibrary.github.io/newt/demos/birds1/README.md)
     4  [birds.csv](https://caltechlibrary.github.io/newt/demos/birds1/birds.csv)
     3  [build.sh](https://caltechlibrary.github.io/newt/demos/birds1/build.sh)
    13  [page.tmpl](https://caltechlibrary.github.io/newt/demos/birds1/page.tmpl)
    38  htdocs/index.html

See <https://caltechlibrary.github.io/newt/demos/birds1/htdocs/index.html>

# [birds 2](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/make-birds2.bash "this website requires a machine")

- Built with SQL (Postgres + PostgREST), Browser side JavaScript
- 2 directories, 6 files, 177 total line count

Lines    Files
------   --------------
    33   [README.md](https://caltechlibrary.github.io/newt/demos/birds2/README.md)
     4   [birds.csv](https://caltechlibrary.github.io/newt/demos/birds2/birds.csv)
    50   [setup.sql](https://caltechlibrary.github.io/newt/demos/birds2/setup.sql)
     3   [postgrest.conf](https://caltechlibrary.github.io/newt/demos/birds2/postgrest.conf)
    63   htdocs/[sightings.js](https://caltechlibrary.github.io/newt/demos/birds2/htdocs/sightings.js)
    24   htdocs/index.html

Dynamic site, requires hosting

# [birds 3](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/make-birds3.bash "this website requires a machine")

- Built with SQL (Postgres + PostgREST), Pandoc, Newt
- 1 directory, 7 files, 162 total line count

Lines   Files
------  ---------------
    33  [README.md](https://caltechlibrary.github.io/newt/demos/birds3/README.md)
     4  birds.csv
    50  setup.sql
    25  [birds.yaml](https://caltechlibrary.github.io/newt/demos/birds3/birds.yaml)
    40  [page.tmpl](https://caltechlibrary.github.io/newt/demos/birds3/page.tmpl)
     7  [post_result.tmpl](https://caltechlibrary.github.io/newt/demos/birds3/post_result.tmpl)
     3  postgrest.conf

Dynamic site, requires hosting, **no JavaScript**

# Newt's YAML file

- environment variable used to access JSON sources
- route definitions
  - (optional) variable definitions (path and form data)
  - request routing details (e.g. path, method)
  - JSON source details (e.g. api URL, method, content type)
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

# Weaknesses in my proposal

- Newt is **an experimental prototype** (June/July 2023)
- Newt doesn't support file uploads
- Newt doesn't eliminate learning curves, e.g. Postgres and SQL; Pandoc; using HTTP methods; YAML

# Newt's greatest advantage

> A mature foundation

- 20th Century tech
  - SQL (1974), HTTP (1991), HTML (1993), Postgres (1996)
- 21st Century tech
  - JSON (2001), YAML (2001), Pandoc (2006), PostgREST (2014)

# An unexpected result

> Newt can potentially scale big!

- Newt can be scaled wide (not stateful)
- Pandoc server can be scaled wide (not stateful)
- PostgREST can be scale wide (not stateful)
- Postgres can be clustered (stateful)

# Next steps for Newt?

- Testing Solr/Elasticsearch as alternate JSON sources
- I am building staff facing applications this Summer (2023)
- I hope to move beyond my proof of concept in Fall/Winter (2023)

# Sameday, maybe ...

- Have Newt delegate file uploads to an S3 like service
- Enhance the Newt's YAML to generate SQL models
- Explore integrating SQLite3 support

# If Newt had a community ...

- share SQL code
- share Pandoc templates
- share YAML files
- improve Newt

# Conclusions

- "Off the shelf" microservices can make application construction easier
- Orchestrating the data pipeline in YAML seems reasonable
- SQL may turn some people off
- Pandoc templates are simple to learn and avoid embedding business logic
- You still have all the HTML5 goodness available in the front-end

# Additional resources 

- [Postgres](https://postgres.org) + [PostgREST](https://postgrest.org)
- [Newt](https://github.com/caltechlibrary/newt)
- [Solr](https://solr.apache.org/), [Opensearch](https://www.opensearch.org/)
- Compiling Pandoc or PostgREST requires Haskell
  - Install Haskell [ghcup](https://https://www.haskell.org/ghcup/)

# Thank you!

- This Presentation <https://caltechlibrary.github.io/newt/presentation/>
- Project: <https://github.com/caltechlibrary/newt>
- Email: rsdoiel@caltech.edu

