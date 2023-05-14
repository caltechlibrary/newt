
# Newt, a new take on making web applications

> confessions about my metadata obsession and simplification
> for the data scientist, archivist and librarian

R. S. Doiel
Caltech Library, Digital Library Development
rsdoiel@caltech.edu
https://caltechlibrary.github.io/newt

# Hello, my name is Robert, welcome to a rabbit warren 

- Today, Complexity => LAMP and its legacy
	- Apache 2 and NginX (web sites and single sign on)
	- Solr, Opensearch, LunrJS and Pagefind (Search)
	- Java and Ruby on Rails (Archivesspace)
	- Perl (EPrints), PHP (Drupal/Islandora), Python (Invenio, RDM)
	- Proxy systems (resolvers, Open Athens)
	- Static site generators (Pandoc, Github, R Studio)
	- Custom apps and SAAS (DIBs,EBSCO, Folio)
- Scaling
	- Scaling up, programmable infrastructure (the siren song of Google, AWS and Azure)
	- Scaling down (pack only what you need) <- this is what interests me
- Concepts
	- Static sites and there generators (a rebillion)
	- Structured data JSON, HTML and SQL
	- Front-end and Back-end (web server, web browser)
	- Microservices (choose only a few)
- Goal build a read write web app
	- SQL
	- Pandoc templates
	- Simple data flow described in a spreadsheet (CSV file)
- My Toolbox
	- Text Editor
	- Spreadsheet
	- Web browser
	- Web server
	- Pandoc
	- Postgres
	- Postgres + PostgREST = A Microservice for to manage data
	- Pandoc server = a template engine = A Microservice to transform data
	- Newt = A Microservice to route data flow
- What I am I packing?
	- Text editor, web browser, a spreadsheet 
	- Pandoc
	- Postgres
	- PostgREST
	- Newt
- Minimize the source Luke 
	- HTML + Pandoc
	- SQL (Postgres dialect)
	- How to map a URL route using a PathDSL in a spreadsheet 
- Finding our way out of the rabbit warren
	- Birds 1, static site (read only)
	- Birds 2, Postgres + PostgREST (read write, SQL, HTML, complicated JavaScript required)
	- Birds 3, Postgres + PostgREST and Pandoc server (read write, SQL, HTML, simplified JavaScript requirements)
	- Birds 4, Postgres+ PostgREST, Pandoc server and Newt (read write, SQL, HTML, no JavaScript)
- And then what? (conclusion)
	- How far can I take SQL and Pandoc?
	- What about file uploads?
	 

# Concepts

- Web servers
- Static site generators
- Promise of microservices
- Scaling down, a few off the self microservices