# VXT Apply Plan Hook Execution Design

Date: 2026-05-30
Status: Draft

## Summary

`vxt` already parses document hooks and surfaces them as planned metadata, but it
does not execute them.

This slice adds one explicit post-write workflow API so callers can:

- write planned artifacts
- then execute `after:write` hooks
- through an injected executor

The core design goal is to keep `WritePlan(...)` pure while still making
workflow automation possible for consumers such as `vx`.

## Goals

- add explicit automatic execution for `after:write` hooks
- keep hook execution outside the default write path
- use an injected `HookExecutor`
- expose write success and hook failure separately
- preserve current planned-hook ordering

## Non-Goals

- shell executor built into `vxt`
- trust or sandbox policy
- hook rollback or transactional semantics
- multi-event hook lifecycle framework
- parallel hook execution
- retry logic
- streaming command output capture

## Public API

Recommended initial API:

```go
type HookContext struct {
	Event       string
	Plan        Plan
	WriteReport write.WriteReport
}

type HookExecutor interface {
	Execute(ctx HookContext, hook PlannedHook) error
}

type ApplyResult struct {
	WriteResult WriteResult
	HookErrors  []error
}

func ApplyPlan(p Plan, target write.OutputTarget, executor HookExecutor) ApplyResult
```

Notes:

- `WritePlan(...)` remains unchanged
- `ApplyPlan(...)` is the explicit workflow API
- `PlannedHook` remains the public hook model

## Behavior

`ApplyPlan(...)` should behave like this:

1. write the plan through the existing write path
2. if write fails, stop and return the write failure
3. if write succeeds, execute only hooks with event `after:write`
4. execute hooks sequentially in the order they appear in `Plan.PlannedHooks`
5. collect hook failures in `HookErrors`

## Error Semantics

Write and hook failures must be distinguished clearly.

Rules:

- if write fails:
  - `WriteResult.Err` is non-nil
  - no hooks run
- if write succeeds but hooks fail:
  - `WriteResult.Err` remains nil
  - `HookErrors` contains one entry per failed hook
- no rollback is attempted

This is deliberate: a failed `after:write` hook does not mean the write itself
failed.

## Hook Context

Initial `HookContext` remains intentionally small:

- `Event`
- `Plan`
- `WriteReport`

This is enough for future consumers such as `vx` to decide:

- what command to run
- where outputs were written
- what artifacts were created or updated

It deliberately avoids adding:

- environment maps
- working directory rules
- arbitrary metadata bags

## Why This Boundary

This keeps `vxt` clean:

- `WritePlan(...)` stays a pure write primitive
- hook automation is opt-in through `ApplyPlan(...)`
- `vxt` does not become a shell runtime by default
- `vx` can implement the real executor later without duplication

## Success Criteria

- callers can automatically execute `after:write` hooks
- `WritePlan(...)` behavior remains unchanged
- hook execution is explicit and opt-in
- write errors and hook errors are clearly separated
- planned hook order is preserved

## Open Questions

- whether `ApplyResult` should later gain helper methods like `Err()` or
  `HasErrors()`
- whether future events such as `before:write` should share the same API
- whether `vx` should supply a shell executor plus trust policy on top of this
