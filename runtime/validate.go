package runtime

import (
	"fmt"

	"github.com/alfariiizi/vxt/diag"
	"github.com/alfariiizi/vxt/model"
	"github.com/alfariiizi/vxt/schema"
)

type ValidationResult struct {
	Document    *model.CompiledDocument
	Input       map[string]any
	Diagnostics []diag.Diagnostic
}

func ValidateDocument(doc *model.CompiledDocument, input map[string]any) ValidationResult {
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

		if err := schema.ValidateValueAgainstTypes(decl.TypeName, value, doc.Types); err != nil {
			result.Diagnostics = append(result.Diagnostics, diag.Diagnostic{
				Code:     diag.CodeTypeMismatch,
				Severity: diag.SeverityError,
				Message:  fmt.Sprintf("input %q: %v", decl.Name, err),
			})
		}
	}

	return result
}
