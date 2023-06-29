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
updateDate: 2023-06-28
pubDate: 2023-07-14
place: UCLA
date: July 14, 2023
section-titles: false
toc: true
keywords: [ "code4lib", "microservice", "Postgres", "PostgREST", "Pandoc" ]
url: "https://caltechlibrary.github.io/newt/presentation"
---

# The experiment

- Web applications for Libraries, Archives and Museums
  1. using a scale small approach
  2. largely from off the shelf parts
  3. composition through configuration

# Three abstractions, three cognative shifts

- A Template engine => [Pandoc](https://pandoc.org)
- A JSON source to manage data => [Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org)
- A data router, form validator => [Newt](https://github.com/caltechlibrary/newt/)

- Use templates to generate HTML (Pandoc)
- Model data in SQL (Postgres) get JSON API (PostgREST)
- Orchestrate our microservice conversation with YAML (Newt)

# A microservice conversation

1. web browser => Newt
2. Newt => Postgres+PostgREST
3. Newt => Pandoc
4. Newt => web browser

# Why model data with SQL?

> PostgreSQL+PostgREST, my code savings plan

- Good at describing structured data
- Supports queries, views, functions, triggers, events, ...
- Allows us to model our data once in one language

> **Minimize the source Luke!**

# My scale small experiment

Let's compare three implementations of a bird sighting website

- simple model
- bird, place and date sighted

# [birds 1](../demos/make-birds1.bash "this is a static website")

CSV file, Pandoc, 2 directories, 5 files, **75 total line count**, static site
(I've included my typical project setup with a README)

Lines   Files
------  ---------------
    26  README.md
     4  birds.csv <-- this is used in each of the demos
     6  build.sh
     7  page.tmpl
    32  htdocs/index.html

# [birds 2](../demos/make-birds2.bash "this website requires a machine")

SQL (Postgres + PostgREST), Browser JavaScript, 2 directories, 8 files, **232 total line count**, dynamic site

Lines    Files
------   --------------
    29   README.md
     4   birds.csv <-- from birds1
    34   setup.sql
    60   models.sql
    15   models_test.sql
     3   postgrest.conf
    24   htdocs/index.html
    63   htdocs/sightings.js

# [birds 3](../demos/make-birds3.bash "this website requires a machine")

SQL (Postgres + PostgREST), Pandoc, Newt, 1 directory, 7 files, **225 total line count**, dynamic site

Lines    Files
------   ---------------
    43   README.md
     4   birds.csv <-- from birds1
    34   setup.sql <-- from birds2
    60   models.sql <-- from birds2
    15   models_test.sql <-- from birds2
     3   postgrest.conf <-- from birds2
    23   birds.yaml
    36   page.tmpl
     7   post_result.tmpl
# Insights from experiment

- "Off the shelf" microservices can make application construction easier
- Orchestrating the data pipeline in YAML seems reasonable
- SQL turns some people off
  - models could be bootstraped from Newt's YAML using the form/path validation DSL
- Pandoc templates are simple to learn, should include examples
- Newt stack plays well with HTML5 approaches and best practices
- I realized one unexpected result ...

# The unexpected result

- Newt does not maintain "syncronous state"
- PostgREST maintains very little "syncronous state"
- Pandoc server maintains no "syncronous state"
- Postgres can be clustered

> The Newt stack can scale big

# Newt has weaknesses

- Newt is **an experimental prototype** (June/July 2023, six weeks old)
- Newt doesn't support file uploads

# Newt stack has strengths

> A very mature foundation

- 20th Century tech
  - SQL (1974), HTTP (1991), HTML (1993), Postgres (1996)
- 21st Century tech
  - JSON (2001), YAML (2001), Pandoc (2006), PostgREST (2014)

# Next steps for Newt?

1. Explore SQL generation from Newt's YAML
2. Test with Solr/Elasticsearch as alternate JSON sources
3. Build staff facing applications this Summer (2023)
4. Hopefully move beyond my proof of concept in Winter (2023), Spring (2024)

# Newt's someday, maybe ...

- Have Newt delegate file uploads to an S3 like service (minio via a file stream?)
- Explore integrating SQLite3 support
- A Newt community to share YAML, SQL and Pandoc templates

# Related resources

- Newt <https://github.com/caltechlibrary/newt>
- Postgres <https://postgres.org> + PostgREST <https://postgrest.org>
  - PostgREST Community Tutorials <https://postgrest.org/en/stable/ecosystem.html>
- Compiling Pandoc or PostgREST requires Haskell
  - Install Haskell GHCup <https://www.haskell.org/ghcup/>

# Thank you!

- This Presentation <https://caltechlibrary.github.io/newt/presentation/>
- Project: <https://github.com/caltechlibrary/newt>
- Email: rsdoiel@caltech.edu
