package compat

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/albert-einshutoin/mockport/internal/adapter"
)

// ManifestFileBytesFromJSON parses checked-in manifest JSON and returns the
// canonical encoded bytes used by scripts/gen-compat-manifests.
func ManifestFileBytesFromJSON(data []byte) ([]byte, error) {
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	return ManifestFileBytes(manifest)
}

// ManifestFileBytes returns canonical checked-in manifest JSON (indent + trailing newline).
// Output matches scripts/gen-compat-manifests so CI can diff file bytes deterministically.
func ManifestFileBytes(m Manifest) ([]byte, error) {
	if err := m.Validate(); err != nil {
		return nil, err
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return nil, err
	}
	data = append(data, '\n')
	return data, nil
}

// WriteManifestFiles writes {name}.json for each adapter using FromMetadata(metadata).
func WriteManifestFiles(outDir string, adapters []adapter.Adapter) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", outDir, err)
	}
	for _, adapterImpl := range adapters {
		manifest := FromMetadata(adapterImpl.Metadata())
		data, err := ManifestFileBytes(manifest)
		if err != nil {
			return fmt.Errorf("%s manifest: %w", adapterImpl.Name(), err)
		}
		path := filepath.Join(outDir, adapterImpl.Name()+".json")
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
	}
	return nil
}
