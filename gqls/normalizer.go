package gqls

import (
	"cmp"
	"slices"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

func Normalize(schema *ast.Schema) {
	if schema == nil {
		return
	}

	for _, v := range schema.Directives {
		normalizeDirectives(v)
	}
	for _, v := range schema.Types {
		normalizeDefinition(v)
	}
}

func normalizeDirectives(def *ast.DirectiveDefinition) {
	if def == nil {
		return
	}
	def.Description = strings.Trim(def.Description, " ")
	slices.Sort(def.Locations)
}

func normalizeDefinition(def *ast.Definition) {
	if def == nil {
		return
	}
	def.Description = strings.Trim(def.Description, " ")

	// directives
	if def.Directives != nil {
		slices.SortFunc(def.Directives, func(a, b *ast.Directive) int {
			return cmp.Compare(a.Name, b.Name)
		})
	}

	// interfaces
	if def.Interfaces != nil {
		slices.Sort(def.Interfaces)
	}

	// fields
	if def.Fields != nil {
		slices.SortFunc(def.Fields, func(a, b *ast.FieldDefinition) int {
			return cmp.Compare(a.Name, b.Name)
		})
		for _, f := range def.Fields {
			f.Description = strings.Trim(f.Description, " ")
		}
	}

	// types
	if def.Types != nil {
		slices.Sort(def.Types)
	}

	// enums
	if def.EnumValues != nil {
		slices.SortFunc(def.EnumValues, func(a, b *ast.EnumValueDefinition) int {
			return cmp.Compare(a.Name, b.Name)
		})
		for _, f := range def.EnumValues {
			f.Description = strings.Trim(f.Description, " ")
		}
	}
}
