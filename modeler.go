package newt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

// isSupportedElementType checks if the element type is supported by Newt, returns true if OK false is it is not
func isSupportedElementType(eType string) bool {
	eType = strings.TrimSpace(strings.ToLower(eType))
	for _, sType := range []string{"text", "textarea", "button", "checkbox", "radio", "date", "phone", "url", "email", "number", "integer", "select", "submit", "reset", "orcid" } {
		if eType == sType {
			return true
		}
	}
	return false
}

// removeElement removes an element from a model
func removeElement(model *Model, in io.Reader, out io.Writer, eout io.Writer, elementId string) error {
	elemFound := false
	for i, elem := range model.Elements {
		if elem.Id == elementId {
			model.Elements = append(model.Elements[:i], model.Elements[(i+1):]...)
			model.isChanged = true
			elemFound = true
		}
	}
	if !elemFound {
		return fmt.Errorf("failed to find %s.%s", model.Id, elementId)
	}
	return nil
}

// removeModelRoutesAndTemplates removes the model and related routes and templates
func removeModelRoutesAndTemplates(ast *AST, in io.Reader, out io.Writer, eout io.Writer, modelId string) error {
	// Step 1: Remove the model
	modelFound := false
	for i, model := range ast.Models {
		if model.Id == modelId {
			ast.Models = append(ast.Models[:i], ast.Models[(i+1):]...)
			modelFound = true
		}
	}
	if !modelFound {
		return fmt.Errorf("failed to find %s", modelId)
	}
	for _, action := range []string{"create", "read", "update", "delete", "list"} {
		id := mkName(modelId, action, "")
		// Step 2: For each model remove the related routes
		if err := ast.RemoveRouteById(id); err != nil {
			fmt.Fprintf(eout, "%s\n", err)
		} else {
			ast.isChanged = true
		}
		// Step 3: For each model remove the related templates
		if err := ast.RemoveTemplateById(id); err != nil {
			fmt.Fprintf(eout, "%s\n", err)
		} else {
			ast.isChanged = true
		}
	}
	return nil
}

// saveModelsRoutesAndTemplates takes a existing config structure,
// regenerates routes and templates if needed, writes the updated structure
// to disk
func saveModelsRoutesAndTemplates(configName string, ast *AST) error {
	eBuf := bytes.NewBuffer([]byte{})
	hasError := false
	// Make sure that we have routes and template handling for each model
	for _, model := range ast.Models {
		if err := setupWebFormHandling(ast, model, "create"); err != nil {
			fmt.Fprintf(eBuf, "%s\n", err)
			hasError = true
		}
		if err := setupWebFormHandling(ast, model, "update"); err != nil {
			fmt.Fprintf(eBuf, "%s\n", err)
			hasError = true
		}
		if err := setupWebFormHandling(ast, model, "delete"); err != nil {
			fmt.Fprintf(eBuf, "%s\n", err)
			hasError = true
		}
		if err := setupReadHandling(ast, model, "read"); err != nil {
			fmt.Fprintf(eBuf, "%s\n", err)
			hasError = true
		}
		if err := setupReadHandling(ast, model, "list"); err != nil {
			fmt.Fprintf(eBuf, "%s\n", err)
			hasError = true
		}
	}
	if err := ast.SaveAs(configName); err != nil {
		fmt.Fprintf(eBuf, "%s\n", err)
		hasError = true
	}
	if hasError {
		return fmt.Errorf("%s", eBuf.Bytes())
	}
	return nil
}

// normalizeInputType takes a string and returns a normlized input type
// (e.g. lowercased and spaces trim) and if the type is an alias
// then a pattern is returned along with the normalized type.
// e.g. DATE -> date,""
//      eMail -> email, ""
//      ORCID -> orcid, "[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9A-Z]"
func normalizeInputType(inputType string) (string, string) {
	patternMap := map[string]string{
		// This is where we put the aliases to common regexp validated patterns
		"orcid": "[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9]-[0-9][0-9][0-9][0-9A-Z]",
	}
	val := strings.TrimSpace(strings.ToLower(inputType))
	if pattern, ok := patternMap[val]; ok {
		return val, pattern
	}
	return val, ""
}

