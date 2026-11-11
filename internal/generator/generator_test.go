package generator_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Delnegend/common-facade/internal/config"
	"github.com/Delnegend/common-facade/internal/generator"
)

const sampleValidYAML = `$schema: "./README.schema.json"
name: Differ
repo: Delnegend/differ
tagline: "Pixel-perfect image comparison CLI."
badges:
  ci: true
  release: true
  license: true
quick_start: |
  ` + "```bash" + `
  ./differ -diff a.png b.png
  ` + "```" + `
highlights:
  - "Sub-pixel antialiasing detection"
  - "Zero external dependencies"
documentation:
  architecture: "Core diffing algorithms and pipeline."
  development: "How to run benchmarks and test suites."
  deployment: "Releases and container builds."
license:
  type:
    - MIT
    - APACHE-2.0
  holder: "Delnegend"
`

func TestGenerator_WriteAndCheck(t *testing.T) {
	tmpDir := t.TempDir()
	ymlPath := filepath.Join(tmpDir, "README.yml")
	if err := os.WriteFile(ymlPath, []byte(sampleValidYAML), 0o644); err != nil {
		t.Fatalf("failed to write test YAML: %v", err)
	}

	cfg, err := config.Parse([]byte(sampleValidYAML))
	if err != nil {
		t.Fatalf("failed to parse config: %v", err)
	}

	year := time.Now().UTC().Year()
	filesToGenerate := generator.FilesToGenerate(cfg, year)

	written, err := generator.WriteFiles(tmpDir, filesToGenerate)
	if err != nil {
		t.Fatalf("failed to write files: %v", err)
	}
	if len(written) != 3 {
		t.Fatalf("expected 3 files written, got %d", len(written))
	}

	mismatched, err := generator.CheckDrift(tmpDir, filesToGenerate)
	if err != nil {
		t.Fatalf("unexpected error during drift check: %v", err)
	}
	if len(mismatched) > 0 {
		t.Errorf("expected no drift, got: %v", mismatched)
	}

	for _, p := range []string{filepath.Join(tmpDir, "README.md"), filepath.Join(tmpDir, "LICENSE-MIT"), filepath.Join(tmpDir, "LICENSE-APACHE")} {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("expected file %s to exist", p)
		}
	}
}

func TestResolveConfigPath(t *testing.T) {
	tmpDir := t.TempDir()
	ymlPath := filepath.Join(tmpDir, "README.yml")
	if err := os.WriteFile(ymlPath, []byte("name: test"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Resolve direct file
	resolved, err := generator.ResolveConfigPath(ymlPath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != ymlPath {
		t.Errorf("expected %s, got %s", ymlPath, resolved)
	}

	// Resolve directory
	resolvedDir, err := generator.ResolveConfigPath(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolvedDir != ymlPath {
		t.Errorf("expected %s, got %s", ymlPath, resolvedDir)
	}
}
