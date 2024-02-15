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
updateDate: 2024-02-14
#pubDate: TBD
#place: TBD
#date: TBF
section-titles: false
toc: true
keywords: [ "web service", "micro service", "Postgres", "PostgREST", "Mustache" ]
url: "https://caltechlibrary.github.io/newt/presentation2"
---

# Goal of the second prototype

Is Newt's toolbox and "off the shelf" software enough to simplify our applications?

# First, combine the right abstractions

- Off the shelf web services
- Aligning our expectations
- Composing application using data pipelines

# Off the shelf

- [Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org)o =>  JSON data source
- Newt Mustache a simple template engine => Transform JSON into web pages
- Newt new implementation of a data router ties it all together
- YAML used to express routes, pipelines and data models

# Assembling our with YAML

- Newt's Second Prototype has a new YAML syntax
	- CITATION.cff for application metadata
	- GitHub YAML issue template syntax used to model data
- Newt generates SQL and templates from this YAML
- Introducing Newt Mustache, a light weight template engine
	- Mapping your templates in the pipeline

# Demo time

Let's take a look at the demo ...

# So what did we see?

1. static site derived from CSV data file
    - CSV file holds our sighting list
    - Mustache templates can render our content
2. dynamic, SQL and browser JavaScript
    - (YAML and code generation)
    - (bird sighting list held in SQL database)
    - (JavaScript browser side assembles web site)
3. dynamic, SQL and no browser JavaScript
    - (bird sighting list held in SQL database)
    - (Mustache templates transform content)
    - (Data routing provides the website)

What were the trade-offs?

# [A static bird](https://github.com/caltechlibrary/newt/tree/main/demos/birds1), static site, Mustache templates

CSV file, Mustache, 2 directories, 4 files, **?? total line count**, static site

Lines   Files
------  ---------------
    ??  README.md
     ?  birds.csv <-- this is the seed data used in all the implementations
     ?  page.tmpl
    ??  htdocs/index.html


# [A dynamic site requiring browser side JS](https://github.com/caltechlibrary/newt/tree/main/demos/birds2), dynamic site, browser JavaScript

SQL (Postgres + PostgREST), Browser JavaScript, 2 directories, 9 files, **??? total line count**, dynamic site

Lines    Files
------   --------------
    ??   README.md
     ?   birds.csv <-- from birds1, loaded into SQL database
    ??   birds.yaml <--- use to generate the models only
    ??   setup.sql <-- code generation, newtgen
    ??   models.sql <-- code generation, newtgen
    ??   models_test.sql <-- generated code, newtgen
     ?   postgrest.conf <-- generated code, newtgen
    ??   htdocs/index.html   <-- hand coded
    ??   htdocs/sightings.js <-- hand coded

# [birds version 3](https://github.com/caltechlibrary/newt/tree/main/demos/birds3), dynamic site, no browser JavaScript

SQL (Postgres + PostgREST), Newt Mustache, Newt, 1 directory, 9 files, **225 total line count**, dynamic site

Lines    Files
------   ---------------
    43   README.md
     ?   birds.csv <-- from birds1, loaded into SQL database
    ??   birds.yaml <--- Newt models and routing
    34   setup.sql <-- code generation, newtgen
    60   models.sql <-- code generation, newtgen
    15   models_test.sql <-- code generation, newtgen
     3   postgrest.conf <-- code generation, newtgen
    ??   birds_tmpl.yaml <-- hand coded, maps our templates to request in Newt Mustache
    36   page.tmpl
     7   post_result.tmpl

# Insights from experiment

- A few "Off the shelf" web services can make application construction easier
- A few "Off the shelf" can also shrink what we maintain
- SQL turns some people off, so let's use a code generator
- Modify SQL, even when you know it, is often easier than typing it into an empty file
- Mustache templates are simple to learn, see the [Mustache Web site](https://mustache.github.io)
- Mustache templates are easy to debug either from Newt Mustache or the Mustache cli
- Newt really caters to the newly vintage of HTML5 and progressive enhancement

- I encountered an unexpected result ...

# An unexpected result

- Running more web services has it's own cognitive price, are you OK with that?
- Containers to run many web services just adds to complexity
- Thinking about an echo system of applications that can use the same web services is promising
  - Can we safely aggregate our Newt YAML files?
- The web services need to provide a "developer" model for debugging
- Writing tests with [curl](https://curl.se/) gets tedious, newtgen needs to generate this code
- Newt can be used as a large static site generator, ex. feeds v2

# Newt has weaknesses

- Newt is **an experimental working prototype** (February 2024)
- You really want to minimize your pipelines to one to three stages
- An app with too many routes is hard to maintain, avoid doing that
- While Newt uses YAML, you need to learn the YAML syntax for three Newt tools
  - Newt data router and code generator share one syntax
  - Can Newt Mustache and pdbundler share one syntax?
  - Can the code generator produce the Newt Mustache YAML?
- Newt doesn't support file uploads
  - That would be nice to broaden Newt's application domain

# Newt stack has strengths

> A very mature foundation

- 20th Century tech
  - SQL (1974), HTTP (1991), HTML (1993), Postgres (1996)
- 21st Century tech
  - JSON (2001), YAML (2001), Pandoc (2006), Mustache templates (2009), PostgREST (2014)

# Next steps for Newt Project?

1. Consider allowing the Newt data router to access non-localhost resources
2. Build public facing applications, Spring to Winter 2024
3. Improve documentation and add tutorials

# Newt's someday, maybe ...

- Integrating S3 protocol support would allow file uploads to be addressed
- Explore light weight options to use SQLite 3 databases
- Explore a web service that would enable non-linear conditional pipelines
- Visually programming to generate all the Newt YAML files
- A Newt community to share YAML, SQL, templates and web components

# Related resources

- Newt <https://github.com/caltechlibrary/newt>
- Go 1.22 and it's pattern language for HTTP handlers, see <https://pkg.go.dev/net/http#hdr-Patterns>
- Postgres <https://postgres.org> + PostgREST <https://postgrest.org>
  - PostgREST Community Tutorials <https://postgrest.org/en/stable/ecosystem.html>
- Compiling Pandoc or PostgREST requires Haskell
  - Install Haskell GHCup <https://www.haskell.org/ghcup/>
  - [Quick recipe, Compile PostgREST (M1)](https://rsdoiel.github.io/blog/2023/07/05/quick-recipe-compiling-PostgREST-M1.html)
  - [Quick recipe, Compiling Pandoc (M1)](https://rsdoiel.github.io/blog/2023/07/05/quick-recipe-compiling-Pandoc-M1.html)
- [Mustache templates](https://mustache.github.io) and programming language implementations
  - [Go Mustache package used by Newt Project](https://github.com/cbroglie/mustache), also provides a mustache cli

# Thank you!

- This Presentation <https://caltechlibrary.github.io/newt/presentation2/>
- Project: <https://github.com/caltechlibrary/newt>
- Email: rsdoiel@caltech.edu


