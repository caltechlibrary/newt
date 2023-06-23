
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

