package gqlops

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"github.com/vektah/gqlparser/v2/parser"
)

type aliasRecord struct {
	alias     string
	ambiguous bool
}

func normalizeMappedOperation(
	schema *ast.Schema,
	path string,
	existingContent []byte,
	incomingContent []byte,
	stableName string,
) ([]byte, error) {
	existing, err := parseQueryDocument(path, existingContent)
	if err != nil {
		return nil, fmt.Errorf("parse existing operation: %w", err)
	}
	incoming, err := parseQueryDocument(path, incomingContent)
	if err != nil {
		return nil, fmt.Errorf("parse incoming operation: %w", err)
	}

	existingOperation := existing.Operations[0]
	incomingOperation := incoming.Operations[0]
	incomingOperation.Name = stableName
	incomingOperation.VariableDefinitions = stableVariableOrder(
		existingOperation.VariableDefinitions,
		incomingOperation.VariableDefinitions,
	)

	aliases := collectAliases(existingOperation.SelectionSet, nil)
	applyAliases(incomingOperation.SelectionSet, nil, aliases)
	applyFragmentAliases(existing, incoming)

	return formatAndValidate(schema, path, incoming, stableName)
}

func normalizeNewOperation(
	schema *ast.Schema,
	path string,
	incomingContent []byte,
	stableName string,
) ([]byte, error) {
	incoming, err := parseQueryDocument(path, incomingContent)
	if err != nil {
		return nil, err
	}
	incoming.Operations[0].Name = stableName
	incoming.Operations[0].VariableDefinitions = stableVariableOrder(
		nil,
		incoming.Operations[0].VariableDefinitions,
	)

	return formatAndValidate(schema, path, incoming, stableName)
}

func parseQueryDocument(path string, content []byte) (*ast.QueryDocument, error) {
	query, err := parser.ParseQuery(&ast.Source{Name: path, Input: string(content)})
	if err != nil {
		return nil, fmt.Errorf("parse operation %q: %w", path, err)
	}
	if len(query.Operations) != 1 {
		return nil, fmt.Errorf("operation file %q contains %d operations; want 1", path, len(query.Operations))
	}
	return query, nil
}

func stableVariableOrder(
	existing ast.VariableDefinitionList,
	incoming ast.VariableDefinitionList,
) ast.VariableDefinitionList {
	incomingByName := make(map[string]*ast.VariableDefinition, len(incoming))
	for _, definition := range incoming {
		incomingByName[definition.Variable] = definition
	}

	ordered := make(ast.VariableDefinitionList, 0, len(incoming))
	for _, definition := range existing {
		if current, ok := incomingByName[definition.Variable]; ok {
			ordered = append(ordered, current)
			delete(incomingByName, definition.Variable)
		}
	}

	newNames := make([]string, 0, len(incomingByName))
	for name := range incomingByName {
		newNames = append(newNames, name)
	}
	sort.Strings(newNames)
	for _, name := range newNames {
		ordered = append(ordered, incomingByName[name])
	}

	return ordered
}

func collectAliases(
	selections ast.SelectionSet,
	path []string,
) map[string]aliasRecord {
	aliases := make(map[string]aliasRecord)
	walkSelections(selections, path, func(key string, field *ast.Field) {
		if field.Alias == "" || field.Alias == field.Name {
			return
		}
		record, exists := aliases[key]
		if exists && record.alias != field.Alias {
			record.ambiguous = true
			aliases[key] = record
			return
		}
		aliases[key] = aliasRecord{alias: field.Alias}
	})
	return aliases
}

func applyAliases(
	selections ast.SelectionSet,
	path []string,
	aliases map[string]aliasRecord,
) {
	walkSelections(selections, path, func(key string, field *ast.Field) {
		record, exists := aliases[key]
		if exists && !record.ambiguous {
			field.Alias = record.alias
		}
	})
}

func applyFragmentAliases(existing, incoming *ast.QueryDocument) {
	existingByName := make(map[string]*ast.FragmentDefinition, len(existing.Fragments))
	for _, fragment := range existing.Fragments {
		existingByName[fragment.Name] = fragment
	}
	for _, fragment := range incoming.Fragments {
		previous, exists := existingByName[fragment.Name]
		if !exists {
			continue
		}
		path := []string{"fragment:" + fragment.Name, "on:" + fragment.TypeCondition}
		applyAliases(
			fragment.SelectionSet,
			path,
			collectAliases(previous.SelectionSet, path),
		)
	}
}

func walkSelections(
	selections ast.SelectionSet,
	path []string,
	visit func(string, *ast.Field),
) {
	for _, selection := range selections {
		switch current := selection.(type) {
		case *ast.Field:
			currentPath := appendPath(path, "field:"+fieldSignature(current))
			visit(strings.Join(currentPath, "/"), current)
			walkSelections(current.SelectionSet, currentPath, visit)
		case *ast.InlineFragment:
			walkSelections(
				current.SelectionSet,
				appendPath(path, "on:"+current.TypeCondition),
				visit,
			)
		}
	}
}

func fieldSignature(field *ast.Field) string {
	arguments := make([]string, 0, len(field.Arguments))
	for _, argument := range field.Arguments {
		arguments = append(arguments, argument.Name)
	}
	sort.Strings(arguments)
	return field.Name + "(" + strings.Join(arguments, ",") + ")"
}

func appendPath(path []string, element string) []string {
	result := make([]string, len(path), len(path)+1)
	copy(result, path)
	return append(result, element)
}

func formatAndValidate(
	schema *ast.Schema,
	path string,
	query *ast.QueryDocument,
	stableName string,
) ([]byte, error) {
	var output bytes.Buffer
	formatter.NewFormatter(&output).FormatQueryDocument(query)
	content := output.Bytes()

	parsed, err := parseOperation(schema, path, content)
	if err != nil {
		return nil, err
	}
	if parsed.name != stableName {
		return nil, fmt.Errorf(
			"normalized operation %q has name %q; want %q",
			path,
			parsed.name,
			stableName,
		)
	}
	return content, nil
}
