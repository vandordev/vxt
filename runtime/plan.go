package runtime

import (
	"github.com/alfariiizi/vxt/diag"
	planpkg "github.com/alfariiizi/vxt/plan"
	"github.com/alfariiizi/vxt/render"
)

type PlanResult struct {
	Plan        planpkg.Plan
	Diagnostics []diag.Diagnostic
}

func PlanDocument(validated ValidationResult) PlanResult {
	result := PlanResult{}
	if len(validated.Diagnostics) > 0 {
		result.Diagnostics = append(result.Diagnostics, validated.Diagnostics...)
		return result
	}

	for _, file := range validated.Document.Files {
		content, diags := render.RenderDocumentBody(file, validated.Input)
		if len(diags) > 0 {
			result.Diagnostics = append(result.Diagnostics, diags...)
			return result
		}

		result.Plan.Files = append(result.Plan.Files, planpkg.FileOutput{
			Path:    file.Path,
			Content: content,
			Mode:    file.Mode,
		})
	}

	for _, hook := range validated.Document.Hooks {
		result.Plan.PlannedHooks = append(result.Plan.PlannedHooks, planpkg.HooksFromDecls(hook.Event, hook.Run))
	}

	return result
}
