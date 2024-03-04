
# Newt router explained

## Overview

The Newt Router is a data router. A data router takes a request then pass it on to a list one or more services returning the final result. The sequence of contecting one service and using the result as input to another service is called a pipeline. The Unix family of operating systems shells support this concept of piping the output of one program into another program. The Newt router provides a similar service but for routine between one or more web services. Newt was created to be able to route request between PostgREST and a template engine. The second Newt prototype supports this concept for any web service that is reachable via localhost. Typically this is still PostgREST and the Newt Mustache template engine. PostgREST returns JSON data and Newt Mustache can take the JSON data and render an HTML page using a Mustache template. You could swap PostgREST out for Solr and do the same thing. Newt Router supports pipelines of one or more services. The last service to respond has its results passed back to the requesting web browser.

Newt Router does using to concepts, "routes" and "pipelines". A "route" describes the HTTP method and URL path of a request. The Newt Router using the Newt YAML file can map a requets to a pipeline. Using this approach of defines routes and their pipelines you can compose your web application.

Newt Router can also function as a static web server. This is helpful where your final output is HTML and you wish to also include any related static assets like CSS, images or JavaScript.

## A simple example

Let's say we have a database of music albums and reviews.  Each album includes a rating of "interesting". The range is a zero (uninteresting) to five star rating (most interesting). This information is stored in a Postgres database and made available to the Newt Router via PostgREST. Our web application needs to be able to use the PostgREST JSON API to manage our list of albums and reviews. I am going to assume you have a Postgres "view" called "interesting_album_view" defined and that is available via PostgREST via a GET request at the URL "http://localhost:3000/interesting_albums_list".

### Prep work

Before we can run through the tutorial somethings need to be up and running.
Postgres 16 and PostgREST 12 needs to be installed. The Postgres database for the demo needs to be installed along with configuration for allowing PostgREST to provide the database content. Here's an example SQL file you can run to create your Postgres database configured for using with PostgREST. You can retrieve the SQL code to set things up from <https://github.com/caltechlibrary/newt/blob/main/demos/album_reviews/setup_album_reviews_demo.sql>. You need to change the password in the file before running the following commands.

~~~shell
createdb album_reviews
psql < setup_album_reviews_demo.sql
~~~

Similarly you can get an example "postgrest.conf" from <https://github.com/caltechlibrary/newt/blob/main/demos/album_reviews/postgrest.conf>. You'll need to set the password to match the one you used in the SQL file setting up the Postgres and PostgREST.

### Building our application

Let's create a Newt YAML file called "album_reviews.yaml". Type in the following using your favorite text editor.

~~~yaml
applications:
  newtrouter:
    port: 8010
  postgrest:
    port: 3000
    app_path: postgrest
    conf_path: postgrest.conf
routes:
  - id: interesting_album_view
    request: GET /api/{$}
    pipeline:
      - service: GET http://localhost:3000/interesting_albums_list
        description: Contact PostgREST and get back the intersting album list
~~~

You can check to make sure you've typed in everything correctly using the command.

~~~shell
newt -verbose check album_reviews.yaml
~~~

That should show some output like

~~~text
WARNING: album_reviews.yaml has no models defined
PostgREST configuration is postgrest.conf
PostgREST will be run with the command "postgrest postgrest.conf"
Newt Router configured, port set to :8010
route interesting_album_view defined, request path GET /api/{$}, pipeline size 1
~~~

You can ignore the "WARNING" about models because we've already set that up in our Postgres database.  

Low let's run the Newt Router.

~~~shell
newt run album_reviews.yaml
~~~

Point your web browser at the url <http://localhost:8010/api/>.

You should get back a JSON response that originated from Postgres+PostgREST. It's not very useful yet so 
we have more todo.

You should press "ctrl-c" to exit the newt command before continuing. This will stop both the
Newt Router and the PostgREST service since we will be making changes in the second part.


## Improving our application by adding routes and pipelines

We can do improve our web application by creating a Mustache templates and setting up Newt Mustache as part of our pipelines.

Add these two Mustache templates from <https://github.com/caltechlibrary/newt/blob/main/demos/album_reviews/review_list.tmpl> and <https://github.com/caltechlibrary/newt/blob/main/demos/album_reviews/review_submitted.tmpl>.

Update the album_reviews.yaml to look like this.

~~~yaml
applications:
  newtrouter:
    port: 8010
  newtmustache:
    port: 8011
  postgrest:
    app_path: postgrest
    conf_path: postgrest.conf
routes:
  - id: interesting_album_view
    request: GET /api/{$}
    pipeline:
      - service: GET http://localhost:3000/interesting_albums_list
        description: Contact PostgREST and get back the intersting album list
  - id: add_a_review
    request: POST /add_a_review
    pipeline:
      - service: POST http://localhost:3000/rpc/add_a_review
        description: Add a review via PostgREST function album_reviews.add_a_review
      - service: POST http://localhost:8011/review_submitted
        description: Display the submitted review with link back to list
  - id: show_reviews
    request: GET /{$}
    pipeline:
      - service: GET http://localhost:3000/interesting_albums_list
        description: Contact PostgREST and get back the intersting album list
      - service: POST http://localhost:8011/review_list
        description: Convert the JSON into HTML, show the list and web form
templates:
  - request: /review_list
    template: review_list.tmpl
  - request: /review_submitted
    template: review_submitted.tmpl
~~~

You can check your updated Newt YAML file with the following command.

~~~shell
newt -verbose check album_reviews.yaml
~~~

This time you should see results like the following.

~~~text
WARNING: album_reviews.yaml has no models defined
PostgREST configuration is postgrest.conf
PostgREST will be run with the command "postgrest postgrest.conf"
Newt Router configured, port set to :8010
route interesting_album_view defined, request path GET /api/{$}, pipeline size 1
route add_a_review defined, request path POST /add_a_review, pipeline size 2
route show_reviews defined, request path GET /{$}, pipeline size 2
Newt Mustache configured, port set to 8011
2 Mustache Templates are defined
http://localhost8011/review_list points at review_list.tmpl
http://localhost8011/review_submitted points at review_submitted.tmpl
~~~

If this matches what you've see we're ready to run Newt again with our updated templates and YAML file.

~~~shell
newt run album_reviews.yaml
~~~

Point your web browser at <http://localhost:8010>. You should see an empty album list. You can add a review in the form and that should take you to a page showing the added review. Click the link to return to the list and see the revised results.

## What have we done exactly?

First we've built up a simple web appliction through defining some routes with data pipelines. In our first iteration we just verified that we could conntect to PostgREST in our first pipeline. It was uninstesting because what is returned in a JSON and specifically an empty JSON list. Not very exciting.

In the second iteration we use the Newt Router to run two routes. The first one listed our review list like before but this time the results were displayed as a web page (i.e. HTML markup). This was accomplished by adding Newt Mustache to our first router as the last stage of our pipeline.

We also added support for a second route that handled the web form submission then returned and displayed the results.
