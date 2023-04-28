DROP SCHEMA IF EXISTS birds CASCADE;
CREATE SCHEMA birds;

--
-- This file contains the create statements for our bird table.
--
DROP TABLE IF EXISTS birds.sighting;
CREATE TABLE birds.sighting
(
  bird_name VARCHAR(255),
  place TEXT,
  sighted DATE
);

-- bird_view will become an end point in PostgREST
CREATE VIEW bird_view (bird_name, place, sighted) AS 
  SELECT bird_name, place, sighted FROM birds.sighting ORDER BY sighted, bird_name;

-- record_bird is a stored procedure and will save a new bird sighting
CREATE PROCEDURE record_bird(name VARCHAR(256), description TEXT, dt DATE)
LANGUAGE SQL
AS $$
  INSERT INTO birds.sighting (bird_name, place, sighted) VALUES (name, description, dt);
$$;

--
-- The following additional steps are needed for PostgREST access.
--
DROP ROLE IF EXISTS birds_anonymous;
CREATE ROLE birds_anonymous nologin;

GRANT USAGE ON SCHEMA birds TO birds_anonymous;
GRANT SELECT ON birds.sighting TO birds_anonymous;

DROP ROLE IF EXISTS birds;
CREATE ROLE birds NOINHERIT LOGIN PASSWORD 'replace_me';
GRANT birds_anonymous TO birds;

