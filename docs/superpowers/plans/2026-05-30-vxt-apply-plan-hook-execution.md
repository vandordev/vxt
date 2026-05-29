# VXT Apply Plan Hook Execution Implementation Plan

> **For agentic workers:** REQUIRED: Use `superpowers:subagent-driven-development` (if subagents available) or `superpowers:executing-plans` to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add an explicit `ApplyPlan(...)` workflow API that writes a plan and then executes `after:write` hooks through an injected executor, without changing the semantics of `WritePlan(...)`.

**Architecture:** Keep the existing write path intact and layer hook execution on top through a small runtime workflow unit. Model execution through `HookContext`, `HookExecutor`, and `ApplyResult`, and keep hook running strictly sequential and opt-in.

**Tech Stack:** Go, existing `runtime` package, existing `write.OutputTarget`, existing `PlannedHook`/`Plan` model, `go test`.

---

## Chunk 1: Public Contract for Hook Execution

### Task 1: Add the explicit runtime hook execution contract

**Files:**
- Modify: `runtime/types.go`
- Create: `runtime/apply.go`
- Create: `runtime/apply_test.go`
- Test: `runtime/apply_test.go`

- [ ] **Step 1: Write the failing contract tests**

Create `runtime/apply_test.go` with initial compile-level tests that require:
- `HookContext`
- `HookExecutor`
- `ApplyResult`
- `ApplyPlan(...)`

Use a tiny fake executor:

```go
type recordingExecutor struct {
	calls []runtime.PlannedHook
}

func (r *recordingExecutor) Execute(ctx runtime.HookContext, hook runtime.PlannedHook) error {
	r.calls = append(r.calls, hook)
	return nil
}
```

Add one first test like:

```go
func TestApplyPlanRunsAfterWriteHooks(t *testing.T) {
	plan := runtime.Plan{
		PlannedHooks: []runtime.PlannedHook{
			{Event: "after:write", Run: "echo one"},
		},
	}
	target := write.NewMemoryTarget()
	exec := &recordingExecutor{}

	result := runtime.ApplyPlan(plan, target, exec)

	if result.WriteResult.Err != nil {
		t.Fatalf("unexpected write err: %v", result.WriteResult.Err)
	}
	if len(exec.calls) != 1 {
		t.Fatalf("got %d hook calls", len(exec.calls))
	}
}
```

- [ ] **Step 2: Run the focused apply tests and verify they fail**

Run:

```bash
go test ./runtime -run TestApplyPlan -v
```

Expected:
- FAIL because the new hook-application API does not exist yet

- [ ] **Step 3: Add the public contract**

