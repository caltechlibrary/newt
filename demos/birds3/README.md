
# Birds 3, a demo of PostgreSQL 15, PostgREST 11 and Newt

This directory holds our demo.

## Setup Database

1. Start psql and connect to the birds database
2. Run [setup.sql](setup.sql)
3. Run a select query and confirm the data loaded
4. Quit psql, you are ready to setup PostgREST

~~~
psql
\c birds
\i setup.sql
\copy birds.sighting from 'birds.csv' with (FORMAT CSV, HEADER);
SELECT * FROM birds.sighting;
\q
~~~

## Startup PostgREST

1. start PostgREST 'postgrest postgrest.conf'
2. Using curl make sure it is available 'http://localhost:3000/bird_view'

## Setup Newt

1. Create a [birds.yaml](birds.yaml) file holding the routes for our birds application
2. Start up Newt service using the YAML file.
3. Point your web browser at the Newt and see what happens


