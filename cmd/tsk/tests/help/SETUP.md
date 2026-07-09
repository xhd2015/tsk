# Scenario

**Feature**: `--help` / `-h` at root and nested dispatch surfaces usage on stdout

```
# empty args or -h/--help prints topHelp; subcommands print handler help via lessflags.ErrHelp
tsk [no args | -h | --help | <cmd> --help] -> stdout usage; stderr empty; exit 0
```

## Preconditions

- Root has no top-level flags besides help aliases.
- Nested dispatch commands (`topic`, `label`, `clarify`) expose subcommand lists in their help text.

```go
func Setup(t *testing.T, req *Request) error {
	ensureHelpersUsed()
	return nil
}
```