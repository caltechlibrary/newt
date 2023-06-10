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

	// Caltech Library Packages
	"github.com/caltechlibrary/wsfn"

	// 3rd Party
	"gopkg.in/yaml.v3"
)

// Config holds a configuration for Newt.
type Config struct {
	// Port is the name of the localhost port Newt will listen on.
	Port string `json:"port,omitempty" yaml:"port,omitempty"`
	// Routes is the YAML filename that data API routes are read from
	FName string `json:"route_file,omitempty" yaml:"route_file,omitempty"`
	// Env is a list of environment variables that can be passed
	// through to the RouteDSL when rendering JSON data API calls or
	// calls to Pandoc server.
	Env []string `json:"env,omitempty" yaml:"env,omitempty"`
	// Htdocs holds any static files you want to make available through
	// Newt router.
	Htdocs string `json:"htdocs,omitempty" yaml:"htdocs,omitempty"`
	// Routes holds an array of maps of route definitions
	Routes []*Route `json:"routes,omitempty" yaml:"routes,omitempty"`
}

// LoadConfig loads a configuration return a Config object and error value.
func LoadConfig(configFName string) (*Config, error) {
	cfg := new(Config)
	if configFName != "" {
		src, err := os.ReadFile(configFName)
		if err != nil {
			return nil, fmt.Errorf("failed to read %q, %s", configFName, err)
		}
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
	// Finally make sure we have cfg.Htdocs or cfg.Routes set.
	if cfg.Routes == nil && cfg.Htdocs == "" {
		return nil, fmt.Errorf("NEWT_ROUTES and NEWT_HTDOCS are undefined.")
	}
	return cfg, nil
}

// Run is a runner for Newt URL router and static file server
func Run(in io.Reader, out io.Writer, eout io.Writer, args []string, dryRun bool) int {
	configFName := ""
	if len(args) > 0 {
		configFName = args[0]
	}
	cfg, err := LoadConfig(configFName)
	if err != nil {
		fmt.Fprintf(eout, "%s\n", err)
		return 1
	}

	router := new(Router)
	if cfg != nil {
		router.Configure(cfg)
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
