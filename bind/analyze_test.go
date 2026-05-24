package bind

import (
	"testing"

	"github.com/vandordev/vxt/runtime"
	"github.com/vandordev/vxt/source"
)

func TestAnalyzeDocumentDerivesPublicGoTypesAndInput(t *testing.T) {
	src := source.Source{
		ID: "service.vxt",
		Text: "@template service\n" +
			"@type Repository {\n" +
			"  driver: string\n" +
			"}\n" +
			"@type Entity {\n" +
			"  name: string\n" +
			"  package_name: string\n" +
			"  tags: string[]\n" +
			"  repository?: Repository\n" +
			"}\n" +
			"@input entity Entity\n",
	}

	compiled := runtime.CompileDocument(src)
	if len(compiled.Diagnostics) > 0 {
		t.Fatalf("unexpected compile diagnostics: %#v", compiled.Diagnostics)
	}

	got, err := analyzeDocument("servicevxt", compiled.Document)
	if err != nil {
		t.Fatalf("unexpected analyze error: %v", err)
	}

	if got.PackageName != "servicevxt" {
		t.Fatalf("got package %q", got.PackageName)
	}
	if len(got.Types) != 2 {
		t.Fatalf("got %d types", len(got.Types))
	}

	repository := got.Types[0]
	if repository.Name != "Repository" {
		t.Fatalf("got repository type %q", repository.Name)
	}
	if len(repository.Fields) != 1 {
		t.Fatalf("got %d repository fields", len(repository.Fields))
	}
	if repository.Fields[0].GoName != "Driver" || repository.Fields[0].GoType != "string" {
		t.Fatalf("unexpected repository field: %#v", repository.Fields[0])
	}

	entity := got.Types[1]
	if entity.Name != "Entity" {
		t.Fatalf("got entity type %q", entity.Name)
	}
	if len(entity.Fields) != 4 {
		t.Fatalf("got %d entity fields", len(entity.Fields))
	}
	assertField(t, entity.Fields[0], "Name", "name", "string")
	assertField(t, entity.Fields[1], "PackageName", "package_name", "string")
	assertField(t, entity.Fields[2], "Tags", "tags", "[]string")
	assertField(t, entity.Fields[3], "Repository", "repository", "*Repository")

	if len(got.InputFields) != 1 {
		t.Fatalf("got %d input fields", len(got.InputFields))
	}
	assertField(t, got.InputFields[0], "Entity", "entity", "Entity")
}

func TestGoTypeForFieldMapsOptionalAndArrayShapes(t *testing.T) {
	tests := []struct {
		name  string
		field runtime.TypeFieldDecl
		want  string
	}{
		{
			name:  "primitive string",
			field: runtime.TypeFieldDecl{Name: "name", TypeName: "string"},
			want:  "string",
		},
		{
			name:  "named type",
			field: runtime.TypeFieldDecl{Name: "entity", TypeName: "Entity"},
			want:  "Entity",
		},
		{
			name:  "string array",
			field: runtime.TypeFieldDecl{Name: "tags", TypeName: "string", Array: true},
			want:  "[]string",
		},
		{
			name:  "optional named type",
			field: runtime.TypeFieldDecl{Name: "repository", TypeName: "Repository", Optional: true},
			want:  "*Repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := goTypeForField(tt.field)
			if got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}

func assertField(t *testing.T, got analyzedField, wantName, wantSchema, wantType string) {
	t.Helper()
	if got.GoName != wantName || got.SchemaName != wantSchema || got.GoType != wantType {
		t.Fatalf("got field %#v want name=%q schema=%q type=%q", got, wantName, wantSchema, wantType)
	}
}
