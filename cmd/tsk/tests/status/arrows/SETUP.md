# Scenario

**Feature**: compact pipeline uses directional arrows on main flow and summary fork

```
create -> tsk status --color <id> -> full diagram with ▼ main flow, branch arrows, refine loop
```

## Steps

1. Create task at `create` stage.
2. Run `tsk status --color <id>` (full pipeline diagram).

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "pipeline arrows"
	id := createTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = statusArgs(id, "--color")
	return nil
}
```