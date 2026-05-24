# VXT Go Bindings Design

Date: 2026-05-24
Status: Draft

## Summary

`vxt` remains a library-first product. It should not grow its own end-user CLI.
Human-facing binary workflows belong to `vx`.

To make `vxt` usable from Go in a strict, type-safe way, `vxt` should gain an
optional Go binding generator that turns document-mode `.vxt` templates into
generated Go packages under a `.vxt/` directory.

The generated package should:

- expose public Go structs derived from `@type`
- expose one canonical `Input` struct for `@input`
- embed template source and local `@use` dependencies
- provide typed wrappers around `Compile`, `Validate`, `Plan`, and `Write`
- avoid exposing `map[string]any` in user-facing generated APIs

This keeps `vxt` reusable as an independent library while allowing `vx` to call
the same binding/codegen APIs later from its CLI workflows.

## Goals

- keep `vxt` library-only
- let `.vxt` documents become source-of-truth for template schema and behavior
- generate self-contained Go bindings for document-mode templates
- give Go consumers a strict, ergonomic API with no user-facing `map[string]any`
- keep `vx` free to call the same library generator later

## Non-Goals

- shipping a standalone `vxt` CLI
- single-file binding generation in the first iteration
- package-registry `@use` resolution
- merge/patch support
- parameterized partial binding APIs
- preserving every runtime internal detail as generated public API
- multi-language binding generation

## Product Boundary

The binding generator should live inside `vxt` as an optional library package,
not in the root package and not as a separate binary.

Recommended package direction:

- `github.com/vandordev/vxt/bind`
  or
- `github.com/vandordev/vxt/codegen`

`vx` should later act as a thin CLI wrapper around this library capability.

## Initial Scope

The initial binding generator scope is intentionally narrow:

- document mode only
- one document equals one canonical binding unit
- generated output written under `.vxt/`
- generated files are intended to be committed
- local-path `@use` is resolved and embedded into the generated package
- generated package is self-contained at runtime

## Generated Layout

Example:

```text
my-template/
  service.vxt
  schema.vxt
  .vxt/
    service_gen.go
```

Rules:

- source `.vxt` files remain the source of truth
- `.vxt/` contains generated Go source only
- generated files should not be edited manually
- generated package name follows the main template, for example `servicevxt`

## Generated Public API

For one document template, the generated package should expose:

- public named structs derived from `@type`
- a canonical root `Input` struct derived from `@input`
- `Compile()`
- `Validate(input Input)`
- `Plan(input Input)`
- `Write(input Input, target write.OutputTarget)`
- `PlanDetailed(input Input)`
- `WriteDetailed(input Input, target write.OutputTarget)`

Example target shape:

```go
package servicevxt

type Entity struct {
	Name          string
	PackageName   string
	HasRepository bool
}

type Input struct {
	Entity Entity
}

func Compile() runtime.CompileResult
func Validate(input Input) runtime.ValidationResult
func Plan(input Input) (runtime.Plan, error)
func Write(input Input, target write.OutputTarget) (write.WriteReport, error)
func PlanDetailed(input Input) runtime.PlanResult
func WriteDetailed(input Input, target write.OutputTarget) runtime.WriteResult
```

Notes:

- `Input` is the canonical public contract
- helper parameter-list APIs are intentionally out of scope initially
- diagnostics remain available through detailed variants
- convenience wrappers should remain error-centric for everyday Go usage

## Type Mapping Rules

Initial `.vxt` to Go mapping:

- `string` -> `string`
- `bool` -> `bool`
- `int` -> `int`
- `float` -> `float64`
- `TypeName` -> generated public struct with the same name
- `TypeName[]` -> `[]TypeName`
- optional field -> pointer type

Field naming:

- generated Go fields are exported and use PascalCase
- original schema names are preserved internally for runtime conversion

Example:

```vxt
@type Entity {
  name: string
  package_name: string
  tags: string[]
  repository?: Repository
}
```

becomes:

```go
type Entity struct {
	Name        string
	PackageName string
	Tags        []string
	Repository  *Repository
}
```

## Source Embedding

Generated bindings should embed:

- the main document source
- any local-path `@use` definition sources needed by that document

This allows the generated package to compile and plan without requiring the
original `.vxt` files to exist on disk at runtime.

The generated package should internally construct:

- embedded `source.Source` values
- embedded resolver data for local-path `@use`

## Runtime Relationship

Generated bindings are adapters over `vxt/runtime`. They must not duplicate
engine behavior.

The generator should produce:

- typed public structs
- internal conversion from typed input to runtime input model
- wrapper functions that call existing runtime APIs

The generator should not:

- copy parser logic
- reimplement validation logic
- fork planner or writer behavior

## Why This Boundary

This design keeps product responsibilities clean:

- `vxt` stays reusable and independent
- `vx` stays the CLI and ecosystem surface
- Go consumers get a strict typed API
- template documents remain the authoritative source

This also allows future workflows such as:

- `vx bind`
- `vx generate`
- `vx new`

without forcing `vxt` itself to become a CLI product.

## Success Criteria

- one document `.vxt` template can generate one self-contained Go binding package
- user-facing generated APIs do not require `map[string]any`
- local-path `@use` works without the original source files at runtime
- host Go code can call `Plan` and `Write` in a type-safe way
- `vx` can later call the same generator as a library without code duplication

## Open Questions

- exact package name derivation rules for templates with unusual names
- whether generated bindings should include drift/version metadata
- whether future single-file bindings should share the same generator package
- whether the package name should be `bind` or `codegen`
