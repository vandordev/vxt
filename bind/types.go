package bind

import "github.com/vandordev/vxt/source"

// Request defines one binding generation request.
type Request struct {
	PackageName string
	Document    source.Source
	Uses        map[string]source.Source
}

// File is one generated file artifact.
type File struct {
	Path    string
	Content string
}

// Output is the generated file set for one request.
type Output struct {
	Files []File
}
