
# Type DSL, a domain specific language describing a variables' type

Type DSL evolved from thinking about mapping of one URL path to another. Many web frameworks implement the concept of a "route" which is similar to a file path but may include a placeholder notation for values bound to variable names. The variable names are then available to route handler functions. If you are extracting values and binding them to varaible names it is extremely helpful to have type information.

Newt needed to support embedded typed variables in it mapping ability. It also was needed by Newt to support generating SQL bootstrap code from a simplified Model description.  Newt uses the type information to vet inputs bound to a varaible name as well as when generating bootstrep SQL for a Newt+Postgres+PostgREST+Pandoc (N3P) application.

I surveyed the route descriptions available in several Python and JavaScript frameworks. There was no consensus about how to declare variables and nothing that would allow vetting the variables by type. Newt implemented a Type DSL to provide a way to remedy this situation. 

## Blogging URLs, a use case

TypeDSL should be able to handle simple mappings such as those seen in blogs. Blog paths are often predictable.  A home page is at `/`, a feed of items might be at `/feed/` and individual blog posts might be found in a path formed by embedding four digit year, two digit month, two digit day and a title slug. Let's first look at how that would be described as a Mustache template -- `/blog/${year}/${month}/${day}/${title-slug}`.  This is pretty easy to read. 

How do we know when a path value matches a route?  For a literal path you can simply perform a string comparison but for a path with embedded variables you need to vet the variables to make sure they make sense. A type check needs to be done. In our simple Mustache version though there is no type information. The place holder `${year}` might hold an integer or maybe something completely unrelated.  What if we added an annotation about the variables type? A "year" certainly can be validate. Years are normally four digit numbers.  Likewise month and day could have a simple validation based on being an integer with two digits allowing for leading zeros. These are common enough date formats that many languages provide via standard libraries. Creating types for these types becomes a matter of wrapping our implementation of referencing the variables with a type system. Here's how we might express that in YAML.

~~~
var:
  year: Year
  month: Month
  day: Day
  title-slug: String
~~~

Now we can validate our path against the types assigned.

~~~
/blog/${year}/${month}/${day}/${title-slug}
~~~

Now we know the type and validation method to apply to the embedded variable. Let's explore some values and see if they can be validated.

path value                      year month day title-slug    is it valid?
------------------------------  ---- ----- --- ------------- -----------------
/blog/2023/02/28/my-daily-post  2023 02    28  my-daily-post yes
/blog/2023/02/31/my-daily-post  2023 02    31  my-daily-post no, invalid day
/blog/2023/28/02/my-daily-post  2023 28    01  my-daily-post no, invalid month

Knowing the type lets the router know that path is not valid and reject it if it is not without contacting other services. 

While the TypeDSL specification is not intended to describe a canonical set of types an implementation of a TypeDSL certain would.  It should be easy to implement the type system with one or two functions per type.  I suspect as TypeDSL are implemented a canonical core set of type will emerge.

