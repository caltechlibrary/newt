
# A PathDSL, a routing DSL

One of the challenges of implement a a URL router that is general purpose in front of PostgREST is how you describe the requested router coming from the front end web server and transform that into the route you need to use to get data form PostgREST.  While a pipeline conceptually is easy to describe in our table (e.g. req_router, req_method, req_content_type, req_data, actions) how you describe the route mapping needs to be easy to learn, easy to read and easy to maintain.  This is an exploration of an approach to a routing DSL for `newt`, the URL router designed to work with PostgREST and Pandoc server.

NOTE: PathDSL does not describing implementation in a specific programming language. Ideally it should be "easy to implement" in the language of your choice. 

## Why another yet another routing DSL

Typically server side frameworks include among other things a routing DSL. You see this in frameworks such as Flask, jDango, Ruby on Rails, Sinatra, etc.  In our stack based on Postgres+PostgREST and a front end web server there is now framework to map human friendly URLs to a data API like PostgREST.

There are numerous implementations and there doesn't appear to be a specific strong case that suggests adopting an existing DSL for routing using our approach to simplifying our back-end.

## Blogging URLs, a use case

An example would be translating a traditional blog path to the PostgREST API route.  The PostgREST API has a shallow API that is well suited to a data API. If you can simply describe the route in "blog form" and translate it to a PostgREST API call our `newt` router could be relatively simple.  

For the purposes of discussion I am going to delimit our router dynamic elements with curly brackets.

How do form a "route" description to list blogs posts expressed with the path `/blog/{YEAR}/{MONTH}/{DAY}` and translate that into a PostgREST API call that needs to be in the form of `http://user:secret@localhost:3000/blog?year={YEAR}&month={MONTH}&day={DAY}`? 

In principle you could simply use regular expression notation to identify matching parameters quickly answer the question of if a given URL path matches our route. The problem I have with regular expressions is they quickly become unreadable and difficulty to maintain with even a moderately complex expressions. Additionally returning an array of matched parts isn't the most developer friendly way to reuse any information you want to extract from your matched route. Regular expression need to remain in our toolbox but should not be used to define our routes.

Assuming a route has some embedded data we want to reuse how to me bind those to developer friendly variable names?  If we can delimit variable names in the route expression then we can build a dictionary of mapped names and values which we return when we have a matching route.

In many template languages variables are identified by some starting delimiter and sometimes and terminating one. Some examples I've commonly seen are `${varname}`, `{varname}`, `$(varname)`, `<VARNAME>`. Many languages that include formatted strings use curly braces to delimit the variable names. I will continue with curly brackets as I explore a minimal useful routing DSL.

# Exploring Routes

Routes describe a description of a URL path that match the characteristics described in PathDSL.  A route must be uniquely identifiable to be useful in the larger context of a URL routers like `newt`.

# How do we identify a route match?

Our evaluation function that takes a path value and route returns two pieces of information. First is a boolean indicating if the path value provided match the route expressed. If the route matched we also want to returned name and value map (aka dictionary) found when we evaluated if the route matched.

It is possible the described route is explicit or literal. In the case of a match our evaluation function should return true and an empty name/value map.  The following are examples of explicit routes which would not be associated with any variable names or values.

- `/about.html`, retrieve the content for "about.html"
- `/blog/feed/`, retrieve the content of a "blog feed" 
- `/`, retrieve the homepage of a site (e.g. "index.html")

NOTE: Paths are evaluated from left to right.

##e Routes with variable names and values

A more complex route would include one or more variable expressions. Variable expressions are delimited by an opening curly bracket `{` and closed by a closing curly bracket `}`. This indicates where in the path the variable may be found (assuming the route matches). The value associated with the variable name is a sub string of the path. Using examples of a blog URL paths the route might look like

1. `/blog/{YEAR}/{MONTH}/{DAY}`, a route to blog posts on a given day, this route evaluate against a request path such that `YEAR`, `MONTH`, `DAY` will contain the values is the requested path.  These could be held in a map or dictionary where the names they keys.
2. `/blog/{YEAR}/{MONTH}/{DAY}/{TITLE_SLUG}.html`, this route would result with the variables `YEAR`, `MONTH`, `DAY` and `TITLE_SLUG`. Note the title slug would not be expected to include ".html" as that was expressed as a literal part of the string.
3. `/blog/{YEAR}/{MONTH}/{DAY}/{TITLE_SLUG}.{EXT}`, In this example the "." separating the `TITLE_SLUG` from the `EXT` is a literal. It indicates the separator between two variables.

