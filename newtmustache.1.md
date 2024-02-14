---
title: "newtmustache(1) user manual | 0.0.6 3b0b6fd"
pubDate: 2024-02-13
author: "R. S. Doiel"
---

# NAME

newtmustache

# SYNOPSIS

newtmustache [OPTIONS]

# DESCRIPTION

**newtmustache** is a web service that provides a Mustache template rendering inspired by Pandoc server.

Like Pandoc web server there is no configuration file. There are a few command line options, e.g.
port, template directory and timeout.

# Template engine API

**newtmustache** accepts POST of JSON content and maps them to a template name expressed in the 
request URL. If I had a template called `template/list_objects.tmpl`
then the url would be formed like `http://localhost:3032/list_objects.tmpl`. The JSON
encoded post would then be sent through the "list_objects.tmpl" template and returned to the browser.
The JSON object POSTed does not need a wrapping object like required by Pandoc server.  **newtmustache**
reads in the templates found in the template directory at start up. It creates a map between the 
basename of the template and the URL handler built from the that name. As such the templates are fixed
at startup and do not need to be passed along with your data object.

This improves of Pandoc web servers in a few ways. 

- Now wrapping object
- Template parse errors are known earlier own and are visible from the output log
- Parsing of the templates happens once at startup
- Partial templates can be supported because the startup phase of **newtmustache** handles resolving partials


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

-templates
: (default directory name is "templates") Pick an alternative location for finding templates

--timeout SECONDS
: Timeout in seconds, after which a template rendering is aborted.  Default: 3.

# The templates

Mustache templates are documented at <https://mustache.github.io>. The template engined
used is based on Go package <https://github.com/cbroglie/mustache>.


