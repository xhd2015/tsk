# Scenario

**Feature**: `tsk label list` on empty store prints `0 labels`

```
tsk label list -> 0 labels
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"label", "list"}
	return nil
}
```
