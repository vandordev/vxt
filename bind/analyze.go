package bind

import (
	"errors"

	"github.com/vandordev/vxt/runtime"
)

type analyzedDocument struct {
	PackageName string
	Template    string
	Types       []analyzedType
	InputFields []analyzedField
}

type analyzedType struct {
	Name   string
	Fields []analyzedField
}

type analyzedField struct {
	GoName     string
	SchemaName string
	GoType     string
	TypeName   string
	Optional   bool
	Array      bool
	Primitive  bool
}

func analyzeDocument(packageName string, doc *runtime.CompiledDocument) (analyzedDocument, error) {
	if doc == nil {
		return analyzedDocument{}, errors.New("bind: missing compiled document")
	}

	result := analyzedDocument{
		PackageName: packageName,
		Template:    doc.Template,
	}

	for _, typ := range doc.Types {
		analyzed := analyzedType{Name: typ.Name}
		for _, field := range typ.Fields {
			analyzed.Fields = append(analyzed.Fields, analyzedField{
				GoName:     toExportedGoName(field.Name),
				SchemaName: field.Name,
				GoType:     goTypeForField(field),
				TypeName:   field.TypeName,
				Optional:   field.Optional,
				Array:      field.Array,
				Primitive:  isPrimitiveType(field.TypeName),
			})
		}
		result.Types = append(result.Types, analyzed)
	}

	for _, input := range doc.Inputs {
		result.InputFields = append(result.InputFields, analyzedField{
			GoName:     toExportedGoName(input.Name),
			SchemaName: input.Name,
			GoType:     goTypeForInput(input, doc.Types),
			TypeName:   input.TypeName,
			Primitive:  isPrimitiveType(input.TypeName),
		})
	}

	return result, nil
}
