# Scenario

**Feature**: `tsk note add` without `--id` errors

```
tsk note add hello -> Error: --id required
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"note", "add", "hello"}
	return nil
}
```
