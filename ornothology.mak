#
# This Makefile will generate the ornothology Newt application used to test the second prototype.
#

#PROJECT_FILES = ornothology_setup.sql ornothology_models.sql ornothology.conf \
#                create_sighting.tmpl read_sighting.tmpl update_sighting.tmpl delete_sighting.tmpl list_sighting.tmpl
PROJECT_FILES = ornothology_setup.sql ornothology_models.sql ornothology_models_test.sql ornothology.conf

build: $(PROJECT_FILES)

run: $(PROJECT_FILES)
	postgrest ornothology.conf &
	newt ornothology.yaml &


ornothology_setup.sql: ornothology.yaml .FORCE
	./bin/newtgenerator ornothology.yaml postgres setup >ornothology_setup.sql

ornothology_models.sql: ornothology.yaml .FORCE
	./bin/newtgenerator ornothology.yaml postgres models >ornothology_models.sql

ornothology_models_test.sql: ornothology.yaml .FORCE
	./bin/newtgenerator ornothology.yaml postgres models_test >ornothology_models_test.sql

ornothology.conf: ornothology.yaml
	./bin/newtgenerator ornothology.yaml postgrest >ornothology.conf

create_sighting.tmpl: ornothology.yaml
	./bin/newtgenerator ornothology.yaml mustache create sighting >create_sighting.tmpl

update_sighting.tmpl: ornothology.yaml
	./bin/newtgenerator ornothology.yaml mustache update sighting >update_sighting.tmpl

read_sighting.tmpl: ornothology.yaml
	./bin/newtgenerator ornothology.yaml mustache read sighting >read_sighting.tmpl

delete_sighting.tmpl: ornothology.yaml
	./bin/newtgenerator ornothology.yaml mustache delete sighting >delete_sighting.tmpl

list_sighting.tmpl: ornothology.yaml
	./bin/newtgenerator ornothology.yaml mustache list sighting >list_sighting.tmpl


.FORCE:
