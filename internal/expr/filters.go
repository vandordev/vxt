package expr

import (
	"fmt"
	"strings"
	"unicode"
)

func applyFilter(value any, filter string) (any, error) {
	input := fmt.Sprint(value)

	switch filter {
	case "snake":
		return strings.Join(words(input), "_"), nil
	case "upper_snake":
		return strings.ToUpper(strings.Join(words(input), "_")), nil
	case "kebab":
		return strings.Join(words(input), "-"), nil
	case "pascal":
		return pascalCase(input), nil
	case "camel":
		pascal := pascalCase(input)
		if pascal == "" {
			return "", nil
		}
		runes := []rune(pascal)
		runes[0] = unicode.ToLower(runes[0])
		return string(runes), nil
	case "lower":
		return strings.ToLower(input), nil
	case "upper":
		return strings.ToUpper(input), nil
	default:
		return nil, fmt.Errorf("unsupported filter %q", filter)
	}
}

func pascalCase(input string) string {
	parts := words(input)
	var out strings.Builder
	for _, part := range parts {
		runes := []rune(part)
		if len(runes) == 0 {
			continue
		}
		runes[0] = unicode.ToUpper(runes[0])
		out.WriteString(string(runes))
	}
	return out.String()
}

func words(input string) []string {
	var out []string
	for _, segment := range splitSeparators(input) {
		out = append(out, splitCase(segment)...)
	}
	return out
}

func splitSeparators(input string) []string {
	var segments []string
	var current strings.Builder

	flush := func() {
		if current.Len() == 0 {
			return
		}
		segments = append(segments, current.String())
		current.Reset()
	}

	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current.WriteRune(r)
			continue
		}
		flush()
	}
	flush()

	return segments
}

func splitCase(segment string) []string {
	runes := []rune(segment)
	if len(runes) == 0 {
		return nil
	}

	var out []string
	start := 0
	for i := 1; i < len(runes); i++ {
		prev := runes[i-1]
		curr := runes[i]
		var next rune
		if i+1 < len(runes) {
			next = runes[i+1]
		}

		if startsWord(prev, curr, next) {
			out = append(out, strings.ToLower(string(runes[start:i])))
			start = i
		}
	}
	out = append(out, strings.ToLower(string(runes[start:])))

	return out
}

func startsWord(prev, curr, next rune) bool {
	if unicode.IsLower(prev) && unicode.IsUpper(curr) {
		return true
	}
	if unicode.IsDigit(prev) != unicode.IsDigit(curr) {
		return true
	}
	return unicode.IsUpper(prev) && unicode.IsUpper(curr) && next != 0 && unicode.IsLower(next)
}
