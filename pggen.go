package newt

import (
	"fmt"
	"io"

	// 3rd Party Templates
	"github.com/cbroglie/mustache"
)

//
// This file contains code for generating PostgreSQL dialect of SQL.
//

var (
	ModelToPgSQLType = map[string]string{}

	ModelToPgSQLFuncType = map[string]string{}
)

// pgSetup renders setup SQL for namespace.
func pgSetup(out io.Writer, namespace string) error {
	txt := `--
-- Following would not normally be include in a project's Git repository.
-- It contains a secret.  What I would recommend is writing a short
-- shell script that could generate this in a file, use that, then
-- checking in the shell script to version control since the secret
-- would not be stored in the file!
--

-- Uncomment these two lines if you have used Postgres' createdb yet.
-- DROP DATABASE IF EXISTS {{namespace}};
-- CREATE DATABASE EXISTS {{namespace}};

-- Make sure we are in the {{namespace}} namespace/database
\c {{namespace}}

--
-- Setup a Postgres "schema" (a.k.a. namespace) for
-- working with PostgREST
--
DROP SCHEMA IF EXISTS {{namespace}} CASCADE;
CREATE SCHEMA {{namespace}};

--
-- The following additional steps are needed for PostgREST access
-- are {{namespace}} schema and database.
--
DROP ROLE IF EXISTS {{namespace}}_anonymous;
CREATE ROLE {{namespace}}_anonymous NOLOGIN;

--
-- NOTE: The "CREATE ROLE" line is the problem line for
-- checking into your source control system. It contains a secret!
-- **DO NOT** store secrets in your SQL if you can avoid it!
--
DROP ROLE IF EXISTS {{namespace}};
CREATE ROLE {{namespace}} NOINHERIT LOGIN PASSWORD '{{secret}}';
GRANT {{namespace}}_anonymous TO {{namespace}};
`
	tmpl, err := mustache.ParseString(txt)
	if err != nil {
		return err
	}
	data := map[string]string{
		"namespace": namespace,
		"secret": "__change_me_password_goes_here__",
	}
	return tmpl.FRender(out, data)
}

