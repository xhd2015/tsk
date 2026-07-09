# Scenario

**Feature**: agent format never emits ANSI even when --color is passed

```
tsk create -> tsk status --format=agent --color <id>  # color ignored
```

## Steps

1. Create task at `create`.
2. Run `tsk status --format=agent --color <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "agent no ansi"
	id := createTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = agentStatusArgs(id, "--color")
	return nil
}
```
