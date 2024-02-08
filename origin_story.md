
# Newt's Origin story

## a small epiphany

Newt came from a small epiphany. I was writing yet another flask app and the de ja vu was intense. It got me thinking about how often I write the same code over and over except for minor variations. Why in the world was I doing this?  That caused me to step back, take stock and see if there was a simpler way.

Modern applications like Invenio RDM are created by aggregating many web services. This is a feature of service oriented architecture also known as micro service architecture. Invenio RDM is much cleaner than the venerable EPrints 3 repository system but it remains highly complex. Was this really needed? Was it an artifact of Zenodo? Zenodo is huge. Huge numbers of concurrent interactions (bots and humans) along with a massive amount of data and objects to manage it suggests complexity if only from scale.  At its core RDM does two things really well. Metadata management and object storage. These are non-trivial at Zenodo scale and RDM does a really good job at managing them. But having migrated from EPrints to RDM recently I've come to question why do we built software systems so complex? My library is unlikely to need an application that scales as large as Zenodo. Our needs are far more modest.

I reviewed the other custom metadata management applications I've built, collaborated on or maintained over the last several decades. I feel like they all could be simplified. **It's time to see how much software we can create without writing new software.**

At Caltech Library we build allot of flask apps. Flask compare to many frameworks is a nice balance of flexibility and simple concepts. But looking at our applications I also see that the should be writing less. The middleware I've needs to be scaled back. It should be taking on less responsibility and be more narrowly focus. All the application use a database for content storage (e.g. SQLite3, MySQL and Postgres).  My code spends allot of time modeling data, transforming it on the way to or from the database as well as enforcing access rules. If I use an ORM I'm also doing data modeling in my application. That's shouldn't be part of the middleware at all. This problem is not unique to my code or the code we write as a group in Caltech Library. I saw the same problem when I worked at USC. I see the same problem when using platforms like Drupal and WordPress.

How do we trim down our middleware. Looking at EPrints, RDM, and our custom applications where at Caltech I see several areas that middleware is covering that can be rethought or re-aligned.

1. Data modeling 
2. Data validation
3. Data management
4. User management and permissions regiments
5. Data retrieval 
6. Data transformation for use by the web browser

I realized that all these features often overlap with other software in the systems we build on.  For historical or evolutionary reasons we've been putting those features into middleware. I believe these features often are made redundant depending on the how we use the other systems we build on. 

So here's my propositions

- data modeling (e.g. scheme design) should only happen in the database
- managing data should only happen in the database
- the database needs to be accessible via URL and return JSON content instead of a DB connector used to send SQL down the wire
- user access and permissions should be pushed to the outer layer of a web application (e.g. Apache 2 integrated with Shibboleth)
- if the database produces JSON we still need to turn that into HTML before sending to the browser, this should be done with a simple, stateless template engine
- we need a means of routing data round systems, e.g. taking our JSON API results and sending through the template engine

It is the last bit that our middleware should focus on not the rest. Can data routing become an "off the shelf" component used to build a web application? That is the itch that Newt is trying to scratch.

Proposed components for testing a the simple data router concept

- Postgres database
- PostgREST to turns Postgres into a JSON API managed from Postgres
- Pandoc running as a web service (a very light weight template engine)
- Solr for full text search leveraging its JSON API support

The data router needs to be configurable such that any given URL supported in your application can pass content between the database or search engine and through a template engine handing back the results.

Flask (like many other frameworks) suggest a excellent path foreword.  Flasks lets you assign a "route" to a function. A "route" is expressed with a path like string. It may include replaceable elements that the function can then use to do it's job (e.g. a record id passed in as part of the URL's path).  If you can identify a path using a simple description string and you pair that with a data pipeline sequence then now have a data router. The code you write is configuration for the path and how it should proceed through a pipeline sequence. 

Example 1: contact the JSON API for the search engine, take the results and run them through the template service. 

Example 2: contact the database via JSON API, perform an action, take the results and run them through the template engine.

I was able to illustrate this approach in the initial Newt prototype I showed at SoCal Code4Lib in [July 2023](presentation/). The challenges was my chosen YAML syntax was overly complicated. 

## The 2nd prototype

A couple things fell out of that demonstration and subsequent talks with other developer colleagues.

- Postgres with PostgREST providing JSON APIs is very cool (expected, because yes it is very cool)
- People were not interested in data routing and pipelines (unexpected)

These then were followed by some different flavors of questions and comments. 

Some other comments came in different flavors?

- Do all my data modeling in SQL? 
    - But I don't like SQL! 
    - But I don't know enough SQL!
    - I don't feel comfortable writing SQL!
    - Writing SQL isn't fun!
- I don't need a pipeline the browser can assemble the page!
    - React/Angular/(replace with your favor JS framework) means I don't need templates.
    - Who cares about HTML, JavaScript can deliver everything
    - Do all my modeling in SQL? I don't know SQL well enough!
    - But aren't templates the old thing people used to do?

