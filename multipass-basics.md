
# Multipass basics

Multipass is a Ubuntu specific VM manager making it easy to develop
applications targetting Ubuntu Server and Desktop environments. It is
available on Linux, MacOS and Windows. See [multipass.run](https://multipass.run) for installation instructions.

## Basic commands

The [newstack-demo](https://github.com/caltechlibrary/newstack-demo) provides a [newstack-init.yaml]() to setup a Multipass environment for exploring a web appllication approach focusing on SQL, HTML, JavaScript using PostgreSQL 15, PostgREST 11 and a front end web server.

Creating our developmenet VM.

~~~
multipass launch --name newstack --cloud-init newstack-init.yaml
~~~

Seeing if our VM is running

~~~
multipass list
~~~

Copy [setup-birds.bash](setup-birds.bash) to our VM to save some type setting up our
demo.

~~~
multipass transfer setup-birds.bash newstack:.
~~~

Access our VM's shell

~~~
multipass shell newstack
~~~

Stopping our VM

~~~
multipass stop newstack
~~~

Removing the VM completely

~~~
multipass delete newstack && multipass purge
~~~


