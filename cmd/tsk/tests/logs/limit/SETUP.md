# Scenario

**Feature**: `--limit 1` keeps the last matching mutation

```
add one; add two; logs --limit 1 -> only second add
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "one", "", nil)
	id2 := addTask(t, req, "two", "", nil)
	req.TaskID = id2
	req.Args = []string{"logs", "--limit", "1"}
	return nil
}
```
