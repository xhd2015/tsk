# Scenario

**Feature**: progress show prints the latest entry

```
create -> add in-progress -> add blocked -> show
```

## Steps

1. `tsk add "demo"`.
2. `tsk progress add --id 1 --status in-progress "first"`.
3. `tsk progress add --id 1 --status blocked "second"`.
4. `tsk progress show --id 1`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "demo", "", nil)
	runTskOK(t, req, "progress", "add", "--id", "1", "--status", "in-progress", "first")
	runTskOK(t, req, "progress", "add", "--id", "1", "--status", "blocked", "second")
	req.Args = []string{"progress", "show", "--id", "1"}
	return nil
}
```
