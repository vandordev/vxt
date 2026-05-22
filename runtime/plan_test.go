package runtime_test

import (
	"testing"

	"github.com/alfariiizi/vxt/runtime"
	"github.com/alfariiizi/vxt/source"
)

func TestPlanDocumentReturnsFileArtifactAndPlannedHooks(t *testing.T) {
	src := source.Source{
		ID: "plan-doc.vxt",
		Text: "@template hello\n" +
			"@input name string\n" +
			"@hook after:write \"echo later\"\n" +
			"@file \"hello.txt\"\n" +
			"Hello {{ name }}\n" +
			"@endfile\n",
	}

	compiled := runtime.CompileDocument(src)
	validated := runtime.ValidateDocument(compiled.Document, map[string]any{"name": "Fariz"})
	result := runtime.PlanDocument(validated)

	if len(result.Plan.Files) != 1 {
		t.Fatal("expected one planned file")
	}
	if len(result.Plan.PlannedHooks) != 1 {
		t.Fatal("expected planned hook metadata")
	}
}
