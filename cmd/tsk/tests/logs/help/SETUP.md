# Scenario

**Feature**: `tsk logs -h` prints usage

```
tsk logs -h -> Usage + --all + --limit + --json
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"logs", "-h"}
	return nil
}
```
