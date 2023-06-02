
# Newt, my new take on the web stack

Newt is two things. It is short a "new take" on building web applications. It is also an experimental microservice for working with other off the shelf [microservices](https://en.wikipedia.org/wiki/Microservices). 

Currently I am exploring using Newt to integrate [Postgres](https://postgresql.org), [PostgREST](https://postgrest.org) and [Pandoc](https://pandoc.org). 

The Newt application can be thought of as both a data router and a light weight static file server.

Newt can route a request to a JSON data API and then optional send that result through Pandoc for further processing. Newt runs as a localhost service only. To use in a production setting it would site behind a traditional web server like Apache 2 or NginX.

## Motivation

My belief is that many web services used by Archives, Libraries and Museums can benefit from a simplified and consistent back-end. If the back-end is "easy" then the limited developer resources can be focused on the front-end. **An improved front-end offers opportunities to provide a more humane user experience for staff and patrons.** 

## Newt's approach

From front-end to back-end

- A front end web server (e.g. Apache 2, NginX) can provide access control where appropriate (e.g. Single Sign on with OAuth2 or Shibboleth)
- Newt provides static file services but more importantly serves as a data router. I can map a request to a JSON data API, take those results then send them through Pandoc for transformation.
- Postgres+PostgREST is an example of a JSON data API
- Pandoc server provides a templating engine to transform data sources

All these can be treated as "off the shelf". Aside from configuration they can run as traditional services on most POSIX systems.  Your application is implemented using SQL. It is enhanced by using Pandoc templates used to turn JSON into HTML (or other desired formats). 

## Beyond Postgres+PostgREST

Newt can talk to any JSON data API as long as it returns its content in JSON. This means Newt can also route content to Solr or Elasticsearch before taking the results and sending them through Pandoc. This more elaborate setup would look like this from front-end to back-end.

- A front end web server (e.g. Apache 2, NginX) integrated providing access control where appropriate (e.g. Single Sign on with OAuth2 or Shibboleth)
- Newt URL router (uses a CSV file for defining routes)
- Search Engine (e.g. Solr, Opensearch)
- Postgres+PostgREST (provides a data source API, is programmed in SQL)
- Pandoc server (provides a templating engine to transform data sources)

Once again the routes define where the requests go (e.g. Solr or PostgREST). Pandoc templates transform the JSON responses into HTML result pages.

The main additional work beyond the initial scenario is creating the Pandoc templates you need to display the results from querying Solr or Elasticsearch.


## Exploring Newt's friends

- Pandoc
- Postgres and PostgREST

### Pandoc

Pandoc can be thought of as a powerful template engine. It is familiar to many people who building "static" websites. It is known for being able to convert many structured text formats from one to another. The [birds 1 demo](birds1/), a simple bird listing site,  is built using Pandoc with content written Markdown and a list of birds sighted maintained in a CSV file. Pandoc is used to render both into the HTML needed to display the website.

### Birds 1 Demo files

(see <https://github.com/caltechlibrary/newt/tree/main/birds1>)

- README.md, demo read me
- birds.csv, our list of bird sightings
- birds.md, additional text descriptions
- build.sh, a shell script that uses Pandoc to render site
- htdocs, the website document directory
    - index.html
- page.tmpl, the Pandoc page template

### Birds 1 Demo

Pros                             Cons
-------------------------------- -----------------------------------
Pandoc is widely known           Updates require changing birds.csv
Pandoc is easily scripted        A static web server is needed

### Postgres + PostgREST

Postgres and PostgREST can be combined to provide a dynamic data source for our [birds 2 demo](birds2/) website. In [birds 2 demo](birds2/) PostgREST runs on localhost and we use another web server to presents the static files. The difference from birds 1 demo is that our list of birds is held in the Postgres database and provided via a JSON API call using JavaScript in the web browser.

The back-end is primarily written using SQL. This includes setting up access by PostgREST. The font-end is more complex as JavaScript is required in the browser to assemble pages results from our JSON data API.

#### Birds 2 demo

(see <https://github.com/caltechlibrary/newt/tree/main/birds2>)

- README.md, demo read me
- birds.csv, our sightings data, loads into SQL database
- htdocs, the website document directory
    - index.html
    - sightings.js, JavaScript support is required in the web browser to render the page and handle updates
- postgrest.conf, configuration for PostgREST
- setup.sql, SQL setting up data models and JSON data API

Pros                              Cons
--------------------------------  ------------------------------------
Simple back-end, just SQL code    JavaScript in the browser is complex
                                  You need a static file server
                                  You need to know some SQL


### Assembling responses with Newt

(see <https://github.com/caltechlibrary/newt/tree/main/birds3>)

Managing page assembly browser side is a chore. We can skip the complex JavaScript browser side if we let Newt do the routing of requests to PostgREST and then sending those results through Pandoc. The birds 3 demo shows how this works. With just SQL and Pandoc templates we can build a web application.

#### Birds 3 Demo

- README.md, demo read me
- birds-routes.csv, CSV holding data routing instructions
- birds.csv, CSV holding data to be loaded via SQL
- birds.yaml, configuration for Newt
- htdocs
    - index.html, static homepage, links feature pages
    - add-bird.html, a webform for adding bird sighthings
    - about.html, static page, credits
- page.tmpl, Pandoc template for listing birds
- postgrest.conf, configurtion for PostgREST
- setup.sql, SQL setting of data models and JSON data API

This is very similar to both demo 1 and 2. Missing is build.sh from
demo 1. We don't need it since we're running Pandoc as a web service.
There is an added configuration file for Newt. The Pandoc template
script like in demo 1 and unlike demo 2 there is no sightings.js file to
send to the web browser. Also Newt provides our static file service so when developing  we don't need Apache 2 or NginX.

Pros                                 Cons
-----------------------------------  ------------------------------------
Simple back-end, just SQL code       Like demo 2 you need to know some SQL
No JavaScript required browser side
Newt provides the static web server

## More about Newt

Configuring Newt can be done in a YAML file or through the shell's environment. A minimum setup requires a pointer to the CSV file container your route definitions. The environment variable needed is `NEWT_ROUTES` in the YAML file the attribute is `newt_routes`. If you also wish to have static file service support then you would set `NEWT_HTDOCS` in the environment or `newt_htdocs` in the YAML configuration file. 

- [Building web Postgres, PostgREST, Pandoc and Newt](building with postgres-postgrest-pandoc-and-newt.md)

## Conclusion

In all our demos we've divided the tasks through a series of flexible microservices.

Postgres
: provides data storage, back-end configuration and data services

PostgREST
: Turns a Postgres database into a DATA API

Pandoc
: Run in server mode is a powerful template engine

Newt
: Is a data router and static file server. Newt knows how to talk to a data source like PostgREST or Solr and then run the results through Pandoc.

The "coding" of the back-end is reduced to SQL and Pandoc templates. The front-end can be as complex as you like given the full capabilities of your web browser to handle HTML, CSS and JavaScript.

