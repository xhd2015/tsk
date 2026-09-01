# Scenario

**Feature**: filter progress list by --status

```
create -> add in-progress -> add blocked -> list --status blocked
```

## Steps

1. `tsk add "demo"`.
2. `tsk progress add --id 1 --status in-progress "in-progress"`.
3. `tsk progress add --id 1 --status blocked "blocked"`.
4. `tsk progress list --id 1 --status blocked`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "demo", "", nil)
	runTskOK(t, req, "progress", "add", "--id", "1", "--status", "in-progress", "in-progress")
	runTskOK(t, req, "progress", "add", "--id", "1", "--status", "blocked", "blocked")
	req.Args = []string{"progress", "list", "--id", "1", "--status", "blocked"}
	return nil
}
```
