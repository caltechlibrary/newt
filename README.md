
# Newt

Newt ("NEW Take"), is an experimental [microservice](https://en.wikipedia.org/wiki/Microservices) working with other "off the shelf" microservices. The goal of the project is to provide a rapid devleopment platform for Library, Archive and Museum applications.

The Newt command runs in two modes. It's primarily a data router microservice designed to site behind your favorite web server (e.g. Apache 2, NginX) and route requests from the browser to a JSON API and optionally take the result and run it through Pandoc to generate web content.  Newt can also be used to generate SQL to bootstrap a JSON API built with PostgREST and Postgres.  The specifics of your application are defined by your JSON API microservice and the templates you render the JSON result with. Newt can also servce static web content as this is convient while you're developing a Newt based application (e.g. also make it easy to keep a project in a single Git repository).

A typical Newt app will include the following

- YAML file defining models and routes (i.e. the path part of a URL and http method)
    - Routes definitions include a PATH and http method, content type and target JSON API or S3 Object store end points
    - The models are used to generate SQL suitable for PostgREST+Postgres JSON API microservice
- SQL used by Postgres to manage data model instances of your application JSON API
- A set of Pandoc templates used to transform PostgREST JSON API results another text format (e.g. HTML, SVG, CSV, RSS, plain text)
- A directory of static content (e.g. any browser side JavaScript, CSS, image assets, HTML)

As a developer you spend you time by defining your application with the Newt configuration file and Newt generating bootstrap SQL code.
Data is managed by Postgres via the SQL code and accessed via PostgREST providing the JSON API. You likely have some Pandoc templates to display JSON results from the JSON API. You'll have some static web content you also want to integrated. 

Actually data processing winds up defined via SQL and runs in the Postgres database itself. Newt, PostgREST+Postgres and Pandoc replace all the middleware micro service code you used to write. They already are written and you can just use them as is. Limitting the time you spend writing micro services and the "back end" means you have more time to develop the human facing side using your "front end" developer skills.

A production Newt application is typically composed the following way

1. Front facing web server (e.g. Apache 2, NginX) providing single sign-on integration if needed and also static file hosting
2. Newt providing data routing between Front facing web server and a JSON API data source or S3 Object store and Pandoc server 
3. One or more JSON API data source (e.g. PostgREST+PostgREST, Solr, Elastic/Opensearch, Minio as S3 object store)
4. Pandoc running in server model providing a template engine to transform content from one format to another

Why YAML and why define the back end via SQL?   YAML is a mostly human friendly data representation. It is suited to a declaritive model of configuration. For Newt it provides a declaritive way to express simple data models and routing information.  SQL is a declaritive data management language. It is turing complete. Postgres database provides means of implementing functions, procedures, views, triggers and events. Taken together the encapsulates the whole data management life cycle. An advantage to using SQL is you declare the results you want (e.g. via SELECT statement) and database manager takes care of optimizing the result.  If you need to scale up, you cluster Postgres, run parallel instances of PostgREST, Newt and Pandoc server.  You don't need rewrite a bunch of middleware code it already scales with these off the shelf programs.

Note the Newt stack also provides a clean division of responsibilities.

1. Front Facing web server 
    - provides single sign-on integration (e.g. Shibboleth)
    - static file hosting
    - URL rewriting and redirection if needed
2. Newt provides an data flow or orchestration layer mapping requests to JSON API and assembling results via Pandoc templates if required
3. PostgREST+Postgres provide the JSON API data source and data management
    - receives request vetted by Newt
    - applies database access rules
    - assembles a response as a JSON document
    - database enforce ACID requirements and data consistancy requirements, it manages the application data
    - entire data processing programs can be run inside Postgres, including via events and triggers
4. Pandoc provides JSON or S3 object transformations
    - Pandoc never writes to disc
    - All actions are atomic, receives a request returns a result
5. S3 Object store provides access to data that does fit well directly in Postgres (e.g. a uploaded image or video file)


Postgres is really good at managing data.  This includes the full life cycle of create, reading, update, delete, listing as well as changes triggered an action or happening on a schedule (events). SQL provides a single langauge to implement all these things that are specific to the type and shape of data you're modeling.  In addition if you feel you need a different language (e.g. Python, Lua) this can be integrate directly into the Postgres database environment. Postgres also provides for full management of content access with a robust and flexible permission systems. But Postgres doesn't have a built in web friendly API. 

PostgREST turns Postgres into a JSON API service based on what you've built inside Postgres. The combination of PostgRESt+Postgres gives you a robust micro service for managing content. The combo works well with proven web servers like Apache 2 and NginX. Both PostgREST and Postgres are "cloud friendly" these days and they can scale "wide" with a high degree of parallelization expected at the largest scale of web applications.

You can build rich application directly on PostgREST+Postgres providing the JSON API and some static HTML, CSS and JavaScript. Many internet companies take a similar approach. The dump a whole bunch of JavaScript down to your web browser then make the web browser render the web page.  When that get too slow they then pre-rendering the whole message and cache the result. This is large part (along with survalence advantages) is what gave rise to single page applications and their propogation on phones and tablets.  I can't count the number of times I've experienced a phone "app" (really a thinly valed website application) take forever to start, update and become responsive. Single page applications also rarely lend themselve to re-use or easy integration with other systems and services.

The ticket to solve the limitation of JSON for human consumption while keeping it for development and data management is to use a template engine. Pandoc is most excellent at converting one structured document to another type. This includes processing JSON. It even includes a minimalist template language.  Pandoc can run as a robust secure microservice with configuration. What is needed is an easy way to map a JSON API response data source, optionally get a template and then send it Pandoc server for a result to pass back to the web browser. That problem is why I wrote Newt.

## Origin story

Newt came about when I realized that all I needed was a data router that could map a web browser request to the JSON API provided by PostgREST+Postgres and Pandoc running as a service. This happens to generalizes to nicely to a data source mapping (e.g. JSON API or S3 Object store) and to transformation engine (e.g. Pandoc running in server mode). All that was missing was something to wrote the requests from the browser to the data source through the transformation engine before returning them to the browser.

I demonstrated Newt concept to my colleagues my prototype included supoprt for JSON API via PostgREST+Postgres and Pandoc running in server mode. I didn't get huge rounds of applause but did get some encouragement to explore it further. Similarly when I demonstrated Newt at a the SocCal Code4Lib group and people more interested in finding out about PostgREST and Postgres. So today I think of Newt as a demonstration of using PostgREST+Postgres and Pandoc running in server mode. Those are the pieces doing the heavy lifting. Newt is just there to easy knitt those services together in a simple declaritive way. I also think it is hard to step back from Perl, PHP, Python, Ruby and Go when you spent so much effort learning the languages and dev environment so when someone says you can just use "off the shelf", some SQL, maybe a light wweight markup it sort is an anti-climax.

## Systems design choices

When people have gotten what I'm proposing there are very legimate critisms that I think need consideration

1. Why SQL and not a "better" lanauge?
2. Pandoc as template engine, Why?
3. How do I handle file uploads?

The first I think is because SQL is unavoidable even in the current reliance on ORM.  For Libraries, Archives and Museums managing metadata about our collections is really really critical.  When we rely on ORM to abstract how we're storing the information we ducking responsibiltiy for something that is critical to our organizations.  Modern SQL databases like PostgreSQL provide good support for JSON objects which is effective in representing some of our heirarchly metadata relationships. We acrue massive ammounts of metadata for our organizational and staff sizes. It's in the nature of our discpline since we acrue and curate content over time.  As a result I think getting comforatable with SQL and PostgreSQL in particular brings huge advantages.  Postgres also supports data replacation and working on copies of our metadata corpus is liberating as we can operate on the data using data science techniques to enhance and evaluate the metadata we manage.

The second question I run into is why Pandoc?  Each of the programming languages I used over the years have one or more template languages to make rendering database results easily as text or structured text (e.g. HTML). None of these langauge are compatible. On the other hand in the date science community and for those who work with static site generation Pandoc and by extension Pandoc templates have become common place. The populatiry of R-Studio and the result rMarkdown as well as GitHub flavored Markdown have also increased adoption of Pandoc. Since Pandoc can run as a web service (nice if you're processing allot of document conversations) why not use it. I look at Pandoc as sort of a neutral useful player in document transformation. It just works and runs well on Windows, macOS and Linux/POSIX systems.

Newt was inspired by using PostgREST with Postgres but also by the simplicity of running Pandoc as a service. I wanted a minimal data router that would read it's configuration at startup and just run until it was shutdown. Newt doesn't maintain state once it is up and running. This means like Pandoc server or PostgREST it can scale wide but just adding another instance.

### Arriving at Newt's minimal feature set

When I encountered peoples reaction to building a web service with SQL I was surprised but in hind sight I shouldn't have been. I've learned SQL as I needed spread out over a long period of time. I first learned the basic SELECT statement the original MySQL implemented. Then added INSERT, UPDATE and DELETE. For a long time that was all that I needed. One thing I've noticed about those statements is they are not hard to learn but it is helpful to look at an existing example and modify it to your needs. That when I realized that I should provide some generate SQL for basic CRUD and list operations by adding model description to Newt's YAML defining the routes. When you start a Newt based project Newt can take your model and generate the SCHEME statements and some basic functions, views suitable to setup PostgREST and Postgres. All you need is a running Postgres with a development (aka admin) account with appropriate permissions. This bootstrap is intended to help you with the scaffolding needed for a basic metadata management application such we is the bread and butter of Libraries, Archives and Museum developers.

Finally I believe decoupling the front and back end allows for a more flexible approach to development. Typically data models are slow changing and the back end should also be changing at a slower rate. On the other hand front end web design changes frequently as new and simpler ways of presenting content and discovered. The front end likewise benefits through a predictable evolution of the back end (e.g. new routes become available, old ones still work). Newt works with those observations. You define the back end, get it working and then focus on front end. 

Today's web browsers, 2024, easily deliver on the HTML5 promises of 2008 but many of us don't have the time to focus on the front end becasue we so busy managing complex back end systems. We need this to change, today's web browser is a vary capable platform[1]

[1]: [You don't need JavaScript for that](https://www.htmhell.dev/adventcalendar/2023/2/) is a good reminder of what the browser has built in today.

Newt provides a clear division of labor

1. Web Server is responsible for access control and interacting with Newt
2. Newt provides data routing to Postgres+PostgREST and optional Pandoc template engine
3. Postgres+PostgREST provide a JSON API
4. Postgres is responsible for manaing data and it's models via SQL, views, functions, procedures, events and triggers
5. The web browser talks to the web server and render the application's html friendly interface

The person developing with Newt writes SQL to define the backend, may write Pandoc templates if that is desired and builds the front end with
standard static web assets (e.g. HTML pages, CSS, JavaScript). Newt can support traditional websites and single page applications. It just saves writing a whole bunch of services that already exist.

## Newt as a stack

- Newt is a data router, form validator, static file microservice
- [Pandoc](https://pandoc.org) server provides a template engine to render JSON as something else
- JSON data sources
    - [Postgres](https://postgres.org) + [PostgREST](https://postgrest.org), SQL to JSON API
    - Solr, Elasticsearch/Opensearch full text search engines
- Apache 2 or NginX for integrating SSO or controlling web access (e.g. via Shibboleth, OAuth2)

The Newt prototype type runs on localhost and assumes all the JSON data sources also run on localhost.  It
also assumes that S3 object stores run locally (e.g. Minio). 

In development you don't necessarily needs to be running Apache or NginX as a web server. You could point
your browser at Newt's own localhost service directly.  In a production setting you should put Newt behind
a robust web server like Apache 2, NginX and use that web server to proxy to Newt running on localhost. You
want to does inorder to provide SSL support (which Newt doesn't supply) as well access crontrol if needed.

The stack is deliberately built "off the shell". The person developing software is primarily editing SQL, creating and
updating Pandoc templates (if used) and providing HTML, CSS and JavaScript the constitutions the human user
interface.

## Orchestrating your app with Newt

Newt is configured with a YAML file. Currently the configuration file uses three main attributes

htdocs
: The path to the htdocs directory holding any static assets (e.g. CSS, JavaScript, HTML files, image assets)

env
: A list of environment variables available to Newt's routes and models

routes
: An object describing the mapping of an HTTP request to JSON data source and and optional Pandoc server processing

Under routes you will also use the following attributes

var
: (optional) A list of variable names and types used in validating a request path or web from submission

req_path
: A expression describing a URL path received by Newt (typical made by a web browser or proxy for the web browser)

req_method
: An HTTP method (e.g. GET, POST, PUT, PATCH, DELETE) related to the req_path being handled

api_url
: The URL expression used by Newt to contact the JSON data source for the route described by req_path and req_method

api_method
: The HTTP method (e.g. GET, POST, PUT, PATCH, DELETE) of the JSON data source associated api_url for the given route

api_content_type
: The HTTP content type expression used when submitting the request to the JSON data source

s3_url
: The URL to connect to the S3 object store including the base path

s3_method
: The allowed action(s), create, read, update, delete object

s3_content_type
: The s3 storage content type to use (e.g. application/json, image/jpg)

pandoc_template
: (optional) The path to the pandoc template used to process the results of the JSON data source request results

Additional attributes maybe added to the YAML file in the future.

## Data flow steps

1. Web browser => (Web Server proxy) => Newt
2. Newt => data source (e.g. Postgres + PostgREST or S3 Object store)
3. Newt => Pandoc (optional step)
4. Newt => (Web Server proxy) => Web browser

The Pandoc processing is optional. In this way you could expose selected elements of your JSON data source to the front-end web service by mapping the api request but skipping the Pandoc template used to transform the data. You can also use Pandoc templates to reshape your data into other formats or restructure a complex JSON object into a simplified one. Pandoc template engine is very flexible.

## Other JSON data sources

Newt is intended to work via local host JSON data sources. If you need to work with external ones you need to proxy to them. This can be done many ways and also configured as part of your main web server (e.g. defining a local host reverse proxy via virtual hosting support in the web server configuration). 

If the JSON data source is available via localhost then Newt can talk to it. Beside Postgres+PostgREST a common JSON data source would be Solr or ElasticSearch/OpenSearch. Newt doesn't know specifically about PostgREST, it just talks to the JSON data source you point at for a given route. This allows you to build up Newt applications from various existing JSON API that might be available on your system.

### Handling errors

Newt vets the initial request before contacting the JSON data source. If the request has a problem it will return an appropriate HTTP status code and message.  If the request to the JSON data source has a problem, it will pass through the HTTP status code and message provided by the JSON data source.  Likewise if Pandoc server has a problem Newt will forward that HTTP status code and message. If either the JSON data source or Pandoc server is unavailable Newt will return a "Gateway" http status code and message.

### Static file support

Newt first checks if a request is matched in one of the defined routes. If not it'll try to service the content from the "htdocs" location. If the file is not found then a 404 is returned.  If the "htdocs" value is not set then a 404 will get returned if no routes match. 

Note Newt's static file services are very basic. You can't configure mime type responses or modify behavior via "htaccess" files. If Newt is running behind a traditional web server like Apache 2 or NginX then you could use that service to host your static content providing additional flexibilty.

### Handling secrets, scaling and limitations

Newt's YAML file does not explicitly contain any secrets. This was intentional.  You may need to pass sensitive data to your JSON data source for access (e.g. credentials like a username and password). This should be provided via the environment and the YAML file needs to include these environment variable names in the "env" attribute.  The environment variables can be used to contact the JSON data source. There is still a risk in that theoretically that data source could echo return sensitive information. Newt can't prevent that. Newt is naive in its knowledge of the JSON data source content it receives or hands ofdf the JSON to Pandoc.  You're responsible for securing sensitive information at the database JSON data source level. Follow the recommendations in the Postgres community around securing Postgres.

While Newt was conceived as a  small scale web application platform for Libraries, Archives and Museums it is capable of scaling big as long as your JSON data source can scale big.  Using the Newt stack elements can all be run easily behind load balancers and in parallel across machines because they require no synchronized shared of data between them or instances of them. Postgres itself can be configured in a HA cluster to support high availability and high demand.

Presently Newt does not support file uploads. If you need that you'll need to write your own service to handle it. 

Newt runs exclusively as a localhost service. In a production setting you'd run Newt behind a traditional web server like Apache 2 or NginX. The front-end web service can provide access control via basic auth or single sign-on (e.g. Shibboleth). Newt plays nicely in a container environment, running as a system service or invoked from the command line.

## Taking advantage of JSON data sources

Newt prototype works with a JSON source that can be accessed via URL and HTTP methods (e.g. GET, POST, PUT, PATCH and DELETE) using http protocol. Many systems have JSON data API today (2023). This includes existing Library, Archive and Museum applications like Invenio RDM and ArchivesSpace as well as search engines like Solr and Elasticsearch.  This means Newt can be used to extend existing systems that provide a localhost JSON data source. If the data source is not available via localhost then a proxy is required.

## Motivation

My belief is that many web services used by Archives, Libraries and Museums can benefit from a simplified and consistent back-end. If the back-end is "easy" then the limited developer resources can be focused on the front-end. **An improved front-end offers opportunities to provide a more humane experience for staff and patrons.**. 

I've written many web applications over the years. Newt is focused on providing very specific glue to other microservices which already perform the core of an application (e.g. Postgres+PostgREST provide a data management engine and Pandoc server provides a template engine).  Newt works around a declarative data model for configuration and generating basic SQL scripts for use with Postgres and PostgREST. That means the first step in your application is defining things. The rest of the time your focus is on UI via Pandoc templates and browser side code in the form of JavaScript and CSS.  Newt functions primarily as a data processing pipeline orchestrating requests through Postgres+PostgREST and optionally through processing the result via Pandoc. Newt's configuration langauge is YAML. YAML was picked because it is widely use. Newt implements "domain specific language" or DSL on top of YAML to support rendering a minimally useful SQL source file for use by Postgres and PostgREST providing the data manage layer of your application. You evolve your application by enhancing the SQL codebase Postgres will run, creating or modifying Pandoc templates and augumenting it with browser side behaviors via HTML, CSS and JavaScript. Stick with behind a high performance web server like NginX or Apache 2 and you can layer in your single sign-on system (e.g. Shibboelth, OAuth) or use the web service's native access control features (e.g. basic auth).


## Newt stack, front-end to back-end

- A front end web server (e.g. Apache 2, NginX) can provide access control where appropriate (e.g. single sign-on via Shibboleth)
- Newt provides static file services but more importantly serves as a data router. It can validate and map a request to a JSON source, take those results then send them through Pandoc for transformation.

- JSON data source(s) provide the actual metadata management
    - Postgres+PostgREST is an example of a JSON source integrated with a SQL server
    - Solr, Elasticsearch or Opensearch can also function as a JSON source oriented towards search
    - ArchivesSpace, Invenio RDM are examples of systems that can function as a JSON sources
    - CrossRef, DataCite, ORCID are examples of services that function as JSON sources
- Pandoc server provides a templating engine to transform data sources

All these can be treated as "off the shelf". I.e. we're not writing them from scratch but we're accessing them via configuration. Even using PostgREST with Postgres the "source" boils down to SQL used to define the data models hosted by the SQL service.  Your application is implemented using SQL and configured with YAML and Pandoc templates.

## Getting Newt, Pandoc, PostgREST and Postgres

Newt is an experimental prototype (June/July 20230). It is only distributed in source code form.  You need a working Go language environment, git, make and Pandoc to compile Newt from source code. See [INSTALL.md](INSTALL.md) for details.

Pandoc is available from <https://pandoc.org>, Postgres is available from <https://postgres.org> and PostgREST is available from <https://postgrest.org>.  If you want to compile the latest Pandoc or PostgREST (both are written in Haskell), I recommend using GHCup <https://www.haskell.org/ghcup/>

## Newt source repository

Newt is a project of Caltech Library's Digital Library Development group. It is hosted on GitHub at <https://github.com/caltechlibrary/newt>. If you have questions, problems or concerns regarding Newt you can use GitHub issue tracker to communicate with the development team. It is located at <https://github.com/caltechlibrary/newt/issues>.

### "Someday" ideas to explore

- consider JSON data sources requiring https protocol support may be added (e.g. JSON data sources like CrossRef, DataCite, ORCID use https protocol)
- object storage (e.g. file upload) via S3 protocol (e.g. via minio)
- implement Newt in Haskell and directly integrate Pandoc library
- integrate Lua filters and Pandoc partial support as part of the Newt pipeline
- a Newt based IDE (as demonstration app) for building Newt based projects

## Documentation

- [INSTALL](INSTALL.md)
- [user manual](user-manual.md)
- [About](about.md)

