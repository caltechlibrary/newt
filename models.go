package newt

import (
	"fmt"
	"sort"
	"strings"

	// 3rd Party Packages
	//"gopkg.in/yaml.v3"
)

// NewtModel implements a structure that can accomodate the GitHub YAML issue template syntax.
// See <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-issue-forms>
//
// The NewtModel structure is used by Newt to describe data models. It will be used in code generation and in validating
// POST and PUT requests to the data router.
// 
// Code generation will need to render to SQL, HTML, Mustache and Pandoc templates.
type NewtModel struct {
	// This is a Newt specifc set of attributes to place in the form element of HTML. I.e. it could
	// be form "class", "method", "action", "encoding". It is not defined in the GitHub YAML issue template syntax
	// (optional)
	Attributes map[string]interface{} `json:"attributes,omitempty" yaml:"attributes"`

	// Name, A name for the issue form template. Must be unique from all other templates, including Markdown templates.
	//
	// For Newt this should conform to the variable naming conventions, starts with an alphabetical character, may be
	// alpha number without spaces or punctuation other than '_'.
	// (required)
	Name string `json:"name,required" yaml:"name"`

	// Description, A description for the issue form template, which appears in the template chooser interface.
	// (required)
	Description string `json:"description,required" yaml:"description"`

	// Body, Definition of the input types in the form.
	// (required)
	Body []*Element `json:"body,required" yaml:"body"`

	// Title, A default title that will be pre-populated in the issue submission form.
	// (optional)
	Title string `json:"title,omitempty" yaml:"title"`

	//
	// The following are included in the struct for compatibilty with GitHub YAML issue template syntax.
	// They are ignored by Newt as they are not relevant to the general case of rendering an web application.
	//

	// Assignees, People who will be automatically assigned to issues created with this template.
	// (ignored)
	Assignees []string `json:"assignees,omitempty" yaml:"assignees"`

	// Labels, Labels that will automatically be added to issues created with this template. If a label does not
	// already exist in the repository, it will not be automatically added to the issue.
	// (ignored)
	Labels []string `json:"labels,omitempty" yaml:"labels"`

	// Projects, Projects that any issues created with this template will automatically be added to. The format
	// of this key is PROJECT-OWNER/PROJECT-NUMBER.
	//
	// Note: The person opening the issue must have write permissions for the specified projects. If you don't
	// expect people using this template to have write access, consider enabling your project's auto-add workflow.
	// For more information, see "Adding items automatically." 
	// (ignored)
	Projects []string `json:"projects,omitempty" yaml:"projects"`
}

// ToHTMLForm, render as an HTML web form
func (f *NewtModel) ToHTMLForm() string {
	return "ToHTMLForm() not implemented"
}

// ToHTML, render as an HTML display block
func (f *NewtModel) ToHTML() string {
	return "ToHTML() not implemented"
}

// ToHTMLFormTempalte, render an an HTML web form template. Note this can produce
// wither either Mustache or Pandoc template formats depending on the prefix and
// suffix for the variable reference supplied. (e.g. "${", "}" or "$", "$" would 
// the variable values used in Pandoc templates while "{{", "}}" would conform to
// Mustache)
func (f *NewtModel) ToHTMLFormTemplate(prefix string, suffix string) string {
	return "ToHTMLFormTemplate(prefix, suffix) not implemented"
}

// Element implementes the GitHub YAML issue template syntax for an input element.
// The input element YAML is described at <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-githubs-form-schema>
//
// While the syntax most closely express how to setup an HTML representation it is equally
// suitable to expressing, through inference, SQL column type definitions. E.g. a bare `input` type is a `varchar`,
// a `textarea` is a `text` column type, an `input[type=date]` is a date column type.
type Element struct {
	// Type, The type of element that you want to input. It is required. Valid values are
	// checkboxes, dropdown, input, markdown and text area.
	//
	// The "input" type can be very expressive using CSS selector syntax. E.g. `input[type=URL]` would
	// create an input for URLs and validate using HTML basic intrinsic containts. You can also express
	// many other HTML form elements using the input notation. E.g. A button can be expressed as
	// `input[type=button]` where the attribute `value` becomes the inner HTML of `<button></button>`.
	// See MDN developer docs for input, <https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input>
	Type string `json:"type,required" yaml:"type"`

	// Id for the element, except when type is set to markdown. Can only use alpha-numeric characters,
	//  -, and _. Must be unique in the form definition. If provided, the id is the canonical identifier
	//  for the field in URL query parameter prefills.
	Id string `json:"id,omitempty" yaml:"id"`

	// Attributes, a set of key-value pairs that define the properties of the element.
	// This is a required element as it holds the "value" attribute when expressing
	// HTML content. Other commonly use attributes
	Attributes map[string]interface{} `json:"attributes,omitempty" yaml:"attributes"`

	// Validations, A set of key-value pairs that set constraints on the element.
	// Optional, key-value pair example expressed in HTML include
	// `required="true"` and "pattern" followed by a JavaScript RegExp would be
	// another example of validations.
	//
	// See MDN documentation for Intrinsic and basic contraints on input types,
	// <https://developer.mozilla.org/en-US/docs/Web/HTML/Constraint_validation>
	//
	// See MDN documentation for pattern attribute,
	// <https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes/pattern>
	//
	// See MDN documentation for RegExp,
	// <https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Regular_Expressions/Cheatsheet>
	Validations map[string]interface{} `json:"validations,omitempty" yaml:"validations"`
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

// ToHTMLInput renders the HTML markup for an input element. 
// then it'll replace the input's value attribute with markup for
// a template langauge(e.g. Pandoc or Mustache).
//
// Example rendering of rendering HTML
//
// ```
// e := new(Element)
// e.Id = "project_url"
// e.Type = "url"
// e.Attributes = map[string]string {
//     "value": "https://example.org"
// }
// src := e.Element.ToHTMLInput("", "")
// ```
//
// This would ouput something like
//
// ```
// <input type="url" value="https://example.org">
// ```
//
// Now rendering as a pandoc template
//
// ```
// src = e.Element.ToHTMLInput("${", "}")
// ```
//
// This would ouput something like
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
		if (k != "label") && (k != "description")  {
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

