
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
    app_path: postgrest
    conf_path: postgrest.conf
routes:
  - id: interesting_album_view
    request: GET /interesting_albums_list
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
WARNING: demos/album_reviews/album_reviews.yaml has no models defined
PostgREST configuration is postgrest.conf
PostgREST will be run with the command "postgrest postgrest.conf"
Newt Router configured, port set to :8010
route interesting_album_view defined, request path GET /interesting_albums_list, pipeline size 1
~~~

You can ignore the "WARNING" about models because we've already set that up in our Postgres database.  

Low let's run the Newt Router and test with our web browsers.

~~~shell
newt run album_reviews.yaml
~~~

If you are happy with the results you can run PostgREST and Newt Router using the newt command, `newt run album_review.yaml`.


NOTE: You can check to see if you types in the YAML correctly with `newt check albums.yaml`

This Newt YAML file describes how to run PostgREST and the Newt Router. It then defines the routers that the Newt Router will respond to. Right now we only are responding to request from one route. That's a GET request to the URL "http://localhost:8010/interesting_albums_list". When that request is recieved it'll be pased on to PostgREST running on at "http://localhost:3000/interesting_albums_list". PostgREST will hand back it's results to your web browser.


      - description: |
          Take the results from PostgREST and run them through 
          the newtmustache using the template "ablum_list_view.tmpl"
          service: POST http://localhost:3032/album_list_view.tmpl
          content_type: application/json
          timeout: 10
      debug: true
```


What is being described?  First we have routes defined in our application. Our route is `interesting_album_view`. When a web browser contacts Newt via a GET at the designated path it triggers our pipeline property to start processing the request. In this case it is a two stage pipeline.

The first step retrieves the JSON data (i.e. the content is fetched from PostgREST). This is expressed as an HTTP method and URL. There is a content type that will be used when contacting the URL. You can also include a timeout value, in this case we're willing to wait up to 15 seconds. 

The second stage takes the output from the first and sends it through our template engine. Like the first there is a description for us humans, a service property indicating how and what URL we contact. There is a content type and timeout just like before. The output of the first stage is going to be sent as a port to the Mustache template engine as JSON data. The Mustache template engine returns the type based on the template, in this case HTML. 

There is a last property in our router description. `debug: true`. This property will cause the router to display more logging output. We want this to debug our application. It is very verbose. In a production setting this would be skipped or set to `false`. 

When the Newt router starts up it reads the YAML file and sets up to run each pipeline associated with specific request methods and paths. Those settings don't change until you restart the router. The router only reads the file once a startup. That's important to keep in mind. The router only interacts with "localhost". It listens for requests on a port of your choosing. It then can run a pipeline of other web services also running on "localhost".

This is where having the descriptions in the route definition is handy. It is easy to forgot which services are running on which ports. Both are URLs running as "localhost". In this specific case our PostgREST service is available on port 3000 and our Newt Mustache template engine is available on 3032. While the description element is optional it is what keep port mapping a little more transparent.  This is an area Newt could improve in the future but the reason for using a URL is that Newt doesn't need to know what each stage actually is. It just knows I contact this one service and take the output and send it to the next service and all stages of the pipeline are complete or there is an error reported in the pipeline. The result has handed back to the web browser.


## Changes from the first prototype to the second.

- routes include a pipeline rather than fixed stages
- `newt` was replaced by `newtrouter`. It does less. It just routes data now. It does more, you can have any number of stages in our data pipeline now. It doesn't know how to package things.
- `newtmustache` has replaced Pandoc web service as the Newt template engine of choice. Mustache is a popular templating language and well supported in numerious programming languages. It provided easier to debug issue than working with Pandoc server. `newtmustache` does require of configuration. 
- each pipeline stage has its own timeout setting

While there isn't a fixed limit to the number of stages in a pipeline you will want to keep the number limited. While contacting a web service on localhost is generally very quick the time to round trip the communication still accumulates. As a result it is recommend to stick to less than four stages and to explicitly set the timeout values based on your understanding of performance requirements.  If a stage times out the a response will be returned as an HTTP error.

## Misc

If a requested route is not defined in our YAML by then the router will look matching static files. If that fails an HTTP 404 is returned. For a request route to match our YAML the router compares HTTP method, path and content type. If any of these don't match then the route is not considered a match and will return an appropriate HTTP status and code.

The router uses are defined in the request property. The HTTP method and path indicate what can trigger the pipeline being run.

The Newt router will only support HTTP and run on localhost. This is the same approach taken by Pandoc server. It minimize configuration and also discourages use as a front end service (which is beyond the scope of this project).

This prototype does not support file uploads. In theory you could implement a pipe line stage to handle that but again that is beyond the scope of this project. You can try clever techniques browser side and push objects into Postgres via PostgREST but again that is beyond the scope of this project. I don't recommend that. If you need file upload support Newt project isn't your solution yet.

