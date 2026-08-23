# Scenario

**Feature**: `tsk tree --id` prints a pruned branch for one task with notes and progress

```
create task + topic -> progress add + note add -> tree --id
```

## Steps

1. `tsk topic mkdir kb`.
2. `tsk create --topic kb "report"`.
3. `tsk progress add --id 1 --status in-progress "investigating"`.
4. `tsk note add --id 1 "session abc"`.
5. `tsk tree --id 1`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	runTskOK(t, req, "topic", "mkdir", "kb")
	createTask(t, req, "report", "kb", nil)
	runTskOK(t, req, "progress", "add", "--id", "1", "--status", "in-progress", "investigating")
	runTskOK(t, req, "note", "add", "--id", "1", "session abc")
	req.Args = []string{"tree", "--id", "1"}
	return nil
}
```
