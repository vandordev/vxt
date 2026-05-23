package runtime

import (
	"github.com/alfariiizi/vxt/diag"
	"github.com/alfariiizi/vxt/expr"
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
	plannedFilePaths := map[string]struct{}{}
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
		path, diags := render.RenderTemplateSource(source.Source{
			ID:   file.Path,
			Text: file.Path,
		}, validated.Input)
		if len(diags) > 0 {
			result.Diagnostics = append(result.Diagnostics, diags...)
			return result
		}
		if conflictDiag, ok := duplicateFilePathDiagnostic(path, plannedFilePaths); ok {
			result.Diagnostics = append(result.Diagnostics, conflictDiag)
			return result
		}

		content, diags := render.RenderDocumentBodyWithPartials(file, validated.Input, partials)
		if len(diags) > 0 {
			result.Diagnostics = append(result.Diagnostics, diags...)
			return result
		}

		result.Plan.Files = append(result.Plan.Files, planpkg.FileOutput{
			Path:    path,
			Content: content,
			Mode:    file.Mode,
		})
	}

	for _, conditional := range validated.Document.Conditionals {
		value, err := expr.EvalValue(validated.Input, conditional.Condition)
		if err != nil {
			result.Diagnostics = append(result.Diagnostics, diag.Diagnostic{
				Code:     diag.CodeRenderMissingValue,
				Severity: diag.SeverityError,
				Message:  err.Error(),
			})
			return result
		}
		if !expr.IsTruthy(value) {
			continue
		}

		for _, dir := range conditional.Dirs {
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

		for _, file := range conditional.Files {
			path, diags := render.RenderTemplateSource(source.Source{
				ID:   file.Path,
				Text: file.Path,
			}, validated.Input)
			if len(diags) > 0 {
				result.Diagnostics = append(result.Diagnostics, diags...)
				return result
			}
			if conflictDiag, ok := duplicateFilePathDiagnostic(path, plannedFilePaths); ok {
				result.Diagnostics = append(result.Diagnostics, conflictDiag)
				return result
			}

			content, diags := render.RenderDocumentBodyWithPartials(file, validated.Input, partials)
			if len(diags) > 0 {
				result.Diagnostics = append(result.Diagnostics, diags...)
				return result
			}

			result.Plan.Files = append(result.Plan.Files, planpkg.FileOutput{
				Path:    path,
				Content: content,
				Mode:    file.Mode,
			})
		}
	}

	for _, hook := range validated.Document.Hooks {
		result.Plan.PlannedHooks = append(result.Plan.PlannedHooks, planpkg.HooksFromDecls(hook.Event, hook.Run))
	}

	return result
}

func duplicateFilePathDiagnostic(path string, seen map[string]struct{}) (diag.Diagnostic, bool) {
	if _, exists := seen[path]; exists {
		return diag.Diagnostic{
			Code:     diag.CodeOutputConflict,
			Severity: diag.SeverityError,
			Message:  "duplicate file output path: " + path,
		}, true
	}

	seen[path] = struct{}{}
	return diag.Diagnostic{}, false
}
