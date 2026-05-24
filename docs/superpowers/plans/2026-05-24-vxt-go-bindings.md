# VXT Go Bindings Implementation Plan

> **For agentic workers:** REQUIRED: Use `superpowers:subagent-driven-development` (if subagents available) or `superpowers:executing-plans` to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a library-only Go binding generator to `vxt` that turns one document-mode `.vxt` template into a self-contained, typed Go package under `.vxt/`.

**Architecture:** Implement this as a new `bind` package that analyzes one compiled document, derives Go-facing type information, and emits one generated Go file containing embedded sources plus typed wrapper functions over `runtime`. Keep the generator document-first, local-`@use`-aware, and self-contained without introducing any CLI behavior into `vxt`.

**Tech Stack:** Go, existing `runtime` and `source` packages, Go `text/template` or explicit string emission, Go `embed`-style generated string literals, `go test`, and temporary-directory end-to-end tests.

---

## Chunk 1: Generator Contract and Package Skeleton

### Task 1: Add the `bind` package contract before any implementation logic

**Files:**
- Create: `bind/doc.go`
- Create: `bind/generator.go`
- Create: `bind/types.go`
- Create: `bind/generator_test.go`
- Test: `bind/generator_test.go`

- [ ] **Step 1: Write the failing package-contract test**

Create `bind/generator_test.go` with one high-level test that expects the new package to generate a file for a minimal document:

```go
func TestGenerateReturnsOneGoFileForMinimalDocument(t *testing.T) {
	src := source.Source{
		ID: "hello.vxt",
		Text: "@template hello\n" +
			"@input name string\n" +
			"@file \"hello.txt\"\n" +
			"Hello {{ name }}\n" +
			"@endfile\n",
	}

	out, err := bind.Generate(bind.Request{
		PackageName: "hellovxt",
		Document:    src,
	})
	if err == nil {
		t.Fatal("expected generator to be unimplemented initially")
	}
	if len(out.Files) != 0 {
		t.Fatalf("got %d files", len(out.Files))
	}
}
```

- [ ] **Step 2: Run the new test to verify the package does not exist yet**

Run:

```bash
go test ./bind -run TestGenerateReturnsOneGoFileForMinimalDocument -v
```

Expected:
- FAIL with missing package or symbol errors

- [ ] **Step 3: Add the minimal public package skeleton**

Create `bind/doc.go`:

```go
// Package bind generates typed Go bindings from document-mode VXT templates.
package bind
```

Create `bind/types.go` with the first public contract:

```go
package bind

import "github.com/vandordev/vxt/source"

type Request struct {
	PackageName string
	Document    source.Source
	Uses        map[string]source.Source
}

type File struct {
	Path    string
	Content string
}

type Output struct {
	Files []File
}
```

Create `bind/generator.go`:

```go
package bind

import "errors"

var ErrNotImplemented = errors.New("bind: generator not implemented")

func Generate(req Request) (Output, error) {
	_ = req
	return Output{}, ErrNotImplemented
}
```

- [ ] **Step 4: Re-run the targeted test and make sure it passes at the skeleton stage**

Run:

```bash
go test ./bind -run TestGenerateReturnsOneGoFileForMinimalDocument -v
```

Expected:
- PASS because the test now observes the explicit unimplemented error

- [ ] **Step 5: Commit the package skeleton**

Run:

```bash
git add bind/doc.go bind/generator.go bind/types.go bind/generator_test.go
git commit -m "feat: add bind package skeleton"
```

## Chunk 2: Type Analysis and Go API Shape

### Task 2: Derive typed Go shapes from one compiled document

**Files:**
- Modify: `bind/generator.go`
- Modify: `bind/types.go`
- Create: `bind/analyze.go`
- Create: `bind/names.go`
- Create: `bind/analyze_test.go`
- Test: `bind/analyze_test.go`

- [ ] **Step 1: Write failing tests for type and input analysis**

Create `bind/analyze_test.go` with focused tests for:
- named `@type` becomes a public Go struct
- optional fields become pointers
- array fields become slices
- `@input` fields become one canonical root `Input` struct

Use a compiled document fixture like:

```go
src := source.Source{
	ID: "service.vxt",
	Text: "@template service\n" +
		"@type Repository {\n" +
		"  driver: string\n" +
		"}\n" +
		"@type Entity {\n" +
		"  name: string\n" +
		"  package_name: string\n" +
		"  tags: string[]\n" +
		"  repository?: Repository\n" +
		"}\n" +
		"@input entity Entity\n",
}
```

Assert analysis produces a model equivalent to:
- `Repository.Driver string`
- `Entity.Name string`
- `Entity.PackageName string`
- `Entity.Tags []string`
- `Entity.Repository *Repository`
- `Input.Entity Entity`

- [ ] **Step 2: Run the failing analysis tests**

Run:

```bash
go test ./bind -run 'TestAnalyze|TestGoType' -v
```

