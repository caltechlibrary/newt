%newt(1) user manual | 0.0.9 2024-07-29
% R. S. Doiel
% 2024-07-29 efd5fe1

# NAME

newt

# SYNOPSIS

newt [OPTIONS] NEWT_YAML_FILE

# DESCRIPTION

newt provides a handlebarjs template engine as a web service.
newt's YAML configuration file holds a mapping of URL paths
to handlebar templates. If you make an HTTP GET request the service
will return the unrendered template associated with the URL path.
If you make an HTTP POST it takes the POST the data provided will
be used as the template content. Normally the data you provide is in
the form of JSON. You should set the HTTP header "content-type"
to "application/json" as part of your POST request.

# OPTIONS


help
: display help

license
: display license

version
: display version

debug
: turn on debug logging

port
: set the port number, default 0


# EXAMPLE

Running newt using the newt YAML file named "app.yaml".

~~~shell
newt app.yaml
~~~


