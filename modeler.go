package newt

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// refreshTemplates will look at the list of models in the config
// and make sure we have a current set of templates for each one.
func refreshTemplates(ast *AST) error {
	return fmt.Errorf("FIXME: refreshTemplates(ast) not implemented")
}

// refreshRoutes will look at the list of models and templates in the config and
// make sure we have a current set of related routes for model/template set
func refreshRoutes(ast *AST) error {
	return fmt.Errorf("FIXME: refreshRoutes(ast) not implemented")
}

// removeModelRoutesAndTemplates removes the model and related routes and templates
func removeModelRoutesAndTemplates (ast *AST, modelId string) error {
	return fmt.Errorf("FIXME: removeModelRoutesAndTemplates(ast, %q) not implemented", modelId)
}

// saveModelsRoutesAndTemplates takes a existing config structure,
// regenerates routes and templates if needed, writes the updated structure
// to disk
func saveModelsRoutesAndTemplates(configName string, ast *AST) error {
	// FIXME: Make we have templates our models, update template if needed
	if err := refreshTemplates(ast); err != nil {
		return err
	}
	// FIXME: Make sure that we have routes defined for each model and template
	if err := refreshRoutes(ast); err != nil {
		return err
	}
	if err := ast.SaveAs(configName); err != nil {
		return err
	}
	return nil
}



// addModelStub adds an empty model to the ast.Models list. Returns
// a new model list and error value
func addModelStub(ast *AST, modelId string) ([]string, error) {
	modelList := ast.GetModelIds()
	if ! isValidVarname(modelId) {
		return modelList, fmt.Errorf("%q is not a valid model name", modelId)
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
	return modelList, nil
}

func getModelName(modelList []string, modelId string) (string, bool) {
	// NOTE: nRe tests if modelId is a string representation of a positive integer 
	nRe := regexp.MustCompile(`^[0-9]+$`)
	// See if we have been given a model number or a name
	if isDigit := nRe.Match([]byte(modelId)); isDigit {
		mNo, err := strconv.Atoi(modelId)
		if err == nil {
			// Adjust provided integer for zero based index.
			if mNo > 0 {
				mNo--
			} else {
				mNo = 0
			}
			if mNo < len(modelList) {
				return modelList[mNo], true
			}
		}
	}
	if isValidVarname(modelId) {
		return  modelId, true
	}
	return "", false
}

// modifyModelTUI modify a specific model (e.g. add, modify remove model attributes)
func modifyModelTUI(ast *AST, in io.Reader, out io.Writer, eout io.Writer, modelId string) error {
	return fmt.Errorf("FIXME modify model %q not implemented", modelId)
}

// modelerTUI takes configuration and then runs the interactive text user interface modeler.
func modelerTUI(ast *AST, in io.Reader, out io.Writer, eout io.Writer, configName string, newModels []string) error {
	var (
		err error
		answer string
	)
	readBuffer := bufio.NewReader(in)
	modelList := ast.GetModelIds()
	sort.Strings(modelList)
	if len(newModels) > 0 {
		for _, modelId := range newModels {
			modelList, err = addModelStub(ast, modelId)
			if err != nil {
				fmt.Fprintf(eout, "WARNING: %s\n", err)
			}
		}
	}
	for quit := false; ! quit; {
		fmt.Fprintf(out, "Enter menu menu command and model id\n\n")
		if len(modelList) == 0 {
			fmt.Fprintf(out, "\tNO MODELS DEFINED\n")
		} else {
			for i, modelId := range modelList {
				fmt.Fprintf(out, "\t%3d: %s\n", i+1, modelId)
			}
		}
		fmt.Fprintf(out, "\nMenu [a]dd, [m]odify, [r]emove, [s]ave YAML file, [q]uit\n")
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
		modelId, ok := getModelName(modelList, parts[1])
		switch menu {
			case "a":
			    if ! ok {
					fmt.Fprintf(out, "Enter model id to add: ")
					answer = getAnswer(readBuffer, "", false)
					modelId, ok = getModelName(modelList, answer)
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
					modelId, ok = getModelName(modelList, answer)
				}
				if err := modifyModelTUI(ast, in, out, eout, modelId); err != nil {
					fmt.Fprintf(eout, "ERROR (%q): %s\n", modelId, err)
				}
			case "r":
			    if modelId == "" {
					fmt.Fprintf(out, "Enter model id to remove: ")
					answer = getAnswer(readBuffer, "", false)
					modelId, ok = getModelName(modelList, answer)
				}
				if err := removeModelRoutesAndTemplates(ast, modelId); err != nil {
					fmt.Fprintf(eout, "ERROR (%q): %s\n", modelId, err)
				}
				modelList = ast.GetModelIds()
			case "s":
				if err := saveModelsRoutesAndTemplates(configName, ast); err != nil{
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
