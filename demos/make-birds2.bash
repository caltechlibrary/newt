#!/bin/bash

#
# This file sets up a "birds" project folder and generates some of
# the documents needed to build our web application.
#
if [ -d birds2 ]; then
	rm -fR birds2
fi
mkdir -p birds2/htdocs

# Create the database we'll use in the demo.
if ! createdb birds 2>/dev/null; then
	dropdb birds
	createdb birds
fi

# Generate a README
cat <<EOT>birds2/README.md

# Birds 2 Demo

> Postgre+PostgREST and browser JavaScript

## Setup Database

1. Run [setup.sql](setup.sql) to configure PostgREST access to Postgres (normally NOT checked into Git/SVN!)
2. Run [models.sql](models.sql) to create our data models additiojnal PostgREST end points.
3. Run [models_test.sql](models_test.sql) loads some test data and would run any SQL tests on the models.

~~~
psql -f setup.sql
psql -f models.sql
psql -f models_test.sql
~~~

## Startup PostgREST

1. Start PostgREST 'postgrest postgrest.conf'
2. Using curl make sure it is available 'https://localhost:3000/bird_view'

## Startup static web server on local host

1. In another shell session go to the htdocs directory
2. Start a static web server on localhost, e.g. 'python3 -m http http.server'
3. Point your web browser at the static web server and confirm you see our birds 2 website
4. Open your web browser tools and submit a new bird, you should see the transaction where the JavaScript executes, contacts PostgREST and then updates the page

EOT

# Generate our SQL PostgREST access
cat <<EOT>birds2/setup.sql
--
-- Following I would normally not include in a project SQL codebase.
-- It contains a secret.  What I would recommend is writing a short
-- shell script that could generate this in a file, use that, then
-- checking in the shell script to version control since the secret
-- would not be stored in the file!
--

-- Make sure we are in the birds database
\\c birds

--
-- Setup a Postgres "schema" (a.k.a. namespace) for
-- working with PostgREST
--
DROP SCHEMA IF EXISTS birds CASCADE;
CREATE SCHEMA birds;

--
-- The following additional steps are needed for PostgREST access
-- are birds schema and database.
--
DROP ROLE IF EXISTS birds_anonymous;
CREATE ROLE birds_anonymous nologin;

--
-- NOTE: The "CREATE ROLL" line is the problem line for
-- checking into your source control system. It contains a secret!
-- **DO NOT** store secrets in your SQL if you can avoid it!
--
DROP ROLE IF EXISTS birds;
CREATE ROLE birds NOINHERIT LOGIN PASSWORD 'my_secret_password';
GRANT birds_anonymous TO birds;

EOT

# This SQL models our data structures and behaviors
cat <<EOT >birds2/models.sql
--
-- Below is the SQL I would noramally check into a project repository.
-- It does not contain secrets. It contains our data models, views,
-- and functions. This defines the behaviors made available through
-- PostgRESTS.
--

-- Make sure we are in the birds database
\\c birds
SET search_path TO birds, public;

--
-- Data Models
--

--
-- This file contains the create statements for our bird table.
--
-- DROP TABLE IF EXISTS birds.sighting;
CREATE TABLE birds.sighting
(
  bird VARCHAR(255),
  place TEXT,
  sighted DATE
);

--
-- Data Views and Behaviors.
--

-- bird_view will become an end point in PostgREST, '/bird_view'
CREATE OR REPLACE VIEW birds.bird_view AS
  SELECT bird, place, sighted
  FROM birds.sighting ORDER BY sighted ASC, bird ASC;

-- record_bird is a stored procedure and will save a new bird sighting.
-- It becomes the end point '/rpc/record_bird'
CREATE OR REPLACE FUNCTION birds.record_bird(bird VARCHAR, place TEXT, sighted DATE)
RETURNS bool LANGUAGE SQL AS \$\$
  INSERT INTO birds.sighting (bird, place, sighted) VALUES (bird, place, sighted);
  SELECT true;
\$\$;

--
-- PostgREST access and controls.
--

-- Since our Postgres ROLE and SCHEMA exist and our models may change how
-- we want PostgREST to expose our data via JSON API we GRANT or 
-- revoke role permissions here.
-- with our model.
GRANT USAGE  ON SCHEMA birds      TO birds_anonymous;
-- NOTE: We are allowing insert because this is a demo and we are not
-- implementing an account login requirement. Normally I would force
-- a form update via a SQL function or procedure only.
GRANT SELECT, INSERT ON birds.sighting    TO birds_anonymous;
GRANT SELECT ON birds.bird_view   TO birds_anonymous;
GRANT EXECUTE ON FUNCTION birds.record_bird TO birds_anonymous;


EOT

# This SQL file loads test data into our data models.
cat <<EOT >birds2/models_test.sql
--
-- This script is a convienence. It will use the psql client
-- copy command to load the tables with some test data.
--

-- Make sure we are in the birds database
\\c birds


-- Now import our CSV file of birds.csv
\\copy birds.sighting from 'birds.csv' with (FORMAT CSV, HEADER);

-- Make sure the data loaded, query with a select statement.
SELECT * FROM birds.sighting;

EOT

# Generate some test data to load into our models
cat <<EOT>birds2/birds.csv
"bird","place","sighted"
"robin","seen in my backyard","2023-04-16"
"humming bird","seen in my backyard","2023-02-28"
"blue jay","seen on my back porch","2023-01-12"
EOT

# Generate a template of postgrest.conf file.
cat <<EOT>birds2/postgrest.conf
db-uri = "postgres://birds:my_secret_password@localhost:5432/birds"
db-schemas = "birds"
db-anon-role = "birds_anonymous"
EOT

# Generate index.html
cat <<EOT>birds2/htdocs/index.html
<DOCTYPE html lang="en">
<html>
  <body>
    <h1>Welcome to the bird list!</h1>
    <div id="bird-list"></div>
    <h2>Add a bird</h2>
    <div><form>
      <div>
        <label for="bird">Bird</label>
		<input id="bird" id="bird" type="text" value="">
      </div>
      <div>
        <label for="place">Place</label>
		<input id="place" id="place" type="text" value="">
      </div>
      <div>
        <label for="sighted">Sighted on</label>
		<input id="sighted" id="sighted" type="date">
      </div>
      <button id="record-bird" type="submit">Add Bird Sighting</button>
    </form></div>
    <script src="sightings.js"></script>
  </body>
</html>
EOT

# Generate sightings.js
cat <<EOT>birds2/htdocs/sightings.js
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
      parts.push(\` <tr><td>\${obj.bird}</td><td>\${obj.place}</td><td>\${obj.sighted}</td></td>\`);
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
      console.log(\`DEBUG state \${req.readyState}, status \${req.status}, resp \${req.responseText}\`);
    }
    req.send(JSON.stringify(obj));
  }

  /* Main processing for page */
  add_button.addEventListener("click", function(evt) {
    postData(birdRecord(bird_elem, place_elem, sighted_elem), record_url);
	/* Now we need to update our listing! */
	list_elem.innerHTML = '';
	setTimeout(() => {
  		console.log("Delayed for 10 second.");
  		getData(list_elem, list_url, updateList);
	}, "10 second");
    evt.preventDefault();
  });

  getData(list_elem, list_url, updateList);
})(document, window);
EOT

tree birds2
wc -l birds2/*.* birds2/htdocs/*.*

