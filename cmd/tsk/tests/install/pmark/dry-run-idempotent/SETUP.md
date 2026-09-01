# Scenario

**Feature**: dry-run after live install reports skips only

```
tsk install pmark; tsk install --dry-run pmark
  -> skip wrapper + skip PATH checkers; disk unchanged
```

## Steps

1. Live-install pmark.
2. Run dry-run (Assert checks plan + no second mutation).

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "install", "pmark")
	req.Args = []string{"install", "--dry-run", "pmark"}
	return nil
}
```
