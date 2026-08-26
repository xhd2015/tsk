# Scenario

**Feature**: `tsk search` without a query errors

```
tsk search -> Error: query required
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"search"}
	return nil
}
```
