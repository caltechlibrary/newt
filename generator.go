package newt

import (
	"fmt"
)

// Generator holds our Newt Generator structures for rendering code.
// FIXME: Not implemented yet.
type Generator struct {
}

// NewGenerator instaitates a new Generator object form a filename and Config object
// It returns a Generator object and error value.
func NewGenerator(fName string, cfg *Config) (*Generator, error) {
	return nil, fmt.Errorf("NewGenerator() not implemented")
}

// Generator generates the code based on the contents of Generator struct.
// FIXME: Not implemented yet.
func (g *Generator) Generate() error {
	return fmt.Errorf("g.Generate() not implemented")
}
