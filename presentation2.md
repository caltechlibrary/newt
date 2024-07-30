---
title: "Newt, the second prototype"
author: "R. S. Doiel, <rsdoiel@caltech.edu>"
institute: |
  Caltech Library,
  Digital Library Development
description: Talk about Newt Project's second prototype
urlcolor: blue
linkstyle: bold
aspectratio: 169
createDate: 2024-02-14
updateDate: 2024-04-10
pubDate: 2024-04-19
place: SoCal Code4Lib Meetup at USC
date: 2024-04-19
section-titles: false
toc: true
keywords: [ "web service", "micro service", "Postgres", "PostgREST", "Mustache" ]
url: "https://caltechlibrary.github.io/newt/presentation2"
---

# What is Newt?

- A rapid application develop tool
    - for applications that curate metadata
- Audience: Libraries, Archives and Museums

# Goal of Prototype 2: Answer the question.

Is Newt and "off the shelf" software enough to create metadata curation applications?

# High level Concepts

- describe the application you want
- generate the application you described

# Implementation Concepts

- data sources
- data models
- routing requests through data pipelines
- rendering JSON responses via template engine

# Themes

- Pick Simple = (No coding) + (Less coding)
- Compose applications using data pipelines and templates
- Avoid inventing new things

# Off the shelf (no coding)

- [Postgres](https://postgresql.org) and [PostgREST](https://postgrest.org) 
- [Solr](https://solr.apache.org) or [OpenSearch](https://opensearch.org)
- Newt Mustache => Transform JSON into web pages
- Newt Router, ties it all together

# Office the shelf (other data sources)

- ArchivesSpace, RDM -> JSON API
- ORCID, ROR, CrossRef, DataCite -> JSON API

# Assemble app from YAML (less coding)

- The application you want is described in YAML
- Newt generates the code you need
- Customize by editing the generated code

# How are data models described?

- A model is a set of HTML form input types
- Expressed using GitHub YAML Issue Template Syntax
- Model describes HTML and implies SQL

# How do I think things will work?

1. Interactively generate our application's YAML file
2. Interactively define data models
3. Generate our application code
4. Setup Postgres and PostgREST
5. Run our app with Newt

# Steps one and two are interactive

~~~shell
  newt init app.yaml
  newt model app.yaml
~~~

# Step three, generate our code

~~~shell
  newt generate app.yaml
~~~

> Renders SQL, PostgREST conf, Mustache templates

# Step four, setup Postgres and PostgREST

1. Use the generated SQL and configuration
2. Setup and check via `createdb` and `psql`

# Step four, setup Postgres and PostgREST

~~~shell
  createdb app
  psql app -c '\i setup.sql'
  psql app -c '\i models.sql'
  psql app -c '\dt'
~~~

> should this be automated too?

# Step five, run your application and test

~~~shell
  newt run app.yaml
~~~

> Point your web browser at http://localhost:8010 to test

# Can I run a demo?

Not yet, hopefully in late May 2024.

# Second prototype Status

- A work in progress (April 2024)
- Working prototype target date June 2024
- Using internal applications as test bed

# How much is built?

- [X] Newt developer tool
- [X] Router is implemented and working
- [X] Mustache template engine is working
- [ ] Generator development (in progress)
- [ ] Modeler (design stage)

# Insights from prototypes 1 & 2

- "Off the shelf" is simpler
- Lots of typing discourages use

# Insights from prototypes 1 & 2

- SQL turns people off, use a code generator
- Hand typing templates is a turn off, use a code generator
- Large YAML structures benefit from code generation
- Automatic "wiring up" of routes and templates very helpful

# What's next to wrap up prototype 2?

- Debug and improve the code generator
- Implement a data modeler

# Unanswered Questions

- What should be the minimum knowledge needed to use Newt?
- What should come out of the box with Newt?
    - GUI tools?
    - Web components?
    - Ready made apps?

# Someday, maybe ideas

- SQLite 3 database support
- A S3 protocol web service implementing object storage using OCFL
- Web components for library, archive and museum metadata types
- Visual programming would be easier than editing YAML files

# Related resources

- Newt <https://github.com/caltechlibrary/newt>
- Postgres <https://postgres.org> + PostgREST <https://postgrest.org>
- [Mustache](https://mustache.github.io) programming languages support

# Thank you!

- This Presentation
  - pdf: <https://caltechlibrary.github.io/newt/presentation2/newt-p2.pdf>
  - pptx: <https://caltechlibrary.github.io/newt/presentation2/newt-p2.pptx>
- Newt Documentation <https://caltechlibrary.github.io/newt>
- Source Code: <https://github.com/caltechlibrary/newt>
- Email: rsdoiel@caltech.edu

