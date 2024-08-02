
# Newt's `newt generate` Explained

// FIXME: REWRITE NEEEDED based on direction of current prototype

Newt comes with a code generator. It does the heavy lift of knitting your application together from off the shelf parts. The code generator is used to create a new Newt YAML file as well as SQL, TypeScript, Handlebar templates and other miscelanios code bits. You run the code generator via the Newt development tools, `newt` (or `newt.exe` on Windows). The tool is truction around the concept of an actions on an object.

The basic worflow of code generation is as follows

1. create/update your application's Newt YAML file 
2. create/update your data models
3. generate the application
4. run, test and repeat

Your application at this stage should be able to create modeled objects, retrieve modeled objects, update modeled objects, delete modeled objects and list modeled objects. When youare happy with the base functionality you can then turn to refining the application by customizing the generated code, adding custom routes and templates as well as browser side functionality delivered via Newt's data router static file service.

## Generated files for Postgres+PostgREST based projects

The Newt code generator delivers the following code assets for your Postgres+PostgREST based Newt Application

`postgrest.conf`
: A basic PostgREST configuration file for your project

`setup.sql`
: This file is used to configure Postgres to work with PostgREST. This is code you can execute via the psql shell.
Usually you don't check this file into your code (e.g. GitHub) repository as it will need to contain DB credentials for the application's DB account.

`models.sql`
: This file setups the database schema as well as creates Pg/SQL functions for use with PostgREST in your CRUD-L application operations

Depending on what the setting of the base directory and model name it'll create the follow templates.
In this example `<base_dir>` would be replaced by the "base directory" setting your your Newt YAML file and `<model>` would be replaced by the model name.

`<base_dir>/<model>_create_form.hbs`
: This implements the HTML data entry form create a new object based on `<model>`.

`<base_dir>/<model>_create_response.hbs`
: This implements the HTML response from the create object data entry form.

`<base_dir>/<model>_read.hbs`
: This implements the HTML read only view a object previously created or updated.

`<base_dir>/<model>_update_form.hbs`
: This implements the HTML data entry form to update an objected previouly created.

`<base_dir>/<model>_update_response.hbs`
: This implements the HTML response from the update data entry form.

`<base_dir>/<model>_delete_form.hbs`
: This implemented the HTML form to delete an object previously created.

`<base_dir>/<model>_delete_response.hbs`
: This implemented the HTML response from deleting an object.

`<base_dir>/<model>_list.hbs`
: This lists objects previously created or updated.

Once your code is generated you'll first need to setup the Postgres database to work with PostgREST. That is done by using the Postgres `psql` shell to run the SQL code. You first run `setup.sql` and then run `models.sql`.

After that you're ready to test your basic application. The `newt` (or `newt.exe` on Windows) tool can also "run" your generated code. You would use the "run" action and the name of your Newt Project file.

## Minimal functionality

The generated code, after settting up Postgres+PostgREST has a minimum level of functionality. It implements the five basic operations to manage metadata - Create, Read, Update, Delete and List (abbr: CRUD-L) for each object type modeled in your Newt YAML file. Newt assumes each model is independent but multiple objects may be defined in your application. If you need to combine models (not unusual) then you will need to enhance the generated SQL, routes and/or templates to better support that type of integration. For now let us focus on the basics.

## Newt's Approach

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


