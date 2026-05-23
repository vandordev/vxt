package runtime_test

import (
	"testing"

	"github.com/vandordev/vxt/runtime"
)

func TestValidateDocumentInputRequiresDeclaredString(t *testing.T) {
	doc := &runtime.CompiledDocument{
		Inputs: []runtime.InputDecl{{Name: "name", TypeName: "string"}},
	}

	result := runtime.ValidateDocument(doc, map[string]any{})

	if len(result.Diagnostics) == 0 {
		t.Fatal("expected missing-input diagnostic")
	}
}

func TestValidateDocumentInputAgainstLocalNamedType(t *testing.T) {
	doc := &runtime.CompiledDocument{
		Types: []runtime.TypeDecl{{
			Name: "Entity",
			Fields: []runtime.TypeFieldDecl{
				{Name: "name", TypeName: "string"},
			},
		}},
		Inputs: []runtime.InputDecl{{Name: "entity", TypeName: "Entity"}},
	}

	result := runtime.ValidateDocument(doc, map[string]any{
		"entity": map[string]any{"name": "Booking"},
	})

	if len(result.Diagnostics) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", result.Diagnostics)
	}
}

func TestValidateDocumentInputAgainstLocalNamedTypeRequiresField(t *testing.T) {
	doc := &runtime.CompiledDocument{
		Types: []runtime.TypeDecl{{
			Name: "Entity",
			Fields: []runtime.TypeFieldDecl{
				{Name: "name", TypeName: "string"},
			},
		}},
		Inputs: []runtime.InputDecl{{Name: "entity", TypeName: "Entity"}},
	}

	result := runtime.ValidateDocument(doc, map[string]any{
		"entity": map[string]any{},
	})

	if len(result.Diagnostics) == 0 {
		t.Fatal("expected missing field diagnostic")
	}
}

func TestValidateDocumentInputAgainstImportedNamedType(t *testing.T) {
	doc := &runtime.CompiledDocument{
		Types: []runtime.TypeDecl{{
			Name: "Entity",
			Fields: []runtime.TypeFieldDecl{
				{Name: "name", TypeName: "string"},
			},
		}},
		Inputs: []runtime.InputDecl{{Name: "entity", TypeName: "Entity"}},
	}

	result := runtime.ValidateDocument(doc, map[string]any{
		"entity": map[string]any{"name": "ImportedBooking"},
	})

	if len(result.Diagnostics) != 0 {
		t.Fatalf("unexpected diagnostics: %#v", result.Diagnostics)
	}
}
