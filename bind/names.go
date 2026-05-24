package bind

import (
	"strings"
	"unicode"

	"github.com/vandordev/vxt/runtime"
)

func toExportedGoName(name string) string {
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == '-' || unicode.IsSpace(r)
	})
	if len(parts) == 0 {
		return ""
	}

	var b strings.Builder
	for _, part := range parts {
		if part == "" {
			continue
		}
		rs := []rune(part)
		b.WriteRune(unicode.ToUpper(rs[0]))
		for _, r := range rs[1:] {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func goTypeForField(field runtime.TypeFieldDecl) string {
	base := goTypeName(field.TypeName)
	if field.Array {
		base = "[]" + base
	}
	if field.Optional {
		base = "*" + base
	}
	return base
}

func goTypeForInput(input runtime.InputDecl, _ []runtime.TypeDecl) string {
	return goTypeName(input.TypeName)
}

func goTypeName(typeName string) string {
	switch typeName {
	case "string":
		return "string"
	case "bool":
		return "bool"
	case "int":
		return "int"
	case "float":
		return "float64"
	default:
		return typeName
	}
}

func isPrimitiveType(typeName string) bool {
	switch typeName {
	case "string", "bool", "int", "float":
		return true
	default:
		return false
	}
}
