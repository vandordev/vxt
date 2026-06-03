package expr

import "testing"

func TestEvalPathAppliesCaseFilters(t *testing.T) {
	tests := map[string]string{
		"name | snake":       "order_item",
		"name | upper_snake": "ORDER_ITEM",
		"name | kebab":       "order-item",
		"name | pascal":      "OrderItem",
		"name | camel":       "orderItem",
	}

	for expr, want := range tests {
		got, err := EvalPath(map[string]any{"name": "OrderItem"}, expr)
		if err != nil {
			t.Fatalf("EvalPath(%q) error: %v", expr, err)
		}
		if got != want {
			t.Fatalf("EvalPath(%q) = %q, want %q", expr, got, want)
		}
	}
}

func TestEvalPathAppliesCaseFiltersToSeparatedInput(t *testing.T) {
	got, err := EvalPath(map[string]any{"name": "order item"}, "name | pascal")
	if err != nil {
		t.Fatalf("EvalPath error: %v", err)
	}
	if got != "OrderItem" {
		t.Fatalf("got %q", got)
	}
}
