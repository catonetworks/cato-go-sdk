package gqlops

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/validator/rules"
)

const (
	maxOperationBytes = 8 << 20
	maxSchemaBytes    = 64 << 20
)

type document struct {
	key      string
	kind     string
	name     string
	path     string
	relative string
	content  []byte
	hash     string
}

func loadSchema(sdkRoot string) (*ast.Schema, error) {
	path := filepath.Join(sdkRoot, "cato_api.graphqls")
	content, err := readRegularFile(path, maxSchemaBytes)
	if err != nil {
		return nil, err
	}

	schema, parseErr := gqlparser.LoadSchema(&ast.Source{Name: path, Input: string(content)})
	if parseErr != nil {
		return nil, fmt.Errorf("parse schema %q: %w", path, parseErr)
	}

	return schema, nil
}

func loadSDKDocuments(sdkRoot string, schema *ast.Schema) ([]document, error) {
	sourcePattern := filepath.Join(sdkRoot, "sources", "*.gql")
	paths, err := filepath.Glob(sourcePattern)
	if err != nil {
		return nil, fmt.Errorf("scan SDK operations: %w", err)
	}
	sort.Strings(paths)

	documents := make([]document, 0, len(paths))
	var validationErrors []error

	for _, path := range paths {
		content, readErr := readRegularFile(path, maxOperationBytes)
		if readErr != nil {
			validationErrors = append(validationErrors, readErr)
			continue
		}
		parsed, parseErr := parseOperation(schema, path, content)
		if parseErr != nil {
			validationErrors = append(validationErrors, parseErr)
			continue
		}

		key, keyErr := semanticKey(path, parsed.kind)
		if keyErr != nil {
			validationErrors = append(validationErrors, keyErr)
			continue
		}
		parsed.key = key
		parsed.path = path
		parsed.relative = filepath.ToSlash(filepath.Join("sources", filepath.Base(path)))
		parsed.content = content
		parsed.hash = contentHash(content)
		documents = append(documents, parsed)
	}

	if err := errors.Join(validationErrors...); err != nil {
		return nil, err
	}

	return documents, nil
}

func parseOperation(schema *ast.Schema, path string, content []byte) (document, error) {
	query, err := gqlparser.LoadQueryWithRules(schema, string(content), rules.NewDefaultRules())
	if err != nil {
		return document{}, fmt.Errorf("validate operation %q: %w", path, err)
	}
	if len(query.Operations) != 1 {
		return document{}, fmt.Errorf("operation file %q contains %d operations; want 1", path, len(query.Operations))
	}

	operation := query.Operations[0]
	kind := string(operation.Operation)
	if kind != operationKindQuery && kind != operationKindMutation {
		return document{}, fmt.Errorf("operation file %q has unsupported kind %q", path, kind)
	}
	if operation.Name == "" {
		return document{}, fmt.Errorf("operation file %q has no operation name", path)
	}

	return document{kind: kind, name: operation.Name}, nil
}

func validateUniqueDocuments(documents []document) error {
	keys := make(map[string]string, len(documents))
	names := make(map[string]string, len(documents))
	goNames := make(map[string]string, len(documents))
	var validationErrors []error

	for _, current := range documents {
		key := normalizedKey(current.key)
		if previous, exists := keys[key]; exists {
			validationErrors = append(validationErrors, fmt.Errorf(
				"semantic key %q is duplicated by %q and %q",
				current.key,
				previous,
				current.relative,
			))
		} else {
			keys[key] = current.relative
		}

		if previous, exists := names[current.name]; exists {
			validationErrors = append(validationErrors, fmt.Errorf(
				"operation name %q is duplicated by %q and %q",
				current.name,
				previous,
				current.relative,
			))
		} else {
			names[current.name] = current.relative
		}

		goName := goIdentifier(current.name)
		if previous, exists := goNames[goName]; exists {
			validationErrors = append(validationErrors, fmt.Errorf(
				"generated Go identifier %q is duplicated by %q and %q",
				goName,
				previous,
				current.relative,
			))
		} else {
			goNames[goName] = current.relative
		}
	}

	return errors.Join(validationErrors...)
}

func readRegularFile(path string, limit int64) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, fmt.Errorf("inspect %q: %w", path, err)
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("refuse non-regular file %q", path)
	}

	file, err := os.Open(path) // #nosec G304 -- callers supply confined paths; symlinks and non-regular files are rejected above.
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}
	defer func() {
		_ = file.Close()
	}()

	content, err := io.ReadAll(io.LimitReader(file, limit+1))
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", path, err)
	}
	if int64(len(content)) > limit {
		return nil, fmt.Errorf("file %q exceeds %d bytes", path, limit)
	}

	return content, nil
}

func contentHash(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

func confinedPath(root, path string) error {
	absoluteRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolve root %q: %w", root, err)
	}
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolve path %q: %w", path, err)
	}

	relative, err := filepath.Rel(absoluteRoot, absolutePath)
	if err != nil {
		return fmt.Errorf("compare path %q to root %q: %w", path, root, err)
	}
	if relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return fmt.Errorf("path %q escapes root %q", path, root)
	}

	return nil
}
