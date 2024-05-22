package newt

import (
	"bufio"
	"bytes"
	"os"
	"path"
	"testing"
)

func TestNewTemplateEngine(t *testing.T) {
	ast := NewAST()
	fName := path.Join("testdata", "blog.yaml")
	src, err := os.ReadFile(fName)
	if err != nil {
		t.Errorf("failed to read %q, %s, aborting test", fName, err)
		t.FailNow()
	}
	if err := UnmarshalAST(src, ast); err != nil {
		t.Errorf("failed to unpack AST %q, %s, aborting test", fName, err)
		t.FailNow()
	}
	mustache, err := NewTemplateEngine(ast)
	if err != nil {
		t.Error(err)
	}
	if mustache == nil {
		t.Errorf("newt mustache should not be nil from AST %q", fName)
	}
}

func TestTemplateEngine(t *testing.T) {
	ast := NewAST()
	//in := bufio.NewReader([]byte{})
	buf := bytes.NewBuffer([]byte{})
	out := bufio.NewWriter(buf)
	//eout := bufio.NewWriter([]byte{})
	mustache, err := NewTemplateEngine(ast)
	if err != nil {
		t.Error(err)
	}
	if mustache == nil {
		t.Errorf("NewTemplateEngine(ast) failed to create a new mustache engine")
	}
	if ok := mustache.Check(out); ok {
		t.Errorf("expected mustache.Check(out) to be false, nothing has been configured yet. %s", buf.Bytes())
	}
	fName := path.Join("testdata", "blog.yaml")
	src, err := os.ReadFile(fName)
	if err != nil {
		t.Errorf("failed to read %q, %s, aborting test", fName, err)
		t.FailNow()
	}
	if err = UnmarshalAST(src, ast); err != nil {
		t.Errorf("failed to unpack AST %q, %s, aborting test", fName, err)
		t.FailNow()
	}
	mustache, err = NewTemplateEngine(ast)
	if err != nil {
		t.Error(err)
	}
	if mustache == nil {
		t.Errorf("NewTemplateEngine(ast) failed to create a new mustache engine")
	}
	if ok := mustache.Check(out); !ok {
		t.Errorf("expected mustache.Check(out) to be true, %s", buf.Bytes())
	}
}
