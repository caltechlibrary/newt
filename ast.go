/**
 * ast.go holds the data structure that defines Newt applications.
 *
 * @author R. S. Doiel
 */
package newt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	// 3rd Party
	"gopkg.in/yaml.v3"
)

// AST holds a configuration for Newt for the data router and code generator.
type AST struct {
	// Application holds the application specific settings and metadata
	Applications *Applications `json:"applications,omitempty" yaml:"applications,omitempty"`

	// Models holds a list of data models. It is used by
	// both the data router and code generator.
	Models []*Model `json:"models,omitempty" yaml:"models,omitempty"`

	// Routes holds an array of maps of route definitions used by
	// the data router and code generator
	Routes []*Route `json:"routes,omitempty" yaml:"routes,omitempty"`

	// Templates holds an array of mapsthe request to template to request for
	// Newt Mustache
	Templates []*MustacheTemplate `json:"templates,omitempty" yaml:"templates,omitempty"`

	// isChanged is a convience variable for tracking if the data structure has changed.
	isChanged bool `json:"-" yaml:"-"`
}

// Applications holds the runtime information for newt router, generator,
// template engine.
type Applications struct {
	// Newt Router runtime config
	Router *Application `json:"newtrouter,omitempty" yaml:"newtrouter,omitempty"`

	// Newt Mustache runtime config
	NewtMustache *Application `json:"newtmustache,omitempty" yaml:"newtmustache,omitempty"`

	// Newt Generator runtime config
	NewtGenerator *Application `json:"newtgenerator,omitempty" yaml:"newtgenerator,omitempty"`

	// Postgres runtime config, e.g. port number to use for connecting.
	Postgres *Application `json:"postgres,omitempty" yaml:"postgres,omitempty"`

	// PostgREST runtime config
	PostgREST *Application `json:"postgrest,omitempty" yaml:"postgrest,omitempty"`

	// Environment holds a list of OS environment variables that can be made
	// available to the web services.
	Environment []string `json:"environment,omitempty" yaml:"enviroment,omitempty"`

	// Options is a map of name to string values, it is where the
	// environment variables values are stored.
	Options map[string]string `json:"options,omitempty" yaml:"options,omitempty"`
}

// Application implements runtime config for Newt programs
type Application struct {
	// AppPath holds the path to the binary application, e.g. PostgREST
	// This property provides the location of the service to run.
	AppPath string `json:"app_path,omitempty" yaml:"app_path,omitempty"`

	// ConfPath holds teh path to the configuration file (e.g. PostgREST configuration file)
	ConfPath string `json:"conf_path,omitempty" yaml:"conf_path,omitempty"`

	// Namespace holds the Postgres Schema name It is used to generate
	// a setup.sql file using the -pg-setup option in newt cli.
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`

	// Port is the name of the localhost port Newt will listen on.
	Port int `json:"port,omitempty" yaml:"port,omitempty"`

	// Timeout is a duration, it is used to set timeouts and the application.
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`

	// Htdocs holds any static files you want to make available through
	// Newt router.
	Htdocs string `json:"htdocs,omitempty" yaml:"htdocs,omitempty"`

	// DSN, data ast name is a URI connection string
	DSN string `json:"dsn,omitemity" yaml:"dsn,omitempty"`
}

// UnmarshalAST will read []byte of YAML or JSON,
// populate the provided *AST object and return an error.
//
// ```
// src, _ := os.ReadFile("app.yaml")
// ast := new(AST)
//
//	if err := UnmarshalAST(src, ast); err != nil {
//	    // ... handle error
//	}
//
// ```
func UnmarshalAST(src []byte, ast *AST) error {
	if bytes.HasPrefix(src, []byte("{")) {
		if err := json.Unmarshal(src, &ast); err != nil {
			return err
		}
	} else {
		if err := yaml.Unmarshal(src, &ast); err != nil {
			return err
		}
	}
	if ast.Applications == nil {
		ast.Applications = &Applications{
			Router:        &Application{},
			NewtMustache:  &Application{},
			NewtGenerator: &Application{},
			Postgres:      &Application{},
			PostgREST:     &Application{},
			Options:       map[string]string{},
			Environment:   []string{},
		}
	}
	return nil
}

