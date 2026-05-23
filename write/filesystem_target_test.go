package write_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alfariiizi/vxt/write"
)

func TestFilesystemTargetWritesInsideTempDir(t *testing.T) {
	dir := t.TempDir()
	target := write.NewFilesystemTarget(dir)

	written, err := target.WriteFile("hello.txt", []byte("Hello"), "create")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !written {
		t.Fatal("expected file to be written")
	}

	data, err := os.ReadFile(filepath.Join(dir, "hello.txt"))
	if err != nil {
		t.Fatalf("expected file to exist: %v", err)
	}
	if string(data) != "Hello" {
		t.Fatalf("got %q", string(data))
	}
}
