package newt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// removeElement removes an element from a model
func removeElement(model *Model, in io.Reader, out io.Writer, eout io.Writer, elementId string) error {
	elemFound := false
	for i, elem := range model.Body {
		if elem.Id == elementId {
			model.Body = append(model.Body[:i], model.Body[(i+1):]...)
			model.isChanged = true
			elemFound = true
		}
	}
	if ! elemFound {
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
	if ! modelFound {
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

// addElementStub adds an empty element to model.Body list. Returns a new model list and error value
func addElementStub(model *Model, elementId string) ([]string, error) {
	elementList := model.GetElementIds()
	if !isValidVarname(elementId) {
		return elementList, fmt.Errorf("%q is not a valid elmeent id", elementId)
	}
	elem, err := NewElement(elementId)
	if err != nil {
		return elementList, err
	}
	model.Body = append(model.Body, elem)
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

func getIdFromList(list []string, id string) (string, bool) {
	// NOTE: nRe tests if modelId is a string representation of a positive integer
	nRe := regexp.MustCompile(`^[0-9]+$`)
	// See if we have been given a model number or a name
	if isDigit := nRe.Match([]byte(id)); isDigit {
		mNo, err := strconv.Atoi(id)
		if err == nil {
			// Adjust provided integer for zero based index.
			if mNo > 0 {
				mNo--
			} else {
				mNo = 0
			}
			if mNo < len(list) {
				return list[mNo], true
			}
		}
	}
	if isValidVarname(id) {
		return id, true
	}
	return "", false
}

// modifyModelTUI modify a specific model (e.g. add, modify remove model attributes)
func modifyModelTUI(ast *AST, in io.Reader, out io.Writer, eout io.Writer, modelId string) error {
    var (
		answer string	
	)
	readBuffer := bufio.NewReader(in)
	model, ok := ast.GetModelById(modelId)
	if ! ok {
		return fmt.Errorf("cannot find %q", modelId)
	}
	for quit := false; !quit; {
		elementList := model.GetElementIds()
		fmt.Fprintf(out, "Pick a model attribute to change using the menu options\n\n")
		fmt.Fprintf(out, "   [N]ame %q\n", model.Name)
		fmt.Fprintf(out, "   [D]escription %q\n", model.Description)
		fmt.Fprintf(out, "   [E]lements (%s)\n", strings.Join(elementList, ", "))
		fmt.Fprintf(out, "\n   [q]uit editing\n")
		answer = getAnswer(readBuffer, "", true)
		switch answer {
			case "n":
				fmt.Fprintf(out, "\nEnter model name: ")
				answer = getAnswer(readBuffer, "", false)
				if answer != "" {
					model.Name = answer
				}
				fmt.Fprintln(out, "")
			case "d":
				fmt.Fprintf(out, "\nEnter model description: ")
				answer = getAnswer(readBuffer, "", false)
				if answer != "" {
					model.Description = answer
				}
				fmt.Fprintln(out, "")
			case "e":
				if err := modifyElementsTUI(in, out, eout, model); err != nil {
					fmt.Fprintf(eout, "%s\n", err)
				}
			case "q":
				quit = true
		}
	}
	return nil
} 

func modifyAttributesTUI(model *Model, in io.Reader, out io.Writer, eout io.Writer, elementId string) error {
	return fmt.Errorf("modifyAttributesTUI() not implemented")
}

func modifyValidationsTUI(model *Model, in io.Reader, out io.Writer, eout io.Writer, elementId string) error {
	return fmt.Errorf("modifyAttributesTUI() not implemented")
}

func modifyElementTUI(model *Model, in io.Reader, out io.Writer, eout io.Writer, elementId string) error {
	var (
		answer string
	)
	readBuffer := bufio.NewReader(in)
	elem, ok := model.GetElementById(elementId)
	if ! ok {
		return fmt.Errorf("could not find %q element", elementId)
	}
	for quit := false; !quit; {
		fmt.Fprintf(out, "Enter the menu letter to modify element\n\n")
		fmt.Fprintf(out, "id %s\n", elementId)
		fmt.Fprintf(out, "[t]ype %s\n", elem.Type)
		fmt.Fprintf(out, "[a]ttributes\n")
		if len(elem.Attributes) == 0 {
			fmt.Fprintf(out, "\t NO ATTRIBUTES SET\n")
		} else {
			for k, v := range elem.Attributes {
				fmt.Fprintf(out, "\t%s: %s\n", k, v)
			}
		}
		fmt.Fprintf(out, "\n[v]alidations\n")
		if len(elem.Validations) == 0 {
			fmt.Fprintf(out, "\t NO VALIDATIONS SET\n")
		} else {
			for k, v := range elem.Validations {
				fmt.Fprintf(out, "\t%s: %+v\n", k, v)
			}
		}
		fmt.Fprintf(out, "\n")
		fmt.Fprintf(out, "[m]odel identifier is set to %t\n", elem.IsModelIdentifier)
		fmt.Fprintf(out, "\n\n[q]uit editing\n")
		answer = getAnswer(readBuffer, "", true)
		switch answer {
			case "t":
				fmt.Fprintf(out, `Enter type string (e.g. input, input[type=date]): `)
				answer = getAnswer(readBuffer, "", false)
				if answer != "" {
					// FIXME: should probably validate this ..., should probably be a
					// controlled vocabulary of supported type strings. e.g. button, checkbox, ORCID, ROR, etc
					elem.Type = answer
					elem.isChanged = true
				}
			case "a":
				if err := modifyAttributesTUI(model, in, out, eout, elementId); err != nil {
					fmt.Fprintf(eout, "%s\n", err)
				}
			case "v":
				if err := modifyValidationsTUI(model, in, out, eout, elementId); err != nil {
					fmt.Fprintf(eout, "%s\n", err)
				}
			case "m":
				elem.IsModelIdentifier = ! elem.IsModelIdentifier
				elem.isChanged = true
			case "q":
				quit = true
			case "":
				// do nothing
			default:
				fmt.Fprintf(eout, "did not understand %q\n", answer)
		}
	}
	return fmt.Errorf("modifyElementTUI() not implemented")
}

func removeElementFromModel(model *Model, elementId string) error {
	return fmt.Errorf("removeElementFromModel() not implemented")
}

// modifyElementTUI modify a specific model's element list.
func modifyElementsTUI(in io.Reader, out io.Writer, eout io.Writer, model *Model) error {
    var (
		err error
		answer string	
	)
	readBuffer := bufio.NewReader(in)
	// FIXME: Need to support editing model attributes, then allow for modifying model's body to be modified.
	for quit := false; !quit; {
		elementList := model.GetElementIds()
		fmt.Fprintf(out, "Enter menu command and element id\n\n")
		if len(elementList) == 0 {
			fmt.Fprintf(out, "\tNO ELEMENTS DEFINED\n")
		} else {
			for i, id := range elementList {
				fmt.Fprintf(out, "\t%3d: %s\n", i+1, id)
			}
		}
		fmt.Fprintf(out, "\nMenu [a]dd, [m]odify, [r]emove, [q]uit editing\n")
		answer = getAnswer(readBuffer, "", false)
		var (
			menu string
		)
		// Split answer into menu command and optional model name value (enforce that it is always a two cell slice)
		parts := (append(strings.SplitN(answer, " ", 2), "", ""))[0:2]
		menu = parts[0]
		if len(menu) > 1 {
			menu = menu[0:1]
		}
		elementId, ok := getIdFromList(elementList, parts[1])
		switch menu {
		case "a":
			if !ok {
				fmt.Fprintf(out, "Enter element id to add: ")
				answer = getAnswer(readBuffer, "", false)
				elementId, ok = getIdFromList(elementList, answer)
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
				answer = getAnswer(readBuffer, "", false)
				elementId, ok = getIdFromList(elementList, answer)
			}
			if err := modifyElementTUI(model, in, out, eout, elementId); err != nil {
				fmt.Fprintf(eout, "ERROR (%q): %s\n", elementId, err)
			}
		case "r":
			if elementId == "" {
				fmt.Fprintf(out, "Enter element id to remove: ")
				answer = getAnswer(readBuffer, "", false)
				elementId, ok = getIdFromList(elementList, answer)
			}
			if err := removeElementFromModel(model, elementId); err != nil {
				fmt.Fprintf(eout, "ERROR (%q): %s\n", elementId, err)
			}
			elementList = model.GetElementIds()
		case "q":
			quit = true
		case "":
			// do nothing, redisplay list
			elementList = model.GetElementIds()
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
	readBuffer := bufio.NewReader(in)
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
		fmt.Fprintf(out, "Enter menu command and model id\n\n")
		if len(modelList) == 0 {
			fmt.Fprintf(out, "\tNO MODELS DEFINED\n")
		} else {
			for i, modelId := range modelList {
				fmt.Fprintf(out, "\t%3d: %s\n", i+1, modelId)
			}
		}
		fmt.Fprintf(out, "\nMenu [a]dd, [m]odify, [r]emove, [s]ave, [q]uit editing\n")
		answer = getAnswer(readBuffer, "", false)
		var (
			menu string
		)
		// Split answer into menu command and optional model name value (enforce that it is always a two cell slice)
		parts := (append(strings.SplitN(answer, " ", 2), "", ""))[0:2]
		menu = parts[0]
		if len(menu) > 1 {
			menu = menu[0:1]
		}
		modelId, ok := getIdFromList(modelList, parts[1])
		switch menu {
		case "a":
			if !ok {
				fmt.Fprintf(out, "Enter model id to add: ")
				answer = getAnswer(readBuffer, "", false)
				modelId, ok = getIdFromList(modelList, answer)
			}
			if ok {
				modelList, err = addModelStub(ast, modelId)
				if err != nil {
					fmt.Fprintf(eout, "WARNING: %s\n", err)
				}
			}
		case "m":
			if modelId == "" {
				fmt.Fprintf(out, "Enter model id to modify: ")
				answer = getAnswer(readBuffer, "", false)
				modelId, ok = getIdFromList(modelList, answer)
			}
			if err := modifyModelTUI(ast, in, out, eout, modelId); err != nil {
				fmt.Fprintf(eout, "ERROR (%q): %s\n", modelId, err)
			}
		case "r":
			if modelId == "" {
				fmt.Fprintf(out, "Enter model id to remove: ")
				answer = getAnswer(readBuffer, "", false)
				modelId, ok = getIdFromList(modelList, answer)
			}
			if err := removeModelRoutesAndTemplates(ast, in, out, eout, modelId); err != nil {
				fmt.Fprintf(eout, "ERROR (%q): %s\n", modelId, err)
			}
			modelList = ast.GetModelIds()
		case "s":
			if err := saveModelsRoutesAndTemplates(configName, ast); err != nil {
				fmt.Fprintf(eout, "ERROR: %s\n", err)
			}
		case "q":
			quit = true
		case "":
		// do nothing, display list
		default:
			fmt.Fprintf(eout, "\n\nERROR: Did not understand %q\n\n", answer)
		}
	}
	if ast.HasChanges() {
		fmt.Fprintf(out, "Save before exiting (Y/n)? ")
		answer = getAnswer(readBuffer, "y", true)
		if answer == "y" {
			if err := saveModelsRoutesAndTemplates(configName, ast); err != nil {
				return err
			}
		}
	}
	return nil
}
