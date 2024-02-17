
## Newt design choices

Newt needs to help make it trivial to quickly generate metadata curation applications. It needs to leverage existing software and knowledge.

I know Newt needs to play well with the following software.

- [Postgres](https://postgres.org) + [PostgREST](https://postgrest.org)
- [Solr](https://solr.apache.org)

JSON data is the way I want to work with structured data coming out of these system. I also want to minize any additional middleware needed to assemble a base metadata curation app featuring the basic CRUD-L[^201] operatins. I want assembling my web application to be conceptually as easy as working with Unix pipels in the shell. An example would be to retrieve data from a JSON API (e.g. the software above) and then run it through a simple minimal template engine. I can see cases, based on my prior experiece with Pandoc server, there way be more than two stage in the processing sequence. I need a data workflow, a pipeline.  We do this all the time in web applications but often it is browser side, e.g. CrossRef and ORCID data pulled in via JavaScript. I don't want to rely on the web browser so I need a pipeline behind the web server. I need a data router.

Since Pandoc server isn't my ticket I started to look at other simple template languages. Mustache is one I remember from both browser side implementations and server side when I was working with PHP application at USC. Mustache is alive and kicking and is supported by most of the programming languages I'm likely to run into in libraries, archives and museum organizations. There a good, maintained, Mustache package for Go which is my implementation language for Newt. I can take the lessons learned using Pandoc server and apply them to a template engine as part of the Newt project. It will be separate from the data router so it can be easily swapped out if someone doesn't like Mustache.  It needs to take the JSON object or array from a data source and transform it. So the configuration for the Mustache template engine needs to be part of the YAML used in a Newt project.

This all implies three tools for Newt's second prototype

- data router supporting the pipeline concept
- template engine that plays well with the data router
- a code generator that understands Newt's YAML syntax

This is buildable.

It lets us leverage systems like Postgres and PostgREST. Newt can lower the SQL burden through code generation while keeping the data modeling and management in the database (a core Newt philosophy). Using a single YAML for data modeling means one sorce for both web form generation and for SQL generation. Here's what I need at the top level of my YAML

- application metatadata and runtime settings (e.g. port numbers, directories, OS environment to leverage)
- models to describe the web forms and SQL needed from the database
- routes describing what the data router will listen for
  - pipelines identifying the services and contact requirements to process the data
- templates a mapping of template requests to template processing the received JSON data

Developing the second prototype will help fill in the details guided by these observations.

[^201]: CRUD-L is an ancronym for create, read, update, delete and list. These are the basic operations available in a database management system.

