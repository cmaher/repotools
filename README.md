# repotools

CLI tool that batches git, GitHub, filesystem, code metrics, and beads operations into single tool calls. Built for use with Claude Code to minimize permission prompts.

Written in Go.

## Build / Install / Test

```
make build       # builds ./repotools
make install     # copies to ~/bin
make test        # runs go test ./...
```

## Global Flag

`repotools -C <dir> <command> ...` -- change to DIR before running any command.

## Commands

| Command | Description |
|---------|-------------|
| `status` | Git status + recent log in one call |
| `log [base]` | Commits since diverging from base branch |
| `diff [base] [flags]` | Diff vs base branch |
| `ls [base] [-- path...]` | List files at merge base |
| `pr [number] [--only SECTIONS] [--exclude SECTIONS]` | Fetch GitHub PR data |
| `read <file> [start] [end]` | Print numbered lines from a file |
| `multi-ls dir1 dir2 ...` | List contents of multiple directories |
| `multi-find <head_count> [find_opts...] path1 ...` | Find files across multiple directories |
| `loc [flags] <paths...>` | Lines-of-code metrics |
| `fn-spans [flags] <paths...>` | Function/method span extraction |
| `multi-bead [flags]` | Beads issue tracking operations |
| `bead-status` | Beads status overview |

Use `repotools --help` and `repotools <cmd> --help` for details.
