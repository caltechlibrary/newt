
# Overview of `newt` data router

In the first prototype `newt` supported a two stage pipeline. It support either Postgres+PostgREST through Pandoc web service or JSON API like Solr through Pandoc web service round trip. With the second prototype Newt is generalizing this two stages associated with a route selector to general pipeline with steps indicated by a list of URLs, methods and content types. While the notation in YAML has changed conceptually it is similar just allowing for one or my steps in the pipeline. Let's focus on the inidividual route setup[^1].

[^1]: See [Newt YAML syntax](newt_yaml_sentax.md) for complete discription of the supported YAML. 

It is easy to start with a specific example then show it would be notated.

Let's say we have a database of music albums and reviews.  Each album includes a rating on interesting. Zero (uninstersting) to five star rating (most interesting). Previously we've modeled this in our Postgres database using a `view`.  How do we create a page that lists albums in decending orders of interest? Since we're building with Newt we can assume there is a template to list albums available. That template will be the "last stage" in our pipeline. We need to feed the view into that template. The `view` statement mapped to the PostgREST JSON API is our first stage.

How do you representing a route with two stages?

```yaml
routes:
    - id: interesting_album_view
      pipeline:
         - description: Contact PostgREST and get back the intersting album list
           service: GET http://localhost:3000/rpc/album_view
           content_type: application/json
           timeout: 15
         - description: |
             Take the results from PostgREST and run them through 
             the newtmustache use the template "ablub_list_view.tmpl"
           service: POST http://localhost:3032/album_list_view.tmpl
           content_type: application/json
           timeout: 10
      debug: true
```

What is being described?  First we have routes defined in our application. One route is `interesting_album_view` and it is composed from a two stage pipeline.
The first step is the `json_api` type (i.e. the content is fetched from PostgREST). The second stage is the template engine. Finally there is a `debug: true` prperty for this specific route. The `debug` attribute if true causes more verbose logging when the route is requested. This is helpful in debugging your pipeline.

When `newt` starts up it reads the YAML file. In the second stage above the type is set to `newtmustache`.  `newt` will package send the JSON obect object recieve from the previous (in the case PostgREST) JSON service to `newtmustache`. `newtmustache` understands the path `/album_list_view.tmpl` is a template name form the directory it read in at startup. It then take the JSON sent from PostgREST processing it through `/album_list_view.tmpl` returning a web page of results.

It is important to know `newt` only knows about the JSON (or other content) received from the previous service and the URL, content type and timeout for the next service. It's just a rely. No more or less. In the case above the template engine needs to understand the JSON it recieves and apply the template to it.

Let's say we decide to use Pandoc web service instead of `newtmustache`. How does this change the pipeline?  In this case what comes out of PostgREST doesn't include what Pandoc web server expects. Pandoc when running as a web server expect a JSON object that includes the template source to be used as well as the JSON object. One way would be to modify Pandoc server to use a setup similar to `newtmustache` but unless you want to hack Pandoc in Haskell you're kinda out of luck.  Another way would be to add a new stage in the pipeline that could read in all the Pandoc templates, resolve them and be ready to send them to Pandoc server.For the purposes of discussion let's say this web service (a new middleware) runs on port 3020 and accepts the template name like `newtmustache`. That pipeline could look something like this.


```yaml
routes:
    - id: interesting_album_view
      pipeline:
         - description: Contact PostgREST and get back the intersting album list
           service: GET http://localhost:3000/rpc/album_view
           content_type: application/json
           timeout: 15
         - description: |
             Construct a Pandoc web server JSON POST object that includes the
             template source and the JSON object from PostgREST. 
           service: POST http://localhost/3020/album_list_view.tmpl
           content_type: application/json
           timeout: 5
         - description: |
             Take the results from The template packager and send them to Pandoc.
           service: POST http://localhost:3030
           content_type: application/json
           timeout: 10
      debug: true
```

In this version we have a three stage pipeline. First stage gets some results from PostgREST. Second stage turns the results into
something Pandoc web service will understand.  The last stage is the Pandoc web service. Since debug is set true we can see each
stage of the pipeline accessed.

## Changes from the first prototype to the second.

- routes include a pipeline rather than fixed stages
- `newt` does less. It just routes data. It doesn't know how to package things for Pandoc web service
- `newtmustache` was created to allow the support of Mustache templates but also to simplify data packaging to avoid the object expected by Pandoc web service
- each pipeline stage has its own timeout setting

While there isn't a fixed limit to the number of stages in a pipeline you will want to keep the number limitted. While contacting a web service on localhost is generally very quick the time to round trip the communication still accoumulates. As a result it is recommend to stick to less than four stages and to explicitly set the timeout values based on your understanding of performance requirements.  If a stage times out the response will be generate with an HTTP error and the pipeline will not complete execution.

## Misc

If a requested route is not defined in our table then it looks for a static file matching the description and if that fails returns a 404. If the request is otherwise invalid based on what the router knows it'll return other HTTP status codes. For a request route to match it must match the resolved path, the HTTP method and included content type. If any of these don't match then the route is not considered a match and will return an appropriate HTTP status and code.

The Newt router will only support HTTP and run on localhost. This is the same approach taken by Pandoc server. It minimize configuration and also discourages use as a front end service (which is beyond the scope of this project).

This prototype does not support file uploads. In theory you could do some clever things with browser side JavaScript and store the contents in Postgres but I don't recommend that.

