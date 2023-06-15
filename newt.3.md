---
title: newt(1) 
pubDate: 2023-06-15
author: R. S. Doiel
---

# NAME

new variables

# SYNOPSIS

- var
   - NAME : TYPE

# DESCRIPTION

newt YAML files can contain variable definitions in the routes attribute.
Varaibles are declared as a key and value where the key is the variable name and the value is the variable's type.

newt uses variables to vet request paths that reference embedded variables as well as in getting the form data submitted with a request (e.g. GET, POST, PUT, PATCH). 

Varaibles can also be defined under the "model" attribute where they describe the scructure used to generate SQL and Pandoc templates.

# VARIABLE TYPES

Currently newt supports a limited number of types.

Date [FORMAT]
: A date type by default is expressed in a YYYY-MM-DD style, e.g. 2006-01-02 but you can specify the specific format using dsl that Go uses in its time package, see <https://pkg.go.dev/time>. E.g. `Date 2006` would express a four digit year, `Date 01` would express a two digit month, `Date 02` would express a two digit day.

String
: This conforms to a sequence of utf8 characters except for `/` when part of a path (pathes are split first before their variables are evaluated)

Integer
: An integer value (e.g. "10")

Real
: A float value expressed with a decimal point (e.g. "1.0")

Boolean
: A boolean type expressed as "true" or "false"

Basename
: A file's basename (filename without an extension)

Extname
: A file's extension (e.g. ".html", ".txt", ".rss", ".js")

Isbn10
: A ten digit ISBN

Isbn13
: A thirteen digit ISBN

Isbn
: A ISBN (either 10 ro 13 digit)

Issn
: An ISSN

DOI
: A DOI (digital object identifier)

Isni
: An ISNI

ORCID
: An ORCID identifier

Markdown
: A Markdown data type, Markdown text can be render as HTML and saved to another variable name and sent to the JSON source as part of the API request. 

# EXAMPLES

An example of a blog path expressed for a route.

```
routes:
    - var: { "yr": "Date 2006", "mo": "Date 01", "dy": "Date 02" }
    - req_path: /blog/${yr}/${mo}/${dy}
    - req_method: GET
    - api_url: http://localhost:3000/blog?year=${yr}&month=${mo}&day=${dy}
    - api_method: GET
    - api_headers: { "Content-Type": "application/json"     
    - pandoc_template: list_articles_for_date.tmpl
    - res_headers: { "Content-Type": "text/html" }
~~~
```

# SEE ALSO


