/**
 * config.go holds the data structure that is shared between Newt applications.
 *
 * @author R. S. Doiel
 */
package newt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	// 3rd Party
	"gopkg.in/yaml.v3"
)

// Config holds a configuration for Newt for the data router and code generator.
type Config struct {
	// Application holds the application specific settings and metadata
	Applications *Applications `json:"applications,omitempty" yaml:"applications"`

	// Models holds a list of data models. It is used by
	// both the data router and code generator.
	Models []*NewtModel `json:"models,omitempty" yaml:"models,omitempty"`

	// Routes holds an array of maps of route definitions used by
	// the data router and code generator
	Routes []*NewtRoute `json:"routes,omitempty" yaml:"routes,omitempty"`

	// Templates holds an array of mapsthe request to template to request for
	// Newt Mustache
	Templates []*MustacheTemplate `json:"templates,omitempty" yaml:"templates,omitempty"`
}

// Applications holds the runtime information for newt router, generator,
// template engine.
type Applications struct {
	// Newt Router runtime config
	NewtRouter *Application `json:"newtrouter,omitempty" yaml:"newtrouter"`

	// Newt Mustache runtime config
	NewtMustache *Application `json:"newtmustache,omitempty" yaml:"newtmustache"`

	// Newt Generator runtime config
	NewtGenerator *Application `json:"newtgenerator,omitempty" yaml:"newtgenerator"`

	// Postgres runtime config, e.g. port number to use for connecting.
	Postgres *Application `json:"postgres,omitempty" yaml:"postgres,omitempty"`

	// PostgREST runtime config
	PostgREST *Application `json:"postgrest,omitempty" yaml:"postgrest"`

	// Environment holds a list of OS environment variables that can be made
	// available to the web services.
	Environment []string `json:"enviroment,omitempty" yaml:"enviroment"`

	// Options is a map of name to string values, it is where the
	// environment variables values are stored.
	Options map[string]string `json:"options,omitempty" yaml:"options"`
}

// Application implements runtime config for Newt programs
type Application struct {
	// AppPath holds the path to the binary application, e.g. PostgREST
	// This property provides the location of the service to run.
	AppPath string `json:"app_path,omitempty" yaml:"app_path"`

	// ConfPath holds teh path to the configuration file (e.g. PostgREST configuration file)
	ConfPath string `json:"conf_path,omitempty" yaml:"conf_path"`

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
}

// ConfigUnmarshal will read []byte of YAML or JSON,
// populate the provided *Config object and return an error.
//
// ```
// src, _ := os.ReadFile("newt.yaml")
// cfg := new(Config)
//
//	if err := ConfigUnmarshal(src, cfg); err != nil {
//	    // ... handle error
//	}
//
// ```
func ConfigUnmarshal(src []byte, cfg *Config) error {
	if bytes.HasPrefix(src, []byte("{")) {
		if err := json.Unmarshal(src, &cfg); err != nil {
			return err
		}
	} else {
		if err := yaml.Unmarshal(src, &cfg); err != nil {
			return err
		}
	}
	if cfg.Applications == nil {
		cfg.Applications = &Applications{
			NewtRouter:    &Application{},
			NewtMustache:  &Application{},
			NewtGenerator: &Application{},
			Options:       map[string]string{},
			Environment:   []string{},
		}
	}
	return nil
}

// LoadConfig read a configuration file, merges environment variables
// and returns a Config object and error value.
//
// ```
// cfg, err := LoadConfig("newt.yaml")
//
//	if err != nil {
//	    // ... handle error
//	}
//
// ```
func LoadConfig(configFName string) (*Config, error) {
	cfg := new(Config)
	if configFName != "" {
		src, err := os.ReadFile(configFName)
		if err != nil {
			return nil, fmt.Errorf("failed to read %q, %s", configFName, err)
		}
		if err := ConfigUnmarshal(src, cfg); err != nil {
			return nil, fmt.Errorf("failed to read %q, %s", configFName, err)
		}
	}

	if cfg.Applications == nil {
		cfg.Applications = &Applications{
			NewtRouter:    &Application{},
			NewtMustache:  &Application{},
			NewtGenerator: &Application{},
			Options:       map[string]string{},
			Environment:   []string{},
		}
	}
	// Load environment if missing from config file.
	if len(cfg.Applications.Environment) == 0 {
		for _, envar := range cfg.Applications.Environment {
			// YAML settings take presidence over environment, check for conflicts
			if _, conflict := cfg.Applications.Options[envar]; !conflict {
				cfg.Applications.Options[envar] = os.Getenv(envar)
			}
		}
	}
	return cfg, nil
}

// GetModelIds returns a list of model ids
func (cfg *Config) GetModelIds() []string {
	ids := []string{}
	for _, m := range cfg.Models {
		if m.Id != "" {
			ids = append(ids, m.Id)
		}
	}
	return ids
}

// GetModelNames returns a list of model names (not to be confused with Model ids)
func (cfg *Config) GetModelNames() []string {
	names := []string{}
	for _, m := range cfg.Models {
		if m.Name != "" {
			names = append(names, m.Name)
		}
	}
	return names
}

// GetModelById return a specific model by it's id
func (cfg *Config) GetModelById(id string) (*NewtModel, bool) {
	for _, m := range cfg.Models {
		if m.Id == id {
			return m, true
		}
	}
	return nil, false
}

/*
// GetTemplateIds return a list of template ids
func (cfg *Config) GetTemplateIds() []string {
	ids := []string{}
	for _, t := range cfg.Templates {
		if t.Id != "" {
			ids = append(ids, t.Id)
		}
	}
	return ids
}

// GetTemplateById return a a list
func (cfg *Config) GetTemplateById(id string) (*NewtTemplate, bool) {
	for _, t := range cfg.Models {
		if t.Id == id {
			return t, true
		}
	}
	return nil, false
}
*/
