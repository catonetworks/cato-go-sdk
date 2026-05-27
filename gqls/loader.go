package gqls

import (
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/vektah/gqlparser/v2/ast"
)

type Loader struct {
}

// LoadSDL loads a schema from a file or stdin
func LoadSDL(fName string) (s *ast.Schema, err error) {
	input, fName, closeFunc, err := getReader(fName)
	if err != nil {
		return nil, err
	}
	defer closeFunc()
	return ParseSDLSchema(input, fName)
}

// LoadJSON loads a schema from a file or stdin
func LoadJSON(fName string) (s *ast.Schema, err error) {
	input, fName, closeFunc, err := getReader(fName)
	if err != nil {
		return nil, err
	}
	defer closeFunc()
	return ParseJSONSchema(input, fName)
}

// getReader opens a file or stdin
func getReader(fName string) (reader io.Reader, fileName string, closeFunc func(), err error) {
	if fName == "" {
		return os.Stdin, "<stdin>", func() {}, nil
	}

	fileInput, err := os.Open(fName) // #nosec G304
	if err != nil {
		return nil, fName, func() {}, fmt.Errorf("failed to open file '%s': %w", fName, err)
	}
	return fileInput, fName, func() {
		if err := fileInput.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close input file")
		}
	}, nil
}
