package runtime

import (
	planpkg "github.com/alfariiizi/vxt/plan"
	"github.com/alfariiizi/vxt/write"
)

func WritePlan(p planpkg.Plan, target write.OutputTarget) (write.WriteReport, error) {
	report := write.WriteReport{}
	for _, file := range p.Files {
		if err := target.WriteFile(file.Path, []byte(file.Content)); err != nil {
			return report, err
		}
		report.FilesWritten++
	}
	return report, nil
}
