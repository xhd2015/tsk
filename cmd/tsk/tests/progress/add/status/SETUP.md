# Scenario

**Feature**: add a progress entry with --status, then list shows (status)

```
create task -> progress add --status in-progress -> progress list
```

## Steps

1. `tsk create "demo"`.
2. `tsk progress add --id 1 --status in-progress "in-progress investigation"`.
3. `tsk progress list --id 1`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "demo", "", nil)
	runTskOK(t, req, "progress", "add", "--id", "1", "--status", "in-progress", "in-progress investigation")
	req.Args = []string{"progress", "list", "--id", "1"}
	return nil
}
```
