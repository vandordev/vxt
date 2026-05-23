package write

type OutputTarget interface {
	MkdirAll(path string) error
	WriteFile(path string, content []byte) error
}

type WriteReport struct {
	FilesWritten int
	DirsWritten  int
}
