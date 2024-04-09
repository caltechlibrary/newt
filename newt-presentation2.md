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

# Goal: Second Prototype Question

Is Newt and "off the shelf" enough?

# Continue to choose wisely:

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

- Work in progress but coming together (as of 2024-04-08)
- Hope to have a working prototype by June 2024
- Implementing internal applications using prototype as test bead

# Is there a Demo I can run?

I'm somewhere between vaporware and working prototype

# What would a demo look like?

# Step one create our YAML file

~~~
newt init people.yaml
vi people.yaml
~~~

- generates people.yaml
- edit people.yaml, replace models with models.txt

# Step two generate our SQL files

~~~
newtgenerator people.yaml postgres setup >setup.sql
newtgenerator people.yaml postgres models >models.sql
~~~

- newtgenerator
- edit your SQL files as needed

# Step three generate PostgREST config

~~~
newtgenerator people.yaml postgrest >postgrest.conf
~~~

- newtgenerator

# Step four, setup your data source in Postgres

~~~
createdb people
psql people
\i setup.sql
\i models.sql
\dt
\q
~~~

- psql

# Step five, Generate Mustache templates

~~~
newtgenerator people.yaml mustache create_form people >create_people_form.tmpl
newtgenerator people.yaml mustache create_response people >create_people_response.tmpl
newtgenerator people.yaml mustache update_form people >update_people_form.tmpl
newtgenerator people.yaml mustache update_response people >update_people_response.tmpl
newtgenerator people.yaml mustache delete_form people >delete_people_form.tmpl
newtgenerator people.yaml mustache delete_response people >delete_people_response.tmpl
newtgenerator people.yaml mustache read people >read_people.tmpl
newtgenerator people.yaml mustache list people >list_people.tmpl
~~~

- newtgenerator

# Finally, run newt and test

~~~
newt run people.yaml
~~~

- fire up newt, test and debug
- web browser

# What did we see?

- Code generation (SQL, PostgREST config, Mustache Templates)
- Data pipelines (using PostgREST and Newt Mustache)
- Minimally functional web app

# Insights from prototypes 1 & 2

- "Off the shelf" is simpler
- SQL turns people off, use a code generator
- HTML and Mustache need a code generator
- Generating the "wiring up" of routes and templates is helpful

# Lessons learned

- Managing routes and pipelines has a cognitive price
- Keep your pipelines short
- Web services need a "developer" mode for debugging

# What's next?

- Improve the generated code
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


