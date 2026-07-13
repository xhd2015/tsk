# Scenario

**Feature**: `tsk status --plain` stdout is byte-equal to the committed ASCII golden

```
tsk create -> tsk status --plain <id>
```

## Steps

1. Create task at `create` stage.
2. Run `tsk status --plain <id>`.
3. Compare full stdout to `expected.txt` (exact equality; 1:1 map of unicode golden).

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "plain golden"
	id := createTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = statusArgs(id, "--plain")
	return nil
}
```
