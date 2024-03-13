#
# This Makefile will generate the ornithology Newt application used to test the second prototype.
#

#PROJECT_FILES = ornithology_setup.sql ornithology_models.sql ornithology.conf \
#                create_sighting.tmpl read_sighting.tmpl update_sighting.tmpl delete_sighting.tmpl list_sighting.tmpl
PROJECT_FILES = ornithology_setup.sql ornithology_models.sql ornithology_models_test.sql ornithology.conf

build: $(PROJECT_FILES)

run: $(PROJECT_FILES)
	postgrest ornithology.conf &
	newt ornithology.yaml &


ornithology_setup.sql: ornithology.yaml .FORCE
	./bin/newtgenerator ornithology.yaml postgres setup >ornithology_setup.sql

ornithology_models.sql: ornithology.yaml .FORCE
	./bin/newtgenerator ornithology.yaml postgres models >ornithology_models.sql

ornithology_models_test.sql: ornithology.yaml .FORCE
	./bin/newtgenerator ornithology.yaml postgres models_test >ornithology_models_test.sql

ornithology.conf: ornithology.yaml
	./bin/newtgenerator ornithology.yaml postgrest >ornithology.conf

create_sighting.tmpl: ornithology.yaml
	./bin/newtgenerator ornithology.yaml mustache create sighting >create_sighting.tmpl

update_sighting.tmpl: ornithology.yaml
	./bin/newtgenerator ornithology.yaml mustache update sighting >update_sighting.tmpl

read_sighting.tmpl: ornithology.yaml
	./bin/newtgenerator ornithology.yaml mustache read sighting >read_sighting.tmpl

delete_sighting.tmpl: ornithology.yaml
	./bin/newtgenerator ornithology.yaml mustache delete sighting >delete_sighting.tmpl

list_sighting.tmpl: ornithology.yaml
	./bin/newtgenerator ornithology.yaml mustache list sighting >list_sighting.tmpl


.FORCE:
