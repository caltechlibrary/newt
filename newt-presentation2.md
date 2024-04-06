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

- Simple = (No coding) + (Less coding)
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

~~~
newt init people.yaml
vi people.yaml
~~~

- generate people.yaml
- edit people.yaml

# Step two generate our SQL files

~~~
newtgenerator people.yaml postgres setup >setup.sql
newtgenerator people.yaml postgres models >models.sql
newtgenerator people.yaml postgres models_test >models_test.sql
~~~

- newtgenerator
- edit your SQL files as needed

# Step three generate PostgREST config

~~~
newtgenerator people.yaml postgrest >postgrest.conf
~~~

- newtgenerator

# Step four, setup your data sources

~~~
psql <setup.sql
psql <models.sql
psql <models_test.sql
~~~

- psql

# Step five, Mustache templates

~~~
newtgenerator people.yaml mustache create people >create_people.tmpl
newtgenerator people.yaml mustache read people >read_people.tmpl
newtgenerator people.yaml mustache update people >update_people.tmpl
newtgenerator people.yaml mustache delete people >delete_people.tmpl
newtgenerator people.yaml mustache list people >list_people.tmpl
newtgenerator people.yaml mustache page people >page_people.tmpl
~~~

- newtgenerator

# Step six, edit and test your templates

~~~
newtmustache people.yaml
~~~

- text editor
- mocked up JSON data
- curl and newt

# Step secen, testing your app

~~~
newt run people.yaml
~~~

- fire it up, test and debug
- web browser

# What did we see?

- Code generation (SQL, PostgREST config)
- Data pipelines and a bunch-o-services
- Minimal browser expectations

# Insights from prototypes 1 & 2

- "Off the shelf" is simpler
- SQL turns people off, use a code generator
- HTML and Mustache need a code generator

# Lessons learned

- Managing web services has a cognitive price
- Containers add to complexity
- Keep your pipelines short
- Web services need a "developer" mode for debugging

# What's next?

- Mustache template generation
- HTML file generation

# Newt's weaknesses

- Newt is **a work in progress** (April 2024)
- Newt is missing file upload support

# Newt's strengths

> A very mature foundation

- 20th Century tech
  - SQL (1974), HTTP (1991), HTML (1993), Postgres (1996)
- 21st Century tech
  - JSON (2001), YAML (2001), Mustache templates (2009), PostgREST (2014)

# Unanswered Questions

- Should the router support non-localhost URLs?
- What tools should come with the Newt?
- Should Newt move beyond a Caltech Library experiment?

# My wish list ...

- SQLite 3 database support
- Visually programming for Newt YAML files
- Web components for gallery, library, archive and museum metadata types
- A simple S3 protocol web service that implements storing object using OCFL

# Related resources

- Newt <https://github.com/caltechlibrary/newt>
- Postgres <https://postgres.org> + PostgREST <https://postgrest.org>
- [Mustache](https://mustache.github.io) programming languages support
- Go 1.22, pattern language in HTTP handlers, see <https://pkg.go.dev/net/http#hdr-Patterns>

# Thank you!

- This Presentation <https://caltechlibrary.github.io/newt/presentation2/>
- Newt Documentation <https://caltechlibrary.github.io/newt>
- Project: <https://github.com/caltechlibrary/newt>
- Email: rsdoiel@caltech.edu