> NOTE: the exclamations are my editorial not additions only. Everyone I encountered was very supportive and kind

After taking this all in and have a long, long think about 2024 arrived. These are my current conclusions as we move our way into the new year like an atmospheric river.

When people saw Postgres combined with PostgREST and my controversial suggestion to use SQL to model your data the discussion shifted from frameworks to data modeling and what JSON would results. Granted people weren't looking forward to writing or learning more SQL. That's something I've noticed through most of my career. The people who like SQL tended to leave development behind and get paid to manage large scale database deployments.

Shifting the discussion to where you model your data is the important bit. The fact remains that the database engine does a much better job of managing performance and your data then any middleware code you're likely to write regardless of framework or programming language.  I like writing in Go. It is a reasonably performing language but if the database setup up right running SQL queries will beat my processing data outside the database most days of the week.

Data routing shouldn't be either surprising or exciting. In fact it should be boring. It should just be there available to use.  It's like the street I drive down. Structural engineers may get excited about streets and bridges but most people to use them and don't think allot about them. Whatever we use for data routing should be like that if it is successful. In that sense the prototype was successful.

The SQL conundrum is a problem. While database engines like Postgres can integrate non SQL language for functions and procedures you really need to know SQL to take advantage of that. Otherwise we're stuck writing complex middleware that tries to be a database rather than use the database we have.

I see two possible solutions to the SQL problem (and yes, though I like SQL, I see it as a problem for others). First with SQL you don't need to learn the whole thing at once. In fact I think it is a bad thing if you try. SQL is actually made of several domain specific languages.

DDL, data definition language
: This language describes the data structure of the objects you're going to work with, e.g. table rows and columns. CREATE and ALTER are part of SQL DDL.

DQL, data query language
: This is the language you used to retrieve or list objects you stored, e.g. rows from a table. SELECT and WITH are part of the SQL DQL.

DML, data modification language
: This is the language you used to create, update and delete objects. e.g. modify rows or columns in table. INSERT, UPDATE and DELETE are part of the SQL DML.

If you can get comfortable these three domain languages you'll know enough SQL to shift your data modeling focus to the database.

Still if you really hate SQL you'll not be convinced. For those who wish never to touch SQL I offer an alternative. Code generation.

Most applications written to manage metadata have operations that falling into the acronym CRUD-L. Create, Read, Upload, Delete and List. That's a good thing because if I understand the data model well enough I can calculate all the common expressions of those operations. No AI needed :wink:.

One of the nice things about picking Postgres and PostgREST as the initial target of the first stage of my data pipeline is that configuring it for the common case is boiler plate. Likewise generating the table schema, views and some CRUD functions are pretty straight forward. If all you need is CRUD-L operations I think you can avoid SQL in your application because I can generate that for you. I just need a reasonable way of expressing a model in YAML (the configuration language I am currently supporting in Newt).

## A recent insight

A colleague of mine recently demonstrated an innovative use of GitHub issue templates to trigger GitHub workflows for our library's staff. What impressed me about the demonstration (besides my colleague's cool application) was the YAML used to express that an issue is. After some thinking I realized that my first pass a using YAML to describe data models and variable types was all wrong. I knew it was too complex to begin with but I had also been thinking about the problem the wrong way.

GitHub YAML issue templates (abbr: GHYT) describe the data model of issues where elements of the issue are defined by HTML input elements.  When you think about HTML input elements already suggest a mapping into SQL space.  E.g. a plain `input` element is a string. `textarea` is a `text` element. An HTML5 data element expressed `input[type=date]` maps to a SQL date type. This isn't accidental. HTML web forms we designed to make it easier to get data into a database. If a table row is analogous to an object and an input element is analogous to a column I have a clear mapping for generating SQL without resorting to the developer knowing SQL data types. If we take an additional step forward and require JSON column support in our SQL database (MySQL, Postgres and SQLite have for the last decade) and we require the form to use JSON encoding then even elements like checkbox or select lists can be easily to send over the wire.

By extrapolation if you have web components that implement common data clusters used by libraries, archives and museums then the input forms of those elements can package up the parts and send an array of objects down the wire too. 

The second prototype for Newt is focusing and developing an improved YAML syntax. The plan is to used [CITATION.cff](https://citation-file-format.github.io/) for application metadata while using GHYT for models and GHYT input types to describe route variables used to validate a request before sending it through a pipeline. These changes will bring a renewed focus on code generation targeting SQL for Postgres+PostgREST, Mustache template generation for `newtmustache` service as well as HTML blocks and components for web form generation and composition.

The result is a new Newt YAML syntax is expected to less back end and middleware code written by a human. It also means those who don't like SQL can avoid it a little longer.

