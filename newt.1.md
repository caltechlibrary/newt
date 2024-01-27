---
title: "newt(1) user manual | 0.0.6 cef35c8"
pubDate: 2024-01-26
author: "R. S. Doiel"
---

# NAME

newt

# SYNOPSIS

newt [CONFIG_FILE]

# DESCRIPTION

**newt** is a microservice designed to work along side Postgres, PostgREST, and Pandoc server. It provides URL routing and data flow between the microservices based on a list of "routes" described in a YAML file.  **newt** first sends requests to a JSON data source then processes the results via Pandoc running as a web service.  While **newt** was created to work specifically with PostgREST it can actually talk to any JSON data source that can be specified by a URL and HTTP method (e.g. Solr, Elasticsearch, Opensearch).

Before contacting the JSON data source the request path and any form data is validated based on the request path and any variables defined for the route. If validation is successful values are extracted from the request path along with form data. These are then used to make a request to a JSON data source (e.g. PostgREST) described in our route definition.

When the data source replies the results can be fed through Pandoc running as a web service based on a template filename associated with the route. If no template file is specified then the results of the JSON data source is passed directly back to the web browser (or requesting service).

Additionally **newt** can function as a static content web service.  This is handy when developing a **newt** based project. A typical setup might include running Postgres, PostgREST and Pandoc server along with **newt** as you develop your project. Since **newt** always works as a "localhost" service you will need to proxy to it when deploying to a production setting (e.g. via Apache2 or NginX).

**newt**'s configuration uses a declaritive model expressed in YAML.  It can also allow environment variables read at start up to be part of the data for mapping JSON data source requests. This is particularly helpful for supplying access credentials. You do not express secrets in the **newt** YAML configuration file. This follows the best practice used when working with container services and Lambda like systems.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-dry-run
: Load YAML configuration and report any errors found

# CONFIGURATION

**newt** looks for four environment variables at startup.

NEWT_PORT
: (optional) The port Newt will listen for requests

NEWT_CONFIG
: (optional) The name of a YAML file holding newt configuration

NEWT_ENV
: (optional) The names of environment variables **newt** can make available when setting up route handling.

NEWT_HTDOCS
: (optional) The directory holding static content that Newt will serve alonside any defined data routes specified in the configuration.

Example environment expressed in sh.

~~~
export NEWT_PORT="8001"
export NEWT_ROUTES="routes.yaml"
export NEWT_ENV="DB_NAME;DB_USER;DB_PASSWORD"
export NEWT_HTDOCS="/var/local/html/htdocs"
~~~

These can also be expressed directly in the YAML configuration file using the following attribute names -- "port", "routes", "env", and "htdocs". 

~~~
port: 8001
routes: routes.yaml
env: [ "DB_NAME", "DB_USER", "DB_PASSWORD" ]
htdocs: /var/local/html/htdocs
~~~

The environment will load first then the configuration file.  The configuration file takes precedence over the environment.

**newt** does not contain secrets.  These should be passed in via the environment. This follows the practices that have become commonplace when using containers and Lamdda like services.

Avoiding secrets allows your routes in the YAML file to be safely included in your project source repository along side any Pandoc templates and SQL source included in your project's source code repository.

# Routing data

For **newt** to function as a data router it needs information about which requests will be serviced and how to map them to a JSON source before (optionally) sending to Pandoc.

The routes are held in YAML file under the "routes" attribute.  The following attributes are supported for a route.

var
: One or more variable name and type pairs used in request path or form data. The types allow for data validation.

req_path
: This is the URL path to watch for incoming requests, it may be a literal path or one containing variable declarations used in forming a HTTP call to a JSON source.

req_method
: This is the HTTP method to listen for -- "GET", "POST", "PUT", "PATCH" or "DELETE".

api_url
: This is the URL used to connect to the JSON data source (e.g. PostgREST, Solr, Elasticsearch). It may contain defined variables embedded in the request path or form form data.

api_method
: This is the HTTP method used to access the JSON data source. Maybe "OPTIONS", "GET", "POST", "PUT", "PATCH" or "DELETE"

api_content_type
: This is the HTTP content type string to send with your JSON data source request, typically it is "application/json". 

