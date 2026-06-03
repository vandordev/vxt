# Document Mode

Document mode is the primary `vxt` authoring format. A document describes a
template contract and the file tree that should be planned from caller input.

Use document mode when one template can produce multiple files, directories, or
typed variants.

## Tutorial Example

This example generates a small Go service package.

```vxt
@template service
@use "./schema.vxt"

@type Options {
  repository: bool
}

@input entity Entity
@input options Options

@dir "internal/{{ entity.name | snake }}"

@partial imports
import "context"
@endpartial

@file "internal/{{ entity.name | snake }}/service.go" mode=overwrite
package {{ entity.name | snake }}

{{ include imports }}

type {{ entity.name | pascal }}Service struct{}

func (s {{ entity.name | pascal }}Service) Get(ctx context.Context) error {
	return nil
}
@endfile

@if options.repository
@file "internal/{{ entity.name | snake }}/repository.go"
package {{ entity.name | snake }}

type {{ entity.name | pascal }}Repository struct{}
@endfile
@endif

@hook after:write "gofmt -w internal/{{ entity.name | snake }}"
```

The referenced definition document can provide shared types:

```vxt
@type Entity {
  name: string
}
```

Load the main document with `runtime.CompileDocumentWithResolver` when it uses
local `@use` definitions:

```go
compiled := runtime.CompileDocumentWithResolver(mainSource, runtime.MapResolver{
	"./schema.vxt": schemaSource,
})
```

## Declaration Walkthrough

`@template service` names the document. The runtime stores the name on the
compiled document, and generated Go bindings use it to derive generated file and
wrapper names.

`@use "./schema.vxt"` references a local definition document. The runtime does
not read files by itself; callers provide a `runtime.SourceResolver`. The
in-memory `runtime.MapResolver` is useful for tests and embedding.

`@type Options { ... }` declares an object type local to the document. Fields
use `name: type` syntax. Current validation supports scalar names such as
`string` and `bool`, local object types, optional fields, and arrays as
implemented by the schema validator.

`@input entity Entity` and `@input options Options` declare required input
values. `runtime.ValidateDocument` checks the input map against these
declarations before planning.

`@dir "internal/{{ entity.name | snake }}"` declares a directory output.
Template expressions in paths are rendered during planning. Expressions can use
case filters such as `snake`, `upper_snake`, `kebab`, `pascal`, `camel`,
`lower`, and `upper`.

`@partial imports ... @endpartial` declares reusable text. Inside a file body,
`{{ include imports }}` inserts the partial content.

`@file "..." mode=overwrite ... @endfile` declares one file output. The path and
body are rendered during planning. Supported write modes are `create`,
`overwrite`, and `skip-if-exists`; omitted mode defaults to `create`.

`@if options.repository ... @endif` conditionally contributes nested `@dir` and
`@file` declarations. Current document-mode conditionals do not support nested
`@if` blocks or `@else`.

`@hook after:write "..."` records planned hook metadata. Planning does not run
the command. `WritePlan` does not run hooks. Only `ApplyPlan` can execute
supported planned hooks, and only through a caller-provided executor.

## Planned Result

With input:

```go
map[string]any{
	"entity": map[string]any{
		"name": "user profile",
	},
	"options": map[string]any{
		"repository": true,
	},
}
```

Planning produces:

- directory `internal/user_profile`
- file `internal/user_profile/service.go`
- file `internal/user_profile/repository.go`
- planned hook metadata for `after:write`

The hook remains metadata unless the caller explicitly uses `ApplyPlan` with a
`runtime.HookExecutor`.

## Directive Reference

`@template <name>` names a document. It applies at document top level. Main
documents must declare it; definition documents loaded by `@use` do not need
one.

`@use "<path>"` references a definition document. It applies at document top
level. The runtime only resolves paths through a caller-provided resolver; it
does not implement package registry semantics.

`@type <Name> { ... }` declares an input object type. It applies at document top
level and in definition documents. It is for validation and generated Go
bindings, not for arbitrary runtime reflection APIs.

`@input <name> <type>` declares one required input. It applies at document top
level. Missing inputs are validation diagnostics.

`@dir "<path>"` declares one directory output. It applies at document top level
or inside an `@if` block. The path can contain template expressions.

`@partial <name> ... @endpartial` declares reusable file-body text. It applies
at document top level. Partials are included with `{{ include name }}` inside
file bodies.

`@file "<path>" [mode=<mode>] ... @endfile` declares one file output. It applies
at document top level or inside an `@if` block. Supported modes are `create`,
`overwrite`, and `skip-if-exists`; omitted mode defaults to `create`.

Template expressions use `{{ path.to.value }}` syntax. Add `| filter` to convert
text when rendering paths or file bodies: `{{ entity.name | snake }}`,
`{{ entity.name | pascal }}`, or `{{ entity.name | camel }}`. Supported filters
are `snake`, `upper_snake`, `kebab`, `pascal`, `camel`, `lower`, and `upper`.

`@if <expr> ... @endif` conditionally contributes nested file or directory
outputs. It applies at document top level. Current parsing rejects nested `@if`
and `@else`.

`@hook <event> "<run>"` records hook metadata. It applies at document top level.
`ApplyPlan` currently executes only `after:write` hooks through an injected
executor; other events remain planned metadata.
