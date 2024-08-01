
# Newt's `newt generate` Explained

// FIXME: Needs to be rewritten given code changes and evolution of AST.

Newt comes with a code generator. It is part of the `newt` command. It uses the Newt YAML file to render Postgres SQL, PostgREST configuration and Handlebar templates suitable to bootstrap your Newt based project.  How does it do this? What are the assumptions?

The Newt code generator works primarily with the "models" property in your Newt YAML file.

Our minimal useful application should be able to do five things - Create, Read, Update, Delete and List (abbr: CRUD-L) for objects you are modeling. It needs to offer these actions for each model defined in our project. The Newt assumes each model is independent but multiple objects may be defined in your application. If you need to combine models (not unusual) then you will need to enhance the generated SQL, routes and templates to better support that type of integration. For now let us focus on the basics.

The Newt philosophy is to model data in YAML and then mangage the data via data management system. It presumes the data management system provides a JSON API for interacting with it. In Newt this is referred to as a "JSON data source".  What is a data management system?  Typical this will be a JSON friendly database combined with a JSON API.  Newt provides code generation for Postgres (the database) plus PostgREST (the JSON API). The code generation includes setting up a simple configuration file, generating SQL to define tables, schema as well as SQL functions. The result can be used to bring up an instance of Postgres+PostgREST for your project.

Newt also supports using Dataset+datasetd as a primary JSON data source.  Dataset can store JSON documents in a pairtree, in MySQL, Postgres or SQLite3 database.
Newt's Dataset+datasetd support assumes you're using SQLite3 as your storage format.


Newt code generation can accomplish the following tasks

1. generate the basic configuration files needed for your JSON data source
2. Provide any additioan SQL or data langauge files needed to implement CRUD-L opperations for your models
3. Configure (wire up) the routes and templates mapping in your Newt YAML file
4. Generate Handlebar templates for the created routes in your Newt YAML file
5. Generate TypeScript middleware suitable for Deno to validate you pipeline inputs

## Bringing up Postgres+PostgREST and the Newt Template Engine

Postgres and PostgREST provide us with a JSON API that can be used in a Newt router pipeline. Following the Newt philosophy data modeling and management happen in the database. That is done via the SQL language which the Newt generator can produce. After the database is setup in Postgres with the scheme, tables, functions and views needed for basic CRUD-L operations you can generate a compatible PostgREST configuration file. Here's the steps what the generated code does based on the models defined in a Newt YAML file.

2. For each Newt model create a table in the database/scheme
3. For each Newt model create a function for each of the CRUD[^71] options[^72]
4. For each Newt model create a SQL function to handle the list view of the CRUD-L
5. Generate a PostgREST configuration file.

When the "postgrest.conf", "setup.sql", and "models.sql" are generated then do the following.

1. Create the Postgres database and schema[^70] if they don't exist
2. Run the "setup.sql" and "models.sql" via the `psql` client.

~~~shell
newt generate app.yaml
dropdb --if-exists myapp
createdb app
psql app -c "\\i setup.sql"
psql app -c "\\i models.sql"
~~~

NOTE: The generator always generates the files using predefined names. If generated files already exist it will rename the old version with a ".bak" file extension before writing the newly generated code. If you've customized the generated files then the new version will NOT have those customizations. If you are customizing the generated files you have a couple options. First is to rename them or put them in a sub directory. Another would be to develop a workflow where you diff your modifications and try to apply the resulting patch file after regenerating the files.

NOTE: The Newt generator will alter (add or update) routes associated with the generated SQL functions and the Mustache templates. It will also alter the templates defined your Newt application YAML file.

When we ran `newt generate app.yaml` Mustache templates were also created in the current directory along side the "setup.sql", "models.sql" and "postgrest.conf" files.  The Mustache templates are named according to a model name and their role. When we create, update or delete objects we've modeled there are two round trips needed in the request and response cycle with our browser.  That also means there needs to be two templates, one presenting the web form and the other showing the processed result.  The model name and role are concatenated with an underscore, "\_", form the base of the filename followed by the ".tmpl" extension. The read and list operations only require a single template each. If I have a model with the id of "bird" then the resulting template files would be as follows.

- `views/bird_create_form.hbs` (contains the web form)
- `views/bird_create_response.hbs` (response for the submitted form)
- `views/bird_read.hbs` (display a single object we've modeled)
- `views/bird_update_form.hbs` (contains the web form)
- `views/bird_update_response.hbs` (response for the submitted form)
- `views/bird_delete_form.hbs` (contains the web form)
- `views/bird_delete_response.hbs` (response for the submitted form)
- `views/bird_list.hbs` (list objects we've modeled)

When we generated the Mustache templates the `templates` properties were updated to include these templates for the Newt Mustache template engine. Also
the related routes were updated with a pipeline that interactive with the Postgres+PostgREST JSON API sending that response through the Mustache Template engine.

## Testing our generated application

The `newt` command can run Newt Router, Newt Mustache and PostgREST using the "run" option. We can test our generated application by "running" them and then pointing our web browser at the Newt Router localhost and port to test.

1. Newt Run
2. Point our web browser at the router URL

In a shell/terminal session type the following.

~~~shell
newt run app.yaml
~~~

Now launch your web browser and point it at the URL indicated (e.g. "http://localhost:8010") to test your application.


## Where to go next?

There are three "places" to customize a Newt application. The first would be to customize data processing inside Postgres. Postgres is a full feature object relational database system. Most of the SQL:2023 Core. That provides allow of data management capability. Postgres supports a wide range of data types, more than Newt uses. Postgres also supports many popular programming languages embedded inside Postgres. This includes the [Python programming language](https://www.postgresql.org/docs/current/plpython.html). When you are using a general purpose language like Python inside your database system you have synergy between actions on your data, via SQL triggers, and the functions you write in Python reacting to the data. PL/Python may also be used to reach outside of Postgres to perform other tasks, e.g. feed indexing updates to Solr/OpenSearch/ElasticSearch or send an email.

If you can integrate supports for other services by adding additional "routes" to your Newt YAML file. You can also add steps to your pipeline.  Let's say you have an object modeled in Postgres but you also want to include additional information from another service (e.g. ORCID). You can write a simple web service that takes the object as generated by the PostgREST request and processes it further before handing the result back to Mustache for rendering.

Finally the third options is to enhance your templates with CSS, JavaScript and other static web assets.  Newt's router functions both as a data request router but also as a static file server.  If you wanted to integrate external content with a web form (e.g. ORCID or ROR information) you could do this via JavaScript and have their web browser connect to the external service retrieving additional metadata before the form is submitted.


