
# Newt Mustache explained

## Overview

Newt Mustache is a simple, light weight template engine supporting the use of Mustache templates. If you send a JSON object to a configured Newt Mustache template engine it can run that object through Mustache and hand back a result. This usually means taking a JSON object (e.g. JSON object delivered via PostgREST) and turning that object into web page. That's the type of templates you get from using Newt Generator. Newt Mustache itself just cares about the JSON object it recieves and applying the template configured for the request. Normally the JSON object is sent to Newt Mustache as a HTTP POST action.

## Testing your templates with Newt Mustache

Assuming you have Newt Mustache configuration in your "app.yaml" file (i.e. there is a templates property completed there) and it points to valid Mustache templates then Newt Mustache is ready to test.  I usually use [curl](https://curl.se) to do my initial testing. In the example I'm assuming you have a JSON file named "person.json" that contains JSON markup that fits the templated called "people_read.tmpl" which is going to be accessible with the URL `http://localhost:8040/people_read.tmpl`. I first start Newt Mustahe as a back ground process then use curl to check the template result.

```shell
newtmustache app.yaml &
curl --data '@person.json' http://localhost:8040/people_read.tmpl
```

I should now see HTML markup with the "person.json" content rendered into the template. That's all that Newt Mustache does. It reads in a bunch of templates and requests paths at start, then renders them based on the JSON objects you send them.

It is important to remember that changes you make in your templates will NOT show up in the service until you restart Newt Mustache. This is similar to the constraints regarding Newt Router and PostgREST.

### Why Newt Mustache?

The Newt Mustache template engine came about due to three concerns I encountered with Pandoc web server.  

- Pandoc server usage is hard to debug
- Many people don't know or like Pandoc's template language
- Pandoc server does NOT support partial templates

When building a web application using a template system it is very useful to minimize the number of templates you need to know and work with. One way this is done is to build smaller partial templates that handle specific content elements or types. E.g. a bibliographic citation used on a book review website.

When you run Pandoc from the command line this is readily supported. Unfortunately Pandoc 3.1 server doesn't support this. The server only knows about the full template you provide as part of the JSON POST.  Mistakenly misconfigure the JSON you post to Pandoc and it will happly give you nothing back. Not an error message, not a peep. Technically it's giving back what you requested but it's a pain to fiddle with JSON to get enough of a response to diagnose the problem. That was a show stopper for Newt's second prototype. Time to switch horses. Pandoc server inspired Newt's Mustache template service. Newt's template service is configured from YAML file. Settings are clearer. You can also turn on debugging for a specific template you have concerns about.  Like Pandoc server Newt' Mustache template engine is stateless. You can run as many as you like as long as you have an available port.

Newt Mustache designd for use in Newt's data pipeline.

### Why Mustache templates?

- Mustache is a widely support template language with support include Python, JavaScript, and Go (languages used regularly at Caltech Library)
- Since a Go package provides Mustache template I only need to write a light weight web service to wrap it
- Since I am writing the service I can keep the requires to a minimum, i.e. [Use Newt's YAML file syntax](newtmustache.1.md#newt_config_file).

See [newtmustache](newtmustache.1.md) manual page for details.

