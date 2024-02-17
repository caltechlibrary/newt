
# What is Newt? Why does it matter?

Newt is a project of the [Caltech Library Digital Library Development Group](https://caltechlibrary.github.io/newt). Newt is an on going experiment in rapid application development for libraries, archives and museums (abbr: <abbr title="libraries, archives and museums abbrevation">LAS</abbr>).  How does Newt achieve "rapid" developemnt?

Newt targets metadata curation applications. The LAS community swims in metadata and always needs tooling to better manage it.

Newt's approach to creating a metadata curation application is based on the following characteristics.

- preference for "off the shelf" over writing new code
- configuration over writing new code
- data modeling and management placed squarely in your database management system
- database management systems must provide HTTP accessible JSON API
- data transformation (if needed) via a web service(s) that consumes JSON
- alignment of services avoiding overlapping capabilities
- code generation where appropriate
- write less code, keep it simple and focused

These characteristics and the narrow focus of metadata curation applications suggest a service oriented architecture. If we assume one or more web services can easily be combined. Delivering the results of a web browser request through combining services suggests a data pipeline[^211]. A collection of pipelines maps to the operations needed by web applications. If you know these ahead of time you should be able to calculate the resulting application.

Less code to write, a common proven architecture, a simpler solution.

[^211]: A data pipeline is formed by taking the results from one web service and using it as the input to another web service. It is the web equivalent of Unix pipes. Prior art: [Yahoo! Pipes](https://en.wikipedia.org/wiki/Yahoo!_Pipes)

Newt knitts those characteristics together using YAML description. The YAML describes the data models you want to work with. It also provides the configuration for the Newt applicatoins. It even includes configuration a light weight template engine and mapping the templates needed for the objects you're modeling.

The first step in a Newt based project is to write the YAML that describes the above, generate the code that can run in Postgres and configure PostgREST as well as generate basic Mustache templates to bootstrap your application.

In short, write less code to get a basic CRUD-L[^212] application.

[^212]: CRUD-L, ancronym for create, read, update, delete and list. These form the basic operations on metadata objects.

## What's "off the shelf"?

Software is composed of bits, it doesn't site on myself.  Off the self is software you don't need to write. Software that isn't going to cost you or organization a license fee or subscription charge. The software should be proven, with a good community of developer to help answer questions and a good track recordfor the rest of us humans. Good examples of "off the shift" are

- Firefox web browser
- Postgres, PostgREST
- Apache Solr
- Apache 2 web server
- Shibboleth for single sign-on

You don't write these things you use them. You may have them already installed. They run on macOS, Windows, Linux and the BSDs. Newt doesn't need to provide any of those services it can assume them and take advantage of that.

## What's my point about where data modeling and management take place? Why is that a Newt thing?

Database management systems like Postgres are extremely capable of managing data. The SQL language also is good at describing the data you want to manage.  We really should be taking full advantage of that. We certainly don't want to write code that replicates that ability. Yet we often do write that code and it complicates our programs terribly[^213]. Why? Here's my guesses.

[^213]: When Ruby on Rails gained fame this habit increased. Now it is routing for programming language frameworks to provide tools to model data and manage their evolution. I don't think this has improved things. I think it is a symtop of a bigger problem.

1. SQL is different from general purpose programming languages
2. In the web development genre we're been tought to techniques to cope with that. E.g. embedding SQL or using an ORM[^214]

[^214]: ORM, acronym for "Object Relational Mapper", before JSON column support in the Open Source databases this helped. Now not so much.

The enteria behind the way we've always done things is hard to resist. I suggest we do resist it. The ORM and the other libraries and frameworks that are used in many programming languages like Perl, PHP, Python, Ruby, Java, Go, etc. are a symtom of the challenge of shifting cognative gears between thinking about "how" a computer should perform at task versus "what" you want to computer to provide you. When you try to combine the "how" and "what" closely by creating embedded SQL statements things get messy really fast.  An ORM might mask this but it's there.

A Newt proposition is why accept the cognative clash? It made sense in 1993 but doesn't in 2024.  Let the off the shelf software to do that. Enter [Postgres](https://postgres.org) and [PostgREST](https://postgrest.org). If you think of those two pieces of software as a single service you have SQL managing your data but the rest of your application never touches SQL. It sees JSON and makes request via URLs. Look Ma, No SQL. If Postgres isn't you cup of tea you'll finsh software that does the same for MySQL database and even some projects that try to do that for SQLite 3.  Setup the database with tables, views and functions and you get a JSON API out of the box! It's liberating. 

## "Write SQL? I need to know that? Ugh."

SQL is a problem because it is so different from general purpose langauges many developers use day to day. One approach is suggest a good book and ask your fellow developer to learn it. That's a "no go" for lots of reason for many people.  I think we have an escape hatch. Newt comes with a code generator, include generating SQL code. Generate code doesn't remove the SQL problem but does mitigate it some. Many people find it easier to approach SQL if some is already written. It's a half way approach.

I believe it is worth learning SQL. I don't think you need to learn it all at once. I think it is useful to learn as you need it. Having existing SQL code gives you an opportunity to learns a little at the time. Modify it and see what happens. Newt is happy if you want to regenerate it.  Take your time wrapping your head around SQL's "what you want" orientation versus the "how to do" orientation of general purpose programming languages. Meanwhile let Newt bootstrap your Postgres+PostgREST service with generated SQL and PostgREST configuration.

## Wait, JSON is great but I need a web page!

JSON is a simple and proven way to express structured data. But it is not friendly to humans who expect to see an accessible web page. There are two ways to go right now on the web. The first way is to have a static web page that a web browser can use to then access the JSON data. The web browser then updates the web content to look like the desired results.  The has a huge hidden cost. Well actually many hidden costs but let me focus on one. The way I get JSON into a web page is to **rely** on JavaScript running in the web browser (or more recently a WASM module but that a whole other can of worms). The browser used the JavaScript to the content down, to update the web page contents rendering something that is hopefully what you asked for.  For programmers (I count myself in this group at times) it's catnip. Not only am I telling my server what to do with my program I'll telling that computer in your pocket what to do too! Low the power of the WWW. So that's the problem for the humans?

At least in North America and much of the Central Pacific wireless network access and can be quiet depressing. In 1993 people used modems and desktop computers, in 2024 most of the web is experienced through your "smart phone".  I live in Southern California. It has allot of rich people with latest in clever devices. Our our model networks are horrible. Even with a fancy phone. When you send JavaScript down the virtual wire you'll requiring significant network bandwidth just getting the JavaScript and the content it will assemble to your phone. Then your phone has to run the JavaScript. It is a terriby idea[^213]. Yes, I know that there are companies make bucks off this moddel[^218] but it doesn't make it a good idea[^215]. It is a rotten experience for us humans.

What's the better approach? Glad you asked, it's to have the web browser do only a couple network requests and get HTML, CSS and perhaps a couple media assets and bingo the human can read your content or interact with your web service. Yes, I know this sounds all reto but it works. It works really well. We used to have a fancy term for it, "Progressive Enhancement". If you're interested I recommend reading [Christian Heilmann's blog](https://christianheilmann.com/). I wish my bird was still as red as his.

Newt applications can send JavaScript down the virtual wire but they don't need to. Without enhancement a Newt application should work even using a text browser like [Lynx](https://lynx.browser.org/). A Newt application starts out handing back HTML by using a template engine to transform our JSON into a web page.  To encourage this approach Newt provides Newt Mustache. A simple highly configuration template engine. It expects to receive a JSON object, apply a Mustache template and hand back the results. The Mustache template langauge[^216] is described at their website <https://mustache.github.io/>.

[^218]: Companies like Google, Amazon, Wix and SquareSpace are happy to do business his way.  Why? It is because it drives the web users to the "native app". Native apps mean captured customers.  The open web provides options, they don't want you to have options.

[^215]: Don't take my work, read this <https://infrequently.org/2024/01/performance-inequality-gap-2024/> 

[^216]: Yes, I did just asked you to learn a new language but I like Christian Heilmann's and my facial hair. If I'm lucky you already know Mustache.


## Where does that leave us developers?

The code Newt generates is for a minimal metadata curation application. It is providing you with a bootstrap. It off the ground with something that is sort of working. Albiet not exactually what you want. At best it probably ugly and on cordinated, like a new born Aardvark. Like the Aardvark your app will become adorable if it survives.

You're likely going to spend your time in one of two areas. The back end meaning you mucking about with your database or the front end where you are mocking about with templates and static HTML, CSS, JavaScript and page assets. My hope is that you are freed up to spend time on the front end. Why? Because that is the side of things that us humans experience. While I've always considered myself a "back end" developer, I can exclaim with pride, "Who now cares about the back end? Unless its broken!". Mostly the brack end is doing data management and often boils down to our CRUD-L operations with perhaps some variations on the lists we produced.  If you need a search system I recommend Solr. It is off the shelf too and plays nice with the Newt appraoch. 

If you're lucky enough to have an audience are data analysts then you'll want to address them too. One approach is to dump snapshort of your Postgres database (or subsets) and drop into your static content directory. Then just link to them. You can also reach for JSON API provided by Postgres+PostgREST and a Mustache template to do it that way too. Trade offs either way, but that software engineering generally. Picking the best two of three.

## Newt, three tools for your toolbox

Newt provides a code generator, a Mustache template engine and a data router. With these three "off the shelf" tools you can take advantage of those other "off the shelf" tools like Postgres+PostgREST, Solr and friends.  All three Newt programs use the same YAML file to get their jobs done. The learning curve is primarly picking up Newt's YAML syntax. But that is for another chapter.

> I know I stacked that software someplace ...

### Off the shelf

- [Postgres](https://postgres) + [PostgREST](https://postgrest.org) (data modeling and management)
- [Solr](https://solr.apache.org) full text search engine (search and discovery)
- [Apache 2](https://httpd.apache.org) + [Shibboleth](https://www.shibboleth.net/) (controlling access)

### [Newt tools](https://github.com/caltechlibrary/newt)

- `newtrouter` is a stateless web service (a.k.a. micro service) that routes a web request through a data pipeline built from other web services
- `newtgenerator` is a code generator that can takes a set of data models described in YAML and generates SQL and Mustache templates
- `newtmustache` is a simple stateless template engine inspired by Pandoc server that supports the Mustache template language

These six programs cover allot of ground. They provide the core functionality of many systems built for libraries, archives and museam.

If you're inclined to readup on Postgres then I recommend [The Art of PostgreSQL](https://theartofpostgresql.com/), by Dimitri Fontaine. Gave me allot to think about.


