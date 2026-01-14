# Stamp

A beautiful CLI for managing Architecture Decision Records (ADRs).

<img alt="stamp demo" src="demo.gif" width="600" />

Inspired by [adr-tools](https://github.com/npryce/adr-tools), built with [Go](https://go.dev) and the [Charmbracelet](https://charm.sh) ecosystem for a delightful terminal experience.

## What are ADRs?

Architecture Decision Records are short documents that capture important architectural decisions made along with their context and consequences. They help teams:

- Remember why decisions were made
- Onboard new team members quickly
- Avoid revisiting the same discussions
- Understand the evolution of a system

Learn more at [adr.github.io](https://adr.github.io/).

## Features

- Create, list, and manage ADRs from the command line
- Beautiful terminal output with colored status badges and styled tables
- Link related ADRs together (supersedes, amends, clarifies)
- Rendered markdown viewing with [glamour](https://github.com/charmbracelet/glamour)
- Open ADRs in your favorite editor
- Self-updating binary

## Installation

### Download binary

Download the latest binary for your platform from [GitHub Releases](https://github.com/stef16robbe/stamp/releases/latest).

### Go install

```bash
go install github.com/stef16robbe/stamp/cmd/stamp@latest
```

### Build from source

```bash
git clone https://github.com/stef16robbe/stamp.git
cd stamp
go build -o stamp ./cmd/stamp
```

## Quick Start

```bash
# Initialize ADR directory in your project
stamp init

# Create a new ADR
stamp new "Use PostgreSQL for persistence"

# List all ADRs
stamp list

# View an ADR
stamp show 1

# Update status
stamp status 2 accepted

# Link ADRs together
stamp link 2 1 supersedes

# Edit an ADR
stamp edit 1
```

## Updating

Stamp can update itself to the latest version:

```bash
# Check current version
stamp version

# Update to latest release
stamp update
```

## Configuration

Stamp stores its configuration in `.stamp.yaml` in your project root:

```yaml
directory: docs/adr
```

## ADR Format

ADRs are stored as Markdown files with the following structure:

```markdown
# 1. Title of Decision

Date: 2026-01-13

## Status

Accepted

## Context

[Why is this decision needed?]

## Decision

[What was decided?]

## Consequences

[What are the implications?]
```

## AI disclaimer

This project has been built alongside with [Claude Code](https://github.com/anthropics/claude-code)

## License

See [LICENSE](LICENSE.md)
