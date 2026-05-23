package diag

import "github.com/vandordev/vxt/source"

// Severity describes the impact level of one diagnostic.
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
