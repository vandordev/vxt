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

	removed, err := reconcileOwnedFiles(out, dir)
	if err != nil {
		return report, err
	}
	report.FilesRemoved = append(report.FilesRemoved, removed...)

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

func reconcileOwnedFiles(out Output, dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("bind: read output dir: %w", err)
	}

	active := make(map[string]struct{}, len(out.Files))
	for _, file := range out.Files {
		active[filepath.Base(file.Path)] = struct{}{}
	}

	var removed []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if _, ok := active[name]; ok {
			continue
		}

		target := filepath.Join(dir, name)
		content, err := os.ReadFile(target)
		if err != nil {
			return removed, fmt.Errorf("bind: read candidate stale file: %w", err)
		}
		if !isOwnedBindingFile(out.BindingName, out.PackageName, target, string(content)) {
			continue
		}

		if err := os.Remove(target); err != nil {
			return removed, fmt.Errorf("bind: remove stale file: %w", err)
		}
		removed = append(removed, target)
	}

	return removed, nil
}
