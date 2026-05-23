package model

import "github.com/vandordev/vxt/source"

type InputDecl struct {
	Name     string
	TypeName string
}

type TypeFieldDecl struct {
	Name     string
	TypeName string
	Optional bool
	Array    bool
}

type TypeDecl struct {
	Name   string
	Fields []TypeFieldDecl
}

type FileBlock struct {
	Path string
	Body string
	Mode string
}

type DirBlock struct {
	Path string
}

type PartialDecl struct {
	Name string
	Body string
}

type UseDecl struct {
	Path string
}

type ConditionalBlock struct {
	Condition string
	Files     []FileBlock
	Dirs      []DirBlock
}

type HookDecl struct {
	Event string
	Run   string
}

// CompiledDocument is the stable public object for document mode.
type CompiledDocument struct {
	Source       source.Source
	Template     string
	Types        []TypeDecl
	Inputs       []InputDecl
	Dirs         []DirBlock
	Partials     []PartialDecl
	Uses         []UseDecl
	Conditionals []ConditionalBlock
	Files        []FileBlock
	Hooks        []HookDecl
}
