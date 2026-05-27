package formatter

import (
	"fmt"
	"io"
	"slices"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

const (
	cmtQueries    = "##### Queries ##########################"
	cmtMutations  = "##### Mutations ##########################"
	cmtTypes      = "##### Types ##########################"
	cmtInputs     = "##### Inputs ##########################"
	cmtScalars    = "##### Scalars ##########################"
	cmtEnums      = "##### Enums ##########################"
	cmtInterfaces = "##### Interfaces ##########################"
	cmtUnions     = "##### Unions ##########################"
)

const (
	typeQuery        = "Query"
	typeMutation     = "Mutation"
	typeSubscription = "Subscription"
)

var ourComments = map[string]struct{}{
	cmtQueries:    {},
	cmtMutations:  {},
	cmtTypes:      {},
	cmtInputs:     {},
	cmtScalars:    {},
	cmtEnums:      {},
	cmtInterfaces: {},
	cmtUnions:     {},
}

type Formatter interface {
	FormatSchema(schema *ast.Schema) error
}

//nolint:revive // Ignore "stuttering" name formatter.FormatterOption
type FormatterOption func(*formatter)

// WithIndent uses the given string for indenting block bodies in the output,
// instead of the default, `"\t"`.
func WithIndent(indent string) FormatterOption {
	return func(f *formatter) {
		f.indent = indent
	}
}

// WithComments includes comments from the source/AST in the formatted output.
func WithComments() FormatterOption {
	return func(f *formatter) {
		f.emitComments = true
	}
}

// WithBuiltin includes builtin fields/directives/etc from the source/AST in the formatted output.
func WithBuiltin() FormatterOption {
	return func(f *formatter) {
		f.emitBuiltin = true
		f.emitIntrospection = true
	}
}

// WithNonIntrospectionBuiltin includes builtin fields/directives/etc from the
// source/AST in the formatted output, but excludes the introspection types and
// fields that starts with "__".
func WithNonIntrospectionBuiltin() FormatterOption {
	return func(f *formatter) {
		f.emitBuiltin = true
	}
}

// WithoutDescription excludes GQL description from the source/AST in the formatted output.
func WithoutDescription() FormatterOption {
	return func(f *formatter) {
		f.omitDescription = true
	}
}

// WithCompacted enables compacted output, which removes all unnecessary whitespace.
func WithCompacted() FormatterOption {
	return func(f *formatter) {
		f.compacted = true
	}
}

func NewFormatter(w io.Writer, options ...FormatterOption) Formatter {
	f := &formatter{
		indent: "\t",
		writer: w,
	}
	for _, opt := range options {
		opt(f)
	}
	return f
}

type formatter struct {
	writer io.Writer

	indent                string
	indentSize            int
	emitBuiltin           bool
	emitIntrospection     bool
	emitComments          bool
	omitDescription       bool
	printArgDescription   bool
	printFieldDescription bool
	compacted             bool

	padNext  bool
	lineHead bool

	error error
}

func (f *formatter) writeString(s string) {
	_, err := f.writer.Write([]byte(s))
	if err != nil {
		f.error = err
	}
}

func (f *formatter) writeIndent() {
	if f.lineHead {
		f.writeString(strings.Repeat(f.indent, f.indentSize))
	}
	f.lineHead = false
	f.padNext = false
}

func (f *formatter) WriteNewline() *formatter {
	f.writeString("\n")
	f.lineHead = true
	f.padNext = false

	return f
}

func (f *formatter) WriteWord(word string) *formatter {
	if f.lineHead {
		f.writeIndent()
	}
	if f.padNext {
		f.writeString(" ")
	}
	f.writeString(strings.TrimSpace(word))
	f.padNext = true

	return f
}

func (f *formatter) WriteString(s string) *formatter {
	if f.lineHead {
		f.writeIndent()
	}
	if f.padNext {
		f.writeString(" ")
	}
	f.writeString(s)
	f.padNext = false

	return f
}

func (f *formatter) WriteDescription(s string) *formatter {
	if s == "" || f.omitDescription {
		return f
	}

	f.WriteString(`"""`)
	ss := strings.Split(s, "\n")

	// single line
	if len(ss) == 1 {
		f.WriteString(" " + ss[0] + ` """`).WriteNewline()
		return f
	}

	f.WriteNewline()
	for _, s := range ss {
		f.WriteString(s).WriteNewline()
	}

	f.WriteString(`"""`).WriteNewline()

	return f
}

func (f *formatter) IncrementIndent() {
	f.indentSize++
}

func (f *formatter) DecrementIndent() {
	f.indentSize--
}

func (f *formatter) NoPadding() *formatter {
	f.padNext = false

	return f
}

func (f *formatter) NeedPadding() *formatter {
	f.padNext = true

	return f
}

func (f *formatter) FormatSchema(schema *ast.Schema) error { //nolint: gocyclo
	if schema == nil {
		return nil
	}

	f.FormatCommentGroup(schema.Comment)

	var inSchema bool
	startSchema := func() {
		if !inSchema {
			inSchema = true

			f.WriteWord("schema")

			f.FormatDirectiveList(schema.SchemaDirectives)

			f.WriteString("{").WriteNewline()
			f.IncrementIndent()
		}
	}

	needSchema := (schema.Query != nil && schema.Query.Name != typeQuery) ||
		(schema.Mutation != nil && schema.Mutation.Name != typeMutation) ||
		(schema.Subscription != nil && schema.Subscription.Name != typeSubscription)

	if needSchema && schema.Query != nil {
		startSchema()
		f.WriteWord("query").NoPadding().WriteString(":").NeedPadding()
		f.WriteWord(schema.Query.Name).WriteNewline()
	}
	if needSchema && schema.Mutation != nil {
		startSchema()
		f.WriteWord("mutation").NoPadding().WriteString(":").NeedPadding()
		f.WriteWord(schema.Mutation.Name).WriteNewline()
	}
	if needSchema && schema.Subscription != nil {
		startSchema()
		f.WriteWord("subscription").NoPadding().WriteString(":").NeedPadding()
		f.WriteWord(schema.Subscription.Name).WriteNewline()
	}
	if inSchema {
		f.DecrementIndent()
		f.WriteString("}").WriteNewline()
	} else if len(schema.SchemaDirectives) > 0 {
		// Schema definition is omitted from output, but it has
		// directives. Output them as the schema extension to not loose
		// them
		f.WriteWord("extend").WriteWord("schema")

		f.FormatDirectiveList(schema.SchemaDirectives)

		f.WriteNewline()
	}

	directiveNames := make([]string, 0, len(schema.Directives))
	for name := range schema.Directives {
		directiveNames = append(directiveNames, name)
	}
	sort.Strings(directiveNames)
	for _, name := range directiveNames {
		f.FormatDirectiveDefinition(schema.Directives[name])
	}

	f.formatTypes(schema.Types)
	return f.error
}

func (f *formatter) formatTypes(types map[string]*ast.Definition) { //nolint: gocyclo,funlen
	var scalars, objects, interfaces, unions, enums, inputs, others []string
	for name, def := range types {
		if !f.emitIntrospection && strings.HasPrefix(name, "__") {
			continue
		}
		if !f.emitBuiltin && def.BuiltIn {
			continue
		}

		switch def.Kind {
		case ast.Scalar:
			scalars = append(scalars, name)
		case ast.Object:
			if name != typeQuery && name != typeMutation {
				objects = append(objects, name)
			}
		case ast.Interface:
			interfaces = append(interfaces, name)
		case ast.Union:
			unions = append(unions, name)
		case ast.Enum:
			enums = append(enums, name)
		case ast.InputObject:
			inputs = append(inputs, name)
		default:
			others = append(others, name)
		}
	}

	if query, ok := types["Query"]; ok {
		f.FormatSection(cmtQueries)
		f.FormatDefinition(query, false)
	}

	if mutation, ok := types["Mutation"]; ok {
		f.FormatSection(cmtMutations)
		f.FormatDefinition(mutation, false)
	}

	// types
	if len(objects) > 0 {
		slices.Sort(objects)
		f.FormatSection(cmtTypes)
		for _, name := range objects {
			f.FormatDefinition(types[name], false)
		}
	}

	// inputs
	if len(inputs) > 0 {
		slices.Sort(inputs)
		f.FormatSection(cmtInputs)
		for _, name := range inputs {
			f.FormatDefinition(types[name], false)
		}
	}

	// scalars
	if len(scalars) > 0 {
		slices.Sort(scalars)
		f.FormatSection(cmtScalars)
		for _, name := range scalars {
			f.FormatDefinition(types[name], false)
		}
	}

	// enums
	if len(enums) > 0 {
		slices.Sort(enums)
		f.FormatSection(cmtEnums)
		for _, name := range enums {
			f.FormatDefinition(types[name], false)
		}
	}

	// interfaces
	if len(interfaces) > 0 {
		slices.Sort(interfaces)
		f.FormatSection(cmtInterfaces)
		for _, name := range interfaces {
			f.FormatDefinition(types[name], false)
		}
	}

	// unions
	if len(unions) > 0 {
		slices.Sort(unions)
		f.FormatSection(cmtUnions)
		for _, name := range unions {
			f.FormatDefinition(types[name], false)
		}
	}

	// others
	if len(others) > 0 {
		slices.Sort(others)
		for _, name := range others {
			f.FormatDefinition(types[name], false)
		}
	}
}

func (f *formatter) FormatSection(name string) {
	f.WriteString(name).WriteNewline()
}

func (f *formatter) FormatFieldList(fieldList ast.FieldList, endOfDefComment *ast.CommentGroup) {
	if len(fieldList) == 0 {
		return
	}

	f.WriteString("{").WriteNewline()
	f.IncrementIndent()

	for _, field := range fieldList {
		f.FormatFieldDefinition(field)
	}

	f.FormatCommentGroup(endOfDefComment)

	f.DecrementIndent()
	f.WriteString("}")
}

func (f *formatter) FormatFieldDefinition(field *ast.FieldDefinition) {
	if !f.emitIntrospection && strings.HasPrefix(field.Name, "__") {
		return
	}
	if !f.emitBuiltin &&
		(field.Position != nil && field.Position.Src != nil && field.Position.Src.BuiltIn) {
		return
	}

	f.FormatCommentGroup(field.BeforeDescriptionComment)

	if f.printFieldDescription {
		f.WriteDescription(field.Description)
	}

	f.FormatCommentGroup(field.AfterDescriptionComment)

	f.WriteWord(field.Name).NoPadding()
	f.FormatArgumentDefinitionList(field.Arguments)
	f.NoPadding().WriteString(":").NeedPadding()
	f.FormatType(field.Type)

	if field.DefaultValue != nil {
		f.WriteWord("=")
		f.FormatValue(field.DefaultValue)
	}

	f.FormatDirectiveList(field.Directives)

	f.WriteNewline()
}

func (f *formatter) FormatArgumentDefinitionList(lists ast.ArgumentDefinitionList) {
	if len(lists) == 0 {
		return
	}

	f.WriteString("(")
	for idx, arg := range lists {
		f.FormatArgumentDefinition(arg)

		// Skip emitting (insignificant) comma in case it is the
		// last argument, or we printed a new line in its definition.
		if idx != len(lists)-1 { // if descriptions are turned on, use arg.Description == ""
			f.NoPadding().WriteWord(",")
		}
	}
	f.NoPadding().WriteString(")").NeedPadding()
}

func (f *formatter) FormatArgumentDefinition(def *ast.ArgumentDefinition) {
	f.FormatCommentGroup(def.BeforeDescriptionComment)

	if def.Description != "" && !f.omitDescription && f.printArgDescription {
		f.IncrementIndent()
		f.WriteNewline().IncrementIndent()
		f.WriteDescription(def.Description)
	}

	f.FormatCommentGroup(def.AfterDescriptionComment)

	f.WriteWord(def.Name).NoPadding().WriteString(":").NeedPadding()
	f.FormatType(def.Type)

	if def.DefaultValue != nil {
		f.WriteWord("=")
		f.FormatValue(def.DefaultValue)
	}

	f.NeedPadding().FormatDirectiveList(def.Directives)

	if def.Description != "" && !f.omitDescription && f.printArgDescription {
		f.DecrementIndent()
		f.WriteNewline()
	}
}

func (f *formatter) FormatDirectiveLocation(location ast.DirectiveLocation) {
	f.WriteWord(string(location))
}

func (f *formatter) FormatDirectiveDefinition(def *ast.DirectiveDefinition) {
	if !f.emitBuiltin &&
		(def.Position != nil && def.Position.Src != nil && def.Position.Src.BuiltIn) {
		return
	}

	f.FormatCommentGroup(def.BeforeDescriptionComment)

	f.WriteDescription(def.Description)

	f.FormatCommentGroup(def.AfterDescriptionComment)

	f.WriteWord("directive").WriteString("@").WriteWord(def.Name)

	if len(def.Arguments) != 0 {
		f.NoPadding()
		f.FormatArgumentDefinitionList(def.Arguments)
	}

	if def.IsRepeatable {
		f.WriteWord("repeatable")
	}

	if len(def.Locations) != 0 {
		f.WriteWord("on")

		for idx, dirLoc := range def.Locations {
			f.FormatDirectiveLocation(dirLoc)

			if idx != len(def.Locations)-1 {
				f.WriteWord("|")
			}
		}
	}

	f.WriteNewline()
	if !f.compacted {
		f.WriteNewline()
	}
}

func (f *formatter) FormatDefinition(def *ast.Definition, extend bool) {
	if !f.emitIntrospection && strings.HasPrefix(def.Name, "__") {
		return
	}
	if !f.emitBuiltin && def.BuiltIn {
		return
	}

	f.FormatCommentGroup(def.BeforeDescriptionComment)

	f.WriteDescription(def.Description)

	f.FormatCommentGroup(def.AfterDescriptionComment)

	if extend {
		f.WriteWord("extend")
	}

	switch def.Kind {
	case ast.Scalar:
		f.WriteWord("scalar").WriteWord(def.Name)

	case ast.Object:
		f.WriteWord("type").WriteWord(def.Name)

	case ast.Interface:
		f.WriteWord("interface").WriteWord(def.Name)

	case ast.Union:
		f.WriteWord("union").WriteWord(def.Name)

	case ast.Enum:
		f.WriteWord("enum").WriteWord(def.Name)

	case ast.InputObject:
		f.WriteWord("input").WriteWord(def.Name)
	}

	if len(def.Interfaces) != 0 {
		slices.Sort(def.Interfaces)
		f.WriteWord("implements").WriteWord(strings.Join(def.Interfaces, " & "))
	}

	f.FormatDirectiveList(def.Directives)

	if len(def.Types) != 0 {
		f.WriteWord("=").WriteWord(strings.Join(def.Types, " | "))
	}

	f.FormatFieldList(def.Fields, def.EndOfDefinitionComment)

	f.FormatEnumValueList(def.EnumValues, def.EndOfDefinitionComment)

	f.WriteNewline()
	if !f.compacted {
		f.WriteNewline()
	}
}

func (f *formatter) FormatEnumValueList(
	lists ast.EnumValueList,
	endOfDefComment *ast.CommentGroup,
) {
	if len(lists) == 0 {
		return
	}

	f.WriteString("{").WriteNewline()
	f.IncrementIndent()

	for _, v := range lists {
		f.FormatEnumValueDefinition(v)
	}

	f.FormatCommentGroup(endOfDefComment)

	f.DecrementIndent()
	f.WriteString("}")
}

func (f *formatter) FormatEnumValueDefinition(def *ast.EnumValueDefinition) {
	f.FormatCommentGroup(def.BeforeDescriptionComment)

	f.WriteDescription(def.Description)

	f.FormatCommentGroup(def.AfterDescriptionComment)

	f.WriteWord(def.Name)
	f.FormatDirectiveList(def.Directives)

	f.WriteNewline()
}

func (f *formatter) FormatDirectiveList(lists ast.DirectiveList) {
	if len(lists) == 0 {
		return
	}

	for _, dir := range lists {
		f.FormatDirective(dir)
	}
}

func (f *formatter) FormatDirective(dir *ast.Directive) {
	f.WriteString("@").WriteWord(dir.Name)
	f.FormatArgumentList(dir.Arguments)
}

func (f *formatter) FormatArgumentList(lists ast.ArgumentList) {
	if len(lists) == 0 {
		return
	}
	f.NoPadding().WriteString("(")
	for idx, arg := range lists {
		f.FormatArgument(arg)

		if idx != len(lists)-1 {
			f.NoPadding().WriteWord(",")
		}
	}
	f.WriteString(")").NeedPadding()
}

func (f *formatter) FormatArgument(arg *ast.Argument) {
	f.FormatCommentGroup(arg.Comment)

	f.WriteWord(arg.Name).NoPadding().WriteString(":").NeedPadding()
	f.WriteString(arg.Value.String())
}

func (f *formatter) FormatFragmentDefinitionList(lists ast.FragmentDefinitionList) {
	for _, def := range lists {
		f.FormatFragmentDefinition(def)
	}
}

func (f *formatter) FormatFragmentDefinition(def *ast.FragmentDefinition) {
	f.FormatCommentGroup(def.Comment)

	f.WriteWord("fragment").WriteWord(def.Name)
	f.FormatVariableDefinitionList(def.VariableDefinition)
	f.WriteWord("on").WriteWord(def.TypeCondition)
	f.FormatDirectiveList(def.Directives)

	if len(def.SelectionSet) != 0 {
		f.FormatSelectionSet(def.SelectionSet)
		f.WriteNewline()
	}
}

func (f *formatter) FormatVariableDefinitionList(lists ast.VariableDefinitionList) {
	if len(lists) == 0 {
		return
	}

	f.WriteString("(")
	for idx, def := range lists {
		f.FormatVariableDefinition(def)

		if idx != len(lists)-1 {
			f.NoPadding().WriteWord(",")
		}
	}
	f.NoPadding().WriteString(")").NeedPadding()
}

func (f *formatter) FormatVariableDefinition(def *ast.VariableDefinition) {
	f.FormatCommentGroup(def.Comment)

	f.WriteString("$").WriteWord(def.Variable).NoPadding().WriteString(":").NeedPadding()
	f.FormatType(def.Type)

	if def.DefaultValue != nil {
		f.WriteWord("=")
		f.FormatValue(def.DefaultValue)
	}

	// TODO https://github.com/vektah/gqlparser/v2/issues/102
	//   VariableDefinition : Variable : Type DefaultValue? Directives[Const]?
}

func (f *formatter) FormatSelectionSet(sets ast.SelectionSet) {
	if len(sets) == 0 {
		return
	}

	f.WriteString("{").WriteNewline()
	f.IncrementIndent()

	for _, sel := range sets {
		f.FormatSelection(sel)
	}

	f.DecrementIndent()
	f.WriteString("}")
}

func (f *formatter) FormatSelection(selection ast.Selection) {
	switch v := selection.(type) {
	case *ast.Field:
		f.FormatField(v)

	case *ast.FragmentSpread:
		f.FormatFragmentSpread(v)

	case *ast.InlineFragment:
		f.FormatInlineFragment(v)

	default:
		panic(fmt.Errorf("unknown Selection type: %T", selection))
	}

	f.WriteNewline()
}

func (f *formatter) FormatField(field *ast.Field) {
	f.FormatCommentGroup(field.Comment)

	if field.Alias != "" && field.Alias != field.Name {
		f.WriteWord(field.Alias).NoPadding().WriteString(":").NeedPadding()
	}
	f.WriteWord(field.Name)

	if len(field.Arguments) != 0 {
		f.NoPadding()
		f.FormatArgumentList(field.Arguments)
		f.NeedPadding()
	}

	f.FormatDirectiveList(field.Directives)

	f.FormatSelectionSet(field.SelectionSet)
}

func (f *formatter) FormatFragmentSpread(spread *ast.FragmentSpread) {
	f.FormatCommentGroup(spread.Comment)

	f.WriteWord("...")
	if f.compacted {
		f.NoPadding()
	}
	f.WriteWord(spread.Name)

	f.FormatDirectiveList(spread.Directives)
}

func (f *formatter) FormatInlineFragment(inline *ast.InlineFragment) {
	f.FormatCommentGroup(inline.Comment)

	f.WriteWord("...")
	if inline.TypeCondition != "" {
		f.WriteWord("on").WriteWord(inline.TypeCondition)
	}

	f.FormatDirectiveList(inline.Directives)

	f.FormatSelectionSet(inline.SelectionSet)
}

func (f *formatter) FormatType(t *ast.Type) {
	f.WriteWord(t.String())
}

func (f *formatter) FormatValue(value *ast.Value) {
	f.FormatCommentGroup(value.Comment)

	f.WriteString(value.String())
}

func (f *formatter) FormatCommentGroup(group *ast.CommentGroup) {
	if !f.emitComments || group == nil {
		return
	}
	for _, comment := range group.List {
		if _, ok := ourComments["#"+comment.Text()]; ok {
			continue
		}
		f.FormatComment(comment)
	}
}

func (f *formatter) FormatComment(comment *ast.Comment) {
	if !f.emitComments || comment == nil {
		return
	}
	f.WriteString("#").WriteString(comment.Text()).WriteNewline()
}
