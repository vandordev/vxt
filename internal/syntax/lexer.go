package syntax

import "strings"

type tokenKind string

const (
	tokenText tokenKind = "text"
	tokenExpr tokenKind = "expr"
)

type token struct {
	kind  tokenKind
	value string
}

func lexTemplate(input string) ([]token, error) {
	var tokens []token

	for len(input) > 0 {
		start := strings.Index(input, "{{")
		if start == -1 {
			tokens = append(tokens, token{kind: tokenText, value: input})
			break
		}

		if start > 0 {
			tokens = append(tokens, token{kind: tokenText, value: input[:start]})
		}

		end := strings.Index(input[start+2:], "}}")
		if end == -1 {
			return nil, errUnexpectedEOF
		}

		expr := strings.TrimSpace(input[start+2 : start+2+end])
		tokens = append(tokens, token{kind: tokenExpr, value: expr})
		input = input[start+2+end+2:]
	}

	return tokens, nil
}
