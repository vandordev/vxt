# vxt Documentation Refresh Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Reshape the `vxt` repository documentation so a new reader can understand what `vxt` is, complete a first document-mode flow, and navigate the runtime and bindings APIs without reading source first.

**Architecture:** Keep `README.md` as the repo landing page and move detailed guidance into focused docs under `docs/`. Follow the approved split: README for overview and quick start, `getting-started` and `document-mode` for first-use learning, `runtime-api` and `go-bindings` for usage depth, and `concepts` for product-boundary explanation.

**Tech Stack:** Markdown, existing `vxt` public Go API, current repo docs under `docs/`, Go toolchain for regression verification

---

## File Structure

### Existing files to modify

- `README.md`
  - Repo landing page, quick start, docs map, release/public-package positioning

### New files to create

- `docs/getting-started.md`
  - First successful document-mode walkthrough using `MemoryTarget`, then a short move to filesystem output
- `docs/document-mode.md`
  - Tutorial-first `.vxt` authoring guide plus compact directive reference
- `docs/runtime-api.md`
  - Lifecycle-oriented runtime API reference for compile, validate, plan, write, and apply
- `docs/go-bindings.md`
  - Practical guide for generated typed Go bindings and the `bind` package flow
- `docs/concepts.md`
  - Explanation doc for `vxt` boundaries, `vx` relationship, runtime vs bindings, and hooks model

### Existing files to consult while writing

- `runtime/compile.go`
- `runtime/validate.go`
- `runtime/plan.go`
- `runtime/write.go`
- `runtime/apply.go`
- `runtime/types.go`
- `bind/generator.go`
- `bind/write.go`
- `bind/types.go`
- `docs/releases/v0.1.0.md`

The documentation must describe only behavior that is implemented now.

## Chunk 1: README Reshape

### Task 1: Audit current README against approved structure

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Read the current README and capture sections that should stay**

Review:
- current install instructions
- current public package list
- current quick start snippets
- current generated bindings note

- [ ] **Step 2: Draft the new README outline in place**

Required order:
1. what `vxt` is
2. when to use it
3. what it is not
4. install
5. strong document-mode quick start
6. pipeline overview
7. feature snapshot
8. docs map
9. release/public-package note

- [ ] **Step 3: Rewrite README content to match the outline**

Write content that:
- makes `vxt` clearly library-first
- positions `document mode` as the main quick start
- mentions `ApplyPlan(...)` accurately
- links to the new docs files by path

- [ ] **Step 4: Review README for claims that exceed implemented behavior**

Check specifically:
- hook execution wording
- bindings scope wording
- non-goals wording
- any mention of CLI behavior

- [ ] **Step 5: Commit**

```bash
git add README.md
git commit -m "docs: reshape vxt readme"
```

## Chunk 2: First-Use Learning Docs

### Task 2: Create `docs/getting-started.md`

**Files:**
- Create: `docs/getting-started.md`
- Reference: `README.md`, `runtime/compile.go`, `runtime/validate.go`, `runtime/plan.go`, `runtime/write.go`

- [ ] **Step 1: Write the guide outline**

Required sections:
- audience and goal
- first `.vxt` document
- compile
- validate
- plan
- write with `MemoryTarget`
- move to `FilesystemTarget`
- common first mistakes
- next docs

- [ ] **Step 2: Write a minimal document-mode example**

Use one small `.vxt` document with:
- `@template`
- `@input`
- `@file`

Ensure the Go snippet matches current API names exactly.

- [ ] **Step 3: Add the filesystem transition**

Explain:
- when to switch from `MemoryTarget`
- how `FilesystemTarget` changes behavior
- what write side effects to expect

- [ ] **Step 4: Commit**

```bash
git add docs/getting-started.md
git commit -m "docs: add vxt getting started guide"
```

### Task 3: Create `docs/document-mode.md`

**Files:**
- Create: `docs/document-mode.md`
- Reference: `runtime/types.go`, `runtime/compile_test.go`, `runtime/plan_test.go`, `README.md`

- [ ] **Step 1: Choose one realistic tutorial example**

The example should cover:
- `@template`
- `@input`
- `@type`
- `@dir`
- `@file`
- `@partial`
- `@use`
- `@if`
- `@hook`

Keep it small enough to read in one pass.

- [ ] **Step 2: Write the tutorial explanation in declaration order**

Explain each directive where it first appears.
Do not drift into speculative future syntax.

- [ ] **Step 3: Add a compact directive reference at the end**

Each directive entry should answer:
- what it does
- where it applies
- notable current limitation, if any

- [ ] **Step 4: Commit**

```bash
git add docs/document-mode.md
git commit -m "docs: add vxt document mode guide"
```

