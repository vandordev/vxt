package source

// Span identifies a location within a source.
type Span struct {
	SourceID string
	Start    int
	End      int
	Line     int
	Column   int
}
