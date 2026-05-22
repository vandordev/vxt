package schema_test

import (
	"testing"

	"github.com/alfariiizi/vxt/schema"
)

func TestValidateValueRejectsWrongPrimitiveType(t *testing.T) {
	if err := schema.ValidateValue(schema.TypeString, 42); err == nil {
		t.Fatal("expected validation error")
	}
}
