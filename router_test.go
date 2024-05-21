package newt

import (
	"bufio"
	"bytes"
	"os"
	"path"
	"testing"
)

func TestNewRouter(t *testing.T) {
	ast := NewAST()
	router, err := NewRouter(ast)
	if err != nil {
		t.Error(err)
	}
	if router == nil {
		t.Errorf("NewRouter(ast) failed to create a new router")
	}
}

func TestRouter(t *testing.T) {
	ast := NewAST()
	//in := bufio.NewReader([]byte{})
	buf := bytes.NewBuffer([]byte{})
	out := bufio.NewWriter(buf)
	//eout := bufio.NewWriter([]byte{})
	router, err := NewRouter(ast)
	if err != nil {
		t.Error(err)
	}
	if router == nil {
		t.Errorf("NewRouter(ast) failed to create a new router")
	}
	if ok := router.Check(out); ok {
		t.Errorf("expected router.Check(out) to be false, nothing has been configured yet. %s", buf.Bytes())
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
	router, err = NewRouter(ast)
	if err != nil {
		t.Error(err)
	}
	if router == nil {
		t.Errorf("NewRouter(ast) failed to create a new router")
	}
	if ok := router.Check(out); !ok {
		t.Errorf("expected router.Check(out) to be true, %s", buf.Bytes())
	}
}
