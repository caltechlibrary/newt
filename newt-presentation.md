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

Each listed application is built on a stack. The stacks are complex. Because of the complexity it's hard to sustain them. Some we've outsourced to SAAS providers (e.g. ArchivesSpace). Some we treat as a back boxes (e.g. EPrints). It's just not fun supporting applications at this level of complexity. It takes too much time and energy. It detracts from delivering useful things to our Library, Archives and Caltech Community.

# Why are these things so complex?

> WARNING: gross generalizations ahead

- We want more from our application so more code gets written
- Complexity accrues over time
- A Silicon Valley influenced "best practices"
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
- This is the cloud's siren song

# The alternative, **scale small**

- Pack only what you need
- Simplify! 

# Scaling small

- Limit the moving parts
- Limit the cognitive shifts
- Minimize the toolbox while maximizing how you use it
- Write less code! 
- Remain readable!

# Wait, readable?

- Avoid obfuscation
- Avoid magical knowledge or capabilities
- Less code can mean less to read
- Pay attention to
    - abstractions
    - responsibilities
    - scope

# How minimal can we go?

- Off the self microservices
- SQL 
- Pandoc

# Can we create applications using only SQL and Pandoc?

Just about. Here's the off the shelf microservices I am experimenting with

- [Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org)
- [Pandoc](https://pandoc.org)
- [Newt](https://github.com/caltechlibrary/newt/)

# Simplify through a clear division of labor

- [Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org) => JSON data API, i.e. manages the data
- [Pandoc](https://pandoc.org) =>  a powerful template engine
- [Newt](https://github.com/caltechlibrary/newt/) => provides request router, response assembly and static file services

# How does this work in practice?

Think of a game of telephone
: web browser => Newt => PostgREST => Pandoc => web browser


# How would this work in practice?

1. Model our data using SQL (Postgres)
2. Define our JSON data API using SQL (Postgres+PostgREST)
3. Transform our structured data using Pandoc
4. Use Newt to orchestrate

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
- Manage a CSV file describing data flowing through our microservices

# Web browser knowledge requirements

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
- You don't write middleware to get an API anymore

# Three cognitive shifts

- Write SQL to generate JSON
- Use Pandoc to transform JSON to HTML (or other formats)
- Use a CSV file to describe our data flows
    - maps a request to PostgREST and Pandoc (aka routes)

# Three common data flows

web browser => Newt => PostgREST => Pandoc => web browser

web browser => Newt => PostgREST => web browser

web browser => Newt => static file system => web browser

# Still Helpful to know

- HTML 5 related W3C technologies
  - HTML 5 markup
  - CSS
  - modern JavaScript
- How to integrate static file assets, e.g. html files, images
- Understand how HTTP works, including HTTP methods and Headers

> the front-end can be as simple or as complex as you like

# Is this really simpler?

Let's take a look at three versions of a bird sighting web site.

- [birds 1](birds1/), a static site implementation
- [birds 2](birds2/), a dynamic site implementation, content viewing requires browser JavaScript
- [birds 3](birds3/), a dynamic site implementation, does not require browser JavaScript

# Different birds 1

## [birds 1](birds1/) static site (read only)

- Built with Pandoc from Markdown and CSV file
- Adds bird sightings via updating a CSV file and rebuilding site with Pandoc

# Different birds 2

## [birds 2](birds2/), dynamic site (read/write)

- Built with SQL using Postgres + PostgREST
- Requires the web browser to assemble pages via API calls
- Add birds using a web form requiring JavaScript
- JavaScript has become complex
    - handles fetching data and inserting it into the page
    - handles form prep and submission of our web form
- Solution doesn't work for text only web browsers like Lynx

# Different birds 3

## [birds 3](birds3/), dynamic site (read/write)

- Build from SQL (Postgres + PostgREST) and Pandoc
- Add birds using a simple web form, no JavaScript
- Rendered on server, no JavaScript
- Works even for text web browsers like Lynx

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

- The complicated activities are handled by off the self microservices
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

- Newt is **an experimental prototype** (May 2023)
- SQL has a learning curve
- Pandoc has a learning curve
- Using the HTTP protocol has a learning curve
- Newt doesn't (yet) support file upload handling

# Approach strengths

- We have a mature platform built from
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
- Fix bugs, simplify code, improve performance

## Someday, maybe

- It would be nice if ....
    - Newt routed file uploads to an S3 like service
    - had a better DSL to map requests
    - had a community supporting it

# Thank you!

- Presentation <https://caltechlibrary.github.io/newt/presentation/>
- Project: <https://github.com/caltechlibrary/newt>
- Email: rsdoiel@caltech.edu

