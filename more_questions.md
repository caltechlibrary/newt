
# Questions?

## Where is my development time going to be spent?

The developer writes YAML to generate the back end data management and templates to render web pages for your application. You can enhance the generated code further if you want. I suspect that you'll spend most of your time improving the human experience of your application through improving the HTML markup in the templates, writing some CSS and perhaps enhancing behavior with the JavaScript run in the web browser. If you need to enhance the back end you work in SQL. If you simply need to improve the rendering of your database results then you are working with Mustache templates.

## What about security, integration with single sign-on or other websites or services?

The `newtrouter` is a simple web service providing data routing based on its YAML configuration. It's a team player. In a production setting it should be used behind a front end web server like Apache 2 or NGINX. That latter can be configured to support single sign-on systems like Shibboleth[^80]. The front end web service controls access and handles securing the connection with the web browser. The front end web service proxies to the Newt router. Newt router receives requests and runs the data pipelines on localhost. The data pipelines can be composed of off the shelf software like Postgres+PostgREST, Solr and template engine to turn your JSON into a web page. Having a clear division of responsibilities helps in securing your web application. Since Newt router only knows how to talk to services on localhost you can keep it contained and prevent it from being used for nefarious actions off system. Like Newt router Newt Mustache is constrained to localhost for similar reasons.

Limiting Newt web service applications to localhost keeps them simple. Doing the minimum limits the attack surface for those who want to do mischief.  Neither `newtrouter` or `newtmustache` write to disk or require secrets. They only communicate via localhost using HTTP protocol.

[^80]: Shibboleth is a common single sign-on platform in research libraries, universities and colleges.

## What about "scaling"?

`newtrouter` is just a router. Aside from reading configuration at start up it doesn't maintain state. `newtmustache` functions the same way, read in the configuration and just run. By assigning different ports you can also run many instances of them. This makes it possible to run them in parallel, behind load balancer or even through proxying spread them across many machines. The instances don't share data or coordinate. They start up wait for a request and providing an answer.

So what does this all mean? In principle a Newt based applications can scale big as the slowest element of your pipeline service.

## What is the anatomy of a Newt based web application?

Newt application development is friendly to version control systems (e.g. Git). It consists of a Newt configuration file, along with generated SQL files, HTML templates and any static web assets you've added. A typical disk layout of a Newt project could look like this-

- `/` project folder
  - `app.yaml` would holds the configuration of our Newt tools
  - `postgrest.conf`, the configuration file for PostgREST
  - `htdocs` this directory holds your static content needed by your web application
  - `app_setup.sql` and `*_models.sql` these are the SQL files used by your application to define your models and behaviors in Postgres
  - `*.tmpl` your Mustache templates for turning JSON into HTML
  - `CITATION.cff` or `codemeta.json` for project metadata

