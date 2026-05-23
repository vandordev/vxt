package runtime

import (
	"fmt"

	"github.com/vandordev/vxt/diag"
	"github.com/vandordev/vxt/internal/schema"
)

// ValidationResult captures typed input validation for one compiled document.
type ValidationResult struct {
	Document    *CompiledDocument
	Input       map[string]any
	Diagnostics []diag.Diagnostic
}

// ValidateDocument checks declared document inputs against the compiled document type declarations.
func ValidateDocument(doc *CompiledDocument, input map[string]any) ValidationResult {
	result := ValidationResult{
		Document: doc,
		Input:    input,
	}

	for _, decl := range doc.Inputs {
		value, ok := input[decl.Name]
		if !ok {
			result.Diagnostics = append(result.Diagnostics, diag.Diagnostic{
				Code:     diag.CodeRenderMissingValue,
				Severity: diag.SeverityError,
				Message:  fmt.Sprintf("missing input %q", decl.Name),
			})
			continue
		}

		if err := schema.ValidateValueAgainstTypes(decl.TypeName, value, typeDeclsToInternal(doc.Types)); err != nil {
			result.Diagnostics = append(result.Diagnostics, diag.Diagnostic{
				Code:     diag.CodeTypeMismatch,
				Severity: diag.SeverityError,
				Message:  fmt.Sprintf("input %q: %v", decl.Name, err),
			})
		}
	}

	return result
}
