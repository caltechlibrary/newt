
# Newt Model explained (concept draft)

## Overview

Newt Modeler provides an interactive way for generating and managing your data models expressed in the Newt YAML file. It let's you avoid directly typing YAML. While it is interactive it is still a command you run in your shell. It lets you adding, modify and remove models. When modify a model it lets you add, modify and remove elements. When you are done machine changes at the application, model or element level you can "quit" to indicate your are done making changes. When you quit at the models list level it will prompt you to save the changes if you've made any. 

Since your Newt application flows from your data model(s) the modeler takes care of creating the routes and templates entries for each model you add. Likewise takes care of cleaning up any routing or template references previously generated the modeler when you remove a model.

By using `newt config`, `newt model` and `newt generate` you can create a a deliverable Newt application.  The minimal application will depend on Postgres+PostgREST as well as Mustache templates for rendering the human user interface.

## Running the Newt Modeler

In this short tutorial we are going to create a web application that says describes the contents of a garden.  We will use the `newt` command to configure, model and generate our application.

## Tutorial on the Newt Modeler

We will cover three steps. Creating a new "garden.yaml" file with "config" action, modeling our garden data and generating some code.

### Step 1, create a "garden.yaml" Newt YAML file

We can create or update our "garden.yaml" using our standard `newt config` command. For this tutorial just accept the defaults.

~~~shell
newt config garden.yaml
~~~

### Step 2, defining and managing models for "garden.yaml"

The first stage of the modeler is to present a list of existing models. It then gives you options of "add", "remove" "modify", "config", and "quit". Launch our modeler with `newt model garden.yaml`

~~~shell
newt model garden.yaml
~~~

By default the command `newt config` manages the "applications" attribute in our "garden.yaml" file. It doesn't create any models, routes, or templates.
created an "garden" model. When we use `newt model` the models defined will then trigger the creation of routes and templates attributes needed by the application.

When you launch `net model` you are presented with a text menu that looks like this.

~~~shell
Enter menu letter and id


Menu [a]dd, [m]odify, [r]emove, [q]uit (making changes)
~~~

No models have been defined yet.  Let's add one called "garden" by press "a" key followed by enter key.

You should now see a prompt like below.

~~~shell
Enter model id to add:
~~~

Type in "garden" without the quotes, then press enter key

Now the menu should look like this.

~~~shell
Enter menu letter and id

	1: garden

Menu [a]dd, [m]odify, [r]emove, [q]uit (making changes)
~~~

An empty model called "garden" has been created. If we want to remove it we would type the following ending in pressing the enter key.

~~~shell
remove garden
~~~

At this point the menu would not show any models.  Add "garden" back so we can continue with modifying our newly created "garden" model.

This is what you can type as a short cut.

~~~shell
add garden
~~~

You should see this menu again.

~~~shell
Enter menu letter and id

	1: garden

Menu [a]dd, [m]odify, [r]emove, [q]uit (making changes)
~~~

We now ready to "modify" our model. You'll see the theme of `Menu [a]dd, [m]odify, [r]emove, [q]uit (making changes)` through out the modeling process targeting the specific level you're refining. 

Type the following at this menu.

~~~shell
modify garden
~~~

This will take you to a new menu focused on managing the garden model.

~~~shell
Enter menu letter to manage model

	1: Model Id: garden
	2: Descriptions: ... description of "garden" goes here ...
	3: Elements: oid


Menu: model [d]escription, [e]lements, [q]uit (making changes)
~~~

This menu is slightly different because we have different types of parameters that can be modify in our model.
An individual model contains an identifier. This you can't change but is displayed so you know what you are
modifying.  At the "Menu": line we see that we have choices of description, element and quit. Let's modify
the description first to be something less generatic.

Type the following followed by pressing the enter key.

~~~shell
d This is a model of my garden
~~~

Notice that I typed "d" instead of the whole word "description". This short cut can be used through the modeler where
a single letter is indicated by square brackets, `[d]escription` would mean you could type "d", `[q]uit` would mean
you can type "q" instead of "quit".

The updated menu should look like this.

~~~shell
Enter menu letter to manage model

	1: Model Id: garden
	2: Descriptions: This is a model of my garden
	3: Elements: oid


Menu: model [d]escription, [e]lements, [q]uit (making changes)
~~~

When a model is added one element is created as an example. That element is called "oid" which stands for object identifier. To manage objects we a specific unique key to access them. That is what an object identifier provides. At the SQL level this is analogous to a primary key. 

Elements are the specific data held and managed by your model. You can think of them as the columns in a spreadsheet. We many have lots of objects just as you can have lots of rows in a spreadsheet.

Let's modify our element list. The first thing we'll do is get rid of "oid" and create a new primary key called "garden_id".

