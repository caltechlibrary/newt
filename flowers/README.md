
# Flowers, adding a configuration driven router

With metadata managment provided by Postgres + PostgREST, template engine provided by Pandoc server what's missing is assembling a web page without relying on JavaScript running in the web browser. This can be done by a simple URL request router that understands both PostgREST and Pandoc Server. I call this simple router Newt.

## A simple set of responsiblilities

Newt doesn't need to be complicated, what it is does need to know is how to get data from PostgREST and which template to send with the data to Pandoc srever before returning the result. If you're going to support submitting data (e.g. a web form or upload a file) to our web application then we also want to handle that use case. This is all pretty predictable.  We need a pattern that can be expressed as a literal path or regular expression, the content type to be returned and the HTTP method making the request. If data is being sent to the server we need (e.g. a POST or FILE upload) need to define we want to handle that (e.g. transform a form submission into JSON to send to PostgREST or upload the file to an application directory). Essentially we're defining a set of urls. These rules can be described by a simple table structure and uploaded into Postgres (e.g. via CSV file or SQL statements). Newt just needs to run the configuration defined to respond to each request it recieves. 

Newt life cycle. Newt is configured similarly to PostgREST. A configuration file or the environment defines how to access the Postgres server. The rest of the configuration is retrieved from a Postgres SQL table. The table defines settings like where to find Pandoc templates, where to store uploaded files (if that is allowed),  as well as a table that defines the routes supported by Newt.


