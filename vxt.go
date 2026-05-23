// Package vxt provides a spec-first template compiler/runtime for code and
// file generation.
//
// The canonical product model is a staged pipeline:
// compile -> validate -> plan -> write.
//
// v0.1 intentionally excludes hook execution, trust policy, package semantics,
// and CLI behavior from the core contract.
package vxt

import (
	"github.com/vandordev/vxt/diag"
	"github.com/vandordev/vxt/internal/render"
	"github.com/vandordev/vxt/runtime"
	"github.com/vandordev/vxt/source"
)

// RenderSingleFile is the convenience API for single-file mode in v0.1.
func RenderSingleFile(src source.Source, ctx map[string]any) (string, []diag.Diagnostic) {
	compiled := runtime.CompileSingle(src)
	if len(compiled.Diagnostics) > 0 {
		return "", compiled.Diagnostics
	}

	return render.RenderTemplateSource(compiled.Template.Source, ctx)
}
