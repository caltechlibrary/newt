
# Use Cases

## cold, Curating Caltech People, Groups and Funders

In 2021 a web application called [Acacia](https://github.com/caltechlibrary/acacia) was built with a Go web service that mapped database access to a JSON API. It used static files to access the JSON API. The static assets presented a web based input system for lists of DOI sent to Caltech Library via Email. The web app provided views to process the requests based on a queue table and hosted in a MySQL 8 database. This proved reliable but the burden of relying on JavaScript to assemble the web human interface was problematic. This was particularly true when considering issues of accessibly support. This meant the approach wasn't suitable for public web services offered by the library. Additional building the UI via JavaScript did not reduce development time.

A similar approach was to be taken with cold, built initially in Python but constraints on developer time prevented completion. It this time the cold data model was table oriented. It was found to be easier in a spreadsheet in spite of this high number of columns. Work on cold was postponed.

In 2024 we are revisiting [cold](https://github.com/caltechlibrary/cold). Our data needs no longer warrant a flat table (e.g. alternate names used by authors).  Taking the Newt approach could allow us implement a custom web application of the type we need to cold quickly. It is also likely many more custom applications like Acacia and cold are in our future.

### Project outline

cold is expected to maintain three types of objects. A list of Caltech people, a list of Caltech groups and a list of Caltech funders. It will need to render YAML lists for import into RDM as controlled vocabularies and as a source of what to aggregate in our feeds website. The latter is used to allow faculty, staff and researchers to integrate their publication lists via the centrally provided content management system (they use a JavaScript widget to include data from Caltech Library's feeds). 

Additionally cold is moving beyond a flat spreadsheet oriented data model to a tree oriented one. E.g. A author may use multiple name variations over time. An organization's name may change or evolve or time. It is desirable to track the alternatives. This make continued curation of this data in a CSV file problematic at best. Fortunately storing lists in a database such as Postgres is well understood and the Newt platform using Postgres+PostgREST in a pipeline can support this type of content curation.

For Newt to be helpful to deliver cold the SQL needed to manage the database in Postgres should be generated from the data models supported in cold. Similar the web forms and content display should be deliverable via code generation. The Postgres+PostgREST JSON API means integrating content into other systems that rely on that data is just a web request away without having to change cold or the system using it's data.

### Developer experience

Development cold is scheduled for Late Spring/Summer 2024. It may be desirable to have the second Newt prototype available at that time.


## Evolving feeds beyond v1.5 by leveraging data pipelines to render static content

Moving from our original implementation of feeds.library.caltech.edu to version 1.5 with support for RDM repositories has been problematic. At present (Feb. 2024) It is five months behind schedule and is plagued by bug regressions. Cleaning up the legacy technical dept while maintaining backward compatibility has bogged the project down. The original plan was to launch it in September 2023 with shake down wrapping up in October. We're still shaking out the bugs February 2024.

The feeds process originally was batch oriented with updates aggregated nightly.  The buggy feeds v1.5 has allowed us to run updated aggregation once a twelve hours. That's an improvement. Sadly, feeds v1.5 remains very difficult to debug. Bring down the build time to an hour or two would greatly improve the situation. The performance bottle necks are largely in the time it takes to write aggregations to disk, read them back in then generate the next smaller set of aggregations. This currently is done in a Python orchestrated and Shell.

Newt offers an escape hatch through its pipeline approach. Each level of aggregation can be thought of as a SQL VIEW. This can be generated in JSON directly by Postgres+PostgREST. The results can be stored in feeds database table. Similarly the JSON can be sent through a template engine, either Pandoc or Newt's Mustache engine, to render other formats (e.g. HTML, HTML include, BibTeX, RSS, CSV). This pushes selecting and sorting problems into the database which is better suited than hand coded Python scripts to aggregate content. The resulting JSON can be take straight from Postgres+PostgREST. HTML, HTML include, BibTeX, RSS and CSV can be either written to disk in batch or pushed directly into an S3 bucket. Skipping the render to disk stage would significantly speed up the process but obscure the rendering process more. 

### Project outline

The feeds aggregates CaltechAUTHORS, CaltechTHESIS and CaltechDATA content. Content is presently aggregated by group, people and recent. Feeds v1.5 simplified the earlier implementation by removing support for recent 25 sub directory structures. This improved build time. This was enabled in large part because of IMSS content management system using the JavaScript widgets for content integration not the recent 25 files.  There are a few static assets that also are made available through feeds.library.caltech.edu. These include the JavaScript Builder widget as well as the JavaScript library the offers easy integration for the IMSS content management system.

The current approach needs to be improved by increasing observability in the data processing pipeline and in reducing the over all build times for producing the people and groups parts of the website.

The existing implementation (v1.5) uses database collections representing each repository by implementing a common data model. The dataset collections use Postgres for data management. This has allowed performance improvements through the use of [dsquery](https://caltechlibrary.github.io/dataset/dsquery.1.html) to produce aggregations as CSV files. dsquery results are read and processed to generate the JSON representations for each feed type supported person or group. in the people and groups directory 

I have observed that building the CSV files using dsquery saves consider time. This is because the process is taking advantage of SQL performance in Postgres.

When Newt's 2nd prototype becomes available it provides an opportunity to leverage Newt's data pipeline by creating all the JSON objects in the database then rendering the table results in one pass. This leaves could significantly reduce the site building time when considering Pandoc can be run as a web service and it is already used to render the Markdown, HTML, HTML include, BibTeX formats. The database approach affords us the opportunity to reduce the number of disk reads given the performance characteristics of Postgres and leave the site building process to just write content to the htdocs tree or perhaps directly to the S3 bucket.

It is desirable to produced a "combined" feed that aggregates content across our repository systems. If feeds output is staged in a dataset collection then them a combined repository view is possible is possible. It appears desirable that aggregation happens in the a single feeds dataset collection which also is well suited to the pipeline approach.

Steps to revising the existing feeds v1.5.

- [ ] Decide if this will be undertaken before feeds v2 or to solve the v1.5 bugs
- [ ] Create a "feeds" dataset collection for holding repository aggregations.
  - The dataset collection should use Postgres as the storage platform
- [ ] Audit the URLs in feeds v1.5 and determine the pipe lines needed to reproduce each content type
- [ ] Create the Newt data models for content type
  - Test the models and make sure they can render each content type in the appropriate formats
- [ ] Create the routes necessary to support the required pipe lines
- [ ] Generate the SQL, enhance if necessary, and load into Postgres dataset collections
- [ ] Decide if rendering will be to disk or to database for HTML, HTML Include, BibTeX and RSS
- [ ] Write a script that will trigger each of the pipelines needed to render the content
- [ ] Write a script that will send content to S3 bucket
- [ ] Update builder widget and CL.js to use combined_authors.json for combined authors publications and "combined.json" for aggregation across repositories

### Developer experience

Development of feeds v2 is scheduled for Summer 2024. Feeds v1.5 is on going in the interim. Newt prototype 2 needs to be available before this approach can be tested. A send to bucket service is needed as well as write to htdocs.  It will be important to have the option to write to a local htdocs tree for debugging and development.

Date is not schedule to try this approach before feeds v2.

