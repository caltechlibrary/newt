#!/bin/bash

#
# This file sets up a "birds" project folder and generates some of
# the documents needed to build our web application.
#
if [ -d birds3 ]; then
	rm -fR birds3
fi
mkdir -p birds3/htdocs

# Create the database we'll use in the demo.
if ! createdb birds 2>/dev/null; then
	dropdb birds
	createdb birds
fi

# Generate a README
cat <<EOT >birds3/README.md

# Birds 3 demo

> Newt, Pandoc and Postgres+PostgREST

## Setup Newt

1. Create a [birds.yaml](birds.yaml) file holding the routes for our birds application
2. Start up Newt service using the YAML file.
3. Point your web browser at the Newt and see what happens

## Setup Database

1. Run [setup.sql](setup.sql) to configure PostgREST access to Postgres (normally NOT checked into Git/SVN!)
2. Run [models.sql](models.sql) to create our data models additional PostgREST end points.
3. Run [models_test.sql](models_test.sql) loads some test data and would run any SQL tests on the models.

~~~
psql -f setup.sql
psql -f models.sql
psql -f models_test.sql
~~~

## Start up our microservices

### Startup PostgREST

1. Start PostgREST 'postgrest postgrest.conf'
2. Using curl make sure it is available 'https://localhost:3000/bird_view'

### Startup Pandoc in server mode

In a seperate shell session start Pandoc server

1. 'pandoc server'

### Start up Newt

In a separate shell session start up Newt

1. 'newt birds.yaml'


EOT

# Generate our SQL PostgREST access
cat <<EOT>birds3/setup.sql
--
-- Following would not normally be include in a project's Git repository.
-- It contains a secret.  What I would recommend is writing a short
-- shell script that could generate this in a file, use that, then
-- checking in the shell script to version control since the secret
-- would not be stored in the file!
--

-- Make sure we are in the birds namespace/database
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
-- NOTE: The "CREATE ROLE" line is the problem line for
-- checking into your source control system. It contains a secret!
-- **DO NOT** store secrets in your SQL if you can avoid it!
--
DROP ROLE IF EXISTS birds;
CREATE ROLE birds NOINHERIT LOGIN PASSWORD 'my_secret_password';
GRANT birds_anonymous TO birds;

EOT

# This SQL models our data structures and behaviors
cat <<EOT >birds3/models.sql
--
-- Below is the SQL I would noramally check into a project repository.
-- It does not contain secrets. It contains our data models, views,
-- and functions. This defines the behaviors made available through
-- PostgRESTS.
--

-- Make sure we are in the birds namespace/database
\\c birds
-- Make sure our namespace is first in the search path
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
-- revoke role permissions here to match our model.
GRANT USAGE  ON SCHEMA birds      TO birds_anonymous;
-- NOTE: We are allowing insert because this is a demo and we are not
-- implementing an account login requirement. Normally I would force
-- a form update via a SQL function or procedure only.
GRANT SELECT, INSERT ON birds.sighting    TO birds_anonymous;
GRANT SELECT ON birds.bird_view   TO birds_anonymous;
GRANT EXECUTE ON FUNCTION birds.record_bird TO birds_anonymous;


EOT

# This SQL file loads test data into our data models.
cat <<EOT >birds3/models_test.sql
--
-- This script is a convienence. It will use the psql client
-- copy command to load the tables with some test data.
--

-- Make sure we are in the birds namespace/database
\\c birds


-- Now import our CSV file of birds.csv
\\copy birds.sighting from 'birds.csv' with (FORMAT CSV, HEADER);

-- Make sure the data loaded, query with a view statement.
SELECT * FROM birds.bird_view;

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
port: 8003
routes:
  - req_path: /
    req_method: GET
    api_content_type: application/json
    api_method: GET
    api_url: http://localhost:3000/bird_view
    pandoc_template: page.tmpl
    res_headers: { "content-type": "text/html" }
  - var: { "bird": "String", "place": "String", "sighted": "Date" }
    req_path: /
    req_method: POST
    api_content_type: application/json
    api_method: POST
    api_url: http://localhost:3000/rpc/record_bird
    pandoc_template: post_result.tmpl
    res_headers: { "content-type": "text/html" }
  - req_path: /api/bird_views
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
  <head></head>
  <body>
    <h1>Welcome to the bird list!</h1>
    <div id="bird-list">
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
	</div>

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
Thank you for submitting a bird, <a href="/">View List</a>.
</body>
</html>
EOT

tree birds3
wc -l birds3/*.*

