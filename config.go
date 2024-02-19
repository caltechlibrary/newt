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
	"strconv"
	"strings"

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

	// Environment holds a list of OS environment variables that can be made
	// available to the web services.
	Environment []string `json:"enviroment,omitempty" yaml:"enviroment"`

	// Options is a map of name to string values
  Options map[string] string `json:"options,omitempty" yaml:"options"`
}

// Application implements runtime config for Newt programs
type Application struct {
	// Namespace holds the Postgres Schema name It is used to generate
	// a setup.sql file using the -pg-setup option in newt cli.
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`

	// Port is the name of the localhost port Newt will listen on.
	Port int `json:"port,omitempty" yaml:"port,omitempty"`

	// Htdocs holds any static files you want to make available through
	// Newt router.
	Htdocs string `json:"htdocs,omitempty" yaml:"htdocs,omitempty"`
}


// ConfigUnmarshal will read []byte of YAML or JSON, 
// populate the provided *Config object and return an error.
//
//```
// src, _ := os.ReadFile("newt.yaml")
// cfg := new(Config)
// if err := ConfigUnmarshal(src, cfg); err != nil {
//     // ... handle error
// }
//```
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
	return nil
}

// LoadConfig read a configuration file, merges environment variables
// and returns a Config object and error value.
//
//```
// cfg, err := LoadConfig("newt.yaml")
// if err != nil {
//     // ... handle error
// }
//```
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

	// Load environment if missing from config file.
	if cfg.Application == nil {
		cfg.Application = new(Application)
	}
	if len(cfg.Application.Environment) == 0 {
		s := os.Getenv("NEWT_ENV")
		if s != "" {
			for _, envar := range strings.Split(s, ";") {
				cfg.Application.Environment = append(cfg.Application.Environment, strings.TrimSpace(envar))
			}
		}
	}
	if cfg.Application.Port == 0 {
		val := os.Getenv("NEWT_PORT")
		if val != "" {
			port, err := strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("NEWT_PORT holds an port number, %s", err)
			}
			cfg.Application.Port = port
		}
	}
	if cfg.Application.Htdocs == "" {
		cfg.Application.Htdocs = os.Getenv("NEWT_HTDOCS")
	}

	//
	// Sanity check the configuration.
	//

	// Make sure Htdocs exists
	if cfg.Application.Htdocs != "" {
		if _, err := os.Stat(cfg.Application.Htdocs); err != nil {
			dir, _ := os.Getwd()
			return nil, fmt.Errorf("Can't read %q from %s, %s", cfg.Application.Htdocs, dir, err)
		}
	}
	return cfg, nil
}

