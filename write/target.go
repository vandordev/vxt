package write

type OutputTarget interface {
	WriteFile(path string, content []byte) error
}

type WriteReport struct {
	FilesWritten int
}
