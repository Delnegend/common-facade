# Deployment & Configuration

## Installation

### Homebrew (Linux & macOS)

```bash
brew tap Delnegend/tap
brew install common-facade
```

### From GitHub Release Binary

Download the compiled archive from [GitHub Releases](https://github.com/Delnegend/common-facade/releases/latest):

```bash
URL=$(curl -fsSL https://api.github.com/repos/Delnegend/common-facade/releases/latest | grep -o 'https://[^"]*linux-amd64\.tar\.xz')
curl -fsSL -o common-facade.tar.xz "$URL" && tar -xJf common-facade.tar.xz --strip-components=1
```

## Configuration Specification

All settings are stored in `README.yml` in the project root:

| Key | Type | Description |
|---|---|---|
| `$schema` | String | Path to `README.schema.json` |
| `name` | String | Project display name |
| `repo` | String | GitHub repository slug (`Delnegend/<repo>`) |
| `tagline` | String | Single bold summary sentence |
| `badges` | Object | Badges to show: `ci`, `release`, `license` (all required) |
| `quick_start` | String | Instructions or shell snippets to get started |
| `highlights` | List | Bullet points highlighting key capabilities |
| `documentation` | Object | Required documentation links: `architecture`, `development`, `deployment` |
| `license` | Object | `type: false` or list containing `MIT`, `APACHE-2.0` |

## GitHub Actions Integration

Use the zero-parameter composite action in your workflow:

```yaml
- name: Validate Facade
  uses: Delnegend/actions/common-facade@v1
```
