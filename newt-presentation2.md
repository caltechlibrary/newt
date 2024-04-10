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

# Goal: Answer a question.

Is Newt and "off the shelf" enough?

# How I am proceeding

- Pick Simple = (No coding) + (Less coding)
- Avoid inventing new things
- Compose applications using data pipelines and templates

# Off the shelf (no coding)

- [Postgres](https://postgresql.org) or [PostgREST](https://postgrest.org)
- [Solr](https://solr.apache.org) or [OpenSearch](https://opensearch.org)
- Newt Mustache => Transform JSON into web pages
- Newt Router, ties it all together

# Assembling it with YAML (less coding)

- GitHub YAML issue template syntax described data models
- YAML describes configuration, routes, pipelines
- Template language is now Mustache
- Code generation, "look Mom, no AI!"

# Second prototype status

- A work in progress (April 2024)
- Hope to have a working prototype by June 2024
- Internal applications will serve as test bed

# Is there a Demo I can run?

Not yet, hopefully soon.

# What's working, what's not?

- [X] Router is implemented and working
- [X] Mustache template engine is working
- [ ] Generator developement, in progress

# How do I think things will work?

1. Generate our app YAML
2. Designing our data models
3. Generate SQL, PostgREST config and templates
4. Run generated SQL
5. Run Newt and Test

# How is the data model is described?

- GitHub YAML Issue Template Syntax
  - describes HTML
  - implies SQL

# Step one create our YAML file

~~~
newt init app.yaml
~~~

- Interactively generate app.yaml

# Step two define our data models

~~~
vi app.yaml
~~~

- edit app.yaml

# Step three, generate our SQL files, config

~~~
newtgenerator app.yaml postgres setup >setup.sql
newtgenerator app.yaml postgres models >models.sql
newtgenerator app.yaml postgrest >postgrest.conf
~~~

- newtgenerator
- edit your SQL files if needed

# Step three, generate templates

~~~
newtgenerator app.yaml mustache create_form app >create_app_form.tmpl
newtgenerator app.yaml mustache create_response app >create_app_response.tmpl
newtgenerator app.yaml mustache update_form app >update_app_form.tmpl
newtgenerator app.yaml mustache update_response app >update_app_response.tmpl
newtgenerator app.yaml mustache delete_form app >delete_app_form.tmpl
newtgenerator app.yaml mustache delete_response app >delete_app_response.tmpl
newtgenerator app.yaml mustache read app >read_app.tmpl
newtgenerator app.yaml mustache list app >list_app.tmpl
~~~

- newtgenerator 
- this needs automation

# Step four, run our SQL

~~~
createdb app
psql app
\i setup.sql
\i models.sql
\dt
\q
~~~

- psql
- this could be automated

# Step five, run newt and test

~~~
newt run app.yaml
~~~

- fire up newt, test and debug
- web browser

# What was done?

- Code generation (YAML, SQL, PostgREST config, Mustache Templates)
- Data pipelines (using PostgREST and Newt Mustache)
- Minimally functional web app

# Insights from prototypes 1 & 2

- "Off the shelf" is simpler
- Large YAML structures benefit from code generation
- SQL turns people off, use a code generator
- Mustache/HTML needs a code generator
- Automatic "wiring up" of routes and templates is helpful

# Lessons learned, so far

- Managing routes and pipelines has a cognitive price
- Keep your pipelines short
- Web services need a "developer" mode for debugging

# What's next?

- Find nice way to bootstrap our data models
- Improve the generated code, consolidate actions
- Simplify using Newt (too many steps to type)
- Can Newt be simpler?

# Newt's challenges

- Newt is **a work in progress** (April 2024)
- Newt is missing file upload support

# Unanswered Questions

- What is the minimum knowledge needed to use Newt?
- What should come out of the box with Newt?
    - GUI tools?
    - Web components?
    - Ready made apps?

# My wish list ...

- SQLite 3 database support
- Visually programming would be easier than writing YAML files
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


