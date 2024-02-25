---
title: "newtrouter(1) user manual | 0.0.7-dev b865091"
pubDate: 2024-02-25
author: "R. S. Doiel"
---

# NAME

newtrouter

# SYNOPSIS

newtrouter YAML_CONFIG_FILE

# DESCRIPTION

**newtrouter** is a web service designed to work along side JSON API like that form with Postgres + PostgREST, and a template engine like Newt Mustache. **newtrouter** accepts a request, if it finds a matching route description it runs the request through a data pipeline of web services returning the results of the last one executed to the web browser or requestor. It's just a data router that manages a pipeline of services for each defined request pattern.

In additional content routing newtrouter can also service out static resources. This is handy during development but less useful if you are using a front end web server such as a production setting.

**newtrouter**'s configuration uses a declaritive model expressed in YAML.  It can also allow environment variables read at start up to be part of the data for mapping JSON data source requests. This is particularly helpful for supplying access credentials. You do not express secrets in the **newtrouter** YAML configuration file. This follows the best practice used when working with container services and Lambda like systems.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-dry-run
: Load YAML configuration and report any errors found

# YAML_CONFIG_FILE

**newtrouter** is configured from a YAML file. The YAML should not include secrets. Instead you can pass these in via the enviroment variables identifified tthe `.appliction.environment` property. Here's a summary of the Newt YAML syntax that **newtrouter** uses.

Top level properties for **newtrouter** YAML.

application
: (optional: newtrouter, newtgenerator, newtmustache) holds the run time configuration used by the Newt web service and metadata about the application you're creating.

routes
: (optional: newtrouter, newtgenerator) This holds the routes for the data pipeline (e.g. JSON API and template engine sequence)

## the "application" property

The application properties are optional.

port
: (optional: newtrouter, newtmustache) default is This port number the Newt web services uses to listen for request on localhost

htdocs
: (optional: newtrouter only) Directory that holds your application's static content

environment
: (optional: newtrouter, newtmustache) this is a list of operating system environment that will be available to routes. This is used to pass in secrets (e.g. credentials) need in the pipeline

Routes hosts a list of request descriptions and their data pipelines. This property is only used by Newt router and Newt code generator.

### a route object

`id`
: (required) This identifies the pipeline. It maybe used in code generation. It must conform to variable name rules[^21]

`description`
: (optional, encouraged) This is a human readable description of what you're trying to accomplish in this specific route. It may be used in comments or by documentation generators.

`request [METHOD ][PATH]`
: (required) This is a string that expresses the HTTP method and URL path to used to trigger running the data pipeline. If METHOD is not provided it will match using just the path. This is probably NOT what you want. You can express embedded variables in the PATH element. This is done by using single curl braces around a variable name. E.g. `GET /items/{item_id}` would make `item_id` available in building your service paths in the pipeline. The pattern takes up a whole path segment so `/blog/{year}-{month}-{day}` would not work but `/blog/{year}/{month}/{day}` would capture the individual elements. The Newt router sits closely on top of the Go 1.22 HTTP package route handling. For the details on how Go 1.22 and above request handlers and patterns form see See <https://tip.golang.org/doc/go1.22#enhanced_routing_patterns> and <https://pkg.go.dev/net/http#hdr-Patterns> for explanations.

`pipeline`
: (required) this is a list of URLs to one or more web services visible on localhost. The first stage to fail will abort the pipeline returning an HTTP error status. If done fail then the result of the last stage it returned to the requesting browser.

`debug`
: (optional) if set to true the `newtrouter` service will log verbose results to standard out for this specific pipeline

#### a pipeline object

A pipeline is a list of web services containing a type, URL, method and content types

`service [METHOD ][URL]`
: (required) The HTTP method is included in the URL The URL to be used to contact the web service, may contain embedded variable references drawn from the request path as well as those passed in through `.application.environment`.  All the elements extracted from the elements derived from the request path are passed through strings. These are then used to construct a simple key-value object of variable names and objects which are then passed through the Mustache template representing the target service URL. 

`description`
: (optional, encouraged) This is a description of what this stage of the pipe does. It is used when debug is true in the log output and in program documentation.

`timeout`
: (optional) Set the timeout in seconds for receiving a response from the web server. Remember the time spent at each stage is the cumulative time your browser is waiting for a response. For this reason you may want to set the timeout to a small number.

# EXAMPLES

Running **newtrouter** with a YAML configuration file called "blog.yaml"

~~~
newtrouter blog.yaml
~~~

An example of a YAML file describing blog like application based on Postgres+PostgREST.

~~~
application:
  htdocs: htdocs
  environment:
    - DB_USER
    - DB_PASSWORD
#
# Postgres+PostgREST is listening on port 3000
# Newt Mustache template engine is listening on port 3032
#
# DB_USER and DB_PASSWORD required to access the PostgREST JSON API
# so is passed in via the environment.
routes:
  - id: retrieve_all_posts
    request: GET /archives
    description: This route returns the full blog content
    pipeline:
      - description: |
          Retrieve the blog posts order by descending date
        service: GET http://{DB_USER}:{DB_PASSWORD}@localhost:3000/rpc/view_all_posts
      - description: render the posts using the list_posts.tmpl
        service: POST http://localhost:3032/list_posts.tmpl
  - id: retrieve_year posts
    request: GET /{year}
    description: This route retrieves all the posts in a specific year
    pipeline:
      - description: Retrieve the posts for a specific year
        service: GET http://{DB_USER}:{DB_PASSWORD}@localhost:3000/rpc/year_posts/{year}
      - description: Turn the JSON list into a web page.
        service: POST http://localhost:3032/list_posts.tmpl
  - id: retrieve_month_posts
    request: GET /{year}/{month}
    description: This route retrieves all the posts in a specific year/month
    pipeline:
      - description: Retrieve the posts in the DB for year/month
        service: GET http://{DB_USER}:{DB_PASSWORD}@localhost:3000/rpc/month_posts/{year}/{month}
      - description: Transform monthly list into web page
        service: POST http://localhost:3032/list_posts.tmpl
  - id: retrieve_day_posts
    request: GET /{year}/{month}/{day}
    description: Retrieve all the posts on a specific day
    pipeline:
      - description: Retrieve the posts happending on year/month/day
        service: GET http://{DB_USER}:{DB_PASSWOR}@localhost:3000/rpc/day_posts/{year}/{month}/{day}
      - description: Transform monthly list into web page
        service: POST http://localhost:3032/list_posts.tmpl
  - id: retrieve_recent_posts
    request: GET /
    description: This route retrieves the recent 10 posts.
    pipeline:
      - description: Use the recent_post view to retrieve the recent posts in descending date order
        service: GET http://{DB_USER}:{DB_PASSWORD}@localhost:3000/rpc/recent_posts
      - description: Take the JSON for recent posts and turn it into a web page
        service: GET http://localhost:3032/list_posts.tmpl
  - id: retrieve_a_post
    request: GET /{year}/{month}/{day}/{title-slug}
    description: Retrieve a specific host and display it
    pipeline:
      - description: retrieve the requested post from the blog path
        service: GET http://{DB_USER}:{DB_PASSWORD}@localhost:3000/{year}/{month}/{day}/{title-slug}
      - description: Turn the JSON into a web page
        service: GET http://localhost:3032/blog_post.tmpl
~~~



