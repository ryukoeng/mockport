// Usage: go run ./scripts/gen-compat-manifests [--out compat/manifests]
//
// Writes compat/manifests/{name}.json from each manifest-tracked built-in
// adapter's Metadata() via compat.FromMetadata. Output is deterministic so CI
// can diff against checked-in manifests.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/albert-einshutoin/mockport/internal/builtins"
	"github.com/albert-einshutoin/mockport/internal/compat"
)

func main() {
	outDir := flag.String("out", "compat/manifests", "output directory for manifest JSON files")
	flag.Parse()

	if err := compat.WriteManifestFiles(*outDir, builtins.ManifestAdapters()); err != nil {
		fmt.Fprintf(os.Stderr, "gen-compat-manifests: %v\n", err)
		os.Exit(1)
	}
}
