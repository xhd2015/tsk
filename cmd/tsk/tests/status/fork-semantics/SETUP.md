# Scenario

**Feature**: summary fork separates no-followup and questions branches; satisfied arrows into done rail

```
create -> tsk status --color <id> -> no followup on horizontal branch, questions on spine, satisfied has ►
```

## Steps

1. Create task at `create` stage.
2. Run `tsk status --color <id>` (full pipeline diagram).

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "fork semantics"
	id := createTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = statusArgs(id, "--color")
	return nil
}
```