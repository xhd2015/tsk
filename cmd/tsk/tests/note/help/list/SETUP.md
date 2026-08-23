# Scenario

**Feature**: `tsk note list -h` documents `--id`, `--json`, `--limit`

```
tsk note list -h
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"note", "list", "-h"}
	return nil
}
```
