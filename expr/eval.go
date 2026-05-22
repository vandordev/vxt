package expr

import "fmt"

func EvalPath(ctx map[string]any, path string) (string, error) {
	value, ok := ctx[path]
	if !ok {
		return "", fmt.Errorf("missing value for %q", path)
	}

	return fmt.Sprint(value), nil
}
