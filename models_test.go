package newt

import (
	"testing"
)

func TestModelsPgModelSQL(t *testing.T) {
   src := []byte(`#
# cold.yaml defines the controlled object list and datum application
# from the perspective of Newt, a data route and form validator that
# relies on Postgres+PostgREST and Pandoc server for completing the microservice
# trio that is cold.
# 
# @author: R. S. Doiel <rsdoiel@caltech.edu>
#
namespace: cold
models:
- model: cl_person
  var:
    cl_people_id: String*
    family_name: String
    given_name: String
    orcid: ORCID
    author_id: String
    directory_id: String
    caltech: Boolean
    status: Boolean
    direcotry_person_type: String
    division: String
    updated: Date 2006-01-02
    ror: ROR
- model: cl_group
  var:
    cl_group_id: String*
    name: String
# FIXME: This should be an array of string but currently in CSV is a semi-colon delimited list
    alternative: String
    updated: Timestamp
# Date group was added to Group List
    date: Date 2006-01-02
    email: EMail
    description: Text
    start: String
# approx_starts indicates if the "start" is an approximate (true) or exact (false)
    approx_start: Boolean
# activity is a string value holding either "active" or "inactive"  
    activity: String
    end: String
# approx_end indicates if the "start" is an approximate (true) or exact (false)
    approx_end: Boolean
    website: URL
    pi: String
    parent: String
# prefix holds the DOI prefix associated with the group
    prefix: String
    grid: String
    isni: ISNI
    ringold: String
    viaf: String
    ror: ROR
#
# Controlled vocabulary models, these are mostly key/value pairs
#
- model: cl_subject
  var:
    cl_subject_id: String*
    subject: String
- model: cl_publishers
  var:
    issn: ISSN*
    publisher: String
- model: cl_doi_prefixes
  var:
    prefix: String*
    name: String
`)
   cfg := new(Config)
   if err := Parse(src, cfg); err != nil {
	   t.Error(err)
	   t.FailNow()
   }
   for i, model := range cfg.Models {
	   if _, err := PgModelSQL("test_yaml_source", model); err != nil {
		   t.Errorf("model (%d) %s: %s", i, model.Name, err)
	   }
   }

}

