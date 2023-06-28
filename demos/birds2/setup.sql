--
-- Following would not normally be include in a project's Git repository.
-- It contains a secret.  What I would recommend is writing a short
-- shell script that could generate this in a file, use that, then
-- checking in the shell script to version control since the secret
-- would not be stored in the file!
--

-- Make sure we are in the birds database
\c birds

--
-- Setup a Postgres "schema" (a.k.a. namespace) for
-- working with PostgREST
--
DROP SCHEMA IF EXISTS birds CASCADE;
CREATE SCHEMA birds;

--
-- The following additional steps are needed for PostgREST access
-- are birds schema and database.
--
DROP ROLE IF EXISTS birds_anonymous;
CREATE ROLE birds_anonymous nologin;

--
-- NOTE: The "CREATE ROLE" line is the problem line for
-- checking into your source control system. It contains a secret!
-- **DO NOT** store secrets in your in your repository if you can avoid it!
--
DROP ROLE IF EXISTS birds;
CREATE ROLE birds NOINHERIT LOGIN PASSWORD 'my_secret_password';
GRANT birds_anonymous TO birds;

