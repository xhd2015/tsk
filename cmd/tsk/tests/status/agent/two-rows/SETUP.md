# Scenario

**Feature**: agent format always draws a second back-line row with user_followup, refine, and questions (no satisfied on art)

```
tsk create -> tsk status --format=agent <id>
# row2: refine + user_followup + questions under spine; satisfied not drawn
```

## Steps

1. Create task at `create` (satisfied not in next at this stage).
2. Run `tsk status --format=agent <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "agent two rows"
	id := createTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = agentStatusArgs(id)
	return nil
}
```
