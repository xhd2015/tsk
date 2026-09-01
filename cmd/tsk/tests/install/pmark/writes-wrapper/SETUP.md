# Scenario

**Feature**: `tsk install pmark` writes executable #!/bin/sh forwarder

```
tsk install pmark
  -> stdout "installed ~/.local/bin/pmark"
  -> file body exec tsk project add "$@"
  -> mode executable
  -> ~/.zshrc gains PATH checker
```

## Steps

1. Run `tsk install pmark` with sandboxed HOME.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"install", "pmark"}
	return nil
}
```
