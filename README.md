
# Newt Project

Newt is an experimental set of tools for rapid application development of metadata curation tools. These types of applications are commonly needed in libraries, archives and museums (abbr: LAS). Newt applications implement a service orient architecture. This allows you to take advantage of off the shelf software like Postgres+PostgREST, MariaDB or MySQL with MySQL Rest Service, Solr, Elasticsearch and Opensearch. These software come with a JSON API, a JSON data source. Newt provides a data router that can pair a request with a JSON data pipeline[^11]. Place Newt Mustache as the last service in your pipeline and you can transform your JSON objects into web pages using the Mustache template language. Newt's code generator can take advantage of Newt's YAML file to render data models as the SQL and configuration for Postgres+PostgREST. Newt's code generator can use that same YAML to render Mustache templates too. With Newt's tools and a YAML file you can create a web application for managing metadata which supports that basic CRUD-L[^12] operations needed to curate your metadata. Through in a full text search engine like Solr and you have the core capabilities found in many LAS applications like ArchivesSpace. 

Newt is a narrowly focused rapid application development toolbox.  It does not support the class of web applications that handle file uploads. It's not a replacement for Drupal or WordPress. Newt generators the applications more in line with ArchivesSpace but with simpler data models. If you need file upload support you can either build that as a separate service or use a rich platform like Drupal or Invenio RDM.

Newt applications are well suited to interacting with any existing platform that provides a JSON API. LAS applications like ArchivesSpace and Invenio RDM provide JSON API. It is possible to extended those systems by creating simpler services that can talk to those JSON data sources. Newt is well suited to this "development at the edges" approach. The code you would need to write would be to configure a localhost proxy to those external systems. Similarly that opens the door to interesting integration with services like CrossRef, DataCite, Pub Med Central, ROR and ORCID. If what you need is a simple integration between these systems Newt might provide some useful tools to make that process easier.

