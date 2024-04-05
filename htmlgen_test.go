package newt

import (
	"testing"
)


// TestCheckboxesInputElement check that we can render a set of checkboxes with the
// possible behaviors indicated by the GHYT syntax for input elements of type "checkboxes"
func TestCheckboxesInputElement(t *testing.T) {
       t.Errorf("TestCheckboxesInputElement() not implemented")
}

// TestDropdownInputElement check that we can render a dropdown element with the
// possible behaviors indicated by the GHYT syntax for input elements of type "dropdown"
func TestDropdownInputElement(t *testing.T) {
       t.Errorf("TestDropdownInputElement() not implemented")
}

// TestMarkdownInputElement check that we can render a Markdown element with the
// possible behaviors indicated by the GHYT syntax for input elements of type "markdown"
func TestMarkdownInputElement(t *testing.T) {
       t.Errorf("TestMarkdownInputElement() not implemented")
}


// TestInputElement check that we can render a an input with the
// possible behaviors indicated by the GHYT syntax for input elements of type "input"
func TestInputElement(t *testing.T) {
	var (
		expected string
		got      string
	)
	input := new(Element)
	input.Type = "input"

	expected = "<input>"
	got = input.ToHTMLInput("", "")
	if expected != got {
		t.Errorf("expected %q, got %q", expected, got)
	}

	input.Id = "the_question"
	input.Attributes = map[string]string{
		"value": "hello world!",
		"title": "this is how we ask questions of the world!",
	}
	input.Validations = map[string]interface{}{
		"required": true,
	}

	expected = `<input id="the_question" required title="this is how we ask questions of the world!" value="hello world!">`
	got = input.ToHTMLInput("", "")
	if expected != got {
		t.Errorf("expected\n\t%s, got\n\t%s", expected, got)
	}

	// Test rendering as Mustache template
	expected = `<input id="the_question" required title="this is how we ask questions of the world!" value="{{the_question}}">`
	got = input.ToHTMLInput("{{", "}}")
	if expected != got {
		t.Errorf("expected\n\t%s, got\n\t%s", expected, got)
	}

	// Test rendering as Pandoc template with "$" delimiters
	expected = `<input id="the_question" required title="this is how we ask questions of the world!" value="$the_question$">`
	got = input.ToHTMLInput("$", "$")
	if expected != got {
		t.Errorf("expected\n\t%s, got\n\t%s", expected, got)
	}

	// Test rendering as Pandoc template with "${", "}" delimiters
	expected = `<input id="the_question" required title="this is how we ask questions of the world!" value="${the_question}">`
	got = input.ToHTMLInput("${", "}")
	if expected != got {
		t.Errorf("expected\n\t%s, got\n\t%s", expected, got)
	}

	// Test adding a label attribute
	input.Attributes["label"] = "Is this really the central question?"
	expected = `<label for="the_question">Is this really the central question?</label><input id="the_question" required title="this is how we ask questions of the world!" value="hello world!">`
	got = input.ToHTMLInput("", "")
	if expected != got {
		t.Errorf("expected\n\t%s\n, got\n\t%s\n", expected, got)
	}

}
