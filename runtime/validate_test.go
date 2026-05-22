package runtime_test

import (
	"testing"

	"github.com/alfariiizi/vxt/model"
	"github.com/alfariiizi/vxt/runtime"
)

func TestValidateDocumentInputRequiresDeclaredString(t *testing.T) {
	doc := &model.CompiledDocument{
		Inputs: []model.InputDecl{{Name: "name", TypeName: "string"}},
	}

	result := runtime.ValidateDocument(doc, map[string]any{})

	if len(result.Diagnostics) == 0 {
		t.Fatal("expected missing-input diagnostic")
	}
}
