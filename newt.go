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
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	// Caltech Library Packages
	"github.com/caltechlibrary/wsfn"

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

// ParseConfig will read []byte of YAML or JSON, 
// populate the provided *Config object and return an error.
//
//```
// src, _ := os.ReadFile("newt.yaml")
// cfg := new(Config)
// if err := ParseConfig(src, cfg); err != nil {
//     // ... handle error
// }
//```
func Parse(src []byte, cfg *Config) error {
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
func Load(configFName string) (*Config, error) {
	cfg := new(Config)
	if configFName != "" {
		src, err := os.ReadFile(configFName)
		if err != nil {
			return nil, fmt.Errorf("failed to read %q, %s", configFName, err)
		}
		if err := Parse(src, cfg); err != nil {
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

// RunPostgresSQL is a runner for generating SQL from a Newt YAML file.
func RunPostgresSQL(in io.Reader, out io.Writer, eout io.Writer, args []string, pgSetupSQL bool, pgModelsSQL bool, pgModelsTestSQL bool) int {
	configFName := ""
	if len(args) > 0 {
		configFName = args[0]
	}
	cfg, err := Load(configFName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	if cfg.Models == nil {
		fmt.Fprintln(eout, "-- No modules defined.")
		return 1
	}
	// For each module we generate a table create statement,
	// default view.
	if configFName == "" {
		configFName = "standard input"
	}
	exitCode := 0
	if pgSetupSQL {
		src, err := PgSetupSQL(configFName, cfg.Namespace, "")
		if err != nil {
			fmt.Fprintf(eout, "-- could not create setup for %q based on %q, %s\n", cfg.Namespace, configFName, err)
			exitCode = 1
		} else {
			fmt.Fprintf(out, "%s\n", src)
		}
	}
	if pgModelsSQL {
		modelNames := []string{}
		for i, model := range cfg.Models {
			src, err := PgModelSQL(configFName, model)
			if err != nil {
				fmt.Fprintf(eout, "-- could not create model %q (%d), %s\n", model.Name, i, err)
				exitCode = 1
			} else {
				fmt.Fprintf(out, "%s\n", src)
			}
			modelNames = append(modelNames, model.Name)
		}

		src, err := PgModelPermissions(configFName, cfg.Namespace, modelNames)
		if err != nil {
			fmt.Fprintf(eout, "-- could not permissions for models in %q, %s\n", configFName, err)
			exitCode = 1
		} else {
			fmt.Fprintf(out, "%s\n", src)
		}
	}
	if pgModelsTestSQL {
		for name, model := range cfg.Models {
			src, err := PgModelTestSQL(configFName, model)
			if err != nil {
				fmt.Fprintf(eout, "-- could not create model test %q, %s\n", name, err)
				exitCode = 1
			} else {
				fmt.Fprintf(out, "%s\n", src)
			}
		}
	}
	return exitCode
}

// RunMustache is a runner for a Mustache redner engine service based on the Pandoc server API.
func RunMustache(in io.Reader, out io.Writer, eout io.Writer, port string, timeout int, verbose bool) int {
	exitCode := 0
	err := MustacheServer(out, eout, port, timeout, verbose)
	if err != nil {
		fmt.Fprintf(eout, "error starting server: %s\n", err)
		exitCode = 1
	}
	return exitCode
}

// Run is a runner for Newt URL router and static file server
func Run(in io.Reader, out io.Writer, eout io.Writer, args []string, dryRun bool) int {
	configFName := ""
	if len(args) > 0 {
		configFName = args[0]
	}
	cfg, err := Load(configFName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}
	// Finally make sure we have cfg.Htdocs or cfg.Routes 
	// set to run service.
	if cfg.Routes == nil && cfg.Htdocs == "" {
		if cfg.FName != "" {
			fmt.Fprintf(eout, "routes and htdocs are not set.")
			return 1
		} 
		fmt.Fprintf(eout, "NEWT_ROUTES and NEWT_HTDOCS are undefined.")
		return 1
	}

	router := new(Router)
	if cfg != nil {
		if err := router.Configure(cfg); err != nil {
			fmt.Fprintf(eout, "error configuring router from %q, %s\n", cfg.FName, err)
			return 1
		}
	}
	if err != nil {
		fmt.Fprintf(eout, "error reading routes from %q, %s\n", cfg.FName, err)
		return 1
	}
	if dryRun {
		fmt.Fprintf(out, "configuration and routes successfully processed\n")
		return 0
	}

	appName := path.Base(os.Args[0])
	mux := http.NewServeMux()
	switch {
	case cfg.Htdocs != "" && cfg.Routes != nil:
		log.Printf("%s using %s for static content and %s for router", appName, cfg.Htdocs, cfg.Routes)
		mux.Handle("/", wsfn.RequestLogger(router.Newt(http.FileServer(http.Dir(cfg.Htdocs)))))
	case cfg.Htdocs == "" && cfg.Routes != nil:
		log.Printf("%s using %s for router only", appName, cfg.Routes)
		mux.Handle("/", wsfn.RequestLogger(router.Newt(http.NotFoundHandler())))
	case cfg.Htdocs != "" && cfg.Routes == nil:
		log.Printf("%s using %s for static content only", appName, cfg.Htdocs)
		mux.Handle("/", wsfn.RequestLogger(http.FileServer(http.Dir(cfg.Htdocs))))
	default:
		log.Printf("Not configured, aborting")
		return 1
	}

	log.Printf("%s listening on port %s", appName, cfg.Port)
	http.ListenAndServe(":"+cfg.Port, mux)
	return 0
}


