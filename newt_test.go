package newt

import (
	"path"
	"os"
	"testing"
)

func TestDryRun(t *testing.T){
	in := os.Stdin
	out := os.Stdout
	eout := os.Stderr
	args := []string{}
	args = append(args, path.Join("testdata", "blog.conf"))
	exitCode := Run(in, out, eout, args, true, true)
	if exitCode != 0 {
		t.Errorf("expected exit code zero, got %d\n", exitCode)
	}
}

func TestLiveRun(t *testing.T)  {
	t.Errorf("TestLiveRun() not implemented")
}
