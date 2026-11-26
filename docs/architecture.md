# Architecture

`common-facade` is a standalone CLI tool written in Go that reads `README.yml`, validates its fields against a strict layout specification, and generates standardized `README.md`, `LICENSE-MIT`, and `LICENSE-APACHE` files.

## Workflow Pipeline

1. **Discovery & Resolution**: Resolves target directory or path to `README.yml`.
2. **Schema & Layout Enforcement**: Uses strict AST-level YAML decoding (`KnownFields(true)`) to reject any custom or unrecognized sections. Requires:
   - `name`, `repo`, `tagline`
   - `badges` (`ci`, `release`, `license`)
   - `quick_start`
   - `highlights`
   - `documentation` (`architecture`, `development`, `deployment`)
   - `license` (`type: false` or `["MIT", "APACHE-2.0"]`)
3. **Template Rendering**: Builds Markdown with centered headers, badgers.space badges, and standardized documentation links.
4. **License Generation**: Uses full text templates for MIT and Apache 2.0 licenses, stamping the current UTC year automatically.
5. **Drift Verification (`check`)**: Compares memory-rendered buffers against on-disk files. Exits with code 1 if files are missing or modified.

## Module Structure

The project follows the standard Go `internal/` package layout to keep concerns separated:

- `main.go`: Thin entrypoint dispatching CLI flags and arguments.
- `internal/config`: YAML schema parsing, AST decoding, custom unmarshalers, and structural validation.
- `internal/template`: Markdown and license generation logic and embedded license templates.
- `internal/generator`: Target discovery, file diffing / drift detection (`check`), and file writing.
