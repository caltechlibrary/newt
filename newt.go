/**
 * newt.go an implementation of the Newt data router and Newt Postgres SQL generator.
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

// Config holds a configuration for Newt.
type Config struct {
	// Namespace holds the Postgres Schema name It is used to generate
	// a setup.sql file using the -pg-setup option in newt cli.
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	// Models declare the data structures used in a Newt application.
	// These generally these are tables but can be thought of as objects.
	// Models are used to generate SQL CREATE statements suitable for
	// use with Postgres.
	Models []*ModelDSL `json:"models,omitempty" yaml:"models,omitempty"`
	// Port is the name of the localhost port Newt will listen on.
	Port string `json:"port,omitempty" yaml:"port,omitempty"`
	// Routes is the YAML filename that data API routes are read from
	FName string `json:"route_file,omitempty" yaml:"route_file,omitempty"`
	// Env is a list of environment variables that can be passed
	// through to the RouteDSL when rendering JSON data API calls or
	// calls to Pandoc server or Mustache server.
	Env []string `json:"env,omitempty" yaml:"env,omitempty"`
	// Htdocs holds any static files you want to make available through
	// Newt router.
	Htdocs string `json:"htdocs,omitempty" yaml:"htdocs,omitempty"`
	// Routes holds an array of maps of route definitions
	Routes []map[string]interface{} `json:"routes,omitempty" yaml:"routes,omitempty"`
}

// parseConfig will read []byte of YAML or JSON, 
// populate the provided *Config object and return an error.
//
//```
// src, _ := os.ReadFile("newt.yaml")
// cfg := new(Config)
// if err := parseConfig(src, cfg); err != nil {
//     // ... handle error
// }
//```
func parseConfig(src []byte, cfg *Config) error {
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

// Load read a configuration file and returns a Config object and error value.
// It merges in the Newt specific environment variables.
//
//```
// cfg, err := Load("newt.yaml")
// if err != nil {
//     // ... handle error
// }
//```
func NewtLoad(configFName string) (*Config, error) {
	cfg := new(Config)
	if configFName != "" {
		src, err := os.ReadFile(configFName)
		if err != nil {
			return nil, fmt.Errorf("failed to read %q, %s", configFName, err)
		}
		if err := parseConfig(src, cfg); err != nil {
			return nil, fmt.Errorf("failed to read %q, %s", configFName, err)
		}
	}

	// Load environment if missing from config file.
	if cfg.Port == "" {
		cfg.Port = os.Getenv("NEWT_PORT")
	}
	if cfg.FName == "" {
		cfg.FName = os.Getenv("NEWT_ROUTES")
	}
	if len(cfg.Env) == 0 {
		s := os.Getenv("NEWT_ENV")
		if s != "" {
			for _, envar := range strings.Split(s, ";") {
				cfg.Env = append(cfg.Env, strings.TrimSpace(envar))
			}
		}
	}
	if cfg.Htdocs == "" {
		cfg.Htdocs = os.Getenv("NEWT_HTDOCS")
	}
	if cfg.Port == "" {
		cfg.Port = "8000"
	}
	if strings.HasPrefix(cfg.Port, ":") {
		cfg.Port = cfg.Port[1:]
	}

	//
	// Sanity check the configuration.
	//

	// Make sure port is an integer.
	if _, err := strconv.Atoi(cfg.Port); err != nil {
		return nil, fmt.Errorf("Expected port to be an integer, %s", err)
	}
	// Make sure Routes YAML file is readable
	if cfg.FName != "" {
		if _, err := os.Stat(cfg.FName); err != nil {
			return nil, fmt.Errorf("Can't read %q, %s", cfg.FName, err)
		}
	}
	// Make sure Htdocs exists
	if cfg.Htdocs != "" {
		if _, err := os.Stat(cfg.Htdocs); err != nil {
			dir, _ := os.Getwd()
			return nil, fmt.Errorf("Can't read %q from %s, %s", cfg.Htdocs, dir, err)
		}
	}
	return cfg, nil
}


