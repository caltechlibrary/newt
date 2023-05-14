/**
 * newt.go an implementation of the Newt URL router.
 *
 * @author R. S. Doiel
 */
package newt

import (
	"fmt"
	"io"
	"log"
	"os"
)

type Config struct {
	DBUser          string `json:"db_user,omitempty"`
	DBPassword      string `json:"db_name,omitempty"`
	PostgREST       string `json:"postgrest_url,omitempty"`
	PandocPort      string `json:"pandoc_port,omitempty"`
	PandocTemplates string `json:"pandoc_templates,omitempty"`
	Htdocs          string `json:"htdocs,omitempty"`
}

// Run is a runner for Newt URL router
func Run(in io.Reader, out io.Writer, eout io.Writer, args []string) int {
	cfg := new(Config)
	cfg.DBUser = os.Getenv("DB_USER")
	cfg.DBPassword = os.Getenv("DB_PASSWORD")
	cfg.PostgREST = os.Getenv("POSTGREST_URL")
	cfg.PandocPort = os.Getenv("PANDOC_PORT")
	cfg.PandocTemplates = os.Getenv("PANDOC_TEMPLATES")
	cfg.Htdocs = os.Getenv("HTDOCS")
	if len(args) > 0 {
		// Process a configuration file
		log.Println("Loading configuration file not implemented")
	}
	if len(args) > 1 {
		// Process a CSV file of routes
		log.Println("Loading routes from CSV file is not implemented")
	}
	// Calculate some sane defaults if needed
	if cfg.PostgREST == "" {
		if cfg.DBUser != "" && cfg.DBPassword != "" {
			cfg.PostgREST = fmt.Sprintf("http://%s:%s@localhost:3000", cfg.DBUser, cfg.DBPassword)
		} else if cfg.DBUser != "" {
			cfg.PostgREST = fmt.Sprintf("http://%s@localhost:3000", cfg.DBUser)
		} else {
			cfg.PostgREST = fmt.Sprintf("http://localhost:3000")
		}
	}
	if cfg.PandocPort == "" {
		cfg.PandocPort = "3030"
	}
	if cfg.Htdocs == "" {
		var err error
		cfg.Htdocs, err = os.Getwd()
		if err != nil {
			log.Printf("failed to determine htdocs directory, %s", err)
			return 1
		}
	}
	log.Println("Run() not fully implemented")
	return 0
}
