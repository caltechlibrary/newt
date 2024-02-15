---
title: "pdbundler(1) user manual | 0.0.6 f65586f"
pubDate: 2024-02-15
author: "R. S. Doiel"
---

# NAME

`pdbundler`

# SYNOPSIS

`pdbundler [OPTIONS] YAML_CONFIG_FILE`

# DESCRIPTION

`pdbundler` is web service that takes a JSON object and bundles it for use with a Pandoc's web service. It requires a configuration file that maps requests to a template and options. When the `pdbundler` web service is running it accepts JSON and returns JSON suitable to POST to Pandoc web service. If an error is encounter than the response is an HTTP error response.

In the first prototype the Pandoc web service was integrated into the Newt data router. In the second prototype the Newt data router supports a generalized data pipeline. The Newt data router makes no assumptions tieing it to a template engine. This means the output from Solr or Postgres+PostgREST needs to be transformed into a new JSON expression to be usabled by Pandoc web service. `pdbundler` provides this support. It should run in the pipeline between your JSON data source and just before Pandoc web service.

`pdbundler` is a web service configured from a YAML file. It has a syntax specific to its purposes. `pdbundler` is an example of a narrowly focuses web service or micro service.

Initialization process.

1. `pdbundler` read the YAML configuration file
2. Builds request handlers mapping URL requests with a template's source
3. Start up an http service on local host listening on the specified port for requests.

When the web service is active 

- `pdbundler` receives a request in the form of a HTTP method and URL
- It checks it handlers for a match. If none are found a 404 is returned.
- With a matched request, it builds an object setting the values of `.text` and `.template` accordingly[^1]. 

### Features

- templates are read in at startup and are retained in memory bound to the request path
- the template descriptions can include Pandoc configuration objects to send to Pandoc web service
- configuration options can be set through variables in the URL Path of the request associated with the template

These features are intended to expose the capabilities of Pandoc web service.

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

`request: [METHOD] PATTERN`
: (required) This binds a request to a template. METHOD is usually either a GET or a POST. If METHOD is not defined then a POST is assumed. PATTERN is usually the path to associated with a template name. The PATTERN supports the syntax described in Go 1.22 `http` package, see <https://pkg.go.dev/net/http@master#hdr-Patterns>. Variables expressed in PATTERN are merged into the resulting JSON object. They overwrite what is set in the `.options` settings..

template
: (string, optional) This is path to the primary template used required for this request. The source if read and associated with this request signature. If template is not provided then Pandoc server will default to its internal templates

options
: (object, optional) Are used to form the root object properties in the returned JSON. This is where you would specify the Pandoc web service processing options. The options control the transformation of the text submitted[^2]. NOTE: the `.text` property is replaced by the text received by `pdbundler` and `.template` will be replaced with the resolved source read in at `pdbundler` startup.

## Example YAML configuration

This example shows six different template options. The first three apply one custom `page.tmpl` in different ways. In the last three do the same bu assume the default Pandoc template.

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

[^1]: See <https://pandoc.org/MANUAL.html#running-pandoc-as-a-web-server> for basic information about Pandoc web service. Explore the website to get familiar with Pandoc and its web service.

[^2]: On <https://pandoc.org/MANUAL.html#running-pandoc-as-a-web-server> there is a link, "pandoc-server", to the current release manual page. This covers the details of using Pandoc web service including the POST JSON object setup.

