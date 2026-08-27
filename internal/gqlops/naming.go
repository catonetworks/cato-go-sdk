package gqlops

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
)

const (
	operationKindQuery    = "query"
	operationKindMutation = "mutation"
	minOperationFileParts = 2
)

func semanticKey(fileName, kind string) (string, error) {
	base := filepath.Base(fileName)
	for _, extension := range []string{".gql", ".txt"} {
		base = strings.TrimSuffix(base, extension)
	}

	parts := strings.Split(base, ".")
	if len(parts) < minOperationFileParts {
		return "", fmt.Errorf("operation filename %q must contain a kind and path", fileName)
	}
	if parts[0] != operationKindQuery && parts[0] != operationKindMutation {
		return "", fmt.Errorf("operation filename %q has unsupported prefix %q", fileName, parts[0])
	}
	if kind != operationKindQuery && kind != operationKindMutation {
		return "", fmt.Errorf("operation %q has unsupported kind %q", fileName, kind)
	}

	return kind + "." + strings.Join(parts[1:], "."), nil
}

func normalizedKey(key string) string {
	return strings.ToLower(key)
}

func goIdentifier(name string) string {
	var result []rune
	upperNext := true

	for _, current := range name {
		if !unicode.IsLetter(current) && !unicode.IsDigit(current) {
			upperNext = true
			continue
		}
		if upperNext {
			current = unicode.ToUpper(current)
			upperNext = false
		}
		result = append(result, current)
	}

	return string(result)
}
