
# Newt's Origin story

## a small epiphany

Newt came from a small epiphany. I was writing yet another flask app and the de ja vu was intense. It got me thinking about how often I write the same code over and over except for minor variations. Why in the world was I doing this?  That caused me to step back, take stock and see if there was a simpler way.

Modern applications like Invenio RDM are created by aggregating many web services. This is a feature of service oriented architecture also known as micro service architecture. Invenio RDM is much cleaner than the venerable EPrints 3 repository system but it remains highly complex. Was complexity really needed? Was it an artifact of Zenodo? Zenodo is huge. Huge numbers of concurrent interactions (bots and humans) along with a massive amount of data and objects to manage. It is a complex system because of its scale. At its core RDM does two things really well. Metadata management and object storage. These are non-trivial at Zenodo scale and RDM does a really good job at managing them. But having migrated from EPrints to RDM recently I've come to question why do we built software systems so complex? My library is unlikely to need an application that scales as large as Zenodo. **Our needs are far more modest**.

I reviewed the other custom metadata management applications I've built, collaborated on or maintained over the last several decades. I feel like they all could be simplified. **It's time to see how much software we can create without writing new software.**

At Caltech Library we build allot of flask apps. Flask compared to many frameworks is a nice balance of flexibility and simple concepts. Looking at my applications I see that I should be writing much less. My middleware needs to be scaled back. It should be taking on less responsibility and be more narrowly focus. All my web applications use a database for content storage (e.g. SQLite3, MySQL or Postgres).  My code spends allot of time modeling data, transforming it on the way to or from the database as well as enforcing access rules. While I often use an ORM I'm still modeling data in my application. That's shouldn't be part of the middleware at all. This problem is not unique to my code or the code we write as a group in Caltech Library. I saw the same problem writing software at USC for a quarter century (as a student worker and later as staff). I see the same problem when using platforms like Drupal and WordPress.

How do we trim down our middleware? Looking at EPrints, RDM, and our custom applications where at Caltech I see several areas that middleware is covering that can be rethought or re-aligned.

1. Data modeling
2. Data validation
3. Data management
4. User management and permissions regiments
5. Data retrieval
6. Final data transformation for use by the web browser

I realized that all these features often overlap with other software in the systems we build on.  For historical or evolutionary reasons we've been putting those features into middleware. I believe these features often are made redundant depending on the how we use the other systems we build on.

So here's my propositions

- data modeling should happen os it can propogate to the database, validation services and rendered templates
- managing data should only happen in the database
- database needs to be accessible via URL and return JSON content so that SQL stays out of the middleware completely (the ORM is a symptom not a solution)
- user access and permissions should be pushed to the outer layer of a web application (e.g. Apache 2 integrated with Shibboleth)
- if the database supplies JSON responses we still need to turn that into HTML; use a a simple, stateless template engine
- use a single piece of software for routing data to our JSON API through the template engine

Our middleware efforts should focus that first stage (validation) and last stage (presentation). Everything use could be provided by "off the shelf" components. The middleware focus is in the final composition from the parts. It can use configuration to do that.

Proposed components for testing a the simple router concept

- Postgres database
- PostgREST to turns Postgres into a JSON API managed from Postgres
- A very light weight template engine web service
- Solr for full text search leveraging its JSON API support

The router needs to be configurable such that any given URL supported in your application can pass content between the database or search engine and through a template engine handing back the results.

Flask (like many other frameworks) suggest a excellent path foreword.  Flasks lets you assign a "route" to a function. A "route" is expressed with a path like string. It may include replaceable elements that the function can then use to do it's job (e.g. a record id passed in as part of the URL's path).  If you can identify a path using a simple description string and you pair that with a data pipeline sequence then now have a router. The code you write is configuration for the path and how it should proceed through a pipeline sequence.

Use case 1: contact the JSON API for the search engine, take the results and run them through the template service.

Use case 2: contact the database via JSON API, perform an action, take the results and run them through the template engine.

Use case 3: contact the search engine, get JSON results, run them through the template engine.

I was able to illustrate this approach in the initial Newt prototype I showed at SoCal Code4Lib in [July 2023](presentation/). The biggest challenge wasn't implementing the router it was inventing a YAML syntax. My initial YAML syntax was overly complicated and ugly.

## The 2nd prototype

A couple things fell out of that demonstration and subsequent talks with other developer colleagues.

- Postgres with PostgREST providing JSON APIs is very cool (I expected this response because it is very cool)
- People were not interested in data routing and pipelines (unexpected, routes seemed cool to me so I assumed it was cool to others)

This was expressed in the from of questions and comments.

