package write

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ErrFileExists struct {
	Path string
}

func (e ErrFileExists) Error() string {
	return fmt.Sprintf("file %q already exists", e.Path)
}

type ErrPathEscape struct {
	Path string
}

func (e ErrPathEscape) Error() string {
	return fmt.Sprintf("path %q escapes output root", e.Path)
}

type ErrUnsupportedWriteMode struct {
	Mode string
}

func (e ErrUnsupportedWriteMode) Error() string {
	return fmt.Sprintf("unsupported write mode %q", e.Mode)
}

type FilesystemTarget struct {
	root string
}

func NewFilesystemTarget(root string) *FilesystemTarget {
	return &FilesystemTarget{root: root}
}

func (f *FilesystemTarget) MkdirAll(path string) error {
	cleaned := filepath.Clean(path)
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") || filepath.IsAbs(cleaned) {
		return ErrPathEscape{Path: path}
	}

	fullPath := filepath.Join(f.root, cleaned)
	return os.MkdirAll(fullPath, 0o755)
}

func (f *FilesystemTarget) WriteFile(path string, content []byte, mode string) (bool, error) {
	cleaned := filepath.Clean(path)
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") || filepath.IsAbs(cleaned) {
		return false, ErrPathEscape{Path: path}
	}

	fullPath := filepath.Join(f.root, cleaned)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return false, err
	}

	_, statErr := os.Stat(fullPath)
	exists := statErr == nil
	if statErr != nil && !os.IsNotExist(statErr) {
		return false, statErr
	}

	switch mode {
	case "", "overwrite":
		return true, os.WriteFile(fullPath, content, 0o644)
	case "create":
		if exists {
			return false, ErrFileExists{Path: path}
		}
		return true, os.WriteFile(fullPath, content, 0o644)
	case "skip-if-exists":
		if exists {
			return false, nil
		}
		return true, os.WriteFile(fullPath, content, 0o644)
	default:
		return false, ErrUnsupportedWriteMode{Mode: mode}
	}
}
