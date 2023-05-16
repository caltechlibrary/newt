package newt

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
)

// TypeEval is an function that takes a type expression (everything
// that would be int he curly braces) and a value.
// It returns the extracted value and bool indicating
// is the extraction was succesful.
type TypeEval func(string, string) (interface{}, bool)

// PathDSLExpression holds the attributes need to decode
// a PathDSL expression, match and decode against path values.
type PathDSLExpression struct {
	Src  string   `json:"src"`
	Dirs []string `json:"dirs,omitempty"`
	Base string   `json:"base,omitempty"`
	Ext  string   `json:"ext,omitempty"`
	// VarToTypes maps the variable name to a var defn
	VarToTypes map[string]string `json:"var_to_types,omitempty"`
	// TypeFn maps a type to a tpe eval func
	TypeFn map[string]TypeEval `json:"-"`
}

func (pdsl *PathDSLExpression) String() string {
	src, _ := json.MarshalIndent(pdsl, "", "    ")
	return string(src)
}

// varDefn evaluates a varaible expression returning a var name,
// type expression.
func varDefn(src string) (string, string, error) {
	expr := strings.TrimSuffix(strings.TrimPrefix(src, "{"), "}")
	if strings.Compare(src, expr) == 0 {
		return "", "", fmt.Errorf("missing opening or closing curly braces")
	}
	parts := strings.SplitN(expr, " ", 2)
	vName, tExpr := parts[0], parts[1]
	if vName == "" {
		return "", "", fmt.Errorf("missing variable name")
	}
	if tExpr == "" {
		return vName, "", fmt.Errorf("missing type expression for var %q", vName)
	}
	return vName, tExpr, nil
}

// NewPathDSL takes a PathDSL expression and returns a
// PathDSLExpresion structure and error value.
func NewPathDSL(src string) (*PathDSLExpression, error) {
	pdsl := new(PathDSLExpression)
	pdsl.Src = src
	dir, base := path.Split(src)
	dirs := strings.Split(strings.TrimSuffix(strings.TrimPrefix(dir, "/"), "/"), "/")
	pdsl.Dirs = []string{}
	// We only evalaute the extension if here are two variables defined for the last element of path.
	if strings.Count(base, "{") == 2 {
		parts := strings.SplitN(base, "}", 2)
		pdsl.Base = parts[0] + "}"
		pdsl.Ext = parts[1]
	} else {
		pdsl.Base = base
		pdsl.Ext = ""
	}
	pdsl.VarToTypes = map[string]string{}

	pdsl.TypeFn = map[string]TypeEval{}
	for i, elem := range dirs {
		if strings.HasPrefix(elem, "{") && strings.HasSuffix(elem, "}") {
			varName, typeExpr, err := varDefn(elem)
			if err == nil {
				pdsl.VarToTypes[varName] = typeExpr
			} else {
				return nil, fmt.Errorf("(%d) %q -> %s", i, elem, err)
			}
			pdsl.Dirs = append(pdsl.Dirs, fmt.Sprintf("{%s}", varName))
		} else {
			pdsl.Dirs = append(pdsl.Dirs, elem)
		}
	}
	if strings.HasPrefix(pdsl.Base, "{") && strings.HasSuffix(pdsl.Base, "}") {
		varName, typeExpr, err := varDefn(pdsl.Base)
		if err == nil {
			pdsl.VarToTypes[varName] = typeExpr
		} else {
			return nil, fmt.Errorf("(basename) %q -> %s", pdsl.Base, err)
		}
		pdsl.Base = fmt.Sprintf("{%s}", varName)
	}
	if pdsl.Ext != "" {
		if strings.HasPrefix(pdsl.Ext, "{") && strings.HasSuffix(pdsl.Ext, "}") {
			varName, typeExpr, err := varDefn(pdsl.Ext)
			if err == nil {
				pdsl.VarToTypes[varName] = typeExpr
			} else {
				return nil, fmt.Errorf("(extname) %q -> %s", pdsl.Ext, err)
			}
			pdsl.Ext = fmt.Sprintf("{%s}", varName)
		}
	}
	return pdsl, nil
}

// RegisterType maps a type name to a function. The function needs to
// be of the form of TypeEval.
func (pdsl *PathDSLExpression) RegisterType(tName string, fn TypeEval) error {
	if _, ok := pdsl.TypeFn[tName]; ok {
		return fmt.Errorf("%q previously registered", tName)
	}
	pdsl.TypeFn[tName] = fn
	return nil
}

func varName(src string) string {
	return strings.TrimSuffix(strings.TrimPrefix(src, "{"), "}")
}

// evalElement takes compares a element against a value (from the path value)
// returns a variable name, interface value and bool indicating a successful match
func (pdsl *PathDSLExpression) evalElement(elem string, src string) (string, interface{}, bool) {
	// Check if we workingwith a literal element or a variable defn.
	if strings.HasPrefix(elem, `{`) {
		// handle variable path element
		vName := varName(elem)
		tExpr, ok := pdsl.VarToTypes[vName]
		if !ok {
			return "", nil, false
		}
		fn, ok := pdsl.TypeFn[tExpr]
		if !ok {
			return vName, nil, false
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
	if strings.Compare(elem, src) != 0 {
		return "", "", false
	}
	return "", "", true
}

// Eval takes a path value and compares it with a Path expression.
// It returns a status boolean, map of variable names to values.
func (pdsl *PathDSLExpression) Eval(pathValue string) (map[string]interface{}, bool) {
	dir, base := path.Split(pathValue)
	pDirs := strings.Split(strings.TrimSuffix(strings.TrimPrefix(dir, "/"), "/"), "/")
	pExt := path.Ext(base)
	pBase := strings.TrimSuffix(base, pExt)
	if pdsl.Ext == "" {
		pExt = ""
		pBase = base
	}
	if len(pDirs) != len(pdsl.Dirs) {
		return nil, false
	}
	m := map[string]interface{}{}
	for i, elem := range pdsl.Dirs {
		vName, val, ok := pdsl.evalElement(elem, pDirs[i])
		if !ok {
			return nil, false
		}
		// Check if we need to store the variable
		if vName != "" {
			m[vName] = val
		}
	}
	// Match Basename
	if vName, val, ok := pdsl.evalElement(pdsl.Base, pBase); ok {
		// Check if we need to store the variable
		if vName != "" {
			m[vName] = val
		}
	} else {
		return nil, false
	}
	// Match the extension, if it contains a
	if pdsl.Ext != "" {
		if vName, val, ok := pdsl.evalElement(pdsl.Ext, pExt); ok {
			// Check if we need to store the variable
			if vName != "" {
				m[vName] = val
			}
		} else {
			return nil, false
		}
	}
	return m, true
}
