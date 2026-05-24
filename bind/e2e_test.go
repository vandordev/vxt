package bind_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vandordev/vxt/bind"
	"github.com/vandordev/vxt/source"
)

func TestGeneratedBindingsCompileAndPlan(t *testing.T) {
	repoRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}

	serviceText, err := os.ReadFile(filepath.Join("testdata", "service", "service.vxt"))
	if err != nil {
		t.Fatalf("read service template: %v", err)
	}
	schemaText, err := os.ReadFile(filepath.Join("testdata", "service", "schema.vxt"))
	if err != nil {
		t.Fatalf("read schema template: %v", err)
	}

	out, err := bind.Generate(bind.Request{
		PackageName: "servicevxt",
		Document: source.Source{
			ID:   "service.vxt",
			Path: "service.vxt",
			Text: string(serviceText),
		},
		Uses: map[string]source.Source{
			"./schema.vxt": {
				ID:   "schema.vxt",
				Path: "schema.vxt",
				Text: string(schemaText),
			},
		},
	})
	if err != nil {
		t.Fatalf("generate bindings: %v", err)
	}
	if len(out.Files) != 1 {
		t.Fatalf("got %d files", len(out.Files))
	}

	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(strings.TrimSpace(`
module example.com/consumer

go 1.24

require github.com/vandordev/vxt v0.0.0

replace github.com/vandordev/vxt => `+repoRoot+`
`)+"\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, ".vxt"), 0o755); err != nil {
		t.Fatalf("mkdir .vxt: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, out.Files[0].Path), []byte(out.Files[0].Content), 0o644); err != nil {
		t.Fatalf("write generated file: %v", err)
	}
	mainSource := `package main

import (
	"fmt"

	servicevxt "example.com/consumer/.vxt"
)

func main() {
	plan, err := servicevxt.Plan(servicevxt.Input{
		Entity: servicevxt.Entity{
			Name:          "User",
			PackageName:   "user",
			HasRepository: true,
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(len(plan.Dirs), len(plan.Files))
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(mainSource), 0o644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	cmd := exec.Command("go", "run", ".")
	cmd.Dir = tmpDir
	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v\n%s", err, string(outBytes))
	}
	if strings.TrimSpace(string(outBytes)) != "1 2" {
		t.Fatalf("got output %q", strings.TrimSpace(string(outBytes)))
	}
}

func TestGenerateToDirWritesBindingsAndConsumerCanPlan(t *testing.T) {
	repoRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatalf("resolve repo root: %v", err)
	}

	serviceText, err := os.ReadFile(filepath.Join("testdata", "service", "service.vxt"))
	if err != nil {
		t.Fatalf("read service template: %v", err)
	}
	schemaText, err := os.ReadFile(filepath.Join("testdata", "service", "schema.vxt"))
	if err != nil {
		t.Fatalf("read schema template: %v", err)
	}

	tmpDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(strings.TrimSpace(`
module example.com/consumer

go 1.24

require github.com/vandordev/vxt v0.0.0

replace github.com/vandordev/vxt => `+repoRoot+`
`)+"\n"), 0o644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}

	report, err := bind.GenerateToDir(bind.Request{
		PackageName: "servicevxt",
		Document: source.Source{
			ID:   "service.vxt",
			Path: "service.vxt",
			Text: string(serviceText),
		},
		Uses: map[string]source.Source{
			"./schema.vxt": {
				ID:   "schema.vxt",
				Path: "schema.vxt",
				Text: string(schemaText),
			},
		},
	}, filepath.Join(tmpDir, ".vxt"), bind.WriteOptions{})
	if err != nil {
		t.Fatalf("generate to dir: %v", err)
	}
	if len(report.FilesWritten)+len(report.FilesOverwritten) != 1 {
		t.Fatalf("unexpected write report: %#v", report)
	}

	mainSource := `package main

import (
	"fmt"

	servicevxt "example.com/consumer/.vxt"
)

func main() {
	plan, err := servicevxt.Plan(servicevxt.Input{
		Entity: servicevxt.Entity{
			Name:          "User",
			PackageName:   "user",
			HasRepository: true,
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(len(plan.Dirs), len(plan.Files))
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(mainSource), 0o644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}

	cmd := exec.Command("go", "run", ".")
	cmd.Dir = tmpDir
	outBytes, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v\n%s", err, string(outBytes))
	}
	if strings.TrimSpace(string(outBytes)) != "1 2" {
		t.Fatalf("got output %q", strings.TrimSpace(string(outBytes)))
	}
}

func TestGenerateToDirDryRunReportsChangesWithoutWritingFiles(t *testing.T) {
	serviceText, err := os.ReadFile(filepath.Join("testdata", "service", "service.vxt"))
	if err != nil {
		t.Fatalf("read service template: %v", err)
	}
	schemaText, err := os.ReadFile(filepath.Join("testdata", "service", "schema.vxt"))
	if err != nil {
		t.Fatalf("read schema template: %v", err)
	}

	tmpDir := t.TempDir()
	outDir := filepath.Join(tmpDir, ".vxt")

	report, err := bind.GenerateToDir(bind.Request{
		PackageName: "servicevxt",
		Document: source.Source{
			ID:   "service.vxt",
			Path: "service.vxt",
			Text: string(serviceText),
		},
		Uses: map[string]source.Source{
			"./schema.vxt": {
				ID:   "schema.vxt",
				Path: "schema.vxt",
				Text: string(schemaText),
			},
		},
	}, outDir, bind.WriteOptions{DryRun: true})
	if err != nil {
		t.Fatalf("generate to dir dry-run: %v", err)
	}
	if !report.DryRun {
		t.Fatal("expected dry-run report")
	}
	if len(report.FilesWritten)+len(report.FilesOverwritten) != 1 {
		t.Fatalf("unexpected dry-run report: %#v", report)
	}
	if _, err := os.Stat(outDir); !os.IsNotExist(err) {
		t.Fatalf("expected dry-run to avoid creating output dir, stat err=%v", err)
	}
}
