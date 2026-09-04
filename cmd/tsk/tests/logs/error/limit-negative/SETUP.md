# Scenario

**Feature**: `--limit -1` is an error

```
tsk logs --limit -1 -> Error: --limit must be >= 0; exit 1
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"logs", "--limit", "-1"}
	return nil
}
```
