
# Newt Project

The Newt Project is an experiment in rapid web application development for libraries, archives and museums (abbr: <abbr title="libraries, archives and museums abbrevation">LAS</abbr>).  Newt uses a service oriented architecture forming data pipelines[^1]. The pipelines compose the web application.

[^1]: A data pipeline is formed by taking the results from one web service and using it as the input to another web service. It is the web equivalent of Unix pipes. Prior art: [Yahoo! Pipes](https://en.wikipedia.org/wiki/Yahoo!_Pipes)

Taking this approach minimizes the software you need to write in favor of configuration. This is because "off the shelf" software is available to do the heavy lifting. Example, Postgres+PostgREST gives you an out of the box, web friendly, data management platform in the form of a JSON API. Light weight template engines like Pandoc running as a web service can transform you JSON API output into a web page. Throw in a full text search engine like Solr and you can check off the core features of many LAS software systems. What is missing is the layer to orchestrate the data flowing through a pipe.

The Newt Project is trying to encourage the following characteristics.

- preference for "off the shelf" over writing new code
- data modeling and management placed squarely in your database management system
- data sources accessed as JSON API
- data transformation (if needed) happens in simple stateless template engines
- leverage code generation when appropriate

The Newt Project is exploring this by providing three tools to fill in the gaps.

- `newt` is a stateless web service (a.k.a. micro service) that routes a web request through a data pipeline built from other web services
- `newtgen` is a code generator that can take a set of data models described in YAML and generate SQL, HTML and templates based on those models
- `newtmustache` is a simple stateless template engine inspired by Pandoc server that supports the Mustache template language

The Newt web services and data pipeline concept is being tested using

- [Postgres](https://postgres.org), data management and modeling
- [PostgREST](https://postgrest.org), a service that turns Postgres into a JSON API
- [Solr](https://solr.apache.org), full text search that provides a JSON API

Newt can tie these together through YAML expressing

- application definition (run time information and application metadata)
- routes (web requests differentiated by a HTTP method, headers and URL path)
- models (describe the data as would be input into a web form)

## What type of applications are supported by Newt?

Most LAS applications are focused on managing and curating some sort of metadata record. This is the primary target of Newt. This might be as simple as a controlled vocabulary or as complex as a archival or repository metadata record.

## Motivation

Over the last several decades web applications became very complex. This complexity is expensive in terms of reliability, enhancement, bug fixes and software sustainability.

> A brief historic detour to set context

Databases have been used to generate web pages since the early web.  Databases are well suited to managing data.  When the web became dynamic, databases continued to be use for data persistence. By 1993 the web as an application platform was born[^2] and with it a good platform for providing useful organizational software.


By the mid 1990s the Open Source databases, MySQL and Postgres, were the popular choice for building web applications. It is important to note neither MySQL or Postgres spoke HTTP[^3]. To solve this problem many people wrote software in languages like Perl, PHP and Python that ran inside the popular web server software called Apache. It was a pain to setup but once setup relatively easy to build things that relied on databases.  This led the web to explode with bespoke systems for curating and distributing web content. By the late 1990s and the early 2000s the practice of "mashing up" sites (i.e. content reuse) became the rage. As this increased in popularity web systems specialized further to cater to reuse. [Yahoo Pipes](https://en.wikipedia.org/wiki/Yahoo!_Pipes) was a very interesting expression of the "mashup culture"[^4]. Specialization has continued every sense. Some of these systems have become less bespoke. Eventual bespoke systems gave way to common use cases[^5]. A good example of a common use case is Apache's [Solr](https://solr.apache.org) search engine. Another example was the in-house bespoke content systems gave way to systems like [Drupal](https://drupal.org) and [WordPress](https://wordpress.org).


[^2]: Web applications proceeded to eat all the venerable green screen systems they could find. Eventually they and their corporate sponsors invented the surveillance economy we have today. Sometimes "good ideas" have terrible consequences. Making it easier to produce custom web applications should always be done keeping in mind the necessity for humane and inclusive use. Newt can be both part of a solution but also be used to exacerbate human problems.

[^3]: HTTP being the protocol the communicates with. Essentially at the time RDBMS spoke a dialect of SQL as the unifying language. The web of the time understood HTML and to a certain degree XML. By 2000 people were looking for something simpler than XML to move structured data about. [JSON](https://en.wikipedia.org/wiki/JSON) quickly became the answer.

[^4]: The basic concept was to make it easy to work with "data feeds" and combined them into a useful human friendly web pages. It even included a visual programming language to make it friendly to the non-programmer crowd. 

[^5]: If a use case is solved reliably enough it becomes "off the shelf" software.

> fast forward to 2024, context set

Much of the back end of web applications can largely be assemble from off the shelf software. Middleware however remains complex. I believe this to be a by product of inertia in software development practices and the assumption that what is good for "Google Scale" is good for everyone. 

I think a radical simplification is due.  Most software doesn't need to scale that large. Even in the research and LAS communities we don't routinely write software that scales as large as [Zenodo](https://zenodo.org/).  We don't typically support tens of thousands of simultaneous users. If you accept that premise then we can focus our efforts around orchestrating off the shelf components and put our remaining development efforts into improving the human experience of using our software.

A big key to simplification is realizing that the middleware no longer needs to be responsible for managing data models, access control and users or data transformation. There is no requirement for being bespoke. If you can configure the data routes and your data models the rest can become turn key.

## Off the shelf deliverables

- Postgres combined with PostgREST gives you an out of the box JSON API for managing data
- Solr gives you a powerful, friendly, JSON API for search and discovery
- Pandoc running as a web service provides a simple but powerful template engine
- Apache 2 or NginX when combined with a single sign on system (e.g. Shibboleth) provides access control
- Web browsers now provide a rich software platform in their own right

## This missing bits

With the above list can already can build complex applications that run inside the web browser. The cost is JavaScript is required to render everything. The trouble is relying on JavaScript assemble of content in the web browser is a horrible idea[^6]. A better approach is for the web browser to make a minimum number of requests to a single web service and get back useful results without having to process more than HTML and CSS. 

Taking the better approach in the past has required the writing of more middleware. I think we can avoid that or at least avoid complex middleware.

For over a decade web frameworks developed for programming languages like Go, Java, JavaScript, PHP, Python, and Ruby have relied on a concept of "routes". A "route" describes the URL path and HTTP method used to make a web request (e.g. your web browser requesting to view an HTML page). The mapping of a route to a function simplifies the model of receiving and responding to HTTP requests[^6]. The collection of routes and their functions compose the API your browser uses to navigate through your application. 

There is no requirement for the functions to be simple or complex. It depends on the task they are solving. Historically before single sign-on systems because common the function handling the request was responsible for the whole transaction. It need to hand access control, data validation, data formatting, storing or retrieving data from the database. In the later was used to persist data it was a big part of the code. You had to make sure the request from the public didn't lead to a compromise of your database's security. This was especially true when generate a SQL statement to interact with the database. Those functions that did it all are complex in when narrowly focused on an operation like creating data object, reading back a data object, updating a data object. Libraries and frameworks arose in part to address all the complexity by hiding it.
 
[^6]: That functionality is traditionally called a "handler" or "route handler".

I don't think that is required today. Web services talk to other web services all the time. What if we align access control with our front end web server or or inside the database itself? What if we could treat the response from one web service as the input of another creating pipe lines for web data? What if putting together pipelines was as easy as hosting static web content? I think the answer is "yes we can" to all these propositions.

Newt provides a simple web service that takes a YAML file and pairs the requests it receives with a data pipeline. The last stage of the pipeline executed is returned to the web browser. If there is a failure in the pipeline then an appropriate error is returned to the requesting web browser. Additionally Newt's web service can provide basic data validation by using its knowledge of the data models expressed as YAML.  Newt's code generator uses the same YAML to generate the SQL, templates and HTML needed to form a basic human user interface. There isn't much middleware left to write using this approach. If you do write middleware it can be narrowly focused to a specific stage of your pipelines.

[^6]: See <https://infrequently.org/2024/01/performance-inequality-gap-2024/> for a nice discussion of the problem.

## What comes with the Newt Project?

- [newt](newt.1.md) a [web service](https://en.wikipedia.org/wiki/microservices) designed for working with other "off the shelf" web services. It functions both as a data router and as a static file server. It is responsible for vetting requests against the models described in the YAML file. The same models used to generate the SQL for the database and the templates for user with a template engine.
- [newtgen](newtgen.1.md) is a command line program that reads the Newt YAML file and generates SQL and templates used to build your application.  The generated SQL currently targets Postgres+PostgREST.
- [newtmustache](newmustache.1.md) implements a simple lightweight template engine supporting [Mustache](https://mustache.github.io/) templates. Mustache template language is support by a wide variety of programming languages. It is provided as an alternative to using Pandoc as a template engine.

NOTE: See the [user manual](user_manual.md) for details

## How is Newt speeding up development?

The Newt suite of tools use a common YAML file.

1. `newtgen` can render SQL suitable for bootstrapping your Postgres+PostgREST database, templates and HTML 
2. `newtmustache` provides a simple, stateless, Mustache template engine
3. `newt` provides data routing between other web services that fill specific functions or roles.

Your back end is constructed from "off the shelf" parts. Newt provides the routing. It allows our customization efforts to focus on data modeling in the database and the HTML, CSS and JavaScript that will run in your web browser.

## Where is development time spent?

The developer writes YAML to generate the back end and the other services into an application. The developer generates the SQL and templates needed for the application. They may choose to customize those further. Newt web service provides the data pipeline management. Once the back end is configured the developer can focus on the code that runs in the web browser.

## What about security, single sign-on, integration with other websites or services?

The `newt` program is a simple web service providing data routing based on its YAML configuration. It's a team player. In a production setting it should be used behind a front end web server like Apache 2 or NginX. That latter can be configured to support single sign-on systems like Shibboleth[^7]. The front end web service controls access and handles securing the connection with the web browser. The front end web service proxies to the Newt web service. Newt web services runs the data pipeline on localhost. The data pipelines are off the shelf services like Postgres+PostgREST, Solr and Pandoc server. Having a clear division of responsibilities helps in securing your web application. Since Newt only knows how to talk to services on localhost you can keep it contained and prevent it from being used to doing nefarious things off system. Similarly if you decide to use Newt's template engine, `newtmustache`, it also is restricted to localhost. It is stateless and doesn't use secrets.

Limiting Newt web service applications to localhost keeps them simple. Only doing the minimum limits the attack surface for those who want to do mischief.  Neither `newt` or `newtmustache` write to disk or require secrets. They only communicate via localhost using HTTP protocol.

If you need to integrate a Newt application with an external service (e.g. CrossRef, ORCID or ROR) this can be done browser side or via a proxy mapped to localhost on your server.

[^7]: Shibboleth is a common single sign-on platform in research libraries, universities and colleges.

### What about "scaling"?

`newt` is just a data router. Aside from its configuration read at start up it doesn't maintain state. `newtmustache` doesn't maintain state and the only configuration is the port number it runs on. Both can safely run with many instances in parallel as needed on one or more machines. The instances don't share data or coordinate. Run waiting for a request and providing an answer. So what does this mean?

In principle a Newt based application can scale as wide as needed as long as each element in the pipeline(s) can also scale.

If the elements of your data pipeline can scale then your application can scale too. In our example Postgres can be configured as a <abbr title="high availability">HA</abbre> cluster. PostgREST, Pandoc server, Newt's mustache engine all are stateless and be run in parallel.

While my practical use of Newt is targeting the small it should scale up to meet high volume demands.

### Anatomy of a Newt based web application

Newt application development is friendly to version control systems (e.g. Git). It consists of a Newt configuration file, along with generate SQL files, HTML templates and any static web assets you've added. A typical disk layout of a Newt project could look like this-

- `/` project folder
  - `htdocs` this directory holds your static content needed by your web application
  - `*.sql` these are the SQL files used by your application to define your models and behaviors in Postgres
  - `*.tmpl` template files for either Pandoc or Mustache template engines
  - `newt.yaml` would hold the a Newt configuration file (this is an example name for the configuration file, you can name it whatever you like)
  - `CITATION.cff` and `codemeta.json` for project metadata

> Newt, a type of salamander, doesn't seek attention. It does its own thing. You only notice a salamander if you look carefully.

## About the Newt source repository

Newt is a project of Caltech Library's Digital Library Development group. It is hosted on GitHub at <https://github.com/caltechlibrary/newt>. If you have questions, problems or concerns regarding Newt you can use GitHub issue tracker to communicate with the development team. It is located at <https://github.com/caltechlibrary/newt/issues>.

## "Someday, maybe" exploration

- Integrate S3 protocol to support storing binary or large objects

## Getting help

**The Newt Project is an experiment!!**. The source code for the project is supplied "as is". Newt is most likely broken. At a stretch it could be considered a working prototype. You should not use it for production systems.  However if you'd like to ask a question or have something you'd like to contribute please feel free to file a GitHub issue, see <https://github.com/caltechlibrary/newt/issues>. Just keep in mind it remains an **experiment** as of February 2024.

## Documentation

- [user manual](user_manual.md)
- [INSTALL](INSTALL.md)
- [About](about.md) 

