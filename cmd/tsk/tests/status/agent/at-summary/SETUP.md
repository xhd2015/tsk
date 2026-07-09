# Scenario

**Feature**: agent format at summary shows followup and done next options

```
create -> … -> summary -> tsk status --format=agent <id>
```

## Steps

1. Create task and advance to `summary`.
2. Run `tsk status --format=agent <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "agent at summary"
	id := createTask(t, req, req.Title, "", nil)
	advanceTo(t, req, id, "summary")
	req.TaskID = id
	req.Args = agentStatusArgs(id)
	return nil
}
```
