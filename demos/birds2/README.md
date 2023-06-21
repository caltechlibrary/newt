
# Birds, a demo of PostgreSQL 15, PostgREST 11

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
SELECT * FROM birds.sighting;
\q
~~~

## Startup PostgREST

1. start PostgREST 'postgrest postgrest.conf'
2. Using curl make sure it is available 'https://localhost:3000/bird_view'

## Startup static web server on local host

1. in another shell session go to the htdocs directory
2. start a static web server on localhost, e.g. 'python3 -m http http.server'
3. Point your web browser at the static web server and see what happens



