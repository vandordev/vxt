package bind

import (
	"errors"
	"fmt"

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

	_, err := analyzeDocument(req.PackageName, compiled.Document)
	if err != nil {
		return Output{}, err
	}

	return Output{}, ErrNotImplemented
}
