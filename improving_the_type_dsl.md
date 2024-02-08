
# Improving Newt's model and type DSL

GitHub has defined a YAML issue template DSL[^1] that they use to generate issues with a specific structure[^2]. The GitHub DSL includes a type attribute used to describe the user input expected for the issue. The types map to HTML5 form elements. I believe this DSL could also be used to support generating SQL data models.

The GitHub YAML issue syntax could be extended to better support data types specific to the application domains needed by libraries, archives and museums.

## Background

On Thursday 2024-01-25 Caltech Library's digital library development staff meeting included a discussion of using GitHub YAML issue templates as a means of quickly creating forms that could trigger a GitHub actions. The advantage of this approach is it saves developers time from coding up a bunch of HTML forms and then writing a yet another service to trigger the GitHub actions.  The downside of using this SAAS is vendor lock in.

## Problem

Newt's niche is as a rapid web application development platform for libraries, archives and museums. In the Newt 2023 prototype I implemented a small DSL on top of YAML to describe data models that are rendered as SQL. The SQL can be used to provide a JSON API via Postgres and PostgREST. The problem is the prototype's DSL is too narrowly focused to also describe a web form. If Newt is to be a rapid web application development tool then generating the HTML and SQL implementing a PostgREST JSON API is highly desirable.

## Proposed Solution

If I adopt the syntax supported by GitHub YAML issue templates[^3] (abbr: GHYT). I gain the advantage of our staff already being familiar with the GHYT DSL describing the issue templates. The GitHub syntax appears sufficient to render both SQL models and HTML 5 forms. Adding the SQL support is a matter of mapping GHTY input types to both general purpose HTML and their respective SQL data types.

There is room for further enhancement without breaking GHYT.  Libraries, archives and museums work with structured metadata. The metadata attributes often conform to a known data type. This is particularly true of identifiers (e.g. ORCID, DOI, ROR, arXiv). Using the existing GitHub YAML issue template syntax I could support a wider diversity of types through the value associated in the type attribute of the input element. 

This would require replacing original prototype DSL defined in a Go struct found in [type_dsl.go](type_dsl.go) with a new series of struct that map out GHYT. The current[^4] prototype's ModelDSL struct is minimal. If I enhance the ModelDSL struct to match the syntax described in the GitHub YAML issue schema documentation then I should be able to support the GitHub YAML syntax with minimal modifications to the rest of the Newt data router.  I can also improve the code generation ability with this simpler structure for mapping data and representations.

### Advantages of adopting a common syntax

GitHub's YAML syntax promises the following advantages

- Standardize syntax for library staff who implement web forms
- Leverage GitHub documentation approach with the Newt YAML DSL for modeling
- Supports a protyping a process in GitHub issue templates and GitHub action while leaving the door open to easily self host those processes in our machines

### Disadvantages to adopting a common syntax

- GitHub's YAML issue syntax is in beta and is subject to change
- Newt project will need to tracking any changes to remain compatible

## Footnotes

[^1]: DSL, Domain-specific language, See [Wikipedia definition for details](https://en.wikipedia.org/wiki/Domain-specific_language).

[^2]: See GitHub [Syntax for issue forms](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-issue-forms)

[^3]: See GitHub [Syntax for GitHub's form schema](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema#about-githubs-form-schema)

[^4]: See Newt prototype version [v0.0.6](https://github.com/caltechlibrary/newt/releases/tag/v0.0.6)

