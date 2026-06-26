package compat

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/albert-einshutoin/mockport/internal/builtins"
)

func TestManifestFileBytesRoundTrip(t *testing.T) {
	for _, adapterImpl := range builtins.ManifestAdapters() {
		if adapterImpl.Name() != "stripe" {
			continue
		}
		manifest := FromMetadata(adapterImpl.Metadata())
		data, err := ManifestFileBytes(manifest)
		if err != nil {
			t.Fatalf("ManifestFileBytes: %v", err)
		}
		if data[len(data)-1] != '\n' {
			t.Fatal("manifest file bytes missing trailing newline")
		}
		var round Manifest
		if err := json.Unmarshal(data, &round); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		again, err := ManifestFileBytes(round)
		if err != nil {
			t.Fatalf("ManifestFileBytes round: %v", err)
		}
		if !bytes.Equal(data, again) {
			t.Fatal("manifest file bytes not stable after JSON round-trip")
		}
		return
	}
	t.Fatal("stripe adapter not found in ManifestAdapters()")
}
