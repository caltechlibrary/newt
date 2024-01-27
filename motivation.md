
## Motivation

My belief is that many web services used by archives, libraries and museums can benefit from a simplified and consistent back end. If the back end is "easy" then the limited developer resources can be focused on the front end which is what us humans experience day to day. This lead me to take a [RAD](https://en.wikipedia.org/wiki/Rapid_application_development) tools approach to the problem and to encourge the re-use of existing micro services like Pandoc server, PostgREST, Open Search.

I've written many web applications over the years. Newt is focused on providing very specific glue leveraging existing microservices already used by libraries, archives and museums.  For many of these apps the core of an application is a JSON service (e.g. Invenio-RDM, ArchivesSpace). Newt can be used to extend these if needed. Let's take advantage of that. When we do need a custom application let also take advantage of a similar microservices approach. Build your core application in SQL with PostgREST+Postgres, hand of rendering to Pandoc running as a service. Newt can route your data between them two giving you similar benefits to complicated systems like Invenio but simple enough to be implemented by a single person.

## Newt stack, front to back

- A front end web server (e.g. Apache 2, NginX) can provide access control where appropriate (e.g. single sign-on via Shibboleth)
- Newt provides static file services but more importantly serves as a data router. It can validate and map a request to a JSON source, take those results then send them through Pandoc for transformation.

- JSON data source(s) provide the actual metadata management
  - Postgres+PostgREST is an example of a JSON source integrated with a SQL server
  - Solr, Elasticsearch or Opensearch can also function as a JSON source oriented towards search
  - ArchivesSpace, Invenio RDM are examples of systems that can function as a JSON sources
  - CrossRef, DataCite, ORCID are examples of services that function as JSON sources
- Pandoc server provides a templating engine to transform data sources

All these can be treated as "off the shelf". I.e. we're not writing them from scratch but we're accessing them via configuration. Even using PostgREST with Postgres the "source" boils down to SQL used to define the data models hosted by the SQL service.  Your application is implemented using SQL and configured with YAML and Pandoc templates.

## Taking advantage of JSON and S3 data sources

Newt was inspired by my working with PostgREST, Postgres and Pandoc. I also work allot of S3 object stores. I want Newt to be light weight. I wanted Newt to avoid writing anything to disk. That's possible now working with JSON API as data sources. I am in the planning stages of adding S3 protocol support to allow Newt applications to support a bigger domain space. Current plans are focused on using Minio as an off the shelf microservices to fill that responsibility.

## Getting Newt, Pandoc, PostgREST and Postgres

Newt is an experimental prototype (June/July 2023, and January/February 2024). It is distributed in source code form.  You need a working Go language environment, git, make and Pandoc to compile Newt from source code. See [INSTALL.md](INSTALL.md) for details. Go is available from <https://golang.org>.

Pandoc is available from <https://pandoc.org>. PostgREST is available from <https://postgrest.org>. Both can be built from source using [GHCup](https://www.haskell.org/ghcup/) to install the Haskell tool chain.

Pagefind is a Rust application. It can be installed by going to https://pagefind.app/ and following the installtion instructions there.

Postgres can be obtain through your host system's package manager. If not then you can also go to the <https://postgres.org> website and follow the instruction for installation.

