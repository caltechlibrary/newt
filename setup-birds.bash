#!/bin/bash


#
# This file sets up a "birds" project folder and generates some of
# the documents needed to build our web application.
#
mkdir -p birds/htdocs

# Generate the empty files we'll use in the demo.
touch birds/setup.sql
touch birds/birds.csv
touch birds/postgrest.cnf
touch birds/htdocs/index.html
touch birds/htdocs/birds.js

# Create the database we'll use in the demo.
if ! createdb birds 2>/dev/null; then
	dropdb birds
	createdb birds
fi

# Generate a README
cat <<EOT>birds/README.md

# Birds, a demo of PostgreSQL 15, PostgREST 11

This directory holds our demo.

EOT

# Generate our SQL setup modeling our simple data
cat <<EOT>birds/setup.sql
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
AS \$\$
  INSERT INTO birds.sighting (bird_name, place, sighted) VALUES (name, description, dt);
\$\$;

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

EOT

# Generate some test data to load into our models
cat <<EOT>birds/birds.csv
bird,place,sighted
robin, seen in my backyard,2023-04-16
humming bird, seen in my backyard, 2023-02-28
blue jay, seen on my back porch, 2023-01-12
EOT

# Generate a template of postgrest.conf file.
cat <<EOT>birds/postgrest.conf
db-uri = "postgres://birds:replace_me@localhost:5432/postgres"
db-schemas = "birds"
db-anon-role = "birds_anonymous"
EOT

