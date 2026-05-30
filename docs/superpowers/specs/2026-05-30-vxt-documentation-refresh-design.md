# vxt Documentation Refresh Design

## Summary

This spec defines the first cohesive documentation set for `vxt` as a public Go
library. The goal is to make the repository self-explanatory for a new reader
who lands on GitHub or `pkg.go.dev`, understands what `vxt` is and is not, and
can move from overview to successful use without reverse-engineering the code.

The documentation set must prioritize:

- Go consumers first
- document mode as the primary usage path
- clear library boundaries
- stable links from `README.md` into focused docs under `docs/`

This refresh is intentionally repository-local. `vandor-landing` may later link
to these docs, but it is not the source of truth for technical usage.

## Problem

The current `README.md` is useful but still too compressed for a first-time
reader who needs to understand:

- what `vxt` does
- when to use `vxt`
- how to use the runtime pipeline
- how to author `.vxt` documents
- how generated Go bindings fit into the library
- how `vxt` relates to `vx`

The repo also lacks a complete, intentional docs map. As a result:

- beginners do not have a guided first path
- concepts and reference concerns are mixed together
- document-mode features exist in code but are not explained in one place
- bindings are implemented but not yet documented as a complete workflow

## Audience

Primary audience:

- Go developers evaluating or adopting `vxt` as a library

Secondary audience:

- template authors writing `.vxt` documents
- maintainers or contributors who need to understand product boundaries

## Goals

- Make `README.md` explain `vxt` in 1-2 minutes.
- Provide a first successful path for document-mode usage.
- Provide a clear guide for authoring `.vxt` documents.
- Provide lifecycle-oriented runtime API docs.
- Provide practical docs for typed Go bindings.
- Explain the conceptual boundary between `vxt`, `vx`, runtime mode, and
  bindings.

## Non-Goals

- public website docs
- documentation for internal packages under `internal/`
- advanced multi-language workflow documentation
- release-process overhaul
- contributor onboarding overhaul
- exhaustive reference docs for every future directive or feature concept

## Documentation Strategy

This refresh follows a docs-first-with-sharp-README model.

`README.md` remains the entry point. It should stay scannable and should not
become the full manual.

Detailed usage and conceptual material lives in focused docs under `docs/`.

The set should loosely align with Diátaxis:

- Tutorial:
  - `docs/getting-started.md`
  - most of `docs/document-mode.md`
- How-to:
  - `docs/go-bindings.md`
- Reference:
  - `docs/runtime-api.md`
  - compact directive reference section inside `docs/document-mode.md`
- Explanation:
  - `docs/concepts.md`

## Planned Files

### `README.md`

Purpose:
- front page for GitHub and package readers

Required content:
- what `vxt` is
- when to use it
- what it is not
- install
- document-mode quick start
- pipeline overview
- feature snapshot
- docs map
- release status and supported public packages

### `docs/getting-started.md`

Purpose:
- first successful end-to-end experience for a Go consumer

Required content:
- who this guide is for
- minimal document-mode example
- `CompileDocument`
- `ValidateDocument`
- `PlanDocument`
- `WritePlan`
- `MemoryTarget` first
- short transition to `FilesystemTarget`
- common first mistakes
- where to go next

### `docs/document-mode.md`

Purpose:
- teach how to author `.vxt` documents

Required content:
- one realistic example template
- guided explanation of:
  - `@template`
  - `@input`
  - `@type`
  - `@dir`
  - `@file`
  - `@partial`
  - `@use`
  - `@if`
  - `@hook`
- final planned/rendered result
- compact directive reference
- explicit current limitations

### `docs/runtime-api.md`

Purpose:
- explain the public runtime lifecycle as a library API

Required content:
- lifecycle overview
- `CompileDocument`
- `CompileDocumentWithResolver`
- `ValidateDocument`
- `PlanDocument`
- `WritePlan`
- `WritePlanWithDiagnostics`
- `ApplyPlan`
- `OutputTarget`
- diagnostics model
- guidance on choosing API depth

### `docs/go-bindings.md`

Purpose:
- explain typed Go bindings generated from document-mode `.vxt`

Required content:
- what generated bindings are
- when to use them instead of raw runtime APIs
- generated package shape
- typed `Input`
- `Compile`, `Validate`, `Plan`, `Write`
- `bind.Generate`
- `bind.WriteOutput`
- `bind.GenerateToDir`
- dry-run behavior
- local `@use` embedding behavior
- current non-goals and limitations

### `docs/concepts.md`

Purpose:
- explain `vxt` product boundaries and mental model

Required content:
- `vxt` as a spec-first library
- why `vxt` is library-only
- relation to `vx`
- single-file vs document mode
- runtime API vs generated bindings
- planned hooks vs executed hooks
- what belongs in `vxt` vs `vx`

## README Shape

The README should be ordered for fast comprehension:

1. what `vxt` is
2. when to use it
3. what it is not
4. install
5. strong document-mode quick start
6. pipeline overview
7. feature snapshot
8. docs map
9. release/public-package note

The quick start should favor document mode over single-file rendering. The
single-file path may remain in the repo docs, but it should not define the main
identity of the README.

## Writing Priorities

Implementation order:

1. `README.md`
2. `docs/getting-started.md`
3. `docs/document-mode.md`
4. `docs/runtime-api.md`
5. `docs/go-bindings.md`
6. `docs/concepts.md`

This order ensures the repo landing page and first-use path are fixed before the
deeper reference and explanation docs.

## Style Rules

- Prefer direct, non-marketing language.
- Use current `vxt` terminology consistently.
- Distinguish current implemented behavior from future ideas.
- Avoid claiming support for behavior that is only planned.
- Prefer short runnable examples over abstract descriptions.
- Keep examples aligned with the current public API.

## Acceptance Criteria

This documentation refresh is successful when:

- a new reader can identify `vxt` as a Go library, not a CLI
- a reader can complete a first document-mode flow from the docs
- a reader can understand the difference between runtime API and generated
  bindings
- a reader can find where `.vxt` document directives are explained
- a reader can understand current hook behavior, including `ApplyPlan(...)`
- README links clearly to the rest of the docs set

