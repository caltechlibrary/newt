
# Birds, a demo of PostgreSQL 15, PostgREST 11

This directory holds our demo.

## Setup

1. Start psql and connect to the birds database
1. Run [setup.sql](setup.sql)
2. Load [birds.csv](birds.csv)
3. Run a select query and confirm the data loaded
3. Quit psql and start Postgres
4. Start a localhost static web server for the htdocs directory
5. Point your browser at your static site web server and explore

~~~
psql
\c birds
\i setup.sql
\copy sighting from 'birds.csv' with (FORMAT CSV, HEADER);
SELECT * FROM sighting;
\q
~~~

