package diag

// Code identifies one stable diagnostic category emitted by vxt.
type Code string

const (
	CodeParseUnexpectedEOF   Code = "VXT_PARSE_UNEXPECTED_EOF"
	CodeTypeMismatch         Code = "VXT_TYPE_MISMATCH"
	CodeRenderMissingValue   Code = "VXT_RENDER_MISSING_VALUE"
	CodeOutputConflict       Code = "VXT_OUTPUT_CONFLICT"
	CodeWriteFileExists      Code = "VXT_WRITE_FILE_EXISTS"
	CodeWritePathEscape      Code = "VXT_WRITE_PATH_ESCAPE"
	CodeWriteUnsupportedMode Code = "VXT_WRITE_UNSUPPORTED_MODE"
)
