---
cff-version: 1.2.0
message: "If you use this software, please cite it as below."
type: software
title: "Newt"
abstract: "Newt is an experimental rapid application development. It is
a toolbox for creating metadata curation applications. Newt applications
implement a service orient architecture to minimize the code you have to
write and maintain. The Newt approach emphasizes off the shelf
components. Newt comes with a data router, a stateless template engine
and a developer tool. The developer tool includes data modeler, code
generator and app runner. Data models are described in YAML similar to a
YAML description of a web form. The code generator targets SQL and conf
files for running Postgres plus PostgREST. It can generate handlebar
templates from your modeled data. If manages the data routes and
template mapping. It can generate a TypeScript validation service based
on your data models. Newt composes your metadata curation application
from these services."
authors:
  - family-names: Doiel
    given-names: R. S.
    orcid: ""

repository-code: "git+https://github.com/caltechlibrary/newt"
version: 0.0.9

keywords: [ "web development", "web service", "service oriented
architecture", "micro service architecture", "micro
service", "Pandoc", "Mustache template", "JSON API", "data
router", "static web server", "template engine", "rapid application
development" ]

---

About this software
===================

## Newt 0.0.9

### Authors

- R. S. Doiel



Newt is an experimental rapid application development. It is a toolbox
for creating metadata curation applications. Newt applications implement
a service orient architecture to minimize the code you have to write and
maintain. The Newt approach emphasizes off the shelf components. Newt
comes with a data router, a stateless template engine and a developer
tool. The developer tool includes data modeler, code generator and app
runner. Data models are described in YAML similar to a YAML description
of a web form. The code generator targets SQL and conf files for running
Postgres plus PostgREST. It can generate handlebar templates from your
modeled data. If manages the data routes and template mapping. It can
generate a TypeScript validation service based on your data models. Newt
composes your metadata curation application from these services.


- GitHub: <git+https://github.com/caltechlibrary/newt>
- Issues: <git+https://github.com/caltechlibrary/newt/issues>


### Programming languages

- YAML
- SQL
- Go
- TypeScript

### Operating Systems

- Linux
- Windows
- macOS

### Software Requirements

- Dataset &gt;= 2.1
- Deno &gt;= 1.44
- Pandoc &gt; 3.1
- Golang &gt; 1.22
- A front end web server supporting reverse proxy (e.g. Apache2, NGINX)
