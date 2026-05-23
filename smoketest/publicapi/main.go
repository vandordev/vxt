package main

import (
	"fmt"
	"os"

	"github.com/vandordev/vxt/runtime"
	"github.com/vandordev/vxt/source"
	"github.com/vandordev/vxt/write"
)

func main() {
	src := source.Source{
		ID: "release-smoke.vxt",
		Text: "@template smoke\n" +
			"@input name string\n" +
			"@file \"hello.txt\"\n" +
			"Hello {{ name }}\n" +
			"@endfile\n",
	}

	compiled := runtime.CompileDocument(src)
	if len(compiled.Diagnostics) > 0 {
		fail("compile", compiled.Diagnostics[0].Message)
	}

	validated := runtime.ValidateDocument(compiled.Document, map[string]any{
		"name": "Vandor",
	})
	if len(validated.Diagnostics) > 0 {
		fail("validate", validated.Diagnostics[0].Message)
	}

	planned := runtime.PlanDocument(validated)
	if len(planned.Diagnostics) > 0 {
		fail("plan", planned.Diagnostics[0].Message)
	}
	if len(planned.Plan.Files) != 1 {
		fail("plan", fmt.Sprintf("expected 1 planned file, got %d", len(planned.Plan.Files)))
	}

	target := write.NewMemoryTarget()
	result := runtime.WritePlanWithDiagnostics(planned.Plan, target)
	if result.Err != nil {
		fail("write", result.Err.Error())
	}
	if result.Report.FilesWritten != 1 {
		fail("write", fmt.Sprintf("expected 1 written file, got %d", result.Report.FilesWritten))
	}

	content, ok := target.Files()["hello.txt"]
	if !ok {
		fail("write", "expected hello.txt in memory target")
	}
	if got := string(content); got != "Hello Vandor" {
		fail("write", fmt.Sprintf("expected rendered content %q, got %q", "Hello Vandor", got))
	}

	fmt.Println("public smoke test ok")
}

func fail(stage, msg string) {
	fmt.Fprintf(os.Stderr, "%s failed: %s\n", stage, msg)
	os.Exit(1)
}
