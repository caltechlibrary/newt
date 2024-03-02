--
-- Following would not normally be include in a project's Git repository.
-- It contains a secret.  What I would recommend is writing a short
-- shell script that could generate this in a file, use that, then
-- checking in the shell script to version control since the secret
-- would not be stored in the file!
--

-- Make sure we are in the birds namespace/database
\c birds

--
-- Setup a Postgres "schema" (a.k.a. namespace) for
-- working with PostgREST
--
drop schema if exists birds cascade;
create schema birds;

--
-- The following additional steps are needed for PostgREST access
-- are birds schema and database.
--
drop role if exists birds_anonymous;
create role birds_anonymous nologin;

--
-- NOTE: The "CREATE ROLE" line is the problem line for
-- checking into your source control system. It contains a secret!
-- **DO NOT** store secrets in your SQL if you can avoid it!
--
drop role if exists birds;
create role birds noinherit login password 'my_secret_password';
grant birds_anonymous to birds;

