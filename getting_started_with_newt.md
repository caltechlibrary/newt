---
title: Getting Started with Newt
created: 2024-08-15
draft: true
---

# Getting Started with Newt

By R. S. Doiel

## Introduction

Newt is a rapid application development set of tools. It was inspired by the need to quickly
create metadata curations tools at Caltech Library and based on my experience working for USC's Center for
Scholarly before 2015. It is aimed at empower Librarians and library staff who need to create metadata curation
tools for themselves or projects. 

Before getting started with Newt I am making some assumptions that you understand the following concepts.

- the role of an identifier in uniquely identifying a metadata record
- the concept of a data model
- at a very high level understand of how the web works (e.g. url, web service, web browser)
- how to install software on your computer
- how to start a "terminal" session on your Mac, Windows or Linux computer

With these knowledge prerequists Newt lets an individual quickly build a minimal metadata curation web application
which they can run on their desktop or laptop computer and access with their favorite web browser.

Newt provides an opportunity to explore the following concepts

- the basic web application request and response
- database management system for persisting data
- creating an application as a composition of web services

If you are inclined you're include Newt also provides you a good platform to exploring HTML5, CSS and JavaScript.


## Setting up Newt

Newt doesn't standalone. It builds on existing software. If you are developing with Newt you will
need to have the following software installed in addition to Newt.

- Postgres >= 16
- PostgREST >= 12
- Deno >= 1.45.5
- On POSIX systems (to installed Newt)
    - unzip
    - curl
    - bash
- On Windows
    - Powershell

### Installing Newt

If you are running macOS, Linux or other POSIX system you should be able to install Newt with the following command

~~~shell
curl -fsSL https://caltechlibrary.github.io/newt/installer.sh | sh
~~~

On Windows you use Powershell with the following command.

~~~pwsh
irm https://caltechlibrary.github.io/newt/installer.ps1 | iex
~~~

These command will retrieve the latest version of Newt and install them on your machine in your account.

### Installing Postgres

Go to the Postgres website, <https://www.postgresql.org/>, and follow their instructions to download and
install PosgreSQL on your computer.

### Install PostgREST

Go to the PostgREST website, <https://postgrest.org>, and follow the instructions to download
and install PostgREST.

### Installing Deno

Deno can be installed from the Deno Land website, <https://deno.land>. Like Newt you can use your shell or 
computer. For macOS or Linux I use the following.

~~~shell
curl -fsSL https://deno.land/install.sh | sh
~~~

For Windows 11 I use the following.

~~~pwsh
irm https://deno.land/install.ps1 | iex
~~~

## Creating your first Newt App

FIXME: come up with a simple example applicaition to build. Walk through the steps of model, build, and run.

## Where to go from here

FIXME: Discuss what you can do as "next" steps in your first app.

- manually enhance templates to include custom CSS and JavaScript
- Adding additional models
- Going from a "localhost" system to a staff or public facing one




