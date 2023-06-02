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


type Config struct {
	// Port is the name of the localhost port Newt will listen on.
	Port string `json:"newt_port,omitempty" yaml:"newt_port,omitempty"`
	// Routes is the CSV filename that data API routes are read from
	Routes string `json:"newt_routes,omitempty" yaml:"newt_routes,omitempty"`
	// Env is a list of environment variables that can be passed
	// through to the RouteDSL when rendering JSON data API calls or
	// calls to Pandoc server.
	Env []string `json:"newt_env,omitempty" yaml:"newt_env,omitempty"`
	// Htdocs holds any static files you want to make available through
	// Newt router.
	Htdocs string `json:"newt_htdocs,omitempty" yaml:"newt_htdocs,omitempty"`
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
	// Make sure Routes CSV
	if cfg.Routes != "" {
		if _, err := os.Stat(cfg.Routes); err != nil {
			return nil, fmt.Errorf("Can't read %q, %s", cfg.Routes, err)
		}
	}
	// Make sure Htdocs exists
	if cfg.Htdocs != "" {
		if _, err := os.Stat(cfg.Htdocs); err != nil {
			return nil, fmt.Errorf("Can't read %q, %s", cfg.Htdocs, err)
		}
	}
	// Finally make sure we have cfg.Htdocs or cfg.Routes set.
	if cfg.Routes == "" && cfg.Htdocs == "" {
		return nil, fmt.Errorf("newt_routes and newt_htdocs are missing in configuration.")
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
	if cfg.Routes != "" {
		if err := router.ReadCSV(cfg.Routes); err != nil {
			fmt.Fprintf(eout, "error reading routes from %q, %s\n", cfg.Routes, err)
			return 1
		}
	}
	if dryRun {
		fmt.Fprintf(out, "configuration and routes successfully processed\n")
		return 0
	}


	appName := path.Base(os.Args[0])
	mux := http.NewServeMux()
	switch {
	case cfg.Htdocs != "" && cfg.Routes != "":
		log.Printf("%s using %s for static content and %s for router", appName, cfg.Htdocs, cfg.Routes)
		mux.Handle("/", wsfn.RequestLogger(router.Newt(http.FileServer(http.Dir(cfg.Htdocs)))))
	case cfg.Htdocs == "" && cfg.Routes != "":
		log.Printf("%s using %s for router only", appName, cfg.Routes)
		mux.Handle("/", wsfn.RequestLogger(router.Newt(http.NotFoundHandler())))
	case cfg.Htdocs != "" && cfg.Routes == "":
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
