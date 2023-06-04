---
title: "Newt, assemble web applications with Pandoc, Postgres and PostgREST"
author: "R. S. Doiel, <rsdoiel@caltech.edu>"
institute: |
  Caltech Library,
  Digital Library Development
description: Code4Lib Meet up, Los Angeles
fronttheme: "default"
fontsize: 12pt
urlcolor: blue
linkstyle: bold
aspectratio: 169
createDate: 2023-05-16
pubDate: 2023-07-14
place: UCLA
date: July 14, 2023
section-titles: false
toc: true
keywords: [ "code4lib", "meetup", "Postgres", "webstack", "PostgREST", "Pandoc" ]
url: "https://caltechlibrary.github.io/newt/presentation"
---


# Today, LAMP and its legacy 

Four example systems found in Caltech Library

- EPrints
- ArchivesSpace
- Islandora
- Invenio RDM

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
                                 NodeJS and NPM
Islandora      PHP/SQL           MySQL, Fedora, Apache 2


These are all really complicated pieces of software.

# The problem

Each listed application is built on a stack. The stacks are complex. Because of the complexity it's hard to sustain them. Some we've outsourced to SAAS providers (e.g. ArchivesSpace). Some we treat as a back boxes (e.g. EPrints). Some we continue to run (e.g. Invenio RDM). 

It's just not fun supporting applications at this level of complexity. It takes too much time and energy. It detracts from delivering useful things to our Library, Archives and Caltech Community.

# Why are these things so complex?

> WARNING: gross generalizations ahead

- We want more from our application so more code gets written
- Complexity accrues over time
- A Silicon Valley influenced "best practice"
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
- Scaling big is the cloud's siren song

# The alternative, **scale small**

- Pack only what you need
- Simplify! 

# Scaling small

- Limit the moving parts
- Limit the cognitive shifts
- Minimize the toolbox while maximizing how you use it
- Write less code but remaining readable!

# Building small

- We need to
    - pick the right abstractions
    - pick the right division of responsibilities
    - solve the problem in the desired scope
    - Avoid obfuscation and duplication 
    - Avoid magical capabilities
    - Write less code, read less code

# Why?

> Human time in is a scarce resource

# How minimal can we go?

- Use "off the self" microservices
- SQL 
- Pandoc

# Can we create applications using only SQL and Pandoc?

Here's the "off the shelf" microservices I am experimenting with

