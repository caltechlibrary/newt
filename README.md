
# Newt, a new take of the web stack

Newt is two things. Newt is short for "new take". Newt is also an experimental web service for working with off the shelf [microservices](https://en.wikipedia.org/wiki/Microservices). The specific microservices Newt is intend to help explore are [Postgres](https://postgresql.org), [PostgREST](https://postgrest.org), [Pandoc server](https://pandoc.org). The Newt program can be thought of as both a data router (i.e. make requests to PostgREST and Pandoc) an a light weight static file server. It only runs on localhost so normally would be behind a front-end web server like NginX or Apache 2.

## Motivation

My belief is that many web services used by Archives, Libraries and Museums can benefit from a simplier back-end. If the back-end is "easy" then the limited developer resources can be focused on the front-end. An improved front-end offers opportunities to provide a more humane user experience for staff and patrons.

## A basic Newt approach

From back to front

- Postgres+PostgREST (provides a data API, it is programmed in SQL)
- Pandoc server (provides a templating engine to transform data sources)
- Newt is a URL router (uses a CSV file for defining routes) and provides static file services
- A front end web server (e.g. Apache 2, NginX) can provide access control where appropriate (e.g. Single Signon with OAuth2 or Shibboleth)

All these can be treated as "off the shelf". Aside from configuration they run as traditional services on POSIX systems.  Your application is primarily developed using SQL (defining the behavior our PostgREST), Pandoc templates used to turn JSON into HTML. Any additional static content (e.g. web pages, CSS, JavaScript for the web browser) may also be provided by Newt.

## A more complex stack

Similar to the above example you could include full text search via Solr, Elasticsearch or Opensearch. In that case the stack from back to front would look like.

- Postgres+PostgREST (provides a data source API, is programmed in SQL)
- Search Engine (e.g. Solr, Opensearch)
- Pandoc server (provides a templating engine to transform data sources)
- Newt URL router (uses a CSV file for defining routes)
- A front end web server (e.g. Apache 2, NginX) integrated providing access control where appropriate (e.g. Single Signon with OAuth2 or Shibboleth)

The routes your define for Newt would select the appropriate data source.

The increase in complexity would be the additional routes you define in our
CSV file as well as any additional templating needed to render the search
engine results via Pandoc.


## Explore Newt and friends

### Step 1, Pandoc

Pandoc a powerful template engine often assoicated with building static sites. It can run as a microservice (more on that below).  The [birds 1](birds1/) demo shows a simple bird listing site built from static content in Markdown, a list of birds in a CSV file and stitched together using Pandoc and Pandoc templates. 

### Step 2, Postgres + PostgREST

The combination of Postgres and PostgREST provides a rich data API defined from the data models in Postgres. In [birds 2](birds2/) demo PostgREST runs on localhost and we use a localhost web server to present the static files. The differents is that our list of birds is held in the Postgres database and provides via a JSON API call using JavaScript.

- [Building with Postgres and PostgREST](building-with-postgres-postgrest.md) discusses the approach taken to create the [birds](birds2/) demo
- Extas for setting up a developer environment
    - [setup-birds.bash](setup-birds.bash), a bash script that generates the contents of demo's running code
    - [Multipass basics](multipass-basics.md), multipass runs Ubuntu VM which can be used to run the demo
    - [newt-init.yaml](newt-init.yaml) provides the configuration for a multipass based VM to run the demo
    - [setup-developer-account.bash](setup-developer-account.bash) is a bash script that displays Postgres commands for setting up a super user account for development.

## Step 3, what we've learned togher with Newt

Managing page assembly browser side is a chore. We can skip the complex JavaScript for that and use Newt to route data from PostgREST or a search engine through Pandoc before returning it to the web browser. The [birds 3](birds3/) demo shows how this works. With just SQL and Pandoc templates we can build a web application.

Configuring Newt can be done in a YAML file or through the shell's enviroment. A minimum setup requires a pointer to the CSV file container your route definitions. The environment variable needed is `NEWT_ROUTES` in the YAML file the attribute is `newt_routes`. If you also wish to have static file service support then you would set `NEWT_HTDOCS` in the environment or `newt_htdocs` in the YAML configuration file. 

- [Building web Postgres, PostgREST, Pandoc and Newt](building with postgres-postgrest-pandoc-and-newt.md)

## Conclusion

In all our demos we've devided the tasks through a series of flexible microservices.

Postgres
: provides data storage, back-end configuration and data services

PostgREST
: Turns a Postgres database into a DATA API

Pandoc
: Run in server mode is a powerful template engine

Newt
: Is a data router and stastic file server. Newt knows how to talk to a data source like PostgREST or Solr and then run the results through Pandoc.

The "coding" of the back-end is reduced to SQL and Pandoc templates. The front-end can be as complex as you like given the full capabilities of your web browser to handle HTML, CSS and JavaScript.

