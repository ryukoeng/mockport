package compat

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindRepoRoot(t *testing.T) {
	root, err := FindRepoRoot(".")
	if err != nil {
		t.Fatalf("FindRepoRoot: %v", err)
	}
	if !isFile(filepath.Join(root, "go.mod")) {
		t.Fatalf("go.mod missing under %s", root)
	}
	if !isDir(filepath.Join(root, "compat", "manifests")) {
		t.Fatalf("compat/manifests missing under %s", root)
	}
	if _, err := os.Stat(filepath.Join(root, "compat", "manifests", "stripe.json")); err != nil {
		t.Fatalf("stripe manifest missing: %v", err)
	}
}
