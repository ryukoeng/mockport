package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/albert-einshutoin/mockport/internal/builtins"
	"github.com/albert-einshutoin/mockport/internal/compat"
)

// TestUpdateRepoManifests regenerates compat/manifests/*.json when run with:
//
//	UPDATE_COMPAT_MANIFESTS=1 go test ./scripts/gen-compat-manifests -run TestUpdateRepoManifests
func TestUpdateRepoManifests(t *testing.T) {
	if os.Getenv("UPDATE_COMPAT_MANIFESTS") != "1" {
		t.Skip("set UPDATE_COMPAT_MANIFESTS=1 to regenerate checked-in manifests")
	}
	repoRoot, err := compat.FindRepoRoot(".")
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	outDir := filepath.Join(repoRoot, "compat", "manifests")
	if err := compat.WriteManifestFiles(outDir, builtins.ManifestAdapters()); err != nil {
		t.Fatalf("WriteManifestFiles: %v", err)
	}
}
