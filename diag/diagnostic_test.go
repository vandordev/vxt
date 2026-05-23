package diag_test

import (
	"testing"

	"github.com/vandordev/vxt/diag"
	"github.com/vandordev/vxt/source"
)

func TestDiagnosticIncludesCodeSeverityAndSpan(t *testing.T) {
	src := source.Source{ID: "basic.vxt", Text: "{{ name }}"}
	d := diag.Diagnostic{
		Code:     diag.CodeRenderMissingValue,
		Severity: diag.SeverityError,
		Message:  "missing value",
		Span:     source.Span{SourceID: src.ID, Start: 3, End: 7},
	}

	if d.Code == "" || d.Severity == "" || d.Span.SourceID == "" {
		t.Fatal("diagnostic contract incomplete")
	}
}
