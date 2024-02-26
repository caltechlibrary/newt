package newt

import (
	"fmt"
	"sort"
	"strings"
)

// ToHTMLInput renders the HTML markup for an input element.
// then it'll replace the input's value attribute with markup for
// a template langauge(e.g. Pandoc or Mustache).
//
// # Example rendering of rendering HTML
//
// ```
// e := new(Element)
// e.Id = "project_url"
// e.Type = "url"
//
//	e.Attributes = map[string]string {
//	    "value": "https://example.org"
//	}
//
// src := e.Element.ToHTMLInput("", "")
// ```
//
// # This would ouput something like
//
// ```
// <input type="url" value="https://example.org">
// ```
//
// # Now rendering as a pandoc template
//
// ```
// src = e.Element.ToHTMLInput("${", "}")
// ```
//
// # This would ouput something like
//
// ```
// <input type="url" value="${project_url}">
// ```
func (e *Element) ToHTMLInput(prefix string, suffix string) string {
	renderAsTemplate := !(prefix == "" && suffix == "")
	// Are we rendering standard HTML or a template?
	inputType := strings.ToLower(e.Type)
	switch inputType {
	case "checkboxes":
		return e.toHTMLCheckboxesInput()
	case "dropdown":
		return e.toHTMLDropdownInput()
	case "markdown":
		return e.toHTMLMarkdownInput()
	}
	inputType = strings.TrimSpace(strings.TrimPrefix(inputType, "input"))
	if strings.HasPrefix(inputType, "[") {
		inputType = strings.TrimPrefix(
			strings.TrimSpace(
				strings.TrimSuffix(strings.TrimPrefix(
					strings.ReplaceAll(inputType, " ", ""),
					"["), "]")), "type=")
	}
	// Do we need to generate a label element?
	block := []string{}
	label, ok := e.Attributes["label"]
	if ok {
		block = append(block, fmt.Sprintf("<label for=%q>%s</label>", e.Id, label))
	}
	// Assemble our HTML input element
	attr := []string{}
	if inputType != "" {
		attr = append(attr, fmt.Sprintf("type=%q", inputType))
	}
	if e.Id != "" {
		attr = append(attr, fmt.Sprintf("id=%q", e.Id))
	}
	for k, v := range e.Validations {
		switch k {
		case "required":
			if v.(bool) {
				attr = append(attr, "required")
			}
		case "pattern":
			if v.(string) != "" {
				attr = append(attr, fmt.Sprintf("pattern=%q", v.(string)))
			}
		}
	}
	// I want to enforce an order of the keys to make testing easier.
	keys := []string{}
	for k, _ := range e.Attributes {
		if (k != "label") && (k != "description") {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)
	for _, k := range keys {
		val := e.Attributes[k].(string)
		// Possible special handling
		if renderAsTemplate && k == "value" {
			// Substitue "value" value with variable reference in a the template language
			val = fmt.Sprintf("%s%s%s", prefix, e.Id, suffix)
		}
		attr = append(attr, fmt.Sprintf("%s=%q", k, val))
	}
	if len(attr) > 0 {
		block = append(block, fmt.Sprintf(`<input %s>`, strings.Join(attr, " ")))
	} else {
		block = append(block, "<input>")
	}
	return strings.Join(block, "")
}

// ToHTMLForm, render as an HTML web form
func (f *NewtModel) ToHTMLForm() string {
       return "ToHTMLForm() not implemented"
}

// ToHTML, render as an HTML display block
func (f *NewtModel) ToHTML() string {
       return "ToHTML() not implemented"
}

func (e *Element) toHTMLCheckboxesInput() string {
       return "toHTMLCheckboxesInput() not implement"
}

func (e *Element) toHTMLDropdownInput() string {
       return "toHTMLDropdownInput() not implemented"
}

func (e *Element) toHTMLMarkdownInput() string {
       return "toHTMLMarkdownInput() not implemented"
}
