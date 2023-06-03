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

# Generate birds-routes.csv
cat <<EOT >birds3/birds-routes.csv
req_path,req_method,api_url,api_method,api_content_type,pandoc_template,res_headers
/sightings,GET,http://localhost:3000/sighting,GET,application/json,page.tmpl,"{""content-type"": ""text/html""}"
/sightings,POST,http://localhost:3000/sighting,POST,application/json,redirect.tmpl,"{""content-type"": ""text/html""}"
/bird_views,GET,http://localhost:3000/bird_view,GET,application/json,,"{""content-type"": ""application/json""}"
EOT

# Generate birds.yaml configuration file
cat <<EOT >birds3/birds.yaml
newt_routes: "birds-routes.csv"
newt_htdocs: "htdocs"
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
    <div><form method="POST" action="/sightings">
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
# Create a Pandoc template page.tmpl
cat <<EOT >birds3/redirect.tmpl
<DOCTYPE html lang="en">
<html>
<head>
	<meta http-equiv="refresh" content="0; url="/sightings" />
</head>
<body>
Thank you for submitting a bird, <a href="/sighthings">View List</a>.
</body>
</html>
EOT

# Create index.md
cat <<EOT >birds3/htdocs/index.md

# Birds 3 Demo

- [about.html](/about.html) (static HTML page)
- [sightings](/sightings)   (HTML page)
- [bird views](/bird_views) (JSON output)

EOT

# Create index.html
pandoc --from markdown --to html5 \
	--metadata title="Birds Demo 3" \
	-s \
	birds3/htdocs/index.md >birds3/htdocs/index.html


# Create about.html
cat <<EOT >birds3/htdocs/about.md
#
# About Birds 3 Demo

This demo shows how Newt used with Postres+PostgREST and Pandoc
can provide a complete stack for basic web development that does not
involve file uploads.

EOT

# Create index.html
pandoc --from markdown --to html5 \
	--metadata title="Birds Demo 3" \
	-s \
	birds3/htdocs/about.md >birds3/htdocs/about.html


# Create the setup and run bash script.
cat <<EOT >birds3/setup-for-birds.bash
#!/bin/bash

dropdb birds
createdb birds
psql -d birds -c '\i setup.sql'
pandoc server &
postgrest postgrest.conf
EOT

chmod 775 birds3/*.bash
