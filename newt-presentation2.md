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

# Goal: Answer a question.

Is Newt and "off the shelf" enough?

# How I am proceeding

- Pick Simple = (No coding) + (Less coding)
- Avoid inventing new things
- Compose applications using data pipelines and templates

# Off the shelf (no coding)

- [Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org) =>  JSON data source
- Newt Mustache => Transform JSON into web pages
- Newt Router, ties it all together
- YAML can express routes, pipelines and data models

# Assembling it with YAML (less coding)

- GitHub YAML issue template syntax described data models
- Template language is now Mustache
- More Code generation, "look Mom, no AI!"

# Second prototype status

- Still a work in progress (as of 2024-04-08)
- Hope to have a working prototype by June 2024
- Implementing internal applications using prototype as test bead

# Is there a Demo I can run?

I'm somewhere between vaporware and working prototype

# What's working, what's not?

- [X] Router is implemented and working
- [X] Mustache template engine is working
- [ ] Generator is being debugged and improved

# How do I think things will work?

1. Start by designing our data models
2. Generate the rest of our Newt YAML
3. Generate SQL and templates
4. Refine generated code

# How is the data model is described?

- GitHub YAML Issue Template Syntax
  - can render HTML
  - can map to SQL

# Step one create our YAML file

~~~
newt init app.yaml
vi app.yaml
~~~

- Interactively generate app.yaml
- Create a models description in YAML in app.yaml

# Step two generate our SQL files

~~~
newtgenerator app.yaml postgres setup >setup.sql
newtgenerator app.yaml postgres models >models.sql
~~~

- newtgenerator
- edit your SQL files if needed

# Step three generate PostgREST config

~~~
newtgenerator app.yaml postgrest >postgrest.conf
~~~

- newtgenerator

# Step four, setup your data source in Postgres

~~~
createdb app
psql app
\i setup.sql
\i models.sql
\dt
\q
~~~

- psql

# Step five, Generate Mustache templates

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

- newtgenerator, needs to be made much simpler ...

# Finally, run newt and test

~~~
newt run app.yaml
~~~

- fire up newt, test and debug
- web browser

# What would these steps show us?

- Code generation (SQL, PostgREST config, Mustache Templates)
- Data pipelines (using PostgREST and Newt Mustache)
- Minimally functional web app

# Insights from prototypes 1 & 2

- "Off the shelf" is simpler
- Large YAML structures benefit from code generation
- SQL turns people off, use a code generator
- HTML and Mustache need a code generator
- Generating the "wiring up" of routes and templates is helpful

# Lessons learned, so far

- Managing routes and pipelines has a cognitive price
- Keep your pipelines short
- Web services need a "developer" mode for debugging

# What's next?

- Find nice way to bootstrap our data models
- Improve the generated code, figure out better way to bootstrap modeling
- Simplify using Newt (too many steps to type)
- Can I simplify the Newt YAML further?

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


