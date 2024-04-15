rem
rem This script setups of the model definition for development. It will import mimimal modules needed to compile newt.exe.
rem


rem
rem Check all of irdmtools is installed
rem
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
if exist go.mod del go.mod
if exist go.sum del go.sum
go mod init "github.com/caltechlibrary/newt"
go mod tidy
