
## Origin story

Newt came with a small epiphany. The applications for managing structured data have been generally split into descrete set of services already (e.g. Invenio RDM). These can be thought of as a JSON API for managing data, a presentation layer (i.e. template engine) and a front end built using HTML, CSS and JavaScript. The JSON API always winds up talking to a relational database system such as Postgres, SQLite or MySQL. The middleware mapping the database management system to respond to request is essentially a whole bunch of boiler plat code even if you use an application framework like Flash or Django. I stumbled on PostgREST after I started using Postgres for some library project. PostgREST provides a JSON API based on the contents and configuration of a Postgres database.  There is no need to write a bunch of boiler plat to create yet another JSON API. Just let PostgREST do it for you.  The problem remained that JSON isn't what the web browser needs to provide a human friendly experience. The JSON needs to be tansformed into HTML most of the time.  A good solution for transform structure data from one form to another is Pandoc. When I realized the Pandoc can run as a web service I realized I had been wasting my time with allot of Python code and I could skip the middleware completely if I had a way to map a request (i.e. route a URL request path) to a short data pipeline that talked to PostgREST, got results and then passed them through Pandoc rendering the result as HTML. I just need to describe the routes. 

After experimenting with this idea I realized I could take it a step further. Before sending requests to the JSON API I could vet the requests and make sure the query parameters, URL path or the contents of a PUT or POST were valid. This can be accomplished by validating the data involved. There is also overlap in spefici metadata types used. Example common identifiers like DOI, ORCID, ROR are well formed and you can validate the form easily. So adding validation support server side makes sense before allowing the request to pass through the pipeline. The best way to know what to validate is by recording the "model" in some fashion.

All this can be described in a YAML file who's syntax is easy to pickup by combined existing YAML syntaxes found in CITATION.cff and GitHub issue templates. Only a YAML syntex describing the route in the data pipeline is needed to round out the functionality needed to compose an application from a small collection of the micro services[^2].

Newt uses YAML to express the routes in the pipeline and uses GitHub issue template syntax[^3] to describe a data model. Data models can be associated with one or more rouates. 

[^2]: Newt data routining has been tested with PostgREST+Postgres, Pandoc running as a web service. It should also work with other JSON oriented data sources such as Solr, Elasticsearch and Opensearch

[^3]: GitHub has a simple template langauge for custom issues. It is general enough to be repurposed as a model language in Newt. The language describes the model from the point of view of a web form. Newt takes it a step further by extrapolating the mapping to infer SQL reprenting the data model. This provides a strong alignment between the web form, the JSON API provided by PostgREST+Postgres and the necessary SQL needed to bootstrap PostgREST+Postgres databases.

The venerable solution to transforming approach to that has been template engines. The problem with template engines is the lack of standardization and the problem of becoming too feature rich[^1]. about when I realized that all I needed was a data router that could map a web browser request to the JSON API provided by PostgREST+Postgres and Pandoc running as a service. That setup could replace most of the applications I'd written for the last decade or two. It would fit most of the small web applications I had previously built in PHP or Python for my library. Generalizing the concept of simple data router for a data source and render engine also meant I had an easy integration point for most the institutional software we currently run in our library. So I wrote a data router to do just that.

[^1]: Example, PHP started out as template system for Perl CGI programs before becoming a full featured programming language

I demonstrated the Newt concept to my colleagues with a prototype. The prototype talked to a JSON API provided by PostgREST+Postgres and used Pandoc running in server mode for a rendering engine. I got some polite supportive comments. No one was particularly excited by it. I demonstrated a prototype Newt at a my local SoCal Code4Lib group. There people were excited by PostgREST and Postgres and not so excited about data routing. This was discouraging. I thought I was barking up the wrong tree. Eventually I realized the ambivalence of the router was a type of success. Newt isn't exciting. Newt should never be exciting. It just routes data! You configure it and forget it. It just runs.

The important take away was I had failed to appreciate how Newt successfully shifted the discussion from programming language frameworks, package management and build systems to to modeling data in SQL and using simpler HTML5, CSS and JavaScript for display.
