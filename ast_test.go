package newt

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"testing"
)

// TestNewtAST() tests NewAST() and NewApplications()
func TestNewAST(t *testing.T) {
	ast := NewAST()
	if ast.Applications == nil {
		t.Errorf("ast.Applications should not be nil")
	}
}

// TestUnmarshalAST tests unmarshalling YAML into a Newt AST object
func TestUnmarshalAST(t *testing.T) {
	configFiles := []string{
		path.Join("testdata", "birds.yaml"),
		path.Join("testdata", "blog.yaml"),
		path.Join("testdata", "bundler_test.yaml"),
	}
	for _, fName := range configFiles {
		src, err := os.ReadFile(fName)
		if err != nil {
			t.Errorf("failed to read %q, %s", fName, err)
		} else {
			ast := new(AST)
			if err := UnmarshalAST(src, ast); err != nil {
				t.Errorf("failed tn UnmarshalAST %q, %s", fName, err)
			} else {
				buf := bytes.NewBuffer([]byte{})
				if ok := ast.Check(buf); !ok {
					t.Errorf("UnmarshalAST %q, failed to pass check -> %s", fName, buf.Bytes())
				}
			}
		}
	}

}

// TestLoadAST tests reading on and populating the shared YAML configuration used
// by Newt applications.
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
	}
}

// TestHasChanges reads in our test YAML files, checks for changes, then modifies them and checkes for changes again.
func TestHasChanges(t *testing.T) {
	configFiles := []string{
		path.Join("testdata", "birds.yaml"),
		path.Join("testdata", "blog.yaml"),
		path.Join("testdata", "bundler_test.yaml"),
	}
	for i, fName := range configFiles {
		ast, err := LoadAST(fName)
		if err != nil {
			t.Errorf("failed to load %q, %s", fName, err)
		}
		if ast.HasChanges() {
			t.Errorf("should not have changes after LoadAST(%q)", fName)
		}
		whatChanged := ""
		switch i {
		case 0:
			ast.Applications = nil
			ast.isChanged = true
			whatChanged = fmt.Sprintf("removed .applications from ast for %q", fName)
		case 1:
			modelList := ast.GetModelIds()
			modelId := modelList[0]
			whatChanged = fmt.Sprintf("removed model %q from %q", modelId, fName)
			ast.RemoveModelById(modelId)
		case 2:
			t := ast.Templates[0]
			ast.Templates = ast.Templates[1:]
			ast.isChanged = true
			whatChanged = fmt.Sprintf("removed template %q -> %q for %q", t.Pattern, t.Template, fName)
		}
		if !ast.HasChanges() {
			t.Errorf("%s, did not detect change for %q", whatChanged, fName)
		}
	}
}
