#!/bin/bash

#
# This file sets up a "birds" project folder and generates some of
# the documents needed to build our web application.
#
mkdir -p birds3/htdocs

# Create the database we'll use in the demo.
if ! createdb birds 2>/dev/null; then
	dropdb birds
	createdb birds
fi

# Generate a README
cat <<EOT >birds3/README.md

# Birds 3, a demo of PostgreSQL 15, PostgREST 11 and Newt

This directory holds our demo.

## Setup Database

1. Start psql and connect to the birds database
2. Run [setup.sql](setup.sql)
3. Run a select query and confirm the data loaded
4. Quit psql, you are ready to setup PostgREST

~~~
psql
\\c birds
\\i setup.sql
\\copy birds.sighting from 'birds.csv' with (FORMAT CSV, HEADER);
SELECT * FROM birds.sighting;
\\q
~~~

## Startup PostgREST

1. start PostgREST 'postgrest postgrest.conf'
2. Using curl make sure it is available 'http://localhost:3000/bird_view'

## Setup Newt

1. Create a birds.yaml file holding the routes for our birds application
2. Start up Newt service using the YAML file.
3. Point your web browser at the Newt and see what happens


EOT

# Generate our SQL setup modeling our simple data
cat <<EOT >birds3/setup.sql
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
  bird VARCHAR(255),
  place TEXT,
  sighted DATE
);

-- bird_view will become an end point in PostgREST
CREATE OR REPLACE VIEW birds.bird_view AS
  SELECT bird, place, sighted
  FROM birds.sighting ORDER BY sighted ASC, bird ASC;

-- record_bird is a stored procedure and will save a new bird sighting
CREATE OR REPLACE FUNCTION birds.record_bird(bird VARCHAR, place TEXT, sighted DATE)
RETURNS bool LANGUAGE SQL AS \$\$
  INSERT INTO birds.sighting (bird, place, sighted) VALUES (bird, place, sighted);
  SELECT true;
\$\$;

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

-- Import our CSV data into birds.sighting
\\copy birds.sighting from 'birds.csv' with (FORMAT CSV, HEADER);

EOT

# Generate some test data to load into our models
cat <<EOT >birds3/birds.csv
"bird","place","sighted"
"robin","seen in my backyard","2023-04-16"
"humming bird","seen in my backyard","2023-02-28"
"blue jay","seen on my back porch","2023-01-12"
EOT

# Generate a template of postgrest.conf file.
cat <<EOT >birds3/postgrest.conf
db-uri = "postgres://birds:my_secret_password@localhost:5432/birds"
db-schemas = "birds"
db-anon-role = "birds_anonymous"
EOT

# Generate birds.yaml
cat <<EOT >birds3/birds.yaml
env: [ "DB_NAME", "DB_USER", "DB_PASSWORD" ]
htdocs: htdocs
routes: 
  - req_path: /
    req_method: GET
    api_content_type: application/json
    api_method: GET
    api_url: http://localhost:3000/sighting
    pandoc_template: page.tmpl
    res_headers: { "content-type": "text/html" }
  - var: { "bird": "String", "place": "String", "sighted": "Date" }
    req_path: /
    req_method: POST
    api_content_type: application/json
    api_method: POST
    api_url: http://localhost:3000/sighting
    pandoc_template: post_result.tmpl
    res_headers: '{"content-type": "text/html"}'
  - req_path: /bird_views
    req_method: GET
    api_content_type: application/json
    api_method: GET
    api_url: http://localhost:3000/bird_view
    res_headers: { "content-type": "application/json" }
EOT


# Create a Pandoc template page.tmpl
cat <<EOT >birds3/page.tmpl
<DOCTYPE html lang="en">
<html>
  <head><title>Birds 3 Demo</title></head>
  <body>
    <header>Birds 3 Demo</header>
	<p>
    <h1>Welcome to the bird list!</h1>
    <p>
    <div id="bird-list"></div>
    <h2>Add a bird</h2>
    <div><form name="add_bird" method="POST" action="/">
      <div>
        <label for="bird">Bird</label>
                <input id="bird" name="bird" type="text" value="">
      </div>
      <div>
        <label for="place">Place</label>
                <input id="place" name="place" type="text" value="">
      </div>
      <div>
        <label for="sighted">Sighted on</label>
                <input id="sighted" name="sighted" type="date">
      </div>
      <button id="record-bird" type="submit">Add Bird Sighting</button>
    </form></div>
    <h2>Bird List</h2>
	<table>
		<thead>
			<tr class="header"> <th>bird</th> <th>place</th> <th>sighted</th> </tr>
		</thead>
		<tbody>
\$for(data)\$
			<tr> <td>\$it.bird\$</td> <td>\$it.place\$</td> <td>\$it.sighted\$</td> </tr>
\$endfor\$
		</tbody>
	</table>
	<p>
	<footer></footer>
  </body>
</html>
EOT

#
# Create a Pandoc template post_result.tmpl
#
cat <<EOT >birds3/post_result.tmpl
<DOCTYPE html lang="en">
<html>
<head> <meta http-equiv="refresh" content="0; url="/" /> </head>
<body>
Thank you for submitting a bird, <a href="/sighthings">View List</a>.
</body>
</html>
EOT

tree birds3
wc -l birds3/*.*

