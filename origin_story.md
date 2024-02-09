
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

Example 3: contact the search engine, get JSON results, run them through the template engine.

I was able to illustrate this approach in the initial Newt prototype I showed at SoCal Code4Lib in [July 2023](presentation/). The biggest challenge wasn't implemeting the router it was inventing a YAML syntax. My YAML syntax was overly complicated. 

## The 2nd prototype

A couple things fell out of that demonstration and subsequent talks with other developer colleagues.

- Postgres with PostgREST providing JSON APIs is very cool (expected, because yes it is very cool)
- People were not interested in data routing and pipelines (unexpected)

These then were followed by some different flavors of questions and comments. 

- Do all my data modeling in SQL?  ...
    - But I don't like SQL! 
    - But I don't know enough SQL!
    - I don't feel comfortable writing SQL!
    - Writing SQL isn't fun!
- I don't need a pipeline the browser can assemble the page! ...
    - React/Angular/(replace with your favor JS framework) means I don't need templates.
    - Who cares about HTML, JavaScript can deliver everything!![^1]
    - But aren't templates the old thing people used to use?

[^1]: Example: SquareSpace and Wix deliver webpages as JavaScript.  I think this is to hide easy what they really are providing.

> NOTE: the exclamations are my editorial not additions only. Everyone I encountered was very supportive and kind

After taking this all in and having a long think 2024 arrived. These are my current conclusions.

1. When people saw Postgres combined with PostgREST and my controversial suggestion to use SQL to model your data the discussion shifted from frameworks to data modeling and what JSON would results. Granted people weren't looking forward to writing or learning more SQL. That's something I've noticed through most of my career. The people who like SQL tended to leave development behind and get paid more running Oracle deployments.

2. Shifting the discussion to where you model your data is the important bit. The fact remains that the database engine does a much better job of managing performance and your data then any middleware code you're likely to write regardless of framework or programming language.  I like writing in Go. It is a reasonably efficient language. If the database is setup up correctly running SQL queries will beat my processing data outside the database most days of the week. Arriving at SQL to manage data is generally a win.

3. Data routing shouldn't be either surprising or exciting. In fact it should be boring. It should just be there available to use.  It's like the street I drive down. Structural engineers may get excited about streets and bridges but most people who use them and don't think about them. Date routing should be like streets. Almost taken for granted and usually available when you need them.

4. The SQL conundrum is a problem. While database engines like Postgres can suppoert writing functions and procedeures in non-SQL langauges you still need to know and type at least some SQL. I think in part that is why we're still in the complex middleware quagmire. Avoiding messing about with SQL is a feature. On the plus side if your database hands you JSON and you can hand back JSON your middleware doesn't need to know about SQL and you can avoid the cognative overhead embedding SQL into your non-SQL program.

## Addressing the SQL conundrum

> SQL is a problem because embedding SQL causes a cognative shift.[^2]


[^2]: The price of that shift is similar to when you inline CSS or JavaScript in HTML.

I see several possible solutions to the SQL problem. Separate your SQL out and avoid embedding it. This has the advantage of highlighting the ability to write SQL to test your SQL results (e.g. for queries and views you can write the results into a in-memory table then check them with simple SQL SELECT statements). 

Of course you can also try to learn SQL. This is not as bad as some programming language because SQL is itself made of smaller domain specific languages. Three are important to pickup incrementally over time[^3]. SQL is made of several domain specific languages. Here's the ones you are likely need to need in the order you will need them for development purposes.

[^3]: I would go so far as saying attempting to learn a SQL dialect all at once is a bad idea.


DDL, data definition language
: This language describes the data structure of the objects you're going to work with, e.g. table rows and columns. CREATE and ALTER are part of SQL DDL.

DQL, data query language
: This is the language you used to retrieve or list objects you stored, e.g. rows from a table. SELECT and WITH are part of the SQL DQL. You'll spend allot of time using this one when you write your own SQL.

DML, data modification language
: This is the language you used to create, update and delete objects. e.g. modify rows or columns in a table. INSERT, UPDATE and DELETE are part of the SQL DML.

If you can get comfortable these three domain languages you'll know enough SQL to shift your data modeling focus to the database.

**Do you still want to avoid SQL?** For those who wish never to touch SQL I offer an alternative. Code generation.

Most applications written to manage metadata have operations falling under the acronym CRUD-L. Create, Read, Upload, Delete and List. That's a good thing.  With a the knowledge of a data model you can calculate the SQL needed for those operations. No AI needed :wink:.

Taking a SQL code generation approach, unlike an ORM, let's you avoid dealing with SQL aside from the command to load it in the database[^4]. It has the benefit too of leaving the SQL as an artifact should you need it or change your oppionion of SQL. Code generation lets you avoid thinking about it or even typing it. It doesn't hide it either.

[^4]: Loading a file of SQL commands, say  "myfile.sql" targetting "mydatabase", is trivial with `psql`. Example: `psql mydatabase <myfile.sql`.

Picking Postgres and PostgREST as a target in Newt's data pipeline was deliberate. Setup can be calculated from a configuration file. It relatively easily to load generated SQL or script the setup. It is also possible with other SQL engines but the Postgres+PostgREST really fit the bill nicely.

## A recent insight

A colleague of mine recently demonstrated an innovative use of GitHub issue templates to trigger GitHub workflows for our library's staff. What impressed me about the demonstration (besides my colleague's cool application) was the YAML used to express what an issue was. After some thinking I realized my original YAML for Newt could be suplanted by GitHub YAML issue template syntax. This would simplify documenting and teaching Newt configuration.

GitHub YAML issue templates (abbr: GHYT) describes the data model of an type type. It does this in large part by pointing out elements you'd expected in an HTML web form. HTML input elements already suggest a mapping to SQL data types.  E.g. a plain `input` element is a string -> varchar. `textarea` maps to a `text` data type. An HTML5 data element expressed `input[type=date]` maps to a SQL date type. This isn't accidental. HTML web forms we designed to make it easier to get data into a database. If a table row is analogous to an object and an input element is analogous to a column I have a clear mapping for generating SQL without resorting to the developer knowing SQL data types. If we take an additional step forward and require JSON column support in our SQL database (MySQL, Postgres and SQLite have for the last decade) and we require the form to use JSON encoding then even elements like checkbox or select lists can be easily to send over the wire.

By extrapolation web components can be expressed as JSON objects. This allows Newt to be enriched with data types specific to libraries, archives and museums. A cororlary is the objects stored as a JSON column can be expressed as web components. Since composing a web component take effort it makes sense for Newt to include the common cases used in libraries, archives and museums.

## What's next?

The second prototype for Newt is focusing and developing an improved YAML syntax. The plan is to used [CITATION.cff](https://citation-file-format.github.io/) for application metadata while using GHYT for models and GHYT input types to describe route variables.  These changes also come with renewed focus on code generation targeting SQL for Postgres+PostgREST, Mustache templates for `newtmustache` web service and HTML elements for web forms and display.

The simpler YAML reusing existing YAML domain syntaxes along with code generation should ease or remove burden of writing middleware.

