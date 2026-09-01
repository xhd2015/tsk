# Scenario

**Feature**: `tsk install --dry-run` plans without writing

```
HOME=WorkRoot; empty ~/.local/bin
tsk install --dry-run pmark
  -> [dry-run] would write / would create …
  -> no wrapper, no .zshrc
```

## Steps

1. Run `tsk install --dry-run pmark` under sandboxed HOME.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"install", "--dry-run", "pmark"}
	return nil
}
```
