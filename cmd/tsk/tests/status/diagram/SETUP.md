# Scenario

**Feature**: status at clarification highlights current stage in pipeline

```
create -> advance x2 (clarification) -> tsk status <id>
```

## Steps

1. Advance task to `clarification`.
2. Run `tsk status <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "status check"
	id := createTask(t, req, req.Title, "", nil)
	advanceTask(t, req, id, "")
	advanceTask(t, req, id, "")
	req.TaskID = id
	req.Args = []string{"status", fmt.Sprintf("%d", id)}
	return nil
}
```