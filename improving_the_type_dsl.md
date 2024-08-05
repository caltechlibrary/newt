
# Improving Newt's model and type DSL

GitHub has defined a YAML issue template DSL[^51] that they use to generate issues with a specific structure[^52]. The GitHub DSL includes a type attribute used to describe the user input expected for the issue. The types map to HTML5 form elements. This inspired a new take and a basis for a similar DSL could also be used to support generating SQL data models.

The GitHub YAML issue syntax can be simplified and better support data types specific to the application domains needed by libraries, archives and museums.

## Background

On Thursday 2024-01-25 Caltech Library's digital library development staff meeting included a discussion of using GitHub YAML issue templates as a means of quickly creating forms that could trigger a GitHub actions. The advantage of this approach is it saves developers time from coding up a bunch of HTML forms and then writing a yet another service to trigger the GitHub actions.  The downside of using this SAAS is vendor lock in.

## Problem

Newt's niche is as a rapid web application development platform for libraries, archives and museums. In the Newt 2023 prototype I implemented a small DSL on top of YAML to describe data models that are rendered as SQL. The SQL can be used to provide a JSON API via Postgres and PostgREST. The problem is the prototype's DSL is too narrowly focused to also describe a web form. If Newt is to be a rapid web application development tool then generating the HTML and SQL implementing a PostgREST JSON API is highly desirable. Newt continues to be prototyped and refined in 2024.

## Proposed Solution

If the data model DSL was inspired by GitHub YAML issue template syntax[^53] (abbr: GHYITS). Most of that syntax describes the HTML5 web form elements needed to generate and process a issue request. I've simplified that model further focusing on the HTML5 form element itself.  Both approaches have the advantage of describing HTML5 input elements and can easily infer a SQL data type. It was through experimentation that I realized directly support GHYITS was unnecessary and significantly increased Newt's implementation complexity for little benefit. If compatibly is needed it would be straight forward to create a program that cross walked the more complex GHYITS to the simpler Newt data model.

### Advantages of adopting a common syntax

GitHub's YAML syntax promises the following advantages

- Standardize syntax for library staff who implement web forms
- Leverage GitHub documentation approach with the Newt YAML DSL for modeling
- Supports a prototyping a process in GitHub issue templates and GitHub action while leaving the door open to easily self host those processes in our machines

### Disadvantages to adopting a common syntax

- GitHub's YAML issue syntax is in beta and is subject to change
- Most of the GitHub YAML issue syntax is unrelated to Newt's data modeling needs.
- Newt project would need to track changes to remain compatible
- one for one compatibility between Newt data model and GHYITS brings little advantage as Newt comes with an interactive modeler which removes the challenge of typing and maintaining complex YAML[^54].

## Footnotes

[^51]: DSL, Domain-specific language, See [Wikipedia definition for details](https://en.wikipedia.org/wiki/Domain-specific_language).

[^52]: See GitHub [Syntax for issue forms](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-issue-forms)

[^53]: See GitHub [Syntax for GitHub's form schema](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema#about-githubs-form-schema)

[^54]: See Newt prototype version [v0.0.9](https://github.com/caltechlibrary/newt/releases/tag/v0.0.9)
