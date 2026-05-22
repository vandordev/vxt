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