func pgModels(out io.Writer, namespace string, models []*NewtModel) error {
	// Setup the opening comment, change to the name space and then render models.
	txt := `--
-- Below is the SQL I would noramally check into a project repository.
-- It does not contain secrets. It contains our data models, views,
-- and functions. This defines the behaviors made available through
-- PostgREST.
--
-- Newt Generator creates a simple SQL schema leveraging the JSONB column
-- support in Postgres.  There are four fields inspired by a simplification
-- of RDM records.
--
-- oid: This is a short unique identifier for the object
-- created: This is the date the object was created (timestamp format)
-- updated: This is a timestamp of when an object changed
-- object: This is a jsonb column holding our object
--
-- This simplification allows for simplier code in the 2nd Newt Prototype.
-- In principle the model could be mapped directly into column defs but
-- this is unnecessary in the prototype so skipped for now.
--
-- Make sure we are in the birds namespace/database
\c {{namespace}}
-- Make sure our namespace is first in the search path
SET search_path TO {{namespace}}, public;

--
-- Data Models
--
`
	tmpl, err := mustache.ParseString(txt)
	if err != nil {
		return err
	}
	data := map[string]string{
		"namespace": namespace,
	}
	if err := tmpl.FRender(out, data); err != nil {
		return err
	}

	// Now for the table definitions.
	txt = `
--
-- This file contains the create statements for our bird table.
--
-- DROP TABLE IF EXISTS {{namespace}}.{{model}};
CREATE TABLE {{namespace}}.{{model}} (
	oid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	object JSONB
);

CREATE OR REPLACE FUNCTION {{model}}_update_updated()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated = current_timestamp;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER {{model}}_updated_trigger
BEFORE UPDATE ON {{model}}
FOR EACH ROW
EXECUTE FUNCTION {{model}}_update_updated();

--
-- {{namespace}}.{{model}} CRUD-L operations (functions and views)
--

-- {{namspace}}.{{model}}_create is a stored function to create a new object
-- It becomes the end point '/rpc/create_{{model}}'
-- It turns the generator UUID
CREATE OR REPLACE FUNCTION {{namespace}}.create_{{model}}(object)
RETURNS UUID, TIMESTAMPZ, TIMESTAMPZ, JSONB LANGUAGE SQL AS $$
    INSERT INTO {{namespace}}.{{model}} (object)) VALUES (object)
    RETURNING {{model}}.oid, {{model}}.created, {{model}}.updated, {{model}}.object
$$;

-- {{namespace}}.{{model}}_update is a stored function to update an object.
-- It becomes the end point '/rpc/update_{{model}}'
CREATE OR REPLACE FUNCTION {{namespace}}.update_{{model}}(oid, object)
RETURN UUID, TIMESTAMPZ, TIMESTAMPZ, JSONB LANGAUGE SQL AS $$
    UPDATE {{namespace}}.{{model}} SET {{model}}.object = object WHERE {{model}}.oid = oid 
    RETURNING {{model}}.oid, {{model}}.created, {{model}}.updated, {{model}}.object
$$;

-- {{namespace}}.read_{{model}} will become an end point in PostgREST, '/read_{{model}}'
CREATE OR REPLACE VIEW {{namespace}}.read_{{model}} AS
  SELECT oid, created, updated, object FROM {{namespace}}.{{model}} WHERE oid = ? LIMIT 1;

-- {{namespace}}.delete_{{model}} is a stored function to delete an object.
-- It becomes the end point '/rpc/delete_{{model}}'
CREATE OR REPLACE FUNCTION {{namespace}}.delete_{{model}}(oid)
RETURN UUID, TIMESTAMPZ, TIMESTAMPZ, JSONB LANGAUGE SQL AS $$
   DELETE {{namespace}}.{{model}} WHERE {{model}}.oid = oid
   RETURn oid, created, updated, object;
$$;

-- {{namespace}}.list_{{model}} will become an end point in PostgREST, '/list_{{model}}'
-- It returns a list sorted by descending updated date.
CREATE OR REPLACE VIEW {{namespace}}.list_{{model}} AS
  SELECT oid, created, updated, object FROM {{namespace}}.{{model}} ORDER BY updated DESC;
`
	tmpl, err = mustache.ParseString(txt)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "-- DEBUG Code generated for for %d models\n", len(models))
	data = map[string]string{}
	for _, m := range models {
		data["namespace"] = namespace
		data["model"] = m.Name
		if err := tmpl.FRender(out, data); err != nil {
			return err
		}
	}

	txt = `
--
-- PostgREST access and controls.
--

-- Since our Postgres ROLE and SCHEMA exist and our models may change how
-- we want PostgREST to expose our data via JSON API we GRANT or
-- revoke role permissions here to match our model.
GRANT USAGE ON SCHEMA {{namespace}} TO {{namespace}}_anonymous;
-- NOTE: The generated functions for create, update and delete do not
-- implement an account requirement. In a production application you
-- should modify these functions to check for authorization before
-- allowing changes to be made.
`
	tmpl, err = mustache.ParseString(txt)
	if err != nil {
		return err
	}
	data = map[string]string{
		"namespace": namespace,
	}
	if err := tmpl.FRender(out, data); err != nil {
		return err
	}

	txt = `
GRANT SELECT, INSERT ON {{namespace}}.{{model}} TO {{namespace}}_anonymous;
GRANT EXECUTE ON FUNCTION {{namespace}}.create_{{model}} TO {{namespace}}_anonymous;
GRANT EXECUTE ON FUNCTION {{namespace}}.updated_{{model}} TO {{namespace}}_anonymous;
GRANT EXECUTE ON FUNCTION {{namespace}}.delete_{{model}} TO {{namespace}}_anonymous;
GRANT SELECT ON {{namespace}}.read_{{model}} TO {{namespace}}_anonymous;
GRANT SELECT ON {{namespace}}.list_{{model}} TO {{namespace}}_anonymous;

`
	tmpl, err = mustache.ParseString(txt)
	if err != nil {
		return err
	}
	for _, m := range models {
		data = map[string]string{
			"namespace": namespace,
			"model":     m.Name,
		}
		if err := tmpl.FRender(out, data); err != nil {
			return err
		}
	}
	return nil
}

func pgModelsTest(out io.Writer, models []*NewtModel) error {
	return fmt.Errorf("pgModelsTest() not implemented yet")
}
