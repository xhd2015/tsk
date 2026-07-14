# Scenario

**Feature**: `--exec-if-idle-1h LINE` parsed as single shell command line

```
# shellwords.Parse(LINE) -> argv[] -> exec.Command(argv[0], argv[1:]...)
```

## Context

Exec leaves verify that quoted tokens in LINE survive parsing and reach the shell.

```go
func Setup(t *testing.T, req *Request) error {
	ensureCheckHelpersUsed()
	return nil
}
```