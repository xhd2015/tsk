# Scenario

**Feature**: `tsk show` prints `progress:` line when progress entries exist

```
create -> progress add --status -> show
```

## Steps

1. `tsk create "demo"`.
2. `tsk progress add --id 1 --status in-progress "investigating"`.
3. `tsk show 1`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "demo", "", nil)
	runTskOK(t, req, "progress", "add", "--id", "1", "--status", "in-progress", "investigating")
	req.Args = []string{"show", "1"}
	return nil
}
```
