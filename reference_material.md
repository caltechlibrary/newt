
# Reference Material

These are links to prior art, related efforts and resources for considatiion in how Newt evolves as a prototype.

## Data Modeling

- [Syntax for GitHub form schema](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema)
- [JSON Forms](https://jsonforms.io/docs/)i
    - pretty heavy weight in the deployed results (i.e. renders to either React or Angular)
    - Uses JSON rather than YAML or TOML to describe forms which can be construed as less human friendly
- [YAQL](https://yaql.readthedocs.io/en/latest/getting_started.html), a YAML like query langauge that maps to SQL
- [pg_query](https://github.com/pganalyze/pg_query), a Ruby library to parse SQL and normalize into a data structure
- [htmx](https://htmx.org/), a Web Component like markup implementing the wrapping layer between JSON API and HTML structured markup
- [MDN Web Components](https://developer.mozilla.org/en-US/docs/Web/API/Web_Components), useful for make specialize input elements, like ORCID, DOI, ArXiv identifier entry
- [MDN HTML5 input types](https://developer.mozilla.org/en-US/docs/Learn/Forms/HTML5_input_types)

## SQL DB to REST JSON API

- [PostgREST](https://postgrest.org), build a REST JSON API directly on top of Postgres. This is what started my think about Newt prototype
- [MRS](https://dev.mysql.com/doc/dev/mysql-rest-service/latest/), MySQL REST Service
- [sqlite2rest](https://github.com/nside/sqlite2rest), Automatically RESTful OpenAPI server from SQLite database