// addElementStub adds an empty element to model.Elements list. Returns a new model list and error value
func addElementStub(model *Model, elementId string) ([]string, error) {
	elementList := model.GetElementIds()
	if !isValidVarname(elementId) {
		return elementList, fmt.Errorf("%q is not a valid elmeent id", elementId)
	}
	elem, err := NewElement(elementId)
	if err != nil {
		return elementList, err
	}
	model.Elements = append(model.Elements, elem)
	elementList = model.GetElementIds()
	return elementList, nil
}

// addModelStub adds an empty model to the ast.Models list. Returns
// a new model list and error value
func addModelStub(ast *AST, modelId string) ([]string, error) {
	modelList := ast.GetModelIds()
	if !isValidVarname(modelId) {
		return modelList, fmt.Errorf("%q is not a valid model id", modelId)
	}
	model, err := NewModel(modelId)
	if err != nil {
		return modelList, err
	}
	if err := ast.AddModel(model); err != nil {
		return modelList, err
	}
	modelList = ast.GetModelIds()
	sort.Strings(modelList)
	// Setup templates and routes for model
	return modelList, nil
}

// modifyModelTUI modify a specific model (e.g. add, modify remove model attributes)
func modifyModelTUI(ast *AST, in io.Reader, out io.Writer, eout io.Writer, modelId string) error {
	buf := bufio.NewReader(in)
	model, ok := ast.GetModelById(modelId)
	if !ok {
		return fmt.Errorf("cannot find %q", modelId)
	}
	for quit := false; !quit; {
		elementList := model.GetElementIds()
		//attributeList := model.GetAttributeIds()
		menu, opt := selectMenuItem(in, out,
			fmt.Sprintf("Manage %s model", modelId),
			"Menu: model [d]escription, [e]lements or press enter when done",
			[]string{
				fmt.Sprintf("Model Id:\n\t\t%s\n", model.Id),
				fmt.Sprintf("Description:\n\t\t%s\n", model.Description),
				fmt.Sprintf("Element(s):\n\t\t%s\n", strings.Join(elementList, ",\n\t\t")),
			}, false, "", "", true)
		if len(menu) > 0 {
			menu = menu[0:1]
		}
		switch menu {
		case "d":
			if opt == "" {
				fmt.Fprintf(out, "\nEnter model description: ")
				opt = getAnswer(buf, "", false)
			}
			if opt != "" {
				model.Description = strings.TrimSpace(opt)
				model.isChanged = true
			}
		case "e":
			if err := modifyElementsTUI(in, out, eout, model); err != nil {
				fmt.Fprintf(eout, "%s\n", err)
			}
		case "q":
			quit = true
		case "":
		    quit = true
		default:
			fmt.Fprintf(eout, "failed to underand %q\n", opt)
		}
	}
	return nil
}

// modifyModelAttributesTUI provides a text UI for managing a model's attributes
func modifyModelAttributesTUI(model *Model, in io.Reader, out io.Writer, eout io.Writer) error {
	buf := bufio.NewReader(in)
	if model.Attributes == nil {
		model.Attributes = map[string]string{}
	}
	for quit := false; !quit; {
		attributeList := []string{}
		for k, v := range model.Attributes{
			attributeList = append(attributeList, fmt.Sprintf("%s -> %q", k, v))
		}
		menu, opt := selectMenuItem(in, out,
			fmt.Sprintf("Manage %s attributes", model.Id), TuiStandardMenu,
			attributeList, false, "", "", true)
		if len(menu) > 0 {
			menu = menu[0:1]
		}
		var ok bool
		opt, ok = getIdFromList(attributeList, opt)
		switch menu {
		case "a":
			if opt == "" {
				fmt.Fprintf(out, "Enter attribute name: ")
				opt = getAnswer(buf, "", true)
			}
			if !isValidVarname(opt) {
				fmt.Fprintf(eout, "%q is not a valid attribute name\n", opt)
			} else {
				model.Attributes[opt] = ""
				model.isChanged = true
			}
		case "m":
			if opt == "" {
				fmt.Fprintf(out, "Enter attribute name: ")
				opt = getAnswer(buf, "", true)
			}
			fmt.Fprintf(out, "Enter %s's value: ", opt)
			val := getAnswer(buf, "", false)
			if val != "" {
				model.Attributes[opt] = val
				model.isChanged = true
			}
		case "r":
			if opt == "" {
				fmt.Fprintf(out, "Enter attribute name to remove: ")
				opt = getAnswer(buf, "", true)
				opt, ok = getIdFromList(attributeList, opt)
			}
			if ok {
				if _, ok := model.Attributes[opt]; ok {
					delete(model.Attributes, opt)
					model.isChanged = true
				}
			} 
			if ! ok {
				fmt.Fprintf(eout, "failed to find %q in attributes\n", opt)
			}
		case "q":
			quit = true
		case "":
		    quit = true
		default:
			fmt.Fprintf(eout, "failed to underand %q\n", opt)
		}
	}
	return nil
}


