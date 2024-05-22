package newt

import (
	"fmt"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
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
	if len(ids) > 0 {
		sort.Strings(ids)
	}
	return ids
}

// Get return the first key and value pair
func getValAndLabel(option map[string]string) (string, string, bool) {
	for val, label := range option {
		return val, label, true
	}
	return "", "", false
}


// getValueLabelList takes an array of map[string]string and yours a list of
// strings indicating the value and label
func getValueLabelList(list []map[string]string) []string {
	options := []string{}
	for _, m := range list {
		val, label, ok := getValAndLabel(m)
		if ok {
			options = append(options, fmt.Sprintf("%s %s", val, label))
		}
	}
	return options
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
				if strings.Contains(list[mNo], " ") {
					parts := strings.SplitN(list[mNo], " ", 2)
					return parts[0], true
				}
				return list[mNo], true
			}
		}
	}
	if isValidVarname(id) {
		return id, true
	}
	return "", false
}

// getItemNoFromList returns the zero based item number from list of items.
// If an item number is provided then it will be adjusted to zero based and
// validated, if a string is provided it will the position in the list that
// is returned.
func getItemNoFromList(list []string, id string) (int, bool) {
	// Handle case where ones based number has been provided as a string.
	nRe := regexp.MustCompile(`^[0-9]+$`)
	if isDigit := nRe.Match([]byte(id)); isDigit {
		itemNo, err := strconv.Atoi(id)
		if err == nil {
			itemNo--
			if itemNo >= 0 && itemNo < len(list) {
				return itemNo, true
			}
		}
		return -1, false
	}
	// Handle case where the string needs to be matched.
	for i, val := range list {
		if val == id {
			return i, true
		}
	}
	return -1, false
}

// fNameToNamespace trunes a filename (e.g. app.yaml) into a namespace (e.g. app)
func fNameToNamespace(s string) string {
	bName := path.Base(s)
	ext := path.Ext(bName)
	if len(ext) > 0 {
		return strings.TrimSuffix(bName, ext)	
	}
	return bName
}
