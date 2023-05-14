package newt

import (
	"testing"
	"io"
)

func TestRun(t *testing.T) {
	var (
		in io.Reader
		out io.Writer
		eout io.Writer
		args []string
	)
	exitCode := Run(in, out, eout, args)
	if exitCode != 0 {
		t.Errorf("expected exit code zero, got %d\n", exitCode)
	}
}

