---
title: nte(1) user manual | 0.0.9 ac665bc
pubDate: 2024-07-30
author: R. S. Doiel
---

# NAME

nte

# SYNOPSIS

nte [OPTIONS] YAML_CONFIG_FILE

# DESCRIPTION

**nte** is a web service that provides a Mustache template rendering engine inspired
by Pandoc server.

Unlike Pandoc web server, `nte` expects a YAML_CONFIG_FILE. The format is
described below. That file specifies the runtime configuration. It specifies the request path
to template mapping. It can also specify ancillary information made available to the Mustache
template associated with the request path and template.

The `nte` template engine listens for a POST requests of JSON encoded data.
It  checks requested path to see if that matches the request path described in the Newt YAML
file. If there is a match it processes the request returning the template results matched with
 any data found in the POST. `nte` doesn't respond to any other HTTP methods.

The content of the POST is passed to the template as `.body`, applications options
are passed to the template as `.options`, any vocabulary content read in at startup
is passed to the template as `.vocabulary`. Finally if you've defined a variable
in your request path those will be available to your template as `.vars`.

Like Pandoc web service `nte` does not normally log requests. It's a quick
transaction. If you want to debug your templates use the verbose command line option to turn on
debug output.

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

`-timeout SECONDS`
: Timeout in seconds, after which a template rendering is aborted.  Default: 3.

`-verbose`
: If set provide verbose debugging output for requests

# The templates

Mustache templates are documented at <https://mustache.github.io>. The template engine
used is based on Go package <https://github.com/cbroglie/mustache>.

## Features

- Newt template engine only runs on localhost at the designated port (default is 8011).
- Templates are read in at startup and are retained in memory bound to the request path.
- Vocabulary files are read in at startup and bound to the request path.
- Options are set at startup and mapped into the request path.
- No addition reads are performed once the web service starts listening.
- Variables found expressed in the request path are available in the `.vars`
passed to the template.

# YAML_CONFIG_FILE

This is a list of the Newt YAML syntax relevant to **nte**.

## Top level properties

These are the top level properties in YAML files.

applications
: (optional) holds the run time configuration used by the Newt applications.

templates
: (required by newtmustache)

## The applications property

The applications properties are optional. Some maybe set via command line. See Newt application's manual page for specific ones. These properties lets you override the default settings of Newt programs.

template_engine
: this contains configuration for Newt template engine, e.g. port, base_dir, ext_name.

options
: holds key value pairs of which can be referenced in the values of models, routes and templates.

### newtmustache settings

port
: (all) Port number to used for Newt web service running on localhost

### the "routes" property

Routes hosts a list of request descriptions and their data pipelines. This property is only used by Newt router and Newt code generator.

## templates property

This property is used by Newt template engine. It is ignore by Newt router and code generator.

templates
: (optional: newtmustache) this holds a list of template objects

### template object model

The template objects are used by Newt template engine. If you're not using it you can skip these.

`request PATH`
: (required) This holds the request URL's path. `nte` only listens for POST method.

`template`
: (required: newtmustache only) This is the path to the template associated with request. NOTE: Pandoc web service does not support partial templates. Mustache does support partial templates

`partials`
: (optional, newtmustache only) A list of paths to partial Mustache templates used by `.template`.

`options`
: (optional, newtmustache only) An object is passed to the template as `.options`.

`vocabulary`
: (optional, newtmustache only) This is the filename for a YAML file which is exposed inside the template as `.vocabulary`. You can think of this as options maintained outside the Newt YAML file.

`debug`
: (optional) this turns on debugging output for this template

# EXAMPLES

Example of newtmustache YAML that only runs the template engine by itself.
The paths are used to provide template content.

~~~yaml
applications:
  template_engine:
    port: 8011
	# where to find the templates
	base_dir: testadata
	# where, under base_dir, to find partial templates
	partials: partials
	# template extension to use
	ext_name: .tmpl
templates:
  - id: hello
    request: /hello/{name}
    template: simple
  - id: hello
    request: /hello
    template: simple
    options:
      name: Universe
  - id: hi
    request: /hi/{name}
    template: hithere
    debug: true
  - id: hi
    request: /hi
    template: hithere
    options:
      name: Universe
~~~


