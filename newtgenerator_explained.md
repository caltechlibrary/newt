
# Newt Generator Explained

Newt comes with a code generator called `newtgenerator`. It uses the Newt YAML file to render Postgres SQL, PostgREST configuration and Mustache templates suitable to bookstrap your Newt based project.  How does it do this? What are the assumptions?

The Newt code generator works primarily with the "models" property in your Newt YAML file.

Our minimal useful application should be able to do five things - Create, Read, Update, Delete and List (abbr: CRUD-L). It needs to offer these actions for each model defined in our project. The second Newt prototype assumes each model is indepent. If you need to combine models (not unsual) then you will need to enhance the generated SQL. For now let us focus on the basics.

In the Newt philosophy we model our data and manage our data in the database. For us this means Postgres. By combining Postgres and PostgREST we create a JSON API based on what we've modelled and managed in our database. Our middleware, the stages in our pipelines, do not need to know anything about SQL. They don't ever touch it. This gives us a clean break from SQL and the rest of our system. It avoids the burden of doing the cognitive shifts when implementing middleware. This is a result of Postgrs plus PostgREST automagically giving a JSON API. Clever stuff really.

## Bringing up Postgres+PostgREST

Postgres and PostgREST provide us with a JSON API that can be used in a Newt router pipeline. Following the Newt philosophy data modeling and management happen in the database. That is done via the SQL langauge which the Newt generator can produce. After the database is setup in Postgres with the scheme, tabes, functions and views needed for basic CRUD-L operations you can generate a compatible PostgREST configuration file. Here's the steps I follow after creating my Newt YAML file.

1. Create the Postgres database and schema[^70] if they don't exist
2. For each Newt model create a table in the database/scheme
3. For each Newt model create a function for each of the CRUD[^71] options[^72]
4. For each Newt model create a SQL view to handle the list view of the CRUD-L
5. Generate a PostgREST configuration file.

NOTE: The generator always generates it's output to standard out. Errors, if they occur are sent to standard error. This makes it easy to script the Newt generator in shell or in a Makefile.

### Generating some SQL

Let's turn our attention to setting up Postgres database. Several things need to be accomplished and the Newt generator will generate the SQL needed for all of them. It will even generate the SQL to make managing your content via PostgREST JSON API easier.

OK, now that we know what we're generating lets generate this SQL from our app.yaml. We'll save the SQL output in a file called "app.sql"

```shell
newtgenerator postgres app.yaml >app.sql
```

In the statement above "postgres" tells the generator target the Postgres database SQL dialect. 

Take a look at "app.sql". Notice that the SQL if the code contains comments.  If you've included the `.description` attributes in our models then the SQL also includes them. This means you can continue improving the SQL manually if you like. Another approach would be to create a second SQL file with your modifications and use SQL `alter` or `drop` statements where you need to augment or replace the Newt generated parts.

