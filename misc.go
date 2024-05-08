package newt

import (
	"fmt"
)

//
// These are misc. functions used in various source files defining Newt
//

// mkName assmebles a name or id from a model, action and possible suffix
func mkName(modelId string, action string, suffix string) string {
	return fmt.Sprintf("%s_%s%s", modelId, action, suffix)
}

// inList takes a string and compares it with a list of strings. It
// returns true when a match is found, false otherwise
func inList(target string, list []string) bool {
	for _, val := range list {
		if val == target {
			return true
		}
	}
	return false
}

// getAttributeIds returns a list of attribue keys in a maps[string]interface{} structure
func getAttributeIds(m map[string]string) []string {
	ids := []string{}
	for k, _ := range m {
		if k != "" {
			ids = append(ids, k)
		}
	}
	return ids
}


