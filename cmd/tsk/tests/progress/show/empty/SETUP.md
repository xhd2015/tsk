# Scenario

**Feature**: progress show on a task with no entries prints `no progress`

```
create -> progress show -> no progress
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "demo", "", nil)
	req.Args = []string{"progress", "show", "--id", "1"}
	return nil
}
```
