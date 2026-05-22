package runtime_test

import (
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
