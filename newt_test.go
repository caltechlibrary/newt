package newt

import (
	"io"
	"path"
	"testing"
)

func TestDryRun(t *testing.T) {
	var (
		in   io.Reader
		out  io.Writer
		eout io.Writer
		args []string
	)
	args = append(args, path.Join("testdata", "blog.conf"))
	exitCode := Run(in, out, eout, args, true)
	if exitCode != 0 {
		t.Errorf("expected exit code zero, got %d\n", exitCode)
	}
}

func TestLiveRun(t *testing.T)  {
	t.Errorf("TestLiveRun() not implemented")
}
