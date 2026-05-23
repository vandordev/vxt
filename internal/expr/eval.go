package expr

import (
	"fmt"
	"reflect"
	"strings"
)

func EvalPath(ctx map[string]any, path string) (string, error) {
	value, err := EvalValue(ctx, path)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(value), nil
}

func EvalValue(ctx map[string]any, path string) (any, error) {
	current, ok := ctx[path]
	if ok {
		return current, nil
	}

	parts := strings.Split(path, ".")
	if len(parts) == 0 {
		return nil, fmt.Errorf("missing value for %q", path)
	}

	current, ok = ctx[parts[0]]
	if !ok {
		return nil, fmt.Errorf("missing value for %q", path)
	}

	for _, part := range parts[1:] {
		next, err := access(current, part)
		if err != nil {
			return nil, fmt.Errorf("missing value for %q", path)
		}
		current = next
	}

	return current, nil
}

func IsTruthy(value any) bool {
	if value == nil {
		return false
	}

	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		return typed != ""
	}

	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		return rv.Len() > 0
	case reflect.Ptr, reflect.Interface:
		return !rv.IsNil()
	default:
		return true
	}
}

func access(value any, key string) (any, error) {
	rv := reflect.ValueOf(value)
	for rv.Kind() == reflect.Pointer || rv.Kind() == reflect.Interface {
		if rv.IsNil() {
			return nil, fmt.Errorf("nil value")
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.Map:
		mv := rv.MapIndex(reflect.ValueOf(key))
		if !mv.IsValid() {
			return nil, fmt.Errorf("missing key")
		}
		return mv.Interface(), nil
	case reflect.Struct:
		field := rv.FieldByNameFunc(func(name string) bool {
			return strings.EqualFold(name, key)
		})
		if !field.IsValid() {
			return nil, fmt.Errorf("missing field")
		}
		return field.Interface(), nil
	default:
		return nil, fmt.Errorf("unsupported access")
	}
}
