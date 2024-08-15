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

// pgSetup renders setup SQL for namespace.
func pgSetup(out io.Writer, namespace string) error {
	txt := `--
-- Following would not normally be include in a project's Git repository.
-- It contains a secret.  What I would recommend is writing a short
-- shell script that could generate this in a file, use that, then
-- checking in the shell script to version control since the secret
-- would not be stored in the file!
--

-- Uncomment these two lines if you have used not used Postgres' createdb yet.
-- drop database if exists {{namespace}};
-- create database {{namespace}};

-- Make sure we are in the {{namespace}} namespace/database
-- \c {{namespace}}

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
drop role if exists {{namespace}}_authenticator;
create role {{namespace}}_authenticator noinherit login password '{{secret}}';
grant {{namespace}}_anonymous to {{namespace}}_authenticator;
`
	tmpl, err := mustache.ParseString(txt)
	if err != nil {
		return err
	}
	data := map[string]string{
		"namespace": namespace,
		"secret":    "__change_me_password_goes_here__",
	}
	return tmpl.FRender(out, data)
}

func pgModels(out io.Writer, namespace string, models []*Model) error {
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
-- identifier: This is a short unique identifier for the object
-- created: This is the date the object was created (timestamp format)
-- updated: This is a timestamp of when an object changed
-- object: This is a jsonb column holding our object
--
-- This simplification allows for simplier code in the 2nd Newt Prototype.
-- In principle the model could be mapped directly into column defs but
-- this is unnecessary in the prototype so skipped for now.
--
-- Make sure we are in the birds namespace/database
-- \c {{namespace}}

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
-- Make sure our namespace is first in the search path
set search_path to {{namespace}}, public;

--
-- {{namespace}}.{{model}} table and index definitions
--
drop table if exists {{namespace}}.{{model}};
create table {{namespace}}.{{model}}
(
	identifier uuid primary key default gen_random_uuid(), 
	created timestamp with time zone default now(),
	updated timestamp with time zone default now(), 
	object jsonb
);

--
-- {{namespace}}.{{model}} CRUD-L operations (functions and views)
--

-- {{namespace}}.{{model}}_create is a stored function to create a new object
-- It takes the new object without identifier, created, updated. It will return the
-- assigned uuid.
--
-- This becomes the end point '/rpc/{{model}}_create' in PostgREST
--
--drop function if exists {{namespace}}.{{model}}_create (new_object jsonb);
create or replace function {{namespace}}.{{model}}_create (
    new_object jsonb
) returns jsonb
language plpgsql
as $$
declare
  	id uuid;
  	ts timestamp;
    updated_object jsonb;
begin
    id := gen_random_uuid();
	ts := now();
	-- create record
    insert into {{namespace}}.{{model}}
        (identifier, created, updated, object)
    values
        (id, ts, ts, new_object);
	-- return created record
  	select (jsonb_build_object(
    	'identifier', identifier,
    	'created', created,
    	'updated', updated
    	) || object) 
  	into updated_object
  	from {{namespace}}.{{model}}
  	where identifier = id;
  	return updated_object;
end;
$$
;

--
-- {{namespace}}.{{model}}_update is a stored function to update an object.
-- NOTE: It replaces the whole object!
--
-- It takes the UUID of the object to update and the model's content. You do
-- not need to provide created and updated attributes as those are managed by
-- this function.  The value returned is the revised JSON including all attributes.
--
-- It becomes the end point '/rpc/{{model}}_update'
--
--drop function if exists {{namespace}}.{{model}}_update (id uuid, new_object jsonb);
create or replace function {{namespace}}.{{model}}_update (
    id uuid,
    new_object jsonb
)
returns jsonb
language plpgsql
as $$
declare
    updated_object jsonb;
begin
	-- update record
	update {{namespace}}.{{model}}
	set updated = now(), object = new_object
	where identifier = id;
	-- return updated record
	select (jsonb_build_object(
    	'identifier', identifier,
    	'created', created,
    	'updated', updated
    	) || object) 
	into updated_object
	from {{namespace}}.{{model}}
	where identifier = id;
	return updated_object;
end;
$$
;

--
-- {{namespace}}.{{model}}_read will retrieve the object stored flattening it to include
-- the fields modelled but also include the attributes identifier, created, updated enforce by these procedures.
--
-- It become an end point in PostgREST, '/rpc/{{model}}_read'
--
--drop function if exists {{namespace}}.{{model}}_read (IN id uuid, OUT obj jsonb);
create or replace function {{namespace}}.{{model}}_read (
  IN id uuid,
  OUT obj jsonb
)
language sql
as $$
  select (jsonb_build_object(
    'identifier', identifier,
    'created', created,
    'updated', updated
    ) || object) as obj
  from {{namespace}}.{{model}}
  where identifier = id
  limit 1;
$$
;

--
-- {{namespace}}.{{model}}_delete is a stored function to delete an object.
-- It takes the object UUID (identifier) as a parameter. If successful it returns the
-- identifier deleted.
--
-- It becomes the end point '/rpc/{{model}}_delete'
--
--drop function if exists {{namespace}}.{{model}}_delete (id uuid);
create or replace function {{namespace}}.{{model}}_delete (
    id uuid
) returns uuid
language plpgsql
as $$
begin
    delete from {{namespace}}.{{model}}
    where identifier = id;
    return id;
end;
$$
;

--
-- {{namespace}}.list_{{model}} lists object by descending updated timestamp.
-- It does not take any parameters.
--
-- It will become an end point in PostgREST, '/rpc/{{model}}_list'
--
--drop function if exists {{namespace}}.{{model}}_list ();
create or replace function {{namespace}}.{{model}}_list ()
returns table (
  obj jsonb
)
language sql
as $$
  select (jsonb_build_object(
    'identifier', identifier,
    'created', created,
    'updated', updated
    ) || object)::jsonb as obj
  from {{namespace}}.{{model}}
  order by updated desc;
$$
;

`
	tmpl, err = mustache.ParseString(txt)
	if err != nil {
		return err
	}
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
grant select, insert, update, delete on {{namespace}}.{{model}} to {{namespace}}_anonymous;
grant execute on function {{namespace}}.{{model}}_create to {{namespace}}_anonymous;
grant execute on function {{namespace}}.{{model}}_read to {{namespace}}_anonymous;
grant execute on function {{namespace}}.{{model}}_update to {{namespace}}_anonymous;
grant execute on function {{namespace}}.{{model}}_delete to {{namespace}}_anonymous;
grant execute on function {{namespace}}.{{model}}_list to {{namespace}}_anonymous;

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

func pgModelsTest(out io.Writer, namespace string, models []*Model) error {
	txt := `--
-- {{namespace}}_test.sql tests the models described in {{namespace}}.yaml
--
-- WARNING: It is only an approximate test as it DOES NOT use your actual models!!!
--
-- You should modify this file to meat the needs of testing your Postgres models.
--
-- \c {{namespace}}

select 'Testing {{namespace}}.{{model}}_create' as msg limit 1;
--
-- Test {{namespace}}.{{model}}_create can execute.
--
select {{namespace}}.{{model}}_create(jsonb_build_object(
   'test_number', {{test_no}},
   'test_notes', 'Initial Object Creation'
));

select 'Testing {{namespace}}.{{model}}_read' as msg limit 1;
--
-- Test if {{namespace}}.{{model}}_read can execute.
--
with t as (
  select identifier
  from {{namespace}}.{{model}}
  order by created desc
  limit 1
) select {{namespace}}.{{model}}_read(t.identifier) from t limit 1;


select 'Testing {{namespace}}.{{model}}_update' as msg limit 1;
-- 
-- Test if {{namespace}}.{{model}}_update can execute.
-- 
with t as (
  select identifier
  from {{namespace}}.{{model}}
  order by created desc
  limit 1
) select {{namespace}}.{{model}}_update(t.identifier, jsonb_build_object(
   'test_number', {{test_no}},
   'test_notes', 'Updated Object Action'
)) from t limit 1;


--
-- Test listing all records created, read and updated.
--
select {{namespace}}.{{model}}_list();

--
-- Test listing all records delete and list records.
--
select 'Testing {{namespace}}.{{model}}_delete' as msg limit 1;
-- 
-- Test if {{namespace}}.{{model}}_delete can execute.
-- 
with t as (
  select identifier
  from {{namespace}}.{{model}}
  order by updated desc
  limit 1
) select {{namespace}}.{{model}}_delete(t.identifier) from t;

--
-- Test listing all records after delete.
--
select {{namespace}}.{{model}}_list();

`
	tmpl, err := mustache.ParseString(txt)
	if err != nil {
		return err
	}
	for i, m := range models {
		data := map[string]string{
			"namespace": namespace,
			"model":     m.Id,
			"test_no":   fmt.Sprintf("%d", i),
		}
		if err := tmpl.FRender(out, data); err != nil {
			return err
		}
	}
	return nil
}
