package compat

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/albert-einshutoin/mockport/internal/builtins"
)

// TestCheckedInManifestsMatchFromMetadata ensures compat/manifests/*.json stay
// aligned with runtime FromMetadata output. Regenerate with:
//
//	UPDATE_COMPAT_MANIFESTS=1 go test ./internal/compat -run TestUpdateCheckedInManifests
func TestCheckedInManifestsMatchFromMetadata(t *testing.T) {
	repoRoot, err := FindRepoRoot(".")
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	for _, adapterImpl := range builtins.ManifestAdapters() {
		name := adapterImpl.Name()
		path := filepath.Join(repoRoot, "compat", "manifests", name+".json")
		checkedIn, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read manifest %s: %v", path, err)
		}
		want, err := ManifestFileBytes(FromMetadata(adapterImpl.Metadata()))
		if err != nil {
			t.Fatalf("encode generated manifest %s: %v", name, err)
		}
		got, err := ManifestFileBytesFromJSON(checkedIn)
		if err != nil {
			t.Fatalf("canonicalize checked-in manifest %s: %v", name, err)
		}
		if bytes.Equal(want, got) {
			continue
		}
		t.Fatalf(
			"manifest drift for %s: run UPDATE_COMPAT_MANIFESTS=1 go test ./internal/compat -run TestUpdateCheckedInManifests",
			name,
		)
	}
}

// TestUpdateCheckedInManifests regenerates compat/manifests/*.json when run with:
//
//	UPDATE_COMPAT_MANIFESTS=1 go test ./internal/compat -run TestUpdateCheckedInManifests
func TestUpdateCheckedInManifests(t *testing.T) {
	if os.Getenv("UPDATE_COMPAT_MANIFESTS") != "1" {
		t.Skip("set UPDATE_COMPAT_MANIFESTS=1 to regenerate checked-in manifests")
	}
	repoRoot, err := FindRepoRoot(".")
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	outDir := filepath.Join(repoRoot, "compat", "manifests")
	if err := WriteManifestFiles(outDir, builtins.ManifestAdapters()); err != nil {
		t.Fatalf("WriteManifestFiles: %v", err)
	}
}
