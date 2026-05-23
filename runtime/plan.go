package runtime

import (
	"github.com/alfariiizi/vxt/diag"
	planpkg "github.com/alfariiizi/vxt/plan"
	"github.com/alfariiizi/vxt/render"
	"github.com/alfariiizi/vxt/source"
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

	partials := make(map[string]string, len(validated.Document.Partials))
	for _, partial := range validated.Document.Partials {
		partials[partial.Name] = partial.Body
	}

	for _, dir := range validated.Document.Dirs {
		path, diags := render.RenderTemplateSource(source.Source{
			ID:   dir.Path,
			Text: dir.Path,
		}, validated.Input)
		if len(diags) > 0 {
			result.Diagnostics = append(result.Diagnostics, diags...)
			return result
		}
		result.Plan.Dirs = append(result.Plan.Dirs, planpkg.DirOutput{Path: path})
	}

	for _, file := range validated.Document.Files {
		content, diags := render.RenderDocumentBodyWithPartials(file, validated.Input, partials)
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
