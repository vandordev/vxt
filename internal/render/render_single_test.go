package render_test

import (
	"testing"

	"github.com/vandordev/vxt"
	"github.com/vandordev/vxt/source"
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

func TestRenderSingleFileAppliesCaseFilters(t *testing.T) {
	src := source.Source{
		ID:   "case-filter.vxt",
		Text: "type {{ name | pascal }}Service struct {\n\tfield {{ name | camel }}\n\tkey string // {{ name | snake }}\n}",
	}

	out, diags := vxt.RenderSingleFile(src, map[string]any{"name": "order item"})

	if len(diags) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
	want := "type OrderItemService struct {\n\tfield orderItem\n\tkey string // order_item\n}"
	if out != want {
		t.Fatalf("got %q", out)
	}
}
