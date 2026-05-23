package runtime

import (
	"fmt"
	"strings"

	"github.com/vandordev/vxt/diag"
	"github.com/vandordev/vxt/internal/syntax"
	"github.com/vandordev/vxt/source"
)

// CompileResult captures the outcome of compiling a single-file or document template.
type CompileResult struct {
	Template    *CompiledTemplate
	Document    *CompiledDocument
	Diagnostics []diag.Diagnostic
}

// SourceResolver resolves one @use path into an in-memory source document.
type SourceResolver interface {
	Resolve(path string) (source.Source, error)
}

// MapResolver is a simple in-memory SourceResolver for tests and embedding use cases.
type MapResolver map[string]source.Source

func (m MapResolver) Resolve(path string) (source.Source, error) {
	src, ok := m[path]
	if !ok {
		return source.Source{}, fmt.Errorf("missing source for use path %q", path)
	}
	return src, nil
}

// CompileSingle parses single-file template syntax and returns structured diagnostics on failure.
func CompileSingle(src source.Source) CompileResult {
	_, err := syntax.ParseTemplate(src)
	if err == nil {
		return CompileResult{
			Template: &CompiledTemplate{Source: src},
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

// CompileDocument parses document-mode syntax and returns the compiled document contract.
func CompileDocument(src source.Source) CompileResult {
	doc, err := syntax.ParseDocument(src)
	if err == nil {
		return CompileResult{
			Document: compiledDocumentFromInternal(doc),
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

// CompileDocumentWithResolver parses a document and resolves any referenced definition documents.
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
		result.Document.Types = append(result.Document.Types, typeDeclsFromInternal(imported.Types)...)
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
