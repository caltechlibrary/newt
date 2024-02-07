
# Newt

The Newt Project is an experiment in rapid web application development for libraries, archives and museums.  Newt's approach is to integrate off the shelf applications like Postgres, PostgREST compose web applications for curating metadata. Newt accomplishes this through configuration of two tools via a YAML file. Newt provides a program called `newt` that provides data routing and static file services.  `newtgen` provides code generation to setup Postgres+PostgREST as well as HTML forms and templates. This reduces the back end (server side) software construction process to writing the YAML file and possibly updating some templates.

## Motiviation

In the past decade I've watch the comlexity in web applications grow. While moving away from monolithic software or platform systems (e.g. Drupal, Wordpress and Zope) can be seen as a simplificaiton the reality is we're often still left writing or mainting complex middleware when creating custom software for libraries, archives and museums. I believe we can simplify this process through re-aligning the where processing happens and how we pull the pieces together.

In 2024 JSON API services are common. They are common both as software as a service but also architecturally in building libraries applications. You can see this in applications like Invenio-RDM and ArchivesSpace.  Both Invenio-RDM and ArchivesSpace spend allot of code providing those services. Is it possible to do this without writing code? Is it possible to do this with off the shelf software and a little configuration? At Caltech Library we've been using Postgres combined with PostgREST to provide this capability. This nice thing about this combo, Postgres+PostgREST is that it can be scaled up readily.

