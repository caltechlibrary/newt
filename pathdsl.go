package newt

import (
	"fmt"
	"strings"
	"path"
)

// TypeFunc is an function that takes a type expression (everything
// that would be int he curly braces) and a value.
// It returns the extracted value and bool indicating
// is the extraction was succesful.
type TypeFunc func(string, string) (interface{}, bool)

// PathDSLExpression holds the attributes need to decode
// a PathDSL expression, match and decode against path values.
type PathDSLExpression struct {
	Src string `json:"src"`
	Dirs []string `json:"dirs,omitempty"`
	Base string `json:"base,omitempty"`
	Ext string `json:"ext,omitempty"`
	Types map[string]string `json:"types,omitempty"`
	TypeFn map[string]TypeFunc `json:"-"`
}

// varDefn evaluates a varaible expression returning a var name,
// type expression.
func varDefn(src string) (string, string, error) {
	expr := strings.TrimSuffix(strings.TrimPrefix(src, "{"), "}")
	if strings.Compare(src, s) == 0 {
		return "", "", fmt.Errorf("missing opening or closing curly braces")
	}
	parts := strings.SplitN(expr, 2)
	if parts[0] == "" {
		return "", "", fmt.Eprintf("missing variable name")
	}
	if parts[1] == "" {
		return "", "", fmt.Errorf("missing type expression")
	}
	return parts[0], s, nil
}

// NewPathDSL takes a PathDSL expression and returns a
// PathDSLExpresion structure and error value.
func NewPathDSL(src string) (*PathDSLExpression, error) {
	expr := new(PathDSLExpression)
	expr.Src = src
	dir, base := path.Split(src)
	dirs = strings.Split(strings.TrimPrefix(dir, "/"), "/")
	expr.Dirs = []string{}
	expr.Ext = path.Ext(base)
	expr.Base = strings.TrimSuffix(base, expr.Ext)
	expr.Types = map[string]string{}
	expr.TypeFn = map[string]TypeFunc{}
	for i, elem := range dirs {
		if strings.HasPrefix(elem, "{") && strings.HasSuffix(elem, "}") {
			if varName, typeExpr, err := varDefn(elem); err == nil {
				expr.Types[varName] = typeExpr
			} else {
				return nil, fmt.Errorf("(%d) %q -> %s", i, elem, error)
			}
			expr.Dirs = append(expr.Dirs, fmt.Sprintf("{%s}", varName))
		} else {
			expr.Dirs = append(expr.Dirs, elem)
		}
	}
	if strings.HasPrefix(expr.Base, "{") && strings.HasSuffix(expr.Base, "}") {
		if varName, typeExpr, err := varDefn(expr.Base); err == nil {
			expr.Types[varName] = typeExpr
		} else {
			return nil, fmt.Errorf("(%d) %q -> %s", i, expr.Base, error)
		}		
	}
	if strings.HasPrefix(expr.Ext, "{") && strings.HasSuffix(expr.Ext, "}") {
		if varName, typeExpr, err := varDefn(expr.Ext); err == nil {
			expr.Types[varName] = typeExpr
		} else {
			return nil, fmt.Errorf("(%d) %q -> %s", i, expr.Ext, err)
		}		
	}
	return nil, fmt.Errorf("NewPathDSL() not implemented")
}

// RegisterType maps a type name to a function. The function needs to
// be of the form of TypeFunc.
func (expr *PathDSLExpression) RegisterType(name string, fn TypeFunc) error {
	if _, ok := expr.TypeFn[name]; ok {
		return fmt.Errorf("%q previously registered", name)
	}
	expr.TypeFn[name] = fn
	return nil
}

// Eval takes a path value and compares it with a Path expression.
// It returns a status boolean, map of variable names to values.
func (expr *PathDSLExpression) Eval(pathValue string) (map[string]interface{}, bool) {
	dir, base := path.Split(pathValue)
	pDirs := strings.Split(strings.TrimPrefix(dir, "/"))
	pExt := path.Ext(base)
	pBase := strings.TrimSuffix(base, pExt)
	if len(pDirs) != len(expr.Dirs) {s
		return nil, false
	}
	m := map[string]interface{}{}
	for i, elem := range expr.Dirs {
		name := varName(elem)
		if strings.HasSuffix(elem, '{') {
			fn, ok := expr.TypeFunc[name]
			if ! ok {
				return nil, false
			}
			if expr, ok := expr.Types[name]; ok {
				if val, ok := fn(expr, pDirs[i]); ! ok {
					return nil, false
				} else {
					m[name] = val
				}
			}
		} else if strings.Compare(pDirs[i], elem) != 0 {
			return nil, false
		}
	}
	return m, true
}
