package config_test

import (
	"strings"
	"testing"

	"github.com/Delnegend/common-facade/internal/config"
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
