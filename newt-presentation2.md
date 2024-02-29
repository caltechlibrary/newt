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
updateDate: 2024-02-15
#pubDate: TBD
#place: TBD
#date: TBF
section-titles: false
toc: true
keywords: [ "web service", "micro service", "Postgres", "PostgREST", "Mustache" ]
url: "https://caltechlibrary.github.io/newt/presentation2"
---

# Goal: Second Prototype

Is Newt and "off the shelf" enough?

# Choose wisely: pick your abstractions

- Simple = (No coding) + (Less code)
- Compose applications, use data pipelines
- Align service and deliverables

# Off the shelf (no coding)

- [Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org) =>  JSON data source
- Newt Mustache => Transform JSON into web pages
- Newt Router, ties it all together
- YAML can express routes, pipelines and data models

# Assembling it with YAML (less code)

- Newt's Second Prototype has a new YAML syntax
  - GitHub YAML issue template syntax
- Code generator, "look Mom, no AI!"
- Newt Mustache, "YAML your templates in your pipeline"

# Demo time

A second prototype Newt application

# Step one create our YAML file

- text editor

# Step two generate our files

- newtgenerator
- text editor

# Step three implement our JSON data source

- psql
- postgrest

# Step four define routes and map templates

- text editor

# Step five, fire everything up and test

- PostgREST 
- newt
- web browser

# What did we see?

- Code generation
- Data pipelines and a bunch-o-services
- Minimal browser expectations

# Insights from prototypes 1 & 2

- "Off the shelf" is simpler
- SQL turns people off, use a code generator
- Mustache templates, simple but not too simple
- Two YAML syntaxes is too many
- Unexpected encounters ...

# Unexpected encounters

- Web services' cognitive price
- Containers add to complexity
- Web services need a "developer" mode for debugging
- Newt, **a static site generator platform**

# Newt has weaknesses

- Newt is **an experimental working prototype** (February 2024)
- Pipeline delays accumulate
- Many routes are hard to maintain
- Newt is missing file upload support

# Newt has strengths

> A very mature foundation

- 20th Century tech
  - SQL (1974), HTTP (1991), HTML (1993), Postgres (1996)
- 21st Century tech
  - JSON (2001), YAML (2001), Mustache templates (2009), PostgREST (2014)

# Questions guiding next steps

- Should the router support non-localhost URLs?
- What tools should come from the Newt Project?
- How do we move Newt beyond a Caltech Library experiment?

# My wish list ...

- simple web service with S3 protocol implementing current OCFL
- SQLite 3 database support
- Visually programming for Newt YAML files
- Web components for LAM metadata types

# Related resources

- Newt <https://github.com/caltechlibrary/newt>
- Postgres <https://postgres.org> + PostgREST <https://postgrest.org>
- [Mustache](https://mustache.github.io) programming languages support
- [Go Mustache package](https://github.com/cbroglie/mustache), provides a mustache cli
- Go 1.22, pattern language in HTTP handlers, see <https://pkg.go.dev/net/http#hdr-Patterns>

# Thank you!

- This Presentation <https://caltechlibrary.github.io/newt/presentation2/>
- Newt Documentation <https://caltechlibrary.github.io/newt>
- Project: <https://github.com/caltechlibrary/newt>
- Email: rsdoiel@caltech.edu


