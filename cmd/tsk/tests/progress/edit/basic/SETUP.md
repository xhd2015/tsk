# Scenario

**Feature**: progress edit updates one indexed entry's status and text

```
create -> add -> progress edit --index 1 --status done text -> list
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "demo", "", nil)
	runTskOK(t, req, "progress", "add", "--id", "1", "--status", "in-progress", "investigating")
	runTskOK(t, req, "progress", "edit", "--id", "1", "--index", "1", "--status", "done", "investigation complete")
	req.Args = []string{"progress", "list", "--id", "1"}
	return nil
}
```
