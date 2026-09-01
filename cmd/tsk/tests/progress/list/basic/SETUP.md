# Scenario

**Feature**: list multiple progress entries with status display

```
create -> add in-progress -> add in-progress -> add blocked -> list
```

## Steps

1. `tsk add "demo"`.
2. `tsk progress add --id 1 --status in-progress "in-progress investigation"`.
3. `tsk progress add --id 1 --status in-progress "optimized fetch"`.
4. `tsk progress add --id 1 --status blocked "waiting on upstream"`.
5. `tsk progress list --id 1`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "demo", "", nil)
	runTskOK(t, req, "progress", "add", "--id", "1", "--status", "in-progress", "in-progress investigation")
	runTskOK(t, req, "progress", "add", "--id", "1", "--status", "in-progress", "optimized fetch")
	runTskOK(t, req, "progress", "add", "--id", "1", "--status", "blocked", "waiting on upstream")
	req.Args = []string{"progress", "list", "--id", "1"}
	return nil
}
```
