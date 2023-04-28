---
title: "Building with PostgreSQL and PostgREST"
pubDate: 2023-04-27
author: "R. S. Doiel"
---

# Building with PostgeSQL and PostgREST

I would like to suggest a simplification in building web applications. I will show you how to create a complete web application using SQL, HTML and JavaScript. The application uses a micro service model. Data is modeled and access via a Postgres 15 database. The data base content is managed through a RESTful JSON API implemented via PostgREST 11. A web server is used to serve static content, e.g. HTML and JavaScript that makes calls to the JSON API providing our web application. This demo focuses on the impact in the development process not how you would configure a similarly structure application in a public or production settings.  PostgREST and our web server run on localhost for the purposes of the demonstration.

Who does this simplify creating and managing the back end?

PostgREST provides a generic middle layer. It knows the schema for the database because it knows how Postgres works. Routes with PostgREST are formed by the database name and additional modifiers the either represent an SQL expression or correspond to a SQL defined view, stored procedure or prepared statement. What this means is that the SQL you write in Postgres determines the how the RESTful JSON API is configured by PostgREST. 

What is the minumum SQL I need to know to use Postgres with PostgREST?  A read only service might only need some CREATE statements to define our data model(s) and some SELECT statements to retreive results. A more like scenerio is you'll also want to know how to write INSERT, UPDATE and DELETE statements. If you want to most readable URLs in your JSON API it's nice to know how to write an [SQL view](https://en.wikipedia.org/wiki/View_(SQL)) and [stored procedures](https:https://en.wikipedia.org/wiki/Stored_procedure).  

This should make more sense through demonstrating writing a simple web application that implements a bird list.  Aside from configuring everything the source code I write will be either SQL, HTML or JavaScript.

## Dividing up application responsibilities

The architecture I am suggesting is traditional. The novel feature is relying only on off the shelf middle layer, PostgREST, that is defined and managed directly in our Postgres database.

1. Two micro services provide data manage and transformation
    a. PostreSQL, data modeling, management and transformation
    b. PostgREST turning what we have in Postgres into a JSON API
2. The front end web server provides our human facing interface defined in HTML and JavaScript
3. With the static web server and our JSON API running on localhost our web browser knitts the two together giving us our application

## How do I suggest proceding with development?

The Postgres client [psql](https://www.postgresql.org/docs/current/app-psql.html) provides a rich development environment for interacting with Postgres.  It includes support for access your favorite text editor, load data and executing SQL programs or scripts. It also lets you take definitions you've written interactively and save them to files for reuse later (or debuggin).

Outside of `psql` you'll be using familiar tools like your text editor and web browser. `psql` can also be easily used by your shell or other scripting languages if you want to automate things that way.

## Birds, an application 

Our demonstration application presents and manages a list of birds and their sightings.  It includes being able to list the birds, add a new sighting and search for a specific sighting.

The steps to set things up are

1. Creating a project directory and some empty files
2. Setting up a database called "birds" for our application
3. Modeling our data in a table called "sighting", define some user roles and grant permissions, define a view and a stored procedure
4. Configuring our PostgREST and insuring access works
5. Writing some HTML and JavaScript to take the JSON responses and turn them into something that is human friendly

This is the recipe I used to explore both Postgres and PostgREST.

### Prequisites

I assume you have Postgres, PostgREST installed and available on your development machine. I also assuming you have a recent Python 3 available[2]. I will be using that as my simple web server.

### Steps 1 and 2

After comfirming that Postgres is installed and I have access as a super user from my account I created a directory for the project. I'm going to use two folder `birds` and `birds/htdocs`. I the first directory I'll place all my SQL code as well as my configuration for PostgREST.  In the second directory, `birds/htdocs` I will place my "index.html" file and "birds.js". The later is how I access the JSON API provided by PostgREST.

~~~
mkdir -p birds/htdocs
cd birds
touch setup.sql
touch birds.csv
touch postgrest.cnf
touch htdocs/index.html
touch htdocs/birds.js
createdb birds
~~~

### Step 3

This next order of business is to use the `psql` client and our text editor to create "setup.sql". In this file we need to setup a schema, define out "sighting" table. Setup some roles and grant some permissions. Add a view and stored procedure.

Type `psql -d birds` and press enter. You should now be in the postgres client.  Fromt client we can do lots of things but right now we want to fill in our "setup.sql" with the actions a outlined above. To edit "setup.sql" type
`\e setup.sql`. You should not be in your favorite text editor[1].

[1]: If for some reason the editor isn't set, try setting `PSQL_EDITOR` in your environment before start `psql`.

The completed [setup.sql](setup.sql) should look something like this.

~~~
DROP SCHEMA IF EXISTS birds CASCADE;
CREATE SCHEMA birds;

--
-- This file contains the create statements for our bird table.
--
DROP TABLE IF EXISTS birds.sighting;
CREATE TABLE birds.sighting
(
  bird_name VARCHAR(255),
  place TEXT,
  sighted DATE
);

-- bird_view will become an end point in PostgREST
CREATE VIEW bird_view (bird_name, place, sighted) AS 
  SELECT bird_name, place, sighted FROM birds.sighting ORDER BY sighted, bird_name;

-- record_bird is a stored procedure and will save a new bird sighting
CREATE PROCEDURE record_bird(name VARCHAR(256), description TEXT, dt DATE)
LANGUAGE SQL
AS $$
  INSERT INTO birds.sighting (bird_name, place, sighted) VALUES (name, description, dt);
$$;

--
-- The following additional steps are needed for PostgREST access.
--
DROP ROLE IF EXISTS birds_anonymous;
CREATE ROLE birds_anonymous nologin;

GRANT USAGE ON SCHEMA birds TO birds_anonymous;
GRANT SELECT ON birds.sighting TO birds_anonymous;

DROP ROLE IF EXISTS birds;
CREATE ROLE birds NOINHERIT LOGIN PASSWORD 'replace_me';
GRANT birds_anonymous TO birds;
~~~

After finishing typing this in save it and exit your editor. You should be back at the `psql` prompt.  Type `\i setup.sql` and press enter. This command reads the updated "setup.sql" and executes the instructions it finds. If you get erorrs try editing the file again correcting the erorrs, reload and test.

Next we want to create some test data in a CSV file and load that into
our database. We can do that with our text editor once again. This time
type `\e birds.csv`.

Fill in "birds.csv" as follows.

~~~
bird,place,sighted
robin, seen in my backyard,2023-04-16
humming bird, seen in my backyard, 2023-02-28
blue jay, seen on my back porch, 2023-01-12
~~~

Save and exit your editor. To take this CSV file and populate the "sighting" table run this following command at the `psql` prompt.

~~~
\copy birds.sighting from 'birds.csv' with (format csv, header true);
~~~

That should load three rows of data into our "sighting" table. you can check that with the following SQL statement.

~~~
SELECT * FROM birds.sighting;
~~~

We're now ready to check to see if our view works.

~~~
SELECT * FROM bird_view;
CALL record_bird('owl', 'my front porch', current_date); 
SELECT * FROM bird_view;
~~~

If you've made it this far exist `psql` with `\q`.

### Step 4

We previously setup Postgres to be ready to work with PostgREST. Now we
need a configuration file named "postgrest.conf". Using your text editor
this is what the configuration file should contain.

~~~
db-uri = "postgres://birds:replace_me@localhost:5432/postgres"
db-schemas = "birds"
db-anon-role = "birds_anonymous"
~~~

We are now ready run and test our PostgREST.

~~~
postgrest postgrest.conf
~~~

At this point you should see a response in your terminal something like

~~~
27/Apr/2023:19:36:51 -0700: Attempting to connect to the database...
27/Apr/2023:19:36:51 -0700: Connection successful
27/Apr/2023:19:36:51 -0700: Listening on port 3000
27/Apr/2023:19:36:51 -0700: Listening for notifications on the pgrst channel
27/Apr/2023:19:36:51 -0700: Config reloaded
27/Apr/2023:19:36:51 -0700: Schema cache loaded
~~~

Point your web browser at <http://localhost:3000/sighting> and see if you get back a JSON response like

~~~
[{"bird_name":"robin","place":" seen in my backyard","sighted":"2023-04-16"},
 {"bird_name":"humming bird","place":" seen in my backyard","sighted":"2023-02-28"},
 {"bird_name":"blue jay","place":" seen on my back porch","sighted":"2023-01-12"},
 {"bird_name":"owl","place":"front porch","sighted":"2023-04-27"}]
~~~

If you have gotten this far you have a working RESTful JSON API.

### Step 5

In this section we'll need our text editor to create `htdocs/index.html`
and `htdocs/birds.js`.  We will make server these using the Python 3
command `python3 -m http.server` from within the `htdocs` directory.

First open a new terminal window or shell session and change into our htdocs directory and start the python web server.

~~~
cd htdocs
python3 -m http.server
~~~

Now edit the file "index.html" to look like this.

~~~
<DOCTYPE html>
<html>
<body>
<h1>Welcome to the bird list!</h1>
<div id="bird-list"></div>
<script src="birds.js"></script>
</body>
</html>
~~~

We now need to populate our "bird-list" div element using JavaScript.

We need to take the following steps.

1. Get a handle on the "bird-list" div
2. Contact the JSON API to get back our list of birds
3. When we have read it back generate HTML to insert inside "bird-list"

~~~
/* FIXME: write example JavaScript here. */
~~~

