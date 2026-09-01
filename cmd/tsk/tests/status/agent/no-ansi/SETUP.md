# Scenario

**Feature**: agent format never emits ANSI even when --color is passed

```
tsk add -> tsk status --format=agent --color <id>  # color ignored
```

## Steps

1. Create task at `create`.
2. Run `tsk status --format=agent --color <id>`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "agent no ansi"
	id := addTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = agentStatusArgs(id, "--color")
	return nil
}
```
