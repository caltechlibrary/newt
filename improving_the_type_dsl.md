
# Improving Newt's model and type DSL

GitHub has defined a YAML issue template DSL[^1] that they use to generate issues with a specific structure[^2]. The GitHub DSL includes a type attribute used to describe the user input expected for the issue. The types map to HTML5 form elements. I believe this DSL could also be used to support generating SQL data models.

The GitHub YAML issue syntax could be extended to better support data types specific to the application domains needed by libraries, archives and museums.

## Background

On Thursday 2024-01-25 Caltech Library's digital library development staff meeting included a discussion of using GitHub YAML issue templates as a means of quickly creating forms that could trigger a GitHub actions. The advantage of this approach is it saves developers time from coding up a bunch of HTML forms and then writing a yet another service to trigger the GitHub actions.  The downside of using this SAAS is vendor lock in.

## Problem

Newt's niche is as a rapid web application development platform for libraries, archives and museums. In the Newt 2023 prototype I implemented a small DSL on top of YAML to describe data models that are rendered as SQL. The SQL can be used to provide a JSON API via Postgres and PostgREST. The problem is the prototype's DSL is too narrowly focused to also describe a web form. If Newt is to be a rapid web application development tool then generating the HTML and JavaScript for a web form that maps to SQL needed to implement PostgREST JSON API is highly desirable.

## Proposed Solution

If I adopt the syntax supported by GitHub YAML issue templates[^3] I gain the advantage of our staff already being familiar with the GitHub YAML DSL describing the issues web forms. The GitHub syntax appears sufficient to render both SQL models and HTML 5 forms. Adding the SQL support is a matter of mapping the DSL types used to describe HTML form elements to their respective SQL data types.

Libraries, archives and museums work with structured metadata. The metadata attributes often conform to a known data type. This is particularly true of identifiers (e.g. orcid, doi, ror, arXiv). Using the existing GitHub YAML DSL syntax I could support a wider diversity of types through the value associated in the type attribute. 

The type attribute in the GitHub syntax controls both the rendered Markdown/HTML markup but also representation in SQL used to render the JSON API via Postgres+PostgREST.

Newt's type DSL is implemented using a Go struct found in [type_dsl.go](type_dsl.go). The current[^4] prototype's ModelDSL struct is minimal. If I enhance the ModelDSL struct to match the syntax described in the GitHub YAML issue schema documentation then I should be able to support the GitHub YAML syntax with minimal modifications to Newt's type DSL. An additional tool can then be included with Newt to support rendering HTML or Markdown forms from same YAML model I use to generate the SQL for Postgres+PostgREST.

### Advantages of adopting a common syntax

GitHub's YAML syntax promises the following advantages

- Standardize syntax for library staff who implement web forms
- Leverage GitHub documentation approach with the Newt YAML DSL for modeling
- A process that is prototype using GitHub issue templates and actions could be ported to a Newt base application eliminating the problem of SAAS vendor lock in

### Disadvantages to adopting a common syntax

- GitHub's YAML issue syntax is in beta and not formalized
- Newt project would need to tracking an informal specification changes 

## Footnotes

[^1]: DSL, Domain-specific language, See [Wikipedia definition for details](https://en.wikipedia.org/wiki/Domain-specific_language).

[^2]: See GitHub [Syntax for issue forms](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-issue-forms)

[^3]: See GitHub [Syntax for GitHub's form schema](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema#about-githubs-form-schema)

[^4]: See Newt prototype version [v0.0.6](https://github.com/caltechlibrary/newt/releases/tag/v0.0.6)

