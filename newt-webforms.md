
# Newt as data router, web forms

Often middleware is written to handle URL request paths as [routes](newt-router.md). The other common attribute in a route handle function is validating inputs. This is particularly true with web form submission. 

Newt's [Route DSL](route-dsl.md) does this partially by validating the path. It doesn't do that for URL query paramaters. Since Route DSL does establish a syntax for a simple markup of variable names and their types it could be leverage to also describe the contents of a GET or POST request. This would effectively allow us to validate form input by it's type much like we do with the Route DSL. Since webform don't support array or object input (without resorting to JavaScript) a simple web form could be seen as a map between input element's attribute name and a type like those found in Route DSL.

Here's an example of a our birds web form validation

~~~
{
    "bird" : "String",
    "place": "String",
    "sighted": "Datestamp"
}
~~~

If we use YAML to express the JSON then we'd get soemthing like

~~~
bird: String
place: String
sighted: String
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

## Use case: simple web forms and online magazine publisher

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

When adding an article we need a title, byline, article copy. When editing an article we access it through the article_id genrated on a SQL INSERT. To publish an article we need to set approved to true, set a publication date, volume number, issue and article number.