// modifyElementAttributesTUI provides a text UI for managing a model's element's attributes
func modifyElementAttributesTUI(model *Model, in io.Reader, out io.Writer, eout io.Writer, elementId string) error {
	buf := bufio.NewReader(in)
	elem, _ := model.GetElementById(elementId)
	for quit := false; !quit; {
		var ok bool
		attributeList := []string{}
		for k, v := range elem.Attributes {
			attributeList = append(attributeList, fmt.Sprintf("%s -> %q", k, v))
		}
		menu, opt := selectMenuItem(in, out,
			fmt.Sprintf("Manage %s.%s attributes", model.Id, elementId), TuiStandardMenu,
			attributeList, true, "", "", true)
		if len(menu) > 0 {
			menu = menu[0:1]
		}
		opt, ok = getIdFromList(attributeList, opt)
		switch menu {
		case "a":
			if opt == "" {
				fmt.Fprintf(out, "Enter attribute name: ")
				opt = getAnswer(buf, "", true)
			}
			if !isValidVarname(opt) {
				fmt.Fprintf(eout, "%q is not a valid attribute name\n", opt)
			} else {
				elem.Attributes[opt] = ""
				elem.isChanged = true
			}
		case "m":
			if opt == "" {
				fmt.Fprintf(out, "Enter attribute name: ")
				opt = getAnswer(buf, "", true)
			}
			fmt.Fprintf(out, "Enter %s's value: ", opt)
			val := getAnswer(buf, "", false)
			if val != "" {
				elem.Attributes[opt] = val
				elem.isChanged = true
			}
		case "r":
			if opt == "" {
				fmt.Fprintf(out, "Enter attribute name to remove: ")
				opt = getAnswer(buf, "", true)
				opt, ok = getIdFromList(attributeList, opt)
			}
			if ok {
				if _, ok = elem.Attributes[opt]; ok {
					delete(elem.Attributes, opt)
					elem.isChanged = true
				}
			} 
			if ! ok {
				fmt.Fprintf(eout, "failed to find %q in attributes\n", opt)
			}
		case "q":
			quit = true
		case "":
		    quit = true
		default:
			fmt.Fprintf(eout, "failed to underand %q\n", opt)
		}
	}
	return nil
}

