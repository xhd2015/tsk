# Scenario

**Feature**: colored status at create highlights the create stage

```
tsk add -> tsk status --color 1
```

## Steps

1. Create task only (stage `create`).
2. Run `tsk status --color <id>`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "fresh task"
	id := addTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = statusArgs(id, "--color")
	return nil
}
```