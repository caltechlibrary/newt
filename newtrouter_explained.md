
# Newt router explained

## Overview

The Newt Router is a data router. A data router takes a request then pass it on to a list one or more services returning the final result. The sequence of contacting one service and using the result as input to another service is called a pipeline. The Unix family of operating systems shells support this concept of piping the output of one program into another program. The Newt router provides a similar service but for connecting one or more web services[^1]. Newt was created to be able to route request between PostgREST and a template engine. The second Newt prototype supports this concept for any web service that is reachable via localhost. Typically this is still PostgREST and the Newt Mustache template engine. PostgREST returns JSON data and Newt Mustache can take the JSON data and render an HTML page using a Mustache template. You could swap PostgREST out for Solr and do the same thing. Newt Router supports pipelines of one or more services. The last service to respond has its results passed back to the requesting web browser.

[^1]: Newt Router was inspired by [Yahoo Pipes!](https://en.wikipedia.org/wiki/Yahoo!_Pipes)

Newt Router is organized round two concepts, "routes" and "pipelines". A "route" describes the HTTP method and URL path of a web browser sends to the service. The pipeline is the sequence of services that process that request. The Newt Router uses a YAML file to describe the mapping of requests to pipelines.

The Newt Router can also function as a static web server. If you're pipeline results in HTML output, the static assets (e.g. CSS, JavaScript, image files) can add the polish to your human facing interface.

## A brief tour, albums and reviews

Let's say we have a database of music albums and reviews.  Each database entry includes the album name, the review and a rating of "interesting". The range is a zero (uninteresting) to five star rating (most interesting). This information is stored in a Postgres database and made available to the Newt Router via PostgREST. Our web application needs to be able to use the PostgREST JSON API to manage our list of albums and reviews. I am going to assume you have a Postgres "view" called "interesting_album_view" defined and that is available via PostgREST via a GET request at the URL "http://localhost:3000/interesting_albums_list".

### Step 1. Prep work

Before we can run through the tutorial somethings need to be up and running.  We need Postgres 16 and PostgREST installed. Information about installing [Postgres](https://postgres.org) and [PostgREST](https://postgrest.org) can be found on their respective websites.

The Postgres database needs to be created for our demo and it needs to be configured to allow PostgREST to access it. That can be accomplished by downloading this SQL file <https://github.com/caltechlibrary/newt/blob/main/demos/album_reviews/setup_album_reviews_demo.sql>.

Here's the steps.

1. Download the SQL file (I use curl for this)
2. Edit the SQL file and change the line with "my_secret_password" to something more appropriate
3. Create the "album_reviews" database
4. Run the SQL program to do the setup.

~~~shell
curl -L -o https://github.com/caltechlibrary/newt/blob/main/demos/album_reviews/setup_album_reviews_demo.sql
nano setup_album_reviews_demo.sql
createdb album_reviews
psql album_reviews < setup_album_reviews_demo.sql
~~~

Similarly you can get an example "postgrest.conf" from <https://github.com/caltechlibrary/newt/blob/main/demos/album_reviews/postgrest.conf>. You'll need to set the password to match the one you used in the SQL file setting up the Postgres and PostgREST.

~~~shell
curl -L -o https://github.com/caltechlibrary/newt/blob/main/demos/album_reviews/postgrest.conf
nano postgrest.conf
~~~

Remember when you edit the "postgrest.conf" file you need to have the password match the one you used
in the "setup_album_reviews_demo.sql" file.

### Step 2. Building our application

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

This file is going to define a single "route" with a one stage pipeline that proxies to PostgREST returning the results of our "interesting_albums_list" SQL view. You can check your YAML and to make sure you've typed in everything correctly using the command.

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

For now ignore the "WARNING" about models. The SQL program "setup_album_reviews_demo.sql" already created the database, table and functions to the models are not necessary for this demo.

Now let's run the Newt Router.

~~~shell
newt run album_reviews.yaml
~~~

Point your web browser at the URL <http://localhost:8010/api/>.

You should get back a JSON response that originated from Postgres+PostgREST. It's not very useful yet.
We have more to do.  You should press "ctrl-c" to exit the newt command before continuing. This will stop both the Newt Router and the PostgREST service since we will be making changes in the second part.

It is important to realize that the Newt Router loads it's configuration at startup only. This means if the Newt Router is running and you change the YAML file the running router will not act on those changes.

## Step 3. Improving our application by adding routes and pipelines

We can improve our web application by expanding our pipelines to include generating HTML for the web browser. This can be done with Mustache templates run from Newt Mustache as part of our pipelines.

Download the Mustache templates from <https://github.com/caltechlibrary/newt/blob/main/demos/album_reviews/review_list.tmpl> and <https://github.com/caltechlibrary/newt/blob/main/demos/album_reviews/review_submitted.tmpl>. Save them in your directory.

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

This is allot of YAML. Check your updated Newt YAML file with the following command.

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
newt -verbose run album_reviews.yaml
~~~

Point your web browser at <http://localhost:8010>. You should see an album list with one entry. You can add a review using the web form at the bottom of the page. Then you complete the web form and press the submit button it should take you to a page showing the review you just submitted. You can click the link to return to the list and see it has been updated.

## What have we done exactly?

First we've built up a simple web application through defining some routes using data pipelines. In our first iteration we just verified that we could connect to PostgREST using a simple one stage pipeline. It nice but not really compelling for must of us humnas.

In the second iteration we use the Newt Router to run two additional routes. The first one listed our review list like before but this time the results were displayed as a web page (i.e. HTML markup). This was accomplished by adding Newt Mustache in our pipleines. We use one route to display our list and include a web form for submitting another review. The second route we added displays the results from the submission.

By typing in the Newt YAML file, adding some Mustache templates we have a functional web application that is built on what is provided by Postgres and PostgREST.
