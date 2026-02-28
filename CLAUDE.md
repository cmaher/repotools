# Repotools -- Agent Guide

## Project Structure

Go project.

- `cmd/repotools/main.go` -- entry point
- `src/cli/` -- cobra command wiring (one file per command)
- `src/` -- domain packages: `runner`, `fs`, `git`, `github`, `metrics`, `beads`

## Architecture

Each `src/cli/*.go` is a thin cobra command that calls into a domain package. Domain packages accept `io.Writer` for output, making them testable without capturing stdout.

## Adding a New Command

1. Create a function in the appropriate domain package under `src/`. It should accept `io.Writer` for output.
2. Create `src/cli/newcmd.go` with a cobra command that calls the domain function.
3. Register the command in `src/cli/root.go` via `cmd.AddCommand(...)`.

## Test Conventions

- Standard library `testing` package only -- no test frameworks.
- Use `t.TempDir()` for temp files.
- Domain packages are testable via `io.Writer` injection.

## Build

```
make build     # builds ./repotools
make test      # runs go test ./...
make install   # copies to ~/bin
```

**IMPORTANT:** Run `make install` after every code change so `~/bin/repotools` stays current.

## Dependencies

cobra only. No other external dependencies.
