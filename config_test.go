package newt

import (
	"path"
	"testing"
)

// Test the shared YAML configuration for Newt Router, Newt Generator
// Newt Mustache and Pandoc Bundler.
func TestLoadConfig(t *testing.T) {
	configFiles := []string{
		path.Join("testdata", "birds.yaml"),
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
		if cfg.Applications == nil {
			t.Errorf("cfg.Applications is nil (%q), %+v", fName, cfg)
		}
		ids := cfg.GetModelIds()
		if len(ids) == 0 {
			t.Errorf("expected model ids for %q", fName)
		} else {
			mId := ids[0]
			model, ok := cfg.GetModelById(mId)
			if !ok {
				t.Errorf("expected model for %q in %q, %s", mId, fName, err)
			}
			if model == nil {
				t.Errorf("expceted model content for %q in %q, got nil", mId, fName)
			}
		}
		names := cfg.GetModelNames()
		if len(names) != len(ids) {
			t.Errorf("expected %d model names for %q, got %d", len(ids), fName, len(names))
		}
	}
}
