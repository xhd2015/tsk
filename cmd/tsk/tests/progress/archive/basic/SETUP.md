# Scenario

**Feature**: progress archive sets indexed entry status to archived

```
create -> add -> progress archive --index 1 -> list
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "demo", "", nil)
	runTskOK(t, req, "progress", "add", "--id", "1", "--status", "blocked", "waiting")
	runTskOK(t, req, "progress", "archive", "--id", "1", "--index", "1")
	req.Args = []string{"progress", "list", "--id", "1"}
	return nil
}
```
