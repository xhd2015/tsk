# Scenario

**Feature**: `tsk search` with no hits prints `0 matches` and exits 0

```
create task -> search missing-id -> 0 matches
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "demo task", "", nil)
	req.Args = []string{"search", "01a02dcc-698a-7380-a62e-b7adf6982edf"}
	return nil
}
```
