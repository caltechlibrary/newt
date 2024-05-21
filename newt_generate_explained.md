
# Newt Generator Explained

Newt comes with a code generator as part of the `newt` command. It uses the Newt YAML file to render Postgres SQL, PostgREST configuration and Mustache templates suitable to bootstrap your Newt based project.  How does it do this? What are the assumptions?

The Newt code generator works primarily with the "models" property in your Newt YAML file.

Our minimal useful application should be able to do five things - Create, Read, Update, Delete and List (abbr: CRUD-L) the objects you are modeling. It needs to offer these actions for each model defined in our project. The second Newt prototype assumes each model is independent but multiple objects may be needed in your application. If you need to combine models (not unusual) then you will need to enhance the generated SQL or templates to support integration. For now let us focus on the basics.

In the Newt philosophy we model our data and manage our data in the database. For us this means Postgres. By combining Postgres and PostgREST we create a JSON API based on what we've modeled and managed in our database. Our middleware, the stages in our pipelines, do not need to know anything about SQL. They don't ever touch it. This gives us a clean break from SQL and the rest of our system. It avoids the burden of doing the cognitive shifts when implementing middleware. This is a result of Postgres plus PostgREST automatically giving a JSON API. PostgREST is really clever. Once we have a JSON API or JSON data source to interactive the middle can focus on using that in conjunction with a rendering engine. In Newt's case the rendering engine supports the Mustache template language. The Newt development tool, `newt`, knows how PostgREST and Postgres need to be setup. It knows what is the minimum SQL that needs to be generated to support our CRUD-L operations. It can update the routes in the Newt YAML file to support simple data pipelines that talk to PostgREST and run the results through the Mustache template engine. The generator in `newt` is responsibly for generating a PostgREST configuration file called "postgrest.conf", a SQL file called "setup.sql" and a SQL file called "models.sql" that you can run in Postgres to configure and manage your model's data. The generator also generates a half dozen or so Mustache templates for rendering the HTML of your web application that match the routes for performing the CRUD-L operations.

## Bringing up Postgres+PostgREST and Newt Mustache

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

- `bird_create_form.tmpl` (contains the web form)
- `bird_create_response.tmpl` (response for the submitted form)
- `bird_read.tmpl` (display a single object we've modeled)
- `bird_update_form.tmpl` (contains the web form)
- `bird_update_response.tmpl` (response for the submitted form)
- `bird_delete_form.tmpl` (contains the web form)
- `bird_delete_response.tmpl` (response for the submitted form)
- `bird_list.tmpl` (list objects we've modeled)

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



