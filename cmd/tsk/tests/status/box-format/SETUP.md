# Scenario

**Feature**: every workflow stage appears inside a labeled box row

```
create -> tsk status --color <id>
```

## Steps

1. Create task at `create` stage.
2. Run `tsk status --color <id>` (full pipeline diagram).

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "box format"
	id := createTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = statusArgs(id, "--color")
	return nil
}
```