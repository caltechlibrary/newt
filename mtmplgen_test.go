package newt

import (
	//"fmt"
	"bytes"
	"testing"
)

func TestMTmplGen(t *testing.T) {
	src := []byte(`applications:
  newtgenerator:
    namespace: people # E.g. "people" Namespace to use generating Postgres SQL
models:
  - id: people
    name: People Profiles
    description: |
      This models a curated set of profiles of colleagues
    body:
      - id: people_id
        type: input
        attributes:
          label: A unique person id, no spaces, alpha numeric
          placeholder: ex. jane-do-007
        validations:
          required: true
      - id: display_name
        type: input
        attributes:
          label: (optional) A person display name
          placeholder: ex. J. Doe, journalist
      - id: family_name
        type: input
        attributes:
          label: (required) A person's family name or singular when only one name exists
          placeholder: ex. Doe
        validations:
          required: true
      - id: given_name
        type: input
        attributes:
          label: (optional, encouraged) A person's given name
          placeholder: ex. Jane
      - id: orcid
        type: input
        attributes:
          label: (optional) A person's ORCID identifier
          placeholder: ex. 0000-0000-0000-0000
        validations:
          pattern: "[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]"
      - id: ror
        type: input
        attributes:
          label: (optional) A person's ROR identifying their affiliation
      - id: email
        type: "input[type=email]"
        attributes:
          label: (optional) A person public email address
      - id: website
        type: "input[type=url]"
        attributes:
          label: (optional) A person's public website
          placeholder: ex. https://jane.doe.example.org
`)
	ast := new(AST)
	if err := UnmarshalAST(src, ast); err != nil {
		t.Error(err)
		t.FailNow()
	}
	// out is our output buffer that'll be passed to MTmplGen function
    out := bytes.NewBuffer([]byte{})
	modelId := "people"
	model, ok := ast.GetModelById(modelId)
	if !ok {
		t.Errorf("failed to find model id %q", modelId)
		t.FailNow()
	}
	// Run Model check, should return true
	if ok := model.Check(out); ! ok {
		t.Errorf("expected valid model, failed check\n%s\n", out.Bytes())
		t.FailNow()
	}
	// Test generation of display elements (used in read and list)
	for _, elemId := range model.GetElementIds() {
		if elem, ok := model.GetElementById(elemId); ok {
			if s := MDisplayElemGen(elem); s == "" {
				t.Errorf("expected value for %q, got %q", elem.Id, s)
			}
		} else {
			t.Errorf("failed to get element %q from model %q", elemId, modelId)
			t.FailNow()
		}
	}

	// Test generation of input elements (used in create, update and delete)
	for _, elemId := range model.GetElementIds() {
		if elem, ok := model.GetElementById(elemId); ok {
			if s := MInputElemGen(elem); s == "" {
				t.Errorf("expected value for %q, got %q", elem.Id, s)
			}
			t.Errorf("test of output to MInputElemGen(out, %q, %q) not implemented", model.Id, elem.Id)
		} else {
			t.Errorf("failed to get element %q from model %q", elemId, modelId)
			t.FailNow()
		}
	}


	// Test generating partial template
	if err := MTmplGen(out, model, "search"); err == nil {
		t.Errorf("should have gotten an error for \"search\" as an unsupported action")
	}
	for _, action := range []string{"create", "read", "update", "delete", "list"} {
    	out = bytes.NewBuffer([]byte{})
		if err := MTmplGen(out, model, action); err != nil {
			t.Error(err)
		}
		t.Errorf("test of output to MTmplGen(out, model %q) not implemented", action)
	}
}
