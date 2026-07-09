# Scenario

**Feature**: agent format at clarification shows blocked advance and clarify-confirm next

```
create -> advance x2 -> clarification -> tsk status --format=agent <id>
```

## Steps

1. Create task and advance to `clarification`.
2. Run `tsk status --format=agent <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "agent at clarification"
	id := createTask(t, req, req.Title, "", nil)
	advanceTo(t, req, id, "clarification")
	req.TaskID = id
	req.Args = agentStatusArgs(id)
	return nil
}
```
