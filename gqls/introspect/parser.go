//nolint:goconst
package introspect

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/introspection"
)

type introspectionResponse struct {
	Data introspection.Data `json:"data"`
}

const (
	descrLineLength = 101
)

var builtInDirectives = []string{"skip", "include", "deprecated"}
var builtInScalars = []string{"Int", "Float", "String", "Boolean", "ID"}

// ParseSchema parses a graphql introspection JSON
func ParseSchema(contents []byte, fName string) (s *ast.Schema, err error) {
	var schema ast.Schema

	intro := introspectionResponse{}
	if err := json.Unmarshal(contents, &intro); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON introspection from '%s': %w", fName, err)
	}

	if schema.Directives, err = parseDirectives(intro.Data.Schema.Directives); err != nil {
		return nil, err
	}
	if schema.Types, err = parseTypes(intro.Data.Schema.Types); err != nil {
		return nil, err
	}

	return &schema, nil
}

func parseDirectives(directives []introspection.Directive) (map[string]*ast.DirectiveDefinition, error) {
	if len(directives) == 0 {
		return nil, nil
	}

	defs := make(map[string]*ast.DirectiveDefinition)
	for _, dir := range directives {
		if slices.Contains(builtInDirectives, dir.Name) {
			continue
		}
		def := ast.DirectiveDefinition{
			Name:        dir.Name,
			Description: wrapText(dir.Description, descrLineLength),
			Locations:   directiveLocations(dir.Locations),
		}
		var err error
		if def.Arguments, err = directiveArguments(dir.Args); err != nil {
			return nil, err
		}
		defs[def.Name] = &def
	}

	if len(defs) == 0 {
		return nil, nil
	}

	return defs, nil
}

func parseTypes(types []*introspection.FullType) (map[string]*ast.Definition, error) {
	if len(types) == 0 {
		return nil, nil
	}

	defs := make(map[string]*ast.Definition)
	for _, t := range types {
		def := ast.Definition{
			Kind:        ast.DefinitionKind(t.Kind.String()),
			Description: wrapText(t.Description, descrLineLength),
			Name:        t.Name,
			Interfaces:  getInterfaces(t.Interfaces),
		}
		var err error
		switch t.Kind {
		case introspection.SCALAR:
			if slices.Contains(builtInScalars, t.Name) {
				continue
			}
		case introspection.OBJECT, introspection.INTERFACE:
			if def.Fields, err = typeFields(t.Fields); err != nil {
				return nil, err
			}
		case introspection.INPUTOBJECT:
			if def.Fields, err = typeInputFields(t.InputFields); err != nil {
				return nil, err
			}
		case introspection.ENUM:
			def.EnumValues = enumValues(t.EnumValues)
		case introspection.UNION:
			if def.Types, err = unionTypes(t.PossibleTypes); err != nil {
				return nil, err
			}
		}
		defs[def.Name] = &def
	}

	if len(defs) == 0 {
		return nil, nil
	}

	return defs, nil
}

func unionTypes(trs []introspection.TypeRef) (uts []string, err error) {
	for _, t := range trs {
		if t.Name == nil {
			return nil, fmt.Errorf("name is nil on union type %s", t.Kind)
		}
		uts = append(uts, *t.Name)
	}
	return uts, nil
}

func enumValues(evs []introspection.EnumValue) (enums []*ast.EnumValueDefinition) {
	for _, ev := range evs {
		astEnum := ast.EnumValueDefinition{
			Description: wrapText(ev.Description, descrLineLength),
			Name:        ev.Name,
		}
		addDepricated(ev.IsDeprecated, ev.DeprecationReason, &astEnum.Directives)
		enums = append(enums, &astEnum)
	}
	return enums
}

func getInterfaces(is []introspection.TypeRef) (interfaces []string) {
	for _, i := range is {
		if i.Name != nil {
			interfaces = append(interfaces, *i.Name)
		}
	}
	return interfaces
}

func typeFields(fields []introspection.Field) (dirs []*ast.FieldDefinition, err error) {
	for _, field := range fields {
		fd := ast.FieldDefinition{
			Description: wrapText(field.Description, descrLineLength),
			Name:        field.Name,
			// Arguments    ArgumentDefinitionList // only for objects
			// DefaultValue *Value                 // only for input objects
			// Directives   DirectiveList
		}
		if fd.Type, err = getType(field.Type); err != nil {
			return nil, err
		}
		if fd.Arguments, err = getArguments(field.Args); err != nil {
			return nil, err
		}
		addDepricated(field.IsDeprecated, field.DeprecationReason, &fd.Directives)
		dirs = append(dirs, &fd)
	}
	return dirs, nil
}

func typeInputFields(ivs []introspection.InputValue) (dirs []*ast.FieldDefinition, err error) {
	for _, iv := range ivs {
		fd := ast.FieldDefinition{
			Description: wrapText(iv.Description, descrLineLength),
			Name:        iv.Name,
			// DefaultValue *Value                 // only for input objects
			// Directives   DirectiveList
		}
		if fd.Type, err = getType(iv.Type); err != nil {
			return nil, err
		}
		if fd.DefaultValue, err = getDefaultValue(iv.DefaultValue, iv.Type); err != nil {
			return nil, err
		}
		addDepricated(iv.IsDeprecated, iv.DeprecationReason, &fd.Directives)
		dirs = append(dirs, &fd)
	}
	return dirs, nil
}

