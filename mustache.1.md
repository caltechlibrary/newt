---
title: mustache(1) user manual | 0.0.8 e146d30
pubDate: 2024-05-21
author: R. S. Doiel
---

# NAME

mustache

# SYNOPSIS

mustache [OPTIONS] TEMPLATE_NAME [JSON_FNAME]

# DESCRIPTION

**mustache** is command line template rendering engine. It is provided so
you can debug your mustache templates without resorting to using cURL. It uses
the same template library that Newt's Mustache uses so the results should
match. **mustache}** can read JSON data from a file or from standard 
input. The template name needs to be provided on the command line.
In this way **mustache** can be used to process JSON in a POSIX style
pipe line.

# OPTIONS

The following options are supported by **mustache**.

-h
: display this help message

-license
: display the software license

-version
: display version information

-page
: points at a YAML file that will be used as elements in the Mustache template.

# EXAMPLE

In this example there is a JSON file called "data.json" and a template called "page.tmpl"
and **mustache** is used to run the JSON data through the template. In the first
example the data file is specified as part of the command line arguments in the
second it is read from standard input via file redirection. The Third version
works the same way but data is from a pipe.

~~~shell
mustache page.tmpl data.json
mustache page.tmpl <data.json
cat data.json | mustache page.tmpl
~~~


