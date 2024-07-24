%newt(1) user manual | 0.0.8 2024-07-24
% R. S. Doiel
% 2024-07-24 57d3074

# NAME

newt

# SYNOPSIS

newt [OPTIONS] NEWT_YAML_FILE

# DESCRIPTION

newt provides a handlebarjs template engine as a web service

NEW_YAML_FILE is the YAML file for your Newt application with the templates
property and runtime configuration.

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
: set the port number, default 3032


# EXAMPLE

Running newt using the newt YAML file named "my_app.yaml".

~~~shell
newthandlebars my_app.yaml
~~~


