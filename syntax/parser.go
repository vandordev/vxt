package syntax

import (
	"errors"

	"github.com/alfariiizi/vxt/source"
)

var errUnexpectedEOF = errors.New("unterminated template expression")

type Part struct {
	Text string
	Expr string
}

func ParseTemplate(src source.Source) ([]Part, error) {
	tokens, err := lexTemplate(src.Text)
	if err != nil {
		return nil, err
	}

	parts := make([]Part, 0, len(tokens))
	for _, tok := range tokens {
		switch tok.kind {
		case tokenText:
			parts = append(parts, Part{Text: tok.value})
		case tokenExpr:
			parts = append(parts, Part{Expr: tok.value})
		}
	}

	return parts, nil
}
