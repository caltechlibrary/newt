---
title: "Building with PostgreSQL and PostgREST"
pubDate: 2023-06-01
author: "R. S. Doiel"
---

# Building with PostgreSQL and PostgREST

I would like to suggest a simplification in building web applications. Can you build a complete web application using only SQL, HTML and browser side JavaScript? Can you do it with off the shelf services like Postgres and PostgREST? Our application uses a microservice model. Data is modeled and managed via a Postgres 15 database. The database content is managed through a JSON data API provided by PostgREST 11. Assuming you have a static file server handy you can do the page assembly browser side. In the scenario the index.html and sightings.js files access the JSON data API microservice providing a human friendly view of our modeled data.

How does this approach simplify creating and managing web applications?

The back-end becomes largely declarative broken down into three simple tasks.

1. define your data schema in SQL
2. Populate your database tables with some working content
3. Setup PostgREST access to manage the data base contents

While two different programs (i.e. Postgres and PostgREST) they present a unified view of your data via a predictable JSON API. Running a web server in front of Postgres+PostgREST as a proxy will let you build the rest of your application in the front-end using HTML and JavaScript running in the web browser.  The back-end can be debugged by either interacting directly with Postgres or by using curl to test PostgREST responses.

# A bit more about PostgREST

PostgREST provides a generic middle ware that transforms a Postgres database into a RESTful JSON data API that can be OpenAPI compliant if desired. PostgREST knows the schema for the database because it knows how Postgres works. It provides simplified routes (URL paths) for interacting with a database, queries, views and stored procedures via standard HTTP methods (e.g. GET, POST). If these are proxied via your web server then you can use them to allow the front-end running in the web browser to manage display. The back-end code you right is only what you need to model the data in Postgres and facilitated access via PostgREST.  You make more elaborate end points available in PostgREST via stored functions and views but for many simple applications what is provided out of the box by PostgREST may be enough.

# How much SQL do I need to know to use Postgres with PostgREST?

The bare minimum would be to now how to create a database, create a table and the basic SELECT statement.  In most read only style applications that's all you really need to use.  Have an understanding of Postgres SCHEMA[2] will make understanding how PostgREST works easier.

A more rounded minimum SQL knowledge would include being familiar with DROP, INSERT, UPDATE, and DELETE statements. It is really helpful to  have a basic understanding of VIEWS[3]. An understanding of stored procedures and functions[4] allows you to take full advantage of PostgREST's capabilities.


[2]: Postgres [SCHEMA](https://www.postgresql.org/docs/current/ddl-schemas.html)

[3]: [SQL views](https://en.wikipedia.org/wiki/View_(SQL)), Postgres [specifics](https://www.postgresql.org/docs/current/sql-createview.html)

[4]: [Stored procedures concept](https:https://en.wikipedia.org/wiki/Stored_procedure), Postgres [CREATE PROCEDURE](https://www.postgresql.org/docs/current/sql-createprocedure.html) and [Procedural Language](https://www.postgresql.org/docs/current/plpgsql.html)


## Dividing up application responsibilities

The architecture is traditional. It is normally called a [microservice](https://en.wikipedia.org/wiki/Microservices) architecture. The novel feature presented here is that the services are completely off the shelf.  The only back-end programming to be done is setting up the Postgres database using SQL. The rest is delegated to the web browser using the traditional troika of HTML, CSS and JavaScript.

1. Two microservices provide data manage and transformation
  - PostgreSQL, data modeling, data and role management
  - PostgREST, turning what we have in Postgres into a JSON data API
2. The front end web server and web browser provides our human facing interface through static HTML and JavaScript

## Stepping through Birds 2 Demo

The birds two demo consists of the following files.

- [README.md](birds2/README.md)
- [birds.csv](birds2/birds.csv)
- [htdocs](birds2/htdocs)
    - [index.html](birds2/htdocs/index.html)
    - [sightings.js](birds2/htdocs/sightings.js)
- [postgrest.conf](birds2/postgrest.conf)
- [setup.sql](birds2/setup.sql)
~~~

Basic development cycle can be broken down into the following steps.

1. Create a database ans setup the Postgres SCHEMA (namespace) which will be used to allow PostgREST to access our data (SQL)
2. Creating our database model(s)  (SQL)
3. Loading data from a CSV file (also done via SQL)
4. Startup PostgREST and confirm we can access the content
5. Create and debug our index.html/sightings.js files in the web browser

As you evolve the application you'll likely repeat steps two through four.


## Software Requirements

1. PostgreSQL 15 and PostgREST 11 need to be setup and configured to work together
2. You need a simple web server (e.g. `python3 -m http.server`) for serving static content
3. You need a web browser (e.g. Safari, Firefox, Chrome, Edge)
4. You need your favorite text editor (e.g. micro, vim, Emacs, VSCode, Zed)

## Closing thoughts

Birds 2 demo is a proof of concept showing much simplified back-end development process. Namely define your schema and let Postgres/PostgREST do the rest. It does come at the expense of making the front-end work running the web browser more complicated. It remains to be seen if this is worth the trade off.


