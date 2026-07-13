# Scenario

**Feature**: compact pipeline stays within soft width cap (~40 cols; goldens are source of truth)

```
create -> advance x2 (clarification) -> tsk status --color <id>
```

## Steps

1. Advance task to `clarification`.
2. Run `tsk status --color <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "compact width"
	id := createTask(t, req, req.Title, "", nil)
	advanceTask(t, req, id, "")
	advanceTask(t, req, id, "")
	req.TaskID = id
	req.Args = statusArgs(id, "--color")
	return nil
}
```
