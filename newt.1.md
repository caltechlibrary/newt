---
title: "newt(1) user manual | Version 0.0.1 a191ba4"
pubDate: 2023-05-20
author: "R. S. Doiel"
---

# NAME

newt

# SYNOPSIS

newt [CONFIG_FILE]

# DESCRIPTION

*newt* is a microservice designed to work along side Postgres,
PostgREST, and Pandoc server. It provides URL routing and data flow
between the microserves based on a simple configuration table holding
routing information. It is part of the Newt Project which is exploring
building web services, applications and sites using SQL for data modeling
and define back-end service behaviors along with Pandoc templates used to
generate HTML consumed by the web browser.  newt supports data
hosted in Postgres databases via PostgREST JSON API as well as static
files contained in an "htdocs" directory (e.g. HTML, CSS, JavaScript,
and image assets). 

This goel of Newt project is to be able to assemble an entire backend
from off the self services only requiring data modeling and route
definitions using SQL and a Postgres database. Reducing the back-end
to SQL may simplify application management (it reduces it to a
database administrator activity) and free up developer time to focus
more on front end development and human interaction. It is also
hoped that focusing the back-end on a declaritive model will allow for
a more consistent and reliable back-end.

# CONFIGURATION

newt can be configured through the environment or through
a PostgREST is configuration file. It adds only an optional
uri for the Pandoc server (if used) when it runs on a non-standard
port.

~~~
db-uri = "postgres://birds:my_secret_password@localhost:5432/birds"
db-schemas = "birds"
db-anon-role = "birds_anonymous"
# This is used by Newt to know where to find the Pandoc server
# on localhost.
pandoc-server-port = "3030"
~~~

# EXAMPLE

~~~
	newt postgrest.conf
~~~


