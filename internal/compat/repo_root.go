package compat

import (
	"fmt"
	"os"
	"path/filepath"
)

// FindRepoRoot walks up from start until it finds go.mod and compat/manifests/.
// Tests use this instead of fixed ".." segments so go test works when the cwd
// is the module root or the package directory.
func FindRepoRoot(start string) (string, error) {
	dir, err := filepath.Abs(start)
	if err != nil {
		return "", err
	}
	for {
		if isFile(filepath.Join(dir, "go.mod")) && isDir(filepath.Join(dir, "compat", "manifests")) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("repo root not found from %s", start)
		}
		dir = parent
	}
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
