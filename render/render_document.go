package render

import (
	"github.com/alfariiizi/vxt/diag"
	"github.com/alfariiizi/vxt/model"
	"github.com/alfariiizi/vxt/source"
	"github.com/alfariiizi/vxt/syntax"
)

func RenderDocumentBody(file model.FileBlock, input map[string]any) (string, []diag.Diagnostic) {
	src := source.Source{
		ID:   file.Path,
		Text: file.Body,
	}
	nodes, err := syntax.ParseTemplate(src)
	if err != nil {
		return "", []diag.Diagnostic{{
			Code:     diag.CodeParseUnexpectedEOF,
			Severity: diag.SeverityError,
			Message:  err.Error(),
		}}
	}

	return renderNodes(nodes, input)
}
