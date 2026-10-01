package gqlops

import (
	"crypto/sha256"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"unicode"
)

func handwrittenGoNames(sdkRoot string) (map[string]struct{}, error) {
	paths, err := filepath.Glob(filepath.Join(sdkRoot, "*.go"))
	if err != nil {
		return nil, fmt.Errorf("scan handwritten SDK files: %w", err)
	}

	names := make(map[string]struct{})
	files := token.NewFileSet()
	for _, path := range paths {
		base := filepath.Base(path)
		if base == "client.go" || strings.HasSuffix(base, "_test.go") {
			continue
		}
		parsed, parseErr := parser.ParseFile(files, path, nil, 0)
		if parseErr != nil {
			return nil, fmt.Errorf("parse handwritten SDK file %q: %w", path, parseErr)
		}
		for _, declaration := range parsed.Decls {
			switch current := declaration.(type) {
			case *ast.FuncDecl:
				if current.Recv == nil || receiverIsClient(current.Recv) {
					names[current.Name.Name] = struct{}{}
				}
			case *ast.GenDecl:
				for _, spec := range current.Specs {
					if typed, ok := spec.(*ast.TypeSpec); ok {
						names[typed.Name.Name] = struct{}{}
					}
				}
			}
		}
	}
	return names, nil
}

func receiverIsClient(receivers *ast.FieldList) bool {
	if receivers == nil || len(receivers.List) != 1 {
		return false
	}
	expression := receivers.List[0].Type
	if pointer, ok := expression.(*ast.StarExpr); ok {
		expression = pointer.X
	}
	identifier, ok := expression.(*ast.Ident)
	return ok && identifier.Name == "Client"
}

func stableOperationName(
	cliName string,
	key string,
	indexByName map[string]int,
	reservedGoNames map[string]struct{},
) string {
	if operationNameAvailable(cliName, indexByName, reservedGoNames) {
		return cliName
	}

	kind, _, _ := strings.Cut(key, ".")
	candidate := kind + upperInitial(cliName)
	if operationNameAvailable(candidate, indexByName, reservedGoNames) {
		return candidate
	}

	digest := sha256.Sum256([]byte(normalizedKey(key)))
	return fmt.Sprintf("%s%x", candidate, digest[:4])
}

func operationNameAvailable(
	name string,
	indexByName map[string]int,
	reservedGoNames map[string]struct{},
) bool {
	if _, exists := indexByName[name]; exists {
		return false
	}
	_, reserved := reservedGoNames[goIdentifier(name)]
	return !reserved
}

func upperInitial(value string) string {
	runes := []rune(value)
	if len(runes) == 0 {
		return value
	}
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
