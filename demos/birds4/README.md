
# Birds 4 demo

This demo shows how to build a Newt application from the Newt YAML file.
It uses two hand coded Mustache templates, page.tmpl and post_result.tmpl
as well as some hand coded Postgres SQL files and a hand coded PostgREST
configuration.  Normally you'd want to use Newt's code generator to generate
these things but this is a trivially simple app.

The Newt runner (i.e. `newt`)
to test the generated results.

The birds.yaml defines the application including how to run PostgREST from 
it's generated config.

## Setup Newt

1. Create a [birds.yaml](birds.yaml) file holding the applications configurations,
   data models, routes and template organization for our birds application.
2. Create page.tmpl, post_result.tmpl and postgrest.conf
3. Start up Newt service using the YAML file.
4. Point your web browser at the Newt and see what happens

### Creating the birds.yaml file.

You can use your favorite text editor to create the `birds.yaml`. The syntax for
the file is documented at <https://caltechlibrary.github.io/newt/newt_yaml_syntax.html>.

~~~yaml
# Birds 4 is a demonstration of the 2nd Newt Prototype. The YAML
# file has gone through quite a transformation now.
applications:
  newtrouter:
    port: 8010
    htdocs .
  newtmustache:
    port: 8011
  postgrest:
    app_path: postgrest
    conf_path: birds.conf
    port: 3000
routes:
  - id: bird_view
    request: GET /
    description: Newt home page, show the bird view results
    pipeline:
      - service: GET http://localhost:3000/bird_view
        description: retrieve the JSON of the bird view
      - service: POST http://localhost:8011/list.tmpl
        description: Transform the JSON into HTML
  - id: record_bird
    request: POSTS / 
    description: Record a bird sighting
    pipeline:
      - service: POST http://localhost:3000/rpc/record_bird
        description: Send the form data to PostgREST
      - service: POST http://localhost:8011/post_result.tmpl
        description: Display the HTML results of POST to PostgREST
  - id: api_bird_views
    request: GET /api/bird_view
    description: Display the PostgREST API end point for bird_view
    pipeline:
      - service: http://localhost:3000/bird_view
        description: Return the raw JSON result
templates:
  - request: /post_result.tmpl
    template: create.tmpl
  - request: /list.tmpl
    template: page.tmpl
~~~

### Create your two templates. Generate our code from the birds.yaml

The page.tmpl file would look like.

~~~html
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
{{#data}}
			<tr> <td>{{bird}}</td> <td>{{place}}</td> <td>{{sighted}}</td> </tr>
{{/data}}
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
~~~

The post_result.tmpl would look like this.

~~~html
<DOCTYPE html lang="en">
<html>
<head> <meta http-equiv="refresh" content="0; url="/" /> </head>
<body>
Thank you for submitting a bird, <a href="/">View List</a>.
</body>
</html>
~~~

# the postgrest.conf would look like this

~~~
db-uri = "postgres://birds:my_secret_password@localhost:5432/birds"
db-schemas = "birds"
db-anon-role = "birds_anonymous"
~~~

### Setup Database

1. Run [setup.sql](setup.sql) to configure PostgREST access to Postgres (normally NOT checked into Git/SVN!)
2. Run [models.sql](models.sql) to create our data models additional PostgREST end points.
3. Run [models_test.sql](models_test.sql) loads some test data and would run any SQL tests on the models.

~~~
psql -f setup.sql
psql -f models.sql
psql -f models_test.sql
~~~

### Start up our microservices

1. Startup Newt
2. Use your web browser to explore PostgREST API, see http://localhost:3000 (if using the default configuration)
3. Use your web browser to explore running applications, see http://localhost:8010 (if using the default configuration).

When developing your app you can use the `newt` to start and stop the Newt Router, Newt Mustache engine and PostgREST.
All three will log their output to the console.

~~~
newt birds.yaml
~~~

NOTE: you can stop `newt` by pressing "ctrl-c" in the terminal window where you started it.

In a new terminal window try the following to test PostgREST and Newt services.

On macOS you can use the `open` command, on Linux you'd use `xdg-open`.

macOS

~~~
open http://localhost:3000
open http://localhost:8010
~~~

Linux

~~~
xdg-open http://localhost:3000/bird_view
xdg-open http://localhost:8010
~~~

