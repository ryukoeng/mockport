package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/albert-einshutoin/mockport/internal/builtins"
	"github.com/albert-einshutoin/mockport/internal/compat"
)

func TestWriteManifestsCreatesFiveFiles(t *testing.T) {
	outDir := t.TempDir()
	if err := compat.WriteManifestFiles(outDir, builtins.ManifestAdapters()); err != nil {
		t.Fatalf("WriteManifestFiles: %v", err)
	}
	want := []string{"stripe.json", "openai.json", "github-oauth.json", "slack.json", "line.json"}
	for _, name := range want {
		path := filepath.Join(outDir, name)
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("missing manifest %s: %v", name, err)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", name, err)
		}
		if data[len(data)-1] != '\n' {
			t.Fatalf("%s missing trailing newline", name)
		}
	}
}

func TestWriteManifestsIsDeterministic(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()
	adapters := builtins.ManifestAdapters()
	if err := compat.WriteManifestFiles(dir1, adapters); err != nil {
		t.Fatalf("WriteManifestFiles dir1: %v", err)
	}
	if err := compat.WriteManifestFiles(dir2, adapters); err != nil {
		t.Fatalf("WriteManifestFiles dir2: %v", err)
	}
	names := []string{"stripe.json", "openai.json", "github-oauth.json", "slack.json", "line.json"}
	for _, name := range names {
		a, err := os.ReadFile(filepath.Join(dir1, name))
		if err != nil {
			t.Fatalf("read dir1/%s: %v", name, err)
		}
		b, err := os.ReadFile(filepath.Join(dir2, name))
		if err != nil {
			t.Fatalf("read dir2/%s: %v", name, err)
		}
		if !bytes.Equal(a, b) {
			t.Fatalf("manifest %s differs between runs", name)
		}
	}
}

func TestCheckedInManifestsMatchGenerated(t *testing.T) {
	repoRoot, err := compat.FindRepoRoot(".")
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	generatedDir := t.TempDir()
	if err := compat.WriteManifestFiles(generatedDir, builtins.ManifestAdapters()); err != nil {
		t.Fatalf("WriteManifestFiles: %v", err)
	}
	for _, adapterImpl := range builtins.ManifestAdapters() {
		name := adapterImpl.Name() + ".json"
		checkedIn, err := os.ReadFile(filepath.Join(repoRoot, "compat", "manifests", name))
		if err != nil {
			t.Fatalf("read checked-in %s: %v", name, err)
		}
		generated, err := os.ReadFile(filepath.Join(generatedDir, name))
		if err != nil {
			t.Fatalf("read generated %s: %v", name, err)
		}
		want, err := compat.ManifestFileBytesFromJSON(checkedIn)
		if err != nil {
			t.Fatalf("canonicalize checked-in %s: %v", name, err)
		}
		if !bytes.Equal(want, generated) {
			t.Fatalf(
				"manifest drift for %s: run UPDATE_COMPAT_MANIFESTS=1 go test ./internal/compat -run TestUpdateCheckedInManifests",
				adapterImpl.Name(),
			)
		}
	}
}
