package newt

import (
	"fmt"
	"io"
)

// NewtGenerator holds our Newt Generator structures for rendering code.
type NewtGenerator struct {
	// Namespace is used when generating the SQL/conf for setting up Postgres+PostgREST
	Namespace string

	// Models holds the models used to generator specific code
	Models []*NewtModel

	// Options holds the result environment variables and options that can be used in generator code
	Options map[string]string

	// internal this is the output for code generation, usually resolves to stdout
	out io.Writer
	// internal this is the error output for code generation, usually resolves to stderr
	eout io.Writer

	// Postgres configuration information
	Postgres *Application

	// PostgREST configuration information
	PostgREST *Application
}

// NewGenerator instaitates a new Generator object form a filename and Config object
// It returns a Generator object and error value.
func NewGenerator(cfg *Config) (*NewtGenerator, error) {
	if cfg.Applications == nil || cfg.Applications.NewtGenerator == nil {
		return nil, fmt.Errorf("configuration missing for Newt Generator")
	}
	generator := &NewtGenerator{}
	generator.Namespace = cfg.Applications.NewtGenerator.Namespace
	generator.Models = cfg.Models
	generator.Options = map[string]string{}
	if cfg.Applications.Postgres != nil {
		generator.Postgres = cfg.Applications.Postgres
	}
	if cfg.Applications.PostgREST != nil {
		generator.PostgREST = cfg.Applications.PostgREST
	}
	// NOTE: LoadCondfig handles loading the environment into options. We just need to
	// copy into the NewtGenerator struct.
	if len(cfg.Applications.Options) > 0 {
		for k, v := range cfg.Applications.Options {
			generator.Options[k] = v
		}
	}
	return generator, nil
}

// renderPostgreSQL does what its name implies. It outputs an SQL
// program in Postgres SQL dialect. It does so for a specific type
// request. Possible values are "setup", "models" and "models_test".
//
// The "setup" code type includes a placeholder for your DB credentials.
// It should not be included in your GitHub repository.
//
// The "models" contains all the table definitions, view definitions,
// and functions implementing CRUD-L operations for each model.
//
// The "models_test" contains SQL to test your models and ensure they
// were created successfully.
func (g *NewtGenerator) renderPostgreSQL(action string) error {
	switch action {
	case "setup":
		return pgSetup(g.out, g.Namespace)
	case "models":
		return pgModels(g.out, g.Namespace, g.Models)
	case "models_test":
		return pgModelsTest(g.out, g.Namespace, g.Models)
	default:
		return fmt.Errorf("%q not supported at this time", action)
	}
}

// renderPostgREST does what its name implies. It the configuration
// file used when starting up PostgREST.
func (g *NewtGenerator) renderPostgREST() error {
	port := "5432"
	if g.Postgres != nil && g.Postgres.Port != 0 {
		port = fmt.Sprintf("%d", g.Postgres.Port)
	}
	return prConfig(g.out, g.Namespace, port)
}


// renderMustache will render a mustache template for a given model id. The action corresponds
// to the model id.
func (g *NewtGenerator) renderMustache(action string, modelId string) error {
	for _, model := range g.Models {
		if modelId == model.Id {
			return MTmplGen(g.out, action, model)
		}
	}
	return fmt.Errorf("failed to find model id %q", modelId)
}

// renderHtml will render HTML forms for given action and model id.
func (g *NewtGenerator) renderHtml(action string, modelId string) error {
	return fmt.Errorf("g.renderHtml(%q, %q) not implemented", action, modelId)
}

// validate action from list of actions.
func validateAction(action string, supportedActions []string) error {
	if action == "" {
		return fmt.Errorf("missing action")
	}
	for _, supportedAction := range supportedActions {
		if action == supportedAction {
			return nil
		}
	}
	return fmt.Errorf("%q is not a supported action", action)
}

// validateModelId
func validateModelId(modelId string, models []*NewtModel) error {
	for _, model := range models {
		if modelId == model.Id {
			return nil
		}
	}
	return fmt.Errorf("%q is not a valid model id", modelId)
}

// Generator generates the code based on the contents of Generator struct. It will also verify that the
// needed parameters are provided.
//
// - generatorName is the generator to use
// - action is a parameter that the selected generator can use (e.g. PostgreSQL has setup as well as )
// - modelId references the `.id` attribute of the model needing code generation
func (g *NewtGenerator) Generate(generatorName string, action string, modelId string) error {
	pgActions := []string{ "setup", "models", "models_test" }
	//modelActions := []string{ "create", "read", "update", "delete", "list", "page" }
	templateActions := []string{ 
		"create_form", "create_response", 
		"update_form", "update_response",
		"delete_form", "delete_response",
		"read", "list",
	}
	switch generatorName {
	case "postgres":
		if err := validateAction(action, pgActions); err != nil {
			return err
		}
		return g.renderPostgreSQL(action)
	case "postgrest":
		return g.renderPostgREST()
	case "mustache":
		if err := validateAction(action, templateActions); err != nil  {
			return err
		}
		if err := validateModelId(modelId, g.Models);  err != nil {
			return err
		}
		return g.renderMustache(action, modelId)
	default:
		return fmt.Errorf("%q is not supported at this time", generatorName)
	}
}
