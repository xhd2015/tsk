# Scenario

**Feature**: `--json` emits structured inbox + topics array

```
create inbox task + topic task -> tree --json
```

## Steps

1. `tsk add "solo"` (inbox).
2. `tsk topic mkdir kb`.
3. `tsk add --topic kb "report"`.
4. `tsk tree --json`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	addTask(t, req, "solo", "", nil)
	runTskOK(t, req, "topic", "mkdir", "kb")
	addTask(t, req, "report", "kb", nil)
	runTskOK(t, req, "done", "--force", "1")
	req.Args = []string{"tree", "--json"}
	return nil
}
```
