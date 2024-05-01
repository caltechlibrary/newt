package newt

import (
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
	// Id is a required field for Newt, it is an extension and isn't part of GHYITS.
	Id string `json:"id,required" yaml:"id,omitempty"`

	// This is a Newt specifc set of attributes to place in the form element of HTML. I.e. it could
	// be form "class", "method", "action", "encoding". It is not defined in the GitHub YAML issue template syntax
	// (optional)
	Attributes map[string]interface{} `json:"attributes,omitempty" yaml:"attributes,omitempty"`

	// Name, A name for the issue form template. Must be unique from all other templates, including Markdown templates.
	//
	// For Newt this should conform to the variable naming conventions, starts with an alphabetical character, may be
	// alpha number without spaces or punctuation other than '_'.
	// (required)
	Name string `json:"name,required" yaml:"name,omitempty"`

	// Description, A description for the issue form template, which appears in the template chooser interface.
	// (required)
	Description string `json:"description,required" yaml:"description,omitempty"`

	// Body, Definition of the input types in the form.
	// (required)
	Body []*Element `json:"body,required" yaml:"body,omitempty"`

	// Title, A default title that will be pre-populated in the issue submission form.
	// (optional)
	Title string `json:"title,omitempty" yaml:"title,omitempty"`

	//
	// The following are included in the struct for compatibilty with GitHub YAML issue template syntax.
	// They are ignored by Newt as they are not relevant to the general case of rendering an web application.
	//

	// Assignees, People who will be automatically assigned to issues created with this template.
	// (ignored)
	Assignees []string `json:"assignees,omitempty" yaml:"assignees,omitempty"`

	// Labels, Labels that will automatically be added to issues created with this template. If a label does not
	// already exist in the repository, it will not be automatically added to the issue.
	// (ignored)
	Labels []string `json:"labels,omitempty" yaml:"labels,omitempty"`

	// Projects, Projects that any issues created with this template will automatically be added to. The format
	// of this key is PROJECT-OWNER/PROJECT-NUMBER.
	//
	// Note: The person opening the issue must have write permissions for the specified projects. If you don't
	// expect people using this template to have write access, consider enabling your project's auto-add workflow.
	// For more information, see "Adding items automatically."
	// (ignored)
	Projects []string `json:"projects,omitempty" yaml:"projects,omitempty"`
}

// GetElementIds returns a slice of element ids found in the model's .Body
func (m *NewtModel) GetElementIds() []string {
	ids := []string{}
	for _, elem := range m.Body {
		if elem.Id != "" {
			ids = append(ids, elem.Id)
		}
	}
	return ids
}

// GetElementById returns a Element from the model's .Body.
func (m *NewtModel) GetElementById(id string) (*Element, bool) {
	for _, elem := range m.Body {
		if elem.Id == id {
			return elem, true
		}
	}
	return nil, false
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
	Type string `json:"type,required" yaml:"type,omitempty"`

	// Id for the element, except when type is set to markdown. Can only use alpha-numeric characters,
	//  -, and _. Must be unique in the form definition. If provided, the id is the canonical identifier
	//  for the field in URL query parameter prefills.
	Id string `json:"id,omitempty" yaml:"id,omitempty"`

	// Attributes, a set of key-value pairs that define the properties of the element.
	// This is a required element as it holds the "value" attribute when expressing
	// HTML content. Other commonly use attributes
	Attributes map[string]string `json:"attributes,omitempty" yaml:"attributes,omitempty"`

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
	Validations map[string]interface{} `json:"validations,omitempty" yaml:"validations,omitempty"`
}

