# VXT Bind Dry-Run Design

Date: 2026-05-25
Status: Draft

## Summary

`bind.WriteOutput(...)` and `bind.GenerateToDir(...)` currently apply real
filesystem changes directly.

This slice adds a realistic dry-run mode so callers can preview:

- which files would be created
- which files would be overwritten
- which stale generated files would be removed

without applying any actual filesystem changes.

The dry-run should be realistic, meaning it still inspects the real target
directory and computes actions against current state.

## Goals

- add realistic dry-run support to `bind` write flows
- keep `WriteOutput(...)` as the canonical write surface
- allow convenience dry-run through `GenerateToDir(...)`
- expose both summary lists and per-file action details
- avoid any filesystem side effects when dry-run is enabled

## Non-Goals

- full filesystem diff engine
- recursive inventory of untouched files
- reporting preserve/no-op entries for every file
- transactional rollback semantics
- broader write policy matrix like no-overwrite or no-cleanup

## Public API

Recommended new write options:

```go
type WriteOptions struct {
	DryRun bool
}
```

Recommended updated signatures:

```go
func WriteOutput(out Output, dir string, opts WriteOptions) (WriteReport, error)

func GenerateToDir(req Request, dir string, opts WriteOptions) (WriteReport, error)
```

This is an intentional direct signature upgrade. The API is still young enough
that explicit cleanup of the surface is preferable to carrying compatibility
wrappers too early.

## Report Shape

Recommended action model:

```go
type WriteActionKind string

const (
	WriteActionCreate    WriteActionKind = "create"
	WriteActionOverwrite WriteActionKind = "overwrite"
	WriteActionRemove    WriteActionKind = "remove"
)

type WriteAction struct {
	Path   string
	Action WriteActionKind
}
```

Recommended updated report:

```go
type WriteReport struct {
	DryRun           bool
	OutputDir        string
	FilesWritten     []string
	FilesOverwritten []string
	FilesRemoved     []string
	Actions          []WriteAction
}
```

Rules:

- `FilesWritten`, `FilesOverwritten`, and `FilesRemoved` remain the summary layer
- `Actions` provides ordered per-file details
- `DryRun` makes it explicit whether the report describes a preview or real side
  effects

## Dry-Run Behavior

Dry-run must be realistic:

- still create the action plan against the real filesystem state
- still detect which files would be written, overwritten, or removed
- still honor the existing scoped reconcile ownership rules
- never write any files
- never delete any files

This means dry-run is not a shallow simulation based only on the new output. It
must inspect the target directory exactly like a real run.

## Real Write Behavior

Real write behavior remains unchanged in intent:

- create missing output directory
- write new files
- overwrite existing owned target files
- remove stale owned files for the same binding

The only change is that the report becomes richer and the API becomes options
based.

## Action Semantics

Initial action kinds:

- `create`
- `overwrite`
- `remove`

Deliberate omission:

- no `preserve`
- no `skip`
- no `noop`

This keeps reports focused on change sets, which is the information most useful
for both library consumers and future `vx --dry-run` output.

## Why This Boundary

This design keeps responsibilities clean:

- `WriteOutput(...)` remains the canonical place for filesystem behavior
- `GenerateToDir(...)` stays a convenience wrapper
- `vx` can later surface dry-run using the same library behavior
- no CLI concerns leak into `vxt`

## Success Criteria

- callers can preview write behavior without side effects
- dry-run reports are accurate against the real target directory
- summary lists and action details are both available
- scoped reconcile rules are preserved in dry-run mode
- the convenience generation flow supports dry-run too

## Open Questions

- whether action ordering should be formally guaranteed
- whether future write policies should share the same options object
- whether a later human-facing formatter should live in `vx` rather than `vxt`
