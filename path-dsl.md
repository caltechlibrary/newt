
# A Path DSL

One of the challenges of implement a `newt` URL router is describing the routes and identifying any elements that need to be reused to process the routes.

An example would be translating a traditional blog path to the PostgREST API route.  The PostgREST API has a shallow API that is well suited to a data API. If you can simply describe the route in "blog form" and translate it to a PostgREST API call our `newt` router could be relatively simple.  This is a good canidate for using a domain specific language (DSL).  Here's an example. How do I list blogs posts expressed with the path `/blog/{YEAR}/{MONTH}/{DAY}` and translate that into a PostgREST API call that needs to be in the form of `http://user:secret@localhost:3000/blog?year={YEAR}&month={MONTH}&day={DAY}`? While this may make sense to implemented via URL rewrites in the front end web server it combersome if you want to support different front end web servers (e.g. Apache 2 and NginX). Using a path DSL the mapping could be simplified and allow the mapping to be stored in our Postgres database along side the data models.

What could a path DSL look like that would be simple enough to require minimal documentation, be extensable over time and be easy to read and understand?

# Routes

Routes describe a description of a URL path that match the characteristics described in our path DSL.  A route must be uniquely identifiable to be useful in the larger context of a URL routers.

# What is our path DSL and how do we evaluate a match?

Our evaluation function that takes a path value and path DSL (our route description) and returns two pieces of information. First is a boolean indicating if the path value provided match the path DSL expressed (i.e. matched our route described). The second returned value would be a map or dictionary holding any variable names and values encountered when validating our path value against our path expression. Note an empty map could be returned when the route (i.e. path DSL experssion) is not matched as well as when there are no varaible expressions provided in the route expression.

The folowing are examples of routes that could be matched but would return an empty map of variable names to values.

- `/about.html`, retrieve the content for "about.html"
- `/blog/feed/`, retrieve the content of a "blog feed" 
- `/`, retrieve the homepage of a site (e.g. "index.html")

NOTE: Paths are evaluated from left to right.

## route expressions

A more complex route would include one or more variable expressions. Variable expressions are delimited by an opening curly bracket `{` and closed by a closing curly bracket `}`. Evaluating a path against a route results in two values. The first is that the route matched (boolean) or not. The second return value is a map of the variable names and the resulting values. In additional to identifying the name of a variable to be bound it is easily important to know the start and end points of the resulting string bound to the name. How is this detailed?  Let's explore the following three examples.

1. `/blog/{YEAR}/{MONTH}/{DAY}`, a route to blog posts on a given day, this route evaluate against a request path such that `YEAR`, `MONTH`, `DAY` will contain the values is the requested path.  These could be held in a map or dictionary where the names they keys.
2. `/blog/{YEAR}/{MONTH}/{DAY}/{TITLE_SLUG}.html`, this route would result with the variables `YEAR`, `MONTH`, `DAY` and `TITLE_SLUG`. Note the title slug would not be expected to include ".html" as that was expressed as a literal part of the string.
3. `/blog/{YEAR}/{MONTH}/{DAY}/{TITLE_SLUG}.{EXT}`, In this example the "." seprating the `TITLE_SLUG` from the `EXT` is a literal. It indicates the separator between two variables.

In these examples we can see that we are working with directory like paths. We have delimiters for each path element as well as a possible delimiter between a filename and it's extension. A variable's value is found between two literal values or a starting literal value and the end of the path if no literal value is expressed to the right of the variable placeholder.

Our route DSL embeds the declaration for the path element variable in an example path structure. The variable's value is defined both by the delimiter named between our curly brackets as well as the starting literal and possibly trailing literal string. It functions like a string mask.

## Variable annotations

What about the case were there is no delimiter between two varaibles?
Let's re-use our blog path. What if the blog path is still expressed in a YYYYMMDD format but without the internal delimiting `/`?

- `/blog/{YEAR}{MONTH}{DAY}`, this would be too ambigious to parse

How do we know where a `YEAR` value ends and `MONTH` or `DAY` starts? We need some sort of annotation that would indicate where to delimit the values. For year, month and day values that could be the length of the strings.

- `/blog/{YEAR:string:4}{MONTH:string:2}{DAY:string:2}`, now we know that `YEAR` is four characters long and month and day are two

A variation might be to use a regular expression approach

- `/blog/{YEAR:re:2[0-9][0-9][0-9]}{MONTH:[0-1][0-9]}{DAY:[0-3][0-9]}`, now we know that `YEAR` is four characters long starts with '2' and is followed by three digits and similary and month and day are two long each containing a digit in a specific range.

In principle an annoation is separated by the "type" of annotation (so far I've descrived `str` for fixed length string and `re` for regular expression). I've used `:` colons to separate the value of the variable name, the type and any additional information required by the type such as a regexp string describing a match. In principle this approach gives room to grow types.

What exactly is a type of the input is string and boolean indicating a matched route and a map of variable names to strings?

A type for the purpose of our DSL is implemented by a validator. The validator needs to return three thnigs. First does if the path element matched the type being defined. This is a boolean value, true if it matches, false if it does not. The two other values would be the start and end points of the substract to extract from the map and assign into our map. Start and end points are a positive integer (i.e. an integer that is greater than or equal to zero). An invalid start point would be expressed as -1. An invalid or end of string end point would be expressed as a -1.

Using the validator mechanism our path DSL could support many types. You might define a `year` type as a four digit number, a `month` type as a two digit number from `01` to `12`. A `day` type might validate `01` through `31`. This would give using zero padding on the left. Just as we could create types for date or time elements we could also create types that mapped to identifiers. E.g. An "orcid" type would make sure we had a valid ORCID and the validator would return the ORCID's start and end points in the path.

NOTE: Routes still must start with a literal. The shortest "route" is `/`. The longest route would be expressed as `/{REST:string}` where the varaible `REST` contains the whole path without the leading `/`.

## Keeping the path DSL consistent and unambigous.

## Algorithm for route evolation

- Split the path into parts (e.g. directory elements, basename and extension)
- For directory element, basename or extension
    - if it is a literal compare with path, if match continue to next element
    - if element contains one or more variable expression, evaluating from left to right evaluate any liter part and continue through evalauting each variable expression
        - if variable expression matches populate our map with substring from requested path
- if all literal elements and varaibles match then route matches and true with a populated map is returned otherwise the evaluation returns false and an empty map


