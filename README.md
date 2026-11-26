<div align="center">

# Common Facade

**Opinionated CLI and GitHub Action that generates and validates README.md and LICENSE files across repositories.**

[![CI](https://badgers.space/github/checks/Delnegend/common-facade/main?corner_radius=0&label=build)](https://github.com/Delnegend/common-facade/actions)
[![Release](https://badgers.space/github/release/Delnegend/common-facade?corner_radius=0&color=ea7233)](https://github.com/Delnegend/common-facade/releases)
[![License](https://badgers.space/badge/license/MIT%20%7C%20Apache-2.0/007ec6?corner_radius=0)](LICENSE-MIT)

</div>

---

## Quick Start

```bash
# 1. Generate README.md and LICENSE files from README.yml in current directory
common-facade

# 2. Check for drift in CI without modifying files
common-facade check
```

## Highlights

- **Single source of truth**: Author your metadata, badges, docs, and licenses once in README.yml.
- **Strict standardization**: Required layout with Quick Start, Highlights, Documentation, and License sections.
- **Zero drift in CI**: Opinionated composite GitHub Action verifies project facades with no parameters.
- **Automatic license management**: Current year is always stamped automatically for MIT and Apache-2.0.

## Documentation

- **[Architecture](docs/architecture.md)**: Repository layout, AST-level YAML decoding, and drift detection pipeline.
- **[Development](docs/development.md)**: Building with Go, running test suites, and releasing through GitHub Actions.
- **[Deployment](docs/deployment.md)**: Installation via Homebrew tap, prebuilt release binaries, and GitHub Action integration.

## License

Dual-licensed under [MIT](LICENSE-MIT) or [Apache-2.0](LICENSE-APACHE), at your option.
