package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	// Caltech Library Packages
	"github.com/caltechlibrary/newt"
)

var (
	helpText = `---
title: {app_name}(1) user manual | {version} {release_hash}
pubDate: {release_date}
author: R. S. Doiel
---

# NAME

{app_name}

# SYNOPSIS

{app_name} YAML_CONFIG_FILE

# DESCRIPTION

**{app_name}** is a web service designed to work along side JSON API like that form with Postgres + PostgREST, and a template engine like Newt's Handlebars template engine. **{app_name}** accepts a request, if it finds a matching route description it runs the request through a data pipeline of web services returning the results of the last one executed to the web browser or requester. It's just a data router that manages a pipeline of services for each defined request pattern.

In additional content routing {app_name} can also serve static resources. This is handy during development but less useful if you are using a front end web server such as a production setting.

**{app_name}**'s configuration uses a declarative model expressed in YAML.  It can also allow environment variables read at start up to be part of the data for mapping JSON data source requests. This is particularly helpful for supplying access credentials. You do not express secrets in the **{app_name}** YAML configuration file. This follows the best practice used when working with container services and Lambda like systems.

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

**{app_name}** is configured from a YAML file. The YAML should not include secrets. Instead you can pass these in via the environment variables identified the ` + "`" + `.appliction.environment` + "`" + ` property. Here's a summary of the Newt YAML syntax that **{app_name}** uses.

Top level properties for **{app_name}** YAML.

services
: (optional) the run time service configuration used to compose your Newt application

routes
: (optional: newtrouter, newtgenerator) This holds the routes for the data pipeline (e.g. JSON API and template engine sequence)

## the "services" property

The services property is a list of services and how to run them. Each service may have the following properties.

## the "routes" property

Routes hosts a list of request descriptions and their data pipelines. This property is only used by Newt router and Newt code generator.

### a route object

` + "`" + `id` + "`" + `
: (required) This identifies the pipeline. It maybe used in code generation. It must conform to variable name rules[^21]

` + "`" + `description` + "`" + `
: (optional, encouraged) This is a human readable description of what you're trying to accomplish in this specific route. It may be used in comments or by documentation generators.

` + "`" + `request [METHOD ][PATH]` + "`" + `
: (required) This is a string that expresses the HTTP method and URL path to used to trigger running the data pipeline. If METHOD is not provided it will match using just the path. This is probably NOT what you want. You can express embedded variables in the PATH element. This is done by using single curl braces around a variable name. E.g. ` + "`" + `GET /items/{item_id}` + "`" + ` would make ` + "`" + `item_id` + "`" + ` available in building your service paths in the pipeline. The pattern takes up a whole path segment so ` + "`" + `/blog/{year}-{month}-{day}` + "`" + ` would not work but ` + "`" + `/blog/{year}/{month}/{day}` + "`" + ` would capture the individual elements. The Newt router sits closely on top of the Go 1.22 HTTP package route handling. For the details on how Go 1.22 and above request handlers and patterns form see See <https://tip.golang.org/doc/go1.22#enhanced_routing_patterns> and <https://pkg.go.dev/net/http#hdr-Patterns> for explanations.

` + "`" + `pipeline` + "`" + `
: (required) this is a list of URLs to one or more web services visible on localhost. The first stage to fail will abort the pipeline returning an HTTP error status. If done fail then the result of the last stage it returned to the requesting browser.

` + "`" + `debug` + "`" + `
: (optional) if set to true the ` + "`" + `{app_name}` + "`" + ` service will log verbose results to standard out for this specific pipeline

#### a pipeline object

A pipeline is a list of web services containing a type, URL, method and content types

` + "`" + `service [METHOD ][URL]` + "`" + `
: (required) The HTTP method is included in the URL The URL to be used to contact the web service, may contain embedded variable references drawn from the request path as well as those passed in through ` + "`" + `.application.environment` + "`" + `.  All the elements extracted from the elements derived from the request path are passed through strings. These are then used to construct a simple key-value object of variable names and objects which are then passed through the Handlebars template representing the target service URL.

` + "`" + `description` + "`" + `
: (optional, encouraged) This is a description of what this stage of the pipe does. It is used when debug is true in the log output and in program documentation.

` + "`" + `timeout` + "`" + `
: (optional) Set the timeout in seconds for receiving a response from the web server. Remember the time spent at each stage is the cumulative time your browser is waiting for a response. For this reason you may want to set the timeout to a small number.

# EXAMPLES

Running **{app_name}** with a YAML configuration file called "blog.yaml"

~~~
{app_name} blog.yaml
~~~

An example of a YAML file describing blog like application based on Postgres+PostgREST.

~~~
services:
  - name: router
    path: ndr
    port: 8010
    htdocs: htdocs
  - name: template_engine
    path: nte
    port: 8011
  - name: postgres
    path: postgres
    namespace: blog
    port: 5432
    dsn: postgres://{PGUSER}:{PGPASSWORD}@localhost:5432/blog.yaml
    enviroment:
      - PGUSER
      - PGPASSWORD
  - name: postgrest
    path: postgrest
    conf_path: postgrest.conf
    port: 3000
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


`

	showHelp    bool
	showLicense bool
	showVersion bool
)

func main() {
	appName := path.Base(os.Args[0])
	// NOTE: The following variables are set when version.go is generated
	version := newt.Version
	releaseDate := newt.ReleaseDate
	releaseHash := newt.ReleaseHash
	fmtHelp := newt.FmtHelp

	// Standard Options
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "version", false, "display version")

	// App option(s)
	dryRun, port, htdocs, verbose := false, 0, "", false
	flag.BoolVar(&dryRun, "dry-run", dryRun, "evaluate configuration and routes but don't start web service")
	flag.IntVar(&port, "port", port, fmt.Sprintf("set the port %s will listen on", appName))
	flag.BoolVar(&verbose, "verbose", verbose, fmt.Sprintf("run %s in verbose debug mode", appName))
	flag.StringVar(&htdocs, "htdocs", htdocs, "set htdocs directory holding static files")

	// We're ready to process args
	flag.Parse()
	args := flag.Args()

	in := os.Stdin
	out := os.Stdout
	eout := os.Stderr

	if showHelp {
		fmt.Fprintf(out, "%s\n", fmtHelp(helpText, appName, version, releaseDate, releaseHash))
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintf(out, "%s\n", newt.LicenseText)
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintf(out, "%s %s %s\n", appName, version, releaseHash)
		os.Exit(0)
	}
	os.Exit(newt.RunRouter(in, out, eout, args, dryRun, port, htdocs, verbose))
}