Type the letter "e" and press enter.

You should see this new menu.

~~~shell
Enter menu letter and id

	1: oid

Menu [a]dd, [m]odify, [r]emove, [q]uit (making changes)
~~~

This is a list of our model's elements.  We can remove the "oid" by typing the following followed by pressing the enter key.

~~~shell
r oid
~~~

The list should now look like

~~~shell
Enter menu letter and id


Menu [a]dd, [m]odify, [r]emove, [q]uit (making changes)
~~~

Now let's add our "garden_id" by typing the following and then pressing the enter key

~~~shell
a garden_id
~~~

The menu now should look like

~~~shell
Enter menu letter and id

	1: garden_id

Menu [a]dd, [m]odify, [r]emove, [q]uit (making changes)
~~~

We are ready to modify our garden_id and make it the primary key. Type the following and press the entery key.

~~~
m garden_id
~~~

This takes us to the element menu.

~~~
Select menu item to model properties

	1: id garden_id
	2: type input
	3: pattern
	4: attributes name
		title
	5: primary key false

Menu [t]ype, [p]attern, [a]ttributes, [o]bject id flag, [q]uit (making changes)
~~~

Like the model menu we have options to modify the common properties, e.g. type, pattern, and importantly primary key via `[o]object id flag`.

Enter the letter  "o" and press enter. Notice how number five changes. The menu should now look like

~~~shell
Select menu item to model properties

	1: id garden_id
	2: type input
	3: pattern
	4: attributes name
		title
	5: is object id true

Menu [t]ype, [p]attern, [a]ttributes, [o]bject id flag, [q]uit (making changes)
~~~

Because is object id can be either true (on) or false (off) you are not prompted to set its value. Now let's take a look at attributes.
You see "name" and "title". Name is used when this element's HTML is rendered. The name attribute is set by the contents of name. Similarly
the title attribute (used for hover text in a form element) is all set. Managing an element's attributes is done by pressing "a" and then
the enter key. You you've done that you should see the following.

~~~shell
Enter menu letter and id

	1: name
	2: title

Menu [a]dd, [m]odify, [r]emove, [q]uit (making changes)
~~~

Let's change the title text to explain what garden_id is.

Type the following followed by pressing the enter key.

~~~shell
m title
~~~

The menu should now show this.

~~~shell
Enter title's value: 
~~~

Enter a description of the garden_id as follows finishing by presssing the enter key.

~~~shell
This is the unique identifier for the garden model.
~~~

Now that we're done changing the attributes press "q" and then the enter key. This will take us back "up" the element level menu.

~~~shell
Select menu item to model properties

	1: id garden_id
	2: type input
	3: pattern
	4: attributes name
		title
	5: is object id true

Menu [t]ype, [p]attern, [a]ttributes, [o]bject id flag, [q]uit (making changes)
~~~

Typing "q" and enter again will take us up further to the list of elements in our model.

~~~shell
Enter menu letter and id

	1: garden_id

Menu [a]dd, [m]odify, [r]emove, [q]uit (making changes)
~~~

Typing "q" and enter again will show us a summary of our model.

~~~shell
Enter menu letter to manage model

	1: Model Id: garden
	2: Descriptions: This is a model of my garden
	3: Elements: garden_id


Menu: model [d]escription, [e]lements, [q]uit (making changes)
~~~

Typing "q" and pressing enter here will take us to our list of models.

~~~shell
Enter menu letter and id

	1: garden

Menu [a]dd, [m]odify, [r]emove, [q]uit (making changes)
~~~

Typing "q" and pressing enter one last time you should now be prompted to save your model. 


~~~shell
Save before exiting (Y/n)?
~~~

Answer "y" and press enter.

This will save the Newt YAML file and bring us back to the shell prompt.  If you use the more command
you can see the YAML we've created.

~~~shell
more garden.yaml
~~~

This should return something like this.

~~~yaml
#/usr/bin/env newt check
#
# This was generated by rsdoiel on 2024-05-09 with newt version 0.0.8 8453132.
#
applications:
  newtrouter:
    port: 8010
  newtmustache:
    port: 8011
  newtgenerator:
    namespace: garden.yaml
  postgres:
    port: 5432
    dsn: postgres://{PGUSER}:{PGPASSWORD}@localhost:5432/garden.yaml
  postgrest:
    app_path: postgrest
    conf_path: postgrest.conf
    port: 3000
  enviroment:
    - PGUSER
    - PGPASSWORD
models:
  - id: garden
    description: This is a model of my garden
    elements:
      - type: input
        id: garden_id
        attributes:
          name: garden_id
          title: This is the unique identifier for the garden model.
        is_object_id: true
