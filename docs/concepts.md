# Concepts

`vxt` is a spec-first Go library for code and file generation. It defines the
template contract, runtime lifecycle, diagnostics, output targets, and generated
Go bindings.

The core idea is simple: template authors declare what should exist, and Go
callers decide when and how to compile, validate, plan, write, and optionally
apply post-write hooks.

## Library Boundary

`vxt` is library-only. It does not ship a CLI, watch process, project command,
trust prompt, or package registry.

That boundary keeps the public package focused on deterministic generation
building blocks:

- parse `.vxt` source
- validate caller input
- render a concrete plan
- write the plan through an explicit target
- expose hook metadata
- generate typed Go bindings

Applications built on top of `vxt` own user interaction, trust policy, shell
policy, command orchestration, and remote distribution.

## Relation to vx

`vx` is the expected higher-level Vandor tool surface. `vxt` is the reusable Go
library underneath template compilation and generation behavior.

In practice:

- `vxt` should contain the stable generation primitives.
- `vx` can decide how users discover, install, trust, and run templates.
- `vxt` should not grow CLI behavior just because a CLI might need it.

This split keeps `vxt` embeddable for Go consumers that do not want the full
tooling product.

## Single-File vs Document Mode

Single-file mode renders one template string through `vxt.RenderSingleFile`.
It is useful for small embedded rendering cases.

Document mode is the primary model for repository and code generation. It
supports declared inputs, local types, directories, files, partials, local use
definitions, conditional outputs, and planned hooks.

Choose single-file mode for one rendered string. Choose document mode when the
template describes a file tree or a reusable generation contract.

## Runtime API vs Generated Bindings

The raw runtime API works with `source.Source`, `map[string]any`, and staged
result values. It is best when templates are dynamic or the caller wants direct
control over compile, validate, plan, write, and apply.

Generated Go bindings are created from document-mode templates through the
`bind` package. They produce a typed Go package with an `Input` struct, public
types from `@type`, and wrappers for common runtime stages.

Use runtime APIs for flexibility. Use generated bindings when a template is part
of a Go module and callers should get typed input at compile time.

## Planned Hooks and Executed Hooks

`@hook` declarations become `runtime.PlannedHook` values during planning. They
are metadata in the plan.

`runtime.WritePlan` writes directories and files only. It does not execute
hooks.

`runtime.ApplyPlan` is the explicit post-write path. It first writes the plan,
then executes supported planned hooks through the caller's
`runtime.HookExecutor`. The current supported event is `after:write`.

The executor is intentionally injected by the caller. That is where trust
policy, command allowlists, shell choice, environment, logging, and sandboxing
belong. `vxt` does not decide whether a hook string is safe to run.

## What Belongs in vxt

Good fits for `vxt`:

- document syntax and compile behavior
- validation and diagnostics
- plan shape and write targets
- explicit hook metadata and executor interfaces
- generated Go binding support
- package documentation for Go consumers

Poor fits for `vxt`:

- CLI commands
- remote template registries
- interactive trust prompts
- automatic shell execution
- editor integrations
- non-Go generated binding promises

Those higher-level concerns can be built around `vxt` without becoming part of
the core library contract.