func TestModelsParse(t *testing.T) {
	cfg := new(Config)
	src := []byte(`
#
# cold.yaml defines the controlled object list and datum application
# from the perspective of Newt, a data route and form validator that
# relies on Postgres+PostgREST and Pandoc server for completing the microservice
# trio that is cold.
# 
# @author: R. S. Doiel <rsdoiel@caltech.edu>
#
namespace: cold
models:
- model: cl_person
- var:
  cl_people_id: String*
  family_name: String
  given_name: String
  orcid: ORCID
  author_id: String
  directory_id: String
  caltech: Boolean
  status: Boolean
  direcotry_person_type: String
  division: String
  updated: Date 2006-01-02
  ror: ROR
- model: cl_group
- var:
  cl_group_id: String*
  name: String
  # FIXME: This should be an array of string but currently in CSV is a semi-colon delimited list
  alternative: String
  updated: Timestamp
  # Date group was added to Group List
  date: Date 2006-01-02
  email: EMail
  description: Text
  start: String
  # approx_starts indicates if the "start" is an approximate (true) or exact (false)
  approx_start: Boolean
  # activity is a string value holding either "active" or "inactive"  
  activity: String
  end: String
  # approx_end indicates if the "start" is an approximate (true) or exact (false)
  approx_end: Boolean
  website: Url
  pi: String
  parent: String
  # prefix holds the DOI prefix associated with the group
  prefix: String
  grid: String
  isni: ISNI
  ringold: String
  viaf: String
  ror: ROR
#
# Controlled vocabulary models, these are mostly key/value pairs
- model: cl_subject
  cl_subject_id: String*
  subject: String
- model: cl_publishers
  issn: ISSN*
  publisher: String
- model: cl_doi_prefixes
  prefix: String*
  name: String
#
# cold service web configuration 
#
port: 8000
htdocs: htdocs
#
# cold routes, static content routes are noted only in comments. Dynamic
# ones are notated for orchestration by Newt
routes:
# Root path provides a description of PostgREST information (static asset path)
#- request_path: /
# Returns the version number of the service (static asset path
#- request_path: /cold/version
#
# /cold/people returns a list of all cl people ID
# 
- request_path: /cold/people
  request_method: GET
  api_method: GET
  # cl_people_ids is implemented as a SQL view of cl_person model
  api_url: http://localhost:3000/cl_people_ids
# 
# cold people {cl_people_id} returns a single person record
#
- var:
    cl_people_id: String
  request_path: /cold/people/${cl_people_id}
  request_method: GET
  api_method: GET
  # person_object is implemented as a SQL function of cl_person model
  api_url: http://localhost:3000/rpc/person_object?cl_people_id=${cl_people_id}
#
# /cold/group returns a list of all cl_group_ids
#
- request_path: /cold/group
  request_method: GET
  api_method: GET
  api_url: http://localhost:3000/cl_group_ids
#
# /cold/group/${cl_group_id}
# For a GET returns a group object, a PUT will create the group object, POST will replace the group object and DELETE will remove the group object
- var:
    cl_group_id: String
  request_path: /cold/group/${cl_group_id}
  request_method: GET
  api_method: GET
  # group_object is imlpemented as an SQL function of the cl_group model
  api_url: http://localhost:3000/rpc/group_object?cl_group_id=${cl_group_id}
#
# Crosswalks
#
#
# A cross walk lets you put in a collection name (e.g. people, group), 
# a field name and a value and it returns a list of matching records.
#
# /cold/crosswalk/people/${identifier_name}/${identifier_value}
# Returns a list of "cl_people_id" assocated with that identifier
- var: 
    identifier_name: String
    identifier_value: String
  request_path: /cold/crosswalk/people/${identifier_name}/${identifier_value}
  request_method: GET
  api_content_method: GET
  # people_crosswalk is implemented as an SQL function operating on the cl_person model
  api_url: http://localhost:3000/rpc/people_crosswalk?identifier=${identifier_name}&value=${identifier_value}
#
# /cold/crosswalk/group/${identifier_name}/${identifier_value}
# Returns a list of "cl_group__id" assocated with that identifier
- var:
    identifier_name: String
    identifier_value: String
  request_path: /cold/crosswalk/group/${identifier_name}/${identifier_value}
  api_content_method: GET
  api_url: http://localhost:3000/rpc/group_crosswalk?identifier=${identifier_name}&value=${identifier_value}
#
# Vocabularies
# ------------
#
# **cold** also supports end points for stable vocabularies mapping an 
# indentifier to a normalized name. These are set at compile time because
# they are so slow changing. 
#
# /cold/subject
# : Returns a list of all the subject ids (codes)
- request_path: /cold/subject
  request_method: GET
  api_method: GET
  # subjects is implemented as an SQL view returning all subject ids
  api_url: http://localhost:3000/subjects
# 
# /cold/subject/${subject_id}
# : Returns the normalized text string for that subject id
- var:
    subject_id: String
  request_path: /cold/subject/${subject_id}
  request_method: GET
  api_method: GET
  # get_subject is implemented as an SQL function on cl_subject model
  api_url: http://localhost:3000/rpc/get_subject?subject_id=${subject_id}
#
# /cold/issn
# : Returns a list of issn that are mapped to a publisher name
- request_path: /cold/issn
  request_method: GET
  api_method: GET
  api_url: http://localhost:3000/publishers
#
# /cold/issn/${issn}
#: Returns the normalized publisher name for that ISSN
- var:
    issn: ISSN
  request_path: /cold/issn/${issn}
  request_method: GET
  api_method: GET
  # get_publisher is implemented as a SQL function on cl_publisher model
  api_url: http://localhost:3000/rpc/get_publisher?issn=${issn}
#
# /cold/doi-prefix
# : Returns a list of DOI prefixes that map to a normalize name
- request_path: /cold/doi-prefix
  request_method: GET
  api_method: GET
  # doi-prefixes is implemented as a SQL view on cl_doi_prefix model
  api_url: http://localhost:3000/doi-prefixes
#
# /cold/doi-prefix/${doi_prefix}
# : Returns the normalized publisher name for that DOI prefix
- var:
    doi_prefix: String
  request_path: /cold/doi-prefix/${doi_prefix}
  request_method: GET
  api_method: GET
  # get_doi_prefix is implemented as a SQL function on cl_doi_prefix model
  api_url: http://localhost:3000/rpc/get_doi_prefix?doi_prefix=${doi_prefix}
`)
	if err := Parse(src, cfg); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
