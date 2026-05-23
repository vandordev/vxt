package runtime

import (
	internalmodel "github.com/vandordev/vxt/internal/model"
	internalplan "github.com/vandordev/vxt/internal/plan"
)

// InputDecl declares one named document input and its required type.
type InputDecl = internalmodel.InputDecl

// TypeFieldDecl declares one field inside a named document type.
type TypeFieldDecl = internalmodel.TypeFieldDecl

// TypeDecl declares one named document input type.
type TypeDecl = internalmodel.TypeDecl

// FileBlock defines one document file artifact before planning.
type FileBlock = internalmodel.FileBlock

// DirBlock defines one directory artifact before planning.
type DirBlock = internalmodel.DirBlock

// PartialDecl defines one reusable partial body inside a document template.
type PartialDecl = internalmodel.PartialDecl

// UseDecl references one imported definition document.
type UseDecl = internalmodel.UseDecl

// ConditionalBlock defines one conditional document section.
type ConditionalBlock = internalmodel.ConditionalBlock

// HookDecl records one declared hook in document mode.
type HookDecl = internalmodel.HookDecl

// CompiledTemplate is the compiled single-file template contract exposed by runtime.
type CompiledTemplate = internalmodel.CompiledTemplate

// CompiledDocument is the compiled document template contract exposed by runtime.
type CompiledDocument = internalmodel.CompiledDocument

// DirOutput is one planned directory output.
type DirOutput = internalplan.DirOutput

// FileOutput is one planned file output.
type FileOutput = internalplan.FileOutput

// PlannedHook is one hook surfaced as planned metadata.
type PlannedHook = internalplan.PlannedHook

// Plan is the public runtime plan contract for document generation.
type Plan = internalplan.Plan