JSON API are nice but not necessarily what you want to present to humans when interecting with a web application. While in principle you can write JavaScript that runs browser side and fetches data from the JSON API there are many reasons why this is a bad idea[^1]. Fortunately there are tools that accept JSON as import and a template then return HTML (or other formats you've templated). A good example of that type of tool is Pandoc. It's very good at converted structured data from one format to another. It even can run as a web service that accepts a POST of a JSON object containing your data and template and turns the transformed JSON.  It works really quickly and scales nicely.

[^1]: See <https://infrequently.org/2024/01/performance-inequality-gap-2024/> for a nice discussion of the problem.

This missing bit is mapping the browser request to the JSON API through the template engine returning the result. This lead to the creastion of the Newt web service.

Using this combination of web services the data modeling and management happens in Postgres. The transform happes via the template engine. The routing happens via `newt`. But that still leaves setting up the data base, modeling data in SQL, writing views and functions to managing the data. On the front facing side writing templates by hand can quickly become tedious, particularly for web forms that contain allot of fields.  The Newt project expanded this to include code generation by using a YAML configuration file describing what is needed.

## What does Newt Project provide?

- [newt](newt.1.md) a [web service](https://en.wikipedia.org/wiki/microservices) designed for working with other "off the shelf" micro services. It functions both as a data router and as a static file server. It is responsible for vetting requests against the models described in the YAML file. The same models used to generate the SQL for the database and the templates for user with a template engine.
- [newtgen](newtgen.1.md) is a command line program that reads the Newt YAML file and generates SQL and templates used to build your application.  The SQL generated currently target Postgres+PostgREST.
- [newtmustache](newmustache.1.md) is a recent additional to the prototype suite. When discussing using Pandoc templates with colleagues some people pointed out they didn't like Pandoc tempaltes. Newt's template engine uses [Mustache](https://mustache.github.io/) template markup to render JSON content. It functions like Pandoc in server mode. You POST a request to the template engine that includes both your template and JSON data and it returns the transformed result. This demonstrates that you can swap out the template engine and still use the Newt web service. 

## How is Newt speeding up development?

1. `newtgen` can read Newt's YAML file and renders SQL suitable for bootstraping your Postgres+PostgREST database and JSON API
2. `newtgen` can read Newt's YAML file and render templates 
3. `newt` can read Newt's YAML file and provide a data router and static file sever

This leaves your backend contructed from "off the shelf" services but still allows you to support custom data models and HTML look and feel. Back end development is creating and evolving your YAML file, generating SQL and templates, then using those to test your application in your favorite web browser.

## Where is development time spent?

The developer writes YAML to generate the backend and compose those services. The developer generates and manages any templates needed by the template engine. Newt web service provides the data pipeline based on the routes defined in the YAML. The developer and further enhance the front end using static HTML, CSS and JavaScript.


With Newt your development time is focused on three areas.

1. Modeling your data in YAML
2. Modeling your "routing" in YAML
3. Enhancing the user expereience browser side using HTML5, CSS and if needed JavaScript (a.k.a. traditional front end development)

Newt provides a clean divison of labor between static files and web services that can be combined to form the front end of your web application.

## What about security, single sign-on, integration with other websites or services?

The `newt` command is a simple web service providing data routing based on its YAML configuration. It's a team player.  In a production setting it should be used with front end web server like Apache 2 or NGiNX. Both can be configured to support single sign-on sytems like Shibboleth[^2]. Newt web service then runs the data pipeline via other localhost based services like Postgres+PostgREST, Solr and Pandoc server. Having a clear division of responsibilities helps reasonable securing your web applciation

- front end web server (e.g. Apache 2, NGiNX) enforces access control then proxies to Newt 
- Newt receives requests and maps them to a static resource or a data pipe line formed from a JSON data source and template engine
    - Newt doesn't know about users, roles aore permissions
    - Newt only routes data after validating the request against the defined data models in the YAML file
    - Newt doesn't store secrets
- The the database actually manages the content. It has all the usual access control mechanisms available in databasses. In the case of Postgres+PostgREST can take this a step further by wrapping database access in a Pg/PL SQL procedure or function that checks a only returns data when it deems it is appropriate
- Newt can integrate external JSON API only if you've setup a proxy relationship between a localhost end point that Newt contacts and the remote URL that supplies the extenel service


If you need to integrate with an external JSON service (e.g. CrossRef, ORCID, ROR) this can be done by configuring your front end web service to also provide an internal proxy to the remote service and expose it on localhost.

Limitting Newt web service to localhost keeps is very simple, restricts it accessing external services without setup on your part and without having to impose an aditional access control layer and user/role/password management inside Newt.

[^2]: Shibboleth is a common single sign-on platfor in research libraries, universities and colleges.

### What about "scaling"?

Newt is just a data router. You can safely run many Newt instances in parallel as needed. They can run them on the same machine or separate machines. The instances don't share data or coordinate. They just read their configuration files when they start up and route data accordingly.

A typical Newt stack of PostgREST+Postgres and Pandoc server also can scale up. You can run as many instances of Pandoc server or PostgREST as you need. You can spread them across machines. They are both stateless systems like Newt. A Postgres database provides consistency and can be configured in a high availability cluster. In short Postgres scales.

### Annatomy of a Newt based web application

Newt application development is friendly to version control systems (e.g. Git). It consists of a Newt configuration file, along with generate SQL files, HTML templates and any static web assets you've added. A typical disk layout of a Newt project could look like this-

- `/` project folder
  - `htdocs` this directory holds your static content needed by your web application
  - `*.sql` these are the SQL files used by your application to define your models and behaviors in Postgres
  - `templates` this directory holds your Pandoc or Mustache templates
  - `tests` this directory holds your tests of your data model
  - `app.yaml` would hold the a Newt configuration file (this is an example name for the configuration file)
> Newt, a type of salamander, doesn't seek attention. It does its own thing. You only notice a salamander if you look carefully.

## About the Newt source repository

Newt is a project of Caltech Library's Digital Library Development group. It is hosted on GitHub at <https://github.com/caltechlibrary/newt>. If you have questions, problems or concerns regarding Newt you can use GitHub issue tracker to communicate with the development team. It is located at <https://github.com/caltechlibrary/newt/issues>.

### "Someday, maybe" ideas to explore

- Integrate S3 object store support as a data source
- Support other rendering engines besides Pandoc server

## Documentation

- [INSTALL](INSTALL.md)
- [user manual](user-manual.md)
- [About](about.md)

