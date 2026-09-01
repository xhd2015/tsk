# Scenario

**Feature**: re-running `tsk install pmark` keeps one checker and same wrapper

```
tsk install pmark; tsk install pmark
  -> wrapper unchanged; one checker block in .zshrc
```

## Steps

1. Install pmark twice under sandboxed HOME.
2. Assert via Run of second install (first in Setup).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "install", "pmark")
	req.Args = []string{"install", "pmark"}
	return nil
}
```
