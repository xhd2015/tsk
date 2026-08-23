# Scenario

**Feature**: text is required for progress add

```
tsk progress add --id 1 --status in-progress -> error
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "demo", "", nil)
	req.Args = []string{"progress", "add", "--id", "1", "--status", "in-progress"}
	return nil
}
```