[^11]: A data pipeline is formed by taking the results from one web service and using it as the input to another web service. It is the web equivalent of Unix pipes. Prior art: [Yahoo! Pipes](https://en.wikipedia.org/wiki/Yahoo!_Pipes)

[^12]: CRUD-L, acronym meaning, "Create, Read, Update, Delete and List". These are the basic actions used to manage metadata or objects.

A Newt application encourages the following.

- preference for "off the shelf" over writing new code
- use a database management system for modeling and managing your data
- prefer software that can function as a JSON data source
- data transformation (if needed) is provided as a light weight web service
- code generation where appropriate

In 2024 the first three can be provided by off the self software mentioned previously. Newt provides a few tools to fill in the gaps.

- `newtrouter` is a stateless web service (a.k.a. micro service) that routes a web request through a data pipeline built from other web services
- `newtgenerator` is a code generator understands the data models described in Newt's YAML configuration file. It can generate code to setup Postgres and PostgREST, it can generate Mustache templates
- `newtmustache` is a simple stateless template engine inspired by Pandoc server that supports the Mustache template language

Newt's web services and the data pipeline concept is being tested using

- [Postgres](https://postgres.org), data management and modeling
- [PostgREST](https://postgrest.org), a service that turns Postgres into a JSON API

The Newt YAML ties this together expressing

- application (run time information for Newt Router and Newt Mustache)
- models (descriptions of data as you would provided in a web form)
- routes (web requests differentiated by a HTTP method and URL path that trigger processing in a data pipeline)
- templates (pairs a request with a template to transform a JSON into some other form)


## What type of applications are supported by Newt?

Most LAS applications are focused on managing and curating some sort of metadata records. This is the primary target of Newt. This might be as simple as a controlled vocabulary or as complex as an archival or repository metadata record.

## Motivation

Over the last several decades web applications became very complex. This complexity is expensive in terms of reliability, enhancement, bug fixes and software sustainability.

> A brief historic detour to set context

Databases have been used to generate web pages since the early web.  Databases are well suited to managing data.  When the web became dynamic, databases continued to be use for data persistence. By 1993 the web as an application platform was born[^13] and with it a good platform for providing useful organizational and institutional software.

By the mid 1990s the Open Source databases like MySQL and Postgres were popular choices for building web applications. It is important to note neither MySQL or Postgres spoke HTTP[^14]. To solve this problem many people wrote software in languages like Perl, PHP and Python that ran inside the popular Apache web server. It was a pain to setup but once setup relatively easy to build things that relied on databases.  This led the web to explode with bespoke systems. The enabled in the late 1990s and the early 2000s the practice of "mashing up" sites (i.e. content reuse). As "mashing up" became the rage, bespoke systems took advantage of content reuse too. [Yahoo Pipes](https://en.wikipedia.org/wiki/Yahoo!_Pipes) was a very interesting expression of the "mashup culture"[^15]. Yahoo Pipes inspired Newt's data pipelines.  Eventual the bespoke systems gave way to common use cases[^16]. A good example of a common use case is Apache's [Solr](https://solr.apache.org) search engine. Another example is how bespoke content management systems gave way to [Drupal](https://drupal.org) and [WordPress](https://wordpress.org).

[^13]: Web applications proceeded to eat all the venerable green screen systems they could find. Eventually they and their corporate sponsors invented the surveillance economy we have today. Sometimes "good ideas" have terrible consequences. Making it easier to produce custom web applications should always be done keeping in mind the necessity for humane and inclusive use. Newt can be both part of a solution but also be used to exacerbate human problems. Develop with consideration for others.

[^14]: HTTP being the protocol the communicates with. Essentially at the time RDBMS spoke a dialect of SQL as the unifying language. The web of the time understood HTML and to a certain degree XML. By 2000 people were looking for something simpler than XML to move structured data about. [JSON](https://en.wikipedia.org/wiki/JSON) quickly became the answer.

[^15]: The basic concept was to make it easy to work with "data feeds" and combined them into a useful human friendly web pages. It even included a visual programming language to make it friendly to the non-programmer crowd.

[^16]: If a use case is solved reliably enough it becomes "off the shelf" software.

> fast forward to 2024, context set

Much of the back end of web applications can largely be assemble from off the shelf software. Middleware however remains complex. I believe this to be a by product of inertia in software development practices and the assumption that what is good for "Google Scale" is good for everyone.

I think a radical simplification is due.  Newt in part is intended to spark that conversation. My observation is most software doesn't need to scale large. Even in the research and LAS communities we don't routinely write software that scales as large as [Zenodo](https://zenodo.org/).  We don't typically support tens of thousands of simultaneous users. If you accept that premise then we can focus our efforts around orchestrating off the shelf components and put our remaining development efforts into improving the human experience of using our software. A better human experience is an intended side effect of Newt.

Back in the day whenever we were writing in PHP, Python or Perl we creating middleware. Even Drupal ad WordPress can be thought of as middleware. Middleware sites between a data source (e.g. a database) and the web server (e.g. Apache 2, NGINX).

A big key to simplification is narrowing the focus of our middleware. When our middleware has to implement everything it becomes very complex. Look at Drupal and WordPress. They implement data modeling, data management, user accounts, access management, data transformation.  What if middleware was narrowly focus? Conceptually simpler? Did one or two things really well?. The approach starts to sound an like the old Unix philosophy and a single tool that does one thing really web chained together to form a data processing pipeline. If you step back and look at the web today that is what happening. A catalog systems can import data from trusted sources. Creating a repository record might start by pulling in data from CrossRef or ORCID.  What is that type of integration was as easy as using Unix pipes for forming a data processing workflow?

If you can define the functions of a middleware system from existing system and route one system's output into another as import then we can use to pipeline approach to crate useful applications. We just need a means of describing the pipeline (YAML in the case of Newt) and running the pipeline when requested (Newt Router provides this).

Taking this approach (implementing a service oriented architecture) it becomes easy to assign responsibility to various off the shelf components. Here's an example dividing responsibility with off the shelf software.

1. Apache 2 with Shibboleth provides access control and helps up communicate with a web browser (e.g. SSL support in the public URL)
2. Apache 2 proxies to a Newt Router. Newt Router either maps the request to a data pipeline or provides static content
3. The data pipeline performs our processing. It can include any number of web services running on localhost. Which services depend on what we want to do. If we're managing data then Postgres+PostgREST is a good choice. The pipeline first stage might start there.  Eventually we'll want to turn that JSON into HTML so Newt Mustache provides a reliable means of doing that.

In this example structure we can mix and match database with a JSON API and search engines. Using this approach we can "mashup" our data sources locally from behind our front end web service. The web browser contacts the front in web service, it contacts Newt Router and the results flow back through Newt router and the front end web server to the web browser making the request.

## Off the shelf deliverables

Even without Newt we have allot to work with. Here's a shortlist of ones I use.

- (data management) Postgres combined with PostgREST gives you an out of the box JSON API for managing data
- (full text search) Solr gives you a powerful, friendly, JSON API for search and discovery
- (access control) Apache 2 or NGINX combined with Shibboleth for access control and communicating with the web browser
- (rich client) Web browsers now provide a rich software platform in their own right

With the above list we can build capable applications relying on the sophisticated features of our web browsers. 

## Hidden costs and missing bits

The problem with my off the shelf list so far is that it forces us to rely on JavaScript (or WASM module) running in the web browser to use the JSON API or the search engine that renders JSON results. It sounds like it should be easier because you are not writing anything that runs server side. In practice this is a horrible idea[^17].

[^17]: See <https://infrequently.org/2024/01/performance-inequality-gap-2024/> for a nice discussion of the problem.

What we should do is use Newt to tie those JSON services together, render the results using a template engine web service and hand HTML back to the web browser. Newt Router provides a means to tie the services together, Newt Mustache provides a template engine as web service. The Newt Router can also serve out our static content that might be included in the HTML results created by a Mustache template. Newt provides the missing bits so we don't need to send JavaScript down the virtual wire to assemble web content. Less bandwidth required on the network, few network accesses needed, the web browser knows how to quickly display HTML and CSS. It's a big step forward without writing much code, maybe with out writing any code if Newt's code generator does a sufficient job for your needs.

### A Newt baseline

Web services talk to other web services all the time. This isn't new. It isn't exotic it. LAS systems often do this too. Just thing of how they integrate content from external sources like CrossRef, DataCite, ORCID and ROR. Can we do this on the small scale too? Let me break that down further.

- Can we align access control with our front end web server?
- Can we insist our database management system provides flexible JSON API?
- Can we treat the output of one web service as the input for the next?
- Can we aggregate these into a data pipelines?
- Will that be enough to define a web application?

In 2024 for metadata curation apps I think the answer is "yes we can". Here's an example.

- Apache2 (or NGINX) as a front end web server responsible for controlling access
- Newt Router receives request from Apache2, routes it appropriately to our database JSON API, takes the result and run that through a template engine web service
- Newt Mustache provides the template engine as web service
- Postgres+PostgREST provide the database with JSON API

Newt's code generator is used to create the SQL to setup our Postgres database, a PostgREST configuration file and a set of Mustache templates.


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

The developer writes YAML to generate the back end data management and templates to render those results. You can enhance the generated code further if you want. I suspect that you'll spend most of your time improving the human experience of your application through improving the HTML markup in the templates, writing some CSS and perhaps enhancing behavior with the JavaScript run in the web browser. If you need to enhance the back end you work in SQL. If you simply need to improve the rendering of database results then you are working with Mustache templates. Otherwise you are focused on front end development.


## What about security, integration with single sign-on or other websites or services?

The `newtrouter` is a simple web service providing data routing based on its YAML configuration. It's a team player. In a production setting it should be used behind a front end web server like Apache 2 or NGINX. That latter can be configured to support single sign-on systems like Shibboleth[^18]. The front end web service controls access and handles securing the connection with the web browser. The front end web service proxies to the Newt router. Newt router receives requests and runs the data pipelines on localhost. The data pipelines can be composed of off the shelf software like Postgres+PostgREST, Solr and template engine to turn your JSON into a web page. Having a clear division of responsibilities helps in securing your web application. Since Newt router only knows how to talk to services on localhost you can keep it contained and prevent it from being used for nefarious actions off system. Like Newt router Newt Mustache is constrained to localhost for similar reasons.

Limiting Newt web service applications to localhost keeps them simple. Doing the minimum limits the attack surface for those who want to do mischief.  Neither `newtrouter` or `newtmustache` write to disk or require secrets. They only communicate via localhost using HTTP protocol.

[^18]: Shibboleth is a common single sign-on platform in research libraries, universities and colleges.

### What about "scaling"?

`newtrouter` is just a router. Aside from reading configuration at start up it doesn't maintain state. `newtmustache` functions the same way, read in the configuration and just run. By assigning different ports you can also run many instances of them. This makes it possible to run them in parallel, behind load balancer or even through proxying spread them across many machines. The instances don't share data or coordinate. They start up wait for a request and providing an answer. 

So what does this all mean? In principle a Newt based applications can scale big as the slowest element of your pipeline service. 

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

## Getting help

**The Newt Project is an experiment!!**. The source code for the project is supplied "as is". Newt is most likely broken. At a stretch it could be considered a working prototype. You should not use it for production systems.  However if you'd like to ask a question or have something you'd like to contribute please feel free to file a GitHub issue, see <https://github.com/caltechlibrary/newt/issues>. Just keep in mind it remains an **experiment** as of February 2024.

