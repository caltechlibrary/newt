package newt

import (
	"fmt"
	"io"
)

// NewtGenerator holds our Newt Generator structures for rendering code.
// FIXME: Not implemented yet.
type NewtGenerator struct {
	// Namespace is used when generating the SQL/conf for setting up Postgres+PostgREST
	Namespace string

	// Models holds the models used to generator code
	Models []*NewtModel

	// Options holds the result environment variables and options that can be usedd in generator code
	Options map[string]string

	// internal this is the output for code generation, usually resolves to stdout
	out io.Writer
	// internal this is the error output for code generation, usually resolves to stderr
	eout io.Writer
}

// NewGenerator instaitates a new Generator object form a filename and Config object
// It returns a Generator object and error value.
func NewGenerator(cfg *Config) (*NewtGenerator, error) {
	generator := &NewtGenerator{}
	if cfg.Applications.NewtGenerator != nil {
		generator.Namespace = cfg.Applications.NewtGenerator.Namespace
	}
	generator.Models = cfg.Models
	generator.Options = map[string]string{}
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
func (g *NewtGenerator) renderPostgreSQL(codeType string) error {
	switch codeType {
		case "setup":
			return pgSetup(g.out, g.Namespace)
		case "models":
			return pgModels(g.out, g.Namespace, g.Models)
		case "models_test":
			return pgModelsTest(g.out, g.Models)
		default:
			return fmt.Errorf("%q not supported at this time", codeType)
	}
	return fmt.Errorf("g.renderPostgreSQL(%q) not implemented", codeType)
}

func (g *NewtGenerator) renderPostgREST(codeType string) error {
	return fmt.Errorf("g.renderPostgREST(%q) not implemented", codeType)
}

func (g *NewtGenerator) renderMustache(codeType string) error {
	return fmt.Errorf("g.renderMustache(%q) not implemented", codeType)
}

func (g *NewtGenerator) renderHtml(codeType string) error {
	return fmt.Errorf("g.renderHtml(%q) not implemented", codeType)
}

// Generator generates the code based on the contents of Generator struct.
// FIXME: Not implemented yet.
func (g *NewtGenerator) Generate(target string, codeType string) error {
	switch target {
		case "postgres":
			return g.renderPostgreSQL(codeType)
		case "postgrest":
			return g.renderPostgREST(codeType)
		case "mustache":
			return g.renderMustache(codeType)
		case "html":
			return g.renderHtml(codeType)
		default:
			return fmt.Errorf("%q is not supported at this time", target)
	}
	return nil
}
