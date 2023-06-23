--
-- Below is the SQL I would noramally check into a project repository.
-- It does not contain secrets. It contains our data models, views,
-- and functions. This defines the behaviors made available through
-- PostgRESTS.
--

-- Make sure we are in the birds database
\c birds
SET search_path TO birds, public;

--
-- Data Models
--

--
-- This file contains the create statements for our bird table.
--
-- DROP TABLE IF EXISTS birds.sighting;
CREATE TABLE birds.sighting
(
  bird VARCHAR(255),
  place TEXT,
  sighted DATE
);

--
-- Data Views and Behaviors.
--

-- bird_view will become an end point in PostgREST, '/bird_view'
CREATE OR REPLACE VIEW birds.bird_view AS
  SELECT bird, place, sighted
  FROM birds.sighting ORDER BY sighted ASC, bird ASC;

-- record_bird is a stored procedure and will save a new bird sighting.
-- It becomes the end point '/rpc/record_bird'
CREATE OR REPLACE FUNCTION birds.record_bird(bird VARCHAR, place TEXT, sighted DATE)
RETURNS bool LANGUAGE SQL AS $$
  INSERT INTO birds.sighting (bird, place, sighted) VALUES (bird, place, sighted);
  SELECT true;
$$;

--
-- PostgREST access and controls.
--

-- Since our Postgres ROLE and SCHEMA exist and our models may change how
-- we want PostgREST to expose our data via JSON API we GRANT or 
-- revoke role permissions here.
-- with our model.
GRANT USAGE  ON SCHEMA birds      TO birds_anonymous;
-- NOTE: We are allowing insert because this is a demo and we are not
-- implementing an account login requirement. Normally I would force
-- a form update via a SQL function or procedure only.
GRANT SELECT, INSERT ON birds.sighting    TO birds_anonymous;
GRANT SELECT ON birds.bird_view   TO birds_anonymous;
GRANT EXECUTE ON FUNCTION birds.record_bird TO birds_anonymous;


