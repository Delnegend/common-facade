package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Delnegend/common-facade/internal/config"
	"github.com/Delnegend/common-facade/internal/generator"
	"github.com/Delnegend/common-facade/internal/template"
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

func TestParseConfig_Valid(t *testing.T) {
	cfg, err := config.Parse([]byte(sampleValidYAML))
	if err != nil {
		t.Fatalf("unexpected error parsing valid YAML: %v", err)
	}

	if cfg.Name != "Differ" {
		t.Errorf("expected Name 'Differ', got '%s'", cfg.Name)
	}
	if len(cfg.Highlights) != 2 {
		t.Errorf("expected 2 highlights, got %d", len(cfg.Highlights))
	}
	if cfg.License.Unlicensed {
		t.Errorf("expected project to be licensed")
	}

	readme := template.BuildReadme(cfg)
	if !strings.Contains(readme, "# Differ") {
		t.Errorf("expected README to contain '# Differ'")
	}
	if !strings.Contains(readme, "## Documentation") {
		t.Errorf("expected README to contain '## Documentation'")
	}
	if !strings.Contains(readme, "docs/architecture.md") {
		t.Errorf("expected README to contain docs/architecture.md")
	}
	if !strings.Contains(readme, "docs/development.md") {
		t.Errorf("expected README to contain docs/development.md")
	}
	if !strings.Contains(readme, "docs/deployment.md") {
		t.Errorf("expected README to contain docs/deployment.md")
	}
	if !strings.Contains(readme, "Dual-licensed under [MIT](LICENSE-MIT) or [Apache-2.0](LICENSE-APACHE)") {
		t.Errorf("expected dual license text in README")
	}

	year := time.Now().UTC().Year()
	mit := template.BuildLicenseMIT(cfg, year)
	if !strings.Contains(mit, "Delnegend") {
		t.Errorf("expected Delnegend in MIT license")
	}
	apache := template.BuildLicenseApache(cfg, year)
	if !strings.Contains(apache, "Copyright") || !strings.Contains(apache, "Delnegend") {
		t.Errorf("expected Delnegend in Apache license")
	}
}

func TestParseConfig_Unlicensed(t *testing.T) {
	unlicensedYAML := `$schema: "./README.schema.json"
name: ZTNA
repo: Delnegend/ztna
tagline: "Internal network agent."
badges:
  ci: true
  release: false
  license: true
quick_start: "ztna connect"
highlights:
  - "Encrypted tunnels"
documentation:
  architecture: "Tunnel design."
  development: "Local dev."
  deployment: "Gateway deployment."
license:
  type: false
`
	cfg, err := config.Parse([]byte(unlicensedYAML))
	if err != nil {
		t.Fatalf("unexpected error parsing unlicensed YAML: %v", err)
	}
	if !cfg.License.Unlicensed {
		t.Fatalf("expected Unlicensed to be true")
	}

	readme := template.BuildReadme(cfg)
	if strings.Contains(readme, "[![License]") {
		t.Errorf("unlicensed project should not have license badge")
	}
	if !strings.Contains(readme, "UNLICENSED: internal lab tool.") {
		t.Errorf("expected 'UNLICENSED: internal lab tool.' in README")
	}
}

func TestParseConfig_RejectsCustomSections(t *testing.T) {
	badYAML := sampleValidYAML + "\ncustom_section: bad\n"
	_, err := config.Parse([]byte(badYAML))
	if err == nil {
		t.Fatalf("expected error for custom section, got nil")
	}
	if !strings.Contains(err.Error(), "custom_section") {
		t.Errorf("expected error message to mention 'custom_section', got: %v", err)
	}
}

func TestParseConfig_RejectsLicenseYear(t *testing.T) {
	badYAML := strings.Replace(sampleValidYAML, `holder: "Delnegend"`, "year: 2026\n  holder: \"Delnegend\"", 1)
	_, err := config.Parse([]byte(badYAML))
	if err == nil {
		t.Fatalf("expected error for license.year, got nil")
	}
	if !strings.Contains(err.Error(), "license.year") {
		t.Errorf("expected error message to mention 'license.year', got: %v", err)
	}
}

func TestParseConfig_RejectsMissingDocumentationKey(t *testing.T) {
	badYAML := strings.Replace(sampleValidYAML, "deployment:", "invalid_doc:", 1)
	_, err := config.Parse([]byte(badYAML))
	if err == nil {
		t.Fatalf("expected error for missing documentation key, got nil")
	}
}

func TestFileGenerationAndCheck(t *testing.T) {
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

	// Verify drift check passes
	mismatched, err := generator.CheckDrift(tmpDir, filesToGenerate)
	if err != nil {
		t.Fatalf("unexpected error during drift check: %v", err)
	}
	if len(mismatched) > 0 {
		t.Errorf("expected no drift, got: %v", mismatched)
	}

	// Verify files exist
	for _, p := range []string{filepath.Join(tmpDir, "README.md"), filepath.Join(tmpDir, "LICENSE-MIT"), filepath.Join(tmpDir, "LICENSE-APACHE")} {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("expected file %s to exist", p)
		}
	}
}
