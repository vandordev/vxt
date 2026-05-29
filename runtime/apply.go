package runtime

import "github.com/vandordev/vxt/write"

// ApplyPlan writes a plan and then executes supported hooks through the provided executor.
func ApplyPlan(p Plan, target write.OutputTarget, executor HookExecutor) ApplyResult {
	result := ApplyResult{
		WriteResult: WritePlanWithDiagnostics(p, target),
	}
	if result.WriteResult.Err != nil || executor == nil {
		return result
	}

	for _, hook := range p.PlannedHooks {
		if hook.Event != "after:write" {
			continue
		}
		err := executor.Execute(HookContext{
			Event:       hook.Event,
			Plan:        p,
			WriteReport: result.WriteResult.Report,
		}, hook)
		if err != nil {
			result.HookErrors = append(result.HookErrors, err)
		}
	}

	return result
}
