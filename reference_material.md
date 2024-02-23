
# Reference Material

These are links to prior art, related efforts and resources for consideration in how Newt evolves as a prototype.

## Data Modeling

- [Syntax for GitHub form schema](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema)
- [A React component for building Web forms from JSON Schema](https://github.com/rjsf-team/react-jsonschema-form)
- [JSON Forms](https://jsonforms.io/docs/)
    - pretty heavy weight in the deployed results (i.e. renders to either React or Angular)
    - Uses JSON rather than YAML or TOML to describe forms which can be construed as less human friendly
- [YAQL](https://yaql.readthedocs.io/en/latest/getting_started.html), a YAML like query language that maps to SQL
- [pg_query](https://github.com/pganalyze/pg_query), a Ruby library to parse SQL and normalize into a data structure
- [htmx](https://htmx.org/), a Web Component like markup implementing the wrapping layer between JSON API and HTML structured markup
- [Yafowil, yet another form widget library](http://docs.yafowil.info/)
    - [Describe YAFOWIL forms with YAML](https://yafowil.readthedocs.io/en/latest/yaml.html)
- [MDN Web Components](https://developer.mozilla.org/en-US/docs/Web/API/Web_Components), useful for make specialize input elements, like ORCID, DOI, ArXiv identifier entry
- [MDN HTML5 input types](https://developer.mozilla.org/en-US/docs/Learn/Forms/HTML5_input_types)
- [Generate webforms based on YAML schema with pykwalify](https://github.com/cxong/pykwalify-webform)
- [Flask Open API 3](https://pypi.org/project/flask-openapi3/)


## SQL DB to REST JSON API

- [PostgREST](https://postgrest.org), build a REST JSON API directly on top of Postgres. This is what started my think about Newt prototype
- [MRS](https://dev.mysql.com/doc/dev/mysql-rest-service/latest/), MySQL REST Service
- [sqlite2rest](https://github.com/nside/sqlite2rest), Automatically RESTful OpenAPI server from SQLite database
- [Soul](https://github.com/thevahidal/soul), A SQLite REST and realtime server built on NodeJS (wonder if it runs in Deno?)

## SQL JSON support

- SQL dialects
    - [database guide](https://database.guide/), documentation for various SQL dialects including Postgres and SQLite
- Postgres
    - [Postgres JSON functions and operators](https://www.postgresql.org/docs/16/functions-json.html)
    - [Postgres JSON tutorial](https://www.postgresqltutorial.com/postgresql-tutorial/postgresql-json/)
- SQLite 3
    - SQLite [Overview JSON functions](https://sqlite.org/json1.html)
    - [JSON improvements in SQLite 3.38.0](https://tirkarthi.github.io/programming/2022/02/26/sqlite-json-improvements.html)
    - SQLite [JSON function enhancements (2022)](https://sqlite.org/src/doc/json-enhancements/doc/json-enhancements.md)
    - SQLite cli [docs](https://sqlite.org/cli.html), e.g. "Changing output formats" in section 5 covers `.mode json`

## Data transformation and mashups

- [Pipes](https://www.pipes.digital/docs) is a spiritual successor to [Yahoo! Pipes](https://en.wikipedia.org/wiki/Yahoo!_Pipes)

## Other approaches and background

- [path-to-regexp](https://github.com/pillarjs/path-to-regexp)
- [URLPattern](https://developer.mozilla.org/en-US/docs/Web/API/URLPattern) at MDN
- [URLPattern](https://developer.chrome.com/articles/urlpattern/) at Chrome Developer site
- [Flask Route tutorial](https://pythonbasics.org/flask-tutorial-routes/)
- [router.js](https://github.com/tildeio/router.js/)
- [Azure application gateway routing](https://learn.microsoft.com/en-us/azure/application-gateway/url-route-overview#pathbasedrouting-rule)
- [React Router](https://reactrouter.com/en/main/route/route)
- [Nextjs routing](https://nextjs.org/docs/app/building-your-application/routing)
- [dJango routing](https://www.django-rest-framework.org/api-guide/routers/)
- [HYDRA](https://www.markus-lanthaler.com/hydra/), Hypermedia-Driven Web API
- [HAL](https://stateless.group/hal_specification.html), Hypertext Application Language
    - [JSON-HAL](https://datatracker.ietf.org/doc/html/draft-kelly-json-hal-00)
- [JSON-LD](https://en.wikipedia.org/wiki/JSON-LD)
- [Richardson Maturity Model](https://en.wikipedia.org/wiki/Richardson_Maturity_Model), used to evaluate RESTful-ness in JSON API

## Identifiers for a minimal object definition of unique id with created and updated timestamps

- Go implementations
    - [UUID](https://pkg.go.dev/github.com/google/UUID), v4 UUID, go storage in Postgres, long for URL
    - [shortid](https://pkg.go.dev/github.com/teris-io/shortid)
    - [shortuuid](https://github.com/skorokithakis/shortuuid), Python shortuuid pacakge
    - [shortuuid](https://github.com/lithammer/shortuuid), a Port of the Python package to Go, algorithmically compatible with the Python implementation
    - [ulid](https://github.com/oklog/ulid), a logically sortable unique identifier package, ported from the JavaScript package implementing ulid

