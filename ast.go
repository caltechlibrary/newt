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
	"strings"
	"time"

	// 3rd Party
	//"github.com/cbroglie/mustache"
	"github.com/aymerick/raymond"
	"gopkg.in/yaml.v3"
)

// AST holds a configuration for Newt for the data router and code generator.
type AST struct {
	// AppMetadata holds you app's metadata such as needed to render an "about" page in your final app.
	AppMetadata *AppMetadata `json:"app_metadata,omitempty" yaml:"app_metadata,omitempty"`

	// Application holds the back end application with their specific settings used to compose your Newt Application.
	//FIXME: This should be a list of "application" objects their settings.  This would include start/stop/restart actions
	// and enough metadata to generated appropriate Systemd and Luanchd configurations.
	Applications *Applications `json:"applications,omitempty" yaml:"applications,omitempty"`

	// Models holds a list of data models. It is used by
	// both the data router and code generator.
	Models []*Model `json:"models,omitempty" yaml:"models,omitempty"`

	// Routes holds an array of maps of route definitions used by
	// the data router and code generator
	Routes []*Route `json:"routes,omitempty" yaml:"routes,omitempty"`

	// Templates holds an array of maps the request to template to request for
	// Newt (Handlebars) template engine
	Templates []*Template `json:"templates,omitempty" yaml:"templates,omitempty"`

	// isChanged is a convience variable for tracking if the data structure has changed.
	isChanged bool `json:"-" yaml:"-"`
}

// AppMetadata holds metadata about your Newt Application
// This is primarily used in generated Handlbars partials
type AppMetadata struct {
	AppName string `json:"name,omitempty" yaml:"app_name,omitempty"`
	AppTitle string `json:"title,omitempty" yaml:"app_title,omitempty"`
	CopyrightYear string `json:"copyright_year,omitempty" yaml:"copyright_year,omitempty"`
	CopyrightLink string `json:"copyright_link,omitempty" yaml:"copyright_link,omitempty"`
	CopyrightText string `json:"copyright_text,omitempty" yaml:"copyright_text,omitempty"`
	LogoLink string `json:"logo_link,omitempty" yaml:"logo_link,omitempty"`
	LogoText string `json:"logo_text,omitempty" yaml:"logo_text,omitempty"`
	LicenseLink string `json:"license_link,omitempty" yaml:"license_link,omitempty"`
	LicenseText string `json:"license_text,omitempty" yaml:"license_text,omitempty"`
	CSSPath string `json:"css_path,omitempty" yaml:"css_path,omitempty"`
	HeaderLink string `json:"header_link,omitempty" yaml:"header_link,omitempty"`
	HeaderText string `json:"header_text,omitempty" yaml:"header_text,omitempty"`
	ContactAddress string `json:"contact_address,omitempty" yaml:"contact_address,omitempty"`
	ContactPhone string `json:"contact_phone,omitempty" yaml:"contact_phone,omitempty"`
	ContactEMail string `json:"contact_email,omitempty" yaml:"contact_email,omitempty"`
}

