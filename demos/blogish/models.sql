--
-- Model: blog.posts
-- Based on blogish.yaml, 2023-06-23
--
CREATE TABLE blog.posts (
    year DATE,
    month DATE,
    day DATE,
    title VARCHAR(256) DEFAULT '',
    title_slug VARCHAR(256) DEFAULT '' PRIMARY KEY,
    by_line VARCHAR(256) DEFAULT '',
    content TEXT DEFAULT '',
    pub_date DATE
);

--
-- LIST VIEW: blog.posts 
-- FIXME: You probably want to customized this statement 
-- (e.g. add WHERE CLAUSE, ORDER BY, GROUP BY).
--
CREATE OR REPLACE VIEW blog.posts_list_view AS
    SELECT title, title_slug, by_line, content, pub_date, year, month, day
    FROM blog.posts;


