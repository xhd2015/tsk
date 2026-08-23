# Scenario

**Feature**: invalid progress status errors

```
create -> progress add --status started -> error
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "demo", "", nil)
	req.Args = []string{"progress", "add", "--id", "1", "--status", "started", "legacy value"}
	return nil
}
```