// LoadAST read a configuration file, merges environment variables
// and returns a AST object and error value.
//
// ```
// ast, err := LoadAST("app.yaml")
//
//	if err != nil {
//	    // ... handle error
//	}
//
// ```
func LoadAST(configFName string) (*AST, error) {
	ast := new(AST)
	if configFName != "" {
		src, err := os.ReadFile(configFName)
		if err != nil {
			return nil, fmt.Errorf("failed to read %q, %s", configFName, err)
		}
		if err := UnmarshalAST(src, ast); err != nil {
			return nil, fmt.Errorf("failed to read %q, %s", configFName, err)
		}
	}

	if ast.Applications == nil {
		ast.Applications = &Applications{
			Router:        &Application{},
			NewtMustache:  &Application{},
			NewtGenerator: &Application{},
			Options:       map[string]string{},
			Environment:   []string{},
		}
	}
	// Load environment if missing from config file.
	if len(ast.Applications.Environment) == 0 {
		for _, envar := range ast.Applications.Environment {
			// YAML settings take presidence over environment, check for conflicts
			if _, conflict := ast.Applications.Options[envar]; !conflict {
				ast.Applications.Options[envar] = os.Getenv(envar)
			}
		}
	}
	ast.isChanged = false
	return ast, nil
}

func (ast *AST) HasChanges() bool {
	if ast.isChanged {
		return true
	}
	for _, m := range ast.Models {
		if m.HasChanges() {
			return true
		}
	}
	return false
}

func (ast *AST) Encode() ([]byte, error) {
	// Now output the YAML
	timeStamp := (time.Now()).Format("2006-01-02")
	userName := os.Getenv("USER")
	comment := []byte(fmt.Sprintf(`#/usr/bin/env newt check
#
# This was generated by %s on %s with %s version %s %s.
#
`, userName, timeStamp, path.Base(os.Args[0]), Version, ReleaseHash))
	data := bytes.NewBuffer(comment)
	encoder := yaml.NewEncoder(data)
	encoder.SetIndent(2)
	if err := encoder.Encode(ast); err != nil {
		return nil, fmt.Errorf("failed to generate configuration, %s\n", err)
	}
	return data.Bytes(), nil
}

// SaveAs writes the *AST to a YAML file.
func (ast *AST) SaveAs(configName string) error {
	if _, err := os.Stat(configName); err == nil {
		if err := backupFile(configName); err != nil {
			return err
		}
	}
	fp, err := os.Create(configName)
	if err != nil {
		return err
	}
	defer fp.Close()
	src, err := ast.Encode()
	if err != nil {
		return err
	}
	fmt.Fprintf(fp, "%s", src)
	for _, model := range ast.Models {
		for _, element := range model.Elements {
			element.isChanged = false
		}
		model.isChanged = false
	}
	ast.isChanged = false
	return nil
}

// GetModelIds returns a list of model ids
func (ast *AST) GetModelIds() []string {
	if ast.Models == nil {
		ast.Models = []*Model{}
	}
	ids := []string{}
	for _, m := range ast.Models {
		if m.Id != "" {
			ids = append(ids, m.Id)
		}
	}
	return ids
}

// GetModelById return a specific model by it's id
func (ast *AST) GetModelById(id string) (*Model, bool) {
	for _, m := range ast.Models {
		if m.Id == id {
			return m, true
		}
	}
	return nil, false
}

// AddModel takes a new Model, checks if the model exists in the list (i.e.
// has an existing model id that matches the new model and if not appends
// it so the list.
func (ast *AST) AddModel(model *Model) error {
	// Make sure we have a Models lists to work with.
	if ast.Models == nil {
		ast.Models = []*Model{}
	}
	// Check to see if this is a duplicate, return error if it is
	for i, m := range ast.Models {
		if m.Id == model.Id {
			return fmt.Errorf("failed, model %d is a duplicate model id, %q", i, m.Id)
		}
	}
	ast.Models = append(ast.Models, model)
	ast.isChanged = true
	return nil
}

// UpdateModel takes a model id and new model struct replacing the
// existing one.
func (ast *AST) UpdateModel(id string, model *Model) error {
	// Make sure we have a Models lists to work with.
	if ast.Models == nil {
		return fmt.Errorf("no models defined")
	}
	for i, m := range ast.Models {
		if m.Id == id {
			ast.Models[i] = model
			ast.isChanged = true
			return nil
		}
	}
	return fmt.Errorf("failed to find model %q", id)
}

// RemoveModelById find the model with the model id and remove it
func (ast *AST) RemoveModelById(id string) error {
	// Make sure we have a Models lists to work with.
	if ast.Models == nil {
		return fmt.Errorf("no models defined")
	}
	for i, m := range ast.Models {
		if m.Id == id {
			ast.Models = append(ast.Models[:i], ast.Models[(i+1):]...)
			ast.isChanged = true
			return nil
		}
	}
	return fmt.Errorf("failed to find model %q", id)
}

