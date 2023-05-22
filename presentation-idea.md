
# Building a web application using SQL, static HTML and JavaScript

Caltech Library recently started using Postgres for our SQL based applications.  I stumbled upon an application called PostgREST. PostgREST provides a JSON data API (aka microservice) for Postgres databases.  It does this because PostgREST understands Postgres schema including tables, views, functions and procedures.  I used to write a fair amount of code in Go or Python to build a JSON data API wrapping a SQL database. I now use PostgREST and skip writing wrapper code completely.

In this presentation I will introduce PostgREST through building an application using SQL, static HTML and browser side JavaScript. The infrastructure is provided via Postgres, PostgREST and a static web server of your choice.  This approach fits well in the academic setting where we often build applications that integrate with single sign-on systems like Shibboleth via a front-end web service like Apache 2 or NginX. Postgres+PostgREST also plays nice with container style development as well.

Included will be an exploration of how Postgres+PostgREST provides rich opportunities for "development on the edges" for existing applications built on Postgres such as Invenio RDM.

With just a little SQL, a little HTML and some JavaScript you're on your way to building microservice based on web applications and services.
