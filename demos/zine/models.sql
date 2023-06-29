--
-- Model: zine.arts
-- Based on zine.yaml, 2023-06-23
--
CREATE TABLE zine.articles (
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
-- LIST VIEW: zine.articles 
-- FIXME: You probably want to customized this statement 
-- (e.g. add WHERE CLAUSE, ORDER BY, GROUP BY).
--
CREATE OR REPLACE VIEW zine.articles_list_view AS
    SELECT title, title_slug, by_line, content, pub_date, year, month, day
    FROM zine.articles;


