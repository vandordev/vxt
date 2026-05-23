package runtime

import (
	internalmodel "github.com/vandordev/vxt/internal/model"
)

func compiledDocumentFromInternal(doc *internalmodel.CompiledDocument) *CompiledDocument {
	if doc == nil {
		return nil
	}

	return &CompiledDocument{
		Source:       doc.Source,
		Template:     doc.Template,
		Types:        typeDeclsFromInternal(doc.Types),
		Inputs:       inputDeclsFromInternal(doc.Inputs),
		Dirs:         dirBlocksFromInternal(doc.Dirs),
		Partials:     partialDeclsFromInternal(doc.Partials),
		Uses:         useDeclsFromInternal(doc.Uses),
		Conditionals: conditionalBlocksFromInternal(doc.Conditionals),
		Files:        fileBlocksFromInternal(doc.Files),
		Hooks:        hookDeclsFromInternal(doc.Hooks),
	}
}

func typeDeclsFromInternal(decls []internalmodel.TypeDecl) []TypeDecl {
	out := make([]TypeDecl, 0, len(decls))
	for _, decl := range decls {
		out = append(out, TypeDecl{
			Name:   decl.Name,
			Fields: typeFieldDeclsFromInternal(decl.Fields),
		})
	}
	return out
}

func typeFieldDeclsFromInternal(fields []internalmodel.TypeFieldDecl) []TypeFieldDecl {
	out := make([]TypeFieldDecl, 0, len(fields))
	for _, field := range fields {
		out = append(out, TypeFieldDecl{
			Name:     field.Name,
			TypeName: field.TypeName,
			Optional: field.Optional,
			Array:    field.Array,
		})
	}
	return out
}

func inputDeclsFromInternal(decls []internalmodel.InputDecl) []InputDecl {
	out := make([]InputDecl, 0, len(decls))
	for _, decl := range decls {
		out = append(out, InputDecl{
			Name:     decl.Name,
			TypeName: decl.TypeName,
		})
	}
	return out
}

func dirBlocksFromInternal(blocks []internalmodel.DirBlock) []DirBlock {
	out := make([]DirBlock, 0, len(blocks))
	for _, block := range blocks {
		out = append(out, DirBlock{Path: block.Path})
	}
	return out
}

func partialDeclsFromInternal(decls []internalmodel.PartialDecl) []PartialDecl {
	out := make([]PartialDecl, 0, len(decls))
	for _, decl := range decls {
		out = append(out, PartialDecl{
			Name: decl.Name,
			Body: decl.Body,
		})
	}
	return out
}

func useDeclsFromInternal(decls []internalmodel.UseDecl) []UseDecl {
	out := make([]UseDecl, 0, len(decls))
	for _, decl := range decls {
		out = append(out, UseDecl{Path: decl.Path})
	}
	return out
}

func conditionalBlocksFromInternal(blocks []internalmodel.ConditionalBlock) []ConditionalBlock {
	out := make([]ConditionalBlock, 0, len(blocks))
	for _, block := range blocks {
		out = append(out, ConditionalBlock{
			Condition: block.Condition,
			Files:     fileBlocksFromInternal(block.Files),
			Dirs:      dirBlocksFromInternal(block.Dirs),
		})
	}
	return out
}

func fileBlocksFromInternal(blocks []internalmodel.FileBlock) []FileBlock {
	out := make([]FileBlock, 0, len(blocks))
	for _, block := range blocks {
		out = append(out, FileBlock{
			Path: block.Path,
			Body: block.Body,
			Mode: block.Mode,
		})
	}
	return out
}

func hookDeclsFromInternal(decls []internalmodel.HookDecl) []HookDecl {
	out := make([]HookDecl, 0, len(decls))
	for _, decl := range decls {
		out = append(out, HookDecl{
			Event: decl.Event,
			Run:   decl.Run,
		})
	}
	return out
}

func typeDeclsToInternal(decls []TypeDecl) []internalmodel.TypeDecl {
	out := make([]internalmodel.TypeDecl, 0, len(decls))
	for _, decl := range decls {
		out = append(out, internalmodel.TypeDecl{
			Name:   decl.Name,
			Fields: typeFieldDeclsToInternal(decl.Fields),
		})
	}
	return out
}

func typeFieldDeclsToInternal(fields []TypeFieldDecl) []internalmodel.TypeFieldDecl {
	out := make([]internalmodel.TypeFieldDecl, 0, len(fields))
	for _, field := range fields {
		out = append(out, internalmodel.TypeFieldDecl{
			Name:     field.Name,
			TypeName: field.TypeName,
			Optional: field.Optional,
			Array:    field.Array,
		})
	}
	return out
}

func fileBlockToInternal(block FileBlock) internalmodel.FileBlock {
	return internalmodel.FileBlock{
		Path: block.Path,
		Body: block.Body,
		Mode: block.Mode,
	}
}
