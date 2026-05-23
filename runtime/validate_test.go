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

func TestValidateDocumentInputAgainstLocalNamedType(t *testing.T) {
	doc := &model.CompiledDocument{
		Types: []model.TypeDecl{{
			Name: "Entity",
			Fields: []model.TypeFieldDecl{
				{Name: "name", TypeName: "string"},
			},
		}},
		Inputs: []model.InputDecl{{Name: "entity", TypeName: "Entity"}},
	}

	result := runtime.ValidateDocument(doc, map[string]any{
		"entity": map[string]any{"name": "Booking"},
	})

	if len(result.Diagnostics) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", result.Diagnostics)
	}
}

func TestValidateDocumentInputAgainstLocalNamedTypeRequiresField(t *testing.T) {
	doc := &model.CompiledDocument{
		Types: []model.TypeDecl{{
			Name: "Entity",
			Fields: []model.TypeFieldDecl{
				{Name: "name", TypeName: "string"},
			},
		}},
		Inputs: []model.InputDecl{{Name: "entity", TypeName: "Entity"}},
	}

	result := runtime.ValidateDocument(doc, map[string]any{
		"entity": map[string]any{},
	})

	if len(result.Diagnostics) == 0 {
		t.Fatal("expected missing field diagnostic")
	}
}

func TestValidateDocumentInputAgainstImportedNamedType(t *testing.T) {
	doc := &model.CompiledDocument{
		Types: []model.TypeDecl{{
			Name: "Entity",
			Fields: []model.TypeFieldDecl{
				{Name: "name", TypeName: "string"},
			},
		}},
		Inputs: []model.InputDecl{{Name: "entity", TypeName: "Entity"}},
	}

	result := runtime.ValidateDocument(doc, map[string]any{
		"entity": map[string]any{"name": "ImportedBooking"},
	})

	if len(result.Diagnostics) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", result.Diagnostics)
	}
}
