package render

import (
	"strings"

	"github.com/alfariiizi/vxt/diag"
	"github.com/alfariiizi/vxt/expr"
	"github.com/alfariiizi/vxt/model"
	"github.com/alfariiizi/vxt/syntax"
)

func RenderSingle(tpl *model.CompiledTemplate, ctx map[string]any) (string, []diag.Diagnostic) {
	parts, err := syntax.ParseTemplate(tpl.Source)
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

		value, evalErr := expr.EvalPath(ctx, part.Expr)
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
