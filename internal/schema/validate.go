package schema

import (
	"fmt"
	"reflect"

	"github.com/vandordev/vxt/internal/model"
)

func ValidateValue(typeName string, value any) error {
	switch typeName {
	case TypeString:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected string")
		}
	case TypeNumber:
		switch value.(type) {
		case int, int8, int16, int32, int64, float32, float64:
			return nil
		default:
			return fmt.Errorf("expected number")
		}
	case TypeInteger:
		switch value.(type) {
		case int, int8, int16, int32, int64:
			return nil
		default:
			return fmt.Errorf("expected integer")
		}
	case TypeBoolean:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected boolean")
		}
	case "bool":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected boolean")
		}
	default:
		return fmt.Errorf("unsupported type %q", typeName)
	}

	return nil
}

func ValidateValueAgainstTypes(typeName string, value any, decls []model.TypeDecl) error {
	if err := ValidateValue(typeName, value); err == nil {
		return nil
	}

	typeDecl, ok := findTypeDecl(typeName, decls)
	if !ok {
		return fmt.Errorf("unsupported type %q", typeName)
	}

	fields, ok := value.(map[string]any)
	if !ok {
		return fmt.Errorf("expected object for type %q", typeName)
	}

	for _, field := range typeDecl.Fields {
		fieldValue, exists := fields[field.Name]
		if !exists {
			if field.Optional {
				continue
			}
			return fmt.Errorf("missing required field %q", field.Name)
		}

		if field.Array {
			rv := reflect.ValueOf(fieldValue)
			if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
				return fmt.Errorf("field %q: expected array", field.Name)
			}
			for i := 0; i < rv.Len(); i++ {
				if err := ValidateValueAgainstTypes(field.TypeName, rv.Index(i).Interface(), decls); err != nil {
					return fmt.Errorf("field %q: %w", field.Name, err)
				}
			}
			continue
		}

		if err := ValidateValueAgainstTypes(field.TypeName, fieldValue, decls); err != nil {
			return fmt.Errorf("field %q: %w", field.Name, err)
		}
	}

	return nil
}

func findTypeDecl(name string, decls []model.TypeDecl) (model.TypeDecl, bool) {
	for _, decl := range decls {
		if decl.Name == name {
			return decl, true
		}
	}
	return model.TypeDecl{}, false
}
