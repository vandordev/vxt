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

func TestPlanDocumentReturnsDirectoryArtifacts(t *testing.T) {
	src := source.Source{
		ID: "plan-dir-doc.vxt",
		Text: "@template hello\n" +
			"@input entity_name string\n" +
			"@dir \"src/modules/{{ entity_name }}\"\n",
	}

	compiled := runtime.CompileDocument(src)
	validated := runtime.ValidateDocument(compiled.Document, map[string]any{"entity_name": "booking"})
	result := runtime.PlanDocument(validated)

	if len(result.Plan.Dirs) != 1 {
		t.Fatalf("got %d dirs", len(result.Plan.Dirs))
	}
	if result.Plan.Dirs[0].Path != "src/modules/booking" {
		t.Fatalf("got dir path %q", result.Plan.Dirs[0].Path)
	}
}

func TestPlanDocumentRendersLocalPartialIncludes(t *testing.T) {
	src := source.Source{
		ID: "plan-partial-doc.vxt",
		Text: "@template demo\n" +
			"@partial imports\n" +
			"import \"context\"\n" +
			"@endpartial\n" +
			"@file \"demo.go\"\n" +
			"{{ include imports }}\n" +
			"package demo\n" +
			"@endfile\n",
	}

	compiled := runtime.CompileDocument(src)
	validated := runtime.ValidateDocument(compiled.Document, map[string]any{})
	result := runtime.PlanDocument(validated)

	if len(result.Diagnostics) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", result.Diagnostics)
	}
	if len(result.Plan.Files) != 1 {
		t.Fatalf("got %d files", len(result.Plan.Files))
	}
	if result.Plan.Files[0].Content != "import \"context\"\npackage demo" {
		t.Fatalf("got content %q", result.Plan.Files[0].Content)
	}
}
