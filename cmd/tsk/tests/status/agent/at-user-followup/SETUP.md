# Scenario

**Feature**: agent format at user_followup marks followup[doing] and next refine/done

```
… -> summary -> tsk followup <id> … -> tsk status --format=agent <id>
# user_followup[doing] on row 2; spine through summary bare; (done)
```

## Steps

1. Create task, advance to `summary`, run `tsk followup`.
2. Run `tsk status --format=agent <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "agent at user followup"
	id := createTask(t, req, req.Title, "", nil)
	advanceTo(t, req, id, "summary")
	runTskOK(t, req, "followup", fmt.Sprintf("%d", id), "please revise scope")
	req.TaskID = id
	req.Args = agentStatusArgs(id)
	return nil
}
```
