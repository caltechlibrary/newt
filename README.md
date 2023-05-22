
# Newt, a new take of the web stack

Newt is short for "new take". It is my personal new take on building web applications from a few off the shelf [microservices](https://en.wikipedia.org/wiki/Microservices). The microservices newt explores are [Postgres](https://postgresql.org), [PostgREST](https://postgrest.org), [Pandoc server](https://pandoc.org). The Newt program is a minimalist URL request router that services as a bridge between the other microservices and your frontend web server.

## Motivation

My belief is that many web services used by Archives, Libraries and Museums can benefit from a simplier back-end. If the back-end is "easy" then the limited developer resources can be focused on the front end creating a more humane user experience for staff and patrons.

## A basic Newt stack setup

From back to front

- Postgres+PostgREST (provides a data source API, is programmed in SQL)
- Pandoc server (provides a templating engine to transform data sources)
- Newt URL router (uses a CSV file for defining routes)
- A front end web server (e.g. Apache 2, NginX) integrated providing access control where appropriate (e.g. Single Signon with OAuth2 or Shibboleth)

All these can be treated as "off the shelf". Aside from configuration they run as traditional services on POSIX systems.  Your application is primarily developed using SQL (defining the behavior our PostgREST), Pandoc templates used to turn JSON into HTML. Any additional static content (e.g. web pages, CSS, JavaScript for the web browser) are provided by the front end web server. This is just an example. 

## Two more example "stacks"

Similar to the above example you could include full text search via Solr or Opensearch. In that case the stack from back to front would look like.

- Postgres+PostgREST (provides a data source API, is programmed in SQL)
- Search Engine (e.g. Solr, Opensearch)
- Pandoc server (provides a templating engine to transform data sources)
- Newt URL router (uses a CSV file for defining routes)
- A front end web server (e.g. Apache 2, NginX) integrated providing access control where appropriate (e.g. Single Signon with OAuth2 or Shibboleth)

Or you could integrate static file services with a localhost web services so that Newt presents all content to the Front end web server.

- Postgres+PostgREST (provides a data source API, is programmed in SQL)
- Search Engine (e.g. Solr, Opensearch)
- Back end static file server
- Pandoc server (provides a templating engine to transform data sources)
- Newt URL router (uses a CSV file for defining routes)
- A front end web server (e.g. Apache 2, NginX) integrated providing access control where appropriate (e.g. Single Signon with OAuth2 or Shibboleth)

It's basically include as little as you want or as much as you want. Newt manages contacting the data source, sends it to Pandoc server and returns the result to the font-end web server.

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

Managing page assembly browser side is a chore. Since we have Pandoc and it can run as a microservice we can also run Newt. Newt knows how to talk to a data source and send the results to Pandoc server before returning he transformed data to a web browser (or more normally font-end web server). The [birds 3](birds3/) demo shows how this works. With just SQL and Pandoc tempaltes we can build a web application.

- [Building with Postgres, PostgREST, Pandoc Server and Newt](


The second demo, [bees](bees/), builds on the birds demo by adding Pandoc server as a template engine for rendering content retrieved from Postgres via PostgREST. Pandoc service provide much of the same capability of the Pandoc command line program does. It is simple and safe to run as a localhost web service.  In bees I am presenting a bee spotting application with reports rendered via data from PostgREST and processed by Pandoc server.

- [Building with Postgres, PostgREST and Pandoc server](building-with-postgres-postgrest-and-pandoc-server.md)
- Extras for setting up demo (assumes you've already setup birds)
    - [setup-bees.bash](setup-bees.bash), a bash script that generates the contents of demo's running code

## Introducing Newt, a minimalist URL router

There are times when relying on JavaScript in the web browser to assemble a web page is not appropriate or just too much work. Newt URL router was developed to eleminate assembling the content provided by Postgres, PostgREST and Pandoc server.   Newt URL router takes care of mapping a public facing URL path to PostgREST and can optional send the content from PostgREST through Pandoc via Pandoc server. 

Newt URL router just needs to know where to contact PostgREST and Pandoc server and optional if you have a directory of static content you wish to expose too. Newt reads a table (e.g. a CSV file) to get the description of the routes, the transforms needed to query PostgREST and any template or Pandoc options desired if you want to run the retrieved content through Pandoc server. Newt takes care of coordinating the other microservices return content directly back to the web browser (in the developer setting) or to a front-end web server like Apache2 or NginX (e.g. in a production setting).

- [Building web Postgres, PostgREST, Pandoc and Newt](building with postgres-postgrest-pandoc-and-newt.md)
- Extras for setting up demo (assumes you've already setup birds and bees)
    - [setup-flowers.bash](setup-flowers.bash), a bash script that generates the contents of demo's running code

In this final demo, now including the Newt URL router, we use SQL to define our back-end service, Pandoc (via Pandoc server) to format our dynamic web poages and static files to complete our website implementation.

## Conclusion

In all our demos we've devided the tasks through a series of flexible microservices can can be used individually or in combination.

- Postgres, data storage, back-end configuration and data services
- PostgREST, a JSON data API
- Pandoc server as a template engine
- Newt a URL router for PostgREST, Pandoc server and static files
- (in production setting) a front end web server providing access control and user authentication (e.g. Apache2, NginX)

The "coding" left to someone developing a website or service can be as minimal as knowning some SQL, HTML and Pandoc.

