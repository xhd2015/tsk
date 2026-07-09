# Scenario

**Feature**: agent format at done is terminal with done[doing] and blocked/empty next

```
create -> … -> done -> tsk status --format=agent <id>
```

## Steps

1. Advance task through workflow to `done`.
2. Run `tsk status --format=agent <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "agent at done"
	id := createTask(t, req, req.Title, "", nil)
	advanceToDone(t, req, id)
	req.TaskID = id
	req.Args = agentStatusArgs(id)
	return nil
}
```
