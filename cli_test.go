/**
 * cli_test.go is a test runner for the cli.
 */
package newt

import (
	"os"
	"path"
	"testing"
)

func TestDryRunRouter(t *testing.T) {
	in := os.Stdin
	out := os.Stdout
	eout := os.Stderr
	args := []string{}
	args = append(args, path.Join("testdata", "blog.yaml"))
	exitCode := RunNewtRouter(in, out, eout, args, true)
	if exitCode != 0 {
		t.Errorf("expected exit code zero, got %d\n", exitCode)
	}
}

