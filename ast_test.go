package newt

import (
	"path"
	"testing"
)

// Test the shared YAML configuration for Newt Router, Newt Generator
// Newt Mustache and Pandoc Bundler.
func TestLoadAST(t *testing.T) {
	configFiles := []string{
		path.Join("testdata", "birds.yaml"),
		path.Join("testdata", "blog.yaml"),
		path.Join("testdata", "bundler_test.yaml"),
	}
	for _, fName := range configFiles {
		ast, err := LoadAST(fName)
		if err != nil {
			t.Errorf("failed to load %q, %s", fName, err)
		}
		if ast == nil {
			t.Errorf("something went wrong, ast is nil for %q", fName)
		}
		if ast.Applications == nil {
			t.Errorf("ast.Applications is nil (%q), %+v", fName, ast)
		}
		ids := ast.GetModelIds()
		if len(ids) == 0 {
			t.Errorf("expected model ids for %q", fName)
		} else {
			mId := ids[0]
			model, ok := ast.GetModelById(mId)
			if !ok {
				t.Errorf("expected model for %q in %q, %s", mId, fName, err)
			}
			if model == nil {
				t.Errorf("expceted model content for %q in %q, got nil", mId, fName)
			}
		}
		names := ast.GetModelNames()
		if len(names) != len(ids) {
			t.Errorf("expected %d model names for %q, got %d", len(ids), fName, len(names))
		}
	}
}
