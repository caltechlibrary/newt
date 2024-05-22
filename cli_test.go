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
	exitCode := RunRouter(in, out, eout, args, true, 8011, path.Join("testdata", "htdocs"), false)
	if exitCode != 0 {
		t.Errorf("expected exit code zero, got %d\n", exitCode)
	}
}

func TestRunNewtGenerator(t *testing.T) {
	t.Errorf("TestRunNewtGenerator() not implemented")
}

func TestRunTemplateEngine(t *testing.T) {
	t.Errorf("TestRunTemplateEngine() not implemented")
}

func TestRunRouter(t *testing.T) {
	t.Errorf("TestRunTemplateEngine() not implemented")
}

func TestRunStaticWebServer(t *testing.T) {
	t.Errorf("TestRunStaticWebServer() not implemented")
}

func TestRunNewtCheckYAML(t *testing.T) {
	t.Errorf("TestRunCheckYAML() not implemented")
}

func TestRunNewt(t *testing.T) {
	t.Errorf("TestRunNewt() not implemented")
}

func TestRunNewtApplications(t *testing.T) {
	t.Errorf("TestRunNewthApplications() not implemented")
}

func TestRunMustacheCLI(t *testing.T) {
	t.Errorf("TestRunMustacheCLI() not implemented")
}

func TestRunNewtConfig(t *testing.T) {
	t.Errorf("TestRunNewtConfig() not implemented")
}

func TestRunModeler(t *testing.T) {
	t.Errorf("TestRunModeler() not implemented")
}
