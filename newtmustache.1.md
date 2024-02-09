---
title: "newtmustache(1) user manual | 0.0.6 c529101"
pubDate: 2024-02-08
author: "R. S. Doiel"
---

# NAME

newtmustache

# SYNOPSIS

newtmustache

# DESCRIPTION

**newtmustache** is a microservice that provides a Mustache template rendering inspired by Pandoc server.

There is no configuration file. There are command line options to specify port and timeout.

# API

## Root endpoint

The root (/) endpoint accepts only POST method requests. Other methods will yield a appropriate http status code.

##  Request

The body of the POST request should be a JSON object, with the following fields.  Template is always required. It must
also include at least one of either metadata or variables. If both are set then the maps are merged before processing
with the template.

template (string)
: String contents of a document template, see Mustache <https://mustache.github.io/mustache.5.html>.

data (JSON map)
: String-valued metadata.

content_type (optional, default is text/plain)
: Set the rendered content type, default is text/plain. You probably want text/html for web pages.

## Response

It returns a response body and headers set by the default http writer provided in the http go package.
If there request can't be fullfilled then an http status code and text will be returned.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

--port NUM
: HTTP port on which to run the server.  Default: 3030.

--timeout SECONDS
: Timeout in seconds, after which a conversion is killed.  Default: 2.

# The templates

Mustache templates are documented at <https://mustache.github.io>. The template engined
used is based on Go package <https://github.com/cbroglie/mustache>.


