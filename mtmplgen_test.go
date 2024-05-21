package newt

import (
	"fmt"
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
    elements:
      - id: people_id
        type: text
        attributes:
          label: A unique person id, no spaces, alpha numeric
          placeholder: ex. jane-do-007
          required: true
        is_object_id: true
      - id: display_name
        type: text
        attributes:
          label: (optional) A person display name
          placeholder: ex. J. Doe, journalist
      - id: family_name
        type: text
        attributes:
          label: (required) A person's family name or singular when only one name exists
          placeholder: ex. Doe
          required: true
      - id: given_name
        type: text
        attributes:
          label: (optional, encouraged) A person's given name
          placeholder: ex. Jane
      - id: orcid
        type: text
        attributes:
          label: (optional) A person's ORCID identifier
          placeholder: ex. 0000-0000-0000-0000
          pattern: "[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]"
      - id: ror
        type: text
        attributes:
          label: (optional) A person's ROR identifying their affiliation
      - id: email
        type: email
        attributes:
          label: (optional) A person public email address
      - id: website
        type: url
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
	if ok := model.Check(out); !ok {
		t.Errorf("expected valid model, failed check\n%s\n", out.Bytes())
		t.FailNow()
	}
	// Test generation of display elements (used in read and list)
	for _, elemId := range model.GetElementIds() {
		if elem, ok := model.GetElementById(elemId); ok {
			if s := MDisplayElemGen(elem); s == "" {
				t.Errorf("expected value for %q, got %q", elem.Id, s)
			}
			if err := testMDisplayElemGen(elem); err != nil {
				t.Error(err)
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
			// Need to add checks to test output of MInputElemGen()
			if err := testMInputGen(elem); err != nil {
				t.Error(err)
			}
		} else {
			t.Errorf("failed to get element %q from model %q", elemId, modelId)
			t.FailNow()
		}
	}

	// Test generating partial template
	if err := MTmplGen(out, model, "search"); err == nil {
		t.Errorf("should have gotten an error for \"search\" as an unsupported action")
	}
	for _, action := range []string{"create_form", "create_response", "read", "update_form", "update_response", "delete_form", "delete_response", "list"} {
		out = bytes.NewBuffer([]byte{})
		if err := MTmplGen(out, model, action); err != nil {
			t.Error(err)
		}
		// FIXME: Need to add checks to test output of MInputGen()
		//t.Errorf("test of output to MTmplGen(out, model %q) not implemented", action)
	}
}

func testMDisplayElemGen(elem *Element) error {
	return fmt.Errorf("testMDisplayElemGen(%+v) not implemented", elem)
}

func testMInputGen(elem *Element) error {
	return fmt.Errorf("testMInputGen(%+v) not implemented", elem)
}

func testMTmplGen(model *Model, action string) error {
	return fmt.Errorf("testMTmplGen(%+v, %q) not implemented", model, action)
}
