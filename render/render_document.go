package render

import (
	"strings"

	"github.com/alfariiizi/vxt/diag"
	"github.com/alfariiizi/vxt/expr"
	"github.com/alfariiizi/vxt/model"
	"github.com/alfariiizi/vxt/source"
	"github.com/alfariiizi/vxt/syntax"
)

func RenderDocumentBody(file model.FileBlock, input map[string]any) (string, []diag.Diagnostic) {
	src := source.Source{
		ID:   file.Path,
		Text: file.Body,
	}
	parts, err := syntax.ParseTemplate(src)
	if err != nil {
		return "", []diag.Diagnostic{{
			Code:     diag.CodeParseUnexpectedEOF,
			Severity: diag.SeverityError,
			Message:  err.Error(),
		}}
	}

	var out strings.Builder
	for _, part := range parts {
		if part.Expr == "" {
			out.WriteString(part.Text)
			continue
		}

		value, evalErr := expr.EvalPath(input, part.Expr)
		if evalErr != nil {
			return "", []diag.Diagnostic{{
				Code:     diag.CodeRenderMissingValue,
				Severity: diag.SeverityError,
				Message:  evalErr.Error(),
			}}
		}
		out.WriteString(value)
	}

	return out.String(), nil
}
