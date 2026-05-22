# vxt

`vxt` is an independent product under the Vandor organization.

It is a spec-first template compiler/runtime product for code and file
generation. The canonical public model is a staged pipeline:

1. compile
2. validate
3. plan
4. write

For v0.1, `vxt` should prove:

- single-file rendering through a simple convenience API
- document-mode compile, validate, plan, and write
- typed document input validation
- structured diagnostics
- output-target abstraction with filesystem as one adapter

Explicit v0.1 non-goals:

- hook execution
- trust policy
- registry or package semantics
- CLI behavior
- AST manipulation as a public contract
