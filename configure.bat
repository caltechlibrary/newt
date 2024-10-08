rem
rem This script setups of the model definition for development. It will import mimimal modules needed to compile newt.exe.
rem
@echo off
echo "Checking for Go >= 1.22"
go version
echo "Checking for Pandoc >= 3.1"
pandoc --version
echo "Checking for PageFind >= 1.1.0"
pagefind --version
echo "Checking for Postgres >= 16"
psql --version
echo "Checking for PostgREST >= 12"
postgrest --version
echo "Checking for jq >= 1.6"
jq --version
echo "Review version info for each software checked."
echo "Checking for deno >= 1.44"
deno --version
@echo on
