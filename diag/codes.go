package diag

type Code string

const (
	CodeParseUnexpectedEOF Code = "VXT_PARSE_UNEXPECTED_EOF"
	CodeTypeMismatch       Code = "VXT_TYPE_MISMATCH"
	CodeRenderMissingValue Code = "VXT_RENDER_MISSING_VALUE"
	CodeOutputConflict     Code = "VXT_OUTPUT_CONFLICT"
)
