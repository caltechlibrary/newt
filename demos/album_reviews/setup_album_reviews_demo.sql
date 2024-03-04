--
-- Following would not normally be include in a project's Git repository.
-- It contains a secret.  What I would recommend is writing a short
-- shell script that could generate this in a file, use that, then
-- checking in the shell script to version control since the secret
-- would not be stored in the file!
--

-- NOTE: the database album_reviews needs to exist. Create the database using the
-- postgres command `createdb album_reviews`.

-- Make sure we are in the album_reviews namespace/database
\c album_reviews

--
-- Setup a Postgres "schema" (a.k.a. namespace) for
-- working with PostgREST
--
drop schema if exists album_reviews cascade;
create schema album_reviews;

--
-- The following additional steps are needed for PostgREST access
-- are album_reviews schema and database.
--
drop role if exists album_reviews_anonymous;
create role album_reviews_anonymous nologin;

--
-- WARNING: The "CREATE ROLE" line is the problem line for
-- checking into your source control system. It contains a secret!
-- **DO NOT** store secrets in your SQL if you can avoid it!
--
drop role if exists album_reviews;
create role album_reviews noinherit login password 'my_secret_password';
grant album_reviews_anonymous to album_reviews;

--
-- Below is the SQL I would noramally check into a project repository.
-- It does not contain secrets. It contains our data models, views,
-- and functions. This defines the behaviors made available through
-- PostgRESTS.
--

-- Make sure we are in the album_reviews namespace/database
-- Make sure our namespace is first in the search path
set search_path to album_reviews, public;

--
-- Data Models
--

--
-- This file contains the create statements for our bird table.
--
drop table if exists album_reviews.albums;
create table album_reviews.albums
(
  name varchar(255),
  review text,
  rating integer,
  created timestamp default now()
);

--
-- Data Views and Behaviors.
--

-- interesting_albums_list will become an end point in PostgREST, '/interestting_albums_list'
create or replace view album_reviews.interesting_albums_list as
  select name, review, rating
  from album_reviews.albums order by rating desc, name asc;

-- add_a_review is a stored procedure and will save a new album review.
-- It becomes the end point '/rpc/add_a_review'
create or replace function album_reviews.add_a_review(
  new_album varchar, 
  new_review text, 
  newt_rating integer
) returns jsonb
language plpgsql
as $$
declare
  new_object jsonb;
begin
  -- insert our new review
  insert into album_reviews.albums (name, review, rating)
         values (new_album, new_review, new_rating);
  -- retrieve our new review and return it.
  select 
    jsonb_build_object(
        'album', album,
        'review', review,
        'rating', rating)
    into new_object
    from album_reviews.albums
    where name = new_album
    order by created desc, name asc
    limit 1;
  return new_object;
end; 
$$
;

--
-- PostgREST access and controls.
--

-- Since our Postgres ROLE and SCHEMA exist and our models may change how
-- we want PostgREST to expose our data via JSON API we GRANT or
-- revoke role permissions here to match our model.
grant usage on schema album_reviews to album_reviews_anonymous;
-- NOTE: We are allowing insert because this is a demo and we are not
-- implementing an account login requirement. Normally I would force
-- a form update via a SQL function or procedure only.
grant select on album_reviews.interesting_albums_list to album_reviews_anonymous;
grant select, insert on album_reviews.albums to album_reviews_anonymous;
grant execute on function album_reviews.add_a_review to album_reviews_anonymous;


