---
title: "ws(1) user manual | 0.0.8 8f00277"
pubDate: 2024-04-10
author: "R. S. Doiel"
---

# NAME

ws

# SYNOPSIS

ws [OPTIONS] [HTDOCS]

# DESCRIPTION

**ws** is a simple static web server running on localhost. The default port
is 8000 and the current working directory is the default document root. There is no
configuration aside from providing a directory path, optionally changing the port and
setting a verbose mode useful for debugging requests.

This is a minimal web server. No HTTP, no embedded programming languages. No remapping
content types, redirect or other fancy stuff. It provides a quick way to allow your
static content to be available to your web browser over HTTP for development purposes.

Dot files are not served out.

# OPTIONS

The following options are supported by **ws**.

-h, -help
: display this help message

-license
: display the software license

-version
: display version information

-port
: set the port the web server listens on

-verbose
: show verbose logging of requests, e.g. contents of a POST

# EXAMPLE

In the example below the web server would listen for `http://localhost:8080`
and respond with the content in htdocs.

~~~shell
ws -port 8080 htdocs
~~~

An example of using the static file server to debug a form submission by showing
what the forms submits in the log output.

~~~shell
ws -verbose -port 8080 htdocs
~~~



