package model

import "github.com/alfariiizi/vxt/source"

type InputDecl struct {
	Name     string
	TypeName string
}

type FileBlock struct {
	Path string
	Body string
	Mode string
}

type HookDecl struct {
	Event string
	Run   string
}

// CompiledDocument is the stable public object for document mode.
type CompiledDocument struct {
	Source   source.Source
	Template string
	Inputs   []InputDecl
	Files    []FileBlock
	Hooks    []HookDecl
}
