# TODO

Ideas for future development of stamp.

## New Functionality

- [ ] `stamp search <query>` - Full-text search across ADRs (title, content, status)
- [x] `stamp graph` - Generate a visual graph of ADR relationships (Mermaid/Graphviz output)
- [ ] `stamp export` - Export ADRs to HTML, PDF, or a static site for documentation
- [ ] `stamp archive <number>` - Move deprecated/superseded ADRs to an archive folder
- [ ] `stamp template` - Custom templates support (different ADR formats per project)
- [ ] `stamp lint` - Validate ADR format, check for broken links, missing sections
- [ ] `stamp diff <n1> <n2>` - Compare two ADRs side-by-side
- [ ] Interactive mode - Use bubbletea for `stamp new` with prompts for status, related ADRs

## CI/Release Improvements

- [x] Renovate - Auto-update Go dependencies and GitHub Actions
- [x] Release notes automation - Use `release-drafter` for auto-generated notes
- [ ] Code coverage badge - Add codecov.io or coveralls integration
- [x] golangci-lint - More comprehensive static analysis than just `go vet`

## Installation & Distribution

- [x] Homebrew tap - `brew install stef16robbe/tap/stamp`
- [ ] Shell completions - Cobra supports generating bash/zsh/fish/powershell completions
- [ ] Man pages - Auto-generate from Cobra commands
- [ ] asdf plugin - Version manager integration

## Quality of Life

- [ ] `stamp config` - View/edit config from CLI instead of manual YAML editing
- [ ] `--json` flag - Machine-readable output for `list`, `show` (useful for scripting)
- [ ] Git hooks integration - Remind to update ADRs on certain file changes
- [ ] ADR templates in config - Let `.stamp.yaml` define custom sections
- [ ] Date format config - ISO vs locale-specific dates

## Documentation

- [ ] Example workflows - Common ADR patterns (RFC-style, lightweight, etc.)
