package bind

import (
	"strings"
	"testing"

	"github.com/vandordev/vxt/source"
)

func TestGenerateEmitsTypedWrapperPackageSurface(t *testing.T) {
	src := source.Source{
		ID: "service.vxt",
		Text: "@template service\n" +
			"@type Entity {\n" +
			"  name: string\n" +
			"  package_name: string\n" +
			"}\n" +
			"@input entity Entity\n" +
			"@file \"service.go\"\n" +
			"package {{ entity.package_name }}\n" +
			"@endfile\n",
	}

	out, err := Generate(Request{
		PackageName: "servicevxt",
		Document:    src,
	})
	if err != nil {
		t.Fatalf("unexpected generate error: %v", err)
	}
	if len(out.Files) != 1 {
		t.Fatalf("got %d files", len(out.Files))
	}

	content := out.Files[0].Content
	assertContains(t, content, "package servicevxt")
	assertContains(t, content, "type Entity struct")
	assertContains(t, content, "type Input struct")
	assertContains(t, content, "func Compile() runtime.CompileResult")
	assertContains(t, content, "func Validate(input Input) runtime.ValidationResult")
	assertContains(t, content, "func Plan(input Input) (runtime.Plan, error)")
	assertContains(t, content, "func Write(input Input, target write.OutputTarget) (write.WriteReport, error)")
	assertContains(t, content, "func PlanDetailed(input Input) runtime.PlanResult")
	assertContains(t, content, "func WriteDetailed(input Input, target write.OutputTarget) runtime.WriteResult")
}

func assertContains(t *testing.T, content, want string) {
	t.Helper()
	if !strings.Contains(content, want) {
		t.Fatalf("generated content missing %q\n\n%s", want, content)
	}
}
