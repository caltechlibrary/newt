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

# LAMP, its legacy, complexity

App            Languages         Supporting services
---------      ---------         -------------------
ArchivesSpace  Java, Ruby, SQL   MySQL, Solr, Apache or NginX
               JavaScript, CSS
EPrints        Perl, SQL, XML,   MySQL, Apache2 (tight integration),
               EPrints XML,      and Xapian (full text search)
               JavaScript, CSS
Invenio RDM    Python, SQL       Postgres, Redis, Elasticsearch,
               JavaScript, CSS   Docker, Invenio Framework,
                                 Python/PIP, React Framework,
                                 NodeJS/NPM, NginX
Islandora      PHP/SQL           MySQL, Drupal, Fedora, Apache 2
Custom Apps    Python+Flask+ORM  Python/PIP, Flask, MySQL/SQLite2
               JavaScript, CSS   NodeJS/NPM

# Our Legacy of complexity

1. Application are built on a stack
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
: a euphemism for **scaling big**, as used in phrases like "google scale", "amazon scale"

# Scaling

- Scaling big is hard
- Scaling big can make things really complex
- Scaling big lead to scaling small?
  - **Pack only what is needed**

# Scaling approaches

- distributed application design
- containers
- programmable infrastructure
- cache systems and dynamic clustering
- complex systems management

# Scaling small

- Limit the moving parts
- Limit the cognitive shifts
- Try to **Write less code**

# Limit the moving parts

> Three abstractions

- [Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org) => JSON API to manage data, it gives us a JSON source
- [Pandoc](https://pandoc.org) =>  a powerful template engine
- [Newt](https://github.com/caltechlibrary/newt/) => data router, form data validator

# A Microservice conversation

1. web browser => Newt
2. Newt => PostgREST
3. Newt => Pandoc
4. Newt => web browser

# Limit the cognitive shifts

- Write SQL (Postgres), get JSON source (PostgREST)
- Write Pandoc templates, provide JSON, get HTML
- Write a YAML, orchestrate our microservice conversation

# Our Toolbox

- Text editor
- Postgres + PostgREST
- Pandoc
- Newt
- Web browser

# Writing less code

1. Use "off the shelf" microservices
2. Take advantage of SQL
3. Standardize templating with Pandoc

> Postgres+PostgREST+Pandoc+Newt =>
> No need to write any middle-ware

# Why SQL?

> **Minimize the source Luke!**

- SQL is good at describing structured data
- SQL provides views, functions, triggers, ...
- SQL allows us to model our data once

> PostgreSQL+PostgREST, a code savings plan.

# How did I arrive at Newt?

Let's compare three implementations of a bird sighting website

# [birds 1](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/make-birds1.bash) 

- Built with Pandoc from a CSV file and a Pandoc template

~~~
2 directories, 5 files
       5 birds1/README.md
       4 birds1/birds.csv
       3 birds1/build.sh
      13 birds1/page.tmpl
      38 birds1/htdocs/index.html
      63 total
~~~

# [birds 2](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/make-birds2.bash)

- Built with SQL (Postgres + PostgREST), Browser side JavaScript

~~~
2 directories, 6 files
  33 birds2/README.md
   4 birds2/birds.csv
  50 birds2/setup.sql
   3 birds2/postgrest.conf
  24 birds2/htdocs/index.html
  63 birds2/htdocs/sightings.js
 177 total
~~~

# [birds 3](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/make-birds3.bash)

- Built with SQL (Postgres + PostgREST), Pandoc, Newt
- **no JavaScript**

~~~
1 directory, 7 files
  33 birds3/README.md
   4 birds3/birds.csv
  50 birds3/setup.sql
  25 birds3/birds.yaml
  40 birds3/page.tmpl
   3 birds3/postgrest.conf
   7 birds3/post_result.tmpl
 162 total
~~~

# Newt's YAML file

- environment variable used to access JSON sources
- route definitions
  - (optional) variable definitions (path and form data)
  - request routing details (e.g. path, method)
  - JSON source details (e.g. api URL, method, content type)
  - (optional) Pandoc template filename

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

- Newt is **an experimental prototype** (June 2023)
- Newt doesn't support file uploads
- Newt doesn't eliminate Learning curves, e.g.
    1. Postgres and SQL
    2. Pandoc
    3. using HTTP methods
    4. YAML

# Newt's greatest advantage

> A mature foundation

- 20th Century tech
  - SQL (1974), HTTP (1991), HTML (1993), Postgres (1996)
- 21st Century tech
  - JSON (2001), YAML (2001), Pandoc (2006), PostgREST (2014)

# An unexpected result

> Newt can potentially scale big!

- Newt can be scaled wide, it requires minimal state
- Pandoc server can be scaled wide, it retains zero state
- PostgREST can be scale wide
- Postgres (the only part holding state) can be clustered

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
- Orchestrating the data pipeline in YAML seems 
- You need to pick the right ones for the task
- You need to only use what you need
- Having an "off the self" orchestrate like Newt eliminates for of the need to write middleware and frees up time to focus on humane user interfaces


# Additional resources 

- [Postgres](https://postgres.org)
- Install Haskell [ghcup](https://https://www.haskell.org/ghcup/)
  - Build latest [PostgREST](https://postgrest.org)
  - Build latest [Pandoc](https://pandoc.org)
- [Newt](https://github.com/caltechlibrary/newt) proof of concept
- [Solr](https://solr.apache.org/), [Opensearch](https://www.opensearch.org/)
- Alternatives to Postgres+PostgREST
  - [MariaDB](https://mariadb.com/products/enterprise/) + [Maria Max Scale](https://mariadb.com/kb/en/maxscale/)
  - [MySQL](https://dev.mysql.com) + [MRS](https://blogs.oracle.com/mysql/post/introducing-the-mysql-rest-service)

# Thank you!

- This Presentation <https://caltechlibrary.github.io/newt/presentation/>
- Project: <https://github.com/caltechlibrary/newt>
- Email: rsdoiel@caltech.edu

