package newt

import (
	"fmt"
	"io"
	"regexp"
)

// Model implements a structure that can accomodate the GitHub YAML issue template syntax.
// See <https://docs.github.com/en/communities/using-templates-to-encourage-useful-issues-and-pull-requests/syntax-for-issue-forms>
//
// The Model structure is used by Newt to describe data models. It will be used in code generation and in validating
// POST and PUT requests to the data router.
//
// Code generation will need to render to SQL, HTML, Mustache and Pandoc templates.
type Model struct {
	// Id is a required field for Newt, it is an extension and isn't part of GHYITS.
	Id string `json:"id,required" yaml:"id"`

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

	// isChanged is an internal state used by the modeler to know when a model has changed
	isChanged bool `json:"-" yaml:"-"`
}

// HasChanges checks if the model's elements have changed
func (model *Model) HasChanges() bool {
	if model.isChanged {
		return true
	}
	for _, e := range model.Body {
		if e.isChanged {
			return true
		}
	}
	return false
}
// HasElement checks if the model has a given element id
func (model *Model) HasElement(elementId string) bool {
	for _, e := range model.Body {
		if e.Id == elementId {
			return true 
		}
	}
	return false
}

// GetModelIdentifier() returns the element which describes the model identifier.
// Returns the element and a boolean set to true if found.
func (m *Model) GetModelIdentifier() (*Element, bool) {
	for _, e := range m.Body {
		if e.IsModelIdentifier {
			return e, true
		}
	}
	return nil, false
}

// GetElementIds returns a slice of element ids found in the model's .Body
func (m *Model) GetElementIds() []string {
	ids := []string{}
	for _, elem := range m.Body {
		if elem.Id != "" {
			ids = append(ids, elem.Id)
		}
	}
	return ids
}

// GetElementById returns a Element from the model's .Body.
func (m *Model) GetElementById(id string) (*Element, bool) {
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

	// IsModelIdentifier is used by the modeler and config object to know which id values to embed
	// in routes and templates as the identifier for an object.
	IsModelIdentifier bool `json:"is_model_identifier,omitempty" yaml:"is_model_identifier,omitempty"`

	//
	// These fields are used by the modeler to manage the models and their elements
	//
	isChanged bool `json:"-" yaml:"-"`
}

// Check reviews an Element to make sure if is value.
func (e *Element) Check(buf io.Writer) bool {
	ok := true
	if e == nil {
		fmt.Fprintf(buf, "element is nil\n")
		ok = false
	}
	if e.Id == "" {
		fmt.Fprintf(buf, "element missing id\n")
		ok = false
	}
	if e.Type == "" {
		fmt.Fprintf(buf, "element, %q, missing type\n", e.Id)
		ok = false
	}
	return ok
}

// isValidVarname tests a sting confirms to Newt's naming rule.
func isValidVarname(s string) bool {
	if len(s) == 0 {
		return false
	}
	// NOTE: variable names must start with a latter and maybe followed by
	// one or more letters, digits and underscore.
	vRe := regexp.MustCompile(`^([a-zA-Z]|[a-zA-Z][0-9a-zA-Z\_]+)$`)
	return vRe.Match([]byte(s))
}

// NewModel, makes sure model id is valid, populates a Model with the oid element providing
// returns a *Model and error value.
func NewModel(modelId string) (*Model, error) {
	if ! isValidVarname(modelId) {
		return nil, fmt.Errorf("invalid model id, %q", modelId)
	}
	model := new(Model)
	model.Id = modelId
	model.Name = modelId
	model.Description = fmt.Sprintf("... description of %q goes here ...", modelId)
	model.Body = []*Element{}
	// Make the required element ...
	element := new(Element)
	element.Id = "oid"
	element.IsModelIdentifier = true
	element.Type = "input"
	element.Attributes = map[string]string{"readonly": "true"}
	if err := model.InsertElement(0, element); err != nil {
		return nil, err
	}
	return model, nil
}

// Check analyze the model and make sure at least one element exists and the
// model has a single identifier (e.g. "oid")
func (model *Model) Check(buf io.Writer) bool {
	if model == nil {
		fmt.Fprintf(buf, "model is nil\n")
		return false
	}
	if model.Body == nil {
		fmt.Fprintf(buf, "missing %s.body\n", model.Id)
		return false
	}
	// Check to see if we have at least one element in Body
	if len(model.Body) > 0 {
		ok := true
		hasModelId := false
		for i, e := range model.Body {
			// Check to make sure each element is valid
			if ! e.Check(buf) {
				fmt.Fprintf(buf, "error for %s.%s\n", model.Id, e.Id)
				ok = false
			}
			if e.IsModelIdentifier {
				if hasModelId == true {
					fmt.Fprintf(buf, "duplicate model identifier element (%d) %s.%s\n", i, model.Id, e.Id)
					ok = false
				}
				hasModelId = true
			}
		}
		if ! hasModelId {
			fmt.Fprintf(buf, "missing required object identifier for model %s\n", model.Id)
			ok = false
		}
		return ok
	}
	fmt.Fprintf(buf, "Missing elements for model %q\n", model.Id)
	return false
}

// InsertElement will add a new element to model.Body in the position indicated,
// It will also set isChanged to true on additional.
func (model *Model) InsertElement(pos int, element *Element) error {
	if model.Body == nil {
		model.Body = []*Element{}
	}
	if ! isValidVarname(element.Id) {
		return fmt.Errorf("element id is not value")
	}
	if model.HasElement(element.Id) {
		return fmt.Errorf("duplicate element id, %q", element.Id)
	}
	if pos < 0 {
		pos = 0
	}
	if pos > len(model.Body) {
		model.Body = append(model.Body, element)
		model.isChanged = true
		return nil
	}
	elements := append(model.Body[:pos], element)
	model.Body = append(elements, model.Body[(pos+1):]...)
	model.isChanged = true
	return nil
}

// UpdateElement will update an existing element with element id will the new element.
func (model *Model) UpdateElement(elementId string, element *Element) error {
	if ! model.HasElement(elementId) {
		return fmt.Errorf("%q element id not found", elementId)
	}
	for i, e := range model.Body {
		if e.Id == elementId {
			model.Body[i] = element
			model.isChanged = true
			return nil
		}
	}
	return fmt.Errorf("failed to find %q to update", elementId)
}

// RemoveElement removes an element by id from the model.Body
func (model *Model) RemoveElement(elementId string) error {
	if ! model.HasElement(elementId) {
		return fmt.Errorf("%q element id not found", elementId)
	}
	for i, e := range model.Body {
		if e.Id == elementId {
			model.Body = append(model.Body[:i], model.Body[(i+1):]...)
			model.isChanged = true
			return nil
		}
	}
	return fmt.Errorf("%q element id is missing", elementId)
}
