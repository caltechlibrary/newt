
# Supporting complex form elements, history and brain storm

Galleries, libraries, archives and museums have many use cases for individual metadata records containing a lists of items. Example a citation may contain a list of authors or collaborators. 

HTML 5 provides a rich set of input types[^0] but does not include list element type[^1]. The data model of a web form without enhancement is flat. It is a series of key/value pairs where the data when transmitted is urlencoded and sent as text to the receiving web service. Let's put a pin in that for later.

[^0]: See MDN site for details, <https://developer.mozilla.org/en-US/docs/Web/HTML/Element/Input#input_types>

[^1]: For the this exploration I'm going to table discussing file and select element types.

## Simplification and trade offs

Repository and catalog systems current (2024) rely on heavy weight frameworks to provide a decent data entry experience[^2]. This approach that dates back to the early 21st century is not required. CSS improvements like variables and better units, JavaScript improvements including web workers and promise or DOM improvements like support for web components was delivered long ago. I can't remember the last time I felt it necessary to include [modernizer](https://modernizr.com/) or a shim library in a web page. Yet large web systems still build web pages like it 2010.

[^2]: See [ArchivesSpace](https://archivesspace.org)  and [Invenio RDM](https://inveniosoftware.org/products/rdm/)

I think with a couple of heuristics can can skip the CSS and JavaScript frameworks, write less code and still provide a good user experience.

1. Design for the smallest screen used to access your web site or application (often that is a phone today)
2. Use CSS to rearrange elements on the screen as the screen size expands (e.g. your screen might become a wall sized television in a meeting)
3. Limit JavaScript to orchestrating behavior (e.g. defining web components) or access interacting with data services
4. Only write the minimum of CSS to create a readable page 

Using these four heuristics along with vanilla JavaScript and CSS we can avoid whole tool chains of complexity. These are heuristics not rules. As designers and developers we need to make choices. Hopefully we make informed choices.

## The trouble with web components, a fallback problem

In 2010 the trouble with web components was browser support. That isn't the problem in 2024.  The big trouble with web components is JavaScript. JavaScript is used to define them and to run them[^3]. If the browser has JavaScript disabled, it is unavailable or it is not supported by the web browser then you're stuck[^4] with an usable component.

[^3]: WASM loaded via JavaScript can also be used to create UI elements this is exotic in 2024

[^4]: Lynx, Dillo and NetSurf are browsers that don't support JavaScript or WASM (Web Assembly). A network disruptions can also prevent JavaScript from reaching your favorite web browser. Filling out a web form in a subway can be a nice place to experience this phenomena.

So what to do about lack of list elements if web components can't relied upon to solve the problem?  Approach the problem using the ideas from progressive enhancement. First decide how you would edit our lists if web components and JavaScript had never made it into the web? An old school trick for column oriented data where you didn't want the overhead of a table element full of editable cells was to rendering the table as comma delimited output editable with a textarea element.  This was find for small table editing, each line was a row and commas as easy punctuation to explain and worst case the human could quote the cell that needed to continue an embedded comma. Not ideal but doable.  Our list of authors could be edited in this way right up until we have many columns or columns that hold multi dimensional data (e.g. list of attributes and values).

When implementing UI elements in the past I've used JSON to pass data to and from the browser, can we use pretty printed JSON in a textarea as a fallback?

Here's an example of an person object in JSON.

~~~json
{
  "family_name": "Doiel",
  "lived_name": "Robert",
  "identifers": [
    { 
      "orcid": "0000-0003-0900-6903",
      "clpid": "Doiel-R-S"
    }
  ]
}
~~~

This doesn't worked as a simple comma delimited text editable in an textarea element because of the "identifiers" field/column. If JSON is pretty printed and not too complex it can be easy enough to read. It is easy to pick out my name and edit it. There's a catch when it comes to editing JSON. It's trivially easy to accidentally mismatch a double quote, square bracket or curly brace.  Similarly adding or removing a colon or comma can create a problem too. Is there a better way to express this data structure?

YAML is a notation for expressing structured data.  Here's my person object as expressed using YAML.

~~~yaml
fmaily_name: Doiel
lived_name: Robert
idntifiers:
  - orcid: 0000-0003-0900-6903
    clpid: Doiel-R-S
~~~

YAML tends to look like a list. Punctuation still counts as it does in JSON, e.g. colon and dash mean specific things in YAML. Yet if you need to correct the spelling of my name or ORCID you can do so with less worry because we don't need to worry about quoting or about matching braces. I remains easy to read. It resembles a list in Markdown.

Let's look at a list of authors in JSON[^6].

~~~json
[
  {
    "family_name": "Doiel",
    "lived_name": "Robert",
    "identifers": [
      { "orcid": "0000-0003-0900-6903", "clpid": "Doiel-R-S" }
    ]
  },
  {
    "family_name": "Morrell",
    "lived_name": "Thomas",
    "identifiers": [
       { "orcid": "0000-0001-9266-5146", "clpid": "Morrell-T-E" }
    ]
  }
]
~~~

[^6]: I think lists of JSON objects definitely too much to typing even if you are a programmer used to working with JSON.

Now look at the same list expressed in YAML.

~~~yaml
- family_name: Doiel
  lived_name: Robert
  identifers:
  - clpid: Doiel-R-S
    orcid: 0000-0003-0900-6903
- family_name: Morrell
  lived_name: Thomas
  identifiers:
  - clpid: Morrell-T-E
    orcid: 0000-0001-9266-5146
~~

There is a learning curve to YAML. It's more complex than separating field values with a comma.  I need to know the rules about indentation, dashes and colons. If I am including large text blocks I need to know about the pipe character and indentation nuances. I think YAML provides an easy of reading over pretty printed JSON, it provides easier editing than JSON too if you are trying to fix a value in previously created YAML. YAML like JSON is available in port programming languages. YAML can easily be integrated server side (e.g. transforming a SQL query result into YAML is trivial), and browser side too[^7].

[^7]: See <https://yaml.org/> which lists two stable JavaScript libraries

What we want our web component to do then is sit on top of our plain old textarea holding YAML and present an easy to edit interface on top. That way if JavaScript is available we have a nice user experience. If not we can still edit the structure carefully via a standard textarea element. With this approach we can support restricted browsers like Lynx, Dillo and NetSurf while also providing a good experience for those you are using Firefox, Chrome, Edge or Safari.

There are additional benefits for using the textarea apprach to holding a YAML expression of a complex data structure.  I can test servers web form processing using curl or other simple HTTP client library of tool. I don't need to learn a framework or figure out how that impacts the data I send back to the server. I only need to learn how to implemenet web components. If I am clever I will use those components in multiple projects. 

On the database side I take my YAML, flip it to JSON and store the contents in a JSON column. Similarly I can easily take a JSON column and turn that into YAML.

I think this approach leverages the grain of the web while still offering the potential of a good user experience for managing metadata.

