# Scenario

**Feature**: `tsk note edit --status` replaces status on a progress entry

```
create -> progress add (in-progress) -> note edit --label progress --index 1 --status done -> progress list
```

## Steps

1. `tsk create "demo"`.
2. `tsk progress add --id 1 --status in-progress "investigating"`.
3. `tsk note edit --id 1 --label progress --index 1 --status done "investigation complete"`.
4. `tsk progress list --id 1`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "demo", "", nil)
	runTskOK(t, req, "progress", "add", "--id", "1", "--status", "in-progress", "investigating")
	runTskOK(t, req, "note", "edit", "--id", "1", "--label", "progress", "--index", "1", "--status", "done", "investigation complete")
	req.Args = []string{"progress", "list", "--id", "1"}
	return nil
}
```
