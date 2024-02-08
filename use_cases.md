
# Use Cases

## cold, Currating Caltech People, Groups and Funders

In 2021 a web application called [Acacia](https://github.com/caltechlibrary/acacia) was built with a Go web service that mapped database access to a JSON API and use static files to access the JSON API and present a web based input system for lists of DOI sent to Caltech Library via EMail. The web app provided views to process the requests based on a queue table and hosted in a MySQL 8 database. This proved reliable but the burden of relying on JavaScript to assemble the web presentation was problematic in terms of accessibly support. This meant the approach wasn't suitable for public web services offered by the library. Additional building the UI via JavaScript did not reduce development time.

A similar approach was to be taken with cold, built initially in Python but contraints on developer time prevented completion of the project and managing simple table oriented data was found to be easier in a spreadsheet inspite of the problematic interface with a model that included a high number of columns.

In 2024 we are revisiting [cold](https://github.com/caltechlibrary/cold) taking the Newt approach knowing that many more custom applications like Acacia and cold are in our future.

### Requirements

cold is expected to maintain three types of objects. A list of Caltech people, a list of Caltech groups and a list of Caltech funders. It will need to render YAML lists for import into RDM as controlled vocabulariees and as a source of what to aggregate in our feeds website. The latter is used to allow facularty, staff and researchers to integrate their publication lists via the centrally provided content management system (they use a JavaScript widget to include data from Caltech Library's feeds). 

Additionally cold is moving beyond a flat spreadsheet oriented data model to a tree oriented one. E.g. A author may use multiple name variations over time. An organization's name may change or evolve or time. It is desirable to track the alternatives. This make continued curation of this data in a CSV file problematic at best. Furtunately storing lists in a database such as Postgres is well understood and the Newt platform using Postgres+PostgREST in a pipeline can support this type of content curation.

For Newt to be helpful to deliver cold the SQL needed to manage the database in Postgres should be generated from the data models supported in cold. Similar the web forms and content display should be deliverable via code generation. The Postgres+PostgREST JSON API means integrating content into other systems that rely on that data is just a web request away without having to change cold or the system using it's data.

### Developer experience

Development cold is scheduled for Spring/Summer 2024


## Evolving feeds beyond 1.5 by leveraging data pipelines to render static content

Moving from our original implementation of feeds.library.caltech.edu to version 1.5 with support for RDM repositories has been problematic and plague by bug regresions. Cleaning of legacy technical dept while maintaining backward compatibility has bogged the project down. What was originally planed to be a month of development in September 2023 has stretched into February 2024 with no end in site. It has become a case of whack a mole.

The feeds process original was batch oriented with updates aggregated nightly.  The buggy feeds v1.5 has allowed us to run updated aggregation once a twelve hours. This remains dismal as well as very difficult to debug. Bring down the build time to an hour or two would greatly improve the situtation. The performance bottle necks are largely in the time it takes to write aggregations to disk, read them back in then great the next smaller set of aggregations. This currently is done in a Python orchestrated and Shell.

Newt offers an escape hatch through its pipeline approach. Each level of aggregation can be thought of as a SQL VIEW. This can be generated in JSON directly by Postgres+PostgREST. The results can be stored in feeds database table. Similarly the JSON can be sent through a template engine, either Pandoc or Newt's Mustache engine, to render other formats (e.g. HTML, HTML include, BibTeX, RSS, CSV). This pushes all the sorting and structing problem into the database which is better suited than hand coded Python scripts to aggregate content. The resulting JSON, HTML, HTML include, BibTeX, RSS and CSV can be either written to disk in batch or pushed directly into an S3 bucket. 

This pipeline approach also means render one section of the feeds website becomes easily isolated by treating the build process as a webhook driven by a client script.

> Both use cases are theoretical at this time because the 2nd Newt prototype has not been implemented


