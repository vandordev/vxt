package diag

import "github.com/alfariiizi/vxt/source"

type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
)

// Diagnostic is a structured issue emitted by vxt.
type Diagnostic struct {
	Code     Code
	Severity Severity
	Message  string
	Span     source.Span
	Hint     string
}
