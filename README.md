
# Newt Project

Newt is an experimental set of tools for rapid application development. Specifically Newt is focused on web based metadata curation tools. These types of applications are commonly needed in galleries, libraries, archives and museums (abbr: GLAM). Newt makes creating these type of applications easier.

How does Newt do that? Newt generates applications implementing a [service oriented architecture](https://en.wikipedia.org/wiki/Service-oriented_architecture).

You can think of a web application as a sequence of requests and responses. Typically when a web browser contacts your web application one of two things will happen. Your app knows the answer and hands back the result. With the service oriented architecture your application has another option. Your application can contact another service and use that result to answer the request from the web browser.  Newt's enhancement to a basic service oriented architecture is to provide an easy means to construct data processing pipelines[^11]. 

With data pipelines we can accept a request and feed that request to one service then take its output and send it to the next service. Newt does this by providing a data router. Newt can manage the request sequence through a simple YAML description of the pipeline. While it is possible to create pipelines using Apache and NginX proxy features in practice quickly becomes an unmanagable configuration problem. You could encampsulate clusters of processes in containers but this too becomes complex to manage. Newt's router can cut through the hairball of configurations and define pipelines per route request. With Newt a pipeline the last service completed hands its result to the Newt Router which returns the result to the web browser.

Why is this important? Much of the "back end" of a web application is already available as off the shelf software. Here is a short list of examples used in Caltech Library.

- [Postgres](https://postgres.org) and [PostgREST](https://postgrest.org) provides a JSON API for data management
- [Solr](https://solr.apache.org) and Elasticsearch, OpenSearch can provide full text search as a JSON web service
- [ArchivesSpace](https://archivesspace.org) provides a JSON API web service
- [Invenio RDM](https://inveniordm.docs.cern.ch/) provides a JSON API web service
- [Cantaloupe IIIF Image server](https://cantaloupe-project.github.io/) an IIIF [API](https://iiif.io/api/image/3.0/)

This is not an exhaustive list. These types of applications can all be integrated into your application through configuring the connection in Newt's YAML file. Newt Router runs the data pipelines.

> Wait, what about my custom metadata needs?

That role can be filled by a JSON data source. In the second Newt prototype our focus is on using Postgres+PostgREST as a JSON data source. Newt's code generation lends a hand here. Using Newt's YAML file the code generator can generate SQL and configuration for setting up Postgres+PostgREST. Newt's code generator can also render Mustache templates used for interacting with the JSON data source. Between the generated SQL and Mustache templates Newt and Postgres+PostgREST provides the basic CRUD-L[^12] operations for managing custom metadata. You can enhance your Newt application further through customizing the generated Mustache templates and SQL or by adding routes and pipelines leveraging other JSON data sources (e.g. Solr).

## Does Newt clean my house or make cocktails?

Newt is a narrowly focused rapid application development toolbox. Newt will not clean your house or make a cocktail. Additionally it does not support the class of web applications that handle file uploads. That means it is not a replacement for Drupal, WordPress, Islandora, etc. Newt is for building applications more in line with ArchivesSpace but with simpler data models. If you need file upload support Newt is not the right solution at this time.

Newt applications are well suited to interacting with other applications that provide a JSON API. A service with a JSON API can be treated as a JSON data source. A JSON data source can easily be run through a pipeline. Many GLAM applications like ArchivesSpace and Invenio RDM provide JSON API. It is possible to extended those systems by creating simpler services that can talk to those JSON data sources. Newt is well suited to this "development at the edges" approach.

What if those systems aren't available on localhost? In that case you need to provide a proxy from these services to localhost (e.g. via Apache 2, NGINX or Python script). You would write your Newt YAML file pointing at the localhost end point. This approach can be used to integrated external systems like ORCID, ROR, CrossRef, DataCite, Open Alex, Pub Med Central.

[^11]: A data pipeline is formed by taking the results from one web service and using it as the input to another web service. It is the web equivalent of Unix pipes. Prior art: [Yahoo! Pipes](https://en.wikipedia.org/wiki/Yahoo!_Pipes)

[^12]: CRUD-L, acronym meaning, "Create, Read, Update, Delete and List". These are the basic actions used to manage metadata or objects.

## How does Newt impact web application development?

A Newt application encourages the following.

- preference for "off the shelf" over writing new code
- modeling your data simply
- use a database management system for managing your data
- prefer software that can function as a JSON data source
- transforming data representations (if needed) by using a light weight template engine
- code generation where appropriate

## If Newt doesn't make cocktails, what is it bringing to the party?

In 2024 there is allot of off the self software to build on. Newt provides a few tools to fill in the gaps.

- `newt` is a development tool for generating code, running Newt Router, Newt Mustache and PostgREST so you can test your application
- `newtrouter` is a stateless web service (a.k.a. micro service) that routes a web requests through a data pipelines built from other web services
- `newtmustache` is a stateless template engine inspired by Pandoc server that supports the Mustache template language and is designed to process data from a JSON data source.

Newt's 2nd prototype is being tested building applications based on

- [Postgres](https://postgres.org), data management
- [PostgREST](https://postgrest.org), a service that turns Postgres into a JSON API

The Newt YAML ties this together expressing

- applications (run time information for Newt Router, Newt Mustache and PostgREST)
- models (descriptions of data as you would provided in a web form)
- routes (web requests differentiated by a HTTP method and URL path that trigger processing in a data pipeline)
- templates (pairs a request with a template to transform a JSON into some other format such as an HTML document)


## What type of applications are supported by Newt?

Most GLAM applications are focused on managing and curating some sort of metadata records. Sometimes these metadata records are quite complex (e.g. ArchivesSpace, RDM records) but often they are simple (e.g. a list of authors, a list of citations).  Newt's primary target is generating applications to manage simple data models. Simple data models are those which can be expressed in HTML5 web forms via HTML native input elements.

## Motivation

Over the last several decades web applications became very complex. This complexity is expensive in terms of reliability, enhancement, bug fixes and software sustainability. Newt address this by reducing the code you write and focusing your efforts on declaring what you want.

> A brief historic detour to set context

Databases have been used to generate web pages since the early web.  Databases are well suited to managing data.  When the web became dynamic, databases continued to be use for data persistence. By 1993 the web as an application platform was born[^13] and with it a good platform for providing useful organizational and institutional software.

By the mid 1990s the Open Source databases like MySQL and Postgres were popular choices for building web applications. It is important to note neither MySQL or Postgres spoke HTTP[^14]. To solve this problem many people wrote software in languages like Perl, PHP and Python that ran inside the popular Apache web server. It was a pain to setup but once setup relatively easy to build things that relied on databases.  This led the web to explode with bespoke systems. By the late 1990s and early 2000s the practice of "mashing up" sites (i.e. content reuse) was the rage. Bespoke systems took advantage of content reuse too. [Yahoo Pipes](https://en.wikipedia.org/wiki/Yahoo!_Pipes) was a very interesting expression of the "mashup culture"[^15]. Yahoo Pipes inspired Newt's data pipelines.  Eventual the bespoke systems gave way to common use cases[^16]. A good example of a common use case is Apache's [Solr](https://solr.apache.org) search engine. Another example is how bespoke content management systems gave way to [Plone](https://plone.org), [Drupal](https://drupal.org) and [WordPress](https://wordpress.org).

[^13]: Web applications proceeded to eat all the venerable green screen systems they could find. Today's web is mired in surveillance tech and complex solutions. It has drifted far from Sir. Tim's vision of sharing science documents. We need to refocus on the "good ideas" and jettison the complexity that came with the surveillance economy. Newt can be part of that solution. Develop your Newt application with consideration for others.

[^14]: HTTP being the protocol the communicates with. Essentially at the time RDBMS spoke a dialect of SQL as the unifying language. The web of the time understood HTML and to a certain degree XML. By 2000 people were looking for something simpler than XML to move structured data about. [JSON](https://en.wikipedia.org/wiki/JSON) quickly became the answer.

[^15]: The basic concept was to make it easy to work with "data feeds" and combined them into a useful human friendly web pages. It even included a visual programming language to make it friendly to the non-programmer crowd.

[^16]: If a use case is solved reliably enough it becomes "off the shelf" software.

> fast forward to 2024, context set

Much of the back end of web applications can largely be assemble from off the shelf software. Middleware however remains complex. I believe this to be a by product of inertia in software development practices and the assumption that what is good for "Google Scale" is good for everyone.

I think a radical simplification is due.  **Newt is intended to spark that conversation**. My observation is most software doesn't need to scale large. Even in the research and GLAM communities we don't routinely write software that scales as large as [Zenodo](https://zenodo.org/).  We don't typically support tens of thousands of simultaneous users. If you accept that premise then we can focus our efforts around orchestrating off the shelf components and put our remaining development efforts into improving the human experience of using our software. A better human experience is an intended side effect of Newt.

A big key to simplification is narrowing the focus of our middleware. When our middleware has to implement everything it becomes very complex. Look at Drupal and WordPress. They implement data modeling, data management, user accounts, access management, data transformation.

I think our web services should be doing less, much less. Our web services should be narrowly focused. Conceptually simpler. Do one or two things really well. Newt enables using simpler discrete services to build our applications.

## Working with off the shelf deliverables

Take the following as a "for instance".

- (data management) Postgres combined with PostgREST gives you an out of the box JSON API for managing data
- (full text search) Solr gives you a powerful, friendly, JSON API for search and discovery
- (access control) Apache 2 or NGINX combined with Shibboleth for access control and communicating with the web browser
- (rich client) Web browsers now provide a rich software platform in their own right

With the above list we can build capable applications relying on the sophisticated features of our web browsers. This is true even without using Newt. There is a problem though.  If we only use the above software to build our application we must rely on JavaScript (or WASM module) running in the web browser to interact with the server. This sounds simpler. In practice this is a terrible idea[^17].

[^17]: See <https://infrequently.org/2024/01/performance-inequality-gap-2024/> for a nice discussion of the problem.

What we should do is use Newt to tie those JSON services together and send rendered HTML back to the web browser. Newt Router provides static file service and a means of pipelining our JSON data source through a template engine. Newt Mustache provides a template engine. Newt provides the missing bits from my original list so we don't need to send JavaScript down the wire to the web browser. The Newt approach uses less bandwidth, fewer network accesses and less computations cycles on your viewing device. The Newt approach takes advantage of what the web browser is really good at without turning your web pages into a web service. Newt YAML describes the system you want. You get the Newt capabilities without writing much code. Maybe without writing any code if Newt's code generator does a sufficient job for your needs.

### A new baseline

Web services talk to other web services all the time. This isn't new. It isn't exotic. Newt scales down this approach to the single application.

- Can we align access control with our front end web server?
- Can we insist on our database management system providing a JSON API?
- Can we treat the output of one web service as the input for the next?
- Can we aggregate these into data pipelines?
- Will that be enough to define our web application?

In Spring 2024 for metadata curation apps I think the answer is "yes we can".

## What comes with the Newt Project?

The primary tools are.

- [newt](newt.1.md) a developer tool for building a Newt based application which includes code generation support
- [newtrouter](newtrouter.1.md) a [web service](https://en.wikipedia.org/wiki/microservices) designed for working with other "off the shelf" web services. It functions both as a router and as a static file server. It does this by routing your request through a YAML defined pipeline and returning the results. Typically this will be a JSON data source and running that output through a template engine like Newt Mustache.
- [newtmustache](newmustache.1.md) implements a simple lightweight template engine supporting [Mustache](https://mustache.github.io/) templates. Mustache template language is well support by a wide variety of programming languages including Python, PHP, Perl and JavaScript.

Some additional tools are also provided. See the [user manual](user_manual.md) for details.

## About the Newt source repository

Newt is a project of Caltech Library's Digital Library Development group. It is hosted on GitHub at <https://github.com/caltechlibrary/newt>. If you have questions, problems or concerns regarding Newt you can use GitHub issue tracker to communicate with the development team. It is located at <https://github.com/caltechlibrary/newt/issues>. The name comes from wanting a "[New t]ake" on web application development.

## Getting help

**The Newt Project is an experiment!!**. The source code for the project is supplied "as is". Newt is a partially implemented prototype (April 2024). However if you'd like to ask a question or have something you'd like to contribute please feel free to file a GitHub issue, see <https://github.com/caltechlibrary/newt/issues>.

