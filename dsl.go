package newt

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
)

const (
	StartVar = "${"
	EndVar   = "}"
)

// EvalType is the function that is envoked with the data type
// expression and value returning the value as a string and a
// bool that indicates a success or failure in the evalutation.
type EvalType func(string, string) (string, bool)

// ModelDSL describe an object's attributes and type. This is
// analagous to a tables's field definitions.
type ModelDSL struct {
	// Name of model, this may be translated into the table name when
	// rendering SQL
	Name string `json:"name,required" yaml:"name,required"`
	// Var a map of key/values where key is a variable name and value
	// is the data type.
	Var map[string]string `json:"var,omitempty" yaml:"var,omitempty"`
}


// RouteDSL holds the attributes need to decode
// a route DSL expression, match and decode against path values.
type RouteDSL struct {
	Src  string   `json:"src"`
	Dirs []string `json:"dirs,omitempty"`
	Base string   `json:"base,omitempty"`
	Ext  string   `json:"ext,omitempty"`
	// VarToType maps the variable name to a var defn
	// This is based on the vars defined in the route.
	VarToType map[string]string `json:"var_to_types,omitempty"`
	// Types maps type implementation description
	Types map[string]string `json:"-"`
	// Type name to function to Eval function (validates a variable's
	// value and extracts a value)
	TypeFn map[string]EvalType `json:"-"`
}

func (rdsl *RouteDSL) String() string {
	src, _ := json.MarshalIndent(rdsl, "", "    ")
	return string(src)
}

func getVarname(elem string) string {
	return strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(elem, StartVar), EndVar))
}

// NewRouteDSL takes a DSL expression and returns a
// DSLExpresion structure and error value.
func NewRouteDSL(src string, varDef map[string]string) (*RouteDSL, error) {
	dir, base := path.Split(src)

	rdsl := new(RouteDSL)
	rdsl.Src = src
	rdsl.Base = base
	rdsl.Ext = ""
	// Include all the types defined dsl_types.go
	rdsl.TypeFn = map[string]EvalType{}
	for k, v := range DataTypes {
		rdsl.TypeFn[k] = v
	}

	// Load a variable definitions.
	rdsl.VarToType = map[string]string{}
	if len(varDef) > 0 {
		for k, v := range varDef {
			if v == "Basename"{ 
				rdsl.Base = fmt.Sprintf("%s%s%s", StartVar, k, EndVar)
			} 
			if v == "Extname" {
				rdsl.Ext = fmt.Sprintf("%s%s%s", StartVar, k, EndVar)
			}
			rdsl.VarToType[k] = v
		}
	}
	
	// We only evalaute the extension if here are two variables 
	// defined for the last element of path.
	dirs := strings.Split(strings.TrimSuffix(strings.TrimPrefix(dir, "/"), "/"), "/")
	rdsl.Dirs = []string{}
	for _, elem := range dirs {
		if strings.HasPrefix(elem, StartVar) && strings.HasSuffix(elem, EndVar) {
			// Check to see if getVarname is defined, other error out
			vName := getVarname(elem)
			if _, ok := rdsl.VarToType[vName]; ok {
				rdsl.Dirs = append(rdsl.Dirs, elem)
			} else {
				return nil, fmt.Errorf("%s is undefined", vName)
			}
		} else {
			rdsl.Dirs = append(rdsl.Dirs, elem)
		}
	}
	if rdsl.Ext != "" {
		if strings.HasPrefix(rdsl.Ext, StartVar) && strings.HasSuffix(rdsl.Ext, EndVar) {
			vName := getVarname(rdsl.Ext)
			if _, ok := rdsl.VarToType[vName]; ! ok {
				return nil, fmt.Errorf("%s is undefined", vName)
			}
		}
	}
	if rdsl.Base != "" {
		if strings.HasPrefix(rdsl.Ext, StartVar) && strings.HasSuffix(rdsl.Ext, EndVar) {
			vName := getVarname(rdsl.Base)
			if _, ok := rdsl.VarToType[vName]; ! ok {
				return nil, fmt.Errorf("%s is undefined", vName)
			}
		}
	}
	return rdsl, nil
}

// evalElement takes compares a element against a value (from the path value)
// returns a variable name, interface value and bool indicating a successful match
func (rdsl *RouteDSL) evalElement(elem string, src string) (string, string, bool) {
	// Check if we workingwith a literal element or a variable defn.
	if strings.HasPrefix(elem, StartVar) {
		// handle variable path element
		vName := getVarname(elem)
		tExpr, ok := rdsl.VarToType[vName]
		if !ok {
			return "", "", false
		}
		fn, ok := rdsl.TypeFn[tExpr]
		if !ok {
			return vName, "", false
		}
		// Now check the type of dDir against the type expression
		val, ok := fn(tExpr, src)
		if !ok {
			// Something went wrong, path does not match.
			return "", "", false
		}
		return vName, val, true
	}
	// handle literal path element
	return "", "", (strings.Compare(elem, src) == 0)
}

// Eval takes a path value and compares it with a Path expression.
// It returns a status boolean, map of variable names to values.
func (rdsl *RouteDSL) Eval(pathValue string) (map[string]string, bool) {
	dir, base := path.Split(pathValue)
	pDirs := strings.Split(strings.TrimSuffix(strings.TrimPrefix(dir, "/"), "/"), "/")
	pExt := path.Ext(base)
	pBase := strings.TrimSuffix(base, pExt)
	if rdsl.Ext == "" {
		pExt = ""
		pBase = base
	}
	m := map[string]string{}

	// Match the extension, if it contains a
	if rdsl.Ext != "" {
		if vName, val, ok := rdsl.evalElement(rdsl.Ext, pExt); ok {
			// Check if we need to store the variable
			if vName != "" {
				m[vName] = val
			}
		} else {
			return nil, false
		}
	}
	// Match Basename
	vName, val, ok := rdsl.evalElement(rdsl.Base, pBase)
	if ok {
		// Check if we need to store the variable
		if vName != "" {
			m[vName] = val
		}
	} else {
		return nil, false
	}
	if len(pDirs) != len(rdsl.Dirs) {
		return nil, false
	}
	for i, elem := range rdsl.Dirs {
		vName, val, ok := rdsl.evalElement(elem, pDirs[i])
		if !ok {
			return nil, false
		}
		// Check if we need to store the variable
		if vName != "" {
			m[vName] = val
		}
	}
	return m, ok
}

// Resolve takes a map of varnames and values and replaces any
// occurrences found in src string resulting to a new string..
func (rdsl *RouteDSL) Resolve(m map[string]string, src string) string {
	res := src[0:]
	for k, v := range m {
		k = StartVar + k + EndVar
		if strings.Contains(res, k) {
			res = strings.ReplaceAll(res, k, v)
		}
	}
	return res
}
