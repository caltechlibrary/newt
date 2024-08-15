
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
- [MDN Web Components](https://developer.mozilla.org/en-US/docs/Web/API/Web_Components), could be useful for making specialize input elements, like ORCID, DOI, ArXiv identifier entry
- [MDN HTML5 input types](https://developer.mozilla.org/en-US/docs/Learn/Forms/HTML5_input_types)
- [Generate web forms based on YAML schema with pykwalify](https://github.com/cxong/pykwalify-webform)
- [Flask Open API 3](https://pypi.org/project/flask-openapi3/)

## Web Form generation

- [Deform](https://docs.pylonsproject.org/projects/deform/en/latest/)

## SQL DB to REST JSON API

- [PostgREST](https://postgrest.org), build a REST JSON API directly on top of Postgres. This is what started my think about Newt prototype
- [MRS](https://dev.mysql.com/doc/dev/mysql-rest-service/latest/), MySQL REST Service
- [sqlite2rest](https://github.com/nside/sqlite2rest), Automatically RESTful OpenAPI server from SQLite database
- [Soul](https://github.com/thevahidal/soul), A SQLite REST and real time server built on NodeJS (wonder if it runs in Deno?)

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

- [Doing More with Less: Orchestrating Serverless Applications without an Orchestrator. with David Liu, Shadi Noghabi, and Sebastian Burckhardt. (to appear) In proceedings of the 20th USENIX Symposium on Networked Systems Design and Implementation (NSDI) 2023. Paper](https://www.amitlevy.com/papers/2023-nsdi-unum.pdf)
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
- [HAR](http://www.softwareishard.com/blog/har-12-spec), HTTP Archive Request
- [GoJa](https://github.com/dop251/goja), Pure Go JavaScript engine
- [HTTPie-go](https://github.com/nojima/httpie-go), A go implementation of HTTPie, useful to debug API interactions
    - Fork with out stale package references, [HTTP-Pie-go](https://github.com/HexmosTech/httpie-go)
- [Lama2](https://github.com/HexmosTech/Lama2), a plain text powered rest client (not to be confused with Meta's LLM)
- [Go's extended html docs](https://pkg.go.dev/golang.org/x/net/html), HTML tokenizer and parser package

## Go implementations for unique identifiers

It appears it may be desirable when generating SQL code and web forms to allow for the automatic creation of three object attributes as some sort of minimally viable object, unique object id with created and updated timestamps

~~~json
{
    "identifier": "xxxxxxxxx",
    "created": "timestamp value for object creation",
    "updated": "timestamp, updates on change"
}
~~~

This could be implemented directory into the generated SQL or have an Identifier service that sits in front of the DB JSON API.

- [UUID](https://pkg.go.dev/github.com/google/UUID), v4 UUID, go storage in Postgres, long for URL
- [shortid](https://pkg.go.dev/github.com/teris-io/shortid)
- [shortuuid](https://github.com/skorokithakis/shortuuid), Python shortuuid package
- [shortuuid](https://github.com/lithammer/shortuuid), a Port of the Python package to Go, algorithmically compatible with the Python implementation, currently at v4 and still actively developed (April 2024), see <https://pkg.go.dev/github.com/lithammer/shortuuid/v4>
- [ulid](https://github.com/oklog/ulid), a logically sortable unique identifier package, ported from the JavaScript package implementing ulid
- [Cuid2](https://github.com/nrednav/cuid2), a port of [Cuid2]() JavaScript implementation
- [Sonyflake](https://github.com/sony/sonyflake), a uid inspired by Twitter's snowflake
- [Mongo BSON ObjectID](https://pkg.go.dev/github.com/mongodb/mongo-go-driver/bson/objectid), Go implementation from Mongo's go client library
- [xid](https://github.com/rs/xid)


Blog posts on identifiers

- [Identify Crisis, how modern applications generate unique identifiers](https://medium.com/javascript-scene/identity-crisis-how-modern-applications-generate-unique-ids-39562736f557)
- [Unique identifiers in three parts](http://antoniomo.com/blog/2017/05/21/unique-ids-in-golang-part-1/)
- [What is a ulid and why should you start using it?](https://dev.to/nejos97/what-is-ulid-and-why-should-you-start-using-it-14j9)

## Web Components

If Newt shipped with a set of LAM oriented web component that could allow for a more complex data model to be addressed by Newt. E.g. lists of authors are a web component would probably have not just the name fields but also identifiers like ORCID and ISNI attached to them. Something worth exploring.

Internet Archives seems to be investing in the [Lit](https://lit.dev) component library heavily.

