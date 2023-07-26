---
cff-version: 1.2.0
message: "If you use this software, please cite it as below."
type: software
title: "Newt"
abstract: "Newt is an experimental microservice providing two stage data
routing, form data validation, simple static file services. Newt can
route data from a request to a JSON data source
(e.g. Postgres+PostgREST) and then take the JSON results and run them
through Pandoc server."
authors:
  - family-names: Doiel
    given-names: R. S.
    orcid: "https://orcid.org/0000-0003-0900-6903"

repository-code: "git+https://github.com/caltechlibrary/newt"
version: 0.0.5

keywords: [ "web
development", "microservice", "postgres", "postgrest", "pandoc" ]

---

About this software
===================

## Newt 0.0.5

### Authors

- R. S. Doiel



Newt is an experimental microservice providing two stage data routing,
form data validation, simple static file services. Newt can route data
from a request to a JSON data source (e.g. Postgres+PostgREST) and then
take the JSON results and run them through Pandoc server.


- GitHub: <git+https://github.com/caltechlibrary/newt>
- Issues: <git+https://github.com/caltechlibrary/newt/issues>


### Programming languages

- YAML
- SQL
- Golang

### Operating Systems

- Linux
- Windows
- macOS

### Software Requiremets

- Postgres &gt;= 15
- PostgREST &gt;= 11
- Pandoc &gt;= 3
- Golang &gt; 1.20
- A front end web server supporting reverse proxy (e.g. Apache2, NginX)
