
# Multipass basics

Multipass is a Ubuntu specific VM manager making it easy to develop
applications targeting Ubuntu Server and Desktop environments. It is
available on Linux, macOS and Windows. See [multipass.run](https://multipass.run) for installation instructions.

## Basic commands

The [newt-demo](https://github.com/caltechlibrary/newt) provides a [newt-init.yaml]() to setup a Multipass environment for exploring a web application approach focusing on SQL, HTML, JavaScript using PostgreSQL 15, PostgREST 11 and a front end web server.

Creating our development VM.

~~~
multipass launch --name newt --cloud-init newt-init.yaml
~~~

Seeing if our VM is running

~~~
multipass list
~~~

Copy [setup-birds1.bash](setup-birds1.bash), [setup-birds2.bash](setup-birds2.bash) or [setup-birds3.bash](setup-birds3.bash) to our VM to save some type setting up our
demo.

~~~
multipass transfer setup-birds1.bash newt:.
multipass transfer setup-birds2.bash newt:.
multipass transfer setup-birds3.bash newt:.
~~~

Access our VM's shell

~~~
multipass shell newt
~~~

Stopping our VM

~~~
multipass stop newt
~~~

Removing the VM completely

~~~
multipass delete newt && multipass purge
~~~


