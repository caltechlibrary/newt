
# Newt

Newt is an experimental [microservice](https://en.wikipedia.org/wiki/Microservices) working with other "off the shelf" microservices. It is a proof of concept of a "new take" on the old web stack. The "new stack" is built on off the shelf microservices. The initial targets are Postgres+PostgREST and Pandoc server. In practice the Postgres+PostgREST can be replaced (or operating along side) any JSON data source (e.g. Solr, Opensearch) available via localhost.

## Typical newt stack

- Newt is a data router, form validator, static file microservice
- [Pandoc](https://pandoc.org) server provides a template engine to render JSON as something else
- JSON data sources
    - [Postgres](https://postgres.org) + [PostgREST](https://postgrest.org), SQL to JSON API
    - Solr, Elasticsearch/Opensearch full text search engines

The Newt prototype type runs on localhost and assumes all the JSON data sources also run on localhost. In the future support for JSON data sources requiring https protocol support may be added (e.g. JSON data sources like CrossRef, DataCite, ORCID use https protocol).  

In a production setting you would normally run Newt behind a front-end web server like Apache 2 or NginX. The advantage of this approach is that the front-end web service can provide https protocol support, authentication, and single sign-on integration. Newt just provides JSON data routing, form validation and limited static file service support.

Newt as an "off the shelf" microservice plays nicely with virtual hosting and containers. You have lots of flexibility to deploy from a single machines to the data center and the cloud.

## Orchestrating your app with Newt

Newt is configured with a YAML file. Currently the configuration file uses three main attributes

htdocs
: The path to the htdocs directory holding any static assets (e.g. CSS, JavaScript, HTML files, image assets)

env
: A list of environment variables available to Newt's routes and models

routes
: An object describing the mapping of an HTTP request to JSON data source and and optional Pandoc server processing

A routes includes the following attributes

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

pandoc_template
: (optional) The path to the pandoc template used to process the results of the JSON data source request results

Additional attributes maybe added to the YAML file in the future.

## Data flow in steps

1. Web browser => Newt
2. Newt => JSON data source
3. Newt => Pandoc
4. Newt => Web browser

The Pandoc processing is optional. In this way you could expose selected elements of your JSON data source to the front-end web service by mapping the api request but skipping the Pandoc template used to transform the data.

### Handling errors

Newt vets the initial request before contacting the JSON data source. If the request has a problem it will return an appropriate HTTP status code and message.  If the request to the JSON data source has a problem, it will pass through the HTTP status code and message provided by the JSON data source.  Likewise if Pandoc server has a problem Newt will forward that HTTP status code and message. If either the JSON data source or Pandoc server is unavailable Newt will return a "Gateway" http status code and message.

### Static file support

Newt first checks if a request is matched in one of the defined routes. If not it'll try to service the content from the "htdocs" location. If the file is not found then a 404 is returned.  If the "htdocs" value is not set then a 404 will get returned if no routes match. 

Note Newt's static file services are very basic. You can't configure mime type responses or modify behavior via "htaccess" files. If Newt is running behind a traditional web server like Apache 2 or NginX then you should use that service to host your static content as it will give you much more control.

### Handling secrets, scaling and limitations

Newt's YAML file does not explicitly contain any secrets. This was intentional.  You may need to pass sensitive data to your JSON data source for access (e.g. credentials like a username and password). This should be provided via the environment and the YAML file needs to include these environment variable names in the "env" attribute.  The environment variables can be used to contact the JSON data source. There is still a risk in that theoretically that data source could echo return sensitive information. Newt can't prevent that. Newt is naive in its knowledge of the JSON data source content it receives and hands of to Pandoc.  You're responsible for securing sensitive information at the database JSON data source level. 

While Newt was conceived as a  small scale web application platform for Libraries, Archives and Museums it is capable of scaling big as long as your JSON data source can scale big.  Using the Newt stack elements can all be run easily behind load balancers and in parallel across machines because they require no synchronized shared of data between them or instances of them. Postgres itself can be configured in a HA cluster to support high availability and high demand.

Presently Newt does not support file uploads. If you need that you'll need to write your own service to handle it. 

Newt runs exclusively as a localhost service. In a production setting you'd run Newt behind a traditional web server like Apache 2 or NginX. The front-end web service can provide access control via basic auth or single sign-on (e.g. Shibboleth). Newt plays nicely in a container environment, running as a system service or invoked from the command line.

## Taking advantage of JSON data sources

Newt prototype works with a JSON source that can be accessed via URL and HTTP methods (e.g. GET, POST, PUT, PATCH and DELETE) using http protocol. Many systems have JSON data API today (2023). This includes existing Library, Archive and Museum applications like Invenio RDM and ArchivesSpace as well as search engines like Solr and Elasticsearch.  This means Newt can be used to extend existing systems that provide a localhost JSON data source.

In the future Newt may be extended to support a JSON data source external to your machine (e.g. CrossRef, DataCite, ORCID). It would require the Newt code base to be updated to support https protocol in additional to the existing http protocol. 


## Motivation

My belief is that many web services used by Archives, Libraries and Museums can benefit from a simplified and consistent back-end. If the back-end is "easy" then the limited developer resources can be focused on the front-end. **An improved front-end offers opportunities to provide a more humane experience for staff and patrons.**. 

I've written many web applications over the years. Newt is focused on providing very specific glue to other microservices which already perform the core of an application (e.g. PostgREST presenting a data management engine and Pandoc server presenting a template engine).  Newt takes a declarative model to configuration. You define things.  From those definitions it delegates processing to JSON data sources and optionally to Pandoc server if JSON needs to be transformed into HTML. It does this by limiting it's responsibility to verifying routes declared in the YAML file and validating paths and form data requests.  Newt only runs on localhost only. It is happy to run behind a more full featured web service like Apache 2 or NginX.

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

## Documentation

- [INSTALL](INSTALL.md)
- [user manual](user-manual.md)
- [About](about.md)

