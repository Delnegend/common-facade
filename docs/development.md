# Development

## Prerequisites

- Go 1.24+
- `just` command runner

## Local Workflow

Run the test suite, code formatting checks, and vet checks:

```bash
just check
```

## Adding Changes

1. Keep dependencies minimal (only `gopkg.in/yaml.v3`).
2. Run `just check` to verify `gofmt`, `go vet`, and unit tests pass cleanly.
3. Keep `README.yml` up to date and run `common-facade .` to regenerate `README.md` and licenses.
