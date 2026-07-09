# Scenario

**Feature**: agent format marks mid-pipeline stage as [doing], past bare, future (name)

```
create -> … -> implementation -> tsk status --format=agent <id>
# implementation[doing]; past bare; (verification) (summary) (done)
```

## Steps

1. Create task and advance to `implementation`.
2. Run `tsk status --format=agent <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "agent marks mid"
	id := createTask(t, req, req.Title, "", nil)
	advanceTo(t, req, id, "implementation")
	req.TaskID = id
	req.Args = agentStatusArgs(id)
	return nil
}
```