// RemoveRouteById find the route with route id and remove it
func (ast *AST) RemoveRouteById(id string) error {
	routeFound := false
	for i, r := range ast.Routes {
		// NOTE: A route id ties one or more requests together, e.g. retrieve a web form (GET), then handle it (POST)
		if r.Id == id {
			ast.Routes = append(ast.Routes[:i], ast.Routes[(i+1):]...)
			ast.isChanged = true
			routeFound = true
		}
	}
	if ! routeFound {
		return fmt.Errorf("failed to find route %s", id)
	}
	return nil
}

// RemoveTemplateById() find the template id and remove it from the .Templates structure
func (ast *AST) RemoveTemplateById(id string) error {
	templateFound := false
	for i, t := range ast.Templates {
		if t.Id == id {
			ast.Templates = append(ast.Templates[:i], ast.Templates[(i+1):]...)
			ast.isChanged = true
			templateFound = true
		}
	}
	if ! templateFound {
		return fmt.Errorf("failed to find template %s", id)
	}
	return nil
}

// GetRouteIds returns a list of Router ids found in ast.Routes
func (ast *AST) GetRouteIds() []string {
	rIds := []string{}
	for _, r := range ast.Routes {
		if r.Id != "" {
			rIds = append(rIds, r.Id)
		}
	}
	return rIds
}

// GetTemplateIds return a list of template ids.
func (ast *AST) GetTemplateIds() []string {
	tIds := []string{}
	for _, t := range ast.Templates {
		if t.Id != "" {
			tIds = append(tIds, t.Id)
		}
	}
	return tIds
}

// GetPrimaryTemplates return a list of primary template filenames
func (ast *AST) GetPrimaryTemplates() []string {
	fNames := []string{}
	for _, t := range ast.Templates {
		if t.Template != "" {
			fNames = append(fNames, t.Template)
		}
	}
	return fNames
}

// GetAllTemplates returns a list of templates, including partials defined
// in the .Templates property. Part template names are indented with a "\t"
func (ast *AST) GetAllTemplates() []string {
	fNames := []string{}
	for _, t := range ast.Templates {
		if t.Template != "" {
			fNames = append(fNames, t.Template)
			for _, p := range t.Partials {
				if p != "" {
					fNames = append(fNames, fmt.Sprintf("\t%s", p))
				}
			}
		}
	}
	return fNames
}

// GetTemplateByPrimary returns the template entry using primary template filename
func (ast *AST) GetTemplateByPrimary(fName string) (*MustacheTemplate, bool) {
	if ast.Templates != nil {
		for _, t := range ast.Templates {
			if t.Template == fName {
				return t, true
			}
		}
	}
	return nil, false
}

// Check reviews the ast *AST and reports and issues, return true
// if no errors found and false otherwise.  The "buf" will hold the error output.
func (ast *AST) Check(buf io.Writer) bool {
	ok := true
	if ast.Applications == nil {
		fmt.Fprintf(buf, "no applications defined\n")
		ok = false
	}
	if ast.Models == nil || len(ast.Models) == 0 {
		fmt.Fprintf(buf, "no models defined\n")
		ok = false
	} else {
		for i, m := range ast.Models {
			if !m.Check(buf) {
				fmt.Fprintf(buf, "model #%d is invalid\n", i)
				ok = false
			}
		}
	}
	if ast.Routes == nil || len(ast.Routes) == 0 {
		if ast.Applications.Router != nil {
			fmt.Fprintf(buf, "no routes defined for Newt Router\n")
			ok = false
		}
	} else {
		if ast.Applications.Router == nil {
			fmt.Fprintf(buf, "application.router not configured")
			ok = false
		} else if ast.Applications.Router.Port == 0 {
			fmt.Fprintf(buf, "application.router.port not set")
			ok = false
		}
		for i, r := range ast.Routes {
			if !r.Check(buf) {
				fmt.Fprintf(buf, "route (#%d) errors\n", i)
				ok = false
			}
		}
	}
	if ast.Templates == nil || len(ast.Templates) == 0 {
		if ast.Applications.NewtMustache != nil {
			fmt.Fprintf(buf, "no templates defined but Newt Mustache enabled\n")
			ok = false
		}
	} else {
		if ast.Applications.NewtMustache == nil {
			fmt.Fprintf(buf, "application.newt_mustache not configured")
			ok = false
		} else if ast.Applications.NewtMustache.Port  == 0 {
			fmt.Fprintf(buf, "application.newt_mustache.port not set")
			ok = false
		}
		for i, t := range ast.Templates {
			if !t.Check(buf) {
				fmt.Fprintf(buf, "template (#%d) error\n", i)
				ok = false
			}
		}
	}
	return ok
}

