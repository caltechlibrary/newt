---
author: R. S. Doiel
createdDate: 2024-02-16
title: The book of Newt
---

-   [[1]{.toc-section-number} The Book of
    Newt](#the-book-of-newt){#toc-the-book-of-newt}
-   [[2]{.toc-section-number} Preface, the Newt Experiment, February
    2024](#preface-the-newt-experiment-february-2024){#toc-preface-the-newt-experiment-february-2024}
    -   [[2.1]{.toc-section-number} Summer, Fall
        2023](#summer-fall-2023){#toc-summer-fall-2023}
    -   [[2.2]{.toc-section-number} Fall 2023 thru January
        2024](#fall-2023-thru-january-2024){#toc-fall-2023-thru-january-2024}
    -   [[2.3]{.toc-section-number} February 2024, a second
        prototype](#february-2024-a-second-prototype){#toc-february-2024-a-second-prototype}
    -   [[2.4]{.toc-section-number} The second
        prototype](#the-second-prototype){#toc-the-second-prototype}
-   [[3]{.toc-section-number} What is Newt? Why does it
    matter?](#what-is-newt-why-does-it-matter){#toc-what-is-newt-why-does-it-matter}
    -   [[3.1]{.toc-section-number} What's "off the
        shelf"?](#whats-off-the-shelf){#toc-whats-off-the-shelf}
    -   [[3.2]{.toc-section-number} What's my point about where data
        modeling and management take place? Why is that a Newt
        thing?](#whats-my-point-about-where-data-modeling-and-management-take-place-why-is-that-a-newt-thing){#toc-whats-my-point-about-where-data-modeling-and-management-take-place-why-is-that-a-newt-thing}
    -   [[3.3]{.toc-section-number} "Write SQL? I need to know that?
        Ugh."](#write-sql-i-need-to-know-that-ugh.){#toc-write-sql-i-need-to-know-that-ugh.}
    -   [[3.4]{.toc-section-number} Wait, JSON is great but I need a web
        page!](#wait-json-is-great-but-i-need-a-web-page){#toc-wait-json-is-great-but-i-need-a-web-page}
    -   [[3.5]{.toc-section-number} Where does that leave us
        developers?](#where-does-that-leave-us-developers){#toc-where-does-that-leave-us-developers}
    -   [[3.6]{.toc-section-number} Newt, three tools for your
        toolbox](#newt-three-tools-for-your-toolbox){#toc-newt-three-tools-for-your-toolbox}
        -   [[3.6.1]{.toc-section-number} Off the
            shelf](#off-the-shelf){#toc-off-the-shelf}
        -   [[3.6.2]{.toc-section-number} Newt
            tools](#newt-tools){#toc-newt-tools}
    -   [[3.7]{.toc-section-number} Newt design
        choices](#newt-design-choices){#toc-newt-design-choices}
-   [[4]{.toc-section-number} Newt YAML
    syntax](#newt-yaml-syntax){#toc-newt-yaml-syntax}
    -   [[4.1]{.toc-section-number} the "application"
        property](#the-application-property){#toc-the-application-property}
    -   [[4.2]{.toc-section-number} the "routes"
        property](#the-routes-property){#toc-the-routes-property}
        -   [[4.2.1]{.toc-section-number} a route
            object](#a-route-object){#toc-a-route-object}
    -   [[4.3]{.toc-section-number} the "models"
        property](#the-models-property){#toc-the-models-property}
        -   [[4.3.1]{.toc-section-number} a model
            object](#a-model-object){#toc-a-model-object}
    -   [[4.4]{.toc-section-number} input
        types](#input-types){#toc-input-types}
    -   [[4.5]{.toc-section-number} Example Newt YAML for router and
        code
        generator](#example-newt-yaml-for-router-and-code-generator){#toc-example-newt-yaml-for-router-and-code-generator}
    -   [[4.6]{.toc-section-number} templates
        property](#templates-property){#toc-templates-property}
        -   [[4.6.1]{.toc-section-number} template object
            model](#template-object-model){#toc-template-object-model}
-   [[5]{.toc-section-number} Newt Generator
    Explained](#newt-generator-explained){#toc-newt-generator-explained}
    -   [[5.1]{.toc-section-number} Bringing up
        Postgres+PostgREST](#bringing-up-postgrespostgrest){#toc-bringing-up-postgrespostgrest}
        -   [[5.1.1]{.toc-section-number} Generating some
            SQL](#generating-some-sql){#toc-generating-some-sql}
        -   [[5.1.2]{.toc-section-number} Generate our PostgREST
            configuration](#generate-our-postgrest-configuration){#toc-generate-our-postgrest-configuration}
    -   [[5.2]{.toc-section-number} Ready for testing our JSON
        API](#ready-for-testing-our-json-api){#toc-ready-for-testing-our-json-api}
    -   [[5.3]{.toc-section-number} Generating some templates for a web
        UI](#generating-some-templates-for-a-web-ui){#toc-generating-some-templates-for-a-web-ui}
-   [[6]{.toc-section-number} Newt Mustache
    explained](#newt-mustache-explained){#toc-newt-mustache-explained}
    -   [[6.1]{.toc-section-number} Overview](#overview){#toc-overview}
    -   [[6.2]{.toc-section-number} Testing your templates with Newt
        Mustache](#testing-your-templates-with-newt-mustache){#toc-testing-your-templates-with-newt-mustache}
        -   [[6.2.1]{.toc-section-number} Why Newt
            Mustache?](#why-newt-mustache){#toc-why-newt-mustache}
        -   [[6.2.2]{.toc-section-number} Why Mustache
            templates?](#why-mustache-templates){#toc-why-mustache-templates}
-   [[7]{.toc-section-number} Newt router
    explained](#newt-router-explained){#toc-newt-router-explained}
    -   [[7.1]{.toc-section-number}
        Overview](#overview-1){#toc-overview-1}
    -   [[7.2]{.toc-section-number} A simple
        example](#a-simple-example){#toc-a-simple-example}
    -   [[7.3]{.toc-section-number} Changes from the first prototype to
        the
        second.](#changes-from-the-first-prototype-to-the-second.){#toc-changes-from-the-first-prototype-to-the-second.}
    -   [[7.4]{.toc-section-number} Misc](#misc){#toc-misc}
-   [[8]{.toc-section-number} Reference
    Material](#reference-material){#toc-reference-material}
    -   [[8.1]{.toc-section-number} Data
        Modeling](#data-modeling){#toc-data-modeling}
    -   [[8.2]{.toc-section-number} SQL DB to REST JSON
        API](#sql-db-to-rest-json-api){#toc-sql-db-to-rest-json-api}
    -   [[8.3]{.toc-section-number} SQL JSON
        support](#sql-json-support){#toc-sql-json-support}
    -   [[8.4]{.toc-section-number} Data transformation and
        mashups](#data-transformation-and-mashups){#toc-data-transformation-and-mashups}
    -   [[8.5]{.toc-section-number} Other approaches and
        background](#other-approaches-and-background){#toc-other-approaches-and-background}

# The Book of Newt

# Preface, the Newt Experiment, February 2024

The Newt Project was started in 2023. The first prototype was
implemented and culminated with a
[presentation](https://caltechlibrary.github.io/newt/presentation/) to
the SoCal Code4Lib group. The lessons learned from the presentation
included

-   Postgres+PostgREST was exciting (expected)
-   SQL was a problem for many people (expected)
-   Pandoc didn't ring a bell with the web developers (a surprise)
-   A service that combined only Postgres+PostgREST with Pandoc server
    wasn't compelling

What was nice about presenting was the excellent and supportive feed
back. It is all positive because it set me on a better path. It helped
me understand why my colleagues were politely ambivolent to my Newt. I
wasn't about the amphibian.

## Summer, Fall 2023

By the time I had wrapped up coding the first prototype I was
successfully generating Postgres SQL and PostgREST configuration files.
I was troubled by my horrible choices in the YAML syntax I had evolved.
You had to know too much about SQL and specifically Postgres. My YAML
syntax was a complete mess.

I continued to experiment with Pandoc server. While it is capable and
cleverly simple it wasn't particularly fun to work with. If you send
invalid JSON you get an error message. This is very helpful. When you
sent it a valid JSON POST and Pandoc returns an empty request body it is
really hard to see where you went wrong. I found that a completely
frustrating experience. I could not ask others to live with that.

My YAML is horrible and Pandoc is out, what next?

## Fall 2023 thru January 2024

I kept purcolating on Newt. Was I building a solution to a problem only
I had? I dreaded writing yet another metadata curation app. I had three
scheduled to be written and that waa three too many. I was left
pondering. Is it me or do web applications really need to be complex?
Surely there is a simpler way. Or di that vanish with the late
Prof. Niklaus Wirth?

I needed a clean slate, a new prototype. What should that look like?

## February 2024, a second prototype

My motivating insight came from a cool demo by my colleague [Tommy
Keswick](https://library.caltech.edu/about/directory) gave for a project
he did with Caltech Archives. He showed the GitHub YAML issue syntax[^1]
needed to use GitHub issues to trigger GitHub workflows. When I saw how
the syntax setup the data needed for the actions it struck a chord. That
is what I needed for describing my data models. I don't know who
designed the GitHub YAML issue template syntax but they did a really
good job. It feels like you're talking about the elements in a web form.
That means as a web developer I'm building on elements and concepts that
are second nature. I reread the documenetation at
[MDN](https://developer.mozilla.org/en-US/docs/Learn/Forms/HTML5_input_types)
for HTML5 form element types. The light buld turned on. I could infer
directly using the GitHub YAML issue template syntax (abbr: GHYTS) the
SQL I needed in Postgres or SQLite3 to hold the data. GHYTS can indeed
be used as a modeling language targeting an SQL database. Of course it
is designed to model an HTML form so that problem gets solved too.

My sage colleage [Mike
Hucka](https://library.caltech.edu/about/directory) pointed to me the
problems of inventing new languages, even new domain spefic languages.
When you invent a new language or new DSL you are increasing cognative
load on your the developer. By using GHYTS I side step that problem. Our
developers already know how to create the issue temaplates on GitHub. I
am repurposing that knowledge for code generation to build a custom
application. Time to start thinking through the specifics of the new
prototype.

## The second prototype

Newt tools need configuration. YAML is the language used to express that
configuration. Newt's second prototype uses CITATION.cff to describe the
application mentadata for the Newt based project. Actually all it needs
is to point to the CITATION.cff[^2] and the metadata about the
application you're creating is complete. The horrible part of the first
YAML for the first prototype was the modeling syntax. GHYTS solves that.
What's left to learn is much easier to describe. It's how to route the
data and how to provide templates support. But I am getting ahead of my
story so I'll leave that for the next section.

# What is Newt? Why does it matter?

Newt is a project of the [Caltech Library Digital Library Development
Group](https://caltechlibrary.github.io/newt). Newt is an on going
experiment in rapid application development for libraries, archives and
museums (abbr:
`<abbr title="libraries, archives and museums abbrevation">`{=html}LAS`</abbr>`{=html}).
How does Newt achieve "rapid" developemnt?

Newt targets metadata curation applications. The LAS community swims in
metadata and always needs tooling to better manage it.

Newt's approach to creating a metadata curation application is based on
the following characteristics.

-   preference for "off the shelf" over writing new code
-   configuration over writing new code
-   data modeling and management placed squarely in your database
    management system
-   database management systems must provide HTTP accessible JSON API
-   data transformation (if needed) via a web service(s) that consumes
    JSON
-   alignment of services avoiding overlapping capabilities
-   code generation where appropriate
-   write less code, keep it simple and focused

These characteristics and the narrow focus of metadata curation
applications suggest a service oriented architecture. If we assume one
or more web services can easily be combined. Delivering the results of a
web browser request through combining services suggests a data
pipeline[^3]. A collection of pipelines maps to the operations needed by
web applications. If you know these ahead of time you should be able to
calculate the resulting application.

Less code to write, a common proven architecture, a simpler solution.

Newt knitts those characteristics together using YAML description. The
YAML describes the data models you want to work with. It also provides
the configuration for the Newt applicatoins. It even includes
configuration a light weight template engine and mapping the templates
needed for the objects you're modeling.

The first step in a Newt based project is to write the YAML that
describes the above, generate the code that can run in Postgres and
configure PostgREST as well as generate basic Mustache templates to
bootstrap your application.

In short, write less code to get a basic CRUD-L[^4] application.

## What's "off the shelf"?

Software is composed of bits, it doesn't site on myself. Off the self is
software you don't need to write. Software that isn't going to cost you
or organization a license fee or subscription charge. The software
should be proven, with a good community of developer to help answer
questions and a good track recordfor the rest of us humans. Good
examples of "off the shift" are

-   Firefox web browser
-   Postgres, PostgREST
-   Apache Solr
-   Apache 2 web server
-   Shibboleth for single sign-on

You don't write these things you use them. You may have them already
installed. They run on macOS, Windows, Linux and the BSDs. Newt doesn't
need to provide any of those services it can assume them and take
advantage of that.

## What's my point about where data modeling and management take place? Why is that a Newt thing?

Database management systems like Postgres are extremely capable of
managing data. The SQL language also is good at describing the data you
want to manage. We really should be taking full advantage of that. We
certainly don't want to write code that replicates that ability. Yet we
often do write that code and it complicates our programs terribly[^5].
Why? Here's my guesses.

1.  SQL is different from general purpose programming languages
2.  In the web development genre we're been tought to techniques to cope
    with that. E.g. embedding SQL or using an ORM[^6]

The enteria behind the way we've always done things is hard to resist. I
suggest we do resist it. The ORM and the other libraries and frameworks
that are used in many programming languages like Perl, PHP, Python,
Ruby, Java, Go, etc. are a symtom of the challenge of shifting cognative
gears between thinking about "how" a computer should perform at task
versus "what" you want to computer to provide you. When you try to
combine the "how" and "what" closely by creating embedded SQL statements
things get messy really fast. An ORM might mask this but it's there.

A Newt proposition is why accept the cognative clash? It made sense in
1993 but doesn't in 2024. Let the off the shelf software to do that.
Enter [Postgres](https://postgres.org) and
[PostgREST](https://postgrest.org). If you think of those two pieces of
software as a single service you have SQL managing your data but the
rest of your application never touches SQL. It sees JSON and makes
request via URLs. Look Ma, No SQL. If Postgres isn't you cup of tea
you'll finsh software that does the same for MySQL database and even
some projects that try to do that for SQLite 3. Setup the database with
tables, views and functions and you get a JSON API out of the box! It's
liberating.

## "Write SQL? I need to know that? Ugh."

SQL is a problem because it is so different from general purpose
langauges many developers use day to day. One approach is suggest a good
book and ask your fellow developer to learn it. That's a "no go" for
lots of reason for many people. I think we have an escape hatch. Newt
comes with a code generator, include generating SQL code. Generate code
doesn't remove the SQL problem but does mitigate it some. Many people
find it easier to approach SQL if some is already written. It's a half
way approach.

I believe it is worth learning SQL. I don't think you need to learn it
all at once. I think it is useful to learn as you need it. Having
existing SQL code gives you an opportunity to learns a little at the
time. Modify it and see what happens. Newt is happy if you want to
regenerate it. Take your time wrapping your head around SQL's "what you
want" orientation versus the "how to do" orientation of general purpose
programming languages. Meanwhile let Newt bootstrap your
Postgres+PostgREST service with generated SQL and PostgREST
configuration.

## Wait, JSON is great but I need a web page!

JSON is a simple and proven way to express structured data. But it is
not friendly to humans who expect to see an accessible web page. There
are two ways to go right now on the web. The first way is to have a
static web page that a web browser can use to then access the JSON data.
The web browser then updates the web content to look like the desired
results. The has a huge hidden cost. Well actually many hidden costs but
let me focus on one. The way I get JSON into a web page is to **rely**
on JavaScript running in the web browser (or more recently a WASM module
but that a whole other can of worms). The browser used the JavaScript to
the content down, to update the web page contents rendering something
that is hopefully what you asked for. For programmers (I count myself in
this group at times) it's catnip. Not only am I telling my server what
to do with my program I'll telling that computer in your pocket what to
do too! Low the power of the WWW. So that's the problem for the humans?

At least in North America and much of the Central Pacific wireless
network access and can be quiet depressing. In 1993 people used modems
and desktop computers, in 2024 most of the web is experienced through
your "smart phone". I live in Southern California. It has allot of rich
people with latest in clever devices. Our our model networks are
horrible. Even with a fancy phone. When you send JavaScript down the
virtual wire you'll requiring significant network bandwidth just getting
the JavaScript and the content it will assemble to your phone. Then your
phone has to run the JavaScript. It is a terriby idea[^7]. Yes, I know
that there are companies make bucks off this moddel[^8] but it doesn't
make it a good idea[^9]. It is a rotten experience for us humans.

What's the better approach? Glad you asked, it's to have the web browser
do only a couple network requests and get HTML, CSS and perhaps a couple
media assets and bingo the human can read your content or interact with
your web service. Yes, I know this sounds all reto but it works. It
works really well. We used to have a fancy term for it, "Progressive
Enhancement". If you're interested I recommend reading [Christian
Heilmann's blog](https://christianheilmann.com/). I wish my bird was
still as red as his.

Newt applications can send JavaScript down the virtual wire but they
don't need to. Without enhancement a Newt application should work even
using a text browser like [Lynx](https://lynx.browser.org/). A Newt
application starts out handing back HTML by using a template engine to
transform our JSON into a web page. To encourage this approach Newt
provides Newt Mustache. A simple highly configuration template engine.
It expects to receive a JSON object, apply a Mustache template and hand
back the results. The Mustache template langauge[^10] is described at
their website <https://mustache.github.io/>.

## Where does that leave us developers?

The code Newt generates is for a minimal metadata curation application.
It is providing you with a bootstrap. It off the ground with something
that is sort of working. Albiet not exactually what you want. At best it
probably ugly and on cordinated, like a new born Aardvark. Like the
Aardvark your app will become adorable if it survives.

You're likely going to spend your time in one of two areas. The back end
meaning you mucking about with your database or the front end where you
are mocking about with templates and static HTML, CSS, JavaScript and
page assets. My hope is that you are freed up to spend time on the front
end. Why? Because that is the side of things that us humans experience.
While I've always considered myself a "back end" developer, I can
exclaim with pride, "Who now cares about the back end? Unless its
broken!". Mostly the brack end is doing data management and often boils
down to our CRUD-L operations with perhaps some variations on the lists
we produced. If you need a search system I recommend Solr. It is off the
shelf too and plays nice with the Newt appraoch.

If you're lucky enough to have an audience are data analysts then you'll
want to address them too. One approach is to dump snapshort of your
Postgres database (or subsets) and drop into your static content
directory. Then just link to them. You can also reach for JSON API
provided by Postgres+PostgREST and a Mustache template to do it that way
too. Trade offs either way, but that software engineering generally.
Picking the best two of three.

## Newt, three tools for your toolbox

Newt provides a code generator, a Mustache template engine and a data
router. With these three "off the shelf" tools you can take advantage of
those other "off the shelf" tools like Postgres+PostgREST, Solr and
friends. All three Newt programs use the same YAML file to get their
jobs done. The learning curve is primarly picking up Newt's YAML syntax.
But that is for another chapter.

> I know I stacked that software someplace ...

### Off the shelf

-   [Postgres](https://postgres) + [PostgREST](https://postgrest.org)
    (data modeling and management)
-   [Solr](https://solr.apache.org) full text search engine (search and
    discovery)
-   [Apache 2](https://httpd.apache.org) +
    [Shibboleth](https://www.shibboleth.net/) (controlling access)

### [Newt tools](https://github.com/caltechlibrary/newt)

-   `newtrouter` is a stateless web service (a.k.a. micro service) that
    routes a web request through a data pipeline built from other web
    services
-   `newtgenerator` is a code generator that can takes a set of data
    models described in YAML and generates SQL and Mustache templates
-   `newtmustache` is a simple stateless template engine inspired by
    Pandoc server that supports the Mustache template language

These six programs cover allot of ground. They provide the core
functionality of many systems built for libraries, archives and museam.

If you're inclined to readup on Postgres then I recommend [The Art of
PostgreSQL](https://theartofpostgresql.com/), by Dimitri Fontaine. Gave
me allot to think about.

## Newt design choices

Newt needs to help make it trivial to quickly generate metadata curation
applications. It needs to leverage existing software and knowledge.

I know Newt needs to play well with the following software.

-   [Postgres](https://postgres.org) +
    [PostgREST](https://postgrest.org)
-   [Solr](https://solr.apache.org)

JSON data is the way I want to work with structured data coming out of
these system. I also want to minize any additional middleware needed to
assemble a base metadata curation app featuring the basic CRUD-L[^11]
operatins. I want assembling my web application to be conceptually as
easy as working with Unix pipels in the shell. An example would be to
retrieve data from a JSON API (e.g. the software above) and then run it
through a simple minimal template engine. I can see cases, based on my
prior experiece with Pandoc server, there way be more than two stage in
the processing sequence. I need a data workflow, a pipeline. We do this
all the time in web applications but often it is browser side,
e.g. CrossRef and ORCID data pulled in via JavaScript. I don't want to
rely on the web browser so I need a pipeline behind the web server. I
need a data router.

Since Pandoc server isn't my ticket I started to look at other simple
template languages. Mustache is one I remember from both browser side
implementations and server side when I was working with PHP application
at USC. Mustache is alive and kicking and is supported by most of the
programming languages I'm likely to run into in libraries, archives and
museum organizations. There a good, maintained, Mustache package for Go
which is my implementation language for Newt. I can take the lessons
learned using Pandoc server and apply them to a template engine as part
of the Newt project. It will be separate from the data router so it can
be easily swapped out if someone doesn't like Mustache. It needs to take
the JSON object or array from a data source and transform it. So the
configuration for the Mustache template engine needs to be part of the
YAML used in a Newt project.

This all implies three tools for Newt's second prototype

-   data router supporting the pipeline concept
-   template engine that plays well with the data router
-   a code generator that understands Newt's YAML syntax

This is buildable.

It lets us leverage systems like Postgres and PostgREST. Newt can lower
the SQL burden through code generation while keeping the data modeling
and management in the database (a core Newt philosophy). Using a single
YAML for data modeling means one sorce for both web form generation and
for SQL generation. Here's what I need at the top level of my YAML

-   application metatadata and runtime settings (e.g. port numbers,
    directories, OS environment to leverage)
-   models to describe the web forms and SQL needed from the database
-   routes describing what the data router will listen for
    -   pipelines identifying the services and contact requirements to
        process the data
-   templates a mapping of template requests to template processing the
    received JSON data

Developing the second prototype will help fill in the details guided by
these observations.

# Newt YAML syntax

Newt programs are configured in YAML files. Newt programs may focus on
some properties and ignore others. The interpretation is specific to the
program.

These are the top level properties in YAML files.

application
:   (optional: newtrouter, newtgenerator, newtmustache) holds the run
    time configuration used by the Newt web service and metadata about
    the application you're creating.

models
:   (optional: newtrouter, newtgenerator) This holds the description of
    the data models in your application. Each model uses the [GitHub
    YAML issue template
    syntax](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/configuring-issue-templates-for-your-repository#creating-issue-forms)
    (abbr: GHYTS)

routes
:   (optional: newtrouter, newtgenerator) This holds the routes for the
    data pipeline (e.g. JSON API and template engine sequence)

templates
:   (optional: newtmustache)

## the "application" property

The application properties are optional. Some maybe set via command
line. The the Newt application manual pages.

namespace
:   (optional: newtgenerator) uses this in the SQL generated for setting
    up Postgres+PostgREST

port
:   (optional: newtrouter, newtmustache) default is This port number the
    Newt web services uses to listen for request on localhost

htdocs
:   (optional: newtrouter only) Directory that holds your application's
    static content

metadata
:   (optional) This holds metadata about your application using the
    [CITATION.cff](https://citation-file-format.github.io/) YAML syntax
    under metadata.

environment
:   (optional: newtrouter, newtmustache) this is a list of operating
    system environment that will be available to routes. This is used to
    pass in secrets (e.g. credentials) need in the pipeline

The fifth attribute in application is special. It can be used in instead
of `metadata`. If you maintain a CITATION.cff file you can point to it
to avoid maintaining it in two places. When the Newt router or code
generated is started up it will copy the contents into the `metadata`
property.

citation
:   (optional) This points at an file (e.g. CITATION.cff). It is used to
    populate the metadata property at startup

## the "routes" property

Routes hosts a list of request descriptions and their data pipelines.
This property is only used by Newt router and Newt code generator.

### a route object

`id`
:   (required) This identifies the pipeline. It maybe used in code
    generation. It must conform to variable name rules[^12]

`description`
:   (optional, encouraged) This is a human readable description of what
    you're trying to accomplish in this specific route. It may be used
    in comments or by documentation generators.

`request [METHOD ][PATH]`
:   (required) This is a string that expresses the HTTP method and URL
    path to used to trigger running the data pipeline. If METHOD is not
    provided it will match using just the path. This is probably NOT
    what you want. You can express embedded variables in the PATH
    element. This is done by using single curl braces around a variable
    name. E.g. `GET /items/{item_id}` would make `item_id` available in
    building your service paths in the pipeline. The pattern takes up a
    whole path segment so `/blog/{year}-{month}-{day}` would not work
    but `/blog/{year}/{month}/{day}` would capture the individual
    elements. The Newt router sits closely on top of the Go 1.22 HTTP
    package route handling. For the details on how Go 1.22 and above
    request handlers and patterns form see See
    <https://tip.golang.org/doc/go1.22#enhanced_routing_patterns> and
    <https://pkg.go.dev/net/http#hdr-Patterns> for explanations.

`pipeline`
:   (required) this is a list of URLs to one or more web services
    visible on localhost. The first stage to fail will abort the
    pipeline returning an HTTP error status. If done fail then the
    result of the last stage it returned to the requesting browser.

`debug`
:   (optional) if set to true the `newt` service will log verbose
    results to standard out for this specific pipeline

#### a pipeline object

A pipeline is a list of web services containing a type, URL, method and
content types

`service [METHOD ][URL]`
:   (required) The HTTP method is included in the URL The URL to be used
    to contact the web service, may contain embedded variable references
    drawn from the request path as well as those passed in through
    `.application.environment`. All the elements extracted from the
    elements derived from the request path are passed through strings.
    These are then used to construct a simple key-value object of
    variable names and objects which are then passed through the
    Mustache template representing the target service URL.

`description`
:   (optional, encouraged) This is a description of what this stage of
    the pipe does. It is used when debug is true in the log output and
    in program documentation.

`timeout`
:   (optional) Set the timeout in seconds for receiving a response from
    the web server. Remember the time spent at each stage is the
    cumulative time your browser is waiting for a response. For this
    reason you may want to set the timeout to a small number.

## the "models" property

Models holds a list of individual models used by our data pipelines. The
models are by Newt code generator and the Newt router. Models defines a
superset of the GitHub YAML issue template syntax (abbr: GHYTS).

### a model object

The model object is based largely on GitHub YAML issue template syntax
with a couple extra properties that are Newt enhancements.

id
:   (required, newt specific) this is the name identifying the model. It
    must conform to variable name rules[^13]

```{=html}
<!--
routing
: (optional, newt specific) this holds a list of route ids associated with this model. It is use in code generation, e.g. to populate a web form's action and model
-->
```
The following properties are based on the GitHub YAML issue template
syntax[^14] (abbr: GHYTS)

name
:   (required: GHYTS, optional: newt) Must be unique to use with GitHub
    YAML issue templates[^15]. In Newt it will be used in populating
    comments in generated SQL

description
:   (required: GHYTS, optional: newt) A human description of the model,
    It will appear in the web form and SQL components generated from the
    model

body
:   (required) A a list of input types. Each input type maps to columns
    in SQL, input element in web forms and or HTML elements in read only
    pages

#### a model's input types

This is based on GitHub YAML issue template (abbr: GHYTS) input
types[^16].

id
:   (required) an identifier for the element. Must conform to variable
    name rules[^17]. It is used to SQL as a column name and in web forms
    for the input property.

type
:   (required) Identifies the type of elements (input, textarea,
    markdown, checkbox, dropdown).

attributes
:   (optional) A key-value list that define properties of the element.
    These used in rendering the element in SQL or HTML.

validations
:   (optional, encouraged) A set of key-value pairs setting constraints
    of the element content. E.g. required, regexp properties, validation
    rule provided with certain identifiers (e.g. DOI, ROR, ORCID).

## input types

Both the routes and models may contain input types. The types supported
in Newt are based on the types found in the GHYTS for scheme[^18]. They
include

markdown
:   (models only) markdown request displayed to the user but not
    submitted to the user but not submitted by forms.

textarea
:   (models only) A multi-line text field

input
:   A single line text field. This conforms to value input types in HTML
    5 and can be expressed using the CSS selector notation. E.g.
    `input[type=data]` would be a date type. This would result in a date
    column type in SQL, a date input type in HTML forms and in
    formatting other HTML elements for display.

dropdown
:   A dropdown menu. In SQL this could render as an enumerated type. In
    HTML it would render as a drop down list

checkboxes
:   A checkbox element. In SQL if the checkbox is exclusive (e.g. a
    radio button) then the result is stored in a single column, if
    multiple checks are allowed it is stored as a JSON Array column.

Newt may add additional types in the future.

## Example Newt YAML for router and code generator

``` yaml
application:
  port: 8011
  htdocs: htdocs
  metadata:
    cff-version: 1.2.0
    message: Demo of Newt YAML file
    type: software
    title: Newt a faster way to build metadata curation applications
    abstract: |
      This is a demonstation of a YAML that can generate a simple
      application to manage people and groups
    version: 0.0.0
    status: concept
    authors:
      - family-names: Doiel
        given-names: R. S.
        orcid: "https://orcid.org/0000-0003-0900-6903"
    keywords:
      - demo
      - newt
      - rapid application development
  environment:
    - DB_USER
    - DB_PASSWORD
    - DB_HOST
models:
  - id: people
    name: People Profiles
    description: |
      This models a curated set of profiles of colleagues
    body:
      - id: people_id
        type: input
        attributes:
          label: A unique person id, no spaces, alpha numeric
          placeholder: ex. jane-do-007
        validations:
          required: true
      - id: display_name
        type: input
        attributes:
          label: (optional) A person display name
          placeholder: ex. J. Doe, journalist
      - id: family_name
        type: input
        attributes:
          label: (required) A person's family name or singular when only one name exists
          placeholder: ex. Doe
        validations:
          required: true
      - id: given_name
        type: input
        attributes:
          label: (optional, encouraged) A person's given name
          placeholder: ex. Jane
      - id: orcid
        type: input
        attributes:
          label: (optional) A person's ORCID identifier
          placeholder: ex. 0000-0000-0000-0000
        validations:
          pattern: "[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]"
      - id: ror
        type: input
        attributes:
          label: (optional) A person's ROR identifing their affiliation
      - id: email
        type: "input[type=email]"
        attributes:
          label: (optional) A person public email address
      - id: website
        type: "input[type=url]"
        attributes:
          label: (optional) A person's public website
          placeholder: ex. https://jane.doe.example.org
routes:
  - id: create_person
    description: Create new person profile
    request: POST /person
    pipeline:
      - description: This will generate a new user in the database
        service: POST "https://{{DB_USER}}@{{DB_HOST}}:3000/rpc/people"
        content_type: application/json
      - description: |
          This sends the results of creating a person to the template engine
        service: POST http://localhost:3032/people_update_result.tmpl
        content_type: application/json
  - id: read_person
    desciption: Update a person's profile
    request: "GET /person/{{people.people_id}}"
    pipeline:
      - description: Retrieve a person's profile
        service: POST "https://{{DB_USER}}@{{DB_HOST}}:3000/person/{{people.people_id}}"
        content_type: application/json
      - description: |
          Render a person's profile
        service: POST http:localhost:3032/profile.tmpl
        content_type: application/json
  - id: update_person
    description: Update person's profile
    request: "PUT /person/{{people.people_id}}"
    pipeline:
      - description: This will update a person record in the database
        service: PUT "https://{{DB_USER}}@{{DB_HOST}}:3000/rpc/people"
        content_type: application/json
      - description: |
          This sends the results of updating a person to the template engine
        service: POST http://localhost:3032/people_update_result.tmpl
        content_type: application/json
  - id: delete_person
    description: Remove person's profile
    request: "DELETE /person/{{people.people_id}}"
    pipeline:
      - description: Remove the person for the database
        service: DELETE "https://{{DB_USER}}@{{DB_HOST}}:3000/people/{{people.people_id}}"
        content_type: application/json
      - description: Displace the result of what happed in the removal
        service: POST http://localhost:3032/removed_person.tmpl
  - id: list_people
    description: List people profiles available
    request: GET /people
    pipeline:
      - description: Retrieve a list of all people profiles available
        service: GET https://{{DB_HOST}}@{{DB_HOST}}:3000/people
        content_type: application/json
      - discription: format a browsable people list linking to individual profiles
        service: POST http://localhost:3030/list_people.tmpl
        content_type: applicatin/json
```

## templates property

This property is used by Newt Mustache. It is ignore by Newt router and
code generator.

templates
:   (optional: newtmustache) this holds a list of template objects

### template object model

The template objects are used by Newt Mustache template engine. If
you're not using it you can skip these.

`request [METHOD ][PATH]`
:   (required) This holds the request HTTP method and path. If the HTTP
    method is missing a POST is assumed

`template`
:   (required: newtmustache only) This is the path to the template
    associated with request. NOTE: Pandoc web service does not support
    partial templates. Mustache does support partial templates

`partials`
:   (optional, newtmustache only) A list of paths to partial Mustache
    templates used by `.template`.

`options`
:   (optional, newtmustache only) An object that can be merged in with
    JSON received for processing by your Mustache template

`debug`
:   (optional) this turns on debugging output for this template

Example of newtmustache YAML:

``` yaml
application:
    port: 3032
templates:
  - request: GET /hello/{name}
    template: testdata/simple.mustache
  - request: GET /hello
    template: testdata/simple.mustache
    options:
      name: Universe
  - request: GET /hi/{name}
    template: testdata/hithere.html
    partials:
      - testdata/name.mustache
    debug: true
  - request: GET /hi
    template: testdata/hithere.html
    partials:
      - testdata/name.mustache
    options:
      name: Universe
```

# Newt Generator Explained

Newt comes with a code generator called `newtgenerator`. It uses the
Newt YAML file to render Postgres SQL, PostgREST configuration and
Mustache templates suitable to bookstrap your Newt based project. How
does it do this? What are the assumptions?

The Newt code generator works primarily with the "models" property in
your Newt YAML file.

Our minimal useful application should be able to do five things -
Create, Read, Update, Delete and List (abbr: CRUD-L). It needs to offer
these actions for each model defined in our project. The second Newt
prototype assumes each model is indepent. If you need to combine models
(not unsual) then you will need to enhance the generated SQL. For now
let us focus on the basics.

In the Newt philosophy we model our data and manage our data in the
database. For us this means Postgres. By combining Postgres and
PostgREST we create a JSON API based on what we've modelled and managed
in our database. Our middleware, the stages in our pipelines, do not
need to know anything about SQL. They don't ever touch it. This gives us
a clean break from SQL and the rest of our system. It avoids the burden
of doing the cognitive shifts when implementing middleware. This is a
result of Postgrs plus PostgREST automagically giving a JSON API. Clever
stuff really.

## Bringing up Postgres+PostgREST

Postgres and PostgREST provide us with a JSON API that can be used in a
Newt router pipeline. Following the Newt philosophy data modeling and
management happen in the database. That is done via the SQL langauge
which the Newt generator can produce. After the database is setup in
Postgres with the scheme, tabes, functions and views needed for basic
CRUD-L operations you can generate a compatible PostgREST configuration
file. Here's the steps I follow after creating my Newt YAML file.

1.  Create the Postgres database and schema[^19] if they don't exist
2.  For each Newt model create a table in the database/scheme
3.  For each Newt model create a function for each of the CRUD[^20]
    options[^21]
4.  For each Newt model create a SQL view to handle the list view of the
    CRUD-L
5.  Generate a PostgREST configuration file.

NOTE: The generator always generates it's output to standard out.
Errors, if they occur are sent to standard error. This makes it easy to
script the Newt generator in shell or in a Makefile.

### Generating some SQL

Let's turn our attention to setting up Postgres database. Several things
need to be accomplished and the Newt generator will generate the SQL
needed for all of them. It will even generate the SQL to make managing
your content via PostgREST JSON API easier.

OK, now that we know what we're generating lets generate this SQL from
our app.yaml. We'll save the SQL output in a file called "app.sql"

``` shell
newtgenerator postgres app.yaml >app.sql
```

In the statement above "postgres" tells the generator target the
Postgres database SQL dialect.

Take a look at "app.sql". Notice that the SQL if the code contains
comments. If you've included the `.description` attributes in our models
then the SQL also includes them. This means you can continue improving
the SQL manually if you like. Another approach would be to create a
second SQL file with your modifications and use SQL `alter` or `drop`
statements where you need to augment or replace the Newt generated
parts.

This is a text that explains Postgres database so I am going to keep
this short. If you want to read about Postgres to better understand it I
recommend the [Postgres website](https://postgres.org) or the book [The
Art of
Postgres](https://theartofpostgresql.com/ "website for the book"). is
heavely commented, those your

What was generated? The specifics will depend on what was in the
"app.yaml" file, specifically the models property. Do you see comments
you the generate SQL file? These were form from the `.description`
properties found in the YAML. By taking the time to add description
properties are making the generated code more human freindly. While
these are optional properties they are really handy when you evolve the
code generator or come back to a project later.

The first section of the SQL generate should have statements to create
the database if it doesn't exist and to create a name space if it
doesn't exist. Next you should see SQL for creating tables, one for each
model. Then you should see the creation of views. The SQL views fill the
role of the "L" in the acronym CRUD-L.

The nice things about having views is it make it easy to see what our
application will see when accessing the JSON API to list objects based
on our models. This parity becomes very helpful when you are extending
your application or trying to debug why the JSON API isn't giving you
what you expected. A view is accessed through the SQL `select`
statement.

After a model's view Newt will generate functions for create, read,
update and delete. Why do this? First it allows us a consistent
abstraction should the underlieing data model change in the future.
Second you can enhance the function to perform additional validation if
needed. Example, let's say you've imported the Python library idutils
into Postgres. You can use that library to validate identifiers like
DOI, ORCID, ROR, arXiv, ISNI, etc. Since the operation is a "function"
rather than an SQL prodecedure you can let the requestor there was a
problem or take actions to mitigate it. A function can keep the
conversation going with the web browser at the other end of the wire.

You might be wondering about how things are name. The functions and view
names are formed from the model name appened with "\_", followed by the
action. Newt trys to give you a robust bootstrap. One of the ways that
it accomplishes that is by naming things predictably. This improves code
readability for us humans.

After defining the models tables, views and functions the end of the
generated SQL sets up the access needed by PostgREST to provide the JSON
API. Essentially that is a block of permission granting statements. When
you read through these you need to understand them and make sure they
comply with your organization's database practices. Newt is only making
a guess what is OK for you. See the PostgREST website for more
information about how Postgres and PostgREST work together,
<https://postgrest.org/en/stable/>.

### Generate our PostgREST configuration

The generator can generate a PostgREST configuration file,
`postgrest.conf`. This is the first step. Next step is to setup our
Postgres database to model our data, manage our data and to allow
PostgREST to access our database appropriately. This is all done via
SQL. So how to we generate. To generate the `postgrest.conf` you use the
newtgenerator with a YAML file named "app.yaml" you would do the
following.

``` shell
newtgenerator postgrest app.yaml >postgrest.conf
```

In the above command the "postgrest" (note the "t" at the end) tells the
generator what needs to be generated. In this a Postgrest configuration
because that is what PostgREST needs.

## Ready for testing our JSON API

At this point we have defined a predictable JSON API and set of related
URLs we can use in our data pipelines. You should now have three files
in our application directory, two were generated by the Newt generator.

-   app.yaml, you created this by typing it in or copying an example and
    modifying it for your needs
-   postgrest.conf was generated by Newt and holds a configuration used
    when you startup PostgREST
-   app.sql is the SQL file that sets up Postgres to work with our
    application

Now we need to set things in motion. You can load your SQL into Postgres
via the Postgres repl called `psql`. You'll need to have administrative
privilages to run the SQL as it will create a new database for your
application. See
<https://www.postgresql.org/docs/current/sql-createuser.html> about
creating database users for details.

On my computer I have Postgres 16 running and my personal account has
admin privilleges so I use `psql` like this to create the database,
schema, tables, views, functions and permissions needed for the my
application.

``` shell
psql < app.sql
```

If that worked then we can try out starting PostgREST and use
[curl](https://curl.se/) to see if our JSON API works. The exact URL
will be dependent on the database and model names setup in Postgres. The
PostgREST webpage explains the JSON API, see
<https://postgrest.org/en/stable/references/api.html>. If I have created
a model named "people" I can start PostgREST then test of it is
available like this.

    postgrest postgrest.conf &
    curl http://localhost:3000/people

NOTE: This first command starts PostgREST as a back ground process, you
need to kill it if the data base changes table definitions, views or
you've added views, tables or functions. See the PostgREST documentation
for details about when you need to restart PostgREST,
<https://postgrest.org/en/stable/references/schema_cache.html>.

That should return a list of people. Since our database tables aren't
populated it should be an empty list. At this point I suggest getting
comfortable with Postgres and PostgREST. While the generator creates the
SQL and configuration needed that doesn't substitute for understanding
how it is working and actually how you might want to use it!

## Generating some templates for a web UI

Newt comes with a light weight template engine called Newt Mustache. It
implements Mustache templates. The Newt generator knows how to generate
those. Newt can generate a template for each of our CRUD-L operations
for each model. To know what template you want to generate you need to
tell the generator you want to generate mustache templates, which model
you generating it for and what action the template will model. As
mentioned previously Newt generator writes the generated code to
standard output and errors if encountered to standard error. Like with
generating SQL and a configuration file this allows for flexibility in
scripting via shell or a Makefile. Here's an example of the commands I
type to create the templates for our people model.

``` shell
newtgenerator mustache people create >people_create.tmpl
newtgenerator mustache people read >people_read.tmpl
newtgenerator mustache people update >people_update.tmpl
newtgenerator mustache people delete >people_delete.tmpl
newtgenerator mustache people list >people_list.tmpl
```

If you examine the resulting templates you'll notice that create, update
and delete include webforms nad use the model types you describe. On the
other hand the templates for read and list do not include webforms just
some standard markup elements. I expect you'll want to enhance these to
meet you applications need but they should function well enough to test
your data pipelines and debug them.

I usually get the back end setup and tested before moving to make the
application pretty and enhancing the browser experience.

There are two approaches to testing your templates. One is to use them
as the last stage of your JSON API. Another is to configure the
templates in your YAML and run Newt Mustache service and use curl. In in
the example below I'm assuming you've mocked up a person record in a
JSON file called person.json. We can then test come of our templates to
see how they fit the bill. I'm assuming you've setup Newt Mustache to
provide templates based on the names of the templates. I am also
assuming Newt Mustache is running on port 8040 in this example.

``` shell
newtmustache app.yaml &
curl --data '@people.json' http://localhost:8040/people_read.html
```

You should get back an HTML page with the content of "person.json" in
it. If so that template is working. The create template should return an
empty web form as it is used to "create" new model instances and a model
idea isn't available. If the "person.json" JSON included an object id
then you should see it in the update form as a hidden field. Update
should not enable creating new objects in must cases. Similarly if you
mockup a list of people in a JSON file called "people.json" then you
should be able to test the list template too.

I generally will work directly with mockup JSON files and Newt Mustache
to get the web interface I want for my application.

It is important to remember that Newt Mustache reads the templates in at
program start. If you revise your templates you need to **restart** Newt
Mustache. In this way Newt Mustache and Newt Router are like PostgREST
server that you will need to be comfortable stoping and start the
services as you continue your development.

# Newt Mustache explained

## Overview

Newt Mustache is a simple, light weight template engine supporting the
use of Mustache templates. If you send a JSON object to a configured
Newt Mustache template engine it can run that object through Mustache
and hand back a result. This usually means taking a JSON object
(e.g. JSON object delivered via PostgREST) and turning that object into
web page. That's the type of templates you get from using Newt
Generator. Newt Mustache itself just cares about the JSON object it
recieves and applying the template configured for the request. Normally
the JSON object is sent to Newt Mustache as a HTTP POST action.

## Testing your templates with Newt Mustache

Assuming you have Newt Mustache configuration in your "app.yaml" file
(i.e. there is a templates property completed there) and it points to
valid Mustache templates then Newt Mustache is ready to test. I usually
use [curl](https://curl.se) to do my initial testing. In the example I'm
assuming you have a JSON file named "person.json" that contains JSON
markup that fits the templated called "people_read.tmpl" which is going
to be accessible with the URL `http://localhost:8040/people_read.tmpl`.
I first start Newt Mustahe as a back ground process then use curl to
check the template result.

``` shell
newtmustache app.yaml &
curl --data '@person.json' http://localhost:8040/people_read.tmpl
```

I should now see HTML markup with the "person.json" content rendered
into the template. That's all that Newt Mustache does. It reads in a
bunch of templates and requests paths at start, then renders them based
on the JSON objects you send them.

It is important to remember that changes you make in your templates will
NOT show up in the service until you restart Newt Mustache. This is
similar to the constraints regarding Newt Router and PostgREST.

### Why Newt Mustache?

The Newt Mustache template engine came about due to three concerns I
encountered with Pandoc web server.

-   Pandoc server usage is hard to debug
-   Many people don't know or like Pandoc's template language
-   Pandoc server does NOT support partial templates

When building a web application using a template system it is very
useful to minimize the number of templates you need to know and work
with. One way this is done is to build smaller partial templates that
handle specific content elements or types. E.g. a bibliographic citation
used on a book review website.

When you run Pandoc from the command line this is readily supported.
Unfortunately Pandoc 3.1 server doesn't support this. The server only
knows about the full template you provide as part of the JSON POST.
Mistakenly misconfigure the JSON you post to Pandoc and it will happly
give you nothing back. Not an error message, not a peep. Technically
it's giving back what you requested but it's a pain to fiddle with JSON
to get enough of a response to diagnose the problem. That was a show
stopper for Newt's second prototype. Time to switch horses. Pandoc
server inspired Newt's Mustache template service. Newt's template
service is configured from YAML file. Settings are clearer. You can also
turn on debugging for a specific template you have concerns about. Like
Pandoc server Newt' Mustache template engine is stateless. You can run
as many as you like as long as you have an available port.

Newt Mustache designd for use in Newt's data pipeline.

### Why Mustache templates?

-   Mustache is a widely support template language with support include
    Python, JavaScript, and Go (languages used regularly at Caltech
    Library)
-   Since a Go package provides Mustache template I only need to write a
    light weight web service to wrap it
-   Since I am writing the service I can keep the requires to a minimum,
    i.e. [Use Newt's YAML file
    syntax](newtmustache.1.md#newt_config_file).

See [newtmustache](newtmustache.1.md) manual page for details.

# Newt router explained

## Overview

In the first Newt prototype supported a two stage pipeline for routing
request through return a web page. It supported either
Postgres+PostgREST through Pandoc web service or JSON API like Solr
through Pandoc web service round trip. With the second prototype the
Newt router has been generalized. Rather than two stages the second
prototype implements a pipeline. This allows for several services to be
tied together each sending a request to the next. This allows the web
services to be more focused in much the same way that Unix programs can
be chained together to form pipelines. Using a route selector the
generalized pipeline become steps indicated by a list of HTTP methods,
URLs and content types. The YAML notation used has been significantly
changed to support this generalization. Let's focus on the individual
route setup[^22].

It is easy to start with a specific example then show it would be
notated.

## A simple example

Let's say we have a database of music albums and reviews. Each album
includes a rating of "interesting". The range is a zero (uninteresting)
to five star rating (most interesting). Previously we've modeled this in
our Postgres database using a `view`. How do we create a page that lists
albums in descending order of interest? Since we're building with Newt
we can assume there is a template to list albums available. That using
that template will be the "last stage" in our pipeline. We need to feed
the view into that template. The `view` statement is implemented in SQL
in Postgres. That is exposed as a JSON API by PostgREST. That's our
first stage, a JSON data source.

How do you representing a route with two stages?

``` yaml
routes:
    - id: interesting_album_view
      request: GET /interesting_albums_list
      pipeline:
         - description: Contact PostgREST and get back the intersting album list
           service: GET http://localhost:3000/rpc/album_view
           content_type: application/json
           timeout: 15
         - description: |
             Take the results from PostgREST and run them through 
             the newtmustache using the template "ablum_list_view.tmpl"
           service: POST http://localhost:3032/album_list_view.tmpl
           content_type: application/json
           timeout: 10
      debug: true
```

What is being described? First we have routes defined in our
application. Our route is `interesting_album_view`. When a web browser
contacts Newt via a GET at the designated path it triggers our pipeline
property to start processing the request. In this case it is a two stage
pipeline.

The first step retrieves the JSON data (i.e. the content is fetched from
PostgREST). This is expressed as an HTTP method and URL. There is a
content type that will be used when contacting the URL. You can also
include a timeout value, in this case we're willing to wait up to 15
seconds.

The second stage takes the output from the first and sends it through
our template engine. Like the first there is a description for us
humans, a service property indicating how and what URL we contact. There
is a content type and timeout just like before. The output of the first
stage is going to be sent as a port to the Mustache template engine as
JSON data. The Mustache template engine returns the type based on the
template, in this case HTML.

There is a last property in our router description. `debug: true`. This
property will cause the router to display more logging output. We want
this to debug our application. It is very verbose. In a production
setting this would be skipped or set to `false`.

When the Newt router starts up it reads the YAML file and sets up to run
each pipeline associated with specific request methods and paths. Those
settings don't change until you restart the router. The router only
reads the file once a startup. That's important to keep in mind. The
router only interacts with "localhost". It listens for requests on a
port of your choosing. It then can run a pipeline of other web services
also running on "localhost".

This is where having the descriptions in the route definition is handy.
It is easy to forgot which services are running on which ports. Both are
URLs running as "localhost". In this specific case our PostgREST service
is available on port 3000 and our Newt Mustache template engine is
available on 3032. While the description element is optional it is what
keep port mapping a little more transparent. This is an area Newt could
improve in the future but the reason for using a URL is that Newt
doesn't need to know what each stage actually is. It just knows I
contact this one service and take the output and send it to the next
service and all stages of the pipeline are complete or there is an error
reported in the pipeline. The result has handed back to the web browser.

## Changes from the first prototype to the second.

-   routes include a pipeline rather than fixed stages
-   `newt` was replaced by `newtrouter`. It does less. It just routes
    data now. It does more, you can have any number of stages in our
    data pipeline now. It doesn't know how to package things.
-   `newtmustache` has replaced Pandoc web service as the Newt template
    engine of choice. Mustache is a popular templating language and well
    supported in numerious programming languages. It provided easier to
    debug issue than working with Pandoc server. `newtmustache` does
    require of configuration.
-   each pipeline stage has its own timeout setting

While there isn't a fixed limit to the number of stages in a pipeline
you will want to keep the number limited. While contacting a web service
on localhost is generally very quick the time to round trip the
communication still accumulates. As a result it is recommend to stick to
less than four stages and to explicitly set the timeout values based on
your understanding of performance requirements. If a stage times out the
a response will be returned as an HTTP error.

## Misc

If a requested route is not defined in our YAML by then the router will
look matching static files. If that fails an HTTP 404 is returned. For a
request route to match our YAML the router compares HTTP method, path
and content type. If any of these don't match then the route is not
considered a match and will return an appropriate HTTP status and code.

The router uses are defined in the request property. The HTTP method and
path indicate what can trigger the pipeline being run.

The Newt router will only support HTTP and run on localhost. This is the
same approach taken by Pandoc server. It minimize configuration and also
discourages use as a front end service (which is beyond the scope of
this project).

This prototype does not support file uploads. In theory you could
implement a pipe line stage to handle that but again that is beyond the
scope of this project. You can try clever techniques browser side and
push objects into Postgres via PostgREST but again that is beyond the
scope of this project. I don't recommend that. If you need file upload
support Newt project isn't your solution yet.

# Reference Material

These are links to prior art, related efforts and resources for
consideration in how Newt evolves as a prototype.

## Data Modeling

-   [Syntax for GitHub form
    schema](https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema)
-   [A React component for building Web forms from JSON
    Schema](https://github.com/rjsf-team/react-jsonschema-form)
-   [JSON Forms](https://jsonforms.io/docs/)
    -   pretty heavy weight in the deployed results (i.e. renders to
        either React or Angular)
    -   Uses JSON rather than YAML or TOML to describe forms which can
        be construed as less human friendly
-   [YAQL](https://yaql.readthedocs.io/en/latest/getting_started.html),
    a YAML like query language that maps to SQL
-   [pg_query](https://github.com/pganalyze/pg_query), a Ruby library to
    parse SQL and normalize into a data structure
-   [htmx](https://htmx.org/), a Web Component like markup implementing
    the wrapping layer between JSON API and HTML structured markup
-   [Yafowil, yet another form widget
    library](http://docs.yafowil.info/)
    -   [Describe YAFOWIL forms with
        YAML](https://yafowil.readthedocs.io/en/latest/yaml.html)
-   [MDN Web
    Components](https://developer.mozilla.org/en-US/docs/Web/API/Web_Components),
    useful for make specialize input elements, like ORCID, DOI, ArXiv
    identifier entry
-   [MDN HTML5 input
    types](https://developer.mozilla.org/en-US/docs/Learn/Forms/HTML5_input_types)
-   [Generate webforms based on YAML schema with
    pykwalify](https://github.com/cxong/pykwalify-webform)
-   [Flask Open API 3](https://pypi.org/project/flask-openapi3/)

## SQL DB to REST JSON API

-   [PostgREST](https://postgrest.org), build a REST JSON API directly
    on top of Postgres. This is what started my think about Newt
    prototype
-   [MRS](https://dev.mysql.com/doc/dev/mysql-rest-service/latest/),
    MySQL REST Service
-   [sqlite2rest](https://github.com/nside/sqlite2rest), Automatically
    RESTful OpenAPI server from SQLite database
-   [Soul](https://github.com/thevahidal/soul), A SQLite REST and
    realtime server built on NodeJS (wonder if it runs in Deno?)

## SQL JSON support

-   SQL dialects
    -   [database guide](https://database.guide/), documentation for
        various SQL dialects including Postgres and SQLite
-   Postgres
    -   [Postgres JSON functions and
        operators](https://www.postgresql.org/docs/16/functions-json.html)
    -   [Postgres JSON
        tutorial](https://www.postgresqltutorial.com/postgresql-tutorial/postgresql-json/)
-   SQLite 3
    -   SQLite [Overview JSON functions](https://sqlite.org/json1.html)
    -   [JSON improvements in SQLite
        3.38.0](https://tirkarthi.github.io/programming/2022/02/26/sqlite-json-improvements.html)
    -   SQLite [JSON function enhancements
        (2022)](https://sqlite.org/src/doc/json-enhancements/doc/json-enhancements.md)
    -   SQLite cli [docs](https://sqlite.org/cli.html), e.g. "Changing
        output formats" in section 5 covers `.mode json`

## Data transformation and mashups

-   [Pipes](https://www.pipes.digital/docs) is a spiritual successor to
    [Yahoo! Pipes](https://en.wikipedia.org/wiki/Yahoo!_Pipes)

## Other approaches and background

-   [path-to-regexp](https://github.com/pillarjs/path-to-regexp)
-   [URLPattern](https://developer.mozilla.org/en-US/docs/Web/API/URLPattern)
    at MDN
-   [URLPattern](https://developer.chrome.com/articles/urlpattern/) at
    Chrome Developer site
-   [Flask Route
    tutorial](https://pythonbasics.org/flask-tutorial-routes/)
-   [router.js](https://github.com/tildeio/router.js/)
-   [Azure application gateway
    routing](https://learn.microsoft.com/en-us/azure/application-gateway/url-route-overview#pathbasedrouting-rule)
-   [React Router](https://reactrouter.com/en/main/route/route)
-   [Nextjs
    routing](https://nextjs.org/docs/app/building-your-application/routing)
-   [dJango
    routing](https://www.django-rest-framework.org/api-guide/routers/)
-   [HYDRA](https://www.markus-lanthaler.com/hydra/), Hypermedia-Driven
    Web API
-   [HAL](https://stateless.group/hal_specification.html), Hypertext
    Application Language
    -   [JSON-HAL](https://datatracker.ietf.org/doc/html/draft-kelly-json-hal-00)
-   [JSON-LD](https://en.wikipedia.org/wiki/JSON-LD)
-   [Richardson Maturity
    Model](https://en.wikipedia.org/wiki/Richardson_Maturity_Model),
    used to evaluate RESTful-ness in JSON API

[^1]: See
    <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-issue-forms>,
    as viewed Feb 2024

[^2]: See <https://citation-file-format.github.io/>

[^3]: A data pipeline is formed by taking the results from one web
    service and using it as the input to another web service. It is the
    web equivalent of Unix pipes. Prior art: [Yahoo!
    Pipes](https://en.wikipedia.org/wiki/Yahoo!_Pipes)

[^4]: CRUD-L, ancronym for create, read, update, delete and list. These
    form the basic operations on metadata objects.

[^5]: When Ruby on Rails gained fame this habit increased. Now it is
    routing for programming language frameworks to provide tools to
    model data and manage their evolution. I don't think this has
    improved things. I think it is a symtop of a bigger problem.

[^6]: ORM, acronym for "Object Relational Mapper", before JSON column
    support in the Open Source databases this helped. Now not so much.

[^7]: When Ruby on Rails gained fame this habit increased. Now it is
    routing for programming language frameworks to provide tools to
    model data and manage their evolution. I don't think this has
    improved things. I think it is a symtop of a bigger problem.

[^8]: Companies like Google, Amazon, Wix and SquareSpace are happy to do
    business his way. Why? It is because it drives the web users to the
    "native app". Native apps mean captured customers. The open web
    provides options, they don't want you to have options.

[^9]: Don't take my work, read this
    <https://infrequently.org/2024/01/performance-inequality-gap-2024/>

[^10]: Yes, I did just asked you to learn a new language but I like
    Christian Heilmann's and my facial hair. If I'm lucky you already
    know Mustache.

[^11]: CRUD-L is an ancronym for create, read, update, delete and list.
    These are the basic operations available in a database management
    system.

[^12]: variable numbers must start with a letter, may contain numbers
    but not spaces or punctuation except the underscore

[^13]: variable numbers must start with a letter, may contain numbers
    but not spaces or punctuation except the underscore

[^14]: See
    <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-issue-forms>,

[^15]: See
    <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-issue-forms>,

[^16]: See
    <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema>

[^17]: variable numbers must start with a letter, may contain numbers
    but not spaces or punctuation except the underscore

[^18]: See
    <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema>

[^19]: In Postgres a "schema" is really a name space for the tables in a
    database. It should be confused with a specific SQL table
    definition, often refered to using the same term "schema". In
    Postgres a "schema" provides a hook to manage table accesses and
    provide a namespace for referencing them es and provide a namespace
    for referencing them. See
    <https://www.postgresql.org/docs/current/ddl-schemas.html>

[^20]: CRUD is an acronym for "create, read, update and delete", CRUD-L
    stands for "create, read, update, delete and list". CRUD-L are the
    basic operations you perform when managing metadata.

[^21]: Newt creates functions for each of our CRUD operations. It does
    this to allow you a clean way to keep the JSON API constent even
    when you data models might evolve in differently.

[^22]: See [Newt YAML syntax](newt_yaml_sentax.md) for complete
    description of the supported YAML.