- Do all my data modeling in SQL ... ?
    - I don't like SQL!
    - I don't know enough SQL!
    - I don't feel comfortable writing SQL!
    - Writing SQL isn't fun!
- I don't need a pipeline! The browser can assemble the page! ...
    - React/Angular/(replace with your favor JS framework) means I don't need templates.
    - Who cares about HTML, JavaScript can deliver everything!![^41]
    - But aren't templates the old thing people used to use? Why use them now?

[^41]: Example: SquareSpace and Wix deliver web pages as JavaScript.  I think this is to hide what they are really selling.

> NOTE: the exclamations and summation are my editorial. Everyone I encountered was very supportive and encouraging.

After taking this all in and having a long think 2024 arrived. These are my current conclusions.

1. When people saw Postgres combined with PostgREST and my controversial suggestion to use SQL to model your data the discussion shifted from frameworks to data modeling and what JSON would result. Granted people weren't looking forward to writing or learning more SQL. That's something I've noticed throughout most of my career. People who like SQL tended to become database administrators or database developers and get paid more when management by the Oracle cool-aide.

2. Shifting the discussion to where you model your data is the important bit. The fact remains that the database engine does a much better job of managing performance and your data then any middleware code you're likely to write regardless of framework or programming language.  I like writing in Go. It is a reasonably efficient language. If the database is setup up correctly running SQL queries will beat my processing data outside the database most days of the week. Arriving at SQL to manage data is generally a win.

3. Data routing shouldn't be either surprising or exciting. **Data routing should be boring.** It should be boring in the same way people take streets, sidewalks and bridges for granted. They are just there for you to use. Structural engineers and Urban planners may get excited about streets, sidewalks and bridges but that's because they help create them. The rest of us get excited when they are not available. Data routing through web services is just like that.

4. The SQL conundrum is a problem. While database engines like Postgres can support writing functions and procedures in non-SQL languages you still need to know and type at least some SQL. I think in part that is why we're still in the complex middleware quagmire. Avoiding messing about with SQL is a feature. On the plus side if your database hands you JSON and you can hand back JSON your middleware doesn't need to know about SQL and you can avoid the cognitive overhead embedding SQL into your non-SQL program.

## Addressing the SQL conundrum

> SQL is a problem because embedding SQL causes a cognitive shift.[^42]
> SQL is a problem because it is significantly different from general purpose programming languages

[^42]: The price of that shift is similar to when you inline CSS or JavaScript in HTML.

In most computer languages we are telling the computer how to do something. This is true even with object oriented languages, e.g. take this object and apply it's method. SQL takes a different approach. The SQL engine is going to do the heavy lifting. Instead we use SQL to express the result we want. In other words the "what" question versus the "how" instructions. When we embed SQL in a procedural or object oriented language we force these two together in close proximity. Newt's approach to take advantage of distance. By using SQL code generation we don't focus on the "what" at the same time we're focusing on the "how".

I feel comfortable with SQL. I've learned it informally[^43]. I've personally found it challenging when embedding in a general purpose language like Python. This is true even when I've used an ORM. It's a constant tug of war between thinking about Python (the how) then thinking about the SQL (what I want). Typically I've wound up assembling the SQL statements and running them directly in the database to gain insight into unexpected behavior. If the SQL is conditionally constructed then it become harder to figure out exactly when was sent to the SQL engine.

[^43]: SQL is built up from smaller domain specific languages. I've only learned each when I needed it.

The result of that experience has convinced me that keeping the SQL separate as SQL is helpful. I can work directly with the database and understand what I'm asking for (e.g. INSERT, SELECT, UPDATE, DELETE). When you do need to embed it using an SQL view can often simplify a complex SQL statement lowering the cognitive overhead.

Two general approaches appear promising to me. Separate your SQL out and don't embedding it. Use a combo like Postgres plus PostgREST to get a JSON API web service. Your middleware doesn't need to include any SQL at all then. It just calls URLs can procedurally processes results. If you can generate the SQL you need in the database because you know you will need to support CRUD-L then that saves time writing SQL to do that.

I'm not advocating avoiding learning SQL. But I am advocating learning as much as you need when you need it. I am also advocating your middleware should be ignorant completely of SQL. Working with human readable generated SQL can afford you the opportunity to learn SQL incrementally. It gives you the opportunity to learn to trust the database[^44] as a partner rather than simply as a place to store things.

[^44]: I would go so far as saying attempting to learn a SQL dialect all at once is a bad idea.

I have come to think about SQL's domain specific languages as service descriptions.

DDL, data definition language
: This language describes the data structure of the objects you're going to work with, e.g. table rows and columns. CREATE and ALTER are part of SQL DDL.

