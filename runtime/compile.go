package runtime

import (
	"strings"

	"github.com/alfariiizi/vxt/diag"
	"github.com/alfariiizi/vxt/model"
	"github.com/alfariiizi/vxt/source"
	"github.com/alfariiizi/vxt/syntax"
)

type CompileResult struct {
	Template    *model.CompiledTemplate
	Diagnostics []diag.Diagnostic
}

func CompileSingle(src source.Source) CompileResult {
	_, err := syntax.ParseTemplate(src)
	if err == nil {
		return CompileResult{
			Template: &model.CompiledTemplate{Source: src},
		}
	}

	return CompileResult{
		Diagnostics: []diag.Diagnostic{{
			Code:     diag.CodeParseUnexpectedEOF,
			Severity: diag.SeverityError,
			Message:  err.Error(),
			Span: source.Span{
				SourceID: src.ID,
				Start:    strings.Index(src.Text, "{{"),
				End:      len(src.Text),
			},
		}},
	}
}