routes:
  - id: garden_create
    request: GET /garden_create
    description: Handle retrieving the webform for garden create
    pipeline:
      - service: POST http://localhost:8011/garden_create
        description: Display a garden for create
  - id: garden_create
    request: POST /garden_create
    description: Handle form submission for garden create
    pipeline:
      - service: POST http://localhost:3000/rpc/garden_create
        description: Access PostgREST API for garden create
      - service: POST http://localhost:8011/garden_create_response
        description: This is an result template for garden create
  - id: garden_update
    request: GET /garden_update/{garden_id}
    description: Handle retrieving the webform for garden update
    pipeline:
      - service: GET http://localhost:3000/rpc/garden_read/{garden_id}
        description: Retrieve garden from PostgREST API before update
      - service: POST http://localhost:8011/garden_update
        description: Display a garden for update
  - id: garden_update
    request: POST /garden_update
    description: Handle form submission for garden update
    pipeline:
      - service: PUT http://localhost:3000/rpc/garden_update/{garden_id}
        description: Access PostgREST API for garden update
      - service: POST http://localhost:8011/garden_update_response
        description: This is an result template for garden update
  - id: garden_delete
    request: GET /garden_delete/{garden_id}
    description: Handle retrieving the webform for garden delete
    pipeline:
      - service: GET http://localhost:3000/rpc/garden_read/{garden_id}
        description: Retrieve garden from PostgREST API before delete
      - service: POST http://localhost:8011/garden_delete
        description: Display a garden for delete
  - id: garden_delete
    request: POST /garden_delete
    description: Handle form submission for garden delete
    pipeline:
      - service: DELETE http://localhost:3000/rpc/garden_delete/{garden_id}
        description: Access PostgREST API for garden delete
      - service: POST http://localhost:8011/garden_delete_response
        description: This is an result template for garden delete
  - id: garden_read
    request: POST /garden_read
    description: Retrieve object(s) for garden read
    pipeline:
      - service: GET http://localhost:3000/rpc/garden_read/{garden_id}
        description: Access PostgREST API for garden read
      - service: POST http://localhost:8011/garden_read
        description: This template handles garden read
  - id: garden_list
    request: POST /garden_list
    description: Retrieve object(s) for garden list
    pipeline:
      - service: GET http://localhost:3000/rpc/garden_list
        description: Access PostgREST API for garden list
      - service: POST http://localhost:8011/garden_list
        description: This template handles garden list
    description: Retrieve object(s) for garden read
    pipeline:
      - service: GET http://localhost:3000/rpc/garden_read/{garden_id}
        description: Access PostgREST API for garden read
      - service: POST http://localhost:8011/garden_read
        description: This template handles garden read
  - id: garden_list
    request: POST /garden_list
    description: Retrieve object(s) for garden list
    pipeline:
      - service: GET http://localhost:3000/rpc/garden_list
        description: Access PostgREST API for garden list
      - service: POST http://localhost:8011/garden_list
        description: This template handles garden list
templates:
  - id: garden_create
    request: /garden_create
    template: garden_create_form.tmpl
    description: Display a garden for create
  - id: garden_create
    request: /garden_create_response
    template: garden_create_response.tmpl
    description: This is an result template for garden create
  - id: garden_update
    request: /garden_update
    template: garden_update_form.tmpl
    description: Display a garden for update
  - id: garden_update
    request: /garden_update_response
    template: garden_update_response.tmpl
    description: This is an result template for garden update
  - id: garden_delete
    request: /garden_delete
    template: garden_delete_form.tmpl
    description: Display a garden for delete
  - id: garden_delete
    request: /garden_delete_response
    template: garden_delete_response.tmpl
    description: This is an result template for garden delete
  - id: garden_read
    request: /garden_read
    template: garden_read.tmpl
    description: This template handles garden read
  - id: garden_list
    request: /garden_list
    template: garden_list.tmpl
    description: This template handles garden list
~~~

While the application would not be useful yet, our garden model only contains an id, we could generate our application's SQL, PostgREST configuration and Mustache templates
with `newt generate garden.yaml`.

As an exercise I recommend adding additionals elements to our garden model such as plant, location and then regenerate the code an notice what changes in the setup.sql and model.sql as well as in our mustache templates.

### Step 3. Generate some code

This last step is one short command but it will generate many files. Type the following command at the shell prompt.

~~~shell
newt generate garden.yaml
~~~

This will generate the following files.

~~~
postgrest.conf
setup.sql
models.sql
garden_create_form.tmpl
garden_create_response.tmpl
garden_delete_form.tmpl
garden_delete_response.tmpl
garden_list.tmpl
garden_read.tmpl
garden_update_form.tmpl
garden_update_response.tmpl
~~~

