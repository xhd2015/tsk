# Scenario

**Feature**: dry-run still requires a name

```
tsk install --dry-run -> Error: name required
```

## Steps

1. Run `tsk install --dry-run` with no name.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"install", "--dry-run"}
	return nil
}
```
