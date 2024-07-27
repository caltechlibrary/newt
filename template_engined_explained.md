
# Newt Handlebars explained

## Overview

Newt Handlebars is a simple, light weight template engine. It supports the [HandlebarsJS](https://handlebarsjs.com) template language. The way it works is you send data encoded as JSON via a POST to the Newt Handlebars service. Newt Handlebars then takes that data, processes it via a Handlebars template and returns the result.

Newt Handlebars does this through Newt's YAML configuration file. In that file there is a `templates` property where you map the request paths to templates. In Handlebars parlance these are called "views". A view can be made up of a main template with reusable parts included as "partials" templates.

You data should encoded as a JSON object with the following properties.

`body`
: (required) Holds the JSON described object received from the POST and encoded as content type `application/json`.

`vars`
: (optional) Holds any path variables found in the request URL path

`options`
: (optional) Holds a key/value map of strings set in the Newt YAML file in the applications property.

`vocabulary`
: (optional) Holds a object that obtain from an external YAML file.

These object properties are mapped to Handlebars global object. They provide a flexible set of properties for use in creating an HTML page (or other type of text content) from the output of the Newt Handlebars service.

The `vars` and `body` values are determined when a POST request is maid. `options`, `vocabulary` are set when the Newt Handlebars service starts up.  If either `options` or `vocabulary` changes you will not see the changes until you restart the Newt Handlebars service.

## Two ways to run Newt Handlebars

Newt provides two options for running Newt Handlebars.  Newt comes with the Newt Handlebars web service, `newthandlebars` (`newthandlebars.exe` on Windows). This is a standalone web service suitable suitable for running from the command line or via your POSIX systems' init or systemd services. A second more convenient way is to run Newt Handlebars is with the `newt` command (`newt.exe` on Windows). The `newt` command is provided to support a fluid development experience. The one command can perform several actions, e.g. "init", "generate" and "run".  In the following documentation I'll be using the `newt` command to run our Newt Handlebars service.

## Getting started with Newt Handlebars

In this short tutorial we are going to create a web application that says hello. Only Newt Handlebars will be used to implement this service.

### Step 1, create a Newt YAML file

Since I am just focusing on Newt Handlebars I recommend typing in the YAML content below and saving it to a file called "hello.yaml".

~~~yaml
applications:
  newtHandlebars:
    port: 8011
  options:
    default_name: There
    place: Planet Earth!
templates:
  - request: /hello
    template: hello.tmpl
  - request: "/hello/{someplace_else}"
    template: hello.tmpl
~~~

This simple YAML file describes how to run Newt Handlebars, as a web service on port 8011. It also describes the path Newt will listen on to run the "hello.tmpl".

You can use `newt check hello.yaml` to verify you've entered the Newt YAML correctly. Here's an example of the running the command and it's output.

~~~shell
newt check hello.yaml
~~~

The output should look something like

~~~text
WARNING: hello.yaml has no models defined
Newt Handlebars configured, port set to 8011
2 Handlebars Templates are defined
http://localhost8011/hello points at hello.tmpl
http://localhost8011/hello/{someplace_else} points at hello.tmpl
~~~

### Step 2, create our "hello.tmpl" template file

Let's look at a simple hello world template in Handlebars configured to work with the Newt Handlebars service.

~~~html
<!DOCTYPE html>
<html lang="en-US">
  <body>
    {{#if body.name}}Hi {{body.name}}{{else}}Hi {{options.default_name}}{{/if}}
    {{#if vars.someplace_else}}from {{vars.someplace_else}}!{{else}} from {{options.place}}!{{/if}}
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
- The native JavaScript version is available in Deno and can be compiled into a binary service via the `deno compile` command so is suitable for inclusion in Newt.
- Since I am writing the service I can keep the requires to a minimum, i.e. [Use Newt's YAML file syntax](newthandlebars.1.md#newt_config_file).

See [newthandlebars](newthandlebars.1.md) manual page for details.
