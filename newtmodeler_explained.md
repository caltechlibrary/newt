
# Newt Modeler explained (draft concept)

## Overview

Newt Modeler provides an interactive way for generating and managing your models without directly typing YAML. It is an interactive command line program in the spirit of `newt init`.  When started it list the existing models (if any are defined) and lets you add, modify, remove models in our YAML file. It also will updated the default routes and templates defined in the YAML file to matching the model's name.  

## Two ways to run Newt Modeler

In this short tutorial we are going to create a web application that says describes the contents of a garden.  Like the Newt Router, Newt Generator and Newt Mustache the Modeler can be used directly with the `newtmodeler` command or via the `newt` command. In the tutorial we'll be using the `newt` command since this is a development activity and that is the tool designed to help develop Newt applications.

## Tutorial on the Newt Modeler

### Step 1, create a "garden.yaml" Newt YAML file

We can create a new YAML file using our standard `newt init` command. For this tutorial just accept the defaults.

~~~shell
newt init garden.yaml
~~~

### Step 2, defining and managing models for "garden.yaml"

The first stage of the modeler is to present a list of existing models. It then gives you options of "add", "remove" "modify", "init", and "quit". Launch our modeler with `newt model garden.yaml`

~~~shell
newt model garden.yaml
~~~

By default the command `newt init` created an "garden" model. We're not going to use this model so let's remove it.

~~~shell
remove garden
~~~

When that action is done the "garden" model will be removed from the "garden.yaml" file. It will also remove the default routes and templates in the "garden.yaml" generated previously from the `newt init` command. 

When "remove garden" is complete there will be no models. We need to add one. Let's add a model called "plant".  we want to add a new model. 

~~~shell
add plant
~~~

Now we should see a single model called "plant" in the list. This has done three things. It has updated the models list and created a minimal model with a "oid" (i.e. object identifier) as the sole property of the model.  The "oid" property is the identifier used to manage the object(s) created using this model. It is required by Newt.  The second thing that happens when we create a new model is routes and templates properties are added to our "garden.yaml" file.  These can be used to generate the files we need to run our Newt application.

An object with only an "oid" isn't very useful. We want to modify our model to include some addititional properties like "plant_name", "plant_description", and "plant_location" in our garden box, "plant_type" (i.e. vegie, fruit, tuber) and "date_planted".  


In next step is we want to "modify" our model to have these properties.

### Step 3, let's define some new properties.

~~~
modify plant
~~~

You should now see a list of propertes in our "plant" model. Right now it only has the "oid" property defined. This property is required and can't be removed.

Let's add the plant name to our model by typing "add name" at the prompt.

#### plant_name

~~~shell
add plant_name
~~~

This will add default property called "plant_name" to the model "plant".  We will want to modify this property. We do that with the "modify plant_name" command.

~~~shell
modify plant_name
~~~

This will step through the "plant_name" property elements. For each element you have an opportunity to accept (pressing enter), modify or remove the element value(s). The first element presented is optional. It is the description of the property being defined.  It is free text but if populated it will be used in generated comments in SQL code and as help text in Mustache HTML template code.

Because "description" is free format text if you have the EDITOR environment variable set it will open your editor and allow you to edit the text using your preferred editor. If you don't have that value set then you will type text similar to how you type text at the shell prompt.  Your edit ability will be limited.  If you are editing in your EDITOR then when you save the text and exit the value will be updated in the model. If you quit the edit without saving the value will remain the same as before.

After description you will be prompted for the property's type. Like in GitHub YAML Issue Template Syntax this corresponds to the HTML notation of an input element's type.  The default value is "input[type=text]".  You can use any valid HTML 5 input element type expression[^1].  Example, if you wanted to have a "date" input you would change this to `input[type=date]`. 

[^1]: See [MDN Element Input documents](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/Input#input_types) for details.

The "plant_name" property will use the default "input" so press enter to accept. 

You will then be prompt for attributes and validations. These are key value pairs. If you press enter then you are accepting the (currently empty) list presented.
if you wanted to remove a key/value pair you would type in the key and it would be removed. If you wanted to add a key value pair you would type in the key, a space, then
the value. The key must start with a letter or underscore followed by one or more alphanumeric or underscore characters.  The value is set to anything after first space following the key before you press enter.

For "plant_name" when prompted for "attributes" just press enter to accept. When prompted for "validations" type in "required true" (without the quotes) and press enter.

~~~shell
required true
~~~

#### plant_description

Next we need to add the "plant_description" property. It's going to be a text field. Unlike "plant_name" it will not be required. We'll enter a description and then accept the defaults for the other elements of the property.

#### plant_type

Plant type is the first element we will create that has an expected set of values.  Our plants in our garden will be one of several types "veggie", "tuber", "spice", "fruit" and "flower". The HTML element type that make sense for this type of limited vocabulary "select" element. Replace "input" with "select" and press enter. This will then prompt for a set of key corresponds to what the "value" will contain in the option element inside "select" element. The value you enter corresponds to the inner HTML of the option element.  

As example "select" expressed as HTML would look like.

~~~html
<select>
   <option value="veggie">Veggie</option>
   <option value="tuber">Tuber</option>
   <option value="spice">Spice</option>
   <option value="fruit">Fruit</option>
   <option value="flower">Flower</option>
</select>
~~~

For building our select box we need to enter the follow text.

~~~shell
veggie Veggie
tuber Tuber
spice Spice
fruit Fruit
flower Flower
~~~

When the list looks right press enter. We're not going to use the attributes or validations in this property.

#### plant_location

This element describes the location in our garden of the plant. For the purposes our tutorial we're going to describes our garden as a grid. We can use letters for the "columns" of locations and numbers for the "rows" in our garden. This will result is location descriptions like "A1" or "B15". We are going to use the attribute property to enforce the values submitted. 

Enter a description, accept the default input type of "input[type=text]". When you get to attributes we need to enter a "pattern" along with an expression that validates the letterand number combination describing locations in our garden box.

~~~shell
pattern [a-z|A-Z][0-9]
~~~

### date_planted

This is a is the date you planted the plant in your garden box. You should create it using the "input[type=cal]" type. It isn't required because you may have not planted it yet.

When you have entered this once you can press enter again to accept the model you've made you should return to the model list.

### Step 4, save our model's properties

### Step 5, updating the routes and templates

### Step 6, ready to test

See [newtmodeler](newtmodeler.1.md) manual page for details.