// Applications holds the runtime information for newt router, generator,
// template engine.
type Applications struct {
	// Newt Router runtime config
	Router *Application `json:"router,omitempty" yaml:"router,omitempty"`

	// TemplateEngine holds Handlebars runtime configuration for Newt template engine
	TemplateEngine *Application `json:"template_engine,omitempty" yaml:"template_engine,omitempty"`

	// Dataset runtime config
	Datasetd *Application `json:"dataset,omitempty" yaml:"dataset,omitempty"`

	// Postgres runtime config, e.g. port number to use for connecting.
	Postgres *Application `json:"postgres,omitempty" yaml:"postgres,omitempty"`

	// PostgREST runtime config
	PostgREST *Application `json:"postgrest,omitempty" yaml:"postgrest,omitempty"`

	// Environment holds a list of OS environment variables that can be made
	// available to the web services.
	Environment []string `json:"environment,omitempty" yaml:"enviroment,omitempty"`

	// Options is a map of name to string values, it is where the
	// the environment variable valuess are stored.
	Options map[string]interface{} `json:"options,omitempty" yaml:"options,omitempty"`
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

	// CName is the name of the dataset collection you wish to use/generate.
	CName string `json:"c_name,omitempty" yaml:"c_name,omitempty"`

	// Port is the name of the localhost port Newt will listen on.
	Port int `json:"port,omitempty" yaml:"port,omitempty"`

	// Timeout is a duration, it is used to set timeouts and the application.
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`

	// Htdocs holds any static files you want to make available through
	// Newt router.
	Htdocs string `json:"htdocs,omitempty" yaml:"htdocs,omitempty"`

	// BaseDir is used by Handlebars, usually holds the "views" directory.
	BaseDir string `json:"base_dir,omitempty" yaml:"base_dir,omitempty"`

	// ExtName is used by Handlebars to set the expected extension (e.g. ".hbs")
	ExtName string `json:"ext_name,omitempty" yaml:"ext_name,omitempty"`

	// PartialsDir is used by Handlebars to find partial templates, usually inside the views directory
	PartialsDir string `json:"partials_dir,omitempty" yaml:"partials_dir,omitempty"`

	// DSN, data ast name is a URI connection string
	DSN string `json:"dsn,omitemity" yaml:"dsn,omitempty"`
}

// NewApplication will create an empty Application struct
func NewApplication() *Application {
	return &Application{}
}

// NewApplications will create an empty Application with top level attributes
func NewApplications() *Applications {
	return &Applications{
		Router:         NewApplication(),
		TemplateEngine: NewApplication(),
		Postgres:       NewApplication(),
		PostgREST:      NewApplication(),
		Datasetd:       NewApplication(),
		Options:        map[string]interface{}{},
		Environment:    []string{},
	}
}

// NewAST will create an empty AST with top level attributes
func NewAST() *AST {
	ast := new(AST)
	ast.Applications = NewApplications()
	return ast
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
		ast.Applications = NewApplications()
	}
	return nil
}

// LoadAST read a YAML file, merges environment variables
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
		ast.Applications = NewApplications()
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
	comment := []byte(fmt.Sprintf(`#!/usr/bin/env newt check
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
	if !routeFound {
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
	if !templateFound {
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
		}
	}
	return fNames
}

// GetTemplateByPrimary returns the template entry using primary template filename
func (ast *AST) GetTemplateByPrimary(fName string) (*Template, bool) {
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

	if ast.Applications.Postgres != nil || ast.Applications.Datasetd != nil {
		if ast.Models == nil || len(ast.Models) == 0 {
			fmt.Fprintf(buf, "no models defined for applications\n")
			ok = false
		} else {
			for i, m := range ast.Models {
				if !m.Check(buf) {
					fmt.Fprintf(buf, "model #%d is invalid\n", i)
					ok = false
				}
			}
		}
	}

	if ast.Applications.Router != nil {
		if ast.Routes == nil || len(ast.Routes) == 0 {
			fmt.Fprintf(buf, "no routes defined for Newt Router\n")
			ok = false
		}
		if ast.Applications.Router.Port == 0 {
			fmt.Fprintf(buf, "application.router.port not set\n")
			ok = false
		}
		for i, r := range ast.Routes {
			if !r.Check(buf) {
				fmt.Fprintf(buf, "route (#%d) errors\n", i)
				ok = false
			}
		}
	}

	if ast.Applications.TemplateEngine != nil {
		if ast.Templates == nil || len(ast.Templates) == 0 {
			fmt.Fprintf(buf, "template engine is defined but not templates are configured\n")
			ok = false
		} else {
			t, err := NewTemplateEngine(ast)
			if err != nil {
				fmt.Fprintf(buf, fmt.Sprintf("application.template_engine not configured, %s\n", err))
				ok = false
			} else if !t.Check(buf) {
				ok = false
			}
		}
	}
	return ok
}