func modifySelectElementTUI(elem *Element, in io.Reader, out io.Writer, eout io.Writer, modelId string) error {
	buf := bufio.NewReader(in)
	if elem.Options == nil {
		elem.Options = []map[string]string{}
	}
	for quit := false; !quit; {
	    optionsList := getValueLabelList(elem.Options)
		menu, opt := selectMenuItem(in, out, 
		    fmt.Sprintf("Manage %s.%s options", modelId, elem.Id),
			"Menu [a]dd, [m]odify option no., [r]emove option no. or press enter when done",
			optionsList,
			true, "", "", true)
		if len(menu) > 0 {
			menu = menu[0:1]
		}
		var (
			val string
			label string
		)
		switch menu {
			case "a":
				if opt == "" {
					fmt.Fprintf(out, "Enter an option value: ")
					val = getAnswer(buf, "", true)
				} else {
					val = strings.TrimSpace(opt)
				}
				if val == "" {
					fmt.Fprintf(eout, "Error, an option value required\n")
				} else {
					fmt.Fprintf(out, "Enter an option label: ")
					label = getAnswer(buf, "", false)
					if label == "" {
						label = val
					}
					option := map[string]string{
						val: label,
					}
					elem.Options = append(elem.Options, option)
					elem.isChanged = true
				}
			case "m":
				pos, ok := -1, false
				if i, err := strconv.Atoi(opt); err == nil {
					// Adjust to a zero based array
					if i > 0 && i <= len(optionsList) {
						pos, ok = (i - 1), true
					}
				} 
				if ! ok {
					fmt.Fprintf(out, "Enter option number: ")
					pos, ok = getDigit(buf, optionsList)

				}
				if ok {
					option := elem.Options[pos]
					val, label, _ := getValAndLabel(option)
					fmt.Fprintf(out, "Enter an option label (for %q): ", val)
					answer := getAnswer(buf, "", false)
					if answer != "" {
						label = answer
					}
					if label == "" {
						label = val
					}
					option = map[string]string{
						val: label,
					}
					// Replace the option in the options list.
					elem.Options[pos] = option
					elem.isChanged = true
				} else {
					fmt.Fprintf(eout, "Failed to identify option %q\n", opt)
				}
			case "r":
				pos, ok := -1, false
				if i, err := strconv.Atoi(opt); err == nil {
					// Adjust to a zero based array
					if i > 0 && i <= len(optionsList) {
						pos, ok = (i - 1), true
					}
				} else {
					fmt.Fprintf(out, "Enter option number: ")
					pos, ok = getDigit(buf, optionsList)
				}
				if ok {
					elem.Options = append(elem.Options[0:pos], elem.Options[(pos+1):]...)
					elem.isChanged = true
				} else {
					fmt.Fprintf(eout, "failed to remove option number (%d) %q\n", pos, opt)
				}
			case "q":
				quit = true
			case "":
		    	quit = true
			default:
				fmt.Fprintf(eout, "did not understand, %s %s\n", menu, opt)
		}
	}
	return nil
}

