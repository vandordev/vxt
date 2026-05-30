# Runtime API

The `runtime` package exposes the document-mode lifecycle as separate library
stages. Keep the stages separate when callers need validation, preview, tests,
or policy checks before writing output.

## Lifecycle

1. Compile source into a compiled document.
2. Validate caller input against declared inputs and types.
3. Plan concrete directories, files, and hook metadata.
4. Write the plan to an output target.
5. Optionally apply the plan with explicit hook execution.

Each stage has its own result type and diagnostics or errors. Stop the pipeline
when a stage reports diagnostics or an error.

## Compile

`runtime.CompileDocument(src source.Source) runtime.CompileResult` parses one
document-mode source.

```go
compiled := runtime.CompileDocument(src)
if len(compiled.Diagnostics) > 0 {
	return compiled.Diagnostics[0]
}
doc := compiled.Document
```

`CompileResult.Document` is set on success. Diagnostics describe parse or
document-shape failures such as a missing `@template` or an unterminated block.

`runtime.CompileDocumentWithResolver(src, resolver)` parses the same main
document and resolves local `@use` definition documents through the supplied
`runtime.SourceResolver`.

```go
compiled := runtime.CompileDocumentWithResolver(src, runtime.MapResolver{
	"./schema.vxt": schemaSource,
})
```

The runtime does not load files, resolve packages, or contact registries by
itself. `MapResolver` is a small in-memory resolver; applications can provide
their own resolver implementation.

## Validate

`runtime.ValidateDocument(doc, input)` checks declared `@input` values against
known scalar and object types.

```go
validated := runtime.ValidateDocument(compiled.Document, map[string]any{
	"name": "Vandor",
})
if len(validated.Diagnostics) > 0 {
	return validated.Diagnostics[0]
}
```

`ValidationResult` keeps the compiled document and the input map for planning.
Diagnostics include missing inputs and type mismatches.

## Plan

`runtime.PlanDocument(validated)` renders planned filesystem artifacts.

```go
planned := runtime.PlanDocument(validated)
if len(planned.Diagnostics) > 0 {
	return planned.Diagnostics[0]
}
```

`PlanResult.Plan` contains:

- `Dirs`: rendered directory paths
- `Files`: rendered file paths, contents, and write modes
- `PlannedHooks`: hook metadata declared by `@hook`

Planning does not write files. It also does not execute hooks.

## Write

`runtime.WritePlan(plan, target)` writes planned directories and files to a
`write.OutputTarget`.

```go
target := write.NewMemoryTarget()
report, err := runtime.WritePlan(planned.Plan, target)
if err != nil {
	return err
}
_ = report
```

`WritePlan` returns the compatibility surface: a `write.WriteReport` and an
error. It does not execute planned hooks.

Use `runtime.WritePlanWithDiagnostics(plan, target)` when the caller needs
structured diagnostics for write failures:

```go
result := runtime.WritePlanWithDiagnostics(planned.Plan, target)
if result.Err != nil {
	return result.Diagnostics[0]
}
```

Write diagnostics distinguish common output failures such as existing files,
path escapes, and unsupported write modes.

## Apply

`runtime.ApplyPlan(plan, target, executor)` writes the plan and then executes
supported planned hooks through the provided executor.

```go
result := runtime.ApplyPlan(planned.Plan, target, executor)
if result.WriteResult.Err != nil {
	return result.WriteResult.Err
}
if len(result.HookErrors) > 0 {
	return result.HookErrors[0]
}
```

`ApplyPlan` is explicit post-write hook execution. It is not called by
`WritePlan`. If the executor is nil, no hook is executed. The current supported
hook event is `after:write`; unsupported events are skipped.

`runtime.HookExecutor` is the trust and execution boundary:

```go
type HookExecutor interface {
	Execute(ctx runtime.HookContext, hook runtime.PlannedHook) error
}
```

The executor decides what a hook string means, whether it is allowed, and how it
is run.

## Output Targets

`write.OutputTarget` is the interface consumed by the runtime write stage:

```go
type OutputTarget interface {
	MkdirAll(path string) error
	WriteFile(path string, content []byte, mode string) (bool, error)
}
```

`write.NewMemoryTarget()` records output in memory. Use it for tests, previews,
and first-use learning.

`write.NewFilesystemTarget(root)` writes under one filesystem root. It creates
needed directories, rejects absolute paths and `..` escapes, and applies file
write modes.

Current write modes are:

- `create`: fail when the target file already exists
- `overwrite`: write regardless of existing file state
- `skip-if-exists`: leave an existing file unchanged and count it as skipped

Document-mode `@file` declarations default to `create`.

## Diagnostics Guidance

Prefer diagnostics-bearing stages when building user-facing tools:

- inspect `CompileResult.Diagnostics` after compile
- inspect `ValidationResult.Diagnostics` after validate
- inspect `PlanResult.Diagnostics` after plan
- use `WritePlanWithDiagnostics` when write errors should be surfaced as
  structured `diag.Diagnostic` values

Use the shorter `WritePlan` form when a plain Go error is enough.
