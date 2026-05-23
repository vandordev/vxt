package write

type OutputTarget interface {
	MkdirAll(path string) error
	WriteFile(path string, content []byte, mode string) (bool, error)
}

type WriteReport struct {
	FilesWritten int
	FilesSkipped int
	DirsWritten  int
}
