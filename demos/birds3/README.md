
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


