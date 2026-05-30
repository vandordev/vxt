# vxt

`vxt` is a spec-first Go library for turning `.vxt` templates into planned
file output. It is designed for applications and tools that want a typed,
inspectable generation pipeline instead of ad-hoc string rendering.

The main path is document mode: a `.vxt` document declares inputs, files,
directories, partials, optional sections, and hook metadata. Go code compiles
the document, validates input, plans output, and writes that plan to a target.

## When to Use It

Use `vxt` when you want to:

- embed code or file generation in a Go program
- validate template inputs before writing output
- inspect planned files and directories before touching disk
- write to memory in tests and to a rooted filesystem in production
- generate typed Go bindings for a document-mode template

## What It Is Not

`vxt` is library-only. It does not ship a CLI or standalone binary.

It also does not provide:

- trust policy
- package registry semantics
- automatic shell execution
- public AST manipulation APIs
- cross-language generated bindings

## Install

```bash
go get github.com/vandordev/vxt
```

## Quick Start: Document Mode

This example compiles a small `.vxt` document, validates its input, builds a
plan, and writes the result into memory.

```go
package main

import (
	"fmt"

	"github.com/vandordev/vxt/runtime"
	"github.com/vandordev/vxt/source"
	"github.com/vandordev/vxt/write"
)

func main() {
	src := source.Source{
		ID: "hello.vxt",
		Text: "@template hello\n" +
			"@input name string\n" +
			"@file \"hello.txt\"\n" +
			"Hello {{ name }}\n" +
			"@endfile\n",
	}

	compiled := runtime.CompileDocument(src)
	if len(compiled.Diagnostics) > 0 {
		panic(compiled.Diagnostics[0].Message)
	}

	validated := runtime.ValidateDocument(compiled.Document, map[string]any{
		"name": "Vandor",
	})
	if len(validated.Diagnostics) > 0 {
		panic(validated.Diagnostics[0].Message)
	}

	planned := runtime.PlanDocument(validated)
	if len(planned.Diagnostics) > 0 {
		panic(planned.Diagnostics[0].Message)
	}

	target := write.NewMemoryTarget()
	report, err := runtime.WritePlan(planned.Plan, target)
	if err != nil {
		panic(err)
	}

	fmt.Println(report.FilesWritten)
	fmt.Println(string(target.Files()["hello.txt"]))
}
```

Use `write.NewFilesystemTarget(root)` when you are ready to write the same plan
under a real output directory.

## Pipeline Overview

The document-mode runtime is intentionally staged:

1. `runtime.CompileDocument` parses a document and returns a compiled contract.
2. `runtime.ValidateDocument` checks caller-provided input against declared
   inputs and document types.
3. `runtime.PlanDocument` renders concrete directory, file, and hook metadata.
4. `runtime.WritePlan` writes planned directories and files to an
   `write.OutputTarget`.
5. `runtime.ApplyPlan` is optional. It writes the plan and then executes
   supported planned hooks through a caller-provided `runtime.HookExecutor`.

`WritePlan` does not execute hooks. Hooks are metadata unless the caller
explicitly chooses `ApplyPlan(...)` and provides an executor. The currently
supported hook event for `ApplyPlan` is `after:write`.

## Feature Snapshot

- document-mode templates with `@template`, `@input`, `@type`, `@dir`,
  `@file`, `@partial`, `@use`, `@if`, and `@hook`
- single-file rendering through `vxt.RenderSingleFile`
- structured diagnostics across compile, validate, plan, and write stages
- memory and filesystem output targets
- explicit post-write hook execution through `ApplyPlan`
- generated Go bindings through `github.com/vandordev/vxt/bind`

## Docs Map

- [Getting started](docs/getting-started.md): first document-mode flow with
  `MemoryTarget`, then filesystem output.
- [Document mode](docs/document-mode.md): `.vxt` authoring tutorial and
  directive reference.
- [Runtime API](docs/runtime-api.md): compile, validate, plan, write, and apply
  lifecycle details.
- [Go bindings](docs/go-bindings.md): generated typed Go package workflow and
  `bind` package usage.
- [Concepts](docs/concepts.md): product boundaries, `vx` relationship, runtime
  vs bindings, and hook model.
- [v0.1.0 release notes](docs/releases/v0.1.0.md): release scope and
  verification checklist.

## Public Packages and Release Status

`vxt` is in an experimental `v0.x` line. Public packages are intended for Go
consumers, but controlled breaking changes may still happen before `v1.0.0`.

Documented public packages include:

- `github.com/vandordev/vxt`
- `github.com/vandordev/vxt/runtime`
- `github.com/vandordev/vxt/bind`
- `github.com/vandordev/vxt/diag`
- `github.com/vandordev/vxt/source`
- `github.com/vandordev/vxt/write`

Implementation packages under `internal/` are not public API.

## License and Trademark

The source code in this repository is licensed under `AGPL-3.0`. See
[LICENSE](LICENSE).

`Vandor` name and brand assets are not licensed under the AGPL. See
[TRADEMARK.md](TRADEMARK.md).
