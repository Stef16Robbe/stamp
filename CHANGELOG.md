# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2025-01-28

### Added

- `stamp graph` command - Generate visual graphs of ADR relationships in Mermaid or Graphviz DOT format
- Homebrew tap installation: `brew install stef16robbe/tap/stamp`

## [0.2.1] - 2025-01-14

### Added

- Curl installer script: `curl -fsSL .../install.sh | sh`

### Changed

- Installer uses `~/.local/bin` instead of `/usr/local/bin` (no sudo required)

## [0.2.0] - 2025-01-14

### Added

- Self-update capability: `stamp update` downloads and installs the latest release
- Version command: `stamp version` and `--version` flag
- GitHub Actions for CI (test, vet, build on push/PR)
- GitHub Actions for releases (goreleaser on version tags)
- Unit tests for `adr` and `config` packages

## [0.1.0] - 2025-01-13

### Added

- Initial release
- `stamp init` - Initialize ADR directory with first ADR
- `stamp new <title>` - Create new ADR
- `stamp list` - List all ADRs in a styled table
- `stamp show <number>` - View ADR with rendered markdown
- `stamp status <number> <status>` - Update ADR status
- `stamp link <source> <target> <relation>` - Link related ADRs
- `stamp edit <number>` - Open ADR in editor
- Bidirectional linking (supersedes/superseded-by, amends/amended-by, clarifies/clarified-by)
- Beautiful terminal output with Charmbracelet ecosystem
