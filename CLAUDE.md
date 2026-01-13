# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Stamp is a CLI tool for managing Architecture Decision Records (ADRs). Built in Go using Cobra for CLI and Charmbracelet ecosystem (glamour, lipgloss) for terminal styling.

## Build and Development Commands

```bash
# Build the project
go build ./...

# Build binary
go build -o stamp ./cmd/stamp

# Run tests
go test ./...

# Run a single test
go test -run TestName ./path/to/package

# Format and vet
go fmt ./... && go vet ./...

# Tidy dependencies
go mod tidy
```

## Architecture

```
cmd/stamp/main.go          # Entry point
internal/
├── cmd/                   # Cobra commands (init, new, list, show, status, link, edit)
│   └── root.go           # Root command, other commands register via init()
├── adr/                   # Core ADR logic
│   ├── adr.go            # ADR struct, parsing, serialization, status types
│   ├── store.go          # Filesystem operations (List, Load, Save, FindByNumber)
│   └── template.go       # NewADR factory
├── config/
│   └── config.go         # .stamp.yaml loading/saving, directory resolution
└── ui/
    ├── styles.go         # Shared lipgloss styles, colors, status badges
    └── frame.go          # Frame/box rendering helpers
```

## Key Patterns

- Commands register themselves to `rootCmd` via `init()` functions
- Config file (`.stamp.yaml`) is searched up the directory tree from cwd
- ADR filenames are zero-padded: `0001-title-slug.md`
- Status is stored as first line in `## Status` section, links as subsequent lines
- All UI styling centralized in `internal/ui/` package
- Uses `github.com/goccy/go-yaml` for YAML (not gopkg.in/yaml.v3)

## ADR Statuses

`Draft`, `Proposed`, `Accepted`, `Deprecated`, `Superseded`, `Rejected`

## Link Relations

`supersedes`, `superseded-by`, `amends`, `amended-by`, `clarifies`, `clarified-by`
