package render

import (
	"reflect"
	"strings"

	"github.com/vandordev/vxt/diag"
	"github.com/vandordev/vxt/internal/expr"
	"github.com/vandordev/vxt/internal/model"
	"github.com/vandordev/vxt/source"
	"github.com/vandordev/vxt/internal/syntax"
)

func RenderSingle(tpl *model.CompiledTemplate, ctx map[string]any) (string, []diag.Diagnostic) {
	return RenderTemplateSource(tpl.Source, ctx)
}

func RenderTemplateSource(src source.Source, ctx map[string]any) (string, []diag.Diagnostic) {
	return renderTemplateSourceWithPartials(src, ctx, nil)
}

func renderTemplateSourceWithPartials(src source.Source, ctx map[string]any, partials map[string]string) (string, []diag.Diagnostic) {
	nodes, err := syntax.ParseTemplate(src)
	if err != nil {
		return "", []diag.Diagnostic{{
			Code:     diag.CodeParseUnexpectedEOF,
			Severity: diag.SeverityError,
			Message:  err.Error(),
		}}
	}

	return renderNodes(nodes, ctx, partials)
}

func renderNodes(nodes []syntax.Node, ctx map[string]any, partials map[string]string) (string, []diag.Diagnostic) {
	var out strings.Builder
	for _, rawNode := range nodes {
		switch node := rawNode.(type) {
		case syntax.TextNode:
			out.WriteString(node.Text)
		case syntax.ExprNode:
			value, evalErr := expr.EvalPath(ctx, node.Expr)
			if evalErr != nil {
				return "", []diag.Diagnostic{{
					Code:     diag.CodeRenderMissingValue,
					Severity: diag.SeverityError,
					Message:  evalErr.Error(),
				}}
			}
			out.WriteString(value)
		case syntax.IfNode:
			value, evalErr := expr.EvalValue(ctx, node.Cond)
			if evalErr != nil {
				return "", []diag.Diagnostic{{
					Code:     diag.CodeRenderMissingValue,
					Severity: diag.SeverityError,
					Message:  evalErr.Error(),
				}}
			}
			var branch []syntax.Node
			if expr.IsTruthy(value) {
				branch = node.Then
			} else {
				branch = node.Else
			}
			rendered, diags := renderNodes(branch, ctx, partials)
			if len(diags) > 0 {
				return "", diags
			}
			out.WriteString(rendered)
		case syntax.EachNode:
			value, evalErr := expr.EvalValue(ctx, node.Collection)
			if evalErr != nil {
				return "", []diag.Diagnostic{{
					Code:     diag.CodeRenderMissingValue,
					Severity: diag.SeverityError,
					Message:  evalErr.Error(),
				}}
			}
			rv := reflect.ValueOf(value)
			if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
				return "", []diag.Diagnostic{{
					Code:     diag.CodeTypeMismatch,
					Severity: diag.SeverityError,
					Message:  "each requires slice or array value",
				}}
			}
			for i := 0; i < rv.Len(); i++ {
				childCtx := cloneContext(ctx)
				childCtx[node.Item] = rv.Index(i).Interface()
				rendered, diags := renderNodes(node.Body, childCtx, partials)
				if len(diags) > 0 {
					return "", diags
				}
				out.WriteString(rendered)
			}
		case syntax.IncludeNode:
			if partials == nil {
				return "", []diag.Diagnostic{{
					Code:     diag.CodeRenderMissingValue,
					Severity: diag.SeverityError,
					Message:  "missing partial for include " + node.Target,
				}}
			}
			body, ok := partials[node.Target]
			if !ok {
				return "", []diag.Diagnostic{{
					Code:     diag.CodeRenderMissingValue,
					Severity: diag.SeverityError,
					Message:  "missing partial for include " + node.Target,
				}}
			}
			rendered, diags := renderTemplateSourceWithPartials(source.Source{
				ID:   node.Target,
				Text: body,
			}, ctx, partials)
			if len(diags) > 0 {
				return "", diags
			}
			out.WriteString(rendered)
		}
	}

	return out.String(), nil
}

func cloneContext(ctx map[string]any) map[string]any {
	cloned := make(map[string]any, len(ctx)+1)
	for key, value := range ctx {
		cloned[key] = value
	}
	return cloned
}
