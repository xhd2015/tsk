# Scenario

**Feature**: agent format facts block includes absolute task directory as `dir:` after topic

```
# create stores task under TSK_HOME/inbox/[id]-<slug>; agent status prints dir:
tsk add "add dark mode" -> tsk status --format=agent <id>
# facts: id → title → stage → terminal → topic → dir (absolute path; inbox topic)
```

## Steps

1. Create task with known title `add dark mode` (stage `create`, inbox).
2. Run `tsk status --format=agent <id>`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "add dark mode"
	id := addTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = agentStatusArgs(id)
	return nil
}
```
