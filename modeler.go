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

// addModelStub adds an empty model to the cfg.Models list. Returns
// a new model list and error value
func addModelStub(cfg *Config, modelId string) ([]string, error) {
	modelList := cfg.GetModelIds()
	if ! isValidVarname(modelId) {
		return modelList, fmt.Errorf("%q is not a valid model name", modelId)
	}
	// Make sure we not adding a duplicate
	_, found := sort.Find(len(modelList), func(i int) int {
   			return strings.Compare(modelId, modelList[i])
	})
	if found {
		return modelList, fmt.Errorf("%q is a duplicate name", modelId)
	}
	modelList = append(modelList, modelId)
	sort.Strings(modelList)
	model := new(NewtModel)
	model.Id = modelId
	model.Name = modelId
	model.Description = fmt.Sprintf("... description of %q goes here ...", modelId)
	model.Body = []*Element{}
	element := new(Element)
	element.Id = "oid"
	element.Type = "input"
	element.Attributes = map[string]string{"readonly": "true"}
	element.Validations = map[string]interface{}{ "retired": true}
	model.Body = append(model.Body, element)
	if cfg.Models == nil {
		cfg.Models = []*NewtModel{}
	}
	cfg.Models = append(cfg.Models, model)
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
func modifyModelTUI(cfg *Config, in io.Reader, out io.Writer, eout io.Writer, modelId string) error {
	return fmt.Errorf("DEBUG modify model %q not implemented", modelId)
}

// modelerTUI takes configuration and then runs the interactive text user interface modeler.
func modelerTUI(cfg *Config, in io.Reader, out io.Writer, eout io.Writer, configName string, newModels []string) error {
	var (
		err error
		answer string
	)
	readBuffer := bufio.NewReader(in)
	modelList := cfg.GetModelIds()
	sort.Strings(modelList)
	if len(newModels) > 0 {
		for _, modelId := range newModels {
			modelList, err = addModelStub(cfg, modelId)
			if err != nil {
				fmt.Fprintf(eout, "WARNING: %s\n", err)
			}
		}
	}
	notSaved := false
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
					modelList, err = addModelStub(cfg, modelId)
					if err != nil {
						fmt.Fprintf(eout, "WARNING: %s\n", err)
					}
				}
				// FIXME: Add templates and routes based on tis new model.
				notSaved = true
			case "m":
			    if modelId == "" {
					fmt.Fprintf(out, "Enter model id to modify: ")
					answer = getAnswer(readBuffer, "", false)
					modelId, ok = getModelName(modelList, answer)
				}
				if err := modifyModelTUI(cfg, in, out, eout, modelId); err != nil {
					fmt.Fprintf(eout, "ERROR (%q): %s\n", modelId, err)
				}
				// FIXME: Update templates because this model changed.
				notSaved = true
			case "r":
			    if modelId == "" {
					fmt.Fprintf(out, "Enter model id to remove: ")
					answer = getAnswer(readBuffer, "", false)
					modelId, ok = getModelName(modelList, answer)
				}
				if err := cfg.RemoveModelById(modelId); err != nil {
					fmt.Fprintf(eout, "ERROR (%q): %s\n", modelId, err)
				}
				// FIXME: Remove templates and routes related to this model too.
				modelList = cfg.GetModelIds()
				notSaved = true
			case "s":
				if err := cfg.SaveConfig(configName); err != nil {
					fmt.Fprintf(eout, "ERROR: %s\n", err)
				}
				notSaved = false
			case "q":
				quit = true
			case "":
			// do nothing, display list
			default:
				fmt.Fprintf(eout, "\n\nERROR: Did not understand %q\n\n", answer)
		}
	}
	if notSaved {
		fmt.Fprintf(out, "Save before exiting (Y/n)? ")
		answer = getAnswer(readBuffer, "y", true)
		if answer == "y" {
			if cfg.SaveConfig(configName); err != nil {
				return fmt.Errorf("failed to save changes to %q, %s\n", configName, err)
			}
		}
	}
	return nil
} 