template
: If included Newt will load the Pandoc template file into memory and use it when results are returned from a JSON data source. The data is provided to the Pandoc template as part of the "body" pandoc template variable.

render_port (optional)
: This is the port number used for the render engine, e.g. Pandoc server or **newtmustache** render. It not specified then port 3030 is assumed.

res_headers
: This is any additional HTTP headers you want to send back to the client.

# Defining variables

Each route can have its own associated set of variables. Variables are "typed".  The code for type definitions includes validation. When a variable is detected in a request path or form data it is vetted using it's associated type. Only if the variables past validation are they allowed to be used to assemble a request to a JSON data source. 

Variables are defined in the YAML "var" attribute. Here's an example "var" definition defining three form variables for a route. The variable names are "bird", "place" and "sighted" with the types "String", "String" and "Date".

~~~yaml
var: [ "bird": "String", "place": "String", "sighted": "Date" ]
~~~

If a web browser injected additional form values they would not get passed along via the JSON data API request, they would be ignored. This is part of the declaritive approach for defining Newt's behavior.

The variables "bird", "place" and "sighted" can be used when specifying a request route.  Variables that are defined in a route are delimited by an opening '${' and closing '}'.  In the following example the URL could represent browsing birds by place and date sighted.

~~~
/birds/${place}/${sighted}
/birds/${place}/${sighted}/${bird}
~~~

This might be used to make a request to a JSON data source (e.g. PostgREST) like this.

~~~
https://localhost:3000/sightings?bird=${bird}&place=${place}&sighted=${sighted}
~~~

The result of the JSON source request could then be processed with a Pandoc template to render an HTML page.

# Variable types

String
: Any sequence of characters. If the variabe is embedded in a path then "/" will be used to delimited path parts and would not be passed into the variables value.

Date
: (default) A year, month, day string like 2006-01-02

Date 2006
: A four digit year (e.g. 2023)

Date 01
: A two digit month (e.g. "01" for January, "10" for October)

Date 02
: A two digit day (e.g. "01" for the first, "11" for the eleventh)

Basename
: A file's basename (filename without an extension)

Extname
: A file's extension (e.g. ".html", ".txt", ".rss", ".js")

Isbn10
: An ten digit ISBN

Isbn13
: A thirteen digit ISBN

Isbn
: An ISBN (either 10 ro 13 digit)

Issn
: An ISSN

DOI
: A DOI (digital object identifier)

Isni
: An ISNI

ORCID
: An ORCID identifier
 
NOTE: The current names associated with types will likely change
as the prototype **newt** evolves. It is planned for them to be
stable if and when we get to a v1 release (e.g. when we're out of the
prototype phase).

# Pandoc, Pandoc templates

Values received from the JSON data source are passed to the Pandoc template bound to the variable name "data". This is done by taking the JSON recieved and forming a front matter document that is then used alongside Pandoc template in the POST request made to Pandoc running in server mode. See <https://pandoc.org/pandoc-server.html> and <https://pandoc.org/MANUAL.html#templates> for details.

# EXAMPLES

Running **newt** with a YAML configuration file called "blog.yaml"

~~~
newt blog.yaml
~~~

An example of a YAML file describing blog like application based on Postgres+PostgREST.

~~~
env: [ "DB_USER", "DB_PASSWORD" ]
htdocs: htdocs
routes:
  - var: [ "yr": "Date 2006", "mo": "Date 01", "dy": "Date 02" }
    req_path: "/blog/${yr}/${mo}/${dy}"
    req_method: GET
    api_url: "http://${DB_USER}:${DB_PASSWORD}@localhost:3000/posts?year=${yr}&month=${mo}&day=${dy}"
    api_method: GET
    api_content_type: "application/json"
    pandoc_template: article_list.tmpl
    res_headers: { "content-type": "text/html" }
  - var: [ "yr": "Year", "mo": "Month", "dy": "Day" }
    req_path: "/blog/${yr}/${mo}/${dy}/${title-slug}"
    req_method: GET
    api_url": "http://${DB_USER}:${DB_PASSWORD}@localhost:3000/posts?year=${yr}&month=${mo}&day=${dy}&title-slug=${title-slug}"
    pandoc_template: article.tmpl
    res_headers: { "content-type": "text/html" }
~~~



