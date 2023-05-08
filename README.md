
# Newt Demo

Newt is my short for "new take" on building web applications from mostly off the shelf micro services. The micro services newt explores are [Postgres](https://postgresql.org), [PostgREST](https://postgrest.org), [Pandoc server](https://pandoc.org), a minimalist request router taking it's configuration from a SQL table and a front-end web server.

The goal of newt is to define a minimal set of micro services suitable for building applications used by Libraries and Archives. Ideally customization is limited to data modeling and url request routing via SQL and traditional front end development using HTML, CSS and JavaScript.


## Introducing Postgres and PostgREST

This first demo demonstrates building a traditional web application using
Postgres 15, PostgREST 11 and a front end web server. The demo focuses
on the development implecations and should be construed as what is
required in a public facing or production deployment.

- [Building with Postgres and PostgREST](building-with-postgres-postgrest.md) discusses the approach taken to create the [birds](birds/) demo
- Extas for setting up a developer environment
    - [setup-birds.bash](setup-birds.bash), a bash script that generates the contents of demo's running code
    - [Multipass basics](multipass-basics.md), multipass runs Ubuntu VM which can be used to run the demo
    - [newt-init.yaml](newst-init.yaml) provides the configuration for a multipass based VM to run the demo
    - [setup-developer-account.bash](setup-developer-account.bash) is a bash script that displays Postgres commands for setting up a super user account for development.

## Adding a template engine with Pandoc server

This second demo, [bees](bees/), builds on the birds demo by adding Pandoc server as a template engine for rendering content retrieved from Postgres via PostgREST. Finally page assembly is still done on the front end via JavaScript calls back into our micro services.

## Newt, a configuration driven router

There are times when relying on JavaScript in the web browser is not suitable or avaialble. How do we keep things simple without writing a bunch of middleware to service this use case?  The answer is by providing a configurable request router that knows how PostgREST and Pandoc server work. The configuration constists of a path pattern, an HTTP Method, content type requested, and possibly some data (e.g. POST from a form or file upload).  The [flower](flower/) demo shows how using a simple general purpose router can be used to assemple more complex web applications that used to require whole platforms like Wordpress or Drupal.

Newt is a configurable URL request router. At it's simplist is will service as a pass through to our JSON API, but more typically it will map a request to the JSON API, take its result and then run it through the Pandoc server template engine returning the assembled result to our front end web server.

Another trick of Newt is that it can take a urlencoded POST (like from a traditional web form) and translate it into a JSON API call and processing that result via Pandoc server before returning the result to our front end web server. Likewise it also knows how to handle FILE uploads if that is desired.

Newt's behavior is defined in a single route table in our SQL database. Newt queries PandgREST when it startes up for it's configuration, where to find the Pandoc templates and what routing is desired.  That is kept in memory to it remain fast to responsive even on minimal hardware. 

## Conclusion

In all our demos we've devided the tasks of building our website via a set of well defined but configurable services. Postgres provides all our metadata services, PostgREST takes care mapping our JSON API requests to/from Postgres, Pandoc server provides a robust in-memory template engine and Newt provides a configurable URL request router that is PostgREST and Pandoc server aware. Access and controll is defered to either our front-end web server (e.g. via single signon like Shibboleth) or to rules execute via stored procedures and triggers in Postgres.  Building our web site then is reduce to simple configurable applicaitions which are customized via SQL on the back-end and the venerable HTML, CSS and if needed JavaScript in the web browser.






