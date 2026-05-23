package runtime

import (
	"errors"

	"github.com/alfariiizi/vxt/diag"
	planpkg "github.com/alfariiizi/vxt/plan"
	"github.com/alfariiizi/vxt/write"
)

type WriteResult struct {
	Report      write.WriteReport
	Diagnostics []diag.Diagnostic
	Err         error
}

func WritePlan(p planpkg.Plan, target write.OutputTarget) (write.WriteReport, error) {
	result := WritePlanWithDiagnostics(p, target)
	return result.Report, result.Err
}

func WritePlanWithDiagnostics(p planpkg.Plan, target write.OutputTarget) WriteResult {
	result := WriteResult{}
	for _, dir := range p.Dirs {
		if err := target.MkdirAll(dir.Path); err != nil {
			result.Err = err
			result.Diagnostics = append(result.Diagnostics, diagnosticFromWriteError(err))
			return result
		}
		result.Report.DirsWritten++
	}
	for _, file := range p.Files {
		written, err := target.WriteFile(file.Path, []byte(file.Content), file.Mode)
		if err != nil {
			result.Err = err
			result.Diagnostics = append(result.Diagnostics, diagnosticFromWriteError(err))
			return result
		}
		if written {
			result.Report.FilesWritten++
			continue
		}
		result.Report.FilesSkipped++
	}
	return result
}

func diagnosticFromWriteError(err error) diag.Diagnostic {
	var fileExists write.ErrFileExists
	if errors.As(err, &fileExists) {
		return diag.Diagnostic{
			Code:     diag.CodeWriteFileExists,
			Severity: diag.SeverityError,
			Message:  err.Error(),
			Hint:     "use mode=overwrite or mode=skip-if-exists if replacement is intended",
		}
	}

	var pathEscape write.ErrPathEscape
	if errors.As(err, &pathEscape) {
		return diag.Diagnostic{
			Code:     diag.CodeWritePathEscape,
			Severity: diag.SeverityError,
			Message:  err.Error(),
			Hint:     "rendered paths must stay inside the output root",
		}
	}

	var unsupportedMode write.ErrUnsupportedWriteMode
	if errors.As(err, &unsupportedMode) {
		return diag.Diagnostic{
			Code:     diag.CodeWriteUnsupportedMode,
			Severity: diag.SeverityError,
			Message:  err.Error(),
			Hint:     "supported modes are create, overwrite, and skip-if-exists",
		}
	}

	return diag.Diagnostic{
		Code:     diag.CodeOutputConflict,
		Severity: diag.SeverityError,
		Message:  err.Error(),
	}
}
