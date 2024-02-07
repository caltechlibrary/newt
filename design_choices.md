
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

