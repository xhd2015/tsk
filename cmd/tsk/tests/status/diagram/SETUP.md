# Scenario

**Feature**: colored status at clarification highlights current stage in compact pipeline

```
create -> advance x2 (clarification) -> tsk status --color <id>
```

## Steps

1. Advance task to `clarification`.
2. Run `tsk status --color <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "status check"
	id := createTask(t, req, req.Title, "", nil)
	advanceTask(t, req, id, "")
	advanceTask(t, req, id, "")
	req.TaskID = id
	req.Args = statusArgs(id, "--color")
	return nil
}
```