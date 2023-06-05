package newt

import (
	"os"
	"path"
	"testing"
)

func TestDryRun(t *testing.T) {
	in := os.Stdin
	out := os.Stdout
	eout := os.Stderr
	args := []string{}
	args = append(args, path.Join("testdata", "blog.yaml"))
	exitCode := Run(in, out, eout, args, true)
	if exitCode != 0 {
		t.Errorf("expected exit code zero, got %d\n", exitCode)
	}
}

func TestLiveRun(t *testing.T) {
	t.Skipf("TestLiveRun() not implemented")
}
