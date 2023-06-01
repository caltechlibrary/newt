
# Newt, my new take on the web stack

Newt is two things. Newt is short for "new take". Newt is also an experimental web server for working with off the shelf [microservices](https://en.wikipedia.org/wiki/Microservices). The specific microservices Newt is intend to help explore are [Postgres](https://postgresql.org), [PostgREST](https://postgrest.org), [Pandoc server](https://pandoc.org). The Newt program can be thought of as both a data router (i.e. make requests to PostgREST and Pandoc) and as a light weight static file server. It runs on localhost only. It would normally would be behind a front-end web server like NginX or Apache 2.


## Motivation

My belief is that many web services used by Archives, Libraries and Museums can benefit from a simplifying back-end. If the back-end is "easy" then the limited developer resources can be focused on the front-end. **An improved front-end offers opportunities to provide a more humane user experience for staff and patrons.**

## A basic Newt approach

From front-end to back-end

- A front end web server (e.g. Apache 2, NginX) can provide access control where appropriate (e.g. Single Signon with OAuth2 or Shibboleth)
- Newt is a data router. It uses simple a CSV file for defining routes. It can also provide static files
- Postgres+PostgREST provides a JSON data API (programming is in SQL)
- Pandoc server (provides a templating engine to transform data sources)

All these can be treated as "off the shelf". Aside from configuration they run as traditional services on POSIX systems.  Your application is primarily developed using SQL (defining the behavior and content via PostgREST), Pandoc templates are used to turn JSON into HTML (or other desired formats). Any additional static content (e.g. HTML pages, CSS, JavaScript files) can also be made available view Newt.

## Beyond Postgres+PostgREST

Newt can talk to any data API as long as it returns its content in JSON. Thie means Newt can also route content to Solr or Elasticsearch before taking the results and sending them through Pandoc. This more elaborate setup would look like this from front-end to back-end.

- A front end web server (e.g. Apache 2, NginX) integrated providing access control where appropriate (e.g. Single Signon with OAuth2 or Shibboleth)
- Newt URL router (uses a CSV file for defining routes)
- Search Engine (e.g. Solr, Opensearch)
- Postgres+PostgREST (provides a data source API, is programmed in SQL)
- Pandoc server (provides a templating engine to transform data sources)

Once again the routes define where the requests go (e.g. Solr or PostgREST). Pandoc templates transform the JSON responses into HTML result pages.

The main additional work beyond the initial scenario is creating the Pandoc templates you need to display the results from querying Solr.


## Explore Newt and friends

### Step 1, Pandoc

Pandoc can be throught of as a powerful template engine. It is familair to many people who building "static" websites. It is known for being able to convert many structured text formats to another. The [birds 1 demo](birds1/), a simple bird listing site,  is built using Pandoc with content in both Markdown and a list of birds sighted in a CSV file. Pandoc does the stitching rendering the content to HTML.

Birds 1 Demo files

- [README.md](birds1/README.md)
- [birds.csv](birds1/birds.csv)
- [birds.md](birds1/birds.md)
- [build.sh](birds1/build.sh) 
- [htdocs](birds1/htdocs)
    - [index.html](birds1/htdocs/index.html)
- [page.tmpl](birds1/page.tmpl)

Pros                             Cons
-------------------------------- -----------------------------------
Pandoc is widely known           Updates require accessing birds.csv
Pandoc is easily scripted        A static web server is needed

### Step 2, Postgres + PostgREST

Postgres and PostgREST can be combined to provide a dynamic solution for our [birds 2 demo](birds2/) website. In [birds 2 demo](birds2/) PostgREST runs on localhost and we use another web server to presents the static files. The difference from birds 1 demo is that our list of birds is held in the Postgres database and provided via a JSON API call using JavaScript.

The back-end is primarily the SQL code setting things up. The font-end is more complex as JavaScript is required in the browser to assemble pages results from our JSON data API.

The files for birds 2 demo are

- [README.md](birds2/README.md)
- [birds.csv](birds2/birds.csv)
- [htdocs](birds2/htdocs/)
    - [index.html](birds2/htdocs/index.html)
    - [sightings.js](birds2/htdocs/sightings.js)
- [postgrest.conf](birds2/postgrest.conf)
- [setup.sql](birds2/setup.sql)


Pros                             Cons
-------------------------------- ------------------------------------
Simple backend, just SQL code    JavaScript in the browser is complex
                                 You need a static file server
                                 You need to know some SQL


## Step 3, let Newt assemble the pages

Managing page assembly browser side is a chore. We can skip the complex JavaScript browser side if we let Newt to route data to PostgREST or a search engine then run the results through Pandoc. The [birds 3](birds3/) demo shows how this works. With just SQL and Pandoc templates we can build a web application.

The files for birds 3 demo are

- [README.md](birds3/README.md)
- [birds-routes.csv](birds3/birds-routes.csv)
- [birds.csv](birds3/birds.csv)
- [birds.yaml](birds3/birds.yaml)
- [htdocs](birds3/htdocs/)
    - [index.html](birds3/htdocs/index.html)
- [page.tmpl](birds3/page.tmpl)
- [postgrest.conf](birds3/postgrest.conf)
- [setup.sql](birds3/setup.sql)

This is very similar to both demo 1 and 2. Missing is build.sh from
demo 1. We don't need it since we're running Pandoc as a web service.
There is an added configuration file for Newt. The Pandoc template
script like in demo 1 and unlike demo 2 there is no sightings.js file to
send to the web browser. Also Newt provides our static file service so indevelopment we don't need Apache 2 or NginX.

Pros                                 Cons
-----------------------------------  ------------------------------------
Simple backend, just SQL code        Like demo 2 you need to know some SQL
No JavaScript required browser side
Newt provides the static web server

## More about Newt

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

