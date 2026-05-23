package runtime

import (
	planpkg "github.com/alfariiizi/vxt/plan"
	"github.com/alfariiizi/vxt/write"
)

func WritePlan(p planpkg.Plan, target write.OutputTarget) (write.WriteReport, error) {
	report := write.WriteReport{}
	for _, dir := range p.Dirs {
		if err := target.MkdirAll(dir.Path); err != nil {
			return report, err
		}
		report.DirsWritten++
	}
	for _, file := range p.Files {
		written, err := target.WriteFile(file.Path, []byte(file.Content), file.Mode)
		if err != nil {
			return report, err
		}
		if written {
			report.FilesWritten++
			continue
		}
		report.FilesSkipped++
	}
	return report, nil
}