Expected:
- FAIL because analysis helpers do not exist

- [ ] **Step 3: Add the internal analysis model and naming helpers**

Create `bind/analyze.go` with internal structs such as:

```go
type analyzedDocument struct {
	PackageName string
	Types       []analyzedType
	InputFields []analyzedField
}

type analyzedType struct {
	Name   string
	Fields []analyzedField
}

type analyzedField struct {
	GoName     string
	SchemaName string
	GoType     string
}
```

Create `bind/names.go` with focused helpers:

```go
func toExportedGoName(name string) string
func goTypeForField(field runtime.TypeFieldDecl) string
func goTypeForInput(input runtime.InputDecl, types []runtime.TypeDecl) string
```

Use the existing public runtime types from `runtime/types.go` as the canonical source model for generation input.

- [ ] **Step 4: Implement minimal compilation-plus-analysis flow**

Update `bind/generator.go` so `Generate`:
- compiles the document with `runtime.CompileDocument` first
- returns an error if compile diagnostics exist
- runs internal analysis on `result.Document`
- still returns no emitted files yet

Do not emit Go code in this task; stop after analysis is available and tested.

- [ ] **Step 5: Re-run the analysis tests**

Run:

```bash
go test ./bind -run 'TestAnalyze|TestGoType' -v
```

Expected:
- PASS

- [ ] **Step 6: Commit the type-analysis chunk**

Run:

```bash
git add bind/generator.go bind/types.go bind/analyze.go bind/names.go bind/analyze_test.go
git commit -m "feat: analyze vxt documents for typed go bindings"
```

## Chunk 3: Go Source Emission for Typed Bindings

### Task 3: Emit one generated Go file with typed structs and wrapper signatures

**Files:**
- Modify: `bind/generator.go`
- Create: `bind/emit.go`
- Create: `bind/templates.go`
- Create: `bind/emit_test.go`
- Test: `bind/emit_test.go`

- [ ] **Step 1: Write failing emission tests**

Create `bind/emit_test.go` with assertions that generated content includes:
- `package <name>`
- public generated structs for named `@type`
- canonical `type Input struct`
- wrapper signatures:
  - `func Compile() runtime.CompileResult`
  - `func Validate(input Input) runtime.ValidationResult`
  - `func Plan(input Input) (runtime.Plan, error)`
  - `func Write(input Input, target write.OutputTarget) (write.WriteReport, error)`
  - `func PlanDetailed(input Input) runtime.PlanResult`
  - `func WriteDetailed(input Input, target write.OutputTarget) runtime.WriteResult`

- [ ] **Step 2: Run the emission tests and confirm they fail**

Run:

```bash
go test ./bind -run 'TestEmit|TestGenerate' -v
```

Expected:
- FAIL because emission code is missing

- [ ] **Step 3: Add a focused emitter**

Create `bind/emit.go` with one internal entry point such as:

```go
func emitGeneratedFile(doc analyzedDocument, assets embeddedAssets) (string, error)
```

Create `bind/templates.go` with either:
- one `text/template` template string
  or
- explicit string builder helpers

Prefer one small emission unit instead of spreading string fragments across multiple files.

- [ ] **Step 4: Generate the typed wrappers and runtime bridge stubs**

The emitted file should include:
- generated structs
- `Input`
- embedded source variables
- placeholder conversion helpers from typed `Input` to `map[string]any`
- wrapper functions that call `runtime`

At this step, implement the conversion minimally for:
- primitive scalar fields
- nested named types
- slices of strings and named types
- pointer optional fields

- [ ] **Step 5: Update `Generate` to return `.vxt/<template>_gen.go`**

Modify `Generate` so it now returns one file:

```go
Output{
	Files: []File{{
		Path:    ".vxt/service_gen.go",
		Content: generated,
	}},
}
```

Use a deterministic filename derived from the main document template name.

- [ ] **Step 6: Re-run emission tests**

Run:

```bash
go test ./bind -run 'TestEmit|TestGenerate' -v
```

Expected:
- PASS

- [ ] **Step 7: Commit the emission chunk**

Run:

```bash
git add bind/generator.go bind/emit.go bind/templates.go bind/emit_test.go
git commit -m "feat: emit typed go binding packages"
```

## Chunk 4: Embedded Sources and Local `@use` Closure

### Task 4: Make generated bindings self-contained with embedded local dependencies

**Files:**
- Modify: `bind/generator.go`
- Create: `bind/assets.go`
- Create: `bind/use_closure.go`
- Create: `bind/use_closure_test.go`
- Test: `bind/use_closure_test.go`

- [ ] **Step 1: Write failing tests for local `@use` embedding**

Create `bind/use_closure_test.go` with a fixture:
- main document uses `@use "./schema.vxt"`
- schema document defines `@type Entity`

Assert generated code includes:
- embedded main document source
- embedded imported source text or equivalent source map entry
- generated `Compile()` path that can resolve the imported local definition without disk access

