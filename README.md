# vxt

`vxt` is a spec-first Go library for code and file generation under the Vandor
organization.

It exposes a staged pipeline:

1. compile
2. validate
3. plan
4. write

`vxt` is a library package today. It does not ship a CLI or standalone binary.

## Release Status

`vxt` `v0.1.0` is the first public Go package release and should be treated as
an experimental `v0.x` API. The supported public packages are deliberate, but
controlled breaking changes may still happen before `v1.0.0`.

## Install

```bash
go get github.com/vandordev/vxt
```

## Stable Public Packages

The intended public contract for `v0.1` is limited to:

- `github.com/vandordev/vxt`
- `github.com/vandordev/vxt/runtime`
- `github.com/vandordev/vxt/diag`
- `github.com/vandordev/vxt/source`
- `github.com/vandordev/vxt/write`

Implementation packages under `internal/` are not public API and may change
without notice.

## Quick Start

### Single-file rendering

```go
package main

import (
	"fmt"
	
	"github.com/vandordev/vxt"
	"github.com/vandordev/vxt/source"
)

func main() {
	out, diags := vxt.RenderSingleFile(source.Source{
		ID:   "hello.vxt",
		Text: "Hello {{ name }}",
	}, map[string]any{
		"name": "Vandor",
	})
	if len(diags) > 0 {
		panic(diags[0].Message)
	}

	fmt.Println(out)
}
```

### Document pipeline

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
		ID: "hello-doc.vxt",
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
}
```

## v0.1 Scope

The current `v0.1` target is intentionally narrow:

- single-file rendering through a convenience API
- document-mode compile, validate, plan, and write
- typed document input validation
- structured diagnostics
- output-target abstraction with filesystem and memory adapters

Current non-goals:

- hook execution
- trust policy
- registry or package semantics
- CLI behavior
- AST manipulation as a public contract

Hooks are surfaced only as planned metadata in document plans. They are not
executed by `vxt` in `v0.1`.

See [docs/releases/v0.1.0.md](docs/releases/v0.1.0.md) for the curated release
scope and verification checklist.

## License and Trademark

The source code in this repository is licensed under `AGPL-3.0`. See
[LICENSE](LICENSE).

`Vandor` name and brand assets are not licensed under the AGPL. See
[TRADEMARK.md](TRADEMARK.md).
