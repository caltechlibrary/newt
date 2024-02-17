
# Newt Project

The Newt Project is an experiment in rapid web application development for libraries, archives and museums (abbr: <abbr title="libraries, archives and museums abbrevation">LAS</abbr>).  Newt uses a service oriented architecture formed data pipelines[^11]. The pipelines compose the web application. Newt comes with several to facilitate implementing this architecture.

[^11]: A data pipeline is formed by taking the results from one web service and using it as the input to another web service. It is the web equivalent of Unix pipes. Prior art: [Yahoo! Pipes](https://en.wikipedia.org/wiki/Yahoo!_Pipes)

Taking this approach minimizes the software you need to write in favor of configuration. This is because "off the shelf" software is available to do the heavy lifting. Example, Postgres+PostgREST provides data management and a JSON API out of the box. MySQL plus MRS does the same. Solr, Elasticsearch and Opensearch all provide full text search in the form of a JSON API. Again these are all out of the box. Pair these with some data routing and a light weight template engine and you can check off most of the core features used to implement LAS software systems.  Newt provides the missing parts -- data router, template engine even a code generator. A YAML file is used to tie all this capability together. 

The Newt is trying to encourage the following characteristics.

- preference for "off the shelf" over writing new code
- data modeling and management placed squarely in your database management system
- data sources accessed as JSON API
- data transformation (if needed) happens in stateless template engines
- code generation where appropriate

Newt provides tools to fill in the gaps.

- `newtrouter` is a stateless web service (a.k.a. micro service) that routes a web request through a data pipeline built from other web services
- `newtgenerator` is a code generator that can takes a set of data models described in YAML and generates SQL and Mustache templates
- `newtmustache` is a simple stateless template engine inspired by Pandoc server that supports the Mustache template language

The Newt web service and data pipeline concept is being tested using

- [Postgres](https://postgres.org), data management and modeling
- [PostgREST](https://postgrest.org), a service that turns Postgres into a JSON API

The Newt YAML ties this together expresses

- application (run time information and application metadata)
- models (describe the data as would be input into a web form)
- routes (web requests differentiated by a HTTP method and URL path)
- templates (Mustache template service)

## What type of applications are supported by Newt?

Most LAS applications are focused on managing and curating some sort of metadata records. This is the primary target of Newt. This might be as simple as a controlled vocabulary or as complex as an archival or repository metadata record.

## Motivation

Over the last several decades web applications became very complex. This complexity is expensive in terms of reliability, enhancement, bug fixes and software sustainability.

> A brief historic detour to set context

Databases have been used to generate web pages since the early web.  Databases are well suited to managing data.  When the web became dynamic, databases continued to be use for data persistence. By 1993 the web as an application platform was born[^12] and with it a good platform for providing useful organizational software.

By the mid 1990s the Open Source databases MySQL and Postgres were popular choices for building web applications. It is important to note neither MySQL or Postgres spoke HTTP[^13]. To solve this problem many people wrote software in languages like Perl, PHP and Python that ran inside the popular Apache web server. It was a pain to setup but once setup relatively easy to build things that relied on databases.  This led the web to explode with bespoke systems for curating and distributing web content. By the late 1990s and the early 2000s the practice of "mashing up" sites (i.e. content reuse) became the rage. As this increased in popularity web systems specialized further to cater to reuse. [Yahoo Pipes](https://en.wikipedia.org/wiki/Yahoo!_Pipes) was a very interesting expression of the "mashup culture"[^14]. Yahoo Pipes inspired Newt's data pipelines.  Specialization has continued ever since. Some of these systems have become less bespoke. Eventual bespoke systems gave way to common use cases[^15]. A good example of a common use case is Apache's [Solr](https://solr.apache.org) search engine. Another example was the bespoke content management systems gave way to systems like [Drupal](https://drupal.org) and [WordPress](https://wordpress.org).

[^12]: Web applications proceeded to eat all the venerable green screen systems they could find. Eventually they and their corporate sponsors invented the surveillance economy we have today. Sometimes "good ideas" have terrible consequences. Making it easier to produce custom web applications should always be done keeping in mind the necessity for humane and inclusive use. Newt can be both part of a solution but also be used to exacerbate human problems. Develop with consideration for others.

[^13]: HTTP being the protocol the communicates with. Essentially at the time RDBMS spoke a dialect of SQL as the unifying language. The web of the time understood HTML and to a certain degree XML. By 2000 people were looking for something simpler than XML to move structured data about. [JSON](https://en.wikipedia.org/wiki/JSON) quickly became the answer.

[^14]: The basic concept was to make it easy to work with "data feeds" and combined them into a useful human friendly web pages. It even included a visual programming language to make it friendly to the non-programmer crowd.

[^15]: If a use case is solved reliably enough it becomes "off the shelf" software.

> fast forward to 2024, context set

Much of the back end of web applications can largely be assemble from off the shelf software. Middleware however remains complex. I believe this to be a by product of inertia in software development practices and the assumption that what is good for "Google Scale" is good for everyone.

I think a radical simplification is due.  Most software doesn't need to scale that large. Even in the research and LAS communities we don't routinely write software that scales as large as [Zenodo](https://zenodo.org/).  We don't typically support tens of thousands of simultaneous users. If you accept that premise then we can focus our efforts around orchestrating off the shelf components and put our remaining development efforts into improving the human experience of using our software.

A big key to simplification is narrowing the focus of our middleware. It is not necessary for each middleware to implement data modeling, access control, user management, or data transformation. Like the Unix philosophy, simple middleware doing a specific task. Let your data base model your data and manage it. Let a full text search engine provide search, Let your front end web server with integrated single sign-on control access. Ideally assembling the back end would be turnkey. Newt attempts to make that possible for metadata curation applications.

## Off the shelf deliverables

Even without Newt we have allot to work with. Here's a shortlist of ones I use.

- (data management) Postgres combined with PostgREST gives you an out of the box JSON API for managing data
- (full text search) Solr gives you a powerful, friendly, JSON API for search and discovery
- (access control) Apache 2 or NginX when combined with a single sign on system (e.g. Shibboleth) provides access control
- (rich client) Web browsers now provide a rich software platform in their own right

## The missing bits

With the above list we can build capable applications relying on the sophisticated features of our web browsers. There is a cost. JavaScript is required to render everything. Relying on JavaScript in the browser to assemble our content from JSON or XML API is a horrible idea[^16]. A better approach is for the web browser to make a minimum number of requests to a web service and get back useful results without having to process more than HTML and CSS.

Taking the better approach in the past has required the writing of complex middleware. Newt's few tools help us avoid writing middleware.

For over a decade web frameworks developed for programming languages like Go, Java, JavaScript, PHP, Python, and Ruby have relied on a concept of "routes". A "route" is described by an HTTP method and URL path. The mapping of a route to a function simplifies the model of receiving and responding to HTTP requests[^17]. The collection of routes and their functions compose the API your browser uses to navigate through your application.

Historically before single sign-on systems became common the function handling a web request was responsible for the whole transaction. A single function would need to handle access control, data validation, data formatting, storing or retrieving data from the database. You had to make sure the request from the public didn't lead to a compromise of your database or operating system. Inertia has ensure many web applications built as middleware still do this. Even successful projects like Drupal and WordPress do this.

If we narrowly focus a function and allow other layers of the system to handle most of the complexity our functions can be simpler. Easier to write. Easier to reason about.  While that was the motivation for many software libraries and frameworks in practice they just wound up empowering more complexity. We could do more so we did. Maybe the solution isn't a programming language and framework but it is something else.

Are we missing an opportunity to change our base assumptions? Middleware doesn't need to do allot to be useful. That was a the promise when the term "micro service" was coined. The key is to make it easier to chain our middleware together. Unix accomplishes that with pipes. That web already has this concept of this but it is less obvious.  Can we make this easy then build towards the small and simple?

### A new baseline

Web services talk to other web services all the time. This isn't new. It isn't exotic it. Library, archive and museums systems do this all the time between organizations. Can we make it easier for doing this locally? In the little application your organization needs if you can build it. The approach is called "service oriented architecture". More recently was called "micro service architecture". Whatever you call it we're already using it at large (e.g. looking up a DOI via CrossRef). What does it take to do that in the small?

- Can we align access control with our front end web server?
- Can we insist our database management system provides flexible JSON API?
- Can we treat the output of one web services as the input for the next?
- Can we aggregate these into a data pipelines?
- Will that be enough to define a web application?

In 2024 you can answer is "yes we can" to all these propositions.

Newt provides a small set of web services and a code generator. These are driven from YAML files. The Newt code generator lets you bootstrap your application by generating SQL for Postgres+PostgREST. It can also generate templates for web pages. The Newt router takes that YAML file and pairs the requests with data pipelines defined in YAML. The last stage of the pipeline executed is returned to the web browser. Newt's code generator and router share the same YAML file keeping implementation and orchestration in sync. There isn't much middleware left to write using this approach. If you do write middleware you can make it simple and laser focused to a specific stage of your data pipeline. 

[^16]: See <https://infrequently.org/2024/01/performance-inequality-gap-2024/> for a nice discussion of the problem.

[^17]: In programmer's jargon we call that a "handler" or "route handler".

## What comes with the Newt Project?

- [newtrouter](newtrouter.1.md) a [web service](https://en.wikipedia.org/wiki/microservices) designed for working with other "off the shelf" web services. It functions both as a router and as a static file server. It does this by routing your request through a YAML defined pipeline and returning the results. Typically this will be a JSON data source and running that output through a template engine like Newt Mustache.
- [newtgenerator](newtgenerator.1.md) is a command line program that reads the Newt YAML file and generates SQL and templates used to build your application.  Currently the generator target SQL for use with Postgres+PostgREST. The template language being targeted is Mustache.
- [newtmustache](newmustache.1.md) implements a simple lightweight template engine supporting [Mustache](https://mustache.github.io/) templates. Mustache template language is well support by a wide variety of programming languages include Python, PHP, and JavaScript.

NOTE: See the [user manual](user_manual.md) for details

## How is Newt speeding up development?

The Newt suite of tools use a common YAML file.

- `newtgenerator` can render SQL suitable for bootstrapping your Postgres+PostgREST database as well as Mustache templates
- `newtmustache` provides a simple, stateless, Mustache template engine
- `newtrouter` provides data routing between other web services that fill specific functions or roles.

Your back end is constructed from "off the shelf" parts. Newt provides the routing. It allows our customization efforts to focus on data modeling in the database and with a template engine it renders content that you view in your web browser.

## Where is development time spent?

The developer writes YAML to generate the back end data management. The YAML is used generate the SQL and Mustache templates needed for your application. The developer is free to enhance the SQL, Mustache templates as needed. That shifts your focus to data modeling in the database or making your application more human friendly browser side.

## What about security, integration with single sign-on or other websites or services?

The `newtrouter` is a simple web service providing data routing based on its YAML configuration. It's a team player. In a production setting it should be used behind a front end web server like Apache 2 or NginX. That latter can be configured to support single sign-on systems like Shibboleth[^18]. The front end web service controls access and handles securing the connection with the web browser. The front end web service proxies to the Newt router. Newt router receives requests and runs the data pipelines on localhost. The data pipelines can be composed of off the shelf software like Postgres+PostgREST, Solr and template engine to turn your JSON into a web page. Having a clear division of responsibilities helps in securing your web application. Since Newt router only knows how to talk to services on localhost you can keep it contained and prevent it from being used for nefarious actions off system. Like Newt router Newt Mustache is constrained to localhost for similar reasons.

Limiting Newt web service applications to localhost keeps them simple. Doing the minimum limits the attack surface for those who want to do mischief.  Neither `newtrouter` or `newtmustache` write to disk or require secrets. They only communicate via localhost using HTTP protocol.

[^18]: Shibboleth is a common single sign-on platform in research libraries, universities and colleges.

### What about "scaling"?

`newtrouter` is just a router. Aside from reading configuration at start up it doesn't maintain state. `newtmustache` functions the same way, read in the configuration and just run. By assigning different ports you can also run many instances of them. This makes it possible to run them in parallel, behind load balancer or even through proxying spread them across many machines. The instances don't share data or coordinate. They start up wait for a request and providing an answer. 

So what does this all mean? In principle a Newt based applications can scale big as it pipeline services allow. 

### Anatomy of a Newt based web application

Newt application development is friendly to version control systems (e.g. Git). It consists of a Newt configuration file, along with generated SQL files, HTML templates and any static web assets you've added. A typical disk layout of a Newt project could look like this-

- `/` project folder
  - `htdocs` this directory holds your static content needed by your web application
  - `*.sql` these are the SQL files used by your application to define your models and behaviors in Postgres
  - `templates` a template holding your template pages
  - `app.yaml` would hold the a Newt router and code generator configuration file
  - `tmpl.yaml` holding the configuration for `newtmustache`, it cannot run on the same port as `newtrouter`.
  - `CITATION.cff` or `codemeta.json` for project metadata

> Newt, a type of salamander, doesn't seek attention. It does its own thing. You only notice it if you look carefully.

## About the Newt source repository

Newt is a project of Caltech Library's Digital Library Development group. It is hosted on GitHub at <https://github.com/caltechlibrary/newt>. If you have questions, problems or concerns regarding Newt you can use GitHub issue tracker to communicate with the development team. It is located at <https://github.com/caltechlibrary/newt/issues>.

## "Someday, maybe" exploration

- Integrate S3 protocol to support storing binary or large objects

## Getting help

**The Newt Project is an experiment!!**. The source code for the project is supplied "as is". Newt is most likely broken. At a stretch it could be considered a working prototype. You should not use it for production systems.  However if you'd like to ask a question or have something you'd like to contribute please feel free to file a GitHub issue, see <https://github.com/caltechlibrary/newt/issues>. Just keep in mind it remains an **experiment** as of February 2024.