In my prototype implementation of TypeDSL I plan to implement the following types and validating primitive data types -- "Integer", "String", "Real", "Boolean", "Date", "Year", "Month", "Day", "Hour", "Minute", "Second".  The Python project [IdUtils](https://idutils.readthedocs.io/en/latest/) also suggests a common set of identifiers used in the Library, Archive and Museum communities.  A subset of these will be implemented for the prototype of TypeDSL (e.g. "ORCID", "DOI", "ROR").

Using the previously defined types our URL provides the validated values
we need to transform to a API URL.

~~~
http://localhost:3000/blog?year=${year}&month=${month}&day=${day}&title_slug=${title-slug}
~~~

If the request value matches both form and type of embedded variables then you have enough information to call a microservice like PostgREST, Solr or Opensearch.

## Use case, representing filenames and extensions in a path

What if your request URL uses file extensions to indicate the choice of template to use when rendering with Pandoc?  Here's three examples where the content might be the same but the format changes (e.g. HTML, Markdown and plain text).

~~~
/blog/2023/02/28/my-daily-post.html
/blog/2023/02/28/my-daily-post.md
/blog/2023/02/28/my-daily-post.txt
~~~

In our prior example each embedded variable was contained in one element of the path.  A filename can be thought of as one element or as a "basename" and "extension". How can we allow for that use case?  In the use case about the "title-slug" corresponds to the "basename" but the file extension indicates which template to send to Pandoc server.

Here's how you might represent that in TypeDSL.

~~~
var:
  year: Year
  month: Month
  day: Day
  title-slug: Basename
  ext: Extname
~~~

We can then use our request path expression.

~~~
`/blog/${year}/${month}/${day}/${title-slug}${ext}`
~~~

Many implementation languages support parsing path into directory components as well as filename and extension. TypeDSL should be able to leverage this. This suggests an algorithmic behavior in our TypeDSL evaluation.

First TypeDSL should split the path into it's component, then it should determine if the component is a literal string or a variable definition. For each component the variable definition needs to conform to it's type. The trailing element in a URL path can also have the types of "basename" (excluding the file extension), extname (for the file extension value). In the case where you want to treat both the basename and extension as a single unit we have the String type shown previously.


## Algorithm for evaluating a path value against a TypeDSL expression

- A TypeDSL expression decode into it's path components parts
    - zero or more directory names
        - directory name can be a literal
        - directory name can be a variable definition
    - zero or one filename
        - a filename can be a literal
        - a filename can be a variable expression
- A path value is split into it's components
    - zero or more directory names
    - zero or one filename and extension
- Compare each path value element against each TypeDSL expression elements
    - if both are literals
        - stop processing if they do not match, return false
    - if TypeDSL component is a variable, valid the path value component against variable type
        - stop of it does not validate, return false
- If comparison completes without return false we have a match

## Variable definitions

In the use cases we've suggested how a variable in a TypeDSL expression can include a variable name as well as a type annotation to use to validate the value bound to the variable name. Let's specify this in more detail.

- variable definition start with two opening curly brace and conclude with a closing curly brace (i.e. basic Mustache template style notation)
- following the opening curly braces is a variable name which is formed from a letter and one or more alphanumeric characters or an underscores or dashes (i.e. what would be a valid attribute name in a JSON object). E.g.  - `a`, `a1`, `a_long_variable_name`, `title-slug` are valid variable names
    - `""`, `1a`, `+1`, `{}`, `()`, `$foo` are not valid variable names
- a variable name is followed by a space and type expression
- a type expression starts with a letter and can be followed by one or more characters excluding a closing curly brace

Here's some example variable definitions.

~~~
var:
  year: Year
  month: Month
  day: Day
  orcid: ORCID
~~~

These would result in the following type maps expressed in JSON

~~~
{
    "year": "Year",
    "month": "Month",
    "day": "Day",
    "orcid": "ORCID"
}
~~~

When a request URL is evaluated against the route template's type it will return a simple object with variable names mapped to string versions of the values. This will allow transforming both data API URLs and calls to Pandoc service via a simple Mustache like template implementation.

The TypeDSL does not define the supported types only that the information can be extracted from a TypeDSL expression as a map between variable names and a type description. It is up to the specific TypeDSL implementation to define how the type information is interpreted.

## Variable decoding

If a path value matches a TypeDSL expression then when the variables and values can be extracted as a map of variable names and values. The only constraint is that the map be expressible as a valid JSON object. E.g.

Given the TypeDSL expression

~~~
var:
  orcid: ORCID
~~~

Request path we're describing

~~~
/people/${orcid}
~~~

and the path value

~~~
/people/0000-0003-0900-6903
~~~

The resulting map would look like this JSON

~~~
{
    "orcid": "0000-0003-0900-6903"
}
~~~

Given the TypeDSL expression

~~~
var:
  year: Year
  month: Month
  day: Day
  title-slug: Basename
  ext: Extname
~~~

Request path expression

~~~
/blog/${year}/${month}/${day}/${title-slug}${ext}
~~~

And the path value

~~~
/blog/2022/11/07/compiling-pandoc-from-source.html
~~~

A TypeDSL implementation should return a map, dictionary or associative array with the values
converted to the type suggested in the variable's type definition. The constraint is that the map can be expressed as a JSON object. E.g.

~~~
{
    "year": 2022,
    "month": 11,
    "day": 7,
    "title-slug": "compiling-pandoc-from-source",
    "ext": ".html"
}
~~~

In this case our types "Month", "Day", "Year" converted the values to JSON numbers and the rest were left as JSON strings.

## SQL generation, a use case

When exploring working building application directory on Postgres via PostgREST one of the reactions I encountered was a resistence to SQL. This I think is unfortunate as modern SQL, exspecially in Postgres, is a rich capable langauge for managing and manipulating data.  One thing I've noticed over the years is that colleagues will often feel more confortable modifying SQL then writing it from scratch.  The Type DSL was initially developed for vetting routes but it occurred to me that is can also be used to describe a data model.  This lead to the "models" attribute in a Newt YAML file. RDMS support table datastructures. These are easily calculated from a simple key/value notation. It then is possible to bootstrap simple data models in Postgres by evaluating a map of variable names and their type descriptions.  Once you know the table structure you can also calculate some basic views and functions for working with the table (e.g. most applications need to support CRUD, create, read, update and delete).

