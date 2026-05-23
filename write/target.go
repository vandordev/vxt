package write

// OutputTarget receives planned filesystem operations from the runtime write stage.
type OutputTarget interface {
	MkdirAll(path string) error
	WriteFile(path string, content []byte, mode string) (bool, error)
}

// WriteReport summarizes concrete write-stage outcomes.
type WriteReport struct {
	FilesWritten int
	FilesSkipped int
	DirsWritten  int
}
