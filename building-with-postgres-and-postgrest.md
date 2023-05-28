---
title: "Building with PostgreSQL and PostgREST"
pubDate: 2023-04-27
author: "R. S. Doiel"
---

# Building with PostgreSQL and PostgREST

I would like to suggest a simplification in building web applications. I will show you how to create a complete web application using SQL, HTML and JavaScript. The application uses a microservice model. Data is modeled and access via a Postgres 15 database. The database content is managed through a JSON data API provided by PostgREST 11. I will use Python 3's http.server to service my static content (index.html and sightings.js). The index.html and sightings.js files access the JSON data API microservice to provide a human friendly view of our application data.

> This demo focuses on the impact in the development process not how you would configure a similarly application in a public or production settings.  As a result PostgREST and our web server run on localhost.

How does this approach simplify creating and managing web applications?

I believe this approach simplifies the number of layers that need to be managed by the developer. It provides clear lines dividing responsibilities of the back end (PostgreSQL and PostgREST) and the front end (the web server in combination with the web browser). It reduces the required number of language context switches to primarily SQL and JavaScript. It provides the option of integration with other layers (e.g. single sign on if you were using a real web server like NginX or Apache 2) and scaling wide via configuring Postgres in a HA setup along with load balancing. Scaling and adding access controls can happen without changing the SQL code modeling data or application behaviors. These enhancement will not be covered in this presentation and are left as an exercise for the reader.


PostgREST provides a generic middle ware that transforms a Postgres database into a RESTful JSON data API that can be OpenAPI compliant if desired. PostgREST knows the schema for the database because it knows how Postgres works. It provides simplified routes (URL paths) for interacting with a database, queries, views and stored procedures. That later two will help keep the URLs simple when used by your front end. Aside from configuration the database and PostgREST you only need to write SQL to develop the back end.


How much SQL do I need to know to use Postgres with PostgREST?

The bare minimum would be for a read only style web application. You would need to write SOME CREATE statements to model the data and understand the SELECT statement.  An improvement would be to understand SQL VIEWS. With SQL VIEWS you can organize more complicated SELECT statements into a single name that becomes an end point in your JSON data API. The demo will use this approach.

A more realistic minimum knowledge for a read write application would be able to write CREATE, DROP, INSERT, UPDATE and DELETE statements. Have an understanding of Postgres VIEWS[2] and SCHEMA[3] (helpful in debugging permission issues) and a basic understanding of stored procedures[4].

