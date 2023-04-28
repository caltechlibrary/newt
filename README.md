
# Newt Demo

Newt, a "new take" on building web applications with less coding.
Newt demonstrates building a traditional web application using
Postgres 15, PostgREST 11 and a front end web server. The demo focuses
on the development implecations and should be construed as what is
required in a public facing or production deployment.

- [Building with Postgres and PostgREST](building-with-postgres-postgrest.md) discusses the approach taken to create the [birds](birds/) demo
- Extas for setting up a developer environment
    - [setup-birds.bash](setup-birds.bash), a bash script that generates the contents of demo's running code
    - [Multipass basics](multipass-basics.md), multipass runs Ubuntu VM which can be used to run the demo
    - [newt-init.yaml](newst-init.yaml) provides the configuration for a multipass based VM to run the demo
    - [setup-developer-account.bash](setup-developer-account.bash) is a bash script that displays Postgres commands for setting up a super user account for development.

