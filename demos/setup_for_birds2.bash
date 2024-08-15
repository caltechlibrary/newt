#!/bin/bash

dropdb birds
createdb birds
psql -d birds -c '\i setup.sql'
pandoc server &
postgrest postgrest.conf
