// Package vxt provides a spec-first template compiler/runtime for code and
// file generation.
//
// The canonical product model is a staged pipeline:
// compile -> validate -> plan -> write.
//
// v0.1 intentionally excludes hook execution, trust policy, package semantics,
// and CLI behavior from the core contract.
package vxt
