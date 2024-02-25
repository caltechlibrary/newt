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
	ModelToPgSQLType = map[string]string{
	}

	ModelToPgSQLFuncType = map[string]string{
	}
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
CREATE ROLE {{namespace}}_anonymous nologin;

--
-- NOTE: The "CREATE ROLE" line is the problem line for
-- checking into your source control system. It contains a secret!
-- **DO NOT** store secrets in your SQL if you can avoid it!
--
DROP ROLE IF EXISTS {{namespace};
CREATE ROLE {{namespace}} NOINHERIT LOGIN PASSWORD '${PASSWORD_GOES_HERE}';
GRANT {{namespace}}_anonymous TO {{namespace}};
`
	tmpl, err := mustache.ParseString(txt)
	if err != nil {
		return err
	}
	data := map[string]string{
		"namespace": namespace,
	}
	return tmpl.FRender(out, data)
}

func pgModels(out io.Writer, namespace string, models []*NewtModel) error {
	// Setup the opening comment, change to the name space and then render models.
	txt := `--
-- Below is the SQL I would noramally check into a project repository.
-- It does not contain secrets. It contains our data models, views,
-- and functions. This defines the behaviors made available through
-- PostgRESTS.
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
	data := map[string]string {
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
--  Need to insure all objects have an object identifier (primary key), creatd and updated timestamps
-- working without interventions for CRUD-L operations.
--  {{fields_def}}
);

--
-- {{namespace}}.{{model}} CRUD-L operations (functions and views)
--

-- {{namspace}}.{{model}}_create is a stored function to create a new object
-- It becomes the end point '/rpc/create_{{model}}'
-- Three 
CREATE OR REPLACE FUNCTION {{namespace}}.create_{{model}}({{fields_ref}})
RETURNS bool LANGUAGE SQL AS $$
  INSERT INTO {{namespace}}.{{model}} ({{fields_ref}})) VALUES ({{fields_ref}});
-- FIXME: Need to return the object key or whole object.
  SELECT true;
$$;

-- {{namespace}}.{{model}}_update is a stored function to update an object.
-- It becomes the end point '/rpc/update_{{model}}'
CREATE OR REPLACE FUNCTION {{namespace}}.update_{{model}}({{fields_ref}})
RETURN book LANGAUGE SQL AS $$
-- FIXME: Need to generste a {{a_field_ref}} = '{{field_value}}' pairing.
   UPDATE {{namespace}}.{{model}} SET {{field_value_list}} WHERE {{filter_for_record}} LIMIT 1;
-- FIXME: Need to return the object key or whole object.
   SELECT true;
$$;

-- {{namespace}}.read_{{model}} will become an end point in PostgREST, '/read_{{model}}'
CREATE OR REPLACE VIEW {{namespace}}.read_{{model}} AS
  SELECT {{fields_refs}} FROM {{namespace}}.{{model}} WHERE {{filter_for_record}} LIMIT 1;

-- {{namespace}}.delete_{{model}} is a stored function to delete an object.
-- It becomes the end point '/rpc/delete_{{model}}'
CREATE OR REPLACE FUNCTION {{namespace}}.delete_{{model}}({{frield_names}}, {{field_values}})
RETURN book LANGAUGE SQL AS $$
   DELETE {{namespace}}.{{model}} WHERE {{filter_for_record}};
   SELECT true;
$$;

-- {{namespace}}.list_{{model}} will become an end point in PostgREST, '/list_{{model}}'
CREATE OR REPLACE VIEW {{namespace}}.list_{{model}} AS
  {{list_view_select}};

`
	tmpl, err = mustache.ParseString(txt)
	if err != nil {
		return err
	}
	data = map[string]string{}
	for _, m := range models {
		data["namespace"] = namespace
		data["model"] = m.Name
		data["fields_def"] = "<<<field definitions goes here>>>" //FIXME: Need a function that can take a model and return SQL field defs.
		data["fields_ref"] = "<<<field references goes here>>>" //FIXME: Need a function that can take a model and return SQL field defs.
		data["filter_for_record"] = "<<<where clause filter goes here>>>"  // FIXME: Need to write the WHERE clause filter to return a specific record.
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
	data = map[string]string {
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
		data = map[string]string {
			"namespace": namespace,
			"model": m.Name,
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
