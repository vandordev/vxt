package bind

import (
	"fmt"
	"go/format"
	"sort"
	"strconv"
	"strings"
)

func emitGeneratedFile(doc analyzedDocument, assets embeddedAssets) (string, error) {
	var b strings.Builder

	writeLine(&b, "package %s", doc.PackageName)
	writeLine(&b, "")
	writeLine(&b, "import (")
	writeLine(&b, "\t\"fmt\"")
	writeLine(&b, "")
	writeLine(&b, "\t\"github.com/vandordev/vxt/runtime\"")
	writeLine(&b, "\t\"github.com/vandordev/vxt/source\"")
	writeLine(&b, "\t\"github.com/vandordev/vxt/write\"")
	writeLine(&b, ")")
	writeLine(&b, "")

	for _, typ := range doc.Types {
		writeLine(&b, "type %s struct {", typ.Name)
		for _, field := range typ.Fields {
			writeLine(&b, "\t%s %s", field.GoName, field.GoType)
		}
		writeLine(&b, "}")
		writeLine(&b, "")
	}

	writeLine(&b, "type Input struct {")
	for _, field := range doc.InputFields {
		writeLine(&b, "\t%s %s", field.GoName, field.GoType)
	}
	writeLine(&b, "}")
	writeLine(&b, "")

	writeLine(&b, "var documentSource = source.Source{")
	writeLine(&b, "\tID: %s,", strconv.Quote(assets.Main.ID))
	writeLine(&b, "\tPath: %s,", strconv.Quote(assets.Main.Path))
	writeLine(&b, "\tText: %s,", strconv.Quote(assets.Main.Text))
	writeLine(&b, "}")
	writeLine(&b, "")

	if len(assets.Uses) == 0 {
		writeLine(&b, "var useSources runtime.MapResolver")
		writeLine(&b, "")
	} else {
		keys := make([]string, 0, len(assets.Uses))
		for k := range assets.Uses {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		writeLine(&b, "var useSources = runtime.MapResolver{")
		for _, key := range keys {
			src := assets.Uses[key]
			writeLine(&b, "\t%s: {ID: %s, Path: %s, Text: %s},", strconv.Quote(key), strconv.Quote(src.ID), strconv.Quote(src.Path), strconv.Quote(src.Text))
		}
		writeLine(&b, "}")
		writeLine(&b, "")
	}

	for _, typ := range doc.Types {
		emitTypeConverter(&b, typ)
	}
	emitInputConverter(&b, doc)
	emitWrappers(&b)

	formatted, err := format.Source([]byte(b.String()))
	if err != nil {
		return "", fmt.Errorf("bind: format generated file: %w", err)
	}
	return string(formatted), nil
}

func emitTypeConverter(b *strings.Builder, typ analyzedType) {
	receiver := strings.ToLower(string(typ.Name[0]))
	writeLine(b, "func %sToRuntimeValue(%s %s) map[string]any {", receiver, receiver, typ.Name)
	writeLine(b, "\tout := map[string]any{}")
	for _, field := range typ.Fields {
		access := receiver + "." + field.GoName
		switch {
		case field.Array && field.Primitive:
			writeLine(b, "\tout[%s] = append([]%s(nil), %s...)", strconv.Quote(field.SchemaName), goTypeName(field.TypeName), access)
		case field.Array && !field.Primitive:
			writeLine(b, "\titems%s := make([]map[string]any, 0, len(%s))", field.GoName, access)
			writeLine(b, "\tfor _, item := range %s {", access)
			writeLine(b, "\t\titems%s = append(items%s, %sToRuntimeValue(item))", field.GoName, field.GoName, lowerFirst(field.TypeName))
			writeLine(b, "\t}")
			writeLine(b, "\tout[%s] = items%s", strconv.Quote(field.SchemaName), field.GoName)
		case field.Optional && field.Primitive:
			writeLine(b, "\tif %s != nil {", access)
			writeLine(b, "\t\tout[%s] = *%s", strconv.Quote(field.SchemaName), access)
			writeLine(b, "\t}")
		case field.Optional && !field.Primitive:
			writeLine(b, "\tif %s != nil {", access)
			writeLine(b, "\t\tout[%s] = %sToRuntimeValue(*%s)", strconv.Quote(field.SchemaName), lowerFirst(field.TypeName), access)
			writeLine(b, "\t}")
		case !field.Primitive:
			writeLine(b, "\tout[%s] = %sToRuntimeValue(%s)", strconv.Quote(field.SchemaName), lowerFirst(field.TypeName), access)
		default:
			writeLine(b, "\tout[%s] = %s", strconv.Quote(field.SchemaName), access)
		}
	}
	writeLine(b, "\treturn out")
	writeLine(b, "}")
	writeLine(b, "")
}

func emitInputConverter(b *strings.Builder, doc analyzedDocument) {
	writeLine(b, "func (in Input) toRuntimeInput() map[string]any {")
	writeLine(b, "\tout := map[string]any{}")
	for _, field := range doc.InputFields {
		access := "in." + field.GoName
		if field.Primitive {
			writeLine(b, "\tout[%s] = %s", strconv.Quote(field.SchemaName), access)
			continue
		}
		writeLine(b, "\tout[%s] = %sToRuntimeValue(%s)", strconv.Quote(field.SchemaName), lowerFirst(field.TypeName), access)
	}
	writeLine(b, "\treturn out")
	writeLine(b, "}")
	writeLine(b, "")
}

func emitWrappers(b *strings.Builder) {
	writeLine(b, "func Compile() runtime.CompileResult {")
	writeLine(b, "\tif useSources != nil {")
	writeLine(b, "\t\treturn runtime.CompileDocumentWithResolver(documentSource, useSources)")
	writeLine(b, "\t}")
	writeLine(b, "\treturn runtime.CompileDocument(documentSource)")
	writeLine(b, "}")
	writeLine(b, "")

	writeLine(b, "func Validate(input Input) runtime.ValidationResult {")
	writeLine(b, "\tcompiled := Compile()")
	writeLine(b, "\tif len(compiled.Diagnostics) > 0 || compiled.Document == nil {")
	writeLine(b, "\t\treturn runtime.ValidationResult{")
	writeLine(b, "\t\t\tDocument:    compiled.Document,")
	writeLine(b, "\t\t\tInput:       input.toRuntimeInput(),")
	writeLine(b, "\t\t\tDiagnostics: compiled.Diagnostics,")
	writeLine(b, "\t\t}")
	writeLine(b, "\t}")
	writeLine(b, "\treturn runtime.ValidateDocument(compiled.Document, input.toRuntimeInput())")
	writeLine(b, "}")
	writeLine(b, "")

	writeLine(b, "func PlanDetailed(input Input) runtime.PlanResult {")
	writeLine(b, "\treturn runtime.PlanDocument(Validate(input))")
	writeLine(b, "}")
	writeLine(b, "")

	writeLine(b, "func Plan(input Input) (runtime.Plan, error) {")
	writeLine(b, "\tresult := PlanDetailed(input)")
	writeLine(b, "\tif len(result.Diagnostics) > 0 {")
	writeLine(b, "\t\treturn runtime.Plan{}, fmt.Errorf(\"%%s\", result.Diagnostics[0].Message)")
	writeLine(b, "\t}")
	writeLine(b, "\treturn result.Plan, nil")
	writeLine(b, "}")
	writeLine(b, "")

	writeLine(b, "func WriteDetailed(input Input, target write.OutputTarget) runtime.WriteResult {")
	writeLine(b, "\tplanned := PlanDetailed(input)")
	writeLine(b, "\tif len(planned.Diagnostics) > 0 {")
	writeLine(b, "\t\treturn runtime.WriteResult{")
	writeLine(b, "\t\t\tDiagnostics: planned.Diagnostics,")
			writeLine(b, "\t\t\tErr:         fmt.Errorf(\"%%s\", planned.Diagnostics[0].Message),")
	writeLine(b, "\t\t}")
	writeLine(b, "\t}")
	writeLine(b, "\treturn runtime.WritePlanWithDiagnostics(planned.Plan, target)")
	writeLine(b, "}")
	writeLine(b, "")

	writeLine(b, "func Write(input Input, target write.OutputTarget) (write.WriteReport, error) {")
	writeLine(b, "\tresult := WriteDetailed(input, target)")
	writeLine(b, "\treturn result.Report, result.Err")
	writeLine(b, "}")
	writeLine(b, "")
}

func writeLine(b *strings.Builder, format string, args ...any) {
	if len(args) == 0 {
		b.WriteString(format)
		b.WriteByte('\n')
		return
	}
	b.WriteString(fmt.Sprintf(format, args...))
	b.WriteByte('\n')
}

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}