[2]: [SQL views](https://en.wikipedia.org/wiki/View_(SQL)), Postgres [specifics](https://www.postgresql.org/docs/current/sql-createview.html)

[3]: Postgres [SCHEMA](https://www.postgresql.org/docs/current/ddl-schemas.html)

[4]: [Stored procedures concept](https:https://en.wikipedia.org/wiki/Stored_procedure), Postgres [CREATE PROCEDURE](https://www.postgresql.org/docs/current/sql-createprocedure.html) and [Procedural Language](https://www.postgresql.org/docs/current/plpgsql.html)


## Dividing up application responsibilities

The architecture I am suggesting is traditional. The novel feature is relying only on off the shelf middle ware implemented by PostgREST. How the JSON data API works is defined and managed directly in our Postgres database after configuring PostgREST to access our Postgres database.

1. Two microservices provide data manage and transformation
  - PostgreSQL, data modeling, data and role management
  - PostgREST, turning what we have in Postgres into a JSON data API
2. The front end web server and web browser provides our human facing interface through static HTML and JavaScript

## Requirements

1. PostgreSQL 15 and PostgREST 11 need to be setup and configured to work together
2. You need a simple web server (e.g. `python3 -m http.server`) for serving static content
3. You need a web browser (e.g. Safari, Firefox, Chrome, Edge)
4. You need your favorite text editor (e.g. micro, vim, Emacs, VSCode, Zed)

## A toolkit for developing with SQL

The Postgres client [psql](https://www.postgresql.org/docs/current/app-psql.html) provides a rich development environment for interacting with Postgres.  It includes support for using your favorite text editor, loading data and executing SQL programs or scripts. The client supports saving interactive work into text files for reuse and debugging.

Outside of `psql` you'll be using familiar tools like your text editor, shell commands and your web browser or [curl](https://en.wikipedia.org/wiki/CURL).

## Birds, an application

Our demonstration application presents and manages a list of birds sightings.  It includes listing the birds sighted and the availability to add to the list. It's a demo so anyone can add bird sightings if they access to localhost.

The steps covered will be

1. Creating a project directory and some empty files
2. Setting up a database called "birds" for our application
3. Modeling our data in a table called "sighting", define some roles and granting permissions, defining a view and stored procedure
4. (re)Configure and (re)start PostgREST
5. Writing some HTML and JavaScript to take the JSON responses and turn them into something that is human friendly

Steps three through five are repeated as necessary to evolve or debug our application.  This is the recipe I used to explore both PostgreSQL and PostgREST when preparing this presentation.

### Prerequisites

I assume you have Postgres, PostgREST installed and available on your development machine. I also assuming you have a recent Python 3 available.  Python provides a simple web server.

Confirm that Postgres is installed, PostgREST is installed and you have access to Postgres as a super user in the `psql` client.

### Step 1

Create a directory for our project. I'm going to use two folder `birds` and `birds/htdocs`. In the first directory I'll place all my SQL code as well as my configuration for PostgREST.  In the second directory, `birds/htdocs` I will place my "index.html" file and "sightings.js".

~~~
mkdir -p birds/htdocs
cd birds
touch setup.sql
touch birds.csv
touch postgrest.cnf
touch htdocs/index.html
touch htdocs/sightings.js
~~~

### Step 2

Create our "birds" database using the command line program `createdb` that comes with Postgres.

~~~
createdb birds
~~~

If you need to remove a stale version of the database you can use the command

~~~
dropdb birds
~~~

To drop a database you need to make sure you are not connected to it (e.g. check any running `psql` session and make sure it is NOT connected to the birds database).

### Step 3

> NOTE: steps 3 through 5 become our "apply shampoo, lather, rinse and repeat" steps for development

This next order of business is to use the `psql` client and our text editor to fill in "setup.sql". In this file we need to define our schema

- model our data, i.e. define our "sighting" table
- setup some roles and grant some permissions
- Add a view and stored function.

Type `psql -d birds` and press enter. You should now be in the Postgres client connected to the "birds" database.  From client we can do lots of things but right now we want to fill in our "setup.sql" with the actions a outlined above. To edit "setup.sql" type `\e setup.sql`. You should not be in your favorite text editor[5].

[5]: If for some reason the editor isn't set see the [psql](https://www.postgresql.org/docs/current/app-psql.html) web page, look for `PSQL_EDITOR` and follow the instructions.

The completed [setup.sql](birds/setup.sql) should look something like this.

> NOTE: I've included some DROP statements so you can re-run this, debug or start over if you wish to experiment.

~~~
-- Make sure we're in the birds database
\c birds

DROP SCHEMA IF EXISTS birds CASCADE;
CREATE SCHEMA birds;
SET search_path TO birds, public;

--
-- This file contains the create statements for our bird table.
--
-- DROP TABLE IF EXISTS birds.sighting;
CREATE TABLE birds.sighting
(
  bird_name VARCHAR(255),
  place TEXT,
  sighted DATE
);

-- bird_view will become an end point in PostgREST
CREATE OR REPLACE VIEW birds.bird_view AS
  SELECT bird_name AS bird, trim(place) AS place, sighted
  FROM birds.sighting ORDER BY sighted ASC, bird_name ASC;

-- record_bird is a stored procedure and will save a new bird sighting
CREATE OR REPLACE FUNCTION birds.record_bird(bird VARCHAR, place TEXT, sighted DATE)
RETURNS bool LANGUAGE SQL AS $$
  INSERT INTO birds.sighting (bird_name, place, sighted) VALUES (bird, place, sighted);
  SELECT true;
$$;

--
-- The following additional steps are needed for PostgREST access.
--
DROP ROLE IF EXISTS birds_anonymous;
CREATE ROLE birds_anonymous nologin;

GRANT USAGE  ON SCHEMA birds      TO birds_anonymous;
-- NOTE: We're allowing insert because this is a demo and we're not
-- implementing a login requirement!!!!
GRANT SELECT, INSERT ON birds.sighting    TO birds_anonymous;
GRANT SELECT ON birds.bird_view   TO birds_anonymous;
GRANT EXECUTE ON FUNCTION birds.record_bird TO birds_anonymous;

DROP ROLE IF EXISTS birds;
CREATE ROLE birds NOINHERIT LOGIN PASSWORD 'my_secret_password';
GRANT birds_anonymous TO birds;
~~~

> NOTE: you should replace "my_secret_password" with an actual password

After finishing typing save the file and exit your editor. You should be back at the `psql` prompt.  Type `\i setup.sql` and press enter. This command reads the updated "setup.sql" and executes the instructions it finds. If you get errors try editing the file again correcting the errors, reload and test.

Next we want to create some test data in a CSV file and load that into
our database. We can do that with our text editor. This time
type `\e birds.csv`.

Fill in "birds.csv" as follows.

~~~
bird,place,sighted
robin, seen in my backyard,2023-04-16
humming bird, seen in my backyard, 2023-02-28
blue jay, seen on my back porch, 2023-01-12
~~~

Save and exit your editor. To take this CSV file and populate the "sighting" table run this following command at the `psql` prompt

> NOTE: `\copy` and not `COPY`. This is necessary so the file is read client side.

~~~
\copy birds.sighting from 'birds.csv' with (format csv, header true);
~~~

That should load three rows of data into our "sighting" table. you can check that with the following SQL statement.

~~~
SELECT * FROM birds.sighting;
~~~

We're now ready to test our view and stored function.

~~~
SELECT * FROM birds.bird_view;
SELECT birds.record_bird('owl', 'my front porch', current_date);
SELECT * FROM birds.bird_view;
~~~

You should see four birds listed, the last being an "owl".  If you've made it this far exit `psql` with `\q`.

### Step 4

We previously setup Postgres to be ready to work with PostgREST. Now we need a configuration file. Create "postgrest.conf" Using your text editor.  The configuration file should contain.

~~~
db-uri = "postgres://birds:my_secret_password@localhost:5432/birds"
db-schemas = "birds"
db-anon-role = "birds_anonymous"
~~~

We are now ready run and test PostgREST.

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

If you have gotten this far, congratulations, you've created a working back end of our birds application!

### Step 5

In this section we'll need our text editor to create `htdocs/index.html` and `htdocs/sightings.js`.  We will use the Python 3 command `python3 -m http.server` to provide a web server.

In another terminal window or shell session change to the `htdocs` directory and fire up our Python web server.

~~~
cd htdocs
python3 -m http.server
~~~

Edit the file "index.html" so it looks like:

~~~
<DOCTYPE html>
<html>
  <body>
    <h1>Welcome to the bird list!</h1>
    <div id="bird-list"></div>
    <h2>Add a bird</h2>
    <div>
      <form>
        <div>
          <label for="bird">Bird</label> <input id="bird" id="bird" type="text" value="">
        </div>
        <div>
          <label for="place">Place</label> <input id="place" id="place" type="text" value="">
        </div>
        <div>
          <label for="sighted">Sighted on</label> <input id="sighted" id="sighted" type="date">
        </div>
        <button id="record-bird" type="submit">Add Bird Sighting</button>
      </form>
    </div>
    <script src="sightings.js"></script>
  </body>
</html>
~~~

We now need to populate our "bird-list" div element via JavaScript.

This is the basic data flow:

1. Get a handle on the "bird-list" div
2. Contact the JSON data API to get back our list of birds
3. When we have read it back generate HTML to insert inside "bird-list"

Here's the JavaScript to make the bird list appear.

~~~
/* sightings.js provides access to our JSON data API run by PostgREST
   and assembles the results before updating the web page. */
(function(document, window) {
  let list_url = 'http://localhost:3000/bird_view',
    record_url = 'http://localhost:3000/rpc/record_bird',
    list_elem = document.getElementById('bird-list'),
    bird_elem = document.getElementById('bird'),
    place_elem = document.getElementById('place'),
    sighted_elem = document.getElementById('sighted'),
    add_button = document.getElementById('record-bird');

  function updateList(elem, src) {
    let bird_list = JSON.parse(src),
      parts = [];
    parts.push('<table>');
    parts.push('  <tr><th>bird</th><th>place</th><th>sighted</th></tr>');
    for (const obj of bird_list) {
      parts.push(` <tr><td>${obj.bird}</td><td>${obj.place}</td><td>${obj.sighted}</td></td>`);
    }
    parts.push('</table>');
    elem.innerHTML = parts.join('\n');
  }

  function birdRecord(bird_elem, place_elem, sighted_elem) {
    return { "bird": bird_elem.value, "place": place_elem.value, "sighted": sighted_elem.value };
  }

  function getData(elem, url, updateFn) {
    /* We use a xhr to retrieve the current list of sightings. */
    const req = new XMLHttpRequest();
    req.addEventListener("load", function(evt) {
      /* Call our page update function */
      updateFn(elem, this.responseText);
    });
    req.open("GET", url);
	req.setRequestHeader('Cache-Control', 'no-cache');
    req.send();
  };

  function postData(obj, url) {
    const req = new XMLHttpRequest();
    req.open("POST", url, true);
    req.setRequestHeader('Content-Type', 'application/json');
    req.onreadystatechange = function() {//Call a function when the state changes.
      console.log(`DEBUG state ${req.readyState}, status ${req.status}, resp ${req.responseText}`);
    }
    req.send(JSON.stringify(obj));
  }

  /* Main processing for page */
  add_button.addEventListener("click", function(evt) {
    postData(birdRecord(bird_elem, place_elem, sighted_elem), record_url);
	/* Now we need to update our listing! */
	list_elem.innerHTML = '';
	console.log("DEBUG cleared list element, waiting");
	setTimeout(() => {
  		getData(list_elem, list_url, updateList);
  		console.log("Delayed for 1 second.");
	}, "1 second");
    evt.preventDefault();
  });

  getData(list_elem, list_url, updateList);
})(document, window);
~~~

Now you can test adding a bird by filling in the web form below "Add a bird"
and pressing the button "Add Bird Sighting".

## Closing thoughts

This is just a proof of concept demo to show what would be involved in writing a read/write web application based on PostgreSQL 15, PostgREST 11 and static pages. In the real world you'd do must more than this like use a web server like NginX or Apache 2 and reverse proxy to the static file system and the PostgREST microservice.  The web server can also proxy for other services too. An interesting one to support would be Markdown rendering via Pandoc server. You still might want to add some middle ware like Opensearch. I leave theses as exercises for the future.  This demo establishes that you can write a web application using just SQL, HTML and JavaScript. If you are interested in exploring this further I recommend going to both the Postgres website and follow their tutorials and documentation for what Postgres can do (it's allot) and go to the PostgREST site and explore both the documentation and examples. I've barely scratched the surface with this demo.


