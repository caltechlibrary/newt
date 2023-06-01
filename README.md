
# Newt, my new take on the web stack

Newt is a program for working with off the shelf [microservices](https://en.wikipedia.org/wiki/Microservices). Specifically it can route a request to a JSON API such as provided by [Postgres](https://postgresql.org) and [PostgREST](https://postgrest.org) then take the result and run it through [Pandoc server](https://pandoc.org). Newt can also function as a light weight static file server.  It runs on localhost only. In a production setting you would run it behind a front-end web server like NginX or Apache 2.


## Motivation

My belief is that many web services used by Archives, Libraries and Museums can benefit from a simple back-end built from a few off the shelf microservices. If the back-end is "easy" then the limited developer resources can be focused on the front-end. **An improved front-end offers opportunities to provide a more humane user experience for staff and patrons.**

## A basic Newt approach

From front-end to back-end

- A front end web server (e.g. Apache 2, NginX) providing single sign on and access control (e.g Apache2 + Shibboleth)
- Newt is used to route data request to a JSON API and run the results through Pandoc. The routed data is specified in a CSV file read at startup. Newt also provides static file services. 
- Postgres+PostgREST (provides a JSON data API, it is programmed in SQL)
- Pandoc server (provides a templating engine to transform data API responses to HTML or other formats)

All these can be treated as "off the shelf". Aside from configuration they run as traditional services on POSIX systems.  Your application is primarily developed using SQL (defining the behavior our PostgREST), Pandoc templates used to turn JSON into HTML. Any additional static content (e.g. web pages, CSS, JavaScript for the web browser) may also be provided by Newt since it includes a static file service ability.

## Elaborating beyond Postgres+PostgREST

Similar to the above example you could have Newt integrate with a full text search service via Solr, Elasticsearch or Opensearch. The requirement is that the data source return JSON data. In that case the stack from back to front would look like.

- A front end web server (e.g. Apache 2, NginX) integrated providing access control where appropriate (e.g. Single Sign on with OAuth2 or Shibboleth)
- Newt URL router (uses a CSV file for defining routes)
- JSON data sources
    - Postgres+PostgREST (provides a data source API, is programmed in SQL)
    - Search Engine (e.g. Solr, Opensearch)
- Pandoc server (provides a templating engine to transform data sources)

The routes your define for Newt would select the appropriate data source.

The increase in complexity would be the additional routes you define in our
CSV file as well as any additional templating needed to render the search
engine results via Pandoc.


## Explore Newt and friends

### Demo 1, Pandoc

Pandoc provides a powerful template engine often associated with building static sites. It can run as a microservice (more on that below).  The [birds 1](birds1/) demo shows a simple bird listing site built from static content in Markdown, a list of birds in a CSV file and stitched together using Pandoc and Pandoc templates. 

### Demo 2, Postgres + PostgREST

The combination of Postgres and PostgREST provides a rich data API defined from the data models in Postgres. In [birds 2](birds2/) demo PostgREST runs on localhost and we use a localhost web server to present the static files. The difference is that our list of birds is held in the Postgres database and provides via a JSON API call using JavaScript.

- [Building with Postgres and PostgREST](building-with-postgres-postgrest.md) discusses the approach taken to create the [birds](birds2/) demo
- Extras for setting up a developer environment
    - [setup-birds2.bash](setup-birds2.bash), a bash script that generates the contents of demo's running code
    - [setup-developer-account.bash](setup-developer-account.bash) is a bash script that displays Postgres commands for setting up a super user account for development.

## Demo 3, pulling everything together with Newt

In demo 2, we performed page assembly browser side. This is a huge chore. We can skip the complex JavaScript for that and use Newt to route data from PostgREST or a search engine through Pandoc before returning it to the web browser. The [birds 3](birds3/) demo shows how this works. With just SQL and Pandoc templates we can build a simple web application.

Configuring Newt can be done in a YAML file or through the shell's environment. A minimum setup requires a pointer to the CSV file container your route definitions. The environment variable needed is `NEWT_ROUTES` in the YAML file the attribute is `newt_routes`. If you also wish to have static file service support then you would set `NEWT_HTDOCS` in the environment or `newt_htdocs` in the YAML configuration file. 

## Conclusion

In the final demo I've delegated back-end tasks  a small collection of microservices. The programming was reduced to SQL and some Pandoc templating.

Postgres
: provides data storage, back-end configuration and data services

PostgREST
: Turns a Postgres database into a DATA API

Pandoc
: Run in server mode is a powerful template engine

Newt
: Is a data router and static file server. Newt knows how to talk to a data source like PostgREST or Solr and then run the results through Pandoc.

With it we can retrieve data from our database and see formatted results. We can use standard web forms to update our database as well.
