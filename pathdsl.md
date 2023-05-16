
# A PathDSL, a routing DSL

How to you describe mapping of one path to another? Many web frame works implement the concept of a "route" which is similar to a file path but may include a placeholder notation to easily extract values from the path and assign them to variables. The is what PathDSL seeks to address.

## Why another yet another routing DSL

Why another DSL for routes?  Surveying the route descriptions available in several Python and JavaScript fameworks each has chosen its own syntax. There isn't a concensus nor is there one that appears entirely suitable.  In addition the route syntaxes I've seen do not capture enough information about the variable they are working with to fully do their job in terms of identifying a path value that matches the expression or in extracting the variable path elements that will be exposed to the frameworks.

In a framework setting this isn't a bad thing as additional validation or conversion can be applied to the variables described by the route syntax.  The Newt router is very simple and does not have the luxury of a general purpose programming language that can further vet a varaible. In addition web applications based around the route concept often embed varaibles which represent known identifiers or data formats.  It would nice to have not just the variable placement in the path but also enough "type" information to apply a validator function before varifying a path value matches the route describe.

One of Newt's routing features is to take an in bound path and transform it into a URL suitable to access another micro service like PostgREST.  It would be nice if there was some visual symetry between describing the route and the template langauge representing the transformed route. Ulitimately these desired properties moved me to thinking about PathDSL that would allow easily implementation of path value matching a route description and a path value that could be easily transformed with a simple template expression into a new URL.

NOTE: This document if focused on discussing the concepts behind PathDSL not how to implement it in a specific language.

## Blogging URLs, a use case

PathDSL should be able to handle simple mappings such as those seen in blogs. Blog paths are often predictable.  A home page is at `/`, a feed of items might be at `/feed/` and individual blog posts might be found in a path formed by embedding four digit year, two digit month, two digit day and a title slug -- `/blog/{year}/{month}/{day}/{title-slug}`.  I've choosen to use curly brackets to delimit the variable elements in the path because handlebars style string formatting is well known.

How do we know when a path value matches a route?  For a literal path you can simply perform a string comparision but for a path with variables embedded you need to vet the variables for a match. This can be though of as a validation of the value held by the variable. What if we annotated the varaible names between the curly braces with some validatin information? A "year" certainly can be validate if we know
that it conforms to a four digit integer. Month and day could have a simple validation based on being an
integer with two digits allowing for leading zeros. These are common enough knowns that many languages provide standard libraries for validating date elements so why not just create a validation type for this
type of data, let's look at what that might look like for a blog post path.

~~~
/blog/{year Year}/{month Month}/{day Day}/{title-slug String}
~~~

Knowing the type (or validation method to use) our we could check a our PathDSL can evaluate if a given route matches the path value provided. Let's look at three path values.

| path value | is it valid? | year | month | day | title-slug |
| /blog/2023/02/28/my-daily-post | yes | 2023 | 2 | 28 | my-daily-post |
| /blog/2023/02/31/my-daily-post | no  | 2023 | 2 | invalid day | my-daily-post |
| /blog/2023/28/02/my-daily-post | no  | 2023 | invalid month | 2 | my-daily-post |

Knowing the type would let the router know that path is not valid even though it matches the general form we expect in a blog post path.

While the PathDSL would not implement specific type validation it would beable to indicate the type
of the variable validation to the host language implementing our router. In principle we could implement
validators for common privitive types like "Integer", "String", "Real", "Date", "Year", "Month", "Day" but
also common types of identifiers like "ORCID", "DOI", "ROR", "Zipcode", or "PhoneNumber". If you added regular expression you would have a rich type (validation) system for transforming URLs into another URL.

Here's an example of what our PathDSL would enable. Path expression --

~~~
`/blog/{year Year}/{month Month}/{day Day}/{title-slug}`
~~~

A transformed version of the input route could then be described simply as a handlebar template.

~~~
http://localhost:3000/blog?year={year}&month={month}&day={day}&title_slug={title-slug}
~~~

If the value matches both form and types of variables then you have enough information to call a microservice like PostgREST.

## Use case, representing filenames and extenions in a path

What if you have a service that can return different documents based on a file extension? How might that be expressed in our PathDSL? Take these three path values.

~~~
/blog/2023/02/28/my-daily-post.html
/blog/2023/02/28/my-daily-post.md
/blog/2023/02/28/my-daily-post.txt
~~~

