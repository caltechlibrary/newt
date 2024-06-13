---
title: newtrouter(1) user manual | 0.0.8 a5cb6e7
pubDate: 2024-06-12
author: R. S. Doiel
---

# NAME

newtrouter

# SYNOPSIS

newtrouter YAML_CONFIG_FILE

# DESCRIPTION

**newtrouter** is a web service designed to work along side JSON API like that form with Postgres + PostgREST, and a template engine like Newt's Mustache template engine. **newtrouter** accepts a request, if it finds a matching route description it runs the request through a data pipeline of web services returning the results of the last one executed to the web browser or requester. It's just a data router that manages a pipeline of services for each defined request pattern.

In additional content routing newtrouter can also service out static resources. This is handy during development but less useful if you are using a front end web server such as a production setting.

**newtrouter**'s configuration uses a declarative model expressed in YAML.  It can also allow environment variables read at start up to be part of the data for mapping JSON data source requests. This is particularly helpful for supplying access credentials. You do not express secrets in the **newtrouter** YAML configuration file. This follows the best practice used when working with container services and Lambda like systems.

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

**newtrouter** is configured from a YAML file. The YAML should not include secrets. Instead you can pass these in via the environment variables identified the `.appliction.environment` property. Here's a summary of the Newt YAML syntax that **newtrouter** uses.

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
applications:
  router:
    port: 8010
    htdocs: htdocs
  template_engine:
    port: 8011
  postgres:
    namespace: blog
    port: 5432
    dsn: postgres://{PGUSER}:{PGPASSWORD}@localhost:5432/blog.yaml
  postgrest:
    app_path: postgrest
    conf_path: postgrest.conf
    port: 3000
  enviroment:
    - PGUSER
    - PGPASSWORD
models:
  - id: post
    description: A blog post or article
    elements:
      - type: text
        id: post_id
        attributes:
          name: post_id
          placeholdertext: e.g. /<YYYY>/<MM>/<DD>/<SLUG>
          title: (required) Enter the path for the blog entry with a unique slug
		  required: true
      - type: text
        id: title
        attributes:
          name: title
          title: (optional) Enter a title for your post
          placeholdertext: ex. My Blog Post for Today
      - type: text
        id: byline
        attributes:
          name: byline
		  title: (optional) Include a byline for your post
          placeholdertext: ex. By Jane Jones, 1912-12-12
      - type: textarea
        id: content
        attributes:
          name: content
          title: (required) Write your post here
          placeholdertext: ex. Something exciting happened today...
          required: true
      - type: date
        id: pubDate
        attributes:
          name: pubDate
          required: "true"
routes:
  - id: post_create
    request: GET /post_create
    description: Handle retrieving the webform for post create
    pipeline:
      - service: POST http://localhost:8011/post_create
        description: Display a post for create
  - id: post_create
    request: POST /post_create
    description: Handle form submission for post create
    pipeline:
      - service: POST http://localhost:3000/rpc/post_create
        description: Access PostgREST API for post create
      - service: POST http://localhost:8011/post_create_response
        description: This is an result template for post create
  - id: post_update
    request: GET /post_update/{oid}
    description: Handle retrieving the webform for post update
    pipeline:
      - service: GET http://localhost:3000/rpc/post_read/{oid}
        description: Retrieve post from PostgREST API before update
      - service: POST http://localhost:8011/post_update
        description: Display a post for update
  - id: post_update
    request: POST /post_update
    description: Handle form submission for post update
    pipeline:
      - service: PUT http://localhost:3000/rpc/post_update/{oid}
        description: Access PostgREST API for post update
      - service: POST http://localhost:8011/post_update_response
        description: This is an result template for post update
  - id: post_delete
    request: GET /post_delete/{oid}
    description: Handle retrieving the webform for post delete
    pipeline:
      - service: GET http://localhost:3000/rpc/post_read/{oid}
        description: Retrieve post from PostgREST API before delete
      - service: POST http://localhost:8011/post_delete
        description: Display a post for delete
  - id: post_delete
    request: POST /post_delete
    description: Handle form submission for post delete
    pipeline:
      - service: DELETE http://localhost:3000/rpc/post_delete/{oid}
        description: Access PostgREST API for post delete
      - service: POST http://localhost:8011/post_delete_response
        description: This is an result template for post delete
  - id: post_read
    request: POST /post_read
    description: Retrieve object(s) for post read
    pipeline:
      - service: GET http://localhost:3000/rpc/post_read/{oid}
        description: Access PostgREST API for post read
      - service: POST http://localhost:8011/post_read
        description: This template handles post read
  - id: post_list
    request: POST /post_list
    description: Retrieve object(s) for post list
    pipeline:
      - service: GET http://localhost:3000/rpc/post_list
        description: Access PostgREST API for post list
      - service: POST http://localhost:8011/post_list
        description: This template handles post list
templates:
  - id: post_create
    request: /post_create
    template: post_create_form.tmpl
    description: Display a post for create
  - id: post_create
    request: /post_create_response
    template: post_create_response.tmpl
    description: This is an result template for post create
  - id: post_update
    request: /post_update
    template: post_update_form.tmpl
    description: Display a post for update
  - id: post_update
    request: /post_update_response
    template: post_update_response.tmpl
    description: This is an result template for post update
  - id: post_delete
    request: /post_delete
    template: post_delete_form.tmpl
    description: Display a post for delete
  - id: post_delete
    request: /post_delete_response
    template: post_delete_response.tmpl
    description: This is an result template for post delete
  - id: post_read
    request: /post_read
    template: post_read.tmpl
    description: This template handles post read
  - id: post_list
    request: /post_list
    template: post_list.tmpl
    description: This template handles post list
~~~



