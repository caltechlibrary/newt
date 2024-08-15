
# Newt Project

Newt is an experimental rapid application development. Specifically Newt is focused on creating web based metadata curation tools. These types of applications are commonly needed in galleries, libraries, archives and museums (abbr: GLAM). Newt makes creating these type of applications easier.

How does Newt do that? Newt generates an application by implementing a [service oriented architecture](https://en.wikipedia.org/wiki/Service-oriented_architecture) using a combination of off the shelf software and generated code.

You can think of a web application as a sequence of requests and responses. Traditionally a web browser contacts your web site or application then one of two things will happen. Your app knows the answer and hands back the result. Alternatively if it doesn't know the answer it tells you it can't do so (e.g. 404 HTTP STATUS CODE). With service oriented architecture your application has another option. Your application can contact another service and use that result to answer the request from the web browser.  Newt's implements a service oriented architecture by orchestrating data processing pipelines[^11].

[^11]: A data pipeline is formed by taking the results from one web service and using it as the input to another web service. It is the web equivalent of Unix pipes. Prior art: [Yahoo! Pipes](https://en.wikipedia.org/wiki/Yahoo!_Pipes)

With data pipelines we can accept a request and feed that request to one service then take its output and send it to the next service. Newt does this by providing a data router. Newt can manage the request sequence through a simple YAML described pipeline. While it is possible to create pipelines using Apache and NginX proxy features in practice that approach quickly becomes an unmanageable configuration problem. You could encapsulate clusters of processes in containers but this too becomes complex to manage. Newt's router can cut through the hairball of configurations and define pipelines per request route. With Newt's pipelines the last service completed hands its result back to Newt's router which returns the result to the web browser.

> The data pipelines are defined in Newt' YAML configuration file. The pipelines are managed by Newt's data router.

Why is this important? Much of the "back end" of a web application is already available as off the shelf software. Here is a short list of examples used by Caltech Library.

- [Postgres](https://postgres.org) and [PostgREST](https://postgrest.org) provides a JSON API for data management
- [Dataset](https://caltechlibrary.github.io/dataset), a simple JSON document store that has an out of the box JSON API provided by [datasetd](https://caltechlibrary.github.io/dataset/datasetd.1.html).
- [Solr](https://solr.apache.org) and Elasticsearch, OpenSearch can provide full text search as a JSON web service
- [ArchivesSpace](https://archivesspace.org) provides a JSON API web service
- [Invenio RDM](https://inveniordm.docs.cern.ch/) provides a JSON API web service
- [Cantaloupe IIIF Image server](https://cantaloupe-project.github.io/) an IIIF [API](https://iiif.io/api/image/3.0/)

__This is not an exhaustive list__. These types of applications can all be integrated into your application through configuring the connection in Newt's YAML file. Newt's router runs the data pipelines and can host static content.

> Wait, what about my custom metadata needs?

Metadata oriented applications share the following operations -- create, retrieve, update, delete and list. These are called CRUD-L[^12] features. Customization tend towards data models.  Newt's configuration file includes simple YAML description of your data models. It uses the data model to render configuration, middleware and templates. Newt's data router ties them all together into an application using service oriented architecture.

Newt provides:

- a simple data modeler
- code generation
- template engine
- data routing

You can extend your application through browser side enhancements, adding additional data routes and pipelines or through customizing the generated code.

[^12]: CRUD-L, acronym meaning, "Create, Read, Update, Delete and List" or alternately "Create, Retrieve, Update, Delete and List". These are the basic actions used to manage metadata.


## Does Newt clean my house or make cocktails?

Newt is a narrowly focused rapid application development toolbox. Newt will not clean your house or make a cocktail. Additionally it does not support the class of web applications that handle file uploads. That means it is not a replacement for Drupal, WordPress, Islandora, etc. Newt is for building applications more in line with ArchivesSpace but with simpler data models. If you need file upload support Newt is not the right solution at this time.

Newt applications are well suited to interacting with other applications that provide a JSON API. A service with a JSON API can be treated as a JSON data source. A JSON data source can easily be run through a pipeline. Many GLAM applications like ArchivesSpace and Invenio RDM provide JSON API. It is possible to extended those systems by creating simpler services that can talk to those JSON data sources. Newt is well suited to this "development at the edges" approach.

> What if those systems aren't available on localhost?

Services not available on localhost must be proxied to integrate with your Newt application. Why is this necessary? Short answer, security. Newt seeks to reduce the attack surface of your web application as much as possible. It does that by only trusting services offered directly on localhost.  Newt's application models presumes you're running behind a "front end web server" like Apache 2, NginX or Lighttd. These systems can be configured to provide access control as well as perform proxy services. The front end web server is your first line of defense against a cracker.

External web services integrate through a proxy setup (e.g ORCID, ROR, CrossRef or DataCite). This can be done via the front end web server or by writing a dedicated proxy service.  Today writing proxy services us easily accomplished in most popular programming languages due to good support for web protocols. This is true of Python, PHP, JavaScript, Go, Rust and many others.

## How does Newt impact web application development?

A Newt application encourages the following.

- modeling your data simply the result is expressed in YAML
- preference for "off the shelf" over writing new code
- use a database management system for managing your data and SQL functions
- prefer software that can function as a JSON data source
- transforming data representations by using a light weight template engine
- code generation where appropriate

## If Newt doesn't make cocktails, what is it bringing to the party?

In 2024 there is allot of off the self software to build on. Newt provides a few tools to fill in the gaps.

- `newt` is a development tool for data modeling, generating code, running Newt router, template engine and PostgREST or dataset collections
- `ndr` is a stateless web service (a.k.a. micro service) that routes a web requests through a data pipelines built from other web services
- `nte` is a stateless template engine inspired by Pandoc server that supports the [Handlebarsjs](https://handlebarsjs.com) template language and is designed to process data from a JSON data source

The Newt YAML configuration ties these together expressing

- data modeling (descriptions of data as you would provided in a web form)
- code generation
- data routes (web requests differentiated by a HTTP method and URL path that trigger processing in a data pipeline)
- template maps (path/template pairs used that can recieve JSON and render a results)

## What type of applications are supported by Newt?

Most GLAM applications are focused on managing and curating some sort of metadata records. Sometimes these metadata records are quite complex (e.g. ArchivesSpace, RDM records) but often they are simple (e.g. a list of authors, a list of citations).  Newt's primary target is generating applications to manage simple data models. Simple data models are those which can be expressed in web forms via HTML5 native input elements.

## Motivation

Over the last several decades web applications became very complex. This complexity is expensive in terms of reliability, enhancement, bug fixes and software sustainability. Newt is an attempt to address this by reducing the code you write and focusing your efforts on declaring what you want.

In 2024 the back end of web applications can largely be assemble from off the shelf software. Middleware however remains complex. I believe this to be a by product of inertia in software development practices and the assumption that what is good for "Google Scale" is good for everyone.

I think a radical simplification is due.  **Newt is intended to spark that conversation**. My observation is most software doesn't need to scale large. In the research and GLAM communities we don't routinely write software that scales as large as [Zenodo](https://zenodo.org/).  We don't typically support tens of thousands of simultaneous users. I think we can focus our efforts around orchestrating off the shelf components and put our remaining development time into improving the human experience using our software. A better human experience is an intended side effect of Newt.

A big key to simplification is narrowing the focus of our middleware. When our middleware has to implement everything it becomes very complex. Look at Drupal and WordPress. They implement data modeling, data management, user accounts, access management, data transformation.

I think our web services should be doing less, much less. Our web services should be narrowly focused. Conceptually simpler. Do one or two things really well. Newt enables using simpler discrete services to build our applications.

## Working with off the shelf deliverables

Take the following as a "for instance".

- (data management) Postgres combined with PostgREST gives you an out of the box JSON API for managing data
- (access control) Apache 2 or NGINX combined with Shibboleth for access control and communicating with the web browser
- (rich client) Web browsers now provide a rich software platform in their own right
- (optionally you can full text search) Solr gives you a powerful, friendly, JSON API for search and discovery

With the above list we can build capable applications relying on the sophisticated features of our web browsers. This is true even without using Newt. There is a problem though.  If we only use the above software to build our application we must rely on JavaScript (or WASM module) running in the web browser to interact with the server. This sounds simple. In practice this is a terrible idea[^17].

[^17]: See <https://infrequently.org/2024/01/performance-inequality-gap-2024/> for a nice discussion of the problem.

What we should do is use Newt to tie those JSON services together and send rendered HTML back to the web browser. Newt's router provides static file service and a means of pipelining our JSON data source through a template engine. Newt provides a Handlebars[^18] template engine for that purpose. Newt provides the missing bits from my original list so we don't need to send JavaScript down the wire to the web browser. The Newt approach uses less bandwidth, fewer network accesses and less computations cycles on your viewing device. The Newt approach takes advantage of what the web browser is really good at without turning your web pages into a web service. Newt YAML describes the system you want. You get the Newt capabilities without writing much code. Maybe without writing any code if Newt's code generator does a sufficient job for your needs.

[^18]: Handlebars is largely a superset of the Handlbars Template Language, see <https://handlebarsjs.com> and <https://mustache.github.io> for details.
### A "newt" baseline

Web services talk to other web services all the time. This isn't new. It isn't exotic. Newt scales down this approach to the single application.

- Can we align access control with our front end web server?
- Can we insist on our database management system providing a JSON API?
- Can we treat the output of one web service as the input for the next?
- Can we aggregate these into data pipelines?
- Will that be enough to define our web application?

In Spring 2024 for metadata curation apps I think the answer is "yes we can".

## What comes with the Newt Project?

The primary tools are.

- [newt](newt.1.md) a developer tool for building a Newt based application which includes modeling, code generation support, and runtime support
- [ndr](ndr.1.md) a [web service](https://en.wikipedia.org/wiki/microservices) designed for working with other "off the shelf" web services. It functions both as a router and as a static file server. It does this by routing your request through a YAML defined pipeline and returning the results. Typically this will be a JSON data source and running that output through a template engine like Newt's **nte**.
- [nte](nte.1.md) provides a Handlebars template engine for transforming JSON data into a more human friendly format.

See the [user manual](user_manual.md) for details.

## About the Newt source repository

Newt is a project of Caltech Library's Digital Library Development group. It is hosted on GitHub at <https://github.com/caltechlibrary/newt>. If you have questions, problems or concerns regarding Newt you can use GitHub issue tracker to communicate with the development team. It is located at <https://github.com/caltechlibrary/newt/issues>. The name comes from wanting a "[New t]ake" on web application development.

## Getting help

**The Newt Project is an experiment!!**. The source code for the project is supplied "as is". Newt is a partially implemented prototype (May 2024). However if you'd like to ask a question or have something you'd like to contribute please feel free to file a GitHub issue, see <https://github.com/caltechlibrary/newt/issues>.

Currently Newt is targeted at Windows on X86, macOS on x86 and ARM (e.g. M1, M2) and Linux on aarch64 (ARM 64) and x86.

## Building from source

Newt is experimental so doesn't have installers yet. If you are compiling from source the following software is required.

1. Git to retrieve the Newt repository
2. Golang >= 1.22 (used to compile command line programs and services except newthandlebars)
3. Deno >= 1.44 (used to build newthandlebars, note not available on Windows ARM or 32bit Linux)
4. Pandoc >= 3.1 (used to build version.go, version.ts and documentation)
5. GNU Make and Bash shell (you can cross compile for Windows and macOS from Linux)

Steps to compile

1. From your home directory clone the repository
2. change into the repository directory
3. Run Make
4. Run make test
5. Run make install

~~~shell
cd $HOME
git clone git@github.com:caltechlibrary/newt src/github.com/caltechlibrary/newt
cd  src/github.com/caltechlibrary/newt
make
make test
make install
~~~
