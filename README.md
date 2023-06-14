
# Newt

Newt is an experimental [microservice](https://en.wikipedia.org/wiki/Microservices) working with other "off the shelf" microservices. It is a proof of concept of a "new take" on the old web stack.

Newt integrates microservices such as a JSON data source with Pandoc. This is accomplished by routing the request first to the JSON data source then sending the results to Pandoc running as a web service. The the results are then sent back to the requesting web browser or service. Newt also provides a light weight static file service.

Newt's configuration can contain a set of "route" definitions. These map a request to the JSON source via an HTTP call.  The route is vetted (e.g. URL variables and form data are validated) before making the JSON source request.  If the request fails validation Newt responds directly to the requesting web browser or service via HTTP status codes and message. If validation succeeds then the JSON source is contacted and the results (optionally) run through Pandoc via another HTTP call. The final results are returned to the requesting web browser (or service).

An typical setup could be  Newt, [Postgres](https://postgresql.org) + [PostgREST](https://postgrest.org) and [Pandoc](https://pandoc.org). An equally valid setup could be Solr or Elasticsearch integrating with Pandoc. They key is realizing there is a JSON data source that can be interacted with via URL and HTTP method (e.g. GET, POST) and that we can then process the result via Pandoc running as a web service.

Newt's scope is limited to validating the requests, routing and returning content. It doesn't know secrets. It doesn't persist application state. This allows Newt instance(s) to be run singularly or in parallel as desired.

Newt runs as a localhost service. In a production setting you'd run Newt behind a traditional web server like Apache 2 or NginX. The front-end web service can provide access control via basic auth or single sign-on (e.g. Shibboleth). It plays nicely in a container environment or running straight up as a system service.


## Motivation

My belief is that many web services used by Archives, Libraries and Museums can benefit from a simplified and consistent back-end. If the back-end is "easy" then the limited developer resources can be focused on the front-end. **An improved front-end offers opportunities to provide a more humane experience for staff and patrons.**. I believe it is possible to build a highly useful back-end using off the shelf services like Postgres+PostgREST, Solr, Elasticsearch, Opensearch and Pandoc without requiring extensive programming much beyond defining your database models in SQL and declaring appropriate route maps in a YAML file.

If an existing application provides a JSON API (e.g. ArchivesSpace, Invenio-RDM) then it too can be integrated and extended with Newt and Pandoc. This should facilitate a "development at the edges" approach to enhancing legacy systems as well as opportunities for creating need systems to meet the growing needs of Libraries, Archives and Museums.

## Taking advantage JSON sources

Newt works with a JSON source that can be accessed via URL and HTTP methods (e.g. GET, POST, PUT, PATCH and DELETE). Many systems have JSON data API today (2023). This includes existing Library, Archive and Museum applications like Invenio RDM and ArchivesSpace as well as search engines like Solr and Elasticsearch.  It also includes many library services like CrossRef, DataCite and ORCID. These can be integrated via Newt in a easily managed may.

## Newt's approach

From front-end to back-end

- A front end web server (e.g. Apache 2, NginX) can provide access control where appropriate (e.g. single sign-on via Shibboleth)
- Newt provides static file services but more importantly serves as a data router. It can validate and map a request to a JSON source, take those results then send them through Pandoc for transformation.
- JSON data source(s) provide the actual metadata management
    - Postgres+PostgREST is an example of a JSON source integrated with a SQL server
    - Solr, Elasticsearch or Opensearch can also function as a JSON source oriented towards search
    - ArchivesSpace, Invenio RDM are examples of systems that can function as a JSON sources
    - CrossRef, DataCite, ORCID are examples of services that function as JSON sources
- Pandoc server provides a templating engine to transform data sources

All these can be treated as "off the shelf". I.e. we're not writing them from skratch but we're accessing them via configuration. Even using PostgREST with Postgres the "source" boils down to SQL used to define the data models hosted by the SQL service.  Your application is implemented using SQL and configured with YAML and Pandoc templates.

## Exploring Newt

Presented are three [demos](https://github.com/caltechlibrary/newt/tree/main/demos). The demo code can be rendered using one of three Bash scripts. The goal of the three demo is to show Pandoc, Postgres+PostgREST work individually and can work together with Newt. The bash scripts will generate all the files needed to run the demo. A "read me" for each demo describes how to run it.

[Birds 1 Demo](https://github.com/caltechlibrary/newt/blob/main/demos/setup-birds1.bash)
: Shows a simple use of Pandoc to render a static bird sightings website

[Birds 2 Demo](https://github.com/caltechlibrary/newt/blob/main/demos/setup-birds2.bash)
: Shows a dynamic bird sightings website using Postgres+PostgREST, but requires JavaScript running in your web browser

[Birds 3 Demo](https://github.com/caltechlibrary/newt/blob/main/demos/setup-birds3.bash)
: Shows a dynamic bird sightings website using Newt with Postgres+PostgREST and Pandoc

The goal of the three demos is to show an evolution towards simplicity while still enhancing capability.


### Birds 1 Demo, Pandoc only, a read only site

This is a simple static website. It introduces Pandoc. Pandoc is a good tool to transforming content.  In this case Pandoc takes a CSV file using a Pandoc tempalte to generating a simple "static" website. Static websites are generally simple to generate and maintain but require direct access to a host machine to update them.

- README.md, demo "read me" file
- birds.csv, our list of bird sightings
- build.sh, a shell script that uses Pandoc to render site (3 lines)
- htdocs, the website document directory
- page.tmpl, a Pandoc page template using by build.sh (13 lines)

Pros                             Cons
-------------------------------- -----------------------------------
Pandoc is widely known           Updates require changing birds.csv
Pandoc is easily scripted        and rebuilding the site.


### Birds 2 Demo, Postgres + PostgREST

Postgres and PostgREST can be combined to provide a JSON source (aka JSON API). This allows the website to be "dynamic". In birds 2 demo PostgREST runs on localhost and we use another web server to presents the static files and to proxy to PostgREST. The significant enhancements over birds 1 demo is our bird list is held in a Postgres database. It can be dynamically updated via a web browser thanks to PostgREST. The web browser uses JavaScript to call back to the JSON API. The JSON API results are rendered in the web browser via JavaScript.

The back-end is written using SQL. PostgREST understands Postgres and implements the JSON API. The web browser uses JavaScript to contact the JSON API and assemble pages results. The complexity is shifted out of the back-end is using an "off the shelf" PostgREST as a microservice. The back-end boils down to the SQL (50 lines) describing the data model.  Most of the complexity now resides in the front-end (aka web browser running JavaScript).

- README.md, demo "read me" file
- birds.csv, our prior list of bird sightings to be loaded into SQL database
- htdocs, the website document directory
- htdocs/sightings.js, JavaScript support is required in the web browser make requests, to render the page (60+ lines of code)
- postgrest.conf, configuration for PostgREST (3 lines)
- setup.sql, SQL setting up data models and JSON data API (50 lines)

Pros                              Cons
--------------------------------  --------------------------------------
Simple back-end, just SQL code    JavaScript in the browser is complex.
for data models.                  You need a separate static file server.
                                  You need to know some SQL.


### Birds 3 demo, assembling responses with Newt

Managing page assembly browser side is a chore. We can skip the complex JavaScript if we let Newt do the routing of requests to PostgREST and then sending those results through Pandoc like we did in our static site demo. With SQL (moding data), YAML (configuring Newt) and Pandoc templates we can build a web application.

- README.md, demo "read me" file
- birds.yaml, YAML configuration file for Newt (25 lines)
- birds.csv, our prior list of bird sightings to be loaded into SQL database
- htdocs, the website document directory
- page.tmpl, Pandoc template for listing birds  (13 lines)
- post_result.tmpl, Pandoc template for POST results (7 lines) from adding birds
- postgrest.conf, configuration for PostgREST
- setup.sql, SQL setting up data models and JSON data API (50 lines)

This is very similar to both demo 1 and 2. Missing is build.sh from demo 1. We don't need it since we're running Pandoc as a microservice.  There is an YAML configuration file for Newt. Two Pandoc templates performs a similar duty as the one used in birds 1 demo. Notice there is no sightings.js in our htdocs directory. From the web browsers point of view there is no need to run JavaScript to submit a standard web form to add a bird sighting.

NOTE: Newt provides our static file service so when developing we can skip Apache 2 or NginX. We would still need a front-end web service if you want to use it for access control or make the service outside of "localhost".

Pros                                   Cons
-------------------------------------  ---------------------------------------------
Simple back-end, just SQL code.        Like demo 2 you need to know some SQL.
No JavaScript required browser side.   Need to know how to configure Newt via YAML.
Newt provides the static web service
and data routing by a YAML file.


## Conclusion

In the birds 3 I've delegated tasks to a series of flexible microservices.

Postgres
: provides data storage and management

PostgREST
: Turns a Postgres database into a our JSON data API service

Pandoc
: Run in server mode is a powerful template engine, it can convert our JSON data into a web page

Newt
: Is a data router and static file server. It translates the web form submission into JSON before sending requests to PostgREST. Newt takes the results and sends that through Pandoc. Newt can also service static files. It could be used to talk to a JSON oriented full text search engine like Solr, Elasticsearch or Opensearch. The limit of services it based only on the number of routes you're willing to define.

Front-end web server
: A front-end web server can provide access control and proxy to any of our microservices, leverage virtual hosting, etc.

In birds 3, the "coding" of the back-end is reduced to SQL, Pandoc templates and a YAML file. You are free to make the front-end as simple or as complex as you like. The microservices and front-end web server effectively snap together like LEGO bricks.

