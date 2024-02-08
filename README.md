
# Newt Project

The Newt Project is an experiment in rapid web application development for libraries, archives and museums (abbr: <abbr title="libraries, archives and museums abbrevation">LAS</abbr>).  Newt uses a service oriented architecture forming data pipelines[^0]. The pipelines compose the web application.

[^0]: A data pipeline is formed by taking the results from one web service and using it as the input to another web service. It is the web equivalent of Unix pipes. Prior art: [Yahoo! Pipes](https://en.wikipedia.org/wiki/Yahoo!_Pipes)

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

Over the last several decade web applications became very complex. This complexity is expensive in terms of reliability, enhancement, bug fixes and software sustainability.

> A brief historic detour to set context

Databases have been used to generate web pages since the early web.  Databases are well suited to managing data.  When the web became dynamic, databases continued to be use for data persistence. The web as an application was born and proceeded to eat all the venerable green screen systems it could find. When that was complete it continued to engorge itself until we have the web today and along with a problematic surveillance economy. Somewhere in that time frame the web became a good platform for providing useful organizational software.

> moving forward in time

In the 1990s the Open Source databases of choice were MySQL and Postgres. Neither MySQL or Postgres spoke HTTP (the protocol the web runs on). To solve this problem many people wrote software in languages like Perl, PHP and Python that ran inside the popular web server software called Apache. It was a pain to setup but once setup relatively easy to build things that relied on databases.  This led the web to explode with bespoke systems for curating and distributing web content. By the late 1990s and the early 2000s the practice of "mashing up" sites (i.e. content reuse) became the rage. As this increased in popularity web systems specialized further to cater to reuse. [Yahoo Pipes](https://en.wikipedia.org/wiki/Yahoo!_Pipes) was a very interesting expression of the "mashup culture". The basic concept was to "data feeds" and combined them into a useful human friendly web pages. Specialization has continued. Some of these systems have become less bespoke. Eventual bespoke systems gave way to common use cases. Common use cases lead to "off the shelf" programs. A good example is Apache's [Solr](https://solr.apache.org) search engine. 

> fast forward to 2024, context set

The back end of web applications can largely be assemble from off the shelf software. While middleware has grown large and complex that can be viewed as a by product of inertia and the assumption that what is good for "Google Scale" is good for everyone. 

I think a radical simplification is due.  Most software doesn't need to scale that large. Even in the research and LAS communities we don't routinely write software that scales as large as [Zenodo](https://zenodo.org/).  We don't typically support tens of thousands of simultaneous users. If you accept that premise then we can focus our efforts around orchestrating off the shelf components and put our remaining development efforts into improving the human experience of using our software.

A big key to simplification is realizing that the middleware no longer needs to be responsible for managing data models, access and users, data transformation and doesn't need to be bespoke. If you can configure the data routes and your data models the rest could become turn key.

## Off the shelf delivers

- Postgres combined with PostgREST gives you an out of the box JSON API for managing data
- Solr gives you a powerful, friendly, JSON API for search and discovery
- Pandoc running as a web service provides a simple but powerful template engine
- Apache 2 or NginX when combined with a single sign on system (e.g. Shibboleth) provides access control
- Web browsers now provide a rich software platform in their own right

## This missing bits

With the above list can already can build complex applications that run inside the web browser. The cost is JavaScript is required to render everything. The trouble is relying on JavaScript assemble of content in the web browser is a horrible idea[^2]. A better approach is then web browser can make a minimum number of requests to a single web service and get back useful results without having to process more than HTML and CSS. 

Taking the better approach in the past has required the writing of more middleware. I think we can avoid that or at least avoid complex middleware.

Over the last many years many web frameworks developed for Go, Java, JavaScript, PHP, Python, Ruby, etc. that rely on a concept of "routes". The routes form the applications expressed through a collection of predictable URLs.  A route can be though of as an HTTP method and a URL path. To build something for the web in these environments you wind up creating a bunch of functions that map a route to a data source, assemble it and then hand it back. Sometimes those functions are well encapsulated (e.g. static file access) requiring minimal coding. More commonly they become complex. They perform data modeling and validating. They may require coordination between data models. Often the middleware gets stuck also imposing a access management system and thus user profile manage. Of course they still need to eventually returning an appropriate response to the web browser. 

What if we align access crontrol with our front end web server or express it at the database level? What if we could treat the output of a data pipeline as easily as we provide static file access?  I think the answer is "yes we can" to both of these propositions. Newt provides a simple web service that takes a YAML file and pairs the requests it receives with a data pipeline. The last stage of the pipeline executed is returned to the web browser. If there is a failure in the pipeline then an appropriate error is returned to the requesting web browser. Additionally Newt's web service can provide basic data validation by using its knowledge of the data models available in the pipeline. Newt's code generator can create the SQL, templates and HTML needed to assemble a basic human user interface. If that is correct there isn't much middleware left to write.

[^2]: See <https://infrequently.org/2024/01/performance-inequality-gap-2024/> for a nice discussion of the problem.

## What comes with the Newt Project?

- [newt](newt.1.md) a [web service](https://en.wikipedia.org/wiki/microservices) designed for working with other "off the shelf" web services. It functions both as a data router and as a static file server. It is responsible for vetting requests against the models described in the YAML file. The same models used to generate the SQL for the database and the templates for user with a template engine.
- [newtgen](newtgen.1.md) is a command line program that reads the Newt YAML file and generates SQL and templates used to build your application.  The generated SQL currently targets Postgres+PostgREST.
- [newtmustache](newmustache.1.md) is a recent additional to the prototype suite. When discussing using Pandoc templates with colleagues some people pointed out they didn't like Pandoc templates. Newt's template engine uses [Mustache](https://mustache.github.io/) template markup to render JSON content. It functions like Pandoc in server mode. You POST a request to the template engine that includes both your template and JSON data and it returns the transformed result. This demonstrates that you can swap out the template engine and still use the Newt web service. 

NOTE: See the [user manual](user_manual.md) for details

## How is Newt speeding up development?

The Newt suite of tools use a common YAML file.

1. `newtgen` can render SQL suitable for bootstrapping your Postgres+PostgREST database, templates and HTML 
2. `newtmustache` provides a simple, stateless, Mustache template engine
3. `newt` provides a data routing

This leaves your back end to be constructed from "off the shelf" parts. Newt allows customization to focus on data models and HTML, CSS and JavaScript that will run in your web browser.

## Where is development time spent?

The developer writes YAML to generate the back end and the other services into an application. The developer generates the SQL and templates needed for the application. They may choose to customize those further. Newt web service provides the data pipeline management. Once the back end is configured the developer can focus on the code that runs in the web browser.

## What about security, single sign-on, integration with other websites or services?

The `newt` program is a simple web service providing data routing based on its YAML configuration. It's a team player. In a production setting it should be used behind a front end web server like Apache 2 or NginX. That latter can be configured to support single sign-on systems like Shibboleth[^2]. The front end web service controls access and handles securing the connection with the web browser. The front end web service proxies to the Newt web service. Newt web services runs the data pipeline on localhost. The data pipelines are off the shelf services like Postgres+PostgREST, Solr and Pandoc server. Having a clear division of responsibilities helps in securing your web application. Since Newt only knows how to talk to services on localhost you can keep it contained and prevent it from being used to doing nefarious things off system. Similarly if you decide to use Newt's template engine, `newtmustache`, it also is restricted to localhost. It is stateless and doesn't use secrets.

Limiting Newt web service applications to localhost keeps them simple. Only doing the minimum limits the attack surface for those who want to do mischief.  Neither `newt` or `newtmustache` write to disk or require secrets. They only communicate via localhost using HTTP protocol.

If you need to integrate a Newt application with an external service (e.g. CrossRef, ORCID or ROR) this can be done browser side or via a proxy mapped to localhost on your server.

[^2]: Shibboleth is a common single sign-on platform in research libraries, universities and colleges.

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

