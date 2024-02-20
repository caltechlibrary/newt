package newt

import (
	"fmt"
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

// Generator generates the code based on the contents of Generator struct.
// FIXME: Not implemented yet.
func (g *NewtGenerator) Generate(target string, codeType string) error {
	return fmt.Errorf("g.Generate(%q, %q) not implemented", target, codeType)
}
