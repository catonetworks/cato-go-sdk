package gqls

import (
	"fmt"
	"io"
	"os"

	"github.com/catonetworks/cato-go-sdk/gqls/formatter"
	"github.com/rs/zerolog/log"
	"github.com/vektah/gqlparser/v2/ast"
)

func Output(schema *ast.Schema, w io.Writer) error {
	f := formatter.NewFormatter(w, formatter.WithIndent("  "), formatter.WithComments())
	return f.FormatSchema(schema)
}

// OutputFile outputs the schema to a file or stdout if fName is "".
func OutputFile(schema *ast.Schema, fName string) error {
	if fName == "" {
		return Output(schema, os.Stdout)
	}

	f, err := os.OpenFile(fName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600) //nolint: gosec
	if err != nil {
		return fmt.Errorf("failed to open output file '%s': %w", fName, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Error().Err(err).Msg("failed to close output file")
		}
	}()
	return Output(schema, f)
}
