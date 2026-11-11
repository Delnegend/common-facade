package template_test

import (
	"strings"
	"testing"
	"time"

	"github.com/Delnegend/common-facade/internal/config"
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

func TestBuildReadme(t *testing.T) {
	cfg, err := config.Parse([]byte(sampleValidYAML))
	if err != nil {
		t.Fatalf("unexpected error parsing valid YAML: %v", err)
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
}

func TestBuildLicense(t *testing.T) {
	cfg, err := config.Parse([]byte(sampleValidYAML))
	if err != nil {
		t.Fatalf("unexpected error parsing valid YAML: %v", err)
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
