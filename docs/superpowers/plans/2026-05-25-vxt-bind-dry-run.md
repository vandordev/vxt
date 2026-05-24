# VXT Bind Dry-Run Implementation Plan

> **For agentic workers:** REQUIRED: Use `superpowers:subagent-driven-development` (if subagents available) or `superpowers:executing-plans` to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add realistic dry-run support to `bind.WriteOutput(...)` and `bind.GenerateToDir(...)` so callers can preview create/overwrite/remove actions without applying filesystem changes.

**Architecture:** Extend the `bind` write layer with `WriteOptions`, richer report types, and a single planning path that computes actions against the real target directory. Real write and dry-run should share the same action planning logic, with dry-run only skipping side effects.

**Tech Stack:** Go, existing `bind` package, `os` and `path/filepath`, temp-directory filesystem tests, `go test`.

---

## Chunk 1: Options and Report Contract Upgrade

### Task 1: Add `WriteOptions`, action kinds, and richer reports

**Files:**
- Modify: `bind/types.go`
- Modify: `bind/write_test.go`
- Test: `bind/write_test.go`

- [ ] **Step 1: Write the failing report-shape tests**

Update `bind/write_test.go` so existing write tests also assert:
- `report.DryRun == false` on normal writes
- `report.Actions` contains entries for created or overwritten files

Add a small focused test for action kinds:
- a normal create should produce one `create` action
- a normal overwrite should produce one `overwrite` action

- [ ] **Step 2: Run the focused writer tests and confirm they fail**

Run:

```bash
go test ./bind -run TestWriteOutput -v
```

Expected:
- FAIL because `DryRun`, `Actions`, or action kinds are missing

- [ ] **Step 3: Extend the public write contract**

Modify `bind/types.go` to add:

```go
type WriteOptions struct {
	DryRun bool
}

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

Update `WriteReport`:

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

- [ ] **Step 4: Re-run the writer tests to keep them red for the right reason**

Run:

```bash
go test ./bind -run TestWriteOutput -v
```

Expected:
- still FAIL because implementation has not been updated yet

- [ ] **Step 5: Commit only if the contract update is isolated and tests are intentionally red**

Do not commit this chunk yet unless the implementation is also green within the same chunk.

## Chunk 2: Canonical Write Planning Path

### Task 2: Route real writes through a unified action planner

**Files:**
- Modify: `bind/write.go`
- Create: `bind/write_plan.go`
- Modify: `bind/write_test.go`
- Test: `bind/write_test.go`

- [ ] **Step 1: Write failing tests for action planning behavior**

Add focused tests that assert normal writes now populate:
- `FilesWritten`
- `FilesOverwritten`
- `FilesRemoved`
- `Actions`

Use the same scenarios already present:
- create
- overwrite
- scoped stale removal

- [ ] **Step 2: Run the writer tests and confirm they fail on missing implementation**

Run:

```bash
go test ./bind -run TestWriteOutput -v
```

Expected:
- FAIL because `WriteOutput` still lacks action tracking or new signature support

- [ ] **Step 3: Add an internal write plan model**

Create `bind/write_plan.go` with focused internal types such as:

```go
type plannedWriteAction struct {
	Path   string
	Action WriteActionKind
}
```

Add an internal planner function like:

```go
func planWriteActions(out Output, dir string) ([]plannedWriteAction, error)
```

This planner should:
- inspect the real target dir
- decide whether each output file is a create or overwrite
- detect stale owned files that would be removed
- return the full ordered change set

- [ ] **Step 4: Update `WriteOutput` to accept options and use the planner**

Change signature:

```go
func WriteOutput(out Output, dir string, opts WriteOptions) (WriteReport, error)
```

Then:
- call the planner first
- initialize `WriteReport` with `DryRun` and `OutputDir`
- execute actions only when `opts.DryRun == false`
- always populate summary lists and `Actions`

For this chunk, focus on the real-write path first. Dry-run execution skipping comes in the next chunk.

- [ ] **Step 5: Re-run the write tests**

Run:

```bash
go test ./bind -run TestWriteOutput -v
```

Expected:
- PASS for normal write behavior with richer report data

- [ ] **Step 6: Commit the planner and contract upgrade**

Run:

```bash
git add bind/types.go bind/write.go bind/write_plan.go bind/write_test.go
git commit -m "feat: add bind write action planning"
```

## Chunk 3: Realistic Dry-Run

### Task 3: Add dry-run behavior with no filesystem side effects

**Files:**
- Modify: `bind/write.go`
- Modify: `bind/write_test.go`
- Test: `bind/write_test.go`

- [ ] **Step 1: Write failing dry-run tests**

Extend `bind/write_test.go` with scenarios for:
- dry-run create
- dry-run overwrite
- dry-run stale removal

Assertions:
- report shows the same create/overwrite/remove actions as a real run would
- `report.DryRun == true`
- target directory contents remain unchanged after the call

- [ ] **Step 2: Run the dry-run tests and confirm they fail**

Run:

```bash
go test ./bind -run TestWriteOutputDryRun -v
```

Expected:
- FAIL because dry-run either does not exist yet or still mutates the filesystem

- [ ] **Step 3: Implement dry-run execution skipping**

Update `WriteOutput` so:
- the planner always runs
- summary lists and `Actions` are always populated
- if `opts.DryRun` is true:
  - no directory creation beyond what is necessary to inspect state
  - no file writes
  - no file removals
  - the report returns immediately after planning

Use the same planned action set as the real write path.

- [ ] **Step 4: Re-run the dry-run tests**

Run:

```bash
go test ./bind -run TestWriteOutputDryRun -v
```

Expected:
- PASS

- [ ] **Step 5: Re-run the broader writer suite**

Run:

```bash
go test ./bind -run TestWriteOutput -v
```

Expected:
- PASS for both real-write and dry-run cases

- [ ] **Step 6: Commit the dry-run chunk**

Run:

```bash
git add bind/write.go bind/write_test.go
git commit -m "feat: add bind dry-run support"
```

## Chunk 4: Convenience Path and E2E Dry-Run

### Task 4: Support dry-run through `GenerateToDir(...)`

**Files:**
- Modify: `bind/generator.go`
- Modify: `bind/e2e_test.go`
- Modify: `README.md`
- Test: `bind/e2e_test.go`, `go test ./...`

- [ ] **Step 1: Write the failing convenience dry-run test**

Extend `bind/e2e_test.go` with a test that:
- calls `GenerateToDir(req, dir, WriteOptions{DryRun: true})`
- targets a temp `.vxt/` directory
- asserts the report contains realistic actions
- asserts the filesystem is unchanged after the call

- [ ] **Step 2: Run the focused e2e dry-run test and confirm it fails**

Run:

```bash
go test ./bind -run TestGenerateToDirDryRun -v
```

Expected:
- FAIL because `GenerateToDir` does not yet accept options

- [ ] **Step 3: Update `GenerateToDir` to accept `WriteOptions`**

Change signature:

```go
func GenerateToDir(req Request, dir string, opts WriteOptions) (WriteReport, error)
```

Implementation:
- call `Generate(req)`
- then call `WriteOutput(out, dir, opts)`

- [ ] **Step 4: Update README examples**

Revise the bindings section in `README.md` so it shows:
- `bind.GenerateToDir(req, ".vxt", bind.WriteOptions{})`
- `bind.GenerateToDir(req, ".vxt", bind.WriteOptions{DryRun: true})`

Keep examples short and library-first.

- [ ] **Step 5: Run full verification**

Run:

```bash
go test ./bind -v
go test ./...
```

Expected:
- PASS for bind package
- PASS for full repo

- [ ] **Step 6: Commit the convenience dry-run flow**

Run:

```bash
git add bind/generator.go bind/e2e_test.go README.md
git commit -m "feat: add bind generate-to-dir dry-run"
```

## Chunk 5: Scope Check and Cleanup

### Task 5: Confirm the dry-run slice stayed within the approved design

**Files:**
- Modify: `docs/superpowers/specs/2026-05-25-vxt-bind-dry-run-design.md` (only if implementation drift needs clarification)
- Modify: `README.md` (only if limitations remain unclear)
- Test: `go test ./...`

- [ ] **Step 1: Compare implementation against the approved slice**

Check:
- options-based API
- realistic preview
- summary lists plus `Actions`
- no side effects in dry-run
- no extra write-policy matrix

- [ ] **Step 2: Add limitation notes only if needed**

Only if necessary, add one short note clarifying:
- dry-run is realistic, not shallow
- action report shows change set only
- preserve/no-op files are intentionally not listed

- [ ] **Step 3: Run final verification**

Run:

```bash
go test ./...
```

Expected:
- PASS

- [ ] **Step 4: Commit only if docs changed in this chunk**

Run:

```bash
git add README.md docs/superpowers/specs/2026-05-25-vxt-bind-dry-run-design.md
git commit -m "docs: clarify bind dry-run scope"
```
