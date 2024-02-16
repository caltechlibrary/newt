// templates.go defines an interface used by mustache.go and pdbundler.go
package newt

import (
//	"net/http"
)

// This is the interface to a template. This allows us to have a common configuration
// YAML support both Newt Mustache and Pandoc Bundler
type Template interface {
	// This should load a template based on the contents of the implementation struct as needed
	ResolveTemplate() error
}
