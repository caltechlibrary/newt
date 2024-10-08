package newt

import (
	"fmt"
	"io"

	// Caltech Library Packages
	"github.com/caltechlibrary/models"

)

// Generator holds our Newt Generator structures for rendering code.
type Generator struct {
	// Namespace is used when generating the SQL/conf for setting up Postgres+PostgREST
	Namespace string

	// Models holds the models used to generator specific code
	Models []*models.Model

	// Options holds the result environment variables and options that can be used in generator code
	Options map[string]interface{}

	// internal this is the output for code generation, usually resolves to stdout
	out io.Writer
	// internal this is the error output for code generation, usually resolves to stderr
	eout io.Writer

	// AppMetadata holds the metadata for the application being generated
	AppMetadata *AppMetadata

	// Postgres configuration information
	Postgres *Service

	// PostgREST configuration information
	PostgREST *Service

	// TemplateEngine configuration information
	TemplateEngine *Service
}

// NewGenerator instaitates a new Generator object form a filename and AST object
// It returns a Generator object and error value.
func NewGenerator(ast *AST) (*Generator, error) {
	if ast.Services == nil {
		return nil, fmt.Errorf("configuration missing for applications")
	}
	generator := &Generator{}
	postgres := ast.GetService("postgres")
	postgREST := ast.GetService("postgrest")
	templateEngine := ast.GetService("template_engine")
	generator.Namespace = ""
	generator.Models = ast.Models
	generator.Options = make(map[string]interface{})
	if ast.AppMetadata != nil {
		generator.AppMetadata = ast.AppMetadata
	}
	if postgres != nil {
		generator.Postgres = postgres
		generator.Namespace = postgres.Namespace
	}
	if postgREST != nil {
		generator.PostgREST = postgREST
	}
	if templateEngine != nil {
		generator.TemplateEngine = templateEngine
	}
	// NOTE: LoadCondfig handles loading the environment into options. We just need to
	// copy into the Generator struct.
	for k, v := range postgres.Options {
		generator.Options[k] = v
	}
	for k, v := range postgREST.Options {
		generator.Options[k] = v
	}
	for k, v := range postgREST.Options {
		generator.Options[k] = v
	}
	for k, v := range templateEngine.Options {
		generator.Options[k] = v
	}
	return generator, nil
}

// renderPostgreSQL does what its name implies. It outputs an SQL
// program in Postgres SQL dialect. It does so for a specific type
// request. Possible values are "setup" and "models".
//
// The "setup" code type includes a placeholder for your DB credentials.
// It should not be included in your GitHub repository.
//
// The "models" contains all the table definitions, view definitions,
// and functions implementing CRUD-L operations for each model.
func (g *Generator) renderPostgreSQL(action string) error {
	switch action {
	case "setup":
		return pgSetup(g.out, g.Namespace)
	case "models":
		return pgModels(g.out, g.Namespace, g.Models)
	default:
		return fmt.Errorf("%q not supported at this time", action)
	}
}

// renderPostgREST does what its name implies. It the configuration
// file used when starting up PostgREST.
func (g *Generator) renderPostgREST() error {
	port := "5432"
	if g.Postgres != nil && g.Postgres.Port != 0 {
		port = fmt.Sprintf("%d", g.Postgres.Port)
	}
	return postgRESTConf(g.out, g.Namespace, port)
}

// renderModelActionTemplate will render a template for a given model id. The action corresponds
// to the model id.
func (g *Generator) renderModelActionTemplate(modelId string, action string) error {
	for _, model := range g.Models {
		if modelId == model.Id {
			return TmplGen(g.out, model, action)
		}
	}
	return fmt.Errorf("failed to find model id %q", modelId)
}

// renderPartialTemplate renders the head, header, nav (stub), footer partial templates
func (g *Generator) renderPartialTemplate(partial string) error {
	switch partial {
	case "head":
		return TmplHeadPartial(g.out, g.AppMetadata.AppTitle, g.AppMetadata.CSSPath)
	case "header":
		return TmplHeaderPartial(g.out, g.AppMetadata.HeaderLink, g.AppMetadata.HeaderText, g.AppMetadata.LogoLink, g.AppMetadata.LogoText)
	case "nav":
		return TmplNavPartial(g.out, "<!-- navigation goes here -->")
	case "footer":
		return TmplFooterPartial(g.out,
			g.AppMetadata.CopyrightYear,
			g.AppMetadata.CopyrightLink,
			g.AppMetadata.CopyrightText,
			g.AppMetadata.LicenseLink,
			g.AppMetadata.LicenseText,
			g.AppMetadata.ContactAddress,
			g.AppMetadata.ContactPhone,
			g.AppMetadata.ContactEMail)
	default:
		return fmt.Errorf("failed, partial %q not supported", partial)
	}
}

// renderHtml will render HTML forms for given action and model id.
func (g *Generator) renderHtml(modelId string, action string) error {
	return fmt.Errorf("g.renderHtml(%q, %q) not implemented", modelId, action)
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
func validateModelId(modelId string, models []*models.Model) error {
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
func (g *Generator) Generate(generatorName string, modelId string, action string) error {
	pgActions := []string{"setup", "models"}
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
	case "template":
		if err := validateAction(action, templateActions); err != nil {
			return err
		}
		if err := validateModelId(modelId, g.Models); err != nil {
			return err
		}
		if err := g.renderModelActionTemplate(modelId, action); err != nil {
			return err
		}
		return nil
	case "partial_template":
		return g.renderPartialTemplate(action)
	default:
		return fmt.Errorf("%q is not supported at this time", generatorName)
	}
}