- [Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org)
- [Pandoc](https://pandoc.org)
- [Newt](https://github.com/caltechlibrary/newt/)

# Simplify through a clear division of labor

- [Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org) => JSON data API, i.e. manages the data
- [Pandoc](https://pandoc.org) =>  a powerful template engine
- [Newt](https://github.com/caltechlibrary/newt/) => provides request router, response assembly and static file services

# How does this work?

Think of a game of telephone

> web browser => Newt => PostgREST => Pandoc => web browser


# How can this work?

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

# Server side knowledge requirements

- SQL
- Pandoc templates
- A CSV file orchestrating our microservices

# Client side knowledge requirements

- HTML
- CSS (optional)
- JavaScript (optional)

# What does this enable?

We can create interactive applications with

- SQL
- Pandoc templates 
- A little routing information

# Why SQL?

- SQL is really good at describing structured data
- SQL also is good at expressing queries
- A deeper knowledge of SQL provides you with
    - data views, functions, procedures
- SQL has rich data types, e.g. JSON columns
- Postgres + PostgREST provides a full featured JSON data API

# PostgreSQL+PostgREST, a code savings plan

> Minimize the source Luke!

- You don't need to learn an ORM and aren't limited by one
- You don't duplicate the SQL models in another language
    - e.g. classes in Python, PHP, Ruby or Java
- You don't write middleware to get an API

# Three cognitive shifts

- Write SQL to generate JSON
- Use Pandoc to transform JSON to HTML (or other formats)
- Use a CSV file to orchestrate

# Three common data flows

web browser => Newt => PostgREST => Pandoc => web browser

web browser => Newt => PostgREST => web browser

web browser => Newt => static file system => web browser

# Is this really simpler?

Let's take a look at three versions of a bird sighting web site.

- [birds 1](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/setup-birds1.bash), a static site implementation
- [birds 2](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/setup-birds2.bash), a dynamic site implementation, content viewing requires browser JavaScript
- [birds 3](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/setup-birds3.bash), a dynamic site implementation, does not require browser JavaScript

# Different birds 1

## [birds 1](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/setup-birds1.bash) static site (read only)

- Built with Pandoc from Markdown and CSV file
- Adds bird sightings via updating a CSV file and rebuilding site with Pandoc

~~~
       5 birds1/README.md
       4 birds1/birds.csv
       3 birds1/build.sh
      13 birds1/page.tmpl
      25 total
~~~

# Different birds 2

## [birds 2](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/setup-birds2.bash), dynamic site (read/write)

- Built with SQL using Postgres + PostgREST
- Requires the web browser to assemble pages via API calls
- Add birds using a web form requiring JavaScript
- JavaScript has become complex
    - handles fetching data and inserting it into the page
    - handles form prep and submission of our web form
- Solution doesn't work for text only web browsers like Lynx

~~~
      32 birds2/README.md
       4 birds2/birds.csv
       3 birds2/postgrest.conf
      50 birds2/setup.sql
      24 birds2/htdocs/index.html
      63 birds2/htdocs/sightings.js
     176 total
~~~

# Different birds 3

## [birds 3](https://raw.githubusercontent.com/caltechlibrary/newt/main/demos/setup-birds3.bash), dynamic site (read/write)

- Build from SQL (Postgres + PostgREST) and Pandoc
- Add birds using a simple web form, **no JavaScript**
- Everything is rendered server side
- Works even for text web browsers like Lynx

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

# Our different birds

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

- Complicated activities are handled by "off the self" microservices
- Main complexity is limited to SQL and our model data
- Leverages our Pandoc knowledge
- Avoids browser side page assembly

# Newt manages data flow

- request => JSON data API => Pandoc => response
- provides a simple DSL for mapping requests to API and Pandoc
- the data flow, or route, can be managed in spreadsheet!

# Developer workflow

1. Model data in Postgres
2. Create/update Pandoc templates
3. Create/update routes CSV file in a spreadsheet
4. (Re)start PostgREST and Newt to (re)load models and routes

**Repeat as needed**

# Our approach trys to minimize newness

- If you've attended a data science workshop you likely know enough SQL
- If you've built a static website you likely know about Pandoc
- I think there is community knowledge of CSV files
- A simple DSL for mapping requests to data sources and Pandoc
- SQL + CSV files + Pandoc => web application

=> Is this useful for Libraries, Archives and Museums?

# Approach weeknesses

- Newt is **an experimental prototype** (June 2023)
- Newt doesn't validate the POST, PATCH or PUT data
- Newt doesn't support file uploads
- Postgres and SQL have a learning curve
- Pandoc has a learning curve
- Using the HTTP protocol has a learning curve

# Approach strengths

- Maturity and communities
  - SQL (1974)
  - HTTP (1991)
  - HTML (1993)
  - Postgres (July 1996)
  - JSON (April 2001)
  - Pandoc (August 2006)
  - PostgREST (June 2014)

# Next steps for Newt?

- I am building staff facing applications, Summer 2023
- Planning to test with Solr/Opensearch as a JSON data source
- Fix bugs, simplify code, add validation, improve performance

## Someday, maybe

- It would be nice if ....
    - Newt validated POST, PUT and PATCH before sending to API
    - Newt could delegate file uploads to an S3 like service
    - Had it's own community supporting it

# Thank you!

- Presentation <https://caltechlibrary.github.io/newt/presentation/>
- Project: <https://github.com/caltechlibrary/newt>
- Email: rsdoiel@caltech.edu

