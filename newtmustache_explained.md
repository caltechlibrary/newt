---
title: newt mustache templated support explained
pubDate: 2024-02-12
updated: 2024-02-14
author: R. S. Doiel
---

# Newt Mustache template supported explained

2024-02-12

The Newt Mustache template engine came about due to three concerns.  

- People don't think about Pandoc as a web service
- Some people don't like Pandoc templates or are unfamiliar with Pandoc
- Pandoc's web server doesn't handle Pandoc's partial template support

When building a web application using a template system it is very useful to minimize the number of templates you need to know and work with. One way this is done is to build smaller partial templates that handle specific content elements or types. E.g. a bibliographic citation used on a book review website.

When you run Pandoc from the command line this is readily supported. Unfortunately Pandoc 3.1 web service doesn't support this. The web service only knows about the full template you provide as part of the JSON POST.  Another problem with the web service is that you must provide both the template and JSON DATA for each request you make to the web service. Where does that template source come from? Typically it is read from disk unless you're using the default templates Pandoc provides. Internally Pandoc web service unpacks the template, parses it then applies the data object you provided. How might we simplify the process?

Ideally a simple template engine is stateless once it has started up. This doesn't mean it can't take advantage of configuration when starting up. If a simple template engine knew where the template was on disk it could read it and keep the template in memory. It could potentially resolve partial templates too. Combine this with mapping the template to a specific URL request we can eliminate the need for a wrapping object like the one required to use Pandoc web service. 

Newt Project provides a web service, [pdbundler](pdbunder.1.md) let's you prep JSON object for use with Pandoc web service. It solves some of the challenges of working with Pandoc web service.

[Newt Mustache](newtmustache.1.md) takes that a step further and integrates a Mustache template engine which will include partial template support resolved at startup.

Why Mustache templates?

- Mustache is a widely support template language with support include Python, JavaScript, and Go (languages used regularly at Caltech Library)
- Since a Go package provides Mustache template I only need to write a light weight web service to wrap it
- Since I am writing the service I can keep the requires to a minimum, i.e. [A simple YAML configuration file](newtmustache.1.md#newt_config_file).

See [newtmustache](newtmustache.1.md) manual page for details.

