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
	cfg, err := LoadConfig(fName)
	if err != nil {
		t.Errorf("Missing %q, aborting test", fName)
		t.FailNow()
	}
	if cfg.Applications == nil {
		t.Errorf("Missing .Applications from %q, aborting test", fName)
		t.FailNow()
	}
	if cfg.Applications.NewtGenerator == nil {
		t.Errorf("Missing .Applications.NewtGenerator from %q, aborting test", fName)
		t.FailNow()
	}
	if cfg.Applications.NewtGenerator.Namespace == "" {
		t.Errorf("Missing .Applications.NewtGenerator.Namespace from %q, aborting test", fName)
		t.FailNow()
	}
	expected := "birds"
	got := cfg.Applications.NewtGenerator.Namespace
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

-- Uncomment these two lines if you have used Postgres' createdb yet.
-- DROP DATABASE IF EXISTS birds;
-- CREATE DATABASE EXISTS birds;

-- Make sure we are in the birds namespace/database
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
CREATE ROLE birds_anonymous NOLOGIN;

--
-- NOTE: The "CREATE ROLE" line is the problem line for
-- checking into your source control system. It contains a secret!
-- **DO NOT** store secrets in your SQL if you can avoid it!
--
DROP ROLE IF EXISTS birds;
CREATE ROLE birds NOINHERIT LOGIN PASSWORD '__change_me_password_goes_here__';
GRANT birds_anonymous TO birds;
`)
    if len(expectedBytes) != len(data) {
	t.Errorf("expected len %d, got len %d", len(expectedBytes), len(data))
    }

    //fmt.Printf("DEBUG buf -> %s\n", buf)	
}