// TemplateEngine defines the `nte` application YAML file. It joins some of the Application struct
// with an array of templates so that "check" can validate the YAML.
type TemplateEngine struct {
	// Port is the name of the localhost port Newt will listen on.
	Port int `json:"port,omitempty" yaml:"port,omitempty"`

	// BaseDir is holds the "views" for that are formed from the templates.
	BaseDir string `json:"base_dir,omitempty" yaml:"base_dir,omitempty"`

	// ExtName is used to set the expected extension (e.g. ".hbs")
	ExtName string `json:"ext_name,omitempty" yaml:"ext_name,omitempty"`

	// PartialsDir is used to find partial templates, usually inside the views directory
	PartialsDir string `json:"partials_dir,omitempty" yaml:"partials_dir,omitempty"`

	// Timeout is a duration, it is used to set timeouts and the application.
	Timeout time.Duration `json:"timeout,omitempty" yaml:"timeout,omitempty"`

	// Templates defined for the service
	Templates []*Template `json:"templates,omitempty" yaml:"templates,omitempty"`
}

// Template hold the request to template mapping for in the TemplateEngine
type Template struct {
	// Id ties a set of one or more template together, e.g. a web form and its response
	Id string `json:"id,required" yaml:"id,omitempty"`

	// Description describes the purpose of the tempalte mapping. It is used to debug Newt YAML files.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Pattern holds a request path, e.g. `/blog_post`. NOTE: the method is ignored. A POST
	// is presumed to hold data that will be processed by the template engine. A GET retrieves the
	// unresolved template.
	Pattern string `json:"request,required" yaml:"request,omitempty"`

	// Template holds a path to the primary template (aka view) file for this route. Path can be relative
	// to the current working directory.
	Template string `json:"template,required" yaml:"template,omitempty"`

	// Debug logs more verbosely if true
	Debug bool `json:"debug,omitempty" yaml:"debug,omitempty"`

	// Document hold the a map of values passed into it from the Newt YAML file in the applications
	// property. These are a way to map in environment or application wide values. These are exposed in
	// the Newt template engine `options`.
	Document map[string]interface{} `json:"document,omitempty" yaml:"document,omitempty"`

	// Vars holds the names of any variables expressed in the pattern, these an be used to replace elements of
	// the output object.
	Vars []string `json:"-" yaml:"-"`

	// Body holds a map of data to process with the template
	Body map[string]interface{} `json:"-" yaml:"-"`

	// The follow are used to simplify individual template invocation.
	// They are populated from the TemplateEngine object.

	/*FIXME: I want to support both Mustache and Handlebars templates.

			 I need to review both mustache and handlebars implementations so figure out an appropriate
			 or wrapper then Tmpl should point to that interface. Gofiber does this with "template views"
			 but I don't want to have to pull in Gofiber's template engine, it's big and provides too many
			 choices to implement smoothly.

	         I need to decide how to specify the template language and if that is per template or engine wide.
			 One approach would be to pick the template engine based on the file extension. Another approoach would
			 be to make it a property of the engine that inherits like the BaseDir, ExtName, etc.
	*/

	// Tmpl points to the compied template
	Tmpl *raymond.Template `json:"-" yaml:"-"`

	// BaseDir is used by holds the "views" directory.
	BaseDir string `json:"-" yaml:"-"`

	// ExtName is used by set the expected extension (e.g. ".hbs")
	ExtName string `json:"-" yaml:"-"`

	// Partials holds partials directory
	PartialsDir string `json:"-" yaml:"-"`
}

