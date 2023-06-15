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

# LAMP and its legacy 

Complex systems used by Caltech Library

- ArchivesSpace
- EPrints
- Invenio RDM
- Islandora

# LAMP and its legacy 

App            Languages         Supporting services
---------      ---------         -------------------
ArchivesSpace  Java, Ruby, SQL   MySQL, Solr, Apache or NginX
EPrints        Perl, SQL, XML,   MySQL, Apache2 (tight integration),
               EPrints XML       and Xapian
Invenio RDM    Python, SQL       Postgres, Redis, Elasticsearch,
               JavaScript        Docker, Invenio Framework,
                                 Python packaging system,
                                 React JavaScript framework,
                                 NodeJS and NPM, NginX
Islandora      PHP/SQL           MySQL, Drupal, Fedora, Apache 2


# The problem

1. Each application was built on a stack
2. The stacks are complex and divergent
3. Sustaining them requires many coping strategies

# The problem

- Our coping strategies are not sustainable
- Complexity is part of the problem

# Why are these things so complex?

- We want more from our application, more code gets written
- We want "enhancements", complexity accrues overtime
- Best Practices like "systems should be designed to scale"

# Why are these things so complex?

- first two are "people problems" (hard)
- last one might be a system design problem (solvable?)
 
# Scale (from computing practice)

scale
: a euphemism for **scaling big**, as used in phrases like "google scale", "amazon scale"

# Scaling big

- Scaling big is hard
- Scaling big can make things really complex
- Scaling big can require larger teams for success

# Scaling big

- What did scaling big deliver?
  - distributed application design
  - containers 
  - programmable infrastructure 
  - cache systems and dynamic clustering
  - complex systems management

# Scaling (from geometry)

scaling
: a linear transformation that enlarges or diminishes objects

# May I suggests? **scale small**

- Pack only what is needed
- Simplify! 

# Scaling small

- Limit the moving parts
- Limit the cognitive shifts
- Minimize the toolbox, maximize using it
- Try to **Write less code!**

# Limit the moving parts

> Simplify through a clear division of labor

- [Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org) => JSON API to manage data, it gives us a JSON source
- [Pandoc](https://pandoc.org) =>  a powerful template engine
- [Newt](https://github.com/caltechlibrary/newt/) => data router, form data validator, static file services

# Limit the moving parts

1. web browser => Newt
2. Newt => PostgREST
3. Newt => Pandoc 
4. Newt => web browser

# Limit the cognitive shifts, client side

> The status quo presists on the front-end

- HTML, CSS, JavaScript, image assets, etc.

# Limit the cognitive shifts, server side 

- SQL
- Pandoc templates
- Our YAML orchestration file

# Limit the cognitive shifts, why emphasize SQL?

- SQL is good at describing structured data
- SQL is good at expressing queries
- SQL provides rich data types
- SQL can define data views, functions, procedures, triggers, ...

# Write less code

> PostgreSQL+PostgREST is my code saving plan

- The SQL runtime is rich in capability
- Don't need to learn an ORM, we're not limited by an ORM
- Don't duplicate the SQL models in another language
- Don't write any middle-ware

> Minimize the source Luke!

# Minimized Toolbox, maximize using it

- Text editor
- Postgres + PostgREST
- Pandoc
- Newt 
- Web browser


# Three cognitive shifts

- Write SQL (Postgres) and get a JSON source (PostgREST)
- Write Pandoc templates to transform JSON to HTML
- Write a YAML file to orchestrate our microservice conversation

# Is this really simpler?

Three versions of a bird sighting website

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
   4 birds2/birds.csv
   3 birds2/postgrest.conf
  32 birds2/README.md
  50 birds2/setup.sql
  24 birds2/htdocs/index.html
  63 birds2/htdocs/sightings.js
 176 total
~~~

# [birds 3](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/make-birds3.bash)

- Built with SQL (Postgres + PostgREST), Pandoc, Newt
- **no JavaScript**

~~~
1 directory, 7 files
   4 birds3/birds.csv
  25 birds3/birds.yaml
  40 birds3/page.tmpl
   3 birds3/postgrest.conf
   7 birds3/post_result.tmpl
  33 birds3/README.md
  50 birds3/setup.sql
 162 total
~~~

# Comparing three birds

version    site type   pros                     cons
-------    ---------   -----------------------  ------------------------------
birds 1    static      easy to conceptualize,   read only
                       no JavaScript required
birds 2    dynamic     read/write data          requires SQL knowledge
                                                requires browser JavaScript
                                                JavaScript is complex
birds 3    dynamic     read/write data          requires SQL knowledge 
                       easy to conceptualize,   requires knowledge of Pandoc
                       no JavaScript required   requires knowledge of YAML 


# Bird 3 build features

1. Model our data using SQL (Postgres)
2. Access via JSON API (PostgREST)
3. Transform our structured data using Pandoc templates (Pandoc server)
4. Orchestrate our microservice conversation via YAML file (Newt)

# How does Newt orchestrate our microservices?

- Newt's YAML file includes descriptions for request routing
    - htdocs directory (if serving static files)
    - environment variable used to access JSON sources
    - route definitions
        - (optional) variable definitions
        - request routing details (e.g. path, method)
        - JSON source description (e.g. api URL, method, content type)
        - (optional) Pandoc template filename
    - port number (for Newt)

# Newt manages data flow

1. web browser => Newt
2. Newt => PostgREST
3. Newt => Pandoc 
4. Newt => web browser

# Developer workflow

1. Model data in Postgres (using psql and a text editor)
2. Create/update Pandoc templates (using a text editor and Pandoc)
3. Create/update routes/form data validation in YAML file (text editor, `newt -dry-run <YAML_FILENAME>`)
4. (Re)start PostgREST and Newt to (re)load models and routes
5. Test with our web browser

**Repeat as needed**

# Minimal prior knowledge

- A data science workshop often will teach you enough SQL
- Static website generation often envolves Pandoc
- Markdown documents often use YAML to represent front matter

# Weaknesses in my proposal

- Newt is **an experimental prototype** (June 2023)
- Newt doesn't support file uploads
- Newt doesn't eliminate Learning curves, e.g.
    1. Postgres and SQL
    2. Pandoc
    3. using HTTP methods
    4. YAML

# Newt's Advantages, a mature foundation

- 20th Century tech
  - SQL (1974), HTTP (1991), HTML (1993), Postgres (1996)
- 21st Century tech
  - JSON (2001), YAML (2001), Pandoc (2006), PostgREST (2014)


# Next steps for Newt?

- I am building staff facing applications this Summer (2023)
- Testing Solr/Elasticsearch as a JSON sources
- I hope to move beyond my proof of concept in the Fall (2023)

# It would be nice if ...

- Newt could delegate file uploads to an S3 like service
- Newt had a community 
    - sharing SQL code
    - sharing Pandoc templates
    - sharing YAML files
    - Improving Newt

# Posibilities, a roadmap beyond proof of concept

- Enhance the Newt's YAML to generate SQL models and Pandoc templates 
- Explore integrating SQLite3 support
- Rewrite Newt in Haskell
    - Integrate Pandoc into Newt avoiding an HTTP call

# An unexpected result of simplification

> Newt can potentially scale really big!

- Newt can be scaled wide (parallel), it requires minimal state (only what's in the configuration file)
- Pandoc server can be scaled wide (it retains zero stat )
- PostgREST can be scale wide (a minimal configuration file)
- Postgres (the only part holding state) can be clustered

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

