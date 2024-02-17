---
title: "newtmustache(1) user manual | 0.0.7-dev 087e3b7"
pubDate: 2024-02-17
author: "R. S. Doiel"
---

# NAME

newtmustache

# SYNOPSIS

newtmustache [OPTIONS] YAML_CONFIG_FILE

# DESCRIPTION

**newtmustache** is a web service that provides a Mustache template rendering inspired by Pandoc server.

Unlike Pandoc web server, `newtmustache` expects a YAML_CONFIG_FILE. The format is
described below. That file specifics the request to template mapping along with any ancillary informtation
to merge into the submitted object for the template to process. In additional to options expressed in
the configuration the pattern describing the request can also be merged into the JSON objects the template
will process.

`newtmustache` expects a JSON object to process. This means you should normally use a POST
in defining your request pattern to be match.  If no method is specified then a POST method will be
assumed.

The request expression is based on the Go 1.22 patterns used in it's `http` package.
<https://pkg.go.dev/net/http@master#hdr-Patterns>. The only modification is if not method is included
a POST instead of a GET will be assumed.

Like Pandoc web service `newtmustache` does not normally log requests. It's a quick transaction.
If you want to debug your templates use the verbose option or turn on debug for specific requests.

# OPTIONS

The following options are supported by **newtmustache**.

-h, -help
: display this help message

-license
: display the software license

-version
: display version information

-port NUMBER
: (default is port is 3032) Set the port number to listen on

-timeout SECONDS
: Timeout in seconds, after which a template rendering is aborted.  Default: 3.

-verbose
: If set provide verbose debugging output for requests

# The templates

Mustache templates are documented at <https://mustache.github.io>. The template engined
used is based on Go package <https://github.com/cbroglie/mustache>.

## Features

- Newt Mustache only runs on localhost at the designated port (default is 3032).
- Templates are read in at startup and are retained in memory bound to the request path.
- No addition reads are performed once the web service starts listening.
- Patterns expressed in the request definitions are available in the object passed to the template

# YAML_CONFIG_FILE

port
: (integer, defaults to 3032) the port number the service should listen on

templates
: (list of template paths and options)

## a template description

An individual template description has four properties.

request
: (string, required) This is the METHOD and PATH that `newtmustache` will listen for to map this template. If no HTTP method is indicated a POST will be assumed to be the target. Go 1.22 HTTP handler patterns's variables maybe used to overwrite attributes in the submitted JSON Object. They are applied after those in `.options` have been applied.

template
: (string, required) This is the path to the primary Mustache template for this request. The source is read and associated with this request signature.

options
: (object, optional) These are additional attributes that can be merged into the JSON processed by the Mustache template. These are merge before any variables taken from the URL path pattern are merge.

debug
: (bool, optional) If set to true then provide verbose log output for this request route

# EXAMPLES

## YAML configuration

This example shows fix different template options. The first three applies a custom `page.tmpl` in different ways. In the last two the default Mustache template is setup.

~~~yaml
port: 3032
templates:
  - request: "POST /custom_page"
    template: page.tmpl
    options:
      title: This is the custom template with this title
  - request: "POST /custom_page_with_title/{title}"
    template: page.tmpl
    options:
      title: This title is overwritten by the one in the request
  - request: "POST /custom_page_include"
    template: page.tmpl
  - request: "POST /default_html5"
    options:
      title: A Page using the default template
  - request: "POST /default_html5_with_title/{title}"
    options:
      title: This title is replaced by the title in the URL
~~~


