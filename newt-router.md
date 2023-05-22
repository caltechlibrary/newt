
# Newt URL Router

To extend the Postgres+PostgREST to support server side web page assembly with Pandoc server we need to be able to map a requested URL path to both PostgREST for content and a template for processing the data with Pandoc server before returning an assembled web page to the front end web server. What information is needed to define this behavior.

Newt URL router leverages your existing PostgREST configuration with
the optional additional value `pandoc-server-port` if you Pandoc server is
not running on a standard port. Route configuration is loaded from 
Postgres+PostgREST. Routes can be manage in a CSV file which is used to 
update a route table in Postgres.

## The routes table fields

route
: route description (see [RouteDSL concept](routedsl.md))

req_method
: required HTTP method (e.g. GET, POST, HEAD, PATCH, PUT)

req_content_type
: requested content type (e.g. "text/html")

req_data
: a request submitted data (e.g. contents of a POST)

api_URL
: Data API API request url (e.g. PostgREST, Solr, Opensearch)

api_method
: http method to use to access the api data source 

api_content_type
: the http content type to use with the api data source

pandoc_template
: pandoc_template name (to use to process with Pandoc server, template will be read at Newt startup), maybe empty

pandoc_query
: pandoc_query a URL query string used to configure the Pandoc server request (see Pandoc user manual)

resp_headers
: additional response headers that might be useful to send (e.g. CORS, web token management)

These fields can be represented as a row in a table (CSV, SQL), one row per route pattern.

If a requested route is not defined in our table then a 404 is returned. It is useful to define a catch all to allow support for static content. That route description might look something like `/{STATIC_FILE_PATH dirname}/{FILENAME basename}{EXT extname}`. If api_path is empty and pandoc template and doc_settings is empty then the static content would just pass through Newt. This would allow you to define a static route that processed Markdown documents into HTML via Pandoc server but also serve out things like images.

When Newt router starts up it reads a configuration file to know hostname/port and where to find the related PostgREST server, Pandoc server, htdocs directory, and pandoc template directory and URL to table of routes (maybe a file path to CSV file, SQLite3 database or URL to PostgREST). When the configuration is read the route table is scanned and any pandoc templates are read into memory. The router then listens for requests and dispatches them based on the routes received, http method and content types requested.

The Newt router is intended to support http only and run on localhost. This is the same approach adopted by PostgREST and Pandoc server. In a development setting that means you only need to run Postgres, PostgREST, Pandoc server and Newt. Newt only provides routing for PostgREST, Pandoc server and static files. 


## Someday, maybe features

The follow could compliment routing for PostgREST to build more comprehensive web applications.

- gateway for file upload storage service (e.g. gateway to Minio or S3)
