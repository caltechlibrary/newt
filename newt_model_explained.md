
# Newt's `newt model`

The model action does most of the heavy lifting short of actual code generation. If you save the model(s) you will cause the models property to be populated but also related propertes under the routes and templates properties.  Like `newt config` if you do not provide a file name the command will assume you're working with "app.yaml".  

It is possible to run `newt model` before running `newt config`.  Whichever order you run them in they both need to be run before you try `newt run`.


## Modeling your data

The `newt` command provides an "model" action. This lets you add, modify or remove a data model from your Newt YAML file.  Like the "config" action it will provide you with the opportunity to review changes before writing them to disk.

The modeling process is more complex than the "config" action. You may have model than one model, you may have many properties per model.  As a result the dialog between you and `newt model` is separated into stages. First you are asked to manage the model(s) by either adding or removing them.  You can then "modify" a model where it will allow you to add or remove a list of model properties.  Each model will have at least one property, the "oid" property. "oid" stands for object indentifier. By default an object identifier is represented as a UUID. In string form these can be long but they allow for millions of object to be managed.  If you prefer a shorter string representaition for your object identifiers these can be selected but you cannot remove the object identifier from the model. It remains required.

NOTE: You can press "control C" to quit the "config" action without writing the YAML to disk.

### Top level menu

The top level menu lets you perform one of four different actions. Add a model, Modify a model, Remove a model and quit.  The top level model will list any models hat are defined already.

When you choose "add model" you will then be prompted for a model identifier. This identifier must start with an alphabetical character followed by one or more alphanumeric characters or underscore characters. E.g. "my_bird_list" is an example of a valid model name while "2bots!" is not a valid name. The reason for these restrictions is the model id is used when generating SQL as well as when we generate Mustache Templates. An example is in Postgres the model id will be used as the Postgres table name. Model names must be unique inside your application for the same reason that Postgres tables require uniqueness in a given Postgres database.

### Adding a model

Let's fireup `newt model` and add a model called "garden".

~~~shell
newt model
~~~

You should now see the following menu.

~~~shell
Enter menu letter and id


Menu [a]dd, [m]odify, [r]emove, [q]uit (making changes)
~~~

Type the letter "a" and press the enter key.

You will see 

~~~shell
Enter model id to add:
~~~

Enter "my_bird_list" without the quotes and then press the enter key.

You should now see

~~~shell
Enter menu letter and id

	1: my_bird_list

Menu [a]dd, [m]odify, [r]emove, [q]uit (making changes)
~~~

You've successfully created an empty model called "my_bird_list".

### Modifying a model

The modify model menu will show you a list of properties associated with the model.  A model must always have an "oid" (i.e. object identifier) property. While you can't remove the object identifier property you can modify it's type. Currently an "oid" defaults to a UUID (native to Postgres) but you many choose to use a [shortuuid/v4](https://github.com/lithammer/shortuuid) stored as a 22 character string, [Mongo BON ObjectID](https://pkg.go.dev/github.com/mongodb/mongo-go-driver/bson/objectid) stored as a 24 character string.

The modifying model view allows you to add a property, modify a property or remove a property.

In the modify model dialog you can choose add a property, modify a property or remove a property. To add a property type "a" and you will be prompted to provide a property identifiers. This, like model identifiers, needs to start with an alphabetical characters followed by one or more alphanumeric characters or underscore. Press enter and you are taken to the property dialog. By default the added property is of "input" type. You can modify it to define different property attributes such as different types of property.

If you want to modify or remove a property you can specify the property by entering the integer to the left of the property name or my typing the property. If you follow this by "m" then you'll be taken to the property modification dialog. If you type "r" it'll remove the property and if you press enter you will be taken back to the property list for the current model.

Typing "s" will save the current model. Typing "q" will save then take you back to the top level dialog. Typing "c" for cancel will return to the top level dialog without saving the changes.


### Property dialog

Modifying property has a similar interface to the models dialog and the modify model dialog. It presents you with a list of current attributes. It differs in that when you ave the options of select the specific attribute of the property to modify. Note that what is presented is tied to the type value of the property. The type corresponds to the basic HTML input element types defined on at [MDN](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/Input). It also includes some customized input specific to the gallery, library, archive and museum metadata domain. These will be expressed in as standard HTML elements the Mustache templates but can be further vetted from within Postgres via Python's idutils package.

When you press enter when listing the property attributes it accepts the current settings and take you back to the modify model dialog. Typing "c" for cancel will   take you back to he modify model dialog without changing the property's attributes.

If you type "m" to modify your model you when will be shown a list of properties associated with the model. There will always be at least one property, the "oid" or object identifier property.  The object identifier property is special. It can be modified but NOT removed. All models have an "oid" property.  By default the "type" is UUID for the object identifier. UUID are a good way to support object identification in a modern SQL database like Postgres.  They have only one drawback. The way we express the value of a UUID tends to be as a long string. This can be unwieldy for URLs. You can choice an alternate identifier type that results in sorter or numeric identifiers. The shorter identifiers limit the total number of objects you can manage but typically these is not a problem (e.g. Caltech Library's Authors repository has 100,000 objects, someday it might have 200,000 objects. Each object could still have a unique identifiers of with a six or seven character string). Finally you can also use an integer value which is incremented with each added object. This tends to be short but comes at the expense of limiting you to a single database instance in many cases.

Typing "s" will save the property settings. Tying "q" will save then exit the property dialog. Typing "c" for cancel will exit the property dialog without making changes.


## Support input types

FIXME: Need to write up the basic supported input types and how they related to presentations in HTML and SQL.

