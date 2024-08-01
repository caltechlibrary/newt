
# the Newt Template Engine

## Overview

Newt provides a simple stateless template engine as web service. It plays nice as a final stage in a data pipeline. The Newt configuration file maps URL paths to template names. The service only recognizes POST which contain optional JSON payloads. If a POST request is recieved the the rendered version using the JSON data submitted is returned.

The template engine is based on the [raymond](https://github.com/aymerick/raymond) Go package. Raymond implemented Handlerbars 3 support.  The template engine does not support Raymond's Handlebar functions since the engine doesn't compile Go code.

## Template objects

The following three objects are available in the templates.

`body`
: (required) Holds the JSON described object received from the POST and encoded as content type `application/json`.

`vars`
: (optional) Holds any path variables found into the request URL path and specified in the request path in the Newt configuration file.

`document`
: (optional) Holds any values set via the template configuration document, passed from the environment via the Newt configuration file.

## Two ways to run the template engine

Newt provides two options for running the template engine.  Newt comes the `nte` web service (`nte.exe` on Windows). This is a standalone web service suitable suitable for running from the command line or via your POSIX systems' init or systemd service manager. A second more convenient way is to run the template engine is with the `newt` command (`newt.exe` on Windows). The `newt` command is provided to support a fluid development experience. This one command can perform several actions, e.g. "init", "generate" and "run".  In the following documentation I'll be using the `newt` command to run our Newt Handlebars service.

## Getting started with Newt Handlebars

In this short tutorial we are going to create a web application that says hello. Only Newt Handlebars will be used to implement this service.

### Step 1, create a Newt YAML file

Since I am just focusing on Newt Handlebars I recommend typing in the YAML content below and saving it to a file called "hello.yaml".

~~~yaml
applications:
  template_engine:
    port: 8011
    base_dir: views
    ext_name: .hbs
    partials: partials
templates:
  - request: /hello
    template: hello
    document:
      default_name: There
      place: Planet Earth!
  - request: "/hello/{someplace_else}"
    template: hello
    document:
      default_name: There
      place: Planetoid Pluto!
~~~

This simple YAML file describes how to run Newt's template engine. The web service listens on port 8011. When the web service receives a POST
to the path "/hello" it will envoke the "views/hello.hbs" template.

You can use `newt check hello.yaml` to verify you've entered the Newt YAML correctly. Here's an example of the running the command and it's output.

~~~shell
newt check hello.yaml
~~~

The output should look something like

~~~text
template engine will listen on port 8011
templates are located in "views"
template extension is set to ".hbs"
2 template(s) paths mapped
~~~

Notice that we're using one template in multiple paths.

### Step 2, create our "hello.tmpl" template file

Let's look at a simple hello world template in Handlebars configured to work with the Newt Handlebars service.

~~~html
<!DOCTYPE html>
<html lang="en-US">
  <body>
    {{#if body.name}}Hi {{body.name}}{{else}}Hi {{document.default_name}}{{/if}}
    {{#if vars.someplace_else}}from {{vars.someplace_else}}!{{else}} from {{document.place}}!{{/if}}
  </body>
</html>
~~~

### Step 3, run our web application

Now we are ready to use `newt` to "run" our hello application.

~~~shell
newt run hello.yaml
~~~

Open another terminal or shell session to test[^1].

[^1]: On Linux, macOS or Windows using LSW I use [tmux](https://github.com/tmux/tmux/wiki) to allow splitting the window. This let's me run the `newt` command in one and use cURL to test the template rendering in another.

I use [cURL](https://curl.se) to test my templates. Let's see what happens when we send an empty POST to our hello URL defined in our "hello.yaml" templates section.

~~~shell
curl -X POST -H 'content-type: application/json' --data '{}' 'http://localhost:8011/hello'
~~~

The curl command above sends a POST (implied by the `--data` option) using a content type of "application/json". Newt Handlebars only knows how to work with JSON data. curl takes care of setting that content type using the `--data`. Normally a POST from a web form is included using urlencoding. But the web has evolved since the 1990s and most API now produce output encoded as JSON (or have the option to do so). Newt Handlebars is designed to support this behavior when it processes requests.

When you run the curl command you above should get back the HTML markup from the template and in the body element see the message "Hi There from Planet Earth!".

What about that second template request path we defined? The one that contained a variable called `someplace_else` in the path.

~~~shell
curl -X POST -H 'content-type: application/json' --data '{}' 'http://localhost:8011/hello/the%20Moon'
~~~

We should get a similar response as before but the "Planet Earth" should be replaced by "the Moon". If you've made it this far we know both templates paths work.  Now let's try sending our JSON object with cURL.

We want to send the following JSON.

~~~json
{
    "name": "Maxine"
}
~~~

This will populate will let us see the name of "Maxine" rather than word "There" in the body of the HTML.

~~~shell
curl -X POST -H 'content-type: application/json' --data '{"name":"Maxine"}' http://localhost:8011/hello
~~~

This time we should get back similar HTML but in the body element we should see the message "Hi Maxine from Planet Earth!".

Now let's try the second path.

~~~shell
curl -X POST -H 'content-type: application/json' --data '{"name":"Maxine"}' 'http://localhost:8011/hello/Mars'
~~~

Now you should see "Hi Maxine from Mars!".

That's the basic idea of Newt Handlebars. Let me cover some other situations you might encounter when developing your Newt Handlebars templates.

If you use an HTTP method beside POST you will get back an HTTP error message. If you use a URL path not defined in your templates you will get back an HTTP error message. These should be 404 type HTTP error message.

If you're POSTing to a defined URL and still running into a problem (say a template or data issue). You will get back an HTTP error. The easiest way to get insight into what is happening is to run the `newt` command using the `--verbose` option. This will output a allot of debug information which hopefully will help you find the problem in your template or in your data.

### Why Handlebars templates?

- Handlebars is a widely support template language with implementations in many languages including JavaScript and Python.
- It is available browser side (it's written in JavaScript) and server side in environments like [Deno](https://deno.land) and [NodeJS](https://nodejs.org/en)

