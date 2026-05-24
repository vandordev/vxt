package bind_test

import (
	"testing"

	"github.com/vandordev/vxt/bind"
	"github.com/vandordev/vxt/source"
)

func TestGenerateReturnsOneGoFileForMinimalDocument(t *testing.T) {
	src := source.Source{
		ID: "hello.vxt",
		Text: "@template hello\n" +
			"@input name string\n" +
			"@file \"hello.txt\"\n" +
			"Hello {{ name }}\n" +
			"@endfile\n",
	}

	out, err := bind.Generate(bind.Request{
		PackageName: "hellovxt",
		Document:    src,
	})
	if err != nil {
		t.Fatalf("unexpected generate error: %v", err)
	}
	if len(out.Files) != 1 {
		t.Fatalf("got %d files", len(out.Files))
	}
	if out.Files[0].Path != ".vxt/hello_gen.go" {
		t.Fatalf("got file path %q", out.Files[0].Path)
	}
}
