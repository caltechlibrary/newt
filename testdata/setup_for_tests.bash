#!/bin/bash

export DB_NAME="birds"
export DB_USERNAME="birds"
export DB_PASSWORD="my_secret_password"

export NEW_ROUTES="birds_routes.csv"

dropdb birds
createdb birds
psql -d birds -c '\i birds_setup.sql'
pandoc server &
postgrest postgrest.conf
