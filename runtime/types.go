package runtime

import (
	"github.com/vandordev/vxt/source"
	"github.com/vandordev/vxt/write"
)

// InputDecl declares one named document input and its required type.
type InputDecl struct {
	Name     string
	TypeName string
}

// TypeFieldDecl declares one field inside a named document type.
type TypeFieldDecl struct {
	Name     string
	TypeName string
	Optional bool
	Array    bool
}

// TypeDecl declares one named document input type.
type TypeDecl struct {
	Name   string
	Fields []TypeFieldDecl
}

// FileBlock defines one document file artifact before planning.
type FileBlock struct {
	Path string
	Body string
	Mode string
}

// DirBlock defines one directory artifact before planning.
type DirBlock struct {
	Path string
}

// PartialDecl defines one reusable partial body inside a document template.
type PartialDecl struct {
	Name string
	Body string
}

// UseDecl references one imported definition document.
type UseDecl struct {
	Path string
}

// ConditionalBlock defines one conditional document section.
type ConditionalBlock struct {
	Condition string
	Files     []FileBlock
	Dirs      []DirBlock
}

// HookDecl records one declared hook in document mode.
type HookDecl struct {
	Event string
	Run   string
}

// CompiledTemplate is the compiled single-file template contract exposed by runtime.
type CompiledTemplate struct {
	Source source.Source
}

// CompiledDocument is the compiled document template contract exposed by runtime.
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

// DirOutput is one planned directory output.
type DirOutput struct {
	Path string
}

// FileOutput is one planned file output.
type FileOutput struct {
	Path    string
	Content string
	Mode    string
}

// PlannedHook is one hook surfaced as planned metadata.
type PlannedHook struct {
	Event string
	Run   string
}

// Plan is the public runtime plan contract for document generation.
type Plan struct {
	Dirs         []DirOutput
	Files        []FileOutput
	PlannedHooks []PlannedHook
}

// HookContext provides structured context for one hook execution.
type HookContext struct {
	Event       string
	Plan        Plan
	WriteReport write.WriteReport
}

// HookExecutor executes one planned hook with explicit context.
type HookExecutor interface {
	Execute(ctx HookContext, hook PlannedHook) error
}

// ApplyResult captures write and hook-execution outcomes for one applied plan.
type ApplyResult struct {
	WriteResult WriteResult
	HookErrors  []error
}
