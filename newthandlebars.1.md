%newt(1) user manual | 0.0.8 2024-07-24
% R. S. Doiel
% 2024-07-24 0961fb8

# NAME

newt

# SYNOPSIS

newt [OPTIONS] NEWT_YAML_FILE

# DESCRIPTION

newt provides a handlebarjs template engine as a web service

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