This is a text that explains Postgres database so I am going to keep this short. If you want to read about Postgres to better understand it I recommend the [Postgres website](https://postgres.org) or the book [The Art of Postgres](https://theartofpostgresql.com/ "website for the book").  is heavely commented, those your

What was generated?  The specifics will depend on what was in the "app.yaml" file, specifically the models property.  Do you see comments you the generate SQL file? These were form from the `.description` properties found in the YAML. By taking the time to add description properties are making the generated code more human freindly. While these are optional properties they are really handy when you evolve the code generator or come back to a project later.

The first section of the SQL generate should have statements to create the database if it doesn't exist and to create a name space if it doesn't exist. Next you should see SQL for creating tables, one for each model. Then you should see the creation of views. The SQL views fill the role of the "L" in the acronym CRUD-L.

The nice things about having views is it make it easy to see what our application will see when accessing the JSON API to list objects based on our models. This parity becomes very helpful when you are extending your application or trying to debug why the JSON API isn't giving you what you expected. A view is accessed through the SQL `select` statement.

After a model's view Newt will generate functions for create, read, update and delete. Why do this?  First it allows us a consistent abstraction should the underlieing data model change in the future. Second you can enhance the function to perform additional validation if needed. Example, let's say you've imported the Python library idutils into Postgres. You can use that library to validate identifiers like DOI, ORCID, ROR, arXiv, ISNI, etc. Since the operation is a "function" rather than an SQL prodecedure you can let the requestor there was a problem or take actions to mitigate it. A function can keep the conversation going with the web browser at the other end of the wire.

You might be wondering about how things are name.  The functions and view names are formed from the model name appened with "\_", followed by the action. Newt trys to give you a robust bootstrap. One of the ways that it accomplishes that is by naming things predictably. This improves code readability for us humans.

After defining the models tables, views and functions the end of the generated SQL sets up the access needed by PostgREST to provide the JSON API. Essentially that is a block of permission granting statements. When you read through these you need to understand them and make sure they comply with your organization's database practices. Newt is only making a guess what is OK for you. See the PostgREST website for more information about how Postgres and PostgREST work together, <https://postgrest.org/en/stable/>.

### Generate our PostgREST configuration

The generator can generate a PostgREST configuration file, `postgrest.conf`. This is the first step. Next step is to setup our Postgres database to model our data, manage our data and to allow PostgREST to access our database appropriately. This is all done via SQL. So how to we generate. To generate the `postgrest.conf` you use the newtgenerator with a YAML file named "app.yaml" you would do the following.

```shell
newtgenerator postgrest app.yaml >postgrest.conf
```

In the above command the "postgrest" (note the "t" at the end) tells the generator what needs to be generated. In this a Postgrest configuration because that is what PostgREST needs.

## Ready for testing our JSON API

At this point we have defined a predictable JSON API and set of related URLs we can use in our data pipelines.  You should now have three files in our application directory, two were generated by the Newt generator.

- app.yaml, you created this by typing it in or copying an example and modifying it for your needs
- postgrest.conf was generated by Newt and holds a configuration used when you startup PostgREST
- app.sql is the SQL file that sets up Postgres to work with our application

Now we need to set things in motion. You can load your SQL into Postgres via the Postgres repl called `psql`. You'll need to have administrative privilages to run the SQL as it will create a new database for your application. See <https://www.postgresql.org/docs/current/sql-createuser.html> about creating database users for details.

On my computer I have Postgres 16 running and my personal account has admin privilleges so I use `psql` like this to create the database, schema, tables, views, functions and permissions needed for the my application.

```shell
psql < app.sql
```

If that worked then we can try out starting PostgREST and use [curl](https://curl.se/) to see if our JSON API works. The exact URL will be dependent on the database and model names setup in Postgres. The PostgREST webpage explains the JSON API, see <https://postgrest.org/en/stable/references/api.html>. If I have created a model named "people" I can start PostgREST then test of it is available like this.

```
postgrest postgrest.conf &
curl http://localhost:3000/people
```

NOTE: This first command starts PostgREST as a back ground process, you need to kill it if the data base changes table definitions, views or you've added views, tables or functions. See the PostgREST documentation for details about when you need to restart PostgREST, <https://postgrest.org/en/stable/references/schema_cache.html>.

That should return a list of people. Since our database tables aren't populated it should be an empty list. At this point I suggest getting comfortable with Postgres and PostgREST. While the generator creates the SQL and configuration needed that doesn't substitute for understanding how it is working and actually how you might want to use it!

[^70]: In Postgres a "schema" is really a name space for the tables in a database. It should be confused with a specific SQL table definition, often refered to using the same term "schema". In Postgres a "schema" provides a hook to manage table accesses and provide a namespace for referencing them es and provide a namespace for referencing them. See <https://www.postgresql.org/docs/current/ddl-schemas.html>

[^71]: CRUD is an acronym for "create, read, update and delete", CRUD-L stands for "create, read, update, delete and list". CRUD-L are the basic operations you perform when managing metadata.

[^72]: Newt creates functions for each of our CRUD operations. It does this to allow you a clean way to keep the JSON API constent even when you data models might evolve in differently.


## Generating some templates for a web UI

Newt comes with a light weight template engine called Newt Mustache. It implements Mustache templates. The Newt generator knows how to generate those. Newt can generate a template for each of our CRUD-L operations for each model. To know what template you want to generate you need to tell the generator you want to generate mustache templates, which model you generating it for and what action the template will model. As mentioned previously Newt generator writes the generated code to standard output and errors if encountered to standard error. Like with generating SQL and a configuration file this allows for flexibility in scripting via shell or a Makefile. Here's an example of the commands I type to create the templates for our people model.

```shell
newtgenerator mustache people create >people_create.tmpl
newtgenerator mustache people read >people_read.tmpl
newtgenerator mustache people update >people_update.tmpl
newtgenerator mustache people delete >people_delete.tmpl
newtgenerator mustache people list >people_list.tmpl
```

If you examine the resulting templates you'll notice that create, update and delete include webforms nad use the model types you describe. On the other hand the templates for read and list do not include webforms just some standard markup elements.  I expect you'll want to enhance these to meet you applications need but they should function well enough to test your data pipelines and debug them.

I usually get the back end setup and tested before moving to make the application pretty and enhancing the browser experience.

There are two approaches to testing your templates. One is to use them as the last stage of your JSON API. Another is to configure the templates in your YAML and run Newt Mustache service and use curl. In in the example below I'm assuming you've mocked up a person record in a JSON file called person.json. We can then test come of our templates to see how they fit the bill. I'm assuming you've setup Newt Mustache to provide templates based on the names of the templates. I am also assuming Newt Mustache is running on port 8040 in this example.

```shell
newtmustache app.yaml &
curl --data '@people.json' http://localhost:8040/people_read.html
```

You should get back an HTML page with the content of "person.json" in it. If so that template is working.  The create template should return an empty web form as it is used to "create" new model instances and a model idea isn't available.  If the "person.json" JSON included an object id then you should see it in the update form as a hidden field. Update should not enable creating new objects in must cases. Similarly if you mockup a list of people in a JSON file called "people.json" then you should be able to test the list template too. 

I generally will work directly with mockup JSON files and Newt Mustache to get the web interface I want for my application.

It is important to remember that Newt Mustache reads the templates in at program start. If you revise your templates you need to **restart** Newt Mustache. In this way Newt Mustache and Newt Router are like PostgREST server that you will need to be comfortable stoping and start the services as you continue your development.

