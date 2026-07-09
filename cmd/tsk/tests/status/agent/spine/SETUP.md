# Scenario

**Feature**: agent format at create shows spine order, create[doing], future marks, facts, no rectangle chrome

```
tsk create -> tsk status --format=agent <id>
# spine: create[doing] -> (in_process) -> … -> (done); facts id/title/stage/terminal/topic/dir
```

## Steps

1. Create task only (stage `create`, inbox).
2. Run `tsk status --format=agent <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "agent spine"
	id := createTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = agentStatusArgs(id)
	return nil
}
```
