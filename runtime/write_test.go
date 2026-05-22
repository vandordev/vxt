package runtime_test

import (
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
