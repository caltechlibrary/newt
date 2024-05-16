
# Thoughts on supporting complex form elements

Galleries, libraries, archives and museums have many use cases for individual metadata records containing a lists of items. Example a citation may contain a list of authors or collaborators.

HTML5 provides a rich set of input types[^0]. HTML5 does not include list input element type[^1]. The data model of a web form without enhancement is flat. It is a series of key/value pairs where the data when transmitted is urlencoded and sent as text to the receiving web service.

[^0]: See MDN site for details, <https://developer.mozilla.org/en-US/docs/Web/HTML/Element/Input#input_types>

[^1]: For the this exploration I'm going to table discussing file and select element types.

Let's put a pin in our a list of authors and look at the technical situation.

## Simplification and trade offs

Repository and catalog systems current (2024) rely on heavy weight frameworks to provide a decent data entry experience[^2]. This approach dating back to the early 21st century is not required. CSS improvements like variables and better units, JavaScript improvements including web workers and promises as well as the DOM improvements supporting for web components were delivered long ago. I can't remember the last time I felt it necessary to include a shim library in a web page to support them. Why do our large web systems still build web pages like they did in 2010.

[^2]: See [ArchivesSpace](https://archivesspace.org)  and [Invenio RDM](https://inveniosoftware.org/products/rdm/)

I think with a couple of heuristics can can skip the CSS and JavaScript frameworks. We can write less code. We can provide a better user experience.

1. Design for the smallest screen used to access your web site or application (often that is a phone today)
2. Use CSS to rearrange elements on the screen as the screen size expands (e.g. your screen might become a wall sized television in a meeting)
3. Limit JavaScript to orchestrating behavior (e.g. defining web components or access interacting with data services)
4. Only write the minimum of CSS and JavaScript to create an accessible page

Using these four heuristics with vanilla HTML, CSS and JavaScript can help us avoid a whole tool chain of complexity. **I am proposing heuristics not rules**. As designers and developers we need to make choices.

When we work with complex data structures then a web component might make sense but we want to think about these heuristic first.

## A solution to creating an editable list of author data

HTML5 is missing an editable list element. In 2024 it is reasonable to implement a web component to provide this feature. We don't need frameworks just JavaScript, CSS, HTML and the DOM. We can create a component the shows an complex person object, we can also create a component that handles lists of person objects. We can do the same for corporate objects or funders as well. The nice thing about this approach is the objects can conform to the expectations of our organization.

## The trouble with web components, a fallback problem

Web components were a big topic back in 2010. In 2010 the trouble with web components was browser support. In 2024 the big trouble with web components is JavaScript. JavaScript is used to define them and to run them[^3]. If the browser has JavaScript disabled, it is unavailable or it is not supported by the web browser then you're stuck[^4] with an unusable component.

[^3]: WASM loaded via JavaScript can also be used to create UI elements this is exotic in 2024

[^4]: Lynx, Dillo and NetSurf are browsers that don't support JavaScript or WASM (Web Assembly). A network disruptions can also prevent JavaScript from reaching your favorite web browser. Filling out a web form in a subway can be a nice place to experience this phenomena.

Are we out of luck creating an editable complex data structure like an editable list of authors? We can take a page from progressive enhancement. First decide how you would edit our lists if web components and JavaScript had never made it into the web? An old school trick for column oriented data for small tables was to use a simple textarea and instructions. This avoided the whole problem of managing a complex table arrangement in your web form.  Using a textarea meant each line represented a row and a comma is used to separate each cell (column). The comma is an easy punctuation to explain. Worst case the human can double quote the cell that needs an embedded comma. Not perfect but doable for a simple table.

Our list of authors could be edited in this way unless it includes multi dimensional data. Example a list of identifier types and their values. So maybe comma delimited lines are too simple a fallback.

When implementing UI elements in the past I've used JSON to pass data to and from the browser, can we use pretty printed JSON in a textarea as a fallback?

Here's an example of an person object in JSON.

~~~json
{
  "family_name": "Doiel",
  "lived_name": "Robert",
  "identifiers": [
    {
      "orcid": "0000-0003-0900-6903",
      "clpid": "Doiel-R-S"
    }
  ]
}
~~~

If JSON is pretty printed and not too complex it can be easy enough to read. There's a catch when it comes to editing JSON. It's trivially easy to accidentally mismatch a double quote, square bracket or curly brace.  Similarly adding or removing a colon or comma can create a problem too. Is there a better way to express this data structure?

## A solution to our a fallback problem

YAML is a notation for expressing structured data.  Here's my person object as expressed using YAML.

~~~yaml
family_name: Doiel
lived_name: Robert
idntifiers:
  - orcid: 0000-0003-0900-6903
    clpid: Doiel-R-S
~~~

YAML tends to look like a list. Punctuation still counts as it does in JSON, e.g. colon and dash mean specific things in YAML. Yet if you need to correct the spelling of my name or ORCID you can do so with less worry because we don't need to worry about quoting or about matching braces. I remains easy to read. It resembles a list in Markdown.

I've represented a person, what about a list of people? Let's look at a list of authors in JSON[^6].

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
~~~

Each person in the list is delimited by the dash where and the person object is indented two more spaces. It’s compact, I can see where my entry stops and Tom’s starts.

There is a learning curve to YAML. It's more complex than separating field values with a comma and rows by a carriage return.  I need to know the rules about indentation, dashes and colons. If I am including large text blocks I need to know about the pipe character and indentation nuances. On the other hand I am not worried about mismatch braces and double quotes. I think YAML provides an easy of reading experience compared to pretty printed JSON. YAML is easier to edit accurately than JSON. YAML is easier to edit than JSON in large part by avoiding the double quote and brace mismatch problem by using indentation. YAML like JSON is available in most programming languages. YAML can easily be integrated server side (e.g. transforming a SQL query result into YAML is trivial), and browser side too[^7].

[^7]: See <https://yaml.org/> which lists two stable JavaScript libraries

Implementing our web component for managing authors can be done by sitting it on top of a textarea holding YAML. This is similar to WYSIWYG editors manages Markdown held in a textarea.  If JavaScript is available (the usual case situation) we have a nice user experience. If JavaScript is unavailable we can still edit the YAML.  With this approach we can support restricted browsers like Lynx, Dillo and NetSurf while also providing a good experience for those you are using Firefox, Chrome, Edge or Safari.

There are additional benefits in this approach. Testing the web service processing our web form data can be done easily with curl, a simple HTTP client library or other tool.  This enhances our options when debugging our web application. Additional I don't need to learn a framework, figure out how that impacts my structured data when it is sent back to the server. I only need to learn how to implement a web components. If I am clever I will use those components in multiple projects.

On the database side I take the YAML, flip it to JSON and store the contents in a JSON column. Similarly I can easily take a JSON column and turn that into YAML before populating the contents of a web form.

This approaches aligns with the historical grain of HTML and HTTP while still offering the potential of a good user experience when managing complex metadata.
