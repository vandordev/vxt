package runtime_test

import (
	"reflect"
	"testing"

	"github.com/vandordev/vxt/runtime"
)

func TestPublicRuntimeTypesDoNotExposeInternalPackagePaths(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		value any
	}{
		{name: "CompiledTemplate", value: runtime.CompiledTemplate{}},
		{name: "CompiledDocument", value: runtime.CompiledDocument{}},
		{name: "InputDecl", value: runtime.InputDecl{}},
		{name: "TypeFieldDecl", value: runtime.TypeFieldDecl{}},
		{name: "TypeDecl", value: runtime.TypeDecl{}},
		{name: "FileBlock", value: runtime.FileBlock{}},
		{name: "DirBlock", value: runtime.DirBlock{}},
		{name: "PartialDecl", value: runtime.PartialDecl{}},
		{name: "UseDecl", value: runtime.UseDecl{}},
		{name: "ConditionalBlock", value: runtime.ConditionalBlock{}},
		{name: "HookDecl", value: runtime.HookDecl{}},
		{name: "DirOutput", value: runtime.DirOutput{}},
		{name: "FileOutput", value: runtime.FileOutput{}},
		{name: "PlannedHook", value: runtime.PlannedHook{}},
		{name: "Plan", value: runtime.Plan{}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			typ := reflect.TypeOf(tc.value)
			if got, want := typ.PkgPath(), "github.com/vandordev/vxt/runtime"; got != want {
				t.Fatalf("got package path %q, want %q", got, want)
			}
		})
	}
}
