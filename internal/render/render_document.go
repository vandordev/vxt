package render

import (
	"github.com/vandordev/vxt/diag"
	"github.com/vandordev/vxt/internal/model"
	"github.com/vandordev/vxt/source"
)

func RenderDocumentBody(file model.FileBlock, input map[string]any) (string, []diag.Diagnostic) {
	return RenderDocumentBodyWithPartials(file, input, nil)
}

func RenderDocumentBodyWithPartials(file model.FileBlock, input map[string]any, partials map[string]string) (string, []diag.Diagnostic) {
	src := source.Source{
		ID:   file.Path,
		Text: file.Body,
	}
	return renderTemplateSourceWithPartials(src, input, partials)
}
