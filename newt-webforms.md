
# Newt as data router

Often middleware is written to handle URL request based on the path element of the URL. Sometimes this is referred to as a "route" or "end point". Newt also supports the concept of [routes](newt-router.md) but carries it further with the concept of a request route being mapped to a JSON data source and continuing processing via Pandoc running in server mode. In custom middleware a common first order of business in a route function is to validate inputs.  Newt also includes this. It does this via annotations in the YAML configuration in the routes, var attribute. A small Model [DSL](https://en.wikipedia.org/wiki/Domain-specific_language "Domain Specific Language") is used to define a varaibles type. The [type definition](model-dsl.md "Model DSL") is what Newt uses to validate the inputs.

Here's an example of a our birds web form validation in the "var" block
of the route definition.

~~~
var:
  bird: String
  place: Text
  sighted: Date
~~~

This would allow validating a simple webform like 

~~~
<form method="POST" action="/">
<div>
  <label for="bird">Bird</label>
  <input id="bird" name="bird" type="text" value="">
</div>
<div>
  <label for="place">Place</label>
  <input id="place" name="place" type="text" value="">
</div>
<div>
  <label for="sighted">Sighted on</label>
  <input id="sighted" name="sighted" type="date">
</div>
<button id="record-bird" type="submit">Add Bird Sighting</button>
</form>
~~~


## Use case: an article submittion form

In this use case we'll add support to our birds demo for a local birder magazine. The idea is to support writups or articles from the local birding group. Something between a blog, magazine or news letter.

Our articles should have a title, a byline (one or more authors names/identifiers), article copy, title slug, created timestamp, updated timestamp, approval status, published datestamp, a volume number, issue number and article number. An article identifier could be assigned when the article is iniitaially submitted. 

The SQL for that would look like

~~~
CREATE TABLE zine (
   article_id SERIAL,
   title TEXT,
   title_slug VARCHAR(256),
   byline VARCHAR(256),
   article_copy TEXT,
   created TIMESTAMPZ NOT NULL DEFAULT NOW(),
   updated TIMESTAMPZ NOT NULL DEFAULT NOW(),
   approved BOOLEAN,
   published DATE,
   volume INTEGER,
   issue INTEGER,
   article_no INTEGER
);
~~~

The var attribute to validate the page for updating the zine table would look something like this.

~~~
var:
  article_id: "integer*",
  title: "string",
  title_slug: "string",
  byline: "string",
  article_copy: "markdown",
  created: "timestamp now",
  updated: "timestamp now",
  approved: "boolean",
  published: "date",
  volume: "integer",
  issue: "integer",
  article_no: "integer"
~~~

When adding an article we need a title, byline, article copy. When editing an article we access it through the article_id genrated on a SQL INSERT. To publish an article we need to set approved to true, set a publication date, volume number, issue and article number within the issue.  

The model dsl closely aligns with the SQL but uses the terminology common to main stream programming jargon that has evolved since the original invention of SQL in the 1970s. 

There are three variables defined with more than a single word. The first one is our article_id. There is an asterisk at the end of "integer". The asterisk indicates this will be the key used to retrieve this model values from the JSON API.  It also allows Newt to easily map our model description from JSON or YAML into Postgres dialect of SQL.  Similarly our two timestamps have a modifier of "now"  This is used when mapping to Posgres SQL timestamps with a default of now.  It can also be used in the route definition to inject a new "now" value like with the updated timestamp without the need for that to be expressed directly in the web form via JavaScript.