The question now aside from the liter prefixes in our example how do we know where the value we're going to associate with our variable name ends if we're extracting sub-strings?

Many language provide path parsing libraries and can return a list of directories, file basename and extension. In the examples above you could map directory elements into the variables assuming the variables consume the whole element of the path.  This limits the types of routes we can express.

Another challenge of this simple mapping based on existing path delimiters and variable names is it lacks even rudimentary validation. Here's two examples that would match example one but are unlikely to both be valid in our blog context - `/blog/2023/02/28` (valid) `/blog/down/is/up` (probably not what we want). It is desirable to have some level of validation for the values we're going to assign to our variable names and in deciding if a route matches.

Ideally we should be able to specify a path as a sequence of variable names
and parse them appropriately.

- `/blog/{YEAR}{MONTH}{DAY}`, this is ambiguous without more information describing our route

## Variable annotations

How do we know where a `YEAR` value starts and ends? Same for `MONTH` or `DAY`? I think an annotation along side our variable name might be the solution. An easy to implement annotation might be a fixed length string.

- `/blog/{YEAR string 4}{MONTH string 2}{DAY string 2}`, now we know that `YEAR` is four characters long and month and day are two

A variation might be to use a regular expression (re) approach

- `/blog/{YEAR re 2[0-9][0-9][0-9]}{MONTH:[0-1][0-9]}{DAY [0-3][0-9]}`, now we know that `YEAR` is four characters long starts with '2' and is followed by three digits and similarly and month and day are two long each containing a digit in a specific range.

Note I've separated our variable names from our annotations using a space. If an annotation was not present then the assumption of that the variable consumes the remainder of the URL path.  When we process the path value and compare with our route the variable declarations become the delimiting factor. This gives us the ability to describe a match using the whole URL if necessary without relying on a language's implementation of path parsing.

What do I mean by "type" if what our route evaluation returns matching status and a map of variable names pointing at string values?  The first space indicates that we have a type rule that needs to be validated. After the second space we have any additional parameters need to evaluate the type. In our current examples this is string length or a regular expression. The type annotation describes a validation rule that is to be applied and that rule application should return the match state sub string holding the matched sub-string. This approach makes it straight forward to add additional types.

Going back to our original blog path I can express both a file path delimited string and a numeric date string path and still populate my variable parts.

- `/blog/{YEAR year}/{MONTH month}/{DAY day}/{TITLE_SLUG basename}{EXT extname}`

The "year" type would check for a four digit sub-string, it then could check if the value is reasonable. Likewise "month" and "day" should make sure the strings are two digits long and in the range of "01" to "12" for months and "01" to "31" for day. The "basename" and "extname" types could evaluate the whole URL path and return an appropriate filename and extension mapped to "TITLE_SLUG" and "EXT". Over time a collection of variable types could validate elements of the URL we wish to map.

The downside of what I've described is adding to yet another routing DSL. But I think the potential in the context of newt is that this trade off is reasonable to eliminate the need to write custom middle ware for many cases where I would like to build a front-end website on a back-end built from micro services like Postgres+PostgREST and Pandoc server.

## Algorithm for route evaluation

- read the path value from left to write
- if prefix is literal compare with literal in route
- else if route has a variable extract validate the value
    - if valid save the sub-string in our map and continue processing path
    - else we don't have a match, return false and empty name/value map

## Reference materials

- [path-to-regexp](https://github.com/pillarjs/path-to-regexp)
- [URLPattern](https://developer.mozilla.org/en-US/docs/Web/API/URLPattern) at MDN, [URLPattern](https://developer.chrome.com/articles/urlpattern/) at Chrome
- [Flask Route tutorial](https://pythonbasics.org/flask-tutorial-routes/)
- [router.js](https://github.com/tildeio/router.js/)
- [Azure application gateway routing](https://learn.microsoft.com/en-us/azure/application-gateway/url-route-overview#pathbasedrouting-rule)
- [React Router](https://reactrouter.com/en/main/route/route)
- [Nextjs routing](https://nextjs.org/docs/app/building-your-application/routing)
- [dJango routing](https://www.django-rest-framework.org/api-guide/routers/)

