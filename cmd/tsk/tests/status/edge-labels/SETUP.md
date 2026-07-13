# Scenario

**Feature**: edge labels sit on correct transitions in the compact pipeline

```
create -> tsk status --color <id> -> claim, research, confirmed, questions, vertical satisfied in order
```

## Steps

1. Create task at `create` stage.
2. Run `tsk status --color <id>` (full pipeline diagram).

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "edge labels"
	id := createTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = statusArgs(id, "--color")
	return nil
}
```