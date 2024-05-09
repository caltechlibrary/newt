---
title: "Newt, a small system experiment"
author: "R. S. Doiel, <rsdoiel@caltech.edu>"
institute: |
  Caltech Library,
  Digital Library Development
description: Code4Lib Meet up, Los Angeles
urlcolor: blue
linkstyle: bold
aspectratio: 169
createDate: 2023-05-16
updateDate: 2023-06-28
pubDate: 2023-07-14
place: UCLA
date: July 14, 2023
section-titles: false
toc: true
keywords: [ "code4lib", "microservice", "Postgres", "PostgREST", "Pandoc" ]
url: "https://caltechlibrary.github.io/newt/presentation"
---

# The experiment

How do we make building web applications for Libraries, Archives and Museums simpler?

# Focus on three abstractions

- A JSON source for managing data => [Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org)
- A template engine => [Pandoc](https://pandoc.org)
- A data router and form validator => [Newt](https://github.com/caltechlibrary/newt/)

# Three implementations a bird sighting website

1. static site derived from CSV data file
    - (CSV file holds our sighting list)
2. dynamic, SQL and browser JavaScript
    - (bird sighting list held in SQL database)
3. dynamic, SQL and no browser JavaScript
    - (bird sighting list held in SQL database)

What are trade-offs?

# [birds version 1](https://github.com/caltechlibrary/newt/tree/main/demos/birds1), static site, Pandoc

CSV file, Pandoc, 2 directories, 5 files, **75 total line count**, static site

Lines   Files
------  ---------------
    26  README.md
     4  birds.csv <-- this is used in each of the demos
     6  build.sh
     7  page.tmpl
    32  htdocs/index.html

# [birds version 2](https://github.com/caltechlibrary/newt/tree/main/demos/birds2), dynamic site, browser JavaScript

SQL (Postgres + PostgREST), Browser JavaScript, 2 directories, 8 files, **232 total line count**, dynamic site

Lines    Files
------   --------------
    29   README.md
     4   birds.csv <-- from birds1
    34   setup.sql
    60   models.sql <-- implements our data models
    15   models_test.sql
     3   postgrest.conf
    24   htdocs/index.html   <-- hand coded
    63   htdocs/sightings.js <-- hand coded

# [birds version 3](https://github.com/caltechlibrary/newt/tree/main/demos/birds3), dynamic site, no browser JavaScript

SQL (Postgres + PostgREST), Pandoc, Newt, 1 directory, 9 files, **225 total line count**, dynamic site

Lines    Files
------   ---------------
    43   README.md
     4   birds.csv <-- from birds1
    34   setup.sql <-- from birds2
    60   models.sql <-- from birds2
    15   models_test.sql <-- from birds2
     3   postgrest.conf <-- from birds2
    23   birds.yaml
    36   page.tmpl
     7   post_result.tmpl

# Insights from experiment

- A few "Off the shelf" microservices can make application construction easier
- Orchestrating the data pipeline in YAML is reasonable
- SQL turns some people off
  - models could be bootstrapped from Newt's YAML
- Pandoc templates are simple to learn and well documented at [pandoc.org](https://pandoc.org)
- Newt stack plays well with HTML5 and front-end best practices
- I encountered an unexpected result ...

# An unexpected result

- Newt like PostgREST and Pandoc **do not require** shared synchronous state
- Postgres can be deployed in a [HA cluster](High-availability "high available cluster")

> The Newt stack can scale really big

# Newt has weaknesses

- Newt is **an experimental prototype** (June/July 2023, six weeks old)
- Newt doesn't support file uploads

# Newt stack has strengths

> A very mature foundation

- 20th Century tech
  - SQL (1974), HTTP (1991), HTML (1993), Postgres (1996)
- 21st Century tech
  - JSON (2001), YAML (2001), Pandoc (2006), PostgREST (2014)

# Next steps for Newt?

1. Test with Solr/Elasticsearch as alternate JSON sources
2. Build staff facing applications this Summer (2023)
3. Explore generating PostgREST configuration/SQL Models from Newt's YAML
4. (hopefully) move beyond my proof of concept in Fall/Winter (2023)

# Newt's someday, maybe ...

- Have Newt delegate file uploads to an S3 like service
  - One approach would be Minio using file streams
- Explore integrating SQLite3 support as a JSON data source
- Consider implementing Newt in Haskell for richer Pandoc integration
- A Newt community to share YAML, SQL and Pandoc templates

# Related resources

- Newt <https://github.com/caltechlibrary/newt>
- Postgres <https://postgres.org> + PostgREST <https://postgrest.org>
  - PostgREST Community Tutorials <https://postgrest.org/en/stable/ecosystem.html>
- Pandoc <https://pandoc.org>
    - Templates <https://pandoc.org/MANUAL.html#templates>
    - Pandoc Server <https://pandoc.org/pandoc-server.html>
- Compiling Pandoc or PostgREST requires Haskell
  - Install Haskell GHCup <https://www.haskell.org/ghcup/>
  - [Quick recipe, Compile PostgREST (M1)](https://rsdoiel.github.io/blog/2023/07/05/quick-recipe-compiling-PostgREST-M1.html)
  - [Quick recipe, Compiling Pandoc (M1)](https://rsdoiel.github.io/blog/2023/07/05/quick-recipe-compiling-Pandoc-M1.html)

# Thank you!

- This Presentation <https://caltechlibrary.github.io/newt/presentation/>
- Project: <https://github.com/caltechlibrary/newt>
- Email: rsdoiel@caltech.edu
