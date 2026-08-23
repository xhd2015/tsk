# Scenario

**Feature**: `tsk tree --id` on inbox task with no notes shows just the leaf

```
create inbox task -> tree --id
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "solo", "", nil)
	req.Args = []string{"tree", "--id", "1"}
	return nil
}
```
