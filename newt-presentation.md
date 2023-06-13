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
updateDate: 2023-06-09
pubDate: 2023-07-14
place: UCLA
date: July 14, 2023
section-titles: false
toc: true
keywords: [ "code4lib", "meetup", "Postgres", "webstack", "PostgREST", "Pandoc" ]
url: "https://caltechlibrary.github.io/newt/presentation"
---

# Newt

> an experiment in simplifying the webstack 

... but first some context

# LAMP and its legacy 

Complex systems found in Caltech Library

- ArchivesSpace
- EPrints
- Invenio RDM
- Islandora

# Required Knowledge

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


These are all really complicated pieces of software.

# The problem

1. Each application was built on a stack
2. The stacks are complex and divergent
3. Coping strategies
	a. Out source running the applications as a "SAAS"
	b. Treat app as blackbox and avoiding customization
    c. Customize and dedicate individuals to the specific app
    d. ...

# The problem

None of our coping strategies is good.

# Why are these things so complex?

- We want more from our application, more code gets written
- Application needs enhancements, complexity accrues overtime
- Best practices like "systems should be designed to scale"

# Let's talk about scale

scale
: a euphemism for **scaling big**

# Scaling big

- Scaling big is hard
- Scaling big makes things really complex
- Scaling big favors large teams

# Scaling big, a reflection

- scaling big => 
  - distributed application design
  - programmable infrastructure 
  - cache systems and dynamic clustering
  - complex systems management

# An alternative, **scale small**

- Pack only what is needed
- Simplify! 

# Scaling small

- Limit the moving parts
- Limit the cognitive shifts
- Minimize the toolbox, maximize how you use it
- Write less code

# Newt, a small system experiment

- Use "off the shelf" microservices
- SQL ([Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org))
- [Pandoc templates](https://pandoc.org/MANUAL.html#templates) ([Pandoc Server](https://pandoc.org/pandoc-server.html))
- YAML configuration ([Newt](https://github.com/caltechlibrary/newt/))

# Simplify through a clear division of labor

- [Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org) => JSON data API, i.e. manages the data
- [Pandoc](https://pandoc.org) =>  a powerful template engine
- [Newt](https://github.com/caltechlibrary/newt/) => data router, form data validator, static file services

# How does this work?

1. web browser => Newt
2. Newt => PostgREST
3. Newt => Pandoc 
4. Newt => web browser

# How does this work?

1. Model our data using SQL (Postgres)
2. Run our JSON data API (Postgres+PostgREST)
3. Transform our structured data using Pandoc templates
4. Newt's YAML file orchestrates our microservice conversation

# How does Newt orchestrate our microservices?

- Newt's YAML file includes descriptions of data routing
    - (optional) route and form variable definitions
    - request routing details
    - how to form the request to the data source API
    - (optional) Pandoc template to process JSON results

# Our minimized Toolbox

- Text editor
- Web browser
- Pandoc
- Postgres + PostgREST
- Newt 

# Client side knowledge

- HTML
- CSS
- JavaScript

> the usual suspects

# Server side knowledge

- SQL
- Pandoc templates
- Our YAML orchestration file

# Why emphasize SQL?

- SQL is good at describing structured data
- SQL is good at expressing queries
- SQL has rich data types, e.g. JSON columns
- SQL can define data views, functions, procedures

# PostgreSQL+PostgREST, a code savings plan

> Minimize the source Luke!

- Don't need to learn an ORM, aren't limited by it
- Don't duplicate the SQL models in another language
- Don't write any middle-ware

# Three cognitive shifts

- Write SQL (Postgres), use JSON (PostgREST)
- Write Pandoc templates to transform JSON to HTML
- Write a YAML file to orchestrate our microservices

# Is this really simpler?

Three versions of a bird sighting website

- [birds 1](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/setup-birds1.bash), a static site implementation
- [birds 2](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/setup-birds2.bash), a dynamic site implementation, content viewing requires browser JavaScript
- [birds 3](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/setup-birds3.bash), a dynamic site implementation, does not require browser JavaScript

# [birds 1](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/setup-birds1.bash) 

- Built with Pandoc from Markdown and YAML file

~~~
2 directories, 5 files
       5 birds1/README.md
       4 birds1/birds.csv
       3 birds1/build.sh
      13 birds1/page.tmpl
      38 birds1/htdocs/index.html
      63 total
~~~

# [birds 2](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/setup-birds2.bash)

- Built with SQL (Postgres + PostgREST), Browser side JavaScript

~~~
2 directories, 6 files
      32 birds2/README.md
       4 birds2/birds.csv
       3 birds2/postgrest.conf
      50 birds2/setup.sql
      24 birds2/htdocs/index.html
      63 birds2/htdocs/sightings.js
     176 total
~~~

# [birds 3](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/setup-birds3.bash)

- Built with SQL (Postgres + PostgREST) and Pandoc
- **no JavaScript**

~~~
2 directories, 7 files
      33 birds3/README.md
       4 birds3/birds.csv
      24 birds3/birds.yaml
      40 birds3/page.tmpl
       7 birds3/post_result.tmpl
       3 birds3/postgrest.conf
      50 birds3/setup.sql
     161 total
~~~

# Three birds

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


# Birds 3 => Postgres+PostgREST, Pandoc and Newt

- Our "off the shelf" microservices limit complexity
- SQL defines data model and API end points
- Pandoc templates transform JSON to HTML
- A YAML file describes the routes and form data validation

# Newt manages data flow

- request => JSON data API => Pandoc => response
- provides a simple DSL for mapping requests to API and Pandoc
- leverage DSL for describing form data validation

# Developer workflow

1. Model data in Postgres
2. Create/update Pandoc templates
3. Create/update routes/form data validation in YAML file
4. (Re)start PostgREST and Newt to (re)load models and routes

**Repeat as needed**

# Minimal new knowledge

- If you've attended a data science workshop you likely know enough SQL
- If you've built a static website you likely know about Pandoc
- If you've built a static website you likely know YAML
- Use YAML and a simple DSL to map requests to data sources and form data validation
- SQL + YAML files + Pandoc => web application

# Weaknesses

- Newt is **an experimental prototype** (June 2023)
- Newt doesn't support file uploads
- Learning curves: 
    1. Postgres and SQL
    2. Pandoc
    3. using HTTP methods
    4. YAML

# There is strength in Maturity

- SQL (1974)
- HTTP (1991)
- HTML (1993)
- Postgres (1996)
- JSON (2001)
- YAML (2001)
- Pandoc (2006)
- PostgREST (2014)

# Next steps for Newt?

- I am building staff facing applications, Summer 2023
- Testing Solr/Elasticsearch as a JSON data source
- Fix bugs, expand types/validation, smoothing rough edges

# It would be nice if ...

- Newt could delegate file uploads to an S3 like service
- Newt had a community 
	- sharing SQL code
	- sharing Pandoc templates
    - sharing YAML files

# Someday, maybe ideas

- Rewrite Newt in Haskell
    - Integrate Pandoc into Newt, skipping running Pandoc server
- Use YAML file to extrapolate SQL, templates and website 

# An unexpected result of simplification

- Newt can be scaled wide (parallel)
- Pandoc server can be scaled wide
- PostgREST can be scale wide
- Postgres (the only part holding state) can be clustered

# Thank you!

- This Presentation <https://caltechlibrary.github.io/newt/presentation/>
- Project: <https://github.com/caltechlibrary/newt>
- Email: rsdoiel@caltech.edu

