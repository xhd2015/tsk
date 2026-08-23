# Scenario

**Feature**: listing progress on a task with no entries is empty

```
create task -> progress list -> 0 entries
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "demo", "", nil)
	req.Args = []string{"progress", "list", "--id", "1"}
	return nil
}
```