DQL, data query language
: This is the language you used to retrieve or list objects you stored, e.g. rows from a table. SELECT and WITH are part of the SQL DQL. You'll spend allot of time using this one when you write your own SQL.

DML, data modification language
: This is the language you used to create, update and delete objects. e.g. modify rows or columns in a table. INSERT, UPDATE and DELETE are part of the SQL DML.

Thinking of these as services lets me step back and simplify my interactions with our data. The nice thing about this approach is that orchestrating services can be thought of completely procedurally. This how Newt approaches SQL code generation. It can express the services it needs (e.g. create a record, retrieve a record, update or delete records) because it knows about the model being used. `newt generate` renders models as human readable SQL files (complete with comments). This gives us a change to both customize or evolve the SQL but also to pick learn it as we explore the generated code.

If your application only requires the basic CRUD-L[^45] then you might be able to skip the SQL completely.

[^45]: Data operations often fall under the acronym CRUD-L. Create, Read, Upload, Delete and List. That's a good thing. With a the knowledge of a data model you can calculate the SQL needed for those operations. No AI needed :wink:.

Taking SQL code generation approach like an ORM let's you postpone thinking about SQL while you're coding in your general purpose language. Code generation, unlike an ORM, let's you avoid the constant cognitive shifts. When you're ready to deal with SQL you focus your time in the database[^46]. Less frequent cognitive shifts makes the task more pleasant.

[^46]: Loading a file of SQL commands, say  "myfile.sql" targeting "mydatabase", is trivial with `psql`. Example: `psql mydatabase <myfile.sql`.

It was stumbling on the Postgres+PostgREST combo that remembered data pipelines.  Newt became obvious when I realized that the setup of Postgres+PostgREST could be calculated from a configuration file.  I quickly realized that the combo of Postgres+PostgREST meant any JSON data source code be a first stage in a pipeline. Of course other SQL engines can perform a similar role[^47] but the Postgres+PostgREST really fit the bill nicely for the systems I'm working with.

[^47]: MySQL can be coupled with MySQL REST Service (MRS) to achieve a similar results a Postgres+PostgREST

## A recent insight

A colleague of mine recently demonstrated an innovative use of GitHub issue templates to trigger GitHub workflows for our library's staff. What impressed me about the demonstration (besides my colleague's cool application) was the YAML used to express what an issue was. After some thinking I realized my original YAML for Newt could be supplanted by GitHub YAML issue template syntax. This would simplify documenting and teaching Newt configuration.

GitHub YAML issue templates syntax (abbr: GHYITS) describes the data model used by GitHub issues. It does this by describing the elements you'd expected in an HTML web form. My realization was I could cut the weight of my original YAML syntax considerable by letting a description of HTML input elements infer a mapping to Postgres SQL types.  E.g. a plain `input` element is a `string` -> `varchar`. A `textarea` maps to a `text` data type. An HTML5 data element expressed `input[type=date]` maps to a SQL date type. This isn't accidental. HTML web forms were designed to make it easier to get data into a database. If a table row is analogous to an object and an input element is analogous to a column I have a clear mapping for generating SQL without resorting to the developer knowing SQL data types. If we take an additional step forward and require JSON column support in our SQL database (MySQL, Postgres and SQLite have for the last decade) and we require the form to use JSON encoding then even elements like checkbox or select lists become easily to send over the wire and can have simpler expressions in SQL.

By extrapolation web components can be expressed as JSON objects. This allows Newt to be enriched with data types specific to libraries, archives and museums. A corollary is the objects stored as a JSON column can be expressed as web components. Write a web component take effort. The duplication of effort can be reduced if Newt to includes the common cases used in libraries, archives and museums.

Expirmenting with GHYITS lead to a simpler data model since most of what GHYITS is not used by Newt.  E.g. form elemnts can be describe by an id and attribute list with a boolean flag to identifier a data model's unique identifier attribute used by in SQL and routing.

## What's next?

The second prototype for Newt is focusing and developing an improved YAML syntax. The plan is to model HTML5 form element simply in YAML as the attributes of a simple data model. The second prototype will include a interactive "modeler" that be used to generating and update the Newt YAML file.  That file functions as a abstract syntax tree decribing a web application built from off the shelf services, Newt Router, and Newt Mustache template engine.

By focusing on modeling data using the existing HTML5 form elements rendering the webform becomes straight forward. It has the advantage too that the basic HTML5 form elements infer a SQL data type.  The form elements can be extended by aliasing and mapping aliases to either inputs with regexp validation or to external web components. This will provide an avenue for exploring more complex models in the future.

Simularly by treating the Newt YAML file as an abstract syntax tree other generators can be created to in prinincple could generate Go or Python application code allowing development of programs using Newt only as a protyping or bootstrapping tool.