Perhaps you're using two microservices behind our theoretical router. These three URL represent the same
content but returned in different formats has as can be rendered by Pandoc server. In this case you want to
be able to extract both the "title-slug" and a file extension. Here's how you might represent that in PathDSL.

~~~
`/blog/{year Year}/{month Month}/{day Day}/{title-slug basename}{ext extname}`
~~~

Many implementation languages support parsing path into directory components, filename and extension. PathDSL should be able to leveraget this. This suggests an algorithmic behavior in our PathDSL evaluation.

First PathDSL should split the path into it's component, then it should determine if the component is a
literal string or a variable definition. For each component the variable definition needs to conform to
it's type. The trailing element in a path can also have the types of "basename" (excluding the file extension), extname (including the file extension without the leading period). If you want trailing element to include both the file name and extension then we can use the "String" type to describe it.


## Algorithm for evaluating a path value against a PathDSL expression

- A PathDSL expression decode into it's path components parts
    - zero or more directory names
        - directory name can be a literal
        - directory name can be a variable definition
    - zero or one filename
        - a filename can be a literal
        - a filename can be a variable expression
- A path value is split into it's components
    - zero or more directory names
    - zero or one filename and extension
- Compare each path value element against each PathDSL expression elements
    - if both are literals
        - stop processing if they do not match, return false
    - if Path DSL component is a variable, valid the path value component against variable type
        - stop of it does not validate, return false
- If comparison completes without return false we have a match

## Variable definitions

In the use cases we've suggested how a variable in a PathDSL expression can include a variable name as well as validation information. Let's specify this in more detail.

- varaible definition start with an opening curly brace and conclude with a closing curly brace
- Following the curly brace is a variable name which is formed from a letter and one or more alphanumeric characters or an underscores or dashes, e.g.
    - `a`, `a1`, `a_long_variable_name`, `title-slug` are valid variable names
    - `""`, `1a`, `+1`, `{}`, `()`, `$foo` are not valid variable names
- a variable name is followed by a space and type expression
- a type expression starts with a letter and can be followed by one or more characters excluding a closing curly brace

Here's some example varaible definitions.

~~~
{year Year}
{month Month}
{day Day}
{orcid Regexp '[0-9][0-9][0-9][0-9]\-[0-9][0-9][0-9][0-9]\-[0-9][0-9][0-9][0-9]\-[0-9][0-9][0-9][0-9]'}
~~~

These would result in the following type maps expressed in JSON

~~~
{
    "year": "Year",
    "month": "Month",
    "day": "Day",
    "orcid": "Regexp '[0-9][0-9][0-9][0-9]\-[0-9][0-9][0-9][0-9]\-[0-9][0-9][0-9][0-9]\-[0-9][0-9][0-9][0-9]'"
}
~~~

The PathDSL does not define the supported types only that the information can be extracted from a PathDSL expression as a map between variable names and a type description. It is up to the specific PathDSL implementation to define how the type informaiton is interpreted.

## Variable decoding

If a path value matches a PathDSL expression then when the variables and values can be extracted as a map of variable names and values. The only constraint is that the map be expressable as a valid JSON object. E.g.

Given the PathDSL expression

~~~
/people/{orcid ORCID}
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

Given the PathDSL expression

~~~
/blog/{year Year}/{month Month}/{day Day}/{title-slug basename}{ext extname}
~~~

And the path value

~~~
/blog/2022/11/07/compiling-pandoc-from-source.html
~~~

A PathDSL implementation should return a map, dictionary or associative array with the values
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

## Reference materials

- [path-to-regexp](https://github.com/pillarjs/path-to-regexp)
- [URLPattern](https://developer.mozilla.org/en-US/docs/Web/API/URLPattern) at MDN
- [URLPattern](https://developer.chrome.com/articles/urlpattern/) at Chrome Developer site
- [Flask Route tutorial](https://pythonbasics.org/flask-tutorial-routes/)
- [router.js](https://github.com/tildeio/router.js/)
- [Azure application gateway routing](https://learn.microsoft.com/en-us/azure/application-gateway/url-route-overview#pathbasedrouting-rule)
- [React Router](https://reactrouter.com/en/main/route/route)
- [Nextjs routing](https://nextjs.org/docs/app/building-your-application/routing)
- [dJango routing](https://www.django-rest-framework.org/api-guide/routers/)

