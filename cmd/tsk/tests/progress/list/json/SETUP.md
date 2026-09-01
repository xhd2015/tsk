# Scenario

**Feature**: progress list --json outputs structured array

```
create -> add --status -> list --json
```

## Steps

1. `tsk add "demo"`.
2. `tsk progress add --id 1 --status in-progress "investigating"`.
3. `tsk progress list --id 1 --json`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "demo", "", nil)
	runTskOK(t, req, "progress", "add", "--id", "1", "--status", "in-progress", "investigating")
	req.Args = []string{"progress", "list", "--id", "1", "--json"}
	return nil
}
```
