--
-- This script is a convienence. It will use the psql client
-- copy command to load the tables with some test data.
--

-- Make sure we are in the birds database
\c birds


-- Now import our CSV file of birds.csv
\copy birds.sighting from 'birds.csv' with (FORMAT CSV, HEADER);

-- Make sure the data loaded, query with a view statement.
SELECT * FROM birds.bird_view;

