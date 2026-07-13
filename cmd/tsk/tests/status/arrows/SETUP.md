# Scenario

**Feature**: compact pipeline uses directional arrows — main ▼ spine, left refine, right no-followup

```
create -> tsk status --color <id> -> ▼ main flow, ► into clarification (refine), ◄ into done (no followup)
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
