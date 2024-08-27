package newt

import (
	//"fmt"
	"bytes"
	"path"
	"testing"
)

// TestPgGenerators test the generation of setup.sql, models.sql and models_test.sql.
func TestPgGenerators(t *testing.T) {
	fName := path.Join("testdata", "birds.yaml")
	ast, err := LoadAST(fName)
	if err != nil {
		t.Errorf("Missing %q, aborting test", fName)
		t.FailNow()
	}
	if ast.Applications == nil {
		t.Errorf("Missing .Applications from %q, aborting test", fName)
		t.FailNow()
	}
	postgres := ast.GetApplication("postgres")
	if postgres == nil {
		t.Errorf("Missing Postres in .Applications from %q", fName)
		t.FailNow()
	}
	if postgres.Namespace == "" {
		t.Errorf("Missing Postgres.Namespace in .Applications from %q, aborting test", fName)
		t.FailNow()
	}
	expected := "birds"
	got := postgres.Namespace
	if expected != got {
		t.Errorf("expected namepsace %q, got %q", expected, got)
	}
	namespace := "birds"
	buf := bytes.NewBuffer([]byte{})
	pgSetup(buf, namespace)
	data := buf.Bytes()
	if len(data) == 0 {
		t.Errorf("expected data from buf after pgSetup(buf, %q)", namespace)
	}
	expectedBytes := []byte(`--
-- Following would not normally be include in a project's Git repository.
-- It contains a secret.  What I would recommend is writing a short
-- shell script that could generate this in a file, use that, then
-- checking in the shell script to version control since the secret
-- would not be stored in the file!
--

-- Uncomment these two lines if you have used not used Postgres' createdb yet.
-- drop database if exists birds;
-- create database birds;

-- Make sure we are in the birds namespace/database
-- \c birds

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
drop role if exists birds_authenticator;
create role birds_authenticator noinherit login password '__change_me_password_goes_here__';
grant birds_anonymous to birds_authenticator;
`)
	if len(expectedBytes) != len(data) {
		t.Errorf("expected src\n%s\ngot\n%s\n", expectedBytes, data)
		t.Errorf("expected len %d, got len %d", len(expectedBytes), len(data))
	}
}
