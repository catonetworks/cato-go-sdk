package gqls

import (
	"bytes"
	"fmt"
	"io"

	"github.com/catonetworks/cato-go-sdk/gqls/introspect"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

// ParseSDLSchema parses a graphql schema
func ParseSDLSchema(input io.Reader, fName string) (s *ast.Schema, err error) {
	// read contents
	contents, err := io.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("failed to load GrapahQL schema from '%s': %w", fName, err)
	}

	// parse schema
	schema, err := gqlparser.LoadSchema(&ast.Source{Name: fName, Input: string(contents)})
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema from '%s': %w", fName, err)
	}

	return schema, nil
}

// ParseJSONSchema parses a graphql introspection JSON
func ParseJSONSchema(input io.Reader, fName string) (s *ast.Schema, err error) {
	// read contents
	contents, err := io.ReadAll(input)
	if err != nil {
		return nil, fmt.Errorf("failed to load JSON introspection from '%s': %w", fName, err)
	}

	// parse schema
	s, err = introspect.ParseSchema(contents, fName)
	if err != nil {
		return nil, err
	}
	// return s, nil

	var buf bytes.Buffer
	if err = Output(s, &buf); err != nil {
		return nil, err
	}

	return ParseSDLSchema(&buf, fName)
}
