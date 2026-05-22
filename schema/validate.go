package schema

import "fmt"

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
	default:
		return fmt.Errorf("unsupported type %q", typeName)
	}

	return nil
}
