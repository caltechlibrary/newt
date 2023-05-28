---
<<<<<<< HEAD
title: "newt(1) user manual | Version 0.0.1 a191ba4"
pubDate: 2023-05-20
=======
title: "newt(1) user manual | Version 0.0.1 3db768b"
pubDate: 2023-05-26
>>>>>>> 17da2eec9b8f816c2d8e8b64c743425cd5e957ba
author: "R. S. Doiel"
---

# NAME

newt

# SYNOPSIS

newt [CONFIG_FILE]

# DESCRIPTION

*newt* is a microservice designed to work along side Postgres,
PostgREST, and Pandoc server. It provides URL routing and data flow
between the microservices based on a list of "routes" described in a
CSV file.  newt is part of the Newt Project which is exploring
building web services, applications and sites using SQL for data modeling
and define back-end service behaviors along with Pandoc templates used to
generate HTML consumed by the web browser.  newt supports data
hosted in Postgres databases via PostgREST JSON API as well as static
files contained in an "htdocs" directory (e.g. HTML, CSS, JavaScript,
and image assets). 

This goal of Newt Project is to be able to assemble an entire backend
from off the self services only requiring data modeling and end point
definitions using SQL and a Postgres database. Reducing the back-end
to SQL may simplify application management (it reduces it to a
database administrator activity) and free up developer time to focus
more on front end development and human interaction. It is also
hoped that focusing the back-end on a declaritive model will allow for
a more consistent and reliable back-end.

# OPTIONS

-help
: display help

-license
: display license

-version
: display version

-dry-run
: Load configuration and routes CSV but don't start web service


# CONFIGURATION

The three things newt needs to know to run are port number,
where to find the "routes" CSV file and a list of any POSIX environment
variables to import and make available inside the router.

{app_mame} can be configured via a POSIX environment.

~~~
NEWT_PORT="3030"
NEWT_ROUTES="routes.csv"
NEWT_ENV="DB_NAME;DB_USER;DB_PASSWORD"
export NEWT_PORT
export NEWT_ROUTES
export NEWT_ENV
~~~

It can also be configured using a configuration file.


~~~
newt_port = "3030"
newt_routes = "routes.csv"
newt_env = [ "DB_NAME", "DB_USER", "DB_PASSWORD" ]
~~~

The environment will load first then the configuration file if
provided. The configuration file takes presidence to the environment.

newt does not have secrets but could use secrets passed
in via the environment. This allows your routes CSV file to be safely
saved along side your SQL source code for your Newt Project.

# EXAMPLES

Configuration from the environment

~~~
	export NEWT_PORT="3030"
	export NEWT_ROUTES="routes.csv"
	export NEWT_ENV="DB_USER;DB_PASSWORD"
	newt
~~~

Configuration from a YAML file called "newt.yaml"

~~~
newt newt.yaml
~~~


