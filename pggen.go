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
-- drop database if exists {{namespace}};
-- create database {{namespace}};

-- Make sure we are in the {{namespace}} namespace/database
\c {{namespace}}

--
-- Setup a Postgres "schema" (a.k.a. namespace) for
-- working with PostgREST
--
drop schema if exists {{namespace}} cascade;
create schema {{namespace}};

--
-- The following additional steps are needed for PostgREST access
-- are {{namespace}} schema and database.
--
drop role if exists {{namespace}}_anonymous;
create role {{namespace}}_anonymous nologin;

--
-- NOTE: The "CREATE ROLE" line is the problem line for
-- checking into your source control system. It contains a secret!
-- **DO NOT** store secrets in your SQL if you can avoid it!
--
drop role if exists {{namespace}};
create role {{namespace}} noinherit login password '{{secret}}';
grant {{namespace}}_anonymous to {{namespace}};
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
set search_path to {{namespace}}, public;

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
-- {{namespace}}.{{model}} table definition
--
drop table if exists {{namespace}}.{{model}};
create table {{namespace}}.{{model}}
(
	oid uuid primary key default gen_random_uuid(), 
	created timestamp with time zone default now(),
	updated timestamp with time zone default now(), 
	object jsonb
);

--
-- {{namespace}}.{{model}} CRUD-L operations (functions and views)
--

-- {{namespace}}.create_{{model}} is a stored function to create a new object
-- It takes the new object without oid, created, updated. It will return the
-- assigned uuid.
--
-- This becomes the end point '/rpc/create_{{model}}' in PostgREST
--
drop function if exists {{namespace}}.create_{{model}} (new_object jsonb);
create or replace function {{namespace}}.create_{{model}} (
    new_object jsonb
) returns uuid
language sql
as $$
    insert into {{namespace}}.{{model}}
        (oid, created, updated, object)
    values
        (gen_random_uuid(), now(), now(), new_object)
    returning oid
$$
;

--
-- {{namespace}}.update_{{model}} is a stored function to update an object.
-- NOTE: It replaces the whole object!
--
-- It takes the UUID of the object to update and the model's content. You do
-- not need to provide created and updated attributes as those are managed by
-- this function.  The value returned is the revised JSON including all attributes.
--
-- It becomes the end point '/rpc/update_{{model}}'
--
drop function if exists {{namespace}}.update_{{model}} (id uuid, new_object jsonb);
create or replace function {{namespace}}.update_{{model}} (
    id uuid,
    new_object jsonb
)
returns jsonb
language plpgsql
as $$
declare
    updated_object jsonb;
begin
  update {{namespace}}.{{model}}
  set updated = now(), object = new_object
  where oid = id;
  select (jsonb_build_object(
    'oid', oid,
    'created', created,
    'updated', updated
    ) || object) 
  into updated_object
  from {{namespace}}.{{model}}
  where oid = id;
  return updated_object;
end;
$$
;

--
-- {{namespace}}.read_{{model}} will retrieve the object stored flattening it to include
-- the fields modelled but also include the attributes oid, created, updated enforce by these procedures.
--
-- It become an end point in PostgREST, '/rpc/read_{{model}}'
--
drop function {{namespace}}.read_{{model}} (IN id uuid, OUT obj jsonb);
create or replace function {{namespace}}.read_{{model}} (
  IN id uuid,
  OUT obj jsonb
)
language sql
as $$
  select (jsonb_build_object(
    'oid', oid,
    'created', created,
    'updated', updated
    ) || object) as obj
  from {{namespace}}.{{model}}
  where oid = id
  limit 1;
$$
;

--
-- {{namespace}}.delete_{{model}} is a stored function to delete an object.
-- It takes the object UUID (oid) as a parameter. If successful it returns the
-- oid deleted.
--
-- It becomes the end point '/rpc/delete_{{model}}'
--
drop function if exists {{namespace}}.delete_{{model}} (id uuid);
create or replace function {{namespace}}.delete_{{model}} (
    id uuid
) returns uuid
language sql
as $$
    delete from {{namespace}}.{{model}}
    where oid = id
    returning id
$$
;

--
-- {{namespace}}.list_{{model}} lists object by descending updated timestamp.
-- It does not take any parameters.
--
-- It will become an end point in PostgREST, '/rpc/list_{{model}}'
--
drop function {{namespace}}.list_{{model}} ();
create or replace function {{namespace}}.list_{{model}} ()
returns table (
  obj jsonb
)
language sql
as $$
  select (jsonb_build_object(
    'oid', oid,
    'created', created,
    'updated', updated
    ) || object) as obj
  from {{namespace}}.{{model}}
$$
;

`
	tmpl, err = mustache.ParseString(txt)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "-- DEBUG Code generated for for %d models\n", len(models))
	data = map[string]string{}
	for _, m := range models {
		data["namespace"] = namespace
		data["model"] = m.Id
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

grant usage on schema {{namespace}} to {{namespace}}_anonymous;

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
grant select, insert on {{namespace}}.{{model}} to {{namespace}}_anonymous;
grant execute on function {{namespace}}.create_{{model}} to {{namespace}}_anonymous;
grant execute on function {{namespace}}.update_{{model}} to {{namespace}}_anonymous;
grant execute on function {{namespace}}.delete_{{model}} to {{namespace}}_anonymous;
grant select on {{namespace}}.read_{{model}} to {{namespace}}_anonymous;
grant select on {{namespace}}.list_{{model}} to {{namespace}}_anonymous;

`
	tmpl, err = mustache.ParseString(txt)
	if err != nil {
		return err
	}
	for _, m := range models {
		data = map[string]string{
			"namespace": namespace,
			"model":     m.Id,
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
