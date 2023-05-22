
# Newt URL Router

To extend the Postgres+PostgREST to support server side web page assembly with Pandoc server we need to be able to map a requested URL path to both PostgREST and Pandoc server. The information required for this is what I am exploring in this document.

1. request information
2. PostgREST API request information
3. Pandoc URL and template data

I conceptualize this with the following.

- request
    - request route
    - request HTTP method
    - request content type
- data api
    - api url
    - api HTTP method
    - api content type
- pandoc request
    - pandoc port number
    - pandoc server query string (to set pandoc options)
    - pandoc template name

This can be easily represented in a table (e.g. a CSV file) and managed from
a spreadsheet. Here's the columns need in a csv file.

```csv
req_route, req_method, req_content_type, api_url, api_method, api_content_type, pandoc_port, pandoc_query_string, pandoc_template
```

Newt itself needs to know four things to run.

1. url to listen on (e.g. http://localhost:4242)
2. path to CSV file with routes
3. path to the directory for static content
4. The names of environment varaibles to store things like PostgREST user and password information when forming a URL for a data API request


## What is a request route?

A route is a URL path similar to a Unix file path. It main be an explicit route or one that describes an expression expressed in a [PathDSL](pathdsl.md "path domain specific language").  The PathDSL enables a request route to be parsed and transformed into a data API request and Pandoc template.

~~~
/about.html
/blog/{year year}/{month month}/{day day}/
~~~

Routes can contain variables that are re-used in forming a data API request or used by Pandoc server if that is defined for the route.

## What does a data API URL look like?

While Newt is inspired by Postges+PostgREST microservice it can work with any data source that can be reached by a URL (e.g. Opensearch).  The data API URL can be composed as a literal string from the variables captured in the request route. It can also include variables defined in the environment (e.g. so you don't have to hardcore a username, password combination in your routes CSV file).

Here's an example of a data URL that uses the route information from our "blog" path as well as
two variables defined in the environment (i.e. DB_USER, DB_PASSWD).

~~~
http://{DB_USER}:{DB_PASSWD}@localhost:3000/blog?year={year}&month={month}&day={day}
~~~

In the template values for "{year}", "{month}", "{day}" came from our request route while
"{DB_USER}" and "{DB_PASSWD}" came from the environment.


## The routes CSV columns defined

api_url, api_method, api_content_type, pandoc_port, pandoc_query_string, pandoc_template

req_route
: route description (see [Path DSL concept](pathdsl.md))

req_method
: required HTTP method (e.g. GET, POST, HEAD, PATCH, PUT)

req_content_type
: requested content type (e.g. "text/html")

api_url
: the URL template used to contact a data API, e.g. PostgREST API or other services like Solr or Opensearch

api_method
: the HTTP method used when contacting the data source

api_content_type
: the mime type requested from the data source (e.g. "applicaiton/json")

pandoc_port
: the port number the Pandoc server is running on

pandoc_query_string
: this is the query string (e.g. `from=json&to=markdown`) used by the Pandoc server to process the request

pandoc_template
: the filename that will be loaded at start up of Newt and sent with the Pandoc service request for Pandoc to process

resp_headers
: (optional) additional response headers returned to by Newt based on the output data. These are expression as a JSON object (e.g. `{"Content-Type": "text/plain"})

These fields are represented as a row in a table. In the Newt prototype these are included in a CSV file. Concievable other sources could be used like a SQL table or Excel file.

## Newt prototype behavior

When Newt router starts up it reads a configuration file, it then reads the routes CSV file indicated in the configuration. It remembers any additional enviroment variables indicated in the configuration. When that is done it starts listening for requests and performs the routing dance of rejecting the request for undefind or invalid routes or contacting the JSON data API and if needed Pandoc server for processing the request.

## Misc

If a requested route is not defined in our table then a 404 is returned. If the request is otherwise invalid based on what the router knows about the request then a suitable HTTP error should be returned (e.g. The URL might be defined but the HTTP method might not match).

The Newt router will only support http and run on localhost. This is the same approach taken by Pandoc server. It minimize configuration and also discorages use as a front end service (which is beyond the scope of this project).

The prototype does not support file uploads. That's something that could be added in the future and would probably function by pass data through to anther service like MinIO or S3.
