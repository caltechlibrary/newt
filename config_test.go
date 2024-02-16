package newt

import (
	"os"
	"path"
	"testing"
)

// Test the shared YAML configuration for Newt Router, Newt Generator
// Newt Mustache and Pandoc Bundler.
func TestLoadConfig(t *testing.T) {
	configFiles := []string{
		path.Join("testdata", "blog.yaml"),
		path.Join("testdata", "bundler_test.yaml"),
	}
	for _, fName := range configFiles {
		cfg, err := LoadConfig(fName)
		if err != nil {
			t.Errorf("failed to load %q, %s", fName, err)
		}
		if cfg == nil {
			t.Errorf("something went wrong, cfg is nil for %q", fName)
		}
		if cfg.Application == nil {
			t.Errorf("cfg.Application is nil (%q), %+v", fName, cfg)
		}
	}
}


