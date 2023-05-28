#!/bin/bash

#
# This file sets up a "birds" project folder and generates some of
# the documents needed to build our web application.
#
mkdir -p birds3/htdocs

# Generate the empty files we'll use in the demo.
touch birds3/setup.sql
touch birds3/birds-routes.csv
touch birds3/birds.yaml
touch birds3/postgrest.conf
touch birds3/htdocs/about.html

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
3. Load [birds.csv](birds.csv)
4. Run a select query and confirm the data loaded
5. Quit psql, you are ready to setup PostgREST

~~~
psql
\\c birds
\\i setup.sql
\\copy sighting from 'birds.csv' with (FORMAT CSV, HEADER);
SELECT * FROM sighting;
\\q
~~~

## Startup PostgREST

1. start PostgREST 'postgrest postgrest.conf'
2. Using curl make sure it is available 'https://localhost:3000/bird_view'

## Setup Newt

1. Create a birds-routes.csv file holding the routes for our birds application
2. Create birds.yaml file with 'newt_routes' and 'new_htdocs' defined.
3. Start up Newt service using the YAML file.
4. Point your web browser at the Newt and see what happens


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
RETURNS bool LANGUAGE SQL AS \$\$
  INSERT INTO birds.sighting (bird_name, place, sighted) VALUES (bird, place, sighted);
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

EOT

# Generate some test data to load into our models
cat <<EOT >birds3/birds.csv
bird_name,place,sighted
robin, seen in my backyard,2023-04-16
humming bird, seen in my backyard, 2023-02-28
blue jay, seen on my back porch, 2023-01-12
EOT

# Generate a template of postgrest.conf file.
cat <<EOT >birds3/postgrest.conf
db-uri = "postgres://birds:my_secret_password@localhost:5432/birds"
db-schemas = "birds"
db-anon-role = "birds_anonymous"
EOT

# Generate birds-routes.csv
cat <<EOT >birds3/birds-routes.csv
req_path,req_method,req_content_type,api_url,api_method,api_content_type,pandoc,pandoc_options,pandoc_template,res_headers
/sightings,GET,text/html,https://localhost:3000/bird_view,GET,application/json,true,from=json&to=html5&standalone=true,page.tmpl,"{""content-type"": ""text/html""}"
EOT

# Generate birds.yaml configuration file
cat <<EOT >birds3/birds.yaml
newt_routes: "birds-routes.csv"
newt_htdocs: "htdocs"
EOT

# Create a Pandoc template
cat <<EOT >birds3/page.tmpl
<DOCTYPE html lang="en"\>
<html>
  <head><title>Birds 3 Demo</title></head>
  <body>
    <header>Birds 3 Demo</header>
	<p>
    <h1>Welcome to the bird list!</h1>
    <p>
<h2>Bird List</h2>
<table>
<thead>
<tr class="header">
<th>bird_name</th>
<th>place</th>
<th>sighted</th>
</tr>
</thead>
<tbody>
\$for(birds)\$
<tr>
<td>\$it.name\$</td>
<td>\$it.place\$</td>
<td>\$it.sighted\$</td>
</tr>
\$endfor\$
</tbody>
</table>
	<p>
	<footer></footer>
  </body>
</html>
EOT


# Create about.html
cat <<EOT >birds3/htdocs/about.html
<!DOCTYPE html lang="en">
<html>
<head>
	<title>Birds 3 Demo</title>
</head>
<body>
<h1>About Birds 3 Demo</h1>
<p>This demo shows how Newt used with Postres+PostgREST and Pandoc
can provide a complete stack for basic web development that does not
involve file uploads.</p>

</body>
</html>
EOT

