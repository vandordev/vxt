package bind

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/vandordev/vxt/runtime"
)

// ErrNotImplemented is returned until the generator is implemented.
var ErrNotImplemented = errors.New("bind: generator not implemented")

// Generate produces typed Go bindings for one document template.
func Generate(req Request) (Output, error) {
	if req.PackageName == "" {
		return Output{}, fmt.Errorf("bind: package name is required")
	}

	compiled := runtime.CompileDocument(req.Document)
	if len(compiled.Diagnostics) > 0 {
		return Output{}, fmt.Errorf("bind: compile failed: %s", compiled.Diagnostics[0].Message)
	}

	useSources, err := collectLocalUseSources(compiled.Document, req.Uses)
	if err != nil {
		return Output{}, err
	}
	if len(useSources) > 0 {
		compiled = runtime.CompileDocumentWithResolver(req.Document, runtime.MapResolver(useSources))
		if len(compiled.Diagnostics) > 0 {
			return Output{}, fmt.Errorf("bind: compile failed: %s", compiled.Diagnostics[0].Message)
		}
	}

	analyzed, err := analyzeDocument(req.PackageName, compiled.Document)
	if err != nil {
		return Output{}, err
	}

	generated, err := emitGeneratedFile(analyzed, embeddedAssets{
		Main: req.Document,
		Uses: useSources,
	})
	if err != nil {
		return Output{}, err
	}

	return Output{
		BindingName: analyzed.Template,
		PackageName: analyzed.PackageName,
		Files: []File{{
			Path:    filepath.Join(".vxt", analyzed.Template+"_gen.go"),
			Content: generated,
		}},
	}, nil
}

// GenerateToDir generates one binding output set and writes it into dir.
func GenerateToDir(req Request, dir string, opts WriteOptions) (WriteReport, error) {
	out, err := Generate(req)
	if err != nil {
		return WriteReport{}, err
	}
	return WriteOutput(out, dir, opts)
}
