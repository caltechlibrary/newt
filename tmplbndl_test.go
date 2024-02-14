package newt

import (
	"os"
	"testing"
	"fmt"
)

func TestTemplateBundler(t *testing.T) {
	src := `port: 3029
templates:
  - request: "POST /custom_page"
    template: page.tmpl
    options:
      from: markdown
      to: html5
      standalone: true
      title: This is the custom template with this title
  - request: "POST /custom_page_with_title/{title}"
    template: page.tmpl
    options:
      from: markdown
      to: html5
      standalone: true
      title: This title is overwritten by the one in the request
  - request: "POST /custom_page_include"
    template: page.tmpl
    options:
      from: markdown
      to: html5
      standalone: false
  - request: "POST /default_html5"
    options:
      from: markdown
      to: html5
      standalone: true
      title: A Page using the default template
  - request: "POST /default_html5_with_title/{title}"
    options:
      from: markdown
      to: html5
      standalone: true
      title: This title is replaced by the title in the URL
  - request: "POST /default_html5_include"
    options:
      from: markdown
      to: html5
      standalone: false
`
	// Write our a YAML so we can test loading it.
	yamlName := "testdata/bundler_test.yaml"
	fp, err := os.Create(yamlName)
	if err != nil {
		t.Errorf("failed to create %q test data file, %s", yamlName, err)
		t.FailNow()
	}
	fmt.Fprintf(fp, "%s", src)
	fp.Close()

	tb, err := NewTemplateBundler(yamlName)
	if err != nil {
		t.Errorf("%s", err)
		t.FailNow()
	}
	expectedPort := ":3029"
	if expectedPort != tb.Port {
		t.Errorf("expected %q, got %q", expectedPort, tb.Port)
		t.FailNow()
	}

	tb.Port = ":8029"
	if tb == nil {
		t.Errorf("tb template bundler should not be nil")
		t.FailNow()
	}
}
