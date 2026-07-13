# Scenario

**Feature**: summary fork separates no-followup and questions; satisfied is vertical; refine is left rail

```
create -> tsk status --color <id> -> no followup on right rail, questions on spine, vertical satisfied, left refine
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
