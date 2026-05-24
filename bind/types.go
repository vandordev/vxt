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
	BindingName string
	PackageName string
	Files       []File
}

// WriteOptions controls bind output write behavior.
type WriteOptions struct {
	DryRun bool
}

// WriteActionKind describes one filesystem change kind.
type WriteActionKind string

const (
	WriteActionCreate    WriteActionKind = "create"
	WriteActionOverwrite WriteActionKind = "overwrite"
	WriteActionRemove    WriteActionKind = "remove"
)

// WriteAction records one concrete or planned file action.
type WriteAction struct {
	Path   string
	Action WriteActionKind
}

// WriteReport describes one bind output write operation.
type WriteReport struct {
	DryRun           bool
	OutputDir        string
	FilesWritten     []string
	FilesOverwritten []string
	FilesRemoved     []string
	Actions          []WriteAction
}