func modifyElementTUI(model *Model, in io.Reader, out io.Writer, eout io.Writer, elementId string) error {
	buf := bufio.NewReader(in)
	elem, ok := model.GetElementById(elementId)
	if !ok {
		return fmt.Errorf("could not find %q element", elementId)
	}
	var (
		menu string
		opt string
	)
	for quit := false; !quit; {
		attributeList := getAttributeIds(elem.Attributes)
		switch elem.Type {
		case "select":
			optionsList := getValueLabelList(elem.Options)
			menu, opt = selectMenuItem(in, out,
			fmt.Sprintf("Manage %s.%s element", model.Id, elementId),
			"Menu [t]ype, [l]abel, [a]ttributes, [o]ptions or press enter when done",
			[]string{
				fmt.Sprintf("id %s", elementId),
				fmt.Sprintf("type %s", elem.Type),
				fmt.Sprintf("label %s", elem.Label),
				fmt.Sprintf("attributes %s", strings.Join(attributeList, ",\n\t\t")),
				fmt.Sprintf("options %s", strings.Join(optionsList, ",\n\t\t")),
			},
			false, "", "", true)
		case "textarea":
			menu, opt = selectMenuItem(in, out,
			fmt.Sprintf("Manage %s.%s element", model.Id, elementId),
			"Menu [t]ype, [l]abel, [a]ttributes or press enter when done",
			[]string{
				fmt.Sprintf("id %s", elementId),
				fmt.Sprintf("type %s", elem.Type),
				fmt.Sprintf("label %s", elem.Label),
				fmt.Sprintf("attributes %s", strings.Join(attributeList, ",\n\t\t")),
			},
			false, "", "", true)
		default:
			menu, opt = selectMenuItem(in, out,
			fmt.Sprintf("Manage %s.%s element", model.Id, elementId),
			"Menu [t]ype, [l]abel, [o]bject identifier, [p]attern, [a]ttributes, or press enter when done",
			[]string{
				fmt.Sprintf("id %s", elementId),
				fmt.Sprintf("type %s", elem.Type),
				fmt.Sprintf("label %s", elem.Label),
				fmt.Sprintf("pattern %s", elem.Pattern),
				fmt.Sprintf("attributes %s", strings.Join(attributeList, ",\n\t\t")),
				fmt.Sprintf("object identifier? %t", elem.IsObjectId),
			},
			false, "", "", true)
		}
		if len(menu) > 0 {
			menu = menu[0:1]
		}
		switch menu {
		case "t":
			if opt == "" {
				fmt.Fprintf(out, `Enter type string (e.g. text, email, date, textarea, select, orcid): `)
				opt = getAnswer(buf, "", false)
			}
			if opt != "" {
				eType := strings.TrimSpace(opt)
				if ok := isSupportedElementType(eType); ! ok {
					fmt.Fprintf(eout, "%q is not a supported element type", opt)
				} else {
					// FIXME: If elem.Type is select I need to provide an choices TUI for values and labels
					elem.Type, elem.Pattern = normalizeInputType(eType)
					elem.isChanged = true
					if elem.Type == "select" {
						if err := modifySelectElementTUI(elem, in, out, eout, model.Id); err != nil {
							fmt.Fprintf(eout, "ERROR (%q): %s\n", elementId, err)
						}
					}
				}
			}
		case "p":
			if opt == "" {
				fmt.Fprintf(out, "Enter the regexp to use: ")
				opt = getAnswer(buf, "", false)
			}
			// FIXME: how do I clear a pattern if I nolonger want to use one?
			if opt != "" {
				// NOTE: remove a pattern by having it match everything, i.e. astrix
				if opt == "*" {
					elem.Pattern = ""
				} else {
					elem.Pattern = strings.TrimSpace(opt)
				}
				model.isChanged = true
			}
		case "l":
			if opt == "" {
				fmt.Fprintf(out, "Enter a label: ")
				opt = getAnswer(buf, "", false)
			}
			if opt != "" {
				// NOTE: remove a pattern by having it match everything, i.e. astrix
				if opt != elem.Label {
					elem.Label = strings.TrimSpace(opt)
				}
				model.isChanged = true
			}
		case "a":
			if err := modifyElementAttributesTUI(model, in, out, eout, elementId); err != nil {
				fmt.Fprintf(eout, "%s\n", err)
			}
		case "o":
			if elem.Type == "select" {
				if err := modifySelectElementTUI(elem, in, out, eout, model.Id); err != nil {
					fmt.Fprintf(eout, "ERROR (%q): %s\n", elem.Id, err)
				}
			} else {
				elem.IsObjectId = !elem.IsObjectId
				elem.isChanged = true
			}
		case "q":
			quit = true
		case "":
			quit = true
		default:
			fmt.Fprintf(eout, "did not understand %q\n", menu)
		}
	}
	return nil
}

func removeElementFromModel(model *Model, elementId string) error {
	for i, elem := range model.Elements {
		if elem.Id == elementId {
			model.Elements = append(model.Elements[0:i], model.Elements[(i+1):]...)
			model.isChanged = true
			return nil
		}
	}
	return fmt.Errorf("%s not found", elementId)
}

