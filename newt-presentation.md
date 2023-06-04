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
updateDate: 2023-06-04
pubDate: 2023-07-14
place: UCLA
date: July 14, 2023
section-titles: false
toc: true
keywords: [ "code4lib", "meetup", "Postgres", "webstack", "PostgREST", "Pandoc" ]
url: "https://caltechlibrary.github.io/newt/presentation"
---


# LAMP and its legacy 

Four example systems found in Caltech Library

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

1. Each listed application is built on a stack
2. The stacks are complex, divergent
3. Coping strategies
	a. SAAS
	b. blackbox
	c. avoid customization


# Why are these things so complex?

- We want more from our application so more code gets written
- Complexity accrues over time
- A "best practice"
  - **Systems should be designed to scale**

# Let's talk about scale

- "scale", a euphemism for **scaling big**
- scaling big => 
  - distributed application design
  - programmable infrastructure 
  - cache systems and dynamic clustering
  - complex systems management

# Scaling big, a reflection

- Scaling big is hard
- Scaling big makes things really complex
- Scaling big favors large teams
- Scaling big is a siren song

# An alternative, **scale small**

- Pack only what you need
- Simplify! 

# Scaling small

- Limit the moving parts
- Limit the cognitive shifts
- Minimize the toolbox, maximizing how you use it
- Write less code

# Building small

- Pick the right abstractions
- Write less code

# Why?

> Human time is a scarce resource

# How minimal can we go?

- Use "off the shelf" microservices
- SQL 
- Pandoc


# "off the shelf" microservices experiment

- [Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org)
- [Pandoc](https://pandoc.org)
- [Newt](https://github.com/caltechlibrary/newt/)

# Simplify through a clear division of labor

- [Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org) => JSON data API, i.e. manages the data
- [Pandoc](https://pandoc.org) =>  a powerful template engine
- [Newt](https://github.com/caltechlibrary/newt/) => data router, static file services

# How does this work?

1. web browser => Newt
2. Newt => PostgREST
3. Newt => Pandoc 
4. Newt => web browser

# What does this enable?

1. Model our data using SQL (Postgres)
2. Define our JSON data API using SQL (Postgres+PostgREST)
3. Transform our structured data using Pandoc
4. Use Newt to orchestrate the microservice conversation

# Required Toolbox

- Text editor
- Spreadsheet (optional)
- Web browser
- Pandoc
- Postgres + PostgREST
- Newt 

# Client side knowledge requirements

- HTML
- CSS (optional)
- JavaScript (optional)

# Server side knowledge requirements

- SQL
- Pandoc templates
- A CSV file orchestrating our microservices

# Why SQL?

- SQL is good at describing structured data
- SQL is good at expressing queries
- SQL has rich data types, e.g. JSON columns
- SQL has data views, functions, procedures

# PostgreSQL+PostgREST, a code savings plan

> Minimize the source Luke!

- You don't need to learn an ORM, aren't limited by it
- You don't duplicate the SQL models in another language
- You don't write middleware to get a data API

# Three cognitive shifts

- Write SQL and use JSON
- Use Pandoc to transform JSON to HTML
- Use a CSV file to orchestrate our microservices

# Is this really simpler?

Three versions of a bird sighting website

- [birds 1](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/setup-birds1.bash), a static site implementation
- [birds 2](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/setup-birds2.bash), a dynamic site implementation, content viewing requires browser JavaScript
- [birds 3](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/setup-birds3.bash), a dynamic site implementation, does not require browser JavaScript

# [birds 1](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/setup-birds1.bash) 

- Built with Pandoc from Markdown and CSV file

~~~
       5 birds1/README.md
       4 birds1/birds.csv
       3 birds1/build.sh
      13 birds1/page.tmpl
      25 total
~~~

# [birds 2](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/setup-birds2.bash)

- Built with SQL (Postgres + PostgREST), Browser side JavaScript

~~~
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
      34 birds3/README.md
       4 birds3/birds-routes.csv
       4 birds3/birds.csv
       2 birds3/birds.yaml
      40 birds3/page.tmpl
       3 birds3/postgrest.conf
       9 birds3/redirect.tmpl
      50 birds3/setup.sql
     146 total
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
                       no JavaScript required


# Birds 3 => Postgres+PostgREST, Pandoc and Newt

- Our "off the shelf" microservices limit complexity
- SQL defines data model and API end points
- Pandoc templates transform JSON to HTML

# Newt manages data flow

- request => JSON data API => Pandoc => response
- provides a simple DSL for mapping requests to API and Pandoc
- the data flow can be managed with a spreadsheet!

# Developer workflow

1. Model data in Postgres
2. Create/update Pandoc templates
3. Create/update routes in CSV file
4. (Re)start PostgREST and Newt to (re)load models and routes

**Repeat as needed**

# Minimal new knowledge

- If you've attended a data science workshop you likely know enough SQL
- If you've built a static website you likely know about Pandoc
- Use a simple DSL used to map requests to data sources and Pandoc
- SQL + CSV files + Pandoc => web application

# Weaknesses

- Newt is **an experimental prototype** (June 2023)
- Newt doesn't accept POST, PATCH, PUT encoded as JSON
- Newt doesn't validate the GET, POST, PATCH or PUT data
- Newt doesn't support file uploads
- Learning curves: Postgres and SQL, Pandoc, using HTTP methods

# Strength in Maturity

- SQL (1974)
- HTTP (1991)
- HTML (1993)
- Postgres (1996)
- JSON (2001)
- Pandoc (2006)
- PostgREST (2014)

# Next steps for Newt?

- I am building staff facing applications, Summer 2023
- Testing Solr/Elasticsearch as a JSON data source
- Fix bugs, improve validation, simplify code

# It would be nice if ....

- Newt validated POST, PUT and PATCH before sending to API
- Newt could delegate file uploads to an S3 like service
- Had it's own community supporting it
	- share SQL code
	- share pandoc templates

# Thank you!

- Presentation <https://caltechlibrary.github.io/newt/presentation/>
- Project: <https://github.com/caltechlibrary/newt>
- Email: rsdoiel@caltech.edu

