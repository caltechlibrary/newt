
# Newt, a new take of the web stack

Newt is short for "new take". Specific it is my new take on building web applications from a few off the shelf micro services. The micro services newt explores are [Postgres](https://postgresql.org), [PostgREST](https://postgrest.org), [Pandoc server](https://pandoc.org), Newt a minimalist URL request router.

My belief is that many web services and sites used by Archives, Libraries and Museums can benefit from a simplier back-end that allows the application developer to focus on the front-end that impacts both our staff and patrons.

The Newt's vision is to define data models and route handling via SQL queries, views, functions and procedure (via Postgres+PostgREST), leverage Pandoc server for templating needs and a light weigth URL router that can use both PostgREST and Pandoc to return dynmamic content or static content from the local file system. In this way both developer and technical staff and develop interactive websites and services without having to know more than SQL, HTML and a little JavaScript.


## Introducing Postgres and PostgREST

To understand Newt you need to first understand how much you can do with just Postgres and PostgREST. The demo [birds](birds/) builds a simple bird spotting web site. It uses PostgREST as a JSON data API and front end web server (specifically `python3 -m http http.server`) to present static web pages that interact with the PostgREST service. The purpose of the demo is to show how to develop using Postgres 15 and PostgREST 11 and static files.

- [Building with Postgres and PostgREST](building-with-postgres-postgrest.md) discusses the approach taken to create the [birds](birds/) demo
- Extas for setting up a developer environment
    - [setup-birds.bash](setup-birds.bash), a bash script that generates the contents of demo's running code
    - [Multipass basics](multipass-basics.md), multipass runs Ubuntu VM which can be used to run the demo
    - [newt-init.yaml](newt-init.yaml) provides the configuration for a multipass based VM to run the demo
    - [setup-developer-account.bash](setup-developer-account.bash) is a bash script that displays Postgres commands for setting up a super user account for development.

## Pandoc server, a powerful template engine

This second demo, [bees](bees/), builds on the birds demo by adding Pandoc server as a template engine for rendering content retrieved from Postgres via PostgREST. Pandoc service provide much the capability of the Pandoc command line program to a simple and safe localhost web service.  In bees we are presenting a bee spotting application with reports rendered via data from POstgREST and processed by Pandoc server.

- [Building with Postgres, PostgREST and Pandoc server](building-with-postgres-postgrest-and-pandoc-server.md)
- Extras for setting up demo (assumes you've already setup birds)
    - [setup-bees.bash](setup-bees.bash), a bash script that generates the contents of demo's running code

## Introducing Newt, a minimalist URL router

There are times when relying on JavaScript in the web browser to assemble a web page is not appropriate. In development it is also a hassle to run a full blown front end web server like Apache 2 or NginX. The Newt URL router tries to remedy this by allowing you to programtically define the routes your website will support what accomodating web page assembly (via PostgREST and Pandoc server) and static file services for content. In this penultimate demo called [flowers](flowers/). A service that tracks flowers. It is build via SQL, Pandoc templates and static files. Unlike our previous two examples the web pages is assembled server side so no JavaScript is required in the web browser.  While Newt doesn't current support file uploads it does support building dynamic websites and services where the data is held in Postgres.

- [Building web Postgres, PostgREST, Pandoc and Newt](building with postgres-postgrest-pandoc-and-newt.md)
- Extras for setting up demo (assumes you've already setup birds and bees)
    - [setup-flowers.bash](setup-flowers.bash), a bash script that generates the contents of demo's running code

In this final demo, now including the Newt URL router, we use SQL to define our back-end service, Pandoc (via Pandoc server) to format our dynamic web poages and static files to complete our website implementation.

## Conclusion

In all our demos we've devided the tasks through a series of flexible micro services can can be used individually or in combination.

- Postgres, data storage, back-end configuration and data services
- PostgREST, a JSON data API
- Pandoc server as a template engine
- Newt a URL router for PostgREST, Pandoc server and static files
- (in production setting) a front end web server providing access control and user authentication (e.g. Apache2, NginX)

The "coding" left to someone developing a website or service can be as minimal as knowning some SQL, HTML and Pandoc.

