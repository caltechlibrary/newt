
# `newtmustache` exampled

The Newt Mustache template engine came about due to three concerns.  

- People don't think about Pandoc as a web service
- Some people don't like Pandoc templates or are unfamiliar with Pandoc
- Pandoc's web server doesn't handle Pandoc's partial template support

When building a web application using a template system it is very useful to minimize the number of templates you need to know and work with. One way this is done is to build smaller partial templates that handle specific content elements or types. E.g. a bibliographic citation used on a book receive site

When you run Pandoc from the command line this is readily supported. Unfortunately Pandoc web service doesn't support this as of version 3.1. The web service only knows about the full template you provide as part of the JSON blob.  Another problem with the web service is that you must provide both the template and JSON DATA for each request you make to the web service. Internally Pandoc web service unpacks the template, parses it then applies the data object you provided. How might we simplify the process?

Ideally a simple template engine is stateless once it has started up. This doesn't mean it can't take advantage of configuration when starting up. If a simple template engine knew that templates were found in a directory it could read all the templates in the directory at start up. It could resolve partial template issues and be ready for any request. If we leverage the URL to indicate the template desired to be used then the JSON received as the POST data doesn't need to include a wrapping object like Pandoc web server does. 

I thought about writing a pre-processor for Pandoc web server but then decided I don't have time. I checked the current status of languages supporting Mustache templates. The list had grown and also includes a nice Go package for working with Mustache templates. That was the seed for Newt's own template engine.

- Mustache is a widely support template language with support include Go, Python, and JavaScript (languages used regularly at Caltech Library)
- Since a Go package provides Mustache template I only need to write a light weight web service to wrap it
- Since I am writing the service I can keep the requires to a minimum, i.e. port and template directory

## `newtmustache` startup sequence

- read in the templates (e.g. read the contents of `templates` directory or the location specified by the command line option`
- parse the read templates, log errors and skip templates which do not parse succesfully
- create a URL path handler for each template found in the directory
- start up the web service listen on the default port or one supplied on the command line


