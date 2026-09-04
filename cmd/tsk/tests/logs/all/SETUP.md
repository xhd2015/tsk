# Scenario

**Feature**: `--all` includes reads

```
tsk add; tsk list; tsk logs --all -> add and list
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "ship logs", "", nil)
	runTskOK(t, req, "list")
	req.TaskID = id
	req.Args = []string{"logs", "--all"}
	return nil
}
```
