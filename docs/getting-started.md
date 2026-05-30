# Getting Started

This guide is for Go developers using `vxt` as a library. It walks through one
document-mode template from source text to written output.

The goal is a first successful flow:

1. write a small `.vxt` document
2. compile it
3. validate typed input
4. plan output
5. write to memory
6. switch to filesystem output

## First Document

Create a document-mode template with a name, one input, and one file output:

```vxt
@template hello
@input name string

@file "hello.txt"
Hello {{ name }}
@endfile
```

`@template` names the document. `@input` declares the value callers must
provide. `@file` declares one planned file and its rendered body.

## Compile

Use `runtime.CompileDocument` to parse the document and get a compiled contract.

```go
src := source.Source{
	ID: "hello.vxt",
	Text: "@template hello\n" +
		"@input name string\n" +
		"@file \"hello.txt\"\n" +
		"Hello {{ name }}\n" +
		"@endfile\n",
}

compiled := runtime.CompileDocument(src)
if len(compiled.Diagnostics) > 0 {
	panic(compiled.Diagnostics[0].Message)
}
```

Compile diagnostics report syntax and document parsing failures. Do not pass a
nil `compiled.Document` to later stages.

## Validate

Use `runtime.ValidateDocument` with the caller-provided input map.

```go
validated := runtime.ValidateDocument(compiled.Document, map[string]any{
	"name": "Vandor",
})
if len(validated.Diagnostics) > 0 {
	panic(validated.Diagnostics[0].Message)
}
```

Validation checks declared `@input` values against the document's known scalar
and object types.

## Plan

Use `runtime.PlanDocument` to render concrete outputs without writing them.

```go
planned := runtime.PlanDocument(validated)
if len(planned.Diagnostics) > 0 {
	panic(planned.Diagnostics[0].Message)
}
```

The plan contains rendered file paths, file contents, directory outputs, and
planned hook metadata.

## Write with MemoryTarget

Start with `write.NewMemoryTarget()` when learning, testing, or previewing
output in-process.

```go
target := write.NewMemoryTarget()
report, err := runtime.WritePlan(planned.Plan, target)
if err != nil {
	panic(err)
}

fmt.Println(report.FilesWritten)
fmt.Println(string(target.Files()["hello.txt"]))
```

`MemoryTarget` records files in memory and does not touch disk. Its `Files`
method returns file contents keyed by planned relative path.

## Move to FilesystemTarget

Switch to `write.NewFilesystemTarget(root)` when the generated files should be
created on disk.

```go
target := write.NewFilesystemTarget("./out")
report, err := runtime.WritePlan(planned.Plan, target)
if err != nil {
	panic(err)
}

fmt.Println(report.FilesWritten)
```

`FilesystemTarget` roots all writes under the directory you provide. It creates
needed directories, rejects absolute paths and `..` escapes, and applies the
file mode from each `@file` declaration.

By default, `@file` uses `mode=create`, so writing fails if the file already
exists. Use `mode=overwrite` or `mode=skip-if-exists` when the document should
request different write behavior.

## Common First Mistakes

- Treating `vxt` as a CLI. It is a Go library; your program calls the runtime.
- Skipping diagnostic checks. Each stage can return diagnostics that should stop
  the pipeline.
- Calling `WritePlan` and expecting hooks to run. `WritePlan` writes only files
  and directories.
- Passing filesystem paths directly from untrusted input. `FilesystemTarget`
  rejects path escapes, but callers still own their generation policy.
- Forgetting `@endfile`. File blocks must be closed.

## Next Docs

- Read [Document mode](document-mode.md) to author richer `.vxt` documents.
- Read [Runtime API](runtime-api.md) for the full lifecycle and diagnostics
  model.
- Read [Go bindings](go-bindings.md) when you want generated typed wrappers
  instead of the raw runtime API.
