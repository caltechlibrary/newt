
# Newt

Newt is an experimental [microservice](https://en.wikipedia.org/wiki/Microservices) working with other "off the shelf" microservices. It is a proof of concept of a "new take" on the old web stack. The "new stack" is built on off the shelf microservices. The initial target was Postgres+PostgREST and Pandoc server. In practice the Postgres+PostgREST can be replaced (or operating along side) and JSON data source (i.e. JSON API whos interface uses http).

## NP^3, the newt stack

- N is the Newt data router, form validator, static file microservice
- P^3 is
    - [Pandoc](https://pandoc.org) server
    - [Postgres](https://postgres.org) + [PostgREST](https://postgrest.org) JSON API

## MP^3 plus

The targetted microservices all run on localhost.  Individually they have very narrow focuses. In a production setting you would normally run them behind a front-end web server like Apache 2 or NginX. The advantage of this approach is that the front-end web service can provide things like authentication or single sign-on integration.  It also allows a NP^3 based application to sit along many other web applications on the same virtual host.

## The way Newt works

Newt is configured with a YAML file. Currently the configuration file uses three main attributes

htdocs
: The path to the htdoc directory holding any static assets (e.g. CSS, JavaScript, HTML files, image assets)

env
: A list of environment varaibles available to Newt's routes

routes
: An object describing the mapping of request information, how to map it to a JSON data source and optional Pandoc server processing of the result

The routes attribute itself is a object with the following attributes

var
: (optional) A list of variable names and types used in validating a request path or web from submission

req_path
: A expression describing a URL path recieved by Newt (typical made by a web browser or proxy for the web browser)

req_method
: An HTTP method (e.g. GET, POST, PUT, PATCH, DELETE) related to the req_path being handled

api_url
: The URL expression used by Newt to contact the JSON data source for the route described by req_path and req_method

api_method
: The HTTP method (e.g. GET, POST, PUT, PATCH, DELETE) of the JSON data source associated api_url for the given route

api_content_type
: The HTTP content type expression used when submitting the request to the JSON data source

pandoc_template
: (optional) The path to the pandoc template used to process the results of the JSON data source request results

Additional attributes maybe added to the YAML file in the future.

## Newt data flow

1. Web browser => Newt
2. Newt => JSON data source => Newt
3. Newt => Pandoc => Newt
4. Newt => Web browser

The Pandoc processing is optional. In this way you could expose selected elements of your JSON data source to the front-end web service by mapping the access but skipping the Pandoc template used to transform the data.

### Handling errors

Newt vets the initial request before contacting the JSON data source. If the request has a problem it will return an appropriate HTTP status code and message.  If the request to the JSON data source has a problem, it will pass through the HTTP status code and message provided by the JSON datasource.  Likewise if Pandoc server has a problem Newt will forward that HTTP status code and message.

### static file support

Newt first checks if a request is matched in one of the route definitions. If that fail then it'll check against what is available in the htdocs directory if defined in the YAML file. It will return a 404 if no match is found.

### NP^3 Opportunities

The Newt stack aims to work with "off the shelf" web services. This includes external web services like CrossRef, DataCite or ORCID. These all provide JSON API that can be accessed over HTTP. 

As Newt's job is to route requests to a JSON data source and (optionally) process the results with Pandoc server it is not limitted to a specific application. A Newt "app" is simply a collection of routes. They migth all be related but do not have to be.

### Handling secrets, scaling and limitations

Newt's YAML file does not explicitly contain any secrets. This was intentional.  You may need to pass sensitive data to your JSON data source for access (e.g. credentials like a username and password). This should be provided via the environment and the YAML file needs to include these environment variable names in the "env" attribute.  The environment variables can be used to contact the JSON data source and theoretically that data source code echo them to Pandoc or the requesting browser. They is no easy way to prevent that so care should be taken to make sure your JSON data source does not echo those values or return so that Pandoc would process them in a template and return back to the web browser.

While Newt was concieved to services small scale web applications in Libraries, Archives and Museums it is capable of scaling big as long as your JSON data source can scale big.  Using the NP^3 stack elements Newt, Pandoc and PostgREST can all be run easily behind load balancers and in parallel across machines because they require no synchronizated shared of data. Postgres itself can be configured in a HA cluster to support high availability and high demand.

Presently Newt does not support file uploads. This may change in the future.That limits Newt's present application targets. If you needed that support today then you'd need to configure your front-end web server to proxy to a service that could provide file upload support.

Newt runs exclusively as a localhost service. In a production setting you'd run Newt behind a traditional web server like Apache 2 or NginX. The front-end web service can provide access control via basic auth or single sign-on (e.g. Shibboleth). It plays nicely in a container environment or running straight up as a system service.

## Taking advantage JSON sources

Newt works with a JSON source that can be accessed via URL and HTTP methods (e.g. GET, POST, PUT, PATCH and DELETE). Many systems have JSON data API today (2023). This includes existing Library, Archive and Museum applications like Invenio RDM and ArchivesSpace as well as search engines like Solr and Elasticsearch.  It also includes many library services like CrossRef, DataCite and ORCID. These can be integrated via Newt in a easily managed may.


## Motivation

My belief is that many web services used by Archives, Libraries and Museums can benefit from a simplified and consistent back-end. If the back-end is "easy" then the limited developer resources can be focused on the front-end. **An improved front-end offers opportunities to provide a more humane experience for staff and patrons.**. 

I've written many web applications over the years. Newt is focused on providing very specific glue to other microservices which already perform the core of an application (e.g. PostgREST presenting a data management engine and Pandoc server presenting a template engine).  Newt takes a declarative model to configuration. You define things.  From those definitions it delagates processing to JSON data sources and optionally to Pandoc server if JSON needs to be transformed into HTML. It does this by limitting it's responsibility to verifying routes declared in the YAML file and validating paths and form data requests.  Newt only runs on localhost so does not provide a means of virtual hosting or other things associated with a front-end web service.

## NP^3 stack

From front-end to back-end

- A front end web server (e.g. Apache 2, NginX) can provide access control where appropriate (e.g. single sign-on via Shibboleth)
- Newt provides static file services but more importantly serves as a data router. It can validate and map a request to a JSON source, take those results then send them through Pandoc for transformation.

- JSON data source(s) provide the actual metadata management
    - Postgres+PostgREST is an example of a JSON source integrated with a SQL server
    - Solr, Elasticsearch or Opensearch can also function as a JSON source oriented towards search
    - ArchivesSpace, Invenio RDM are examples of systems that can function as a JSON sources
    - CrossRef, DataCite, ORCID are examples of services that function as JSON sources
- Pandoc server provides a templating engine to transform data sources

All these can be treated as "off the shelf". I.e. we're not writing them from skratch but we're accessing them via configuration. Even using PostgREST with Postgres the "source" boils down to SQL used to define the data models hosted by the SQL service.  Your application is implemented using SQL and configured with YAML and Pandoc templates.

## Getting Newt, Pandoc, PostgREST and Postgres

Newt is an experimental prototype (June/July 20230). It is only distributed in source code form.  You need a working Go langauge environment, git, make and Pandoc to compile Newt from source code. See [INSTALL.md](INSTALL.md) for details.

Pandoc is available from <https://pandoc.org>, Postgres is available from <https://postgres.org> and PostgREST is available from <https://postgrest.org>.  If you want to compile the latest Pandoc or PostgREST (both are written in Haskell), I recommend using GHCUp <https://www.haskell.org/ghcup/>

## Newt source repository

Newt is a project of Caltech Library's Digital Library Development group. It is hosted on GitHub at <https://github.com/caltechlibrary/newt>. If you have questions, problems or concerns regarding Newt you can use GitHub issue tracker to communicate with the development team. It is located at <https://github.com/caltechlibrary/newt/issues>.

## Documentation

- [INSTALL](INSTALL.md)
- [user manual](user-manual.md)
- [About](about.md)

