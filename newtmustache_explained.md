---
title: newt mustache templated support explained
pubDate: 2024-02-12
updated: 2024-02-14
author: R. S. Doiel
---

# Newt Mustache template supported explained

2024-02-12

The Newt Mustache template engine came about due to three concerns.  

- Pandoc server usage is hard to debug
- Many people don't know or like Pandoc's template language
- Pandoc server does NOT support partial templates

When building a web application using a template system it is very useful to minimize the number of templates you need to know and work with. One way this is done is to build smaller partial templates that handle specific content elements or types. E.g. a bibliographic citation used on a book review website.

When you run Pandoc from the command line this is readily supported. Unfortunately Pandoc 3.1 server doesn't support this. The server only knows about the full template you provide as part of the JSON POST.  Mistakenly misconfigure the JSON you post to Pandoc and it will happly give you nothing back. Not an error message, not a peep. Technically it's giving back what you requested but it's a pain to fiddle with JSON to get enough of a response to diagnose the problem. That was a show stopper for Newt's second prototype. Time to switch horses. Pandoc server inspired Newt's Mustache template service. Newt's template service is configured from YAML file. Settings are clearer. You can also turn on debugging for a specific template you have concerns about.  Like Pandoc server Newt' Mustache template engine is stateless. You can run as many as you like as long as you have an available port.

[Newt Mustache](newtmustache.1.md) designd for use in Newt's data pipeline.

Why Mustache templates?

- Mustache is a widely support template language with support include Python, JavaScript, and Go (languages used regularly at Caltech Library)
- Since a Go package provides Mustache template I only need to write a light weight web service to wrap it
- Since I am writing the service I can keep the requires to a minimum, i.e. [Use Newt's YAML file syntax](newtmustache.1.md#newt_config_file).

See [newtmustache](newtmustache.1.md) manual page for details.

