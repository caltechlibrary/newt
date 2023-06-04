
# Newt, a new take on the web stack

Newt is two things. Newt is a "new take" on building web applications. Newt is an experimental [microservice](https://en.wikipedia.org/wiki/Microservices) working with other "off the shelf" microservices.

Newt integrates [Postgres](https://postgresql.org), [PostgREST](https://postgrest.org) and [Pandoc](https://pandoc.org) by function as a data router and a light weight static file server.

Newt routes a request to a JSON data API (e.g. PostgREST, Solr, Elasticsearch) and then optionally send that result through Pandoc for further processing. Newt runs as a localhost service. Inm production you'd use Newt behind a traditional web server like Apache 2 or NginX.

## Motivation

My belief is that many web services used by Archives, Libraries and Museums can benefit from a simplified and consistent back-end. If the back-end is "easy" then the limited developer resources can be focused on the front-end. **An improved front-end offers opportunities to provide a more humane user experience for staff and patrons.** 

## Newt's approach

From front-end to back-end

- A front end web server (e.g. Apache 2, NginX) can provide access control where appropriate (e.g. Single Sign on with OAuth2 or Shibboleth)
- Newt provides static file services but more importantly serves as a data router. I can map a request to a JSON data API, take those results then send them through Pandoc for transformation.
- Postgres+PostgREST is an example of a JSON data API
- Solr, Elasticsearch or Opensearch can also function as a JSON data API
- Pandoc server provides a templating engine to transform data sources

All these can be treated as "off the shelf". Aside from configuration they can run as traditional services on most POSIX systems.  Your application is implemented using SQL. It is enhanced by using Pandoc templates used to turn JSON into HTML (or other desired formats). 

## Exploring Newt's friends

In the repository three [demos](https://github.com/caltechlibrary/newt/tree/main/demos) are provided to show how Pandoc, Postgres+PostgRESt and Newt can work together. The demos are generated using the a simple Bash scripts (see below). The bash scripts will generate all the file needed to run the demos. A read me for each demo describes how to run it.

[Birds 1 Demo](https://github.com/caltechlibrary/newt/blob/main/demos/setup-birds1.bash)
: Shows a simple use of Pandoc to render a static bird sightings website

[Birds 2 Demo](https://github.com/caltechlibrary/newt/blob/main/demos/setup-birds2.bash)
: Shows a dynamic bird sightings website using Postgres+PostgREST, but requires JavaScript running in your web browser

[Birds 3 Demo](https://github.com/caltechlibrary/newt/blob/main/demos/setup-birds3.bash)
: Shows a dynamic bird sightings website using Newt with Postgres+PostgREST and Pandoc

The goal of the three demos is to show an evolution towards simplicity.


### Birds 1 Demo, Pandoc only

This is a simple static website. It introduces Pandoc and which we will leverage again in bird 3 demo. It also is a good tool to rendering static content with. Static websites are generally simple to generate and maintain.

- README.md, demo read me
- birds.csv, our list of bird sightings
- build.sh, a shell script that uses Pandoc to render site (3 lines)
- htdocs, the website document directory
- page.tmpl, a Pandoc page template using by build.sh (13 lines)

Pros                             Cons
-------------------------------- -----------------------------------
Pandoc is widely known           Updates require changing birds.csv
Pandoc is easily scripted        and rebuilding the site.


### Birds 2 Demo, Postgres + PostgREST

Postgres and PostgREST can be combined to provide a dynamic data source for our birds 2 demo website. In birds 2 demo PostgREST runs on localhost and we use another web server to presents the static files and to proxy to PostgREST. The enhancements over  birds 1 demo is our list of birds is held in a Postgres database which is available to our web browser thanks to proxying to PostgREST. The web browser uses JavaScript to call back to the JSON data API and using the results to render content in the web page.

The back-end is written using SQL. This includes setting up access by PostgREST. The font-end JavaScript is required in the browser to assemble pages results from our JSON data API is complex. The complexity is shifted out of the back-end which is using an off the shelf micro service to the front-end.

- README.md, demo read me
- birds.csv, our sightings data, loaded into SQL database
- htdocs, the website document directory
- htdocs/sightings.js, JavaScript support is required in the web browser to render the page and handle updates (63 lines)
- postgrest.conf, configuration for PostgREST (3 lines)
- setup.sql, SQL setting up data models and JSON data API (50 lines)

Pros                              Cons
--------------------------------  ------------------------------------
Simple back-end, just SQL code    JavaScript in the browser is complex
                                  You need a static file server
                                  You need to know some SQL


### Birds 3 demo, assembling responses with Newt

Managing page assembly browser side is a chore. We can skip the complex JavaScript browser side if we let Newt do the routing of requests to PostgREST and then sending those results through Pandoc like we did in our static site demo. With just SQL and Pandoc templates we can build a web application.

Newt needs to know how to map the front-end requests to PostgREST so to do that it reads a CSV file holding data routing instructions. E.g. turn a request URL into a PostgREST URL and run the results through Pandoc running as a microservice.

- README.md, demo read me
- birds-routes.csv, CSV holding data routing instructions (3 lines)
- birds.csv, CSV holding data to be loaded into our SQL database
- birds.yaml, configuration for Newt
- htdocs, holds our static content
- page.tmpl, Pandoc template for listing birds  (13 lines)
- postgrest.conf, configuration for PostgREST
- setup.sql, SQL setting of data models and JSON data API (50 lines)

This is very similar to both demo 1 and 2. Missing is build.sh from demo 1. We don't need it since we're running Pandoc as a microservice.  There is an added configuration file for Newt. The Pandoc template performs a similar duty as the one used in birds 1 demo. Notice there is no sightings.js in our htdocs directory. From the web browsers point of view there is no need to run JavaScript to submit a standard web form to add a bird sighting.

NOTE: Newt provides our static file service so when developing we can skip Apache 2 or NginX.

Pros                                 Cons
-----------------------------------  ------------------------------------
Simple back-end, just SQL code       Like demo 2 you need to know some SQL
No JavaScript required browser side
Newt provides the static web server


## Conclusion

In the birds 3 demo I've delegated tasks to a series of flexible microservices.

Postgres
: provides data storage and defines how our JSON data API works

PostgREST
: Turns a Postgres database into a our JSON data API service

Pandoc
: Run in server mode is a powerful template engine, it can convert our JSON data into a web page

Newt
: Is a data router and static file server. It translates the web form submission into JSON before sending requests to PostgREST. Newt takes the results and sends that through Pandoc. Newt can also service static files. It could be used to talk to a JSON oriented full text search engine like Solr, Elasticsearch or Opensearch.

Front-end web server
: A front-end web server can provide access control and proxy to any of our microservices, leverage virtual hosting, etc.

The "coding" of the back-end is reduced to SQL and Pandoc templates. You are free to make the front-end as simple or as complex as you like. The microservices and front-end web server effectively snap together like LEGO bricks.



