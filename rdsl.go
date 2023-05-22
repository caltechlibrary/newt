package newt

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
)


// PathDSLType defines the interface for types bound to
// a name and Go struct.
type PathDSLType interface {
	// EvalType takes an expression string and value string, 
	// checks the expression against the value string validating
	// based on the PathDSLType defined.
	// If the value string is accept a normalized value string and true 
	// are returned. If they value string does not match expr or fails
	// type verification then an empty value string and false is return.
	//
	// EvalType works like a test and set.
	EvalType(string,string) (string, bool)
}


// RouteDSLExpression holds the attributes need to decode
// a RouteDSL expression, match and decode against path values.
type RouteDSLExpression struct {
	Src  string   `json:"src"`
	Dirs []string `json:"dirs,omitempty"`
	Base string   `json:"base,omitempty"`
	Ext  string   `json:"ext,omitempty"`
	// VarToType maps the variable name to a var defn
	VarToType map[string]string `json:"var_to_types,omitempty"`
	// Types maps a type name to type implementation
	Types map[string]PathDSLType `json:"-"`
}

func (rdsl *RouteDSLExpression) String() string {
	src, _ := json.MarshalIndent(rdsl, "", "    ")
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

// NewRouteDSL takes a RouteDSL expression and returns a
// RouteDSLExpresion structure and error value.
func NewRouteDSL(src string) (*RouteDSLExpression, error) {
	rdsl := new(RouteDSLExpression)
	rdsl.Src = src
	dir, base := path.Split(src)
	dirs := strings.Split(strings.TrimSuffix(strings.TrimPrefix(dir, "/"), "/"), "/")
	rdsl.Dirs = []string{}
	// We only evalaute the extension if here are two variables defined for the last element of path.
	if strings.Count(base, "{") == 2 {
		parts := strings.SplitN(base, "}", 2)
		rdsl.Base = parts[0] + "}"
		rdsl.Ext = parts[1]
	} else {
		rdsl.Base = base
		rdsl.Ext = ""
	}
	rdsl.VarToTypes = map[string]string{}

	rdsl.TypeFn = map[string]TypeEval{}
	for i, elem := range dirs {
		if strings.HasPrefix(elem, "{") && strings.HasSuffix(elem, "}") {
			varName, typeExpr, err := varDefn(elem)
			if err == nil {
				rdsl.VarToTypes[varName] = typeExpr
			} else {
				return nil, fmt.Errorf("(%d) %q -> %s", i, elem, err)
			}
			rdsl.Dirs = append(rdsl.Dirs, fmt.Sprintf("{%s}", varName))
		} else {
			rdsl.Dirs = append(rdsl.Dirs, elem)
		}
	}
	if strings.HasPrefix(rdsl.Base, "{") && strings.HasSuffix(rdsl.Base, "}") {
		varName, typeExpr, err := varDefn(rdsl.Base)
		if err == nil {
			rdsl.VarToTypes[varName] = typeExpr
		} else {
			return nil, fmt.Errorf("(basename) %q -> %s", rdsl.Base, err)
		}
		rdsl.Base = fmt.Sprintf("{%s}", varName)
	}
	if rdsl.Ext != "" {
		if strings.HasPrefix(rdsl.Ext, "{") && strings.HasSuffix(rdsl.Ext, "}") {
			varName, typeExpr, err := varDefn(rdsl.Ext)
			if err == nil {
				rdsl.VarToTypes[varName] = typeExpr
			} else {
				return nil, fmt.Errorf("(extname) %q -> %s", rdsl.Ext, err)
			}
			rdsl.Ext = fmt.Sprintf("{%s}", varName)
		}
	}
	return rdsl, nil
}

// RegisterType maps a type name to a function. The function needs to
// be of the form of TypeEval.
func (rdsl *RouteDSLExpression) RegisterType(tName string, fn TypeEval) error {
	if _, ok := rdsl.TypeFn[tName]; ok {
		return fmt.Errorf("%q previously registered", tName)
	}
	rdsl.TypeFn[tName] = fn
	return nil
}

func varName(src string) string {
	return strings.TrimSuffix(strings.TrimPrefix(src, "{"), "}")
}

// evalElement takes compares a element against a value (from the path value)
// returns a variable name, interface value and bool indicating a successful match
func (rdsl *RouteDSLExpression) evalElement(elem string, src string) (string, interface{}, bool) {
	// Check if we workingwith a literal element or a variable defn.
	if strings.HasPrefix(elem, `{`) {
		// handle variable path element
		vName := varName(elem)
		tExpr, ok := rdsl.VarToTypes[vName]
		if !ok {
			return "", nil, false
		}
		fn, ok := rdsl.TypeFn[tExpr]
		if !ok {
			return vName, nil, false
		}
		// Now check the type of dDir against the type expression
		val, ok := defn.EvalType(tExpr, src)
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
func (rdsl *RouteDSLExpression) Eval(pathValue string) (map[string]interface{}, bool) {
	dir, base := path.Split(pathValue)
	pDirs := strings.Split(strings.TrimSuffix(strings.TrimPrefix(dir, "/"), "/"), "/")
	pExt := path.Ext(base)
	pBase := strings.TrimSuffix(base, pExt)
	if rdsl.Ext == "" {
		pExt = ""
		pBase = base
	}
	if len(pDirs) != len(rdsl.Dirs) {
		return nil, false
	}
	m := map[string]interface{}{}
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
	// Match Basename
	if vName, val, ok := rdsl.evalElement(rdsl.Base, pBase); ok {
		// Check if we need to store the variable
		if vName != "" {
			m[vName] = val
		}
	} else {
		return nil, false
	}
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
	return m, true
}
