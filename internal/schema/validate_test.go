package schema_test

import (
	"testing"

	"github.com/vandordev/vxt/internal/schema"
)

func TestValidateValueRejectsWrongPrimitiveType(t *testing.T) {
	if err := schema.ValidateValue(schema.TypeString, 42); err == nil {
		t.Fatal("expected validation error")
	}
}

func TestValidateValueAcceptsBoolAlias(t *testing.T) {
	if err := schema.ValidateValue("bool", true); err != nil {
		t.Fatalf("expected bool alias to validate, got %v", err)
	}
}