Extend `runtime/types.go` with:

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
```

Create `runtime/apply.go` with a minimal stub:

```go
func ApplyPlan(p Plan, target write.OutputTarget, executor HookExecutor) ApplyResult {
	return ApplyResult{}
}
```

- [ ] **Step 4: Re-run the apply tests and keep them red for the correct reason**

Run:

```bash
go test ./runtime -run TestApplyPlan -v
```

Expected:
- still FAIL because behavior is not implemented yet

- [ ] **Step 5: Do not commit yet**

Keep moving within this chunk until the behavior is green.

## Chunk 2: Write-Then-Execute Behavior

### Task 2: Implement `after:write` execution after successful writes

**Files:**
- Modify: `runtime/apply.go`
- Modify: `runtime/apply_test.go`
- Test: `runtime/apply_test.go`

- [ ] **Step 1: Extend the failing tests for core behavior**

Add tests for:
- write happens before hooks
- only `after:write` hooks run
- non-`after:write` hooks are ignored in this slice
- hooks run in the order they appear in `Plan.PlannedHooks`

Add a write-bearing plan:

```go
plan := runtime.Plan{
	Files: []runtime.FileOutput{{
		Path:    "hello.txt",
		Content: "hello",
		Mode:    "create",
	}},
	PlannedHooks: []runtime.PlannedHook{
		{Event: "after:write", Run: "echo one"},
		{Event: "after:write", Run: "echo two"},
	},
}
```

Then assert:
- the memory target got one written file
- hooks were called in the same order

- [ ] **Step 2: Run the focused apply tests and verify failure**

Run:

```bash
go test ./runtime -run TestApplyPlan -v
```

Expected:
- FAIL because `ApplyPlan(...)` is still stubbed

- [ ] **Step 3: Implement the minimal workflow in `runtime/apply.go`**

Implementation shape:

```go
func ApplyPlan(p Plan, target write.OutputTarget, executor HookExecutor) ApplyResult {
	result := ApplyResult{
		WriteResult: WritePlanWithDiagnostics(p, target),
	}
	if result.WriteResult.Err != nil {
		return result
	}
	if executor == nil {
		return result
	}

	for _, hook := range p.PlannedHooks {
		if hook.Event != "after:write" {
			continue
		}
		err := executor.Execute(HookContext{
			Event:       hook.Event,
			Plan:        p,
			WriteReport: result.WriteResult.Report,
		}, hook)
		if err != nil {
			result.HookErrors = append(result.HookErrors, err)
		}
	}

	return result
}
```

Keep it strictly sequential. No retries, no rollback.

- [ ] **Step 4: Re-run the apply tests**

Run:

```bash
go test ./runtime -run TestApplyPlan -v
```

Expected:
- PASS for core hook execution behavior

- [ ] **Step 5: Commit the workflow chunk**

Run:

```bash
git add runtime/types.go runtime/apply.go runtime/apply_test.go
git commit -m "feat: add explicit apply plan hook execution"
```

## Chunk 3: Failure Semantics

### Task 3: Distinguish write failure from hook failure

**Files:**
- Modify: `runtime/apply_test.go`
- Test: `runtime/apply_test.go`

- [ ] **Step 1: Add failing tests for error semantics**

Add tests for:
- if write fails, no hooks run
- if hook fails, write remains successful but `HookErrors` is populated
- multiple hook failures are accumulated

Use:
- a custom failing `write.OutputTarget` for write failure
- a custom executor that returns errors for selected hooks

Example executor:

```go
type failingExecutor struct {
	failRuns map[string]error
	calls    []runtime.PlannedHook
}
```

Assertions:
- `WriteResult.Err` is non-nil only on write failure
- `HookErrors` contains errors only for failing hooks
- hook calls stop only if you explicitly design them to stop; in this slice they should continue and accumulate

- [ ] **Step 2: Run the tests and verify failure**

Run:

```bash
go test ./runtime -run 'TestApplyPlan.*(WriteFail|HookFail|Accumulates)' -v
```

Expected:
- FAIL until failure semantics are fully covered

- [ ] **Step 3: Implement or adjust behavior as needed**

Expected final behavior:
- write failure returns immediately and skips hooks
- hook failure appends to `HookErrors`
- later hooks still run

If the current implementation already satisfies this, no production change may be needed; keep the tests.

- [ ] **Step 4: Re-run the focused failure tests**

Run:

```bash
go test ./runtime -run 'TestApplyPlan.*(WriteFail|HookFail|Accumulates)' -v
```

Expected:
- PASS

- [ ] **Step 5: Commit the failure-semantics coverage**

Run:

```bash
git add runtime/apply_test.go
git commit -m "test: cover apply plan hook failure semantics"
```

## Chunk 4: README and Full Verification

### Task 4: Document the explicit workflow API and verify the whole repo

**Files:**
- Modify: `README.md`
- Modify: `runtime/doc.go` (only if package-level runtime docs need one sentence for `ApplyPlan`)
- Test: `go test ./runtime -v`, `go test ./...`

- [ ] **Step 1: Add a short README note for planned-vs-applied hooks**

Update `README.md` in the document pipeline or hooks area so it clearly says:
- `WritePlan(...)` writes only
- hooks are surfaced as planned metadata
- `ApplyPlan(...)` is the explicit opt-in workflow API for automatic `after:write` execution through an injected executor

Do not document a shell executor here because `vxt` does not provide one.

- [ ] **Step 2: Optionally add one line to `runtime/doc.go` if runtime workflow API is not discoverable enough**

Only if needed, add a short sentence that `runtime` exposes both pure write and explicit apply workflows.

- [ ] **Step 3: Run focused runtime verification**

Run:

```bash
go test ./runtime -v
```

Expected:
- PASS

- [ ] **Step 4: Run full repo verification**

Run:

```bash
go test ./...
```

Expected:
- PASS

- [ ] **Step 5: Commit the documentation update**

Run:

```bash
git add README.md runtime/doc.go
git commit -m "docs: describe apply plan hook execution"
```

## Chunk 5: Scope Check

### Task 5: Confirm the slice stayed inside the approved boundary

**Files:**
- Modify: `docs/superpowers/specs/2026-05-30-vxt-apply-plan-hook-execution-design.md` (only if clarification is needed)
- Test: `go test ./...`

- [ ] **Step 1: Compare implementation against the approved design**

Check:
- explicit `ApplyPlan(...)`
- `after:write` only
- injected executor only
- no shell helper in `vxt`
- sequential execution
- no rollback

- [ ] **Step 2: Add clarification only if needed**

Only if readers could misunderstand, add one short note:
- `ApplyPlan(...)` is opt-in
- `WritePlan(...)` remains pure
- `after:write` is the only executed event in this slice

- [ ] **Step 3: Run final verification**

Run:

```bash
go test ./...
```

Expected:
- PASS

- [ ] **Step 4: Commit only if docs changed**

Run:

```bash
git add docs/superpowers/specs/2026-05-30-vxt-apply-plan-hook-execution-design.md
git commit -m "docs: clarify apply plan hook execution scope"
```
