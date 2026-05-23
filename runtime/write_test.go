package runtime_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alfariiizi/vxt/plan"
	"github.com/alfariiizi/vxt/runtime"
	"github.com/alfariiizi/vxt/write"
)

func TestWritePlanToMemoryTarget(t *testing.T) {
	p := plan.Plan{
		Files: []plan.FileOutput{{Path: "hello.txt", Content: "Hello Fariz"}},
	}

	target := write.NewMemoryTarget()
	report, err := runtime.WritePlan(p, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.FilesWritten != 1 {
		t.Fatalf("got %d", report.FilesWritten)
	}
}

func TestWritePlanRejectsPathOutsideFilesystemTargetRoot(t *testing.T) {
	p := plan.Plan{
		Files: []plan.FileOutput{{Path: "../escape.txt", Content: "nope"}},
	}

	target := write.NewFilesystemTarget(t.TempDir())
	_, err := runtime.WritePlan(p, target)
	if err == nil {
		t.Fatal("expected sandbox error")
	}
}

func TestWritePlanToFilesystemTarget(t *testing.T) {
	root := t.TempDir()
	p := plan.Plan{
		Files: []plan.FileOutput{{Path: "nested/hello.txt", Content: "Hello Fariz"}},
	}

	target := write.NewFilesystemTarget(root)
	report, err := runtime.WritePlan(p, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.FilesWritten != 1 {
		t.Fatalf("got %d", report.FilesWritten)
	}

	if _, err := filepath.Abs(filepath.Join(root, "nested", "hello.txt")); err != nil {
		t.Fatalf("expected written path to resolve: %v", err)
	}
}

func TestWritePlanCreatesDirectoriesOnFilesystemTarget(t *testing.T) {
	root := t.TempDir()
	p := plan.Plan{
		Dirs: []plan.DirOutput{{Path: "nested/modules"}},
	}

	target := write.NewFilesystemTarget(root)
	report, err := runtime.WritePlan(p, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.DirsWritten != 1 {
		t.Fatalf("got %d dirs written", report.DirsWritten)
	}

	info, err := os.Stat(filepath.Join(root, "nested", "modules"))
	if err != nil {
		t.Fatalf("expected dir to exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("expected created path to be dir")
	}
}

func TestWritePlanCreateModeFailsWhenFileExists(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "hello.txt")
	if err := os.WriteFile(path, []byte("existing"), 0o644); err != nil {
		t.Fatalf("seed existing file: %v", err)
	}

	p := plan.Plan{
		Files: []plan.FileOutput{{Path: "hello.txt", Content: "new", Mode: "create"}},
	}

	target := write.NewFilesystemTarget(root)
	_, err := runtime.WritePlan(p, target)
	if err == nil {
		t.Fatal("expected create mode error for existing file")
	}
}

func TestWritePlanOverwriteModeReplacesExistingFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "hello.txt")
	if err := os.WriteFile(path, []byte("existing"), 0o644); err != nil {
		t.Fatalf("seed existing file: %v", err)
	}

	p := plan.Plan{
		Files: []plan.FileOutput{{Path: "hello.txt", Content: "new", Mode: "overwrite"}},
	}

	target := write.NewFilesystemTarget(root)
	_, err := runtime.WritePlan(p, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(content) != "new" {
		t.Fatalf("got content %q", string(content))
	}
}

func TestWritePlanSkipIfExistsModeKeepsExistingFile(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "hello.txt")
	if err := os.WriteFile(path, []byte("existing"), 0o644); err != nil {
		t.Fatalf("seed existing file: %v", err)
	}

	p := plan.Plan{
		Files: []plan.FileOutput{{Path: "hello.txt", Content: "new", Mode: "skip-if-exists"}},
	}

	target := write.NewFilesystemTarget(root)
	report, err := runtime.WritePlan(p, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if report.FilesWritten != 0 {
		t.Fatalf("got %d files written", report.FilesWritten)
	}
	if report.FilesSkipped != 1 {
		t.Fatalf("got %d files skipped", report.FilesSkipped)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}
	if string(content) != "existing" {
		t.Fatalf("got content %q", string(content))
	}
}

func TestWritePlanWithDiagnosticsReturnsFileExistsDiagnostic(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "hello.txt")
	if err := os.WriteFile(path, []byte("existing"), 0o644); err != nil {
		t.Fatalf("seed existing file: %v", err)
	}

	p := plan.Plan{
		Files: []plan.FileOutput{{Path: "hello.txt", Content: "new", Mode: "create"}},
	}

	target := write.NewFilesystemTarget(root)
	result := runtime.WritePlanWithDiagnostics(p, target)
	if result.Err == nil {
		t.Fatal("expected write error")
	}
	if len(result.Diagnostics) != 1 {
		t.Fatalf("got %d diagnostics", len(result.Diagnostics))
	}
	if result.Diagnostics[0].Code != "VXT_WRITE_FILE_EXISTS" {
		t.Fatalf("got diagnostic code %q", result.Diagnostics[0].Code)
	}
}

func TestWritePlanWithDiagnosticsReturnsPathEscapeDiagnostic(t *testing.T) {
	p := plan.Plan{
		Files: []plan.FileOutput{{Path: "../escape.txt", Content: "nope", Mode: "create"}},
	}

	target := write.NewFilesystemTarget(t.TempDir())
	result := runtime.WritePlanWithDiagnostics(p, target)
	if result.Err == nil {
		t.Fatal("expected write error")
	}
	if len(result.Diagnostics) != 1 {
		t.Fatalf("got %d diagnostics", len(result.Diagnostics))
	}
	if result.Diagnostics[0].Code != "VXT_WRITE_PATH_ESCAPE" {
		t.Fatalf("got diagnostic code %q", result.Diagnostics[0].Code)
	}
}

func TestWritePlanWithDiagnosticsReturnsUnsupportedWriteModeDiagnostic(t *testing.T) {
	p := plan.Plan{
		Files: []plan.FileOutput{{Path: "hello.txt", Content: "nope", Mode: "merge"}},
	}

	target := write.NewMemoryTarget()
	result := runtime.WritePlanWithDiagnostics(p, target)
	if result.Err == nil {
		t.Fatal("expected write error")
	}
	if len(result.Diagnostics) != 1 {
		t.Fatalf("got %d diagnostics", len(result.Diagnostics))
	}
	if result.Diagnostics[0].Code != "VXT_WRITE_UNSUPPORTED_MODE" {
		t.Fatalf("got diagnostic code %q", result.Diagnostics[0].Code)
	}
}
