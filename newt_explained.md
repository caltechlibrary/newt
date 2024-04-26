
# Newt's `newt` command

The core of Newt is the command `newt`. It handles the following for your Newt based project.

1. Generate an initialize Newt YAML file through an interactive dialog
2. Generate and manage your models through an interactive dialog
3. Create an manage the initial routing in your Newt Application based on your models
4. Generate a PostgREST configuration file, SQL files as well as Mustache templates
5. Manage the generated templates in your Newt YAML file
6. Run Newt Router, Newt Mustache and PostgREST so you can test and develop your application

The `newt` command is a development tool. It simplifies standing up a Newt application built around the Newt Router, Newt Mustache template engine and Postgres+PostgREST.  A `newt` generated application will provide the basic create, read, update, delete and list operations you need for managing the metadata described by your models.

## The Newt YAML file

The Newt YAML file is made of four major properties

applications
: Holds configuration information for PostgREST, Newt Router and Newt Mustache template engine

models
: This holds the data models your applicaiton will implement, this property is inspired by the [GitHub YAML issue template syntax]()

routes
: These are the request Newt Router will manage they are descriptions of HTTP method and path along with any data pipeline processing needed to respond to HTTP request

templates
: This holds the configuration for the Newt Mustache template engine. The template engine accepts POST requests and associates Mustache templates with a given request path.  Mustache templates support the concept of "partial" templates, these are also specified in relationship to the request path and primary template association. Templates are read in with Newt Mustache starts and are not re-read until you restart Newt Mustache.

## Initializing your Newt YAML file

The `newt` command supports an "init" action. This will create or update a Newt YAML file. It does this through a dialog of questions which you provide answers for. Based on the answers then the Newt YAML file will be modified to reflect your answers. Before the Newt YAML file is written to disk you will be prompted if it is OK to save the name file. If the file already exists than the old one will be renamed with a ".bak" file extension.  If a ".bak" file exist it will be replaced.

Here is an example of creating a "app.yaml" file using the "init" action.

~~~shell
newt init app.yaml
~~~

NOTE: You can press "control C" to quit the "init" action without writing the YAML to disk.

## Modeling your data

The `newt` command provides an "model" action. This lets you add, modify or remove a data model from your Newt YAML file.  Like the "init" action it will provide you with the opportunity to review changes before writing them to disk.

The modeling process is more complex than the "init" action. You may have model than one model, you may have many properties per model.  As a result the dialog between you and `newt model` is separated into stages. First you are asked to manage the model(s) by either adding or removing them.  You can then "modify" a model where it will allow you to add or remove a list of model properties.  Each model will have at least one property, the "oid" property. "oid" stands for object indentifier. By default an object identifier is represented as a UUID. In string form these can be long but they allow for millions of object to be managed.  If you prefer a shorter string representaition for your object identifiers these can be selected but you cannot remove the object identifier from the model. It remains required.

### Top level dialog

The top level dialog lets you perform one of four different actions. Add a model, Modify a model, Remove a model and quit.  The top level dialog will list any models hat are defined.

When you choose "add model" you will then be prompted for a model identifier. This identifier must start with an alphabetical character followed by one or more alphanumeric characters or underscore characters. E.g. "my_bird_list" is an example of a valid model name while "2bots!" is not a valid name. In Postgres the model id will be used as the Postgres table name. The model id will also will be used to form the Mustache template names by the Newt code generator. Model names must be unique inside your application.

When you see one or models listed you less an integer to the left.  If you type the integer (or type our the model name) you then will be asked if you want to "modify" or "remove" the model.  Type "m" will let you modify the model, Type "r" to remove it.  When you remove a model you will be taken back to the top level dialog. If you type "m" to modify the model you will be taken to the modify model dialog.

Typing "s" will save the current state of the models to your Newt YAML file.  Typing "q" will save then "quit" the modeler completely. Typing "c" for cancel will exit the modeler without saving changes.

### Modifying a model

The modify model dialog will show you a list of properties associated with the model.  A model must always have an "oid" (i.e. object identifier) property. While you can't remove the object identifier property you can modify it's type. Currently an "oid" defaults to a UUID (native to Postgres) but you many choose to use a [shortuuid/v4](https://github.com/lithammer/shortuuid) stored as a 22 character string, [Mongo BON ObjectID](https://pkg.go.dev/github.com/mongodb/mongo-go-driver/bson/objectid) stored as a 24 character string.

The modifying model view allows you to add a property, modify a property or remove a property.

In the modify model dialog you can choose add a property, modify a property or remove a property. To add a property type "a" and you will be prompted to provide a property identifiers. This, like model identifiers, needs to start with an alphabetical characters followed by one or more alphanumeric characters or underscore. Press enter and you are taken to the property dialog. By default the added property is of "input" type. You can modify it to define different property attributes such as different types of property.

If you want to modify or remove a property you can specify the property by entering the integer to the left of the property name or my typing the property. If you follow this by "m" then you'll be taken to the property modification dialog. If you type "r" it'll remove the property and if you press enter you will be taken back to the property list for the current model.

Typing "s" will save the current model. Typing "q" will save then take you back to the top level dialog. Typing "c" for cancel will return to the top level dialog without saving the changes.


### Property dialog

Modifying property has a similar interface to the models dialog and the modify model dialog. It presents you with a list of current attributes. It differs in that when you ave the options of select the specific attribute of the property to modify. Note that what is presented is tied to the type value of the property. The type corresponds to the basic HTML input element types defined on at [MDN](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/Input). It also includes some customized input specific to the gallery, library, archive and museum metadata domain. These will be expressed in as standard HTML elements the Mustache templates but can be further vetted from within Postgres via Python's idutils package.

When you press entere when listing the property attributes it accepts the current settings and take you back to the modify model dialog. Typing "c" for cancel will tak take you back to he modify model dialog without changing the property's attributes.

If you type "m" to modify your model you when will be shown a list of properties associated with the model. There will always be at least one property, the "oid" or object identifier property.  The object identifier property is special. It can be modified but NOT removed. All models have an "oid" property.  By default the "type" is UUID for the object identifier. UUID are a good way to support object identification in a modern SQL database like Postgres.  They have only one drawback. The way we express the value of a UUID tends to be as a long string. This can be unwieldly for URLs. You can choice an alternate identifier type that results in sorter or numeric identifiers. The shorter identifiers limit the total number of objects you can manage but typically these is not a problem (e.g. Caltech Library's Authors repository has 100,000 objects, someday it might have 200,000 objects. Each object could still have a unique identifiers of with a six or seven character string). Finally you can also use an integer value which is incremented with each added object. This tends to be short but comes at the expense of limitting you to a single database instance in many cases.

Typing "s" will save the property settings. Tying "q" will save then exit the property dialog. Typing "c" for cancel will exit the property dialog without making changings.


## Support input types

FIXME: Need to write up the basic supported input types and how they related to presentations in HTML and SQL.

