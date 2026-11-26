package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/Delnegend/common-facade/internal/config"
	"github.com/Delnegend/common-facade/internal/template"
)

type GeneratedFile struct {
	Name    string
	Content string
}

func FilesToGenerate(cfg *config.Config, year int) map[string]string {
	files := make(map[string]string)
	files["README.md"] = template.BuildReadme(cfg)

	if !cfg.License.Unlicensed {
		hasMIT := contains(cfg.License.Types, "MIT")
		hasApache := contains(cfg.License.Types, "APACHE-2.0")
		if hasMIT && hasApache {
			files["LICENSE-MIT"] = template.BuildLicenseMIT(cfg, year)
			files["LICENSE-APACHE"] = template.BuildLicenseApache(cfg, year)
		} else if hasApache {
			fname := cfg.License.Filename
			if fname == "" {
				fname = "LICENSE-APACHE"
			}
			files[fname] = template.BuildLicenseApache(cfg, year)
		} else if hasMIT {
			fname := cfg.License.Filename
			if fname == "" {
				fname = "LICENSE-MIT"
			}
			files[fname] = template.BuildLicenseMIT(cfg, year)
		}
	}
	return files
}

func ResolveConfigPath(arg string) (string, error) {
	fi, err := os.Stat(arg)
	if err != nil {
		return "", err
	}
	if fi.IsDir() {
		candidate := filepath.Join(arg, "README.yml")
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		candidateYaml := filepath.Join(arg, "README.yaml")
		if _, err := os.Stat(candidateYaml); err == nil {
			return candidateYaml, nil
		}
		return "", fmt.Errorf("no README.yml found in directory '%s'", arg)
	}
	return arg, nil
}

func CheckDrift(outDir string, filesToGenerate map[string]string) ([]string, error) {
	var mismatched []string
	for name, expected := range filesToGenerate {
		dest := filepath.Join(outDir, name)
		content, err := os.ReadFile(dest)
		if err != nil {
			mismatched = append(mismatched, fmt.Sprintf("%s (missing)", name))
			continue
		}
		if string(content) != expected {
			mismatched = append(mismatched, fmt.Sprintf("%s (content changed)", name))
		}
	}
	return mismatched, nil
}

func WriteFiles(outDir string, filesToGenerate map[string]string) ([]string, error) {
	var written []string
	for name, content := range filesToGenerate {
		dest := filepath.Join(outDir, name)
		if err := os.WriteFile(dest, []byte(content), 0o644); err != nil {
			return written, fmt.Errorf("failed to write '%s': %w", dest, err)
		}
		written = append(written, dest)
	}
	return written, nil
}

func Run(args []string) error {
	isCheck := false
	targetPath := "README.yml"

	var cleanArgs []string
	for _, a := range args {
		if a == "check" || a == "--check" {
			isCheck = true
		} else if a == "-h" || a == "--help" {
			fmt.Println("Usage: common-facade [check|--check] [PATH_OR_DIR]")
			fmt.Println("\nOptions:")
			fmt.Println("  check, --check    Validate README.md and LICENSE files without modifying")
			fmt.Println("  -h, --help        Show help")
			return nil
		} else {
			cleanArgs = append(cleanArgs, a)
		}
	}

	if len(cleanArgs) > 0 {
		targetPath = cleanArgs[0]
	}

	resolved, err := ResolveConfigPath(targetPath)
	if err != nil {
		return fmt.Errorf("cannot find configuration: %w", err)
	}

	data, err := os.ReadFile(resolved)
	if err != nil {
		return fmt.Errorf("cannot read file '%s': %w", resolved, err)
	}

	cfg, err := config.Parse(data)
	if err != nil {
		return fmt.Errorf("configuration error in '%s': %w", resolved, err)
	}

	outDir := filepath.Dir(resolved)
	currentYear := time.Now().UTC().Year()

	filesToGenerate := FilesToGenerate(cfg, currentYear)

	if isCheck {
		mismatched, err := CheckDrift(outDir, filesToGenerate)
		if err != nil {
			return err
		}
		if len(mismatched) > 0 {
			fmt.Fprintf(os.Stderr, "Drift detected! The following files differ from %s:\n", filepath.Base(resolved))
			for _, m := range mismatched {
				fmt.Fprintf(os.Stderr, "  - %s\n", m)
			}
			os.Exit(1)
		}
		fmt.Printf("All generated files match %s.\n", filepath.Base(resolved))
		return nil
	}

	written, err := WriteFiles(outDir, filesToGenerate)
	if err != nil {
		return err
	}
	for _, path := range written {
		fmt.Printf("Generated %s\n", path)
	}

	return nil
}

func contains(slice []string, val string) bool {
	return slices.Contains(slice, val)
}
