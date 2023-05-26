/**
 * newt.go an implementation of the Newt URL router.
 *
 * @author R. S. Doiel
 */
package newt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	// 3rd Party
	"gopkg.in/yaml.v3"
)

type Config struct {
	// Port is the name of the localhost port Newt will listen on.
	Port string `json:"newt_port,omitempty" yaml:"newt_port,omitempty"`
	// Routes is the CSV filename that routes are read from
	Routes string `json:"newt_routes,omitempty" yaml:"newt_routes,omitempty"`
	// Env is a list of environment variables that can be passed
	// through to the RouteDSL when rendering JSON data API calls or
	// calls to Pandoc server.
	Env []string `json:"newt_env,omitempty" yaml:"newt_env,omitempty"`
}

// parseConf will parse a simple key equal value relationship like used in
// PostgREST's configuration file.

// LoadConfig loads a configuration return a Config object and error value.
func LoadConfig(configFName string) (*Config, error) {
	cfg := new(Config)
	if configFName != "" {
		src, err := os.ReadFile(configFName)
		if err != nil {
			return nil, fmt.Errorf("failed to read %q, %s", configFName, err)
		}
		//log.Printf("DEBUG src -> %s", src)
		if bytes.HasPrefix(src, []byte("{")) {
			if err := json.Unmarshal(src, &cfg); err != nil {
				return nil, fmt.Errorf("failed to read JSON %q, %s", configFName, err)
			}
		} else {
			if err := yaml.Unmarshal(src, &cfg); err != nil {
				return nil, fmt.Errorf("failed YAML parse %q, %s", configFName, err)
			}
		}
	}
	// Load environment if missing from config file.
	if cfg.Port == "" {
		cfg.Port = os.Getenv("NEWT_PORT")
	}
	if cfg.Routes == "" {
		cfg.Routes = os.Getenv("NEWT_ROUTES")
	}
	if len(cfg.Env) == 0 {
		s := os.Getenv("NEWT_ENV")
		if s != "" {
			for _, envar := range strings.Split(s, ";") {
				cfg.Env = append(cfg.Env, strings.TrimSpace(envar))
			}
		}
	}
	// Sanity check the configuration.
	if cfg.Port == "" {
		cfg.Port = "4040"
	} else {
		if strings.HasPrefix(cfg.Port, ":") {
			cfg.Port = cfg.Port[1:]
		}
		if _, err := strconv.Atoi(cfg.Port); err != nil {
			return nil, fmt.Errorf("Expected port to be an integer, %s", err)
		}
	}
	return cfg, nil
}

// Run is a runner for Newt URL router
func Run(in io.Reader, out io.Writer, eout io.Writer, args []string, dryRun bool) int {
	configFName := ""
	if len(args) > 0 {
		configFName = args[0]
	}
	cfg, err := LoadConfig(configFName)
	if err != nil {
		log.Printf("%s", err)
		return 1
	}

	router := new(Router)
	if err := router.ReadCSV(cfg.Routes); err != nil {
		log.Printf("error reading %q, %s", cfg.Routes, err)
		return 1
	}
	if dryRun {
		log.Printf("configuration and routes successfully processed")
		return 0
	}
	http.HandleFunc("/", router.Handler)
	log.Printf("%s listening on port %s", path.Base(os.Args[0]), cfg.Port)
	http.ListenAndServe(":"+cfg.Port, nil)
	return 0
}
