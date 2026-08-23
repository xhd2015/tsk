# Scenario

**Feature**: `--json` emits structured inbox + topics array

```
create inbox task + topic task -> tree --json
```

## Steps

1. `tsk create "solo"` (inbox).
2. `tsk topic mkdir kb`.
3. `tsk create --topic kb "report"`.
4. `tsk tree --json`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	createTask(t, req, "solo", "", nil)
	runTskOK(t, req, "topic", "mkdir", "kb")
	createTask(t, req, "report", "kb", nil)
	runTskOK(t, req, "done", "--force", "1")
	req.Args = []string{"tree", "--json"}
	return nil
}
```
