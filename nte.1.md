---
title: nte(1) user manual | 0.0.9 3662ba0
pubDate: 2024-08-14
author: R. S. Doiel
---

# NAME

nte

# SYNOPSIS

nte [OPTIONS] YAML_CONFIG_FILE

# DESCRIPTION

**nte** is a web service that provides a template rendering engine inspired
by Pandoc server.

Unlike Pandoc web server, `nte` expects a YAML configuration file.
The format is described below. That file specifies the runtime configuration. It specifies
the request path to template mapping. It can also specify ancillary information made
available to the template associated with the request path and template.

The `nte` template engine listens for a POST requests with JSON encoded data.
It  checks requested path to see if that matches the request path described in the YAML
file. If there is a match it processes the request returning the rendered results matched with
 any data found in the POST. `nte`.

The content of the POST is passed to the template as `.body`, applications options
are merged into a `.document` object along with any addition mappings specified for
a given template.  Finally if you've defined a variables in the path to the template those
are provided via the `.vars` property.

**nte** only supports POST requests with content type "application/json".

# OPTIONS

The following options are supported by **nte**.

`-h`
: display this help message

`-license`
: display the software license

`-version`
: display version information

`-port NUMBER`
: (default is port is 3032) Set the port number to listen on

`-base-dir PATH`
: set the base directory path (where you have your templates).

`-timeout SECONDS`
: Timeout in seconds, after which a template rendering is aborted.  Default: 3.

`-verbose`
: If set provide verbose debugging output for requests

# The templates

The template engine supports the [Handlebars](https://handlebarsjs.com) template language
which is largely a superset of Mustache templates documented at <https://mustache.github.io>.
The template engine used is based on Go package <github.com/aymerick/raymond>.

## Features

- Newt template engine only runs on localhost at a designated port (default is 8011).
- Templates are read in at startup and are retained in memory bound to the request path.
- JSON data is provided to the template in a `.body` object.
- Variables found expressed in the request path are available in the `.vars`
  passed to the template.
- Except for path variables no addition reads are performed once the web service starts listening.

# YAML_CONFIG_FILE

This is a list of the Newt YAML syntax relevant to **nte**.

## Top level properties

These are the top level properties in YAML files.

applications
: (required) holds the run time configuration used by the Newt applications.

templates
: (required) holds a list of template objects

## The __applications__ property

template_engine
: (required) this contains configuration for Newt template engine, e.g. port, base_dir, ext_name.

### The __template engine__ properties

port
: (required) port number to used for to used for Newt Template Engine

base_dir
: (required) base directory holding the primary templates

partial_dir
: (optional) the sub directory holding the partial templates

ext_name
: (optional) the extension used to identify your templates on
disk. e.g. ".hbs" for handlebar templates.

## The __templates__ property

This property is used by Newt template engine. It provides a list of
template objects.

### The __template__ property

The template object is used by Newt template engine to describe
an individual template mapping and it's properties.

`id`
: (required) Unique template identifier. It is used by other parts of Newt.

`description`
: (suggested) A description of template's purpose. Used by other parts of Newt.

`request PATH`
: (required) This holds the request URL's path. `nte`
only listens for POST method. It may include path variables. The request
path must be unique.


`template`
: (required) This is the name of the primary template (without file extension).
The primary template may also include partials and those will be read from
the partials sub directory defined in the template engine property.

`document`
: this will provide template specific data include content verged from
the provided environment (e.g. template engine's options and environment).

`debug`
: (optional) this turns on debugging output for this template

# EXAMPLES

Example of Newt YAML that only runs the template engine by itself.
The paths are used to provide template content.

~~~yaml
applications:
  template_engine:
    port: 8011
	base_dir: testdata/views
	ext_name: .hbs
	partials: partials
templates:
  - id: hello
    request: /hello/{name}
    template: simple
  - id: hello
    request: /hello
    template: simple
    document:
      name: Universe
  - id: hi
    request: /hi/{name}
    template: hithere
    debug: true
  - id: hi
    request: /hi
    template: hithere
    document:
      name: Universe
~~~

NOTE: the template name doesn't require the extension since that is set at the 
template engine level.


