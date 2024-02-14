---
title: "tmplbndl(1) user manual | 0.0.6 f419528"
pubDate: 2024-02-13
author: "R. S. Doiel"
---

# NAME

`tmplbndl`

# SYNOPSIS

`tmplbnld [OPTIONS] YAML_CONFIG_FILE`

# DESCRIPTION

`tmplbndl` is web service that takes a JSON object and bundles it for use with a Pandoc's web service. It requires a configuration file that maps requests to a template and any other options. When the `tmplbndl` web service is running it accepts a request object, processes and returns a HTTP response. If an error is encounter than the response is an HTTP error response. If no error is encounterd then a JSON object is assembled suitable for Pandoc web service to process.

In the first Newt prototype Pandoc web service integrated into `newt`. In the second Newt prototype the `newt` web service is more generalized and has adopted a pipeline approach. There are no assumptions tieing it to a particular template engine. This means the output from Solr or Postgres+PostgREST needs to be transformed into a new JSON Object to be able to use Pandoc web service. `tmplbndl`'s job is to sit in the pipeline between a JSON data source and Pandoc web service.

`tmplbndl` is a web service configured from a YAML file. It has a syntax specific to its purposes. While `tmplbndl` is focused on producing a JSON object for requests to Pandoc web server it also services as an example for other types of bundlers.

Basic idea.

1. At start up `tmplbndl` reads a YAML file
2. It builds request handlers mapping URL requests with a template's source
3. It maps any varaibles described in the request pattern
3. It starts up is web service

When the web service is active 

1. `tmplbndl` receives a request in the form of a HTTP method and URL
2. It checks it handlers for a match. If none are found a 404 is returned.
3. With a matched request, it builds an object setting the values of `.text` and `.template` accordingly. If there are embeded variables in the pattern then those are mapped into the resulting object as well.

### Features

- templates are read in at startup and are retained in memory bound to the request path
- the object composed from the request may also include options supported in Pandoc expected request object

This makes the full power of the Pandoc web service available to you Newt application pipeline.

# OPTIONS

-help
: Display this help message

-license
: Display license information

-version
: Display version information

-port
: Run service on designated port

## Template bundler's YAML

port
: (integer, defaults to 3029) the port number the service should listen on

templates
: (list of template paths and options)

### a template description

An individual template description has four properties.

request
: (string, required) This is the METHOD and PATH that `tmplbndl` will listen for to map this template. You should normally use a POST method to send content to be bundled. Go 1.22 HTTP handler patterns can be used and can be used to override the settings in `.options`.

template
: (string, optional) This is path to the primary template used required for this request. The source if read and associated with this request signature. If template is not provided then Pandoc server will default to its internal templates

options
: (object, optional) These will become root object properties in the returned JSON. This is where you would specify the Pandoc object properties for controlling the the transform of the text submitted[^1]. NOTE: the `.text` property is replaced by the text received by `tmplbndl` and `.template` will be replaced with the resolved source read in at `tmplbndl` startup.


## Example YAML configuration

This example shows fix different template options. The first three applies a custom `page.tmpl` in different ways. In the last three the default Pandoc template is setup.

~~~yaml
port: 3029
templates:
  - request: "POST /custom_page"
    template: page.tmpl
    options:
      from: markdown
      to: html5
      standalone: true
      title: This is the custom template with this title
  - request: "POST /custom_page_with_title/{title}"
    template: page.tmpl
    options:
      from: markdown
      to: html5
      standalone: true
      title: This title is overwritten by the one in the request
  - request: "POST /custom_page_include"
    template: page.tmpl
    options:
      from: markdown
      to: html5
      standalone: false
  - request: "POST /default_html5"
    options:
      from: markdown
      to: html5
      standalone: true
      title: A Page using the default template
  - request: "POST /default_html5_with_title/{title}"
    options:
      from: markdown
      to: html5
      standalone: true
      title: This title is replaced by the title in the URL
  - request: "POST /default_html5_include"
    options:
      from: markdown
      to: html5
      standalone: false
~~~

[^1]: The Pandoc web service is documented in it's man page, see <https://github.com/jgm/pandoc/blob/main/doc/pandoc-server.md> for the current version

