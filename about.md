---
cff-version: 1.2.0
message: "If you use this software, please cite it as below."
type: software
title: "Newt"
abstract: "Newt is an experimental web service providing data routing
and a multistage pipeline. The pipeline can contain services that work
with JSON data. Example object storage, form data validation, JSON
transformation, even static file services. Newt can route data from a
request to a JSON data source (e.g. Postgres+PostgREST or Solr) and then
run those results through a pipeline that may include Newt’s Mustache
template engine."
authors:
  - family-names: Doiel
    given-names: R. S.
    orcid: "https://orcid.org/0000-0003-0900-6903"

repository-code: "git+https://github.com/caltechlibrary/newt"
version: 0.0.6

keywords: [ "web development", "web service", "service oriented
architecture", "micro service
architecture", "microservice", "postgres", "postgrest", "pandoc", "mustache", "json
api" ]

---

About this software
===================

## Newt 0.0.6

### Authors

- R. S. Doiel



Newt is an experimental web service providing data routing and a
multistage pipeline. The pipeline can contain services that work with
JSON data. Example object storage, form data validation, JSON
transformation, even static file services. Newt can route data from a
request to a JSON data source (e.g. Postgres+PostgREST or Solr) and then
run those results through a pipeline that may include Newt’s Mustache
template engine.


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

- Postgres ≥ 16
- PostgREST ≥ 12
- Pandoc ≥ 3.1
- Golang ≥ 1.22
- A front end web server supporting reverse proxy (e.g. Apache2, NginX)
