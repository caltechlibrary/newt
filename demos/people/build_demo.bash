#!/bin/bash

mkdir -p htdocs
cat <<EOT >htdocs/index.html
<!DOCTYPE html>
<html lang="en-us">
<body>
Hello world!

<ul>
<li><a href="create_people">Create a Person record</a></li>
<li><a href="list_people">List people</a></li>
<li><a href="read_people/Doiel-R-S">Read a Person record</a></li>
<li><a href="update_people/Doiel-R-S">Update a Person record </a></li>
<li><a href="delete_people/Doiel-R-S">Delete a Person record</a></li>
</ul>
</body>
</html>
EOT
newt init people.yaml
vi people.yaml
newtgenerator people.yaml postgres setup >setup.sql
newtgenerator people.yaml postgres models >models.sql
newtgenerator people.yaml postgrest >postgrest.conf
dropdb --if-exists people
createdb people
psql -c "\\i setup.sql" people
psql -c "\\i models.sql" people
psql -c "\\dt" people
newtgenerator people.yaml mustache create_form people >create_people_form.tmpl
newtgenerator people.yaml mustache create_response people >create_people_response.tmpl
newtgenerator people.yaml mustache update_form people >update_people_form.tmpl
newtgenerator people.yaml mustache update_response people >update_people_response.tmpl
newtgenerator people.yaml mustache delete_form people >delete_people_form.tmpl
newtgenerator people.yaml mustache delete_response people >delete_people_response.tmpl
newtgenerator people.yaml mustache read people >read_people.tmpl
newtgenerator people.yaml mustache list people >list_people.tmpl
