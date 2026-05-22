package write

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type FilesystemTarget struct {
	root string
}

func NewFilesystemTarget(root string) *FilesystemTarget {
	return &FilesystemTarget{root: root}
}

func (f *FilesystemTarget) WriteFile(path string, content []byte) error {
	cleaned := filepath.Clean(path)
	if cleaned == ".." || strings.HasPrefix(cleaned, "../") || filepath.IsAbs(cleaned) {
		return fmt.Errorf("path %q escapes output root", path)
	}

	fullPath := filepath.Join(f.root, cleaned)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(fullPath, content, 0o644)
}
