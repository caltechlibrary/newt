
# Newt router explained

## Overview

In the first Newt prototype supported a two stage pipeline for routing request through return a web page. It supported either Postgres+PostgREST through Pandoc web service or JSON API like Solr through Pandoc web service round trip. With the second prototype the Newt router has been generalized. Rather than two stages the second prototype implements a pipeline. This allows for several services to be tied together each sending a request to the next. This allows the web services to be more focused in much the same way that Unix programs can be chained together to form pipelines.  Using a route selector the generalized pipeline become steps indicated by a list of HTTP methods, URLs and content types. The YAML notation used has been significantly changed to support this generalization. Let's focus on the individual route setup[^31].

[^31]: See [Newt YAML syntax](newt_yaml_sentax.md) for complete description of the supported YAML. 

It is easy to start with a specific example then show it would be notated.

## A simple example

Let's say we have a database of music albums and reviews.  Each album includes a rating of "interesting". The range is a zero (uninteresting) to five star rating (most interesting). Previously we've modeled this in our Postgres database using a `view`.  How do we create a page that lists albums in descending order of interest? Since we're building with Newt we can assume there is a template to list albums available. That using that template will be the "last stage" in our pipeline. We need to feed the view into that template. The `view` statement is implemented in SQL in Postgres. That is exposed as a JSON API by PostgREST. That's our first stage, a JSON data source.

How do you representing a route with two stages?

```yaml
routes:
    - id: interesting_album_view
      request: GET /interesting_albums_list
      pipeline:
         - description: Contact PostgREST and get back the intersting album list
           service: GET http://localhost:3000/rpc/album_view
           content_type: application/json
           timeout: 15
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

