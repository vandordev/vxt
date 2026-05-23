package runtime_test

import (
	"testing"

	"github.com/alfariiizi/vxt/runtime"
	"github.com/alfariiizi/vxt/source"
)

func TestCompileSingleFileReturnsCompiledTemplate(t *testing.T) {
	src := source.Source{ID: "basic.vxt", Text: "Hello {{ name }}"}

	result := runtime.CompileSingle(src)

	if result.Template == nil {
		t.Fatal("expected compiled template")
	}
	if len(result.Diagnostics) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", result.Diagnostics)
	}
}

func TestCompileDocumentReturnsCompiledDocument(t *testing.T) {
	src := source.Source{
		ID:   "basic-doc.vxt",
		Text: "@template hello\n@input name string\n@file \"hello.txt\"\nHello {{ name }}\n@endfile\n",
	}

	result := runtime.CompileDocument(src)

	if result.Document == nil {
		t.Fatal("expected compiled document")
	}
}

func TestCompileDocumentParsesLocalTypeDeclarations(t *testing.T) {
	src := source.Source{
		ID: "type-basic.vxt",
		Text: "@template demo\n" +
			"@type Entity {\n" +
			"  name: string\n" +
			"}\n" +
			"@input entity Entity\n",
	}

	result := runtime.CompileDocument(src)

	if result.Document == nil {
		t.Fatal("expected compiled document")
	}
	if len(result.Document.Types) != 1 {
		t.Fatalf("got %d types", len(result.Document.Types))
	}
	if result.Document.Types[0].Name != "Entity" {
		t.Fatalf("got type %q", result.Document.Types[0].Name)
	}
}

func TestCompileDocumentParsesDirectoryDeclarations(t *testing.T) {
	src := source.Source{
		ID:   "dir-basic.vxt",
		Text: "@template demo\n@dir \"src/modules/{{ entity_name }}\"\n",
	}

	result := runtime.CompileDocument(src)

	if result.Document == nil {
		t.Fatal("expected compiled document")
	}
	if len(result.Document.Dirs) != 1 {
		t.Fatalf("got %d dirs", len(result.Document.Dirs))
	}
	if result.Document.Dirs[0].Path != "src/modules/{{ entity_name }}" {
		t.Fatalf("got dir path %q", result.Document.Dirs[0].Path)
	}
}

func TestCompileDocumentParsesLocalPartialDeclarations(t *testing.T) {
	src := source.Source{
		ID: "partial-basic.vxt",
		Text: "@template demo\n" +
			"@partial imports\n" +
			"import \"context\"\n" +
			"@endpartial\n" +
			"@file \"demo.go\"\n" +
			"{{ include imports }}\n" +
			"package demo\n" +
			"@endfile\n",
	}

	result := runtime.CompileDocument(src)

	if result.Document == nil {
		t.Fatal("expected compiled document")
	}
	if len(result.Document.Partials) != 1 {
		t.Fatalf("got %d partials", len(result.Document.Partials))
	}
	if result.Document.Partials[0].Name != "imports" {
		t.Fatalf("got partial %q", result.Document.Partials[0].Name)
	}
}

func TestCompileDocumentWithResolverLoadsLocalUseSource(t *testing.T) {
	main := source.Source{
		ID: "use-main.vxt",
		Text: "@template demo\n" +
			"@use \"./use_schema.vxt\"\n" +
			"@input entity Entity\n",
	}

	resolver := runtime.MapResolver(map[string]source.Source{
		"./use_schema.vxt": {
			ID: "use_schema.vxt",
			Text: "@type Entity {\n" +
				"  name: string\n" +
				"}\n",
		},
	})

	result := runtime.CompileDocumentWithResolver(main, resolver)
	if result.Document == nil {
		t.Fatal("expected compiled document")
	}
	if len(result.Document.Types) != 1 {
		t.Fatalf("got %d types", len(result.Document.Types))
	}
	if result.Document.Types[0].Name != "Entity" {
		t.Fatalf("got type %q", result.Document.Types[0].Name)
	}
}
