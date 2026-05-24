package bind

import (
	"strings"
	"testing"

	"github.com/vandordev/vxt/source"
)

func TestGenerateEmbedsLocalUseSources(t *testing.T) {
	main := source.Source{
		ID: "service.vxt",
		Text: "@template service\n" +
			"@use \"./schema.vxt\"\n" +
			"@input entity Entity\n" +
			"@file \"service.go\"\n" +
			"package service\n" +
			"@endfile\n",
	}

	out, err := Generate(Request{
		PackageName: "servicevxt",
		Document:    main,
		Uses: map[string]source.Source{
			"./schema.vxt": {
				ID: "schema.vxt",
				Text: "@type Entity {\n" +
					"  name: string\n" +
					"}\n",
			},
		},
	})
	if err != nil {
		t.Fatalf("unexpected generate error: %v", err)
	}
	if len(out.Files) != 1 {
		t.Fatalf("got %d files", len(out.Files))
	}

	content := out.Files[0].Content
	if !strings.Contains(content, "\"./schema.vxt\":") {
		t.Fatalf("generated content missing local use resolver entry\n\n%s", content)
	}
	if !strings.Contains(content, "@type Entity {") {
		t.Fatalf("generated content missing embedded imported type source\n\n%s", content)
	}
}
