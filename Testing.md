
# Testing setup

Tests are written using Go's testing package.  You need more than that to run meaningful tests since Newt is a microserve integeracting with other microservices.

# My development setup

1. I maintain a `test.evn` file which sets environment variables enabling access to PostgREST and Postgres
2. I start a new shell session, change to `testdata` and run `setup-for-tests.bash`
3. In another shell session I run `make test` and monitory the results


