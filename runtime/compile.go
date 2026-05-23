package runtime

import (
	"fmt"
	"strings"

	"github.com/alfariiizi/vxt/diag"
	"github.com/alfariiizi/vxt/model"
	"github.com/alfariiizi/vxt/source"
	"github.com/alfariiizi/vxt/syntax"
)

type CompileResult struct {
	Template    *model.CompiledTemplate
	Document    *model.CompiledDocument
	Diagnostics []diag.Diagnostic
}

type SourceResolver interface {
	Resolve(path string) (source.Source, error)
}

type MapResolver map[string]source.Source

func (m MapResolver) Resolve(path string) (source.Source, error) {
	src, ok := m[path]
	if !ok {
		return source.Source{}, fmt.Errorf("missing source for use path %q", path)
	}
	return src, nil
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

func CompileDocument(src source.Source) CompileResult {
	doc, err := syntax.ParseDocument(src)
	if err == nil {
		return CompileResult{
			Document: doc,
		}
	}

	return CompileResult{
		Diagnostics: []diag.Diagnostic{{
			Code:     diag.CodeParseUnexpectedEOF,
			Severity: diag.SeverityError,
			Message:  err.Error(),
			Span: source.Span{
				SourceID: src.ID,
				Start:    0,
				End:      len(src.Text),
			},
		}},
	}
}

func CompileDocumentWithResolver(src source.Source, resolver SourceResolver) CompileResult {
	result := CompileDocument(src)
	if result.Document == nil || len(result.Diagnostics) > 0 {
		return result
	}

	for _, use := range result.Document.Uses {
		resolved, err := resolver.Resolve(use.Path)
		if err != nil {
			return compileErrorResult(src, err)
		}
		imported, err := syntax.ParseDefinitionDocument(resolved)
		if err != nil {
			return compileErrorResult(src, err)
		}
		result.Document.Types = append(result.Document.Types, imported.Types...)
	}

	return result
}

func compileErrorResult(src source.Source, err error) CompileResult {
	return CompileResult{
		Diagnostics: []diag.Diagnostic{{
			Code:     diag.CodeParseUnexpectedEOF,
			Severity: diag.SeverityError,
			Message:  err.Error(),
			Span: source.Span{
				SourceID: src.ID,
				Start:    0,
				End:      len(src.Text),
			},
		}},
	}
}