// modifyElementsTUI modify a specific model's element list.
func modifyElementsTUI(in io.Reader, out io.Writer, eout io.Writer, model *Model) error {
	var (
		err    error
		answer string
	)
	buf := bufio.NewReader(in)
	// FIXME: Need to support editing model attributes, then allow for modifying model's body to be modified.
	for quit := false; !quit; {
		elementList := model.GetElementIds()
		menu, opt := selectMenuItem(in, out,
			fmt.Sprintf("Manage %s elements", model.Id), TuiStandardMenu,
			elementList, true, "", "", true)
		if len(menu) > 1 {
			menu = menu[0:1]
		}
		elementId, ok := getIdFromList(elementList, opt)
		switch menu {
		case "a":
			if !ok {
				fmt.Fprintf(out, "Enter element id to add: ")
				opt = getAnswer(buf, "", false)
				elementId, ok = getIdFromList(elementList, opt)
			}
			if ok {
				elementList, err = addElementStub(model, elementId)
				if err != nil {
					fmt.Fprintf(eout, "WARNING: %s\n", err)
				}
			}
		case "m":
			if elementId == "" {
				fmt.Fprintf(out, "Enter element id to modify: ")
				opt = getAnswer(buf, "", false)
				elementId, ok = getIdFromList(elementList, opt)
			}
			if err := modifyElementTUI(model, in, out, eout, elementId); err != nil {
				fmt.Fprintf(eout, "ERROR (%q): %s\n", elementId, err)
			}
		case "r":
			if elementId == "" {
				fmt.Fprintf(out, "Enter element id to remove: ")
				opt = getAnswer(buf, "", false)
				elementId, ok = getIdFromList(elementList, opt)
			}
			if ok {
				if err := removeElementFromModel(model, elementId); err != nil {
					fmt.Fprintf(eout, "ERROR (%q): %s\n", elementId, err)
				}
			} 
		case "q":
			quit = true
		case "":
			quit = true
		default:
			fmt.Fprintf(eout, "\n\nERROR: Did not understand %q\n\n", answer)
		}

	}
	return nil
}

// modelerTUI takes configuration and then runs the interactive text user interface modeler.
func modelerTUI(ast *AST, in io.Reader, out io.Writer, eout io.Writer, configName string, newModelIds []string) error {
	var (
		err    error
		answer string
	)
	buf := bufio.NewReader(in)
	if ast.Models == nil {
		ast.Models = []*Model{}
	}
	modelList := ast.GetModelIds()
	sort.Strings(modelList)
	if len(newModelIds) > 0 {
		for _, modelId := range newModelIds {
			modelList, err = addModelStub(ast, modelId)
			if err != nil {
				fmt.Fprintf(eout, "WARNING: %s\n", err)
			}
		}
	}
	for quit := false; !quit; {
		menu, opt := selectMenuItem(in, out,
			"Manage Models",
			TuiStandardMenu,
			modelList, true, "", "", true)
		if len(menu) > 1 {
			menu = menu[0:1]
		}
		modelId, ok := getIdFromList(modelList, opt)
		switch menu {
		case "a":
			if !ok {
				fmt.Fprintf(out, "Enter model id to add: ")
				opt = getAnswer(buf, "", false)
				modelId, ok = getIdFromList(modelList, opt)
			}
			if ok {
				modelList, err = addModelStub(ast, modelId)
				if err != nil {
					fmt.Fprintf(eout, "WARNING: %s\n", err)
				}
			}
		case "m":
			if !ok {
				fmt.Fprintf(out, "Enter model id to modify: ")
				opt = getAnswer(buf, "", false)
				modelId, ok = getIdFromList(modelList, opt)
			}
			if err := modifyModelTUI(ast, in, out, eout, modelId); err != nil {
				fmt.Fprintf(eout, "ERROR (%q): %s\n", modelId, err)
			}
		case "r":
			if !ok {
				fmt.Fprintf(out, "Enter model id to remove: ")
				opt = getAnswer(buf, "", false)
				modelId, ok = getIdFromList(modelList, opt)
			}
			if err := removeModelRoutesAndTemplates(ast, in, out, eout, modelId); err != nil {
				fmt.Fprintf(eout, "ERROR (%q): %s\n", modelId, err)
			}
			modelList = ast.GetModelIds()
		case "q":
			quit = true
		case "":
			quit = true
		default:
			fmt.Fprintf(eout, "\n\nERROR: Did not understand %q\n\n", answer)
		}
	}
	if ast.HasChanges() {
		fmt.Fprintf(out, "Save before exiting (Y/n)? ")
		answer = getAnswer(buf, "y", true)
		if answer == "y" {
			if err := saveModelsRoutesAndTemplates(configName, ast); err != nil {
				return err
			}
		}
	}
	return nil
}
