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
	preserveVariableBindings(existing, incoming)
	incomingOperation.VariableDefinitions = stableVariableOrder(
		existingOperation.VariableDefinitions,
		incomingOperation.VariableDefinitions,
	)

	aliases := collectAliases(existingOperation.SelectionSet, nil)
	applyAliases(incomingOperation.SelectionSet, nil, aliases)
	applyFragmentAliases(existing, incoming)

	return formatAndValidate(schema, path, incoming, stableName)
}

func preserveVariableBindings(existing, incoming *ast.QueryDocument) {
	existingBindings := collectVariableBindings(existing)
	incomingBindings := collectVariableBindings(incoming)
	incomingNames := make(map[string]struct{}, len(incoming.Operations[0].VariableDefinitions))
	for _, definition := range incoming.Operations[0].VariableDefinitions {
		incomingNames[definition.Variable] = struct{}{}
	}

	renames := make(map[string]string)
	claimedOldNames := make(map[string]string)
	for key, incomingName := range incomingBindings {
		existingName, exists := existingBindings[key]
		if !exists || existingName == incomingName {
			continue
		}
		if _, collision := incomingNames[existingName]; collision {
			continue
		}
		if previous, conflict := renames[incomingName]; conflict && previous != existingName {
			delete(renames, incomingName)
			continue
		}
		if previous, conflict := claimedOldNames[existingName]; conflict && previous != incomingName {
			continue
		}
		renames[incomingName] = existingName
		claimedOldNames[existingName] = incomingName
	}
	if len(renames) == 0 {
		return
	}

	for _, definition := range incoming.Operations[0].VariableDefinitions {
		if stableName, exists := renames[definition.Variable]; exists {
			definition.Variable = stableName
		}
	}
	renameVariablesInSelections(incoming.Operations[0].SelectionSet, renames)
	renameVariablesInDirectives(incoming.Operations[0].Directives, renames)
	for _, fragment := range incoming.Fragments {
		renameVariablesInSelections(fragment.SelectionSet, renames)
		renameVariablesInDirectives(fragment.Directives, renames)
	}
}

func collectVariableBindings(query *ast.QueryDocument) map[string]string {
	bindings := make(map[string]string)
	collectVariableBindingsFromSelections(query.Operations[0].SelectionSet, nil, bindings)
	for _, fragment := range query.Fragments {
		path := []string{"fragment:" + fragment.Name, "on:" + fragment.TypeCondition}
		collectVariableBindingsFromSelections(fragment.SelectionSet, path, bindings)
	}
	return bindings
}

func collectVariableBindingsFromSelections(
	selections ast.SelectionSet,
	path []string,
	bindings map[string]string,
) {
	walkSelections(selections, path, func(key string, field *ast.Field) {
		for _, argument := range field.Arguments {
			if argument.Value != nil && argument.Value.Kind == ast.Variable {
				bindings[key+"/arg:"+argument.Name] = argument.Value.Raw
			}
		}
	})
}

func renameVariablesInSelections(selections ast.SelectionSet, renames map[string]string) {
	for _, selection := range selections {
		switch current := selection.(type) {
		case *ast.Field:
			renameVariablesInArguments(current.Arguments, renames)
			renameVariablesInDirectives(current.Directives, renames)
			renameVariablesInSelections(current.SelectionSet, renames)
		case *ast.InlineFragment:
			renameVariablesInDirectives(current.Directives, renames)
			renameVariablesInSelections(current.SelectionSet, renames)
		case *ast.FragmentSpread:
			renameVariablesInDirectives(current.Directives, renames)
		}
	}
}

func renameVariablesInDirectives(directives ast.DirectiveList, renames map[string]string) {
	for _, directive := range directives {
		renameVariablesInArguments(directive.Arguments, renames)
	}
}

func renameVariablesInArguments(arguments ast.ArgumentList, renames map[string]string) {
	for _, argument := range arguments {
		renameVariablesInValue(argument.Value, renames)
	}
}

func renameVariablesInValue(value *ast.Value, renames map[string]string) {
	if value == nil {
		return
	}
	if value.Kind == ast.Variable {
		if stableName, exists := renames[value.Raw]; exists {
			value.Raw = stableName
		}
	}
	for _, child := range value.Children {
		renameVariablesInValue(child.Value, renames)
	}
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
	return field.Name
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
