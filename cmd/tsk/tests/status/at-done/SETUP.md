# Scenario

**Feature**: colored status at done highlights the done stage

```
create -> advance chain -> done -> tsk status --color <id>
```

## Steps

1. Advance task through workflow to `done`.
2. Run `tsk status --color <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "finished task"
	id := createTask(t, req, req.Title, "", nil)
	advanceToDone(t, req, id)
	req.TaskID = id
	req.Args = statusArgs(id, "--color")
	return nil
}
```