
# Newt

Newt project is a collection of tools providing a rapid web application development platform for libraries, archives and museums. The collection of commands focuses on using "off the shelf"  micro services combined with a web application composition via a YAML file. Newt provides the following commands.

- [newt](newt.1.md) is an experimental [micro service](https://en.wikipedia.org/wiki/microservices) designed for working with other "off the shelf" micro services. The primary purpose is to function as a localhost data router between a data source and rendering engine. It can also function as a static file server when development an application
- [newtmustache](newmustache.1.md) is a extremely light weigth [Mustache](https://mustache.github.io/) template rendering engine inspired by Pandoc server
- [newtpg](newtpg.1.md) is a command line program designed to generate SQL used bootstrap a JSON API built with PostgREST and Postgres.  Combined `newt` and `newtpg`  saves you from building yet another middleware micro service
- **newtform** (planned but not prototyped) is a web form generator mathcing the SQL generated for the JSON API

## Why rapid development tools? Why Newt?

The promise of modern [RAD](https://en.wikipedia.org/wiki/Rapid_application_development "Rapid Application Development") tools is a faster development cycle delivering a minimal viable application sooner.  The RAD system accomplishes this by code generation taking a splecification (a.k. configuration file) and rendering the code described. The system eases development by providing either a simple highlevel configuration language or a GUI that generated the application code. Often a RAD system will include GUI building tools as well along with a means of hooking up the GUI elements to specific actions. In the past these systems tended to be propriary and often required developer specialization to use it.

Newt is a lightweight set of RAD tools. They all use the same configuration langauge, [YAML](https://en.wikipedia.org/wiki/YAML). Why YAML? Must developers and even some none developers already know it.  The Newt specific aspects of YAML are the data structures represented as YAML used by Newt's commands. Newt doesn't have a GUI and doesn't provide layout tools. The reason if there is lots of software that already does that for web development. Instead Newt focuses on composing the conversation between "office the shelf" micro services and static file assets that result in a "web application". As a result Newt is very narrowly focused. It provides a simple data router that can sit behind a traditional front end web server. Typically this is a JSON data source such as that provided by Postgres+PostgREST. It includes a lightweight Mustache template engine that accepts a JSON structure and template rendering the result. It include a SQL code generator for creating JSON API using Postgres and PostgREST. It will include a simple tool for processing the Newt YAML into web forms suitable for including along with your static assets of your application. Finally Newt is focus in the application domain of metadata curation for libraries, archives and museums.

If successful then Newt will assist technically inclined library staff or librarians in standing simple web applications for metadata curation.

## Where is development time spent?

With Newt your development time is focused on three areas.

1. Modeling your data in YAML and using that to generate SQL and HTML 5 web forms
2. Rendering content coming from JSON API using simple templates engines (e.g. Pandoc or Mustache templates)
3. Enhancing the user expereience browser side using HTML5, CSS and if needed JavaScript (a.k.a. traditional front end development)

Newt achieves this division of responsibities. Newt's YAML file is responsible for describing the data models used by your application and for defining the specific routes used to access your data pipeline. Your data pipeline is a small set of micro services. In Newt's protype I've used Postgres and PostgREST to provide a JSON data source and Pandoc running in server mode for a template rendering engine. The generative tools that come with Newt provide the means to render the SQL needed to setup Postgres and PostgREST as well as generating HTML/Markdown for your web forms interacting with the JSON API data source.

## Example steps to demonstrate building with Newt

The following example describes using PostgREST+Postgres to provide a JSON API, Pandoc server as your render engine and Newt to serve static content and route requests to the JSON API and render engine.

1. Create a directory/folder for your project and change into that directory (setup version control if desired)
2. Create a Newt configuration file describing the models and routes your application will need
3. Use `newtpg` and the Newt configuration file to generate SQL to bootstrap your PostgREST+Postgres JSON service
4. Use `newtform` to render HTML/Markdown of your applications' web forms

At this point you should have a skeleton of your application working. The following are iterative steps you
use to refine your application.

1. Add/update your initial SQL data models and test inside Postgres (psql works nicely for this)
2. Add/update the routes in the routes section of the YAML configuration file if needed
3. Add/update any Pandoc templates, HTML, CSS and JavaScript as needed for your application
4. (re)start PostgREST and `newt` then test with your web browser

These steps repeat and are you spend most of your application development time.  Notice that each of the steps
are focused on a specific area of your application. Step one, improving your data modeling and data management. Step two refining the selection of URLs your application responds to. Step three Is were we start front end development by taking the JSON data and rendering it into HTML. It is here where you can also augment the HTML produced by the template engine to do more browser side. The least step refreshes your development view of your application by reloading the JSON API and data router. This is easily done via a shell script.

## What about security, single sign-on, integration with websites or services?

The `newt` command is a simple micro service providing data routing based on its configuration at startup. It's a team player.  At many universities, colleges, research libraries, archives and museums there is an existing single sign-on mechanism like Shibboleth running along with Apache 2 or NginX web servers.  Newt would run behind those services via reverse proxy. Newt itself doesn't know about users, it knows about routing requests. Newt after reading the configuration file doesn't maintain state.  While the log output may contain identifiable information (e.g. IP address of request) or the JSON data source could contain sensitive information Newt doesn't retaining it. It just routes the data and gets out of the way.

A typical `newt` production setup might look like this

1. NginX with Shibboleth controls access to web site resources and where appropriate proxies to `newt`
2. `newt` responds to requests and maps those to a data source (e.g. PostgREST+Postgres JSON API) and gets back a response
3. `newt` can take a data source response and send it to a render engine (e.g. Pandoc in server mode or `newtmustache`)
4. `newt' assembled result is handed back to the NginX web server to pass onto the requesting web browser

In the example securing your application can happen both at the NginX level (e.g. requiring single sign-on) and
at the JSON API level via Postgres's management of PostgREST responses. NginX can also be used to proxy external resources you may wish Newt to route to.

Newt tries to do as little as possible while still providing data routing and static content services. This reduces the attack surface. Newts log output is written to standard out and standard error and not directly to disk to avoid the problem of high traffic logging filling up your disk drive. Logging can easily be captured by your servers logging system (e.g. systemd and it's log handler).  Newt is configured through the environment so does not require storing secrets in the configuration file.  Newt only reads the configuration at startup and can not write it back to disk. In fact Newt can't write to disk at all.  Newt focuses on data routing at the web application level only. It has a limited two stage pipeline for requests processing. For any given route defined in Newt's configuration it can contact a data source (e.g. PostgREST+Postgres JSON API) and optionally send the result through a rendering service (i.e. Pandoc running in server mode) before handing that result back to the requestor. Newt only knows how to speak http and can only communicate via localhost so can't access outside systems with you proxying to them. Newt only listens for routes defined at startup or derived from a designated htdocs directory.

Keeping Newt simple minimizes the attack surface and keeps Newt a team player in your micro service based applications.

### What about "scaling"?

Newt is just a data router. You can safely run many Newt instances in parallel as needed. They can run them on the same machine or separate machines. The instances don't share data or coordinate. They just read their configuration files when they start up and route data accordingly.

A typical Newt stack of PostgREST+Postgres and Pandoc server also can scale up. You can run as many instances of Pandoc server or PostgREST as you need. You can spread them across machines. They are both stateless systems like Newt. A Postgres database provides consistency and can be configured in a high availability cluster. Postgres scales.

When I created Newt I was interested in small scale applications. I created it as a simple stateless micro service. Because it is a simple stateless micro service can scales as wide as you like. It scales in the same way as Pandoc server or PostgREST.

### Annatomy of a Newt based web application

Newt application development is friendly to version control systems (e.g. Git). It consists of a Newt configuration file, SQL files, HTML templates and any static web assets you need. A typical disk layout of a Newt project could look like this-

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

The important take away was I had failed to appreciate how Newt successfully shifted the discussion from programming language frameworks, package management and build systems to to modeling data in SQL and using simpler HTML5, CSS and JavaScript for display.

> Newt, a type of salamander, doesn't seek attention. It does its own thing. You only notice a salamander if you look carefully.

## System design choices

Demonstrating the "Newt stack" has resulted in questions. Here's the big four questions.

1. Why SQL? Why build your data models with SQL?
2. Pandoc as rendering engine, why?
3. Why YAML for configuration?
4. How do I handle file uploads?

My answer to the first question: I think knowing at least some SQL is unavoidable as a web developer. While many have adopted an ORM[1] to generate the SQL to render models and manage data the resulting SQL is often far from ideal. This is inspite of the fact that the ORM concept is decades old. The problem I worry about with ORM (yes, I've use them) isn't inefficiency. The real problem is the ORM obscures the data model and that discourages data re-use. The vast majority of web applications run in institutional settings use SQL databases (e.g. MySQL, Postgres, Oracle, SQLite3). Embracing SQL gives you power to extend those systems and even integrte them.  In 2024 SQL remains the common language to manage data in a database. That's true even have a half century of existence. Even for non-relational data this has become true[2]. SQL may look ugly or quirky but it definately has legs and plans to stick around for a very long time. Let's embrace it!

You may protest, "but some people really hate writing SQL!" Newt does provide a tool to address this. Newt provides `newtpg` that can take a data model described in the Newt configuration file and render SQL suitable for standanding up and managing data in Postgres and PostgREST. Another tool, `newtform` planed for development, will render the same model to HTML 5 and JavaScript so you can keep your model paired between the JSON API that manages the data and your brower's web form for doing data entry. in YAML that can be used to generate Postgres SQL as well as simple HTML 5 web forms.  Of course Newt just provides a tool, you don't have to use it and you're welcome to write SQL and your web forms from scratch.

The second question, why run Pandoc server to render things? Pandoc is good at transforming structured text. JSON is structured text and Pandoc can easily turn that into HTML.  It is also ready to go in the spirit of using "off the shelf" micro services. In the data science and library circles I travel I've seen a huge adoption of Pandoc for static site generation. That lead me to pick PAndoc as a rendering engine for my Newt[4] prototype. If another template language and engine comes on the scene, Newt can be adapted to us it instead[5]. If you really don't like Pandoc templates Newt now provides a similar engine that supports Mustache templates, `newtmustache`.

The third question, "Why YAML for configuration?". Personally I'm enthusiastic about YAML. It is sort of boring. But I think that is is strength. YAML has become umbiquious for describing configuration and simple data structures. It easily converts to JSON. It's declarative. Newt uses YAML largely because it's well known and is specifically known by my colleagues at Caltech Library. Why invent a new language when I can use one that is already known?

Final question is "How do I handle file uploads?". The short answer is Newt doesn't. Eventually I plan to support for S3 protocol storage systems (e.g. Minio, S3) but I haven't had time to implement this and do not need it for the application I am building with Newt.  A longer answer is yes, it is possible but you need to know Postgres really well.  Technically Postgres+PostgreSQL can handle file uploads because you can store files or large objects in Postgres. Personally I don't want to store files in my data base management system. I'd rather store them in an object store like S3. I don't recomment Newt for applications that require handling file uploads unless you want to write your own micro service to implement it.


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
- `newtform` can use the Newt configuration file to render simple HTML/Markdown forms

Here's the data flow steps of `newt` data router.

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
: This is a list of data models used by Newt to generate bootstrap SQL code for PostgREST+Postgres. The markup will be models is based on GitHub YAML issues template syntax.

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

template
: (optional) The path to the pandoc or mustache template used to process the results of the JSON data source request results

render_port
: (optional) Set the port to use for contacting the render engine. If not set it assume 3030 is the port for the render engine.

The **models** attribute holds a list of models expressed in Newt's data model DSL. The original protype DSL is going to be replaced with the YAML described in [Syntax for Github's form schema](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-issue-forms). That syntax will allow the creation of `newtform` to generate web forms as well as be used by `newtpg` to generate Postgres compatible SQL models and behaviors. Here is an example using the old syntax as a placeholder.

### Handling errors

Newt command vets the initial request before contacting the JSON data source. If the request has a problem it will return an appropriate HTTP status code and message.  If the request to the JSON data source has a problem, it will pass through the HTTP status code and message provided by the JSON data source.  Likewise if Pandoc server has a problem Newt will forward that HTTP status code and message. If either the JSON data source or Pandoc server is unavailable Newt will return a "Gateway" http status code and message.

### Static file support

Newt command first checks if a request is matched in one of the defined routes. If not it'll try to service the content from the "htdocs" location if that is defined in the configuration. If the file is not found or an htdocs directory has not been specified a http status of 404 is returned.

Note Newt's static file services are very basic. You can't configure mime type responses or modify behavior via "htaccess" files. If Newt is running behind a traditional web server like Apache 2 or NginX then you could use that service to host your static content providing additional flexibilty.

### Handling secrets, scaling and limitations

Newt's YAML file does not explicitly contain any secrets. This was intentional.  You may need to pass sensitive data to your JSON data source for access (e.g. credentials like a username and password). This should be provided via the environment and the YAML file needs to include these environment variable names in the "env" attribute.  The environment variables can be used to construct the URLs to contact the JSON or S3 data sources. There is still a risk in that theoretically that data source could echo return sensitive information. Newt can't prevent that. Newt is naive in its knowledge of the data source content it receives or hands off to Pandoc.  You're responsible for securing sensitive information at the database or s3 data source level. Follow the recommendations in the Postgres community around securing Postgres.

While Newt was conceived to be used on as a small scale web application platform for libraries, archives and museums it is capable of scaling big as long as your data source(s) can scale big.  Using the initial "Newt stack" elements can all be run easily behind load balancers and in parallel across machines. Newt is transactional. It does not require synchronized or shared of data between instances. Similarly PostgREST and Pandoc services are transactional and do not require shared state to function in parallel. Postgres itself can be configured in a HA cluster to support high availability and high demand. It should be possible to scale a Newt based application as large as those systems can be scaled.

Presently Newt does not supports file uploads. The plan is to integrated support for an S3 object store. That support is still very much in the planning stages.

Newt command runs exclusively as a localhost service. In a production setting you'd run Newt behind a traditional web server like Apache 2 or NginX. The front end web service can provide access control via basic auth or single sign-on (e.g. Shibboleth). Newt plays nicely in a container environment, running as a system service or invoked from the command line.

Postgres is available from <https://postgres.org>.


## About the Newt source repository

Newt is a project of Caltech Library's Digital Library Development group. It is hosted on GitHub at <https://github.com/caltechlibrary/newt>. If you have questions, problems or concerns regarding Newt you can use GitHub issue tracker to communicate with the development team. It is located at <https://github.com/caltechlibrary/newt/issues>.

### "Someday, maybe" ideas to explore

- Integrate S3 object store support as a data source
- Support other rendering engines besides Pandoc server

## Documentation

- [INSTALL](INSTALL.md)
- [user manual](user-manual.md)
- [About](about.md)
