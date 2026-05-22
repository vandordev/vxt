package render_test

import (
	"testing"

	"github.com/alfariiizi/vxt/model"
	"github.com/alfariiizi/vxt/render"
)

func TestRenderDocumentBodyInterpolatesInput(t *testing.T) {
	body, diags := render.RenderDocumentBody(model.FileBlock{
		Path: "hello.txt",
		Body: "Hello {{ name }}",
	}, map[string]any{"name": "Fariz"})

	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
	if body != "Hello Fariz" {
		t.Fatalf("got %q", body)
	}
}
