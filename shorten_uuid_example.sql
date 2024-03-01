--
-- Example of generating a short uuid for URL usage.
--
with t as (
	select oid
	from ornothology.sighting
	order by created desc
	limit 1
) select encode(uuid_send(t.oid::uuid), 'base64') as short_uuid from t;

--
-- shorten_uuid
--
create or replace function ornothology.shorten_uuid(id uuid)
returns varchar
language plpgsql
as $$
declare
   out_id varchar;
begin
   out_id := encode(uuid_send(id::uuid), 'base64');
   return out_id;
end;
$$
;

--
-- test shorten_uuid function
--
with t as (
	select oid
    from ornothology.sighting
    order by created desc
    limit 1
) select ornothology.shorten_uuid(t.oid) from t;

--
-- unshorten_uuid
--
create or replace function ornothology.unshorten_uuid(short_id varchar)
returns uuid
language plpgsql
as $$
declare
  l_id uuid;
begin
	 l_id := substring(decode(short_id, 'base64')::text from 3)::uuid;
	 return l_id;
end;
$$
;

--
-- test shorten_uuid and unshorten_uuid functions
--
with t as (
	select 
		oid::uuid AS long_id,
		encode(uuid_send(oid::uuid),'base64')::varchar as short_id
    from ornothology.sighting
    order by created desc
    limit 1
) select 
	t.long_id,
	ornothology.shorten_uuid(t.long_id::uuid)::varchar AS r1,
	t.short_id,
	ornothology.unshorten_uuid(t.short_id::varchar)::uuid AS r2
from t;

