---
title: "ws(1) user manual | 0.0.7-dev 24f2d1a"
pubDate: 2024-02-22
author: "R. S. Doiel"
---

# NAME

ws

# SYNOPSIS

ws [OPTIONS] [HTDOCS]

# DESCRIPTION

**ws** is a simple static web server running on localhost. The default port
is 8000 and the current working directory is the default document root. There is no
configuration aside from providing a directory path and optionally changing the port.

This is a mimimal web server. No http, no embedded programming languages. No remapping
content types, redirect or other fancy stuff. It provides a quick way to allow your
static content to be available to your web browser over http for development purposes.

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

# EXAMPLE

In the example below the web server would listen for `http://localhost:8080`
and respond with the content in htdocs.

~~~shell
ws -port 8080 htdocs
~~~