func getArguments(introArgs []introspection.InputValue) (args []*ast.ArgumentDefinition, err error) {
	for _, a := range introArgs {
		astArg := ast.ArgumentDefinition{
			Description: wrapText(a.Description, descrLineLength),
			Name:        a.Name,
			// Directives   DirectiveList

		}
		if astArg.Type, err = getType(a.Type); err != nil {
			return nil, err
		}
		if astArg.DefaultValue, err = getDefaultValue(a.DefaultValue, a.Type); err != nil {
			return nil, err
		}
		addDepricated(a.IsDeprecated, a.DeprecationReason, &astArg.Directives)
		args = append(args, &astArg)
	}
	return args, nil
}

func getDefaultValue(defVal *string, introType introspection.TypeRef) (*ast.Value, error) {
	if defVal == nil {
		return nil, nil
	}
	if introType.Kind == introspection.NONNULL {
		if introType.OfType == nil {
			return nil, fmt.Errorf("failed to parse element type: NonNull does not have a subtype")
		}
		introType = *introType.OfType
	}

	valKind := ast.EnumValue
	switch introType.Kind {
	case introspection.SCALAR:
		if introType.Name != nil {
			switch *introType.Name {
			case "String":
				valKind = ast.StringValue
			case "Int":
				valKind = ast.IntValue
			case "Float":
				valKind = ast.FloatValue
			case "Boolean":
				valKind = ast.BooleanValue
			case "ID":
				valKind = ast.EnumValue
			}
		}
	case introspection.LIST:
		valKind = ast.ListValue
	case introspection.ENUM:
		valKind = ast.EnumValue
	}
	astVal := ast.Value{
		Raw:  *defVal,
		Kind: valKind,
	}
	if valKind == ast.StringValue {
		astVal.Raw = strings.Trim(*defVal, `"`)
	}
	return &astVal, nil
}

func addDepricated(isDeprecated bool, pReason *string, directives *ast.DirectiveList) {
	if !isDeprecated {
		return
	}
	directive := &ast.Directive{Name: "deprecated"}

	if pReason != nil && *pReason == "No longer supported" {
		pReason = nil
	}
	if pReason != nil {
		directive.Arguments = []*ast.Argument{{Name: "reason", Value: &ast.Value{Raw: *pReason, Kind: ast.StringValue}}}
	}
	*directives = append(*directives, directive)
}

func getType(introType introspection.TypeRef) (*ast.Type, error) {
	var astType ast.Type

	if introType.Kind == introspection.NONNULL { // marks the subelement as non-null
		astType.NonNull = true
		if introType.OfType == nil {
			return nil, fmt.Errorf("failed to parse element type: NonNull does not have a subtype")
		}
		introType = *introType.OfType
	}

	// Name=nil, OfType should have something
	if introType.Name == nil {
		if introType.OfType == nil {
			return nil, fmt.Errorf("failed to parse element type: Name is nil and does not have a subtype")
		}
		t, err := getType(*introType.OfType)
		if err != nil {
			return nil, err
		}
		astType.Elem = t
		return &astType, nil
	}

	astType.NamedType = *introType.Name
	return &astType, nil
}

func directiveArguments(args []introspection.InputValue) (defs []*ast.ArgumentDefinition, err error) {
	if len(args) == 0 {
		return nil, nil
	}
	for _, a := range args {
		def := ast.ArgumentDefinition{
			Description: wrapText(a.Description, descrLineLength),
			Name:        a.Name,
		}
		if def.Type, err = getType(a.Type); err != nil {
			return nil, err
		}
		if def.DefaultValue, err = getDefaultValue(a.DefaultValue, a.Type); err != nil {
			return nil, err
		}
		defs = append(defs, &def)
	}
	return defs, nil
}

func directiveLocations(locs []string) (dLocs []ast.DirectiveLocation) {
	for _, loc := range locs {
		dLocs = append(dLocs, ast.DirectiveLocation(loc))
	}
	return dLocs
}

func wrapText(s string, maxLen int) string {
	var lines []string

	for _, line := range strings.Split(s, "\n") {
		var lastSpace int
		if line == "" {
			lines = append(lines, line)
			continue
		}

		for i := 0; i < len(line); i++ {
			if i+1 == len(line) {
				lines = append(lines, line)
				break
			}
			if line[i] == ' ' {
				lastSpace = i
			}
			if i+1 >= maxLen {
				if lastSpace == 0 {
					lastSpace = i
				}
				l := line[:lastSpace+1]
				for l != "" {
					if l[len(l)-1] != ' ' {
						break
					}
					l = l[:len(l)-1]
				}
				lines = append(lines, l)
				line = line[lastSpace+1:]
				lastSpace = 0
				i = -1
			}
		}
	}
	return strings.Join(lines, "\n")
}
