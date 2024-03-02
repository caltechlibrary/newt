
# Newt Mustache explained

## Overview

Newt Mustache is a simple, light weight template engine. It supports the [Mustache](https://mustache.github.io) template language. The way it works is you send data encoded as JSON via a POST to the Newt Mustache service. Newt Mustache then takes that data, processes it via a Mustache template and returns the result.

Newt Mustache does this through Newt's YAML configuration file. In that file there is a `templates` property where you map the request paths to a template. You may also define variables embedded in the request path that will be made available to the related template.  Similarly Newt YAML's `.applications.options` is made available to your tempalte. Finally you can can also define "vocabulary" in an other YAML file that can be made available for processing with your template.

In the Mustache template you have the following objects availabe.

`body`
: Holds the JSON described object recieved from the POST and encoded as content type `application/json`.

`options`
: Holds a key/value map of strings set in the Newt YAML file in the applications property.

`vocabulary`
: Holds a object that obtain from an external YAML file.

`vars`
: Holds any path variables found in the request URL path

This provides a flexible set of properties for use in creating an HTML page (or other type of text document) from the output of the Newt Mustache service.

Aside from `vars` and `body` the rest of the information (e.g. templates, options, vocabulary) are read in when Newt Mustache is started. It does not read from disk after that and doesn't write to disk either. Newt Mustache only communicates via HTTP on localhost. If you change the options, update your vocabulary file you will need to restart Newt Mustache (or `newt`) to see those changes.

## Two ways to run Newt Mustache

Newt provides two options for running Newt Mustache.  Newt comes with the Newt Mustache web service, `newtmustache` or `newtmustache.exe`. This is a standalone web service suitable to be run from init or systemd on your POSIX system. A second more convient way is to run Newt Mustache with the `newt` command. This is what I use when developing a Newt application. With `newt` you can run Newt Router, Newt Mustache and PostgREST if configured in your Newt YAML's applications properties.

In this document we will assume you are running Newt Mustache using the `newt` command.

## Step 1, you need a Newt YAML file

We want to start with configuration. We need to configure both the applications property and templates property to run Newt Mustache. Here's a simple example you should type in as save as "hello.yaml".

~~~yaml
applications:
  newtmustache:
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

This simple YAML file descriptions how to run Newt Mustache, as a web service on port 8011. It also describes the path Newt will listen on to run the "hello.tmpl".

Let's look at a simple hello world template in Mustache configured to work with the Newt Mustache service.

~~~html
<!DOCTYPE html>
<html lang="en-US">
  <body>
    {{#body.name}}Hi {{body.name}}{{/body.name}}    
    {{^body.name}}Hi {{options.default_name}}{{/body.name}}
    {{#vars.someplace_else}}from {{vars.someplace_else}}!{{/vars.someplace_else}}
    {{^vars.someplace_else}}from {{options.place}}!{{/vars.someplace_else}}
  </body>
</html>
~~~

If you have these typed in then we're ready to run `newt`.

~~~shell
newt hello.yaml
~~~

Open another terminal or shell session to test[^1].  

[^1]: On Linux, macOS or Windows using LSW I use [tmux](https://github.com/tmux/tmux/wiki) to allow splitting the window. This let's me run the `newt` command in one and use cURL to test the template rendering in another.

I use [cURL](https://curl.se) to test my templates. Let's see what happens when we send an empty POST to our hello URL defined in our "hello.yaml" templates section.

~~~shell
curl --data '{}' 'http://localhost:8011/hello'
~~~

The curl command above sends a POST (implied by the `--data` option) using a content type of "application/json". Newt Mustache only knows how to work with JSON data. curl takes care of setting that content type using the `--data`. Normally a POST from a web form is encluded using urlencoding. But the web has evolved since the 1990s and most API now produce output encoded as JSON (or have the option to do so). Newt Mustache is designed to support this behavior when it processes requests.

When you run the curl command you above should get back the HTML markup from the template and in the body element see the message "Hi There from Planet Earth!".

What about that second template request path we defined? The one that contained a variable called `someplace_else` in the path.

~~~shell
curl --data '{}' 'http://localhost:8011/hello/the%20Moon'
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
curl --data '{"name":"Maxine"}' http://localhost:8011/hello
~~~

This time we should get back similar HTML but in the body element we should see the message "Hi Maxine from Planet Earth!".

Now let's try the second path.

~~~shell
curl --data '{"name":"Maxine"}' 'http://localhost:8011/hello/Mars'
~~~

Now you should see "Hi Maxine from Mars!".

That's the basic idea of Newt Mustache. Let me cover some other situations you might encounter when developing your Newt Mustache templates.

If you use an HTTP beside POST you will get back an HTTP error message. If you use a URL path not defined in your templates you will get back an HTTP error message. These should be 404 type HTTP error message.

If you're POSTing to a defined URL and still running into a problem (say a template or data issue). You will get back an HTTP error. The easiest way to get insight into what is happening is to run the `newt` command using the `--verbose` option. This will output a allot of debug information which hopefully will help you find the problem in your template or in your data.

### Why Mustache templates?

- Mustache is a widely support template language with support include Python, JavaScript, and Go (languages used regularly at Caltech Library)
- Since a Go package provides Mustache template I only need to write a light weight web service to wrap it
- Since I am writing the service I can keep the requires to a minimum, i.e. [Use Newt's YAML file syntax](newtmustache.1.md#newt_config_file).

See [newtmustache](newtmustache.1.md) manual page for details.
