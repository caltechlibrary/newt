
# Newt

Newt is an experimental [microservice](https://en.wikipedia.org/wiki/Microservices) designed for working with other "off the shelf" microservices. The primary purpose of Newt is to function as a localhost data router between a data source and rendering engine. The goal of the project is to create a rapid development platform through existing microservices targetting web applications suitable in libraries, archives, colleges and museums. The name "Newt" comes from the phrase "new take".

Newt comes with three commands. `newt` is a web service it is designed to sit behind your favorite web server (e.g. Apache 2, NginX) and route requests from the browser to a data source (e.g. JSON API) and optionally take the result and run it through rendering engined (e.g. Pandoc running in server mode).  `newtmustache` is a Pandoc server inspired rendering microsevice based implementing Mustache template support. `newtpg` is a command line program designed to generate SQL used bootstrap a JSON API built with PostgREST and Postgres.  Combined `newt` and `newtpg`  saves you from building yet another middleware microservice. Instead your development time is focused instead on three areas.

1. Modeling your data using SQL
2. Rendering content using simple templates (e.g. Pandoc or Mustache templates)
3. Enhancing the user expereience browser side using HTML5, CSS and if needed JavaScript

Newt achieves this division of responsibity through a YAML configuration file that describes your data models and the mapping of requested URL paths to a multistage pipe line of data sources and render engine.  Data sources are typically JSON API. The currently support render engine is Pandoc running as a microservice. Newt was inspired by my work with PostgREST and Postgres which allows you to quickly build a fully featured JSON API in the Postgres database using SQL. PostgREST+Postgres functions as an "off the shelf" data source.  Newt is design to work with "off the shelf" JSON data sources so it also works well with Solr, Elasticsearch and Opensearch.  Support for S3 object stores is in the early planning stages.

## Steps to build an a Newt web application

The following example describes using PostgREST+Postgres to provide a JSON API, Pandoc server as your render engine and Newt to serve static content and route requests to the JSON API and render engine.

1. Create a directory/folder for your project and change into that directory (setup version control if desired)
2. Create a Newt configuration file describing the models and routes your application will need
3. Use `newtpg` and the Newt configuration file to generate SQL to bootstrap your PostgREST+Postgres JSON service
4. Add/update your initial SQL data models and test inside Postgres (psql works nicely for this)
5. Add/update the routes in the routes section of the YAML configuration file if needed
6. Add/update any Pandoc templates, HTML, CSS and JavaScript as needed for your application
7. (re)start PostgREST and `newt` then test with your web browser

You repeat steps 4 through 7. This is where your application development time is spent. As your data models in Postgres stablize your focus as a developer can then shift naturally to providing an effective human facing experience working primarily with the templates and static web assets.

With Newt there's no more writing middleware, no need to reach for an ORM, not even a need to send JavaScript to the browser unless you choose to. Newt routes the request so you can focus on modeling your data in SQL and content presentation via Pandoc templates and static web assets.

## What about security, single sign-on, integration with websites or services?

The `newt` is a simple microservice providing data routing based on its configuration at startup. It's a team player.  At many universities, colleges, research libraries, archives and museums there is an existing single sign-on mechanism like Shibboleth running along with Apache 2 or NginX web servers.  Newt would run behind those services via reverse proxy. Newt itself doesn't know about users, it only routes data. Newt after reading the configuration file doesn't maintain state.  While the log output can contain identifiable information (e.g. IP address of request) or the JSON data source could contain sensitive information Newt doesn't retaining it. It just routes the data and gets out of the way.

A typical `newt` production setup might look like this

1. NginX with Shibboleth controls access to web site resources and where appropriate proxies to `newt`
2. `newt` responds to requests and maps those to a data source (e.g. PostgREST+Postgres JSON API) and gets back a response
3. `newt` can take a data source response and send it to a render engine (e.g. Pandoc in server mode or `newtmustache`)
4. `newt' assembled result is handed back to the NginX web server to pass onto the requesting web browser

In the example securing your application can happen both at the NginX level (e.g. requiring single sign-on) and
at the JSON API level via Postgres's management of PostgREST responses. NginX can also be used to proxy external resources you may wish Newt to route to.

Newt tries to do as little as possible while still providing data routing and static content services. This reduces the attack surface. Newts log output is written to standard out and standard error and not directly to disk to avoid the problem of high traffic logging filling up your disk drive. Logging can easily be captured by your servers logging system (e.g. systemd and it's log handling).  Newt is configured through the environment so does not require storing secrets in the configuration file.  Newt only reads the configuration at startup and can not write it back to disk. In fact Newt can't write to disk at all.  Newt focuses on data routing at the web application level only. It has a limited two stage pipeline for requests processing. For any given route defined in Newt's configuration it can contact a data source (e.g. PostgREST+Postgres JSON API) and optionally send the result through a rendering service (i.e. Pandoc running in server mode) before handing that result back to the requestor. Newt only knows how to speak http and s3 protocols for systems running on localhost. It only allows that for the routes defined when Newt starts up.

Keeping Newt simple minimizes the attack surface and keeps Newt a team player in your microservice based application.

### What about "scaling"?

Newt is just a data router. You can safely run many Newt instances in parallel as needed. They can run them on the same machine or separate machines. The instances don't share data or coordinate. They just read their configuration files when they start up and route data accordingly.

A typical Newt stack of PostgREST+Postgres and Pandoc also can scale up. You can run as many instances of Pandoc server or PostgREST as you need. You can spread them across machines. They are both stateless systems like Newt. The Postgres database provides consistency and itself can be configured in an HA cluster.

When I created Newt I was interested in small scale applications but since Newt is a microservice it scales as wide as you need just like Pandoc server, PostgREST and Postgres.

### Annatomy of a Newt based web application

Newt application development is friendly to version control systems (e.g. Git). It consists of a Newt configuration file, SQL files, Pandoc templates and any static web assets you need. A typical disk layout of a Newt project could look like this-

- `/` project folder
  - `htdocs` this directory holds your static content needed by your web application
  - `*.sql` these are the SQL files used by your application to define your models and behaviors in Postgres
  - `templates` this directory holds your Pandoc or Mustache templates
  - `tests` this directory holds your tests of your data model
  - `application.yaml` would hold the a Newt configuration file (this is an example name for the configuration file)

Data flow and dividing responsibilities in the Newt stack

- front end web server controls access, proxies to Newt
- Newt receives requests and maps them to a static resource or a data source
  - A data source could be a JSON API (e.g. PostgREST+Postgres, Solr, Elasticsearch, Opensearch) accessed via http over localhost
- Newt receives a response from a data source and can send the results to Pandoc server
- Newt hands back a result to your front end web server

## Origin story

Newt came about when I realized that all I needed was a data router that could map a web browser request to the JSON API provided by PostgREST+Postgres and Pandoc running as a service. That setup could replace most of the applications I'd written for the last decade or two. It would fit most of the small web applications I had previously built in PHP or Python for my library. Generalizing the concept of simple data router for a data source and render engine also meant I had an easy integration point for most the institutional software we currently run in our library. So I wrote a data router to do just that.

I demonstrated the Newt concept to my colleagues with a prototype. The prototype talked to a JSON API provided by PostgREST+Postgres and used Pandoc running in server mode for a rendering engine. I got some polite supportive comments. No one was particularly excited by it. I demonstrated a prototype Newt at a my local SoCal Code4Lib group. There people were excited by PostgREST and Postgres and not so excited about data routing. This was discouraging. I thought I was barking up the wrong tree. Eventually I realized the ambivalence of the router was a type of success. Newt isn't exciting. Newt should never be exciting. It just routes data! You configure it and forget it. It just runs.

The important take away was I had failed to appreceate how Newt successfully shifted the discussion from programming language frameworks, package management and build systems to to modeling data in SQL and using simpler HTML5, CSS and JavaScript for display.

> Newt, a type of salamander, doesn't seek attention. It does its own thing. You only notice a salamander if you look carefully.

## System design choices

Demonstrating the "Newt stack" has resulted in questions. I think I can address the big four.

1. Why SQL? Why build your data models with SQL?
2. Pandoc as rendering engine, why?
3. Why YAML for configuration?
4. How do I handle file uploads in my web application if I'm so focused on Postgres and SQL?

Here's my answer the first question. I think knowing at least some SQL is unavoidable as a web developer. While many have adopted an ORM[1] to generate SQL the models and manage data the resulting SQL is often far from ideal. This remains true in 2024 even though ORM implementations have been around decades.  The real problem with the ORM approach isn't inefficiency at all. The real problem is it obscures the data model and that discourages data re-use.  The ORM is a layer of abstraction on a system that itself is a set of abstractions already. I feel you can skip the extra layer, keep things simplier and avoid additional brittleness of tieing your data models to closely to a single application version. Focus on SQL to describe the data model and how to manage it.  In 2024 SQL remains the common language to manage data in a database. Even for non-relational data this has become true[2]. SQL may look ugly or quirky but it definately has legs and plans to stick around for a very long time. Let's embrace it!

For libraries, archives and museums managing metadata about our collections is critical.  Ditching the ORM and focusing on SQL increases our ability to share techniques with non-developer colleagues. The simple act of teaching the SQL SELECT statement can be liberating for someone who has only curated data via a web form or spreadsheets. SQL is well suited to the approach of incremental learning. Learn as much as you need when you need it.

The second question I run into is why Pandoc?  Each of the programming languages I used over the three decades offer some of template language(s). Even PHP which started out as a template language has template languages and frameworks! In the meantime no common template systems has emerged[3] which is language agnostic and widely adopted. Inspite of that designing sites with templates remains a mainstay in web application development. If you're going to use a template language and engine, which one? Inventing a new one doesn't help situation. I want one that can function as a microservice standing on it's own.

In the data science and library circles I travel I've seen a huge adoption of Pandoc for static site generation. When I eventually stumbled on the feature that Pandoc can run as a web service it hit me.  I could Pandoc as a rendering engine.It doesn't even require configuration. That lead me to pick it to be the rendering engine of Newt[4]. If another template language and engine comes on the scene, Newt can be adapted to us it instead[5].

[1]: ORM, object relational mapper. An ORM maps programs objects to a SQL syntax. With the advent of JSON columns in SQL tables this is rearly a problem anymore. Just focus on SQL.

[2]: I've heard of people querying S3 buckets using a SQL SELECT statement and used utilities in the shell to do the same with text files.

[3]: Mustache templates have become common in most langauges but haven't seen to take a hold like JSON did for structured data representation. The people working on PostgREST have another project to embed Mustache tempaltes in Postgres. Something to keep an eye one.

[4]: Using Pandoc for the render engine could change in the future as long as a common way to package the results from template, options and data source response is as simple as it is with Pandoc server.

[5]: The group behind PostgREST is also developing an Postgres embedable Mustache engine, that may make sense too.
Newt's configuration langauge is YAML. YAML was picked because it is widely use in 2023/2024. I don't need to explain it to my colleagues or peers they already use it. Newt implements a "domain specific language" or DSL on top of YAML to support rendering SQL data models targetting Postgres. YAML is also an easy language to use to describe the information needed for data routing in Newt. YAML seemed a good fit for Newt.

### Newt's minimal feature set

- `newt` as a two stage data router
- `newt` as static file service
- `newtpg` can use the Newt configuration file to render simple data models as SQL generator suitable to bootstrap a PostgREST+Postgres JSON API

Here's the data flow steps

1. Web browser => (Web Server proxy) => Newt
2. Newt => data source (e.g. Postgres + PostgREST or S3 Object store)
3. Newt => Pandoc (optional step)
4. Newt => (Web Server proxy) => Web browser

The person developing with Newt writes SQL to define the back end, may write Pandoc templates if that is desired and builds the front end with standard static web assets (e.g. HTML pages, CSS, JavaScript). Newt can support traditional websites and single page applications. It just saves writing a whole bunch of services that already exist.

## Orchestrating your app with Newt

Newt is configured with a YAML file. Currently the configuration file uses five main attributes and a collection of sub attributes.

htdocs
: The path to the htdocs directory holding any static assets (e.g. CSS, JavaScript, HTML files, image assets)

env
: A list of environment variables available to Newt's routes and models (this is how Newt avoids storing secrets)

routes
: An object describing the mapping of an HTTP request to JSON data source and and optional Pandoc server processing

namespace
: This is the schema name used to interact with PostgREST+Postgres

models
: This is a list of data models used by Newt to generate bootstrap SQL code for PostgREST+Postgres

The **htdocs** just points at a standard directory holding your static web content. It has no sub attributes.

The **env** attributes holds a list of environment variable names that can be used by Newt when defining **routes**.

Under **routes** you will also use the following attributes

var
: (optional) A list of variable names and types used in validating a request path or web from submission

req_path
: A expression describing a URL path received by Newt (typical made by a web browser or proxied from the front end web server)

req_method
: An HTTP method (e.g. GET, POST, PUT, PATCH, DELETE) related to the req_path being handled

api_url
: The URL expression used by Newt to contact the JSON data source for the route described by req_path and req_method. Newt can communicate using one of two protocols identified in the URL, `http://` or `s3://`.

api_method
: The HTTP method (e.g. GET, POST, PUT, PATCH, DELETE) of the JSON data source associated api_url for the given route

api_content_type
: The HTTP content type expression used when submitting the request to the JSON data source

pandoc_template
: (optional) The path to the pandoc template used to process the results of the JSON data source request results

pandoc_options
: (optional) Any options to pass when building the request to Pandoc server

The **models** attribute holds a list of models expressed in Newt's data model DSL. Models are optional but can be used to by Newt to generate bootstrap SQL code for use with PostgREST+Postgres. This is very experimental (2024) and likely to change as usage of Newt increases. Each model has a `model` attribute holding the models name (conforming to a variable name found in langauges like JavaScript, Python, or Lua). Each model also contains a `var` attribute which is a list of key/value pairs. The key/value pairs are made from a variable name (key) and type definition (value). The type definitions are mapped to suitable Postgres SQL schema when generating table definitions. Example models used for groups and people metadata. The asterix at the end of a type string indicates this is to be used as a key when looking up the object.

```yaml
namespace: groups_and_people
models:
- model: cl_person
  var:
    family_name: String
    given_name: String
    orcid: ORCID
    ror: ROR
    created: Timestamp
    updated: Timestamp
- model: cl_group
  var:
    cl_group_id: String*
    short_name: String
    display_name: String
    description: Text
    contact: EMail
    created: Timestamp
    updated: Timestamp
    founded: Date 2006-01-02
    disbanded: Date 2006-01-02
    approx_founding: Boolean
    active: Boolean
    website: URL
    ror: ROR
    grid: String
    isni: ISNI
    ringold: String
    viaf: String
```

The models and namespace attributes are used when generate SQL for PostgREST+Postgres.  The type strings are used in both generating SQL and also when embedded in a route definition to vet requests and fail early before contact the data source if the required information is not provided.

### Handling errors

Newt vets the initial request before contacting the JSON data source. If the request has a problem it will return an appropriate HTTP status code and message.  If the request to the JSON data source has a problem, it will pass through the HTTP status code and message provided by the JSON data source.  Likewise if Pandoc server has a problem Newt will forward that HTTP status code and message. If either the JSON data source or Pandoc server is unavailable Newt will return a "Gateway" http status code and message.

### Static file support

Newt first checks if a request is matched in one of the defined routes. If not it'll try to service the content from the "htdocs" location if that is defined in the configuration. If the file is not found or an htdocs directory has not been specified a http status of 404 is returned.

Note Newt's static file services are very basic. You can't configure mime type responses or modify behavior via "htaccess" files. If Newt is running behind a traditional web server like Apache 2 or NginX then you could use that service to host your static content providing additional flexibilty.

### Handling secrets, scaling and limitations

Newt's YAML file does not explicitly contain any secrets. This was intentional.  You may need to pass sensitive data to your JSON data source for access (e.g. credentials like a username and password). This should be provided via the environment and the YAML file needs to include these environment variable names in the "env" attribute.  The environment variables can be used to construct the URLs to contact the JSON or S3 data sources. There is still a risk in that theoretically that data source could echo return sensitive information. Newt can't prevent that. Newt is naive in its knowledge of the data source content it receives or hands off to Pandoc.  You're responsible for securing sensitive information at the database or s3 data source level. Follow the recommendations in the Postgres community around securing Postgres.

While Newt was conceived to be used on as a small scale web application platform for libraries, archives and museums it is capable of scaling big as long as your data source(s) can scale big.  Using the initial "Newt stack" elements can all be run easily behind load balancers and in parallel across machines. Newt is transactional. It does not require synchronized or shared of data between instances. Similarly PostgREST and Pandoc services are transactional and do not require shared state to function in parallel. Postgres itself can be configured in a HA cluster to support high availability and high demand. It should be possible to scale a Newt based application as large as those systems can be scaled.

Presently Newt does not supports file uploads. The plan is to integrated support for an S3 object store. That support is still very much in the planning stages.

Newt runs exclusively as a localhost service. In a production setting you'd run Newt behind a traditional web server like Apache 2 or NginX. The front end web service can provide access control via basic auth or single sign-on (e.g. Shibboleth). Newt plays nicely in a container environment, running as a system service or invoked from the command line.

## Motivation

My belief is that many web services used by archives, libraries and museums can benefit from a simplified and consistent back end. If the back end is "easy" then the limited developer resources can be focused on the front end which is what us humans experience day to day.

I've written many web applications over the years. Newt is focused on providing very specific glue leveraging existing microservices already used by libraries, archives and museums.  For many of these apps the core of an application is a JSON service (e.g. Invenio-RDM, ArchivesSpace). Newt can be used to extend these if needed. Let's take advantage of that. When we do need a custom application let also take advantage of a similar microservices approach. Build your core application in SQL with PostgREST+Postgres, hand of rendering to Pandoc running as a service. Newt can route your data between them two giving you similar benefits to complicated systems like Invenio but simple enough to be implemented by a single person.

## Newt stack, front to back

- A front end web server (e.g. Apache 2, NginX) can provide access control where appropriate (e.g. single sign-on via Shibboleth)
- Newt provides static file services but more importantly serves as a data router. It can validate and map a request to a JSON source, take those results then send them through Pandoc for transformation.

- JSON data source(s) provide the actual metadata management
  - Postgres+PostgREST is an example of a JSON source integrated with a SQL server
  - Solr, Elasticsearch or Opensearch can also function as a JSON source oriented towards search
  - ArchivesSpace, Invenio RDM are examples of systems that can function as a JSON sources
  - CrossRef, DataCite, ORCID are examples of services that function as JSON sources
- Pandoc server provides a templating engine to transform data sources

All these can be treated as "off the shelf". I.e. we're not writing them from scratch but we're accessing them via configuration. Even using PostgREST with Postgres the "source" boils down to SQL used to define the data models hosted by the SQL service.  Your application is implemented using SQL and configured with YAML and Pandoc templates.

## Taking advantage of JSON and S3 data sources

Newt was inspired by my working with PostgREST, Postgres and Pandoc. I also work allot of S3 object stores. I want Newt to be light weight. I wanted Newt to avoid writing anything to disk. That's possible now working with JSON API as data sources. I am in the planning stages of adding S3 protocol support to allow Newt applications to support a bigger domain space. Current plans are focused on using Minio as an off the shelf microservices to fill that responsibility.

## Getting Newt, Pandoc, PostgREST, Postgres and Minio

Newt is an experimental prototype (June/July 2023, and January/February 2024). It is only distributed in source code form.  You need a working Go language environment, git, make and Pandoc to compile Newt from source code. See [INSTALL.md](INSTALL.md) for details. Go is available from <https://golang.org>.

Pandoc is available from <https://pandoc.org>

PostgREST is available from <https://postgrest.org>.

Both Pandoc and PostgREST are written in Haskell, if you're going to compile them from source I recommend using GHCup <https://www.haskell.org/ghcup/>.

Postgres is available from <https://postgres.org>.

Minio is available from <https://github.com/minio/minio> and the Minio website at <https://min.io>.

## About the Newt source repository

Newt is a project of Caltech Library's Digital Library Development group. It is hosted on GitHub at <https://github.com/caltechlibrary/newt>. If you have questions, problems or concerns regarding Newt you can use GitHub issue tracker to communicate with the development team. It is located at <https://github.com/caltechlibrary/newt/issues>.

### "Someday, maybe" ideas to explore

- Integrate S3 object store support as a data source
- Support other rendering engines besides Pandoc server

## Documentation

- [INSTALL](INSTALL.md)
- [user manual](user-manual.md)
- [About](about.md)
