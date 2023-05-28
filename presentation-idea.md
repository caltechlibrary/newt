
# Building a web application using SQL, static HTML and JavaScript

Caltech Library recently started using Postgres for our SQL based applications.  I stumbled upon an application called PostgREST. That project has been active since at least 2014.  PostgREST provides a JSON data API (aka microservice) for Postgres databases.  It does this because PostgREST understands Postgres schema including tables, views, functions and procedures.  I used to write a fair amount of code in Go or Python to build a JSON data API wrapping a SQL database. I now use PostgREST and skip writing wrapper code completely.

In this presentation I will walk through the basics of building a static birds list applicaition first with Pandoc alone, then with Postgres+PostgREST and browser side JavaScript. Finally I will drop the browser side JavaScript in favor building the whole thing with Postgres+PostgREST, Pandoc and Newt.

This last approach fits well in the academic setting where we often build applications that integrate with single sign-on systems like Shibboleth via a front-end web service like Apache 2 or NginX. Postgres+PostgREST also plays nice with container style development. This is true too for Pandoc and Newt.

With just a little SQL, a little HTML you're on your way to building microservice based on web applications and services.
