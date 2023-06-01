DROP DATABASE IF EXISTS test_blog;
CREATE DATABASE test_blog;

\c test_blog

DROP SCHEMA IF EXISTS test_blog CASCADE;
CREATE SCHEMA test_blog;
SET search_path TO test_blog, public;


--
-- This file contains the create statements for our bird table.
--
DROP TABLE IF EXISTS test_blog.posts;
CREATE TABLE test_blog.posts
(
  post_id INT,
  title VARCHAR(255),
  post TEXT,
  created DATE,
  updated DATE
);

--
-- The following additional steps are needed for PostgREST access.
--
DROP ROLE IF EXISTS test_blog_anonymous;
CREATE ROLE test_blog_anonymous nologin;

GRANT USAGE  ON SCHEMA test_blog      TO test_blog_anonymous;
-- NOTE: We're allowing insert because this is a demo and we're not
-- implementing a login requirement!!!!
GRANT SELECT, INSERT ON test_blog.posts    TO test_blog_anonymous;

DROP ROLE IF EXISTS test_blog;
CREATE ROLE test_blog NOINHERIT LOGIN PASSWORD 'my_secret_password';
GRANT test_blog_anonymous TO test_blog;



