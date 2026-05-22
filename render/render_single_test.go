package render_test

import (
	"testing"

	"github.com/alfariiizi/vxt"
	"github.com/alfariiizi/vxt/source"
)

func TestRenderSingleFileInterpolatesValue(t *testing.T) {
	src := source.Source{ID: "basic.vxt", Text: "Hello {{ name }}"}

	out, diags := vxt.RenderSingleFile(src, map[string]any{"name": "Fariz"})

	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
	if out != "Hello Fariz" {
		t.Fatalf("got %q", out)
	}
}