- [ ] **Step 2: Run the local-`@use` tests and verify failure**

Run:

```bash
go test ./bind -run TestGenerateEmbedsLocalUseSources -v
```

Expected:
- FAIL because imported sources are not embedded yet

- [ ] **Step 3: Add closure collection for local `@use` sources**

Create `bind/use_closure.go` with a helper such as:

```go
func collectLocalUseSources(main source.Source, uses map[string]source.Source) (map[string]source.Source, error)
```

Rules:
- only local-path `@use` participates
- every declared local use must exist in `Request.Uses`
- keep the first version flat and deterministic

- [ ] **Step 4: Add embedded asset modeling**

Create `bind/assets.go` with internal modeling for:
- main document source
- imported local source map

Update the emitter so generated code contains:
- one embedded `source.Source` for the main document
- one embedded `runtime.MapResolver`-compatible map for imported local definitions

- [ ] **Step 5: Re-run the `@use` closure tests**

Run:

```bash
go test ./bind -run TestGenerateEmbedsLocalUseSources -v
```

Expected:
- PASS

- [ ] **Step 6: Commit the embedded-source chunk**

Run:

```bash
git add bind/generator.go bind/assets.go bind/use_closure.go bind/use_closure_test.go
git commit -m "feat: embed local use sources in generated bindings"
```

## Chunk 5: End-to-End Generated Package Verification

### Task 5: Prove the generated package compiles and runs as a real consumer artifact

**Files:**
- Create: `bind/e2e_test.go`
- Create: `bind/testdata/service/service.vxt`
- Create: `bind/testdata/service/schema.vxt`
- Modify: `README.md`
- Test: `bind/e2e_test.go`, `go test ./...`

- [ ] **Step 1: Write the end-to-end failing test first**

Create `bind/e2e_test.go` that:
- generates bindings from `bind/testdata/service/service.vxt`
- writes the generated file into a temp module under `.vxt/`
- writes a tiny consumer `main.go` that imports the generated package
- runs `go test` or `go run` inside that temp module

The consumer should instantiate typed input like:

```go
input := servicevxt.Input{
	Entity: servicevxt.Entity{
		Name:          "User",
		PackageName:   "user",
		HasRepository: true,
	},
}
```

and call:

```go
plan, err := servicevxt.Plan(input)
```

- [ ] **Step 2: Run the end-to-end test and confirm it fails**

Run:

```bash
go test ./bind -run TestGeneratedBindingsCompileAndPlan -v
```

Expected:
- FAIL due to missing or invalid generated package details

- [ ] **Step 3: Fix generated output until the temp consumer builds cleanly**

Adjust generation logic so the produced file:
- imports public `runtime`, `source`, and `write` packages correctly
- compiles without referencing repo-internal packages
- converts typed input to runtime input shape correctly
- uses embedded source/resolver data for `Compile` and `Validate`

- [ ] **Step 4: Add a README section for generated bindings**

Update `README.md` with a concise new section:
- what the `bind` package does
- document-mode-only scope
- generated `.vxt/` directory pattern
- one tiny usage example with `Input`, `Plan`, and `Write`

Do not describe CLI commands here. Keep it library-first.

- [ ] **Step 5: Run the full verification suite**

Run:

```bash
go test ./bind -v
go test ./...
```

Expected:
- PASS for bind package tests
- PASS for the full repo

- [ ] **Step 6: Commit the end-to-end chunk**

Run:

```bash
git add bind/e2e_test.go bind/testdata/service/service.vxt bind/testdata/service/schema.vxt README.md
git commit -m "feat: add typed go binding generator"
```

## Chunk 6: Scope Discipline and Follow-Up Notes

### Task 6: Record the first-iteration limits explicitly

**Files:**
- Modify: `docs/superpowers/specs/2026-05-24-vxt-go-bindings-design.md` (only if implementation reveals spec drift)
- Modify: `README.md` (only if the binding limitations are still unclear)
- Test: documentation sanity checked against implemented tests

- [ ] **Step 1: Compare the implementation against the approved spec**

Check:
- document mode only
- one document equals one binding unit
- local-path `@use` only
- no single-file binding support
- no CLI behavior added

If implementation stayed within scope, do not expand the spec.

- [ ] **Step 2: Add or refine limitation notes only if needed**

If readers could misunderstand the feature, add one short limitations note covering:
- local-path `@use` only
- no package-registry imports
- generated files are intended to be committed
- root `Input` is the canonical generated contract

- [ ] **Step 3: Re-run the final repo verification**

Run:

```bash
go test ./...
```

Expected:
- PASS

- [ ] **Step 4: Commit only if documentation changed in this chunk**

Run:

```bash
git add README.md docs/superpowers/specs/2026-05-24-vxt-go-bindings-design.md
git commit -m "docs: clarify go bindings first-iteration scope"
```