// NewTemplateEngine create a new TemplateEngine struct. If a filename
// is provided it reads the file and sets things up accordingly.
func NewTemplateEngine(ast *AST) (*TemplateEngine, error) {
	if ast.Applications.TemplateEngine == nil {
		return nil, fmt.Errorf("template engine is nil")
	}

	// Copy our options so we can expose them in the template's .document
	docvars := map[string]interface{}{}
	// Copy in options to envars
	if ast.Applications.Options != nil && len(ast.Applications.Options) > 0 {
		for k, v := range ast.Applications.Options {
			docvars[k] = v
		}
	}
	te := &TemplateEngine{
		Port:        TEMPLATE_ENGINE_PORT,
		ExtName:     TEMPLATE_ENGINE_EXT_NAME,
		BaseDir:     TEMPLATE_ENGINE_BASE_DIR,
		PartialsDir: TEMPLATE_ENGINE_PARTIALS_DIR,
	}
	if ast.Applications.TemplateEngine.Port != 0 {
		te.Port = ast.Applications.TemplateEngine.Port
	}
	if ast.Applications.TemplateEngine.BaseDir != "" {
		te.BaseDir = ast.Applications.TemplateEngine.BaseDir
	}
	if ast.Applications.TemplateEngine.ExtName != "" {
		te.ExtName = ast.Applications.TemplateEngine.ExtName
	}
	if ast.Applications.TemplateEngine.PartialsDir != "" {
		te.PartialsDir = ast.Applications.TemplateEngine.PartialsDir
	}
	// FIXME: Need to copy in environment variables from ast.Options and set t.Document content.
	errMsgs := []string{}
	if ast.Templates != nil && len(ast.Templates) > 0 {
		// Map in the BaseDir, PartialsDir, and ExtName for the templates.
		for _, t := range ast.Templates {
			t.ExtName = te.ExtName
			t.BaseDir = te.BaseDir
			t.PartialsDir = te.PartialsDir
			if t.Document == nil {
				t.Document = map[string]interface{}{}
			}
			// Copy in options to .document
			if len(docvars) > 0 {
				for k, v := range docvars {
					t.Document[k] = v
				}
			}
		}
		// Add the resulting templates into struct.
		te.Templates = append([]*Template{}, ast.Templates...)
	}
	if len(errMsgs) > 0 {
		return te, fmt.Errorf("%s", strings.Join(errMsgs, "\n"))
	}
	return te, nil
}

// Check makes sure the TemplateEngine struct is populated
func (tEng *TemplateEngine) Check(buf io.Writer) bool {
	if tEng == nil {
		fmt.Fprintf(buf, "template engine not defined\n")
		return false
	}
	errMsgs := []string{}
	ok := true
	if tEng.Port == 0 {
		errMsgs = append(errMsgs, "template engine port not set")
		ok = false
	} else {
		errMsgs = append(errMsgs, fmt.Sprintf("template engine will listen on port %d", tEng.Port))
	}
	if tEng.BaseDir == "" {
		errMsgs = append(errMsgs, "base directory not set for templates")
		ok = false
	}
	if tEng.ExtName == "" {
		errMsgs = append(errMsgs, "template extension is not set")
		ok = false
	}
	if tEng.Templates == nil || len(tEng.Templates) == 0 {
		errMsgs = append(errMsgs, "no templates found")
		ok = false
	} else {
		errMsgs = append(errMsgs, fmt.Sprintf("templates are located in %q", tEng.BaseDir))
		if tEng.PartialsDir != "" {
			errMsgs = append(errMsgs, fmt.Sprintf("partials are located in %q", path.Join(tEng.BaseDir, tEng.PartialsDir)))
		}
		if tEng.ExtName != "" {
			errMsgs = append(errMsgs, fmt.Sprintf("template extension is set to %q", tEng.ExtName))
		}
		//FIXME: add check for helpers
		errMsgs = append(errMsgs, fmt.Sprintf("%d template path(s) mapped", len(tEng.Templates)))
		for i, t := range tEng.Templates {
			tBuf := bytes.NewBuffer([]byte{})
			if !t.Check(tBuf) {
				errMsgs = append(errMsgs, fmt.Sprintf("template (#%d) failed check, %s\n", i, tBuf.Bytes()))
				ok = false
			}
		}
	}
	fmt.Fprintf(buf, "%s\n", strings.Join(errMsgs, "\n"))
	return ok
}

// Check evaluates the *Template and outputs finding. Returns true of no error, false if errors found
func (tmpl *Template) Check(buf io.Writer) bool {
	ok := true
	if tmpl == nil {
		fmt.Fprintf(buf, "template is nil\n")
		return false
	}
	if tmpl.Pattern == "" {
		fmt.Fprintf(buf, "template does not have an associated path/pattern\n")
		ok = false
	}
	if tmpl.Template == "" {
		fmt.Fprintf(buf, "missing path to template for %s\n", tmpl.Pattern)
		ok = false
	} else {
		fmt.Fprintf(buf, "template name %s\n", tmpl.Template)
	}
	return ok
}
