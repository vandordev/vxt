package bind

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteOutput writes generated binding output into one explicit output directory.
func WriteOutput(out Output, dir string) (WriteReport, error) {
	if dir == "" {
		return WriteReport{}, fmt.Errorf("bind: output dir is required")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return WriteReport{}, fmt.Errorf("bind: create output dir: %w", err)
	}

	report := WriteReport{OutputDir: dir}
	for _, file := range out.Files {
		overwrote, absPath, err := writeGeneratedFile(dir, file)
		if err != nil {
			return report, err
		}
		if overwrote {
			report.FilesOverwritten = append(report.FilesOverwritten, absPath)
			continue
		}
		report.FilesWritten = append(report.FilesWritten, absPath)
	}

	return report, nil
}

func writeGeneratedFile(dir string, file File) (overwrote bool, absPath string, err error) {
	target := filepath.Join(dir, filepath.Base(file.Path))
	if _, statErr := os.Stat(target); statErr == nil {
		overwrote = true
	} else if !os.IsNotExist(statErr) {
		return false, "", fmt.Errorf("bind: stat generated file: %w", statErr)
	}

	if err := os.WriteFile(target, []byte(file.Content), 0o644); err != nil {
		return false, "", fmt.Errorf("bind: write generated file: %w", err)
	}
	return overwrote, target, nil
}
