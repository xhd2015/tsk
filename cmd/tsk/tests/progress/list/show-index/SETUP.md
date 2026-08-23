# Scenario

**Feature**: progress list --show-index prefixes 1-based progress indices

```
create -> add x2 -> progress list --show-index
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "demo", "", nil)
	runTskOK(t, req, "progress", "add", "--id", "1", "--status", "in-progress", "first")
	runTskOK(t, req, "progress", "add", "--id", "1", "--status", "blocked", "second")
	req.Args = []string{"progress", "list", "--id", "1", "--show-index"}
	return nil
}
```