## Chunk 3: Runtime and Bindings Docs

### Task 4: Create `docs/runtime-api.md`

**Files:**
- Create: `docs/runtime-api.md`
- Reference: `runtime/compile.go`, `runtime/validate.go`, `runtime/plan.go`, `runtime/write.go`, `runtime/apply.go`, `runtime/types.go`

- [ ] **Step 1: Write the lifecycle overview**

Required lifecycle:
- compile
- validate
- plan
- write
- apply

- [ ] **Step 2: Document compile and validate APIs**

Cover:
- `CompileDocument`
- `CompileDocumentWithResolver`
- `ValidateDocument`

State what each stage returns and when diagnostics matter.

- [ ] **Step 3: Document plan, write, and apply APIs**

Cover:
- `PlanDocument`
- `WritePlan`
- `WritePlanWithDiagnostics`
- `ApplyPlan`

Explicitly explain:
- `WritePlan` does not execute hooks
- `ApplyPlan` only executes supported planned hooks through an injected executor

- [ ] **Step 4: Add `OutputTarget` and diagnostics notes**

Explain:
- `MemoryTarget`
- `FilesystemTarget`
- when to prefer diagnostics-bearing variants

- [ ] **Step 5: Commit**

```bash
git add docs/runtime-api.md
git commit -m "docs: add vxt runtime api guide"
```

### Task 5: Create `docs/go-bindings.md`

**Files:**
- Create: `docs/go-bindings.md`
- Reference: `bind/generator.go`, `bind/write.go`, `bind/types.go`, `bind/e2e_test.go`, `README.md`

- [ ] **Step 1: Write the bindings positioning**

Answer:
- what generated bindings are
- when to use them instead of the raw runtime API
- current document-only scope

- [ ] **Step 2: Document generated package usage**

Cover:
- generated `Input`
- generated public types from `@type`
- wrapper methods such as `Compile`, `Validate`, `Plan`, and `Write`

- [ ] **Step 3: Document the `bind` package flow**

Cover:
- `bind.Generate`
- `bind.WriteOutput`
- `bind.GenerateToDir`
- `bind.WriteOptions{DryRun: true}`

Explain local `@use` embedding accurately.

- [ ] **Step 4: Add limitations and current non-goals**

Include:
- document mode only
- no single-file bindings yet
- no package-style `@use`
- no claim of cross-language generated bindings

- [ ] **Step 5: Commit**

```bash
git add docs/go-bindings.md
git commit -m "docs: add vxt go bindings guide"
```

## Chunk 4: Concepts Doc and Final Integration

### Task 6: Create `docs/concepts.md`

**Files:**
- Create: `docs/concepts.md`
- Reference: `README.md`, `docs/releases/v0.1.0.md`, `runtime/apply.go`, `bind/doc.go`

- [ ] **Step 1: Write the product-boundary explanation**

Required concepts:
- `vxt` as a spec-first library
- why `vxt` is library-only
- relation to `vx`

- [ ] **Step 2: Explain the two primary usage models**

Cover:
- single-file rendering
- document mode
- raw runtime API vs generated bindings

- [ ] **Step 3: Explain hooks and execution model**

Clarify:
- planned hooks
- `WritePlan`
- `ApplyPlan`
- why trust/shell policy belongs outside `vxt`

- [ ] **Step 4: Commit**

```bash
git add docs/concepts.md
git commit -m "docs: add vxt concepts guide"
```

### Task 7: Final docs-map integration and verification

**Files:**
- Modify: `README.md`
- Verify: `docs/getting-started.md`, `docs/document-mode.md`, `docs/runtime-api.md`, `docs/go-bindings.md`, `docs/concepts.md`

- [ ] **Step 1: Re-check README links and docs map**

Ensure every new doc is linked correctly from the README.

- [ ] **Step 2: Review all docs for terminology consistency**

Check:
- `vxt` vs `vx`
- document mode vs single-file mode
- runtime API vs bindings
- planned hooks vs executed hooks

- [ ] **Step 3: Run repository verification**

Run:

```bash
go test ./...
go doc github.com/vandordev/vxt
go doc github.com/vandordev/vxt/runtime
go doc github.com/vandordev/vxt/bind
```

Expected:
- tests pass
- package docs render without obvious symbol drift

- [ ] **Step 4: Manually scan docs for broken claims**

Verify that no document promises:
- CLI behavior inside `vxt`
- package registry semantics
- automatic shell execution beyond the current explicit hook executor model

- [ ] **Step 5: Commit**

```bash
git add README.md docs/getting-started.md docs/document-mode.md docs/runtime-api.md docs/go-bindings.md docs/concepts.md
git commit -m "docs: add first complete vxt documentation set"
```

