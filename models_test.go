package newt

import (
	"path"
	"testing"
)

func TestModelFuncs(t *testing.T) {
	fName := path.Join("testdata", "birds.yaml")
	ast, err := LoadAST(fName)
	if err != nil {
		t.Errorf("failed to load %q, aborting test, %s", fName, err)
		t.FailNow()
	}
	mName := "bird_sighting"
	m, ok := ast.GetModelById(mName)
	if ! ok {
		t.Errorf("failed to find %q in %q, aborting test", mName, fName)
		t.FailNow()
	}
	ids := m.GetElementIds()
	if len(ids) != 3 {
		t.Errorf("model %q should have three element ids, got %+v", mName, ids)
		t.FailNow()
	}
}
