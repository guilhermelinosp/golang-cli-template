# Contributing

## Branch strategy

- `main` — stable, release-ready
- `feat/*` — new features
- `fix/*` — bug fixes
- `chore/*` — maintenance, dependencies, CI

## Commits

All commits must follow [Conventional Commits](https://www.conventionalcommits.org/):

```
feat: add new feature
fix: resolve login issue
docs: update README
refactor: simplify validation logic
chore: bump dependencies
```

The `release.yml` workflow (GoReleaser) uses commit messages to build the changelog.

## Development workflow

```bash
make setup    # once per clone (template bootstrap only)
lefthook install
make fmt      # format
make lint     # golangci-lint
make race     # tests with race detector
make cover    # coverage report
make build    # stamped binary into bin/
```

Architectural rule: `spf13/cobra` may ONLY be imported inside
`internal/cli/cobra/`. `make verify` enforces it — run it before pushing.

## Pull Requests

- PR titles must also follow Conventional Commits
- Keep PRs focused — one feature/fix per PR
- Update tests and documentation as needed
- Ensure all CI checks pass before requesting review

## Issues

- Use the provided issue templates
- Include steps to reproduce for bugs
- Tag with appropriate labels
