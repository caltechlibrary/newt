
# Newt, a new take of the web stack

Newt is short for "new take". It is my personal new take on building web applications from a few off the shelf [microservices](https://en.wikipedia.org/wiki/Microservices). The microservices newt explores are [Postgres](https://postgresql.org), [PostgREST](https://postgrest.org), [Pandoc server](https://pandoc.org). The Newt program is a minimalist URL request router that services as a bridge between the other microservices and your frontend web server.

My belief is that many web services used by Archives, Libraries and Museums can benefit from a simplier back-end. If the back-end is "easy" then the limited developer resources can be focused on the front end creating a better user experiences for staff and patrons.

The Newt's vision is the back-end can be defined simply through data models in a SQL database like Postgres.  The models can be exposed as a microservice using PostgREST via via SQL queries, views, functions and procedure  that PostgREST can expose as a JSON data API. If you need some sort of server side template engine for content provided by PostgREST then Pandoc server is available and sorts many different types of content transformation as well as a template engine. Finally venerable static files can be served by either your front end web server or a minimalist URL router like that implemetned at part of Newt.

I feel if you know SQL, HTML, and perhaps CSS that should be enough to build useful web applications and services without having to know languages like Python, Perl, PHP, Ruby, Go or Java. While I will discuss using browser side JavaScript in the initial exploration of the micro serverice mention above by the time we get to the final demo which features the Newt URL router you'll need to use SQL, HTML and some CSS and you'll have a web application that works equally well in modern browers like Firefox, Safari and Chrome as well as with text mode browsers like Lynx.


## Introducing Postgres and PostgREST

To understand the Newt it is first helpful to understand Postgres and PostgREST. The first demo, [birds](birds/), builds a simple bird spotting web site. It uses PostgREST as a JSON data API and front end web server (specifically `python3 -m http http.server`) to present static web pages that interact with the PostgREST service. The purpose of the demo is to show how to develop using Postgres 15 and PostgREST 11 and static files.

- [Building with Postgres and PostgREST](building-with-postgres-postgrest.md) discusses the approach taken to create the [birds](birds/) demo
- Extas for setting up a developer environment
    - [setup-birds.bash](setup-birds.bash), a bash script that generates the contents of demo's running code
    - [Multipass basics](multipass-basics.md), multipass runs Ubuntu VM which can be used to run the demo
    - [newt-init.yaml](newt-init.yaml) provides the configuration for a multipass based VM to run the demo
    - [setup-developer-account.bash](setup-developer-account.bash) is a bash script that displays Postgres commands for setting up a super user account for development.

## Introducing Pandoc server, a powerful template engine

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

