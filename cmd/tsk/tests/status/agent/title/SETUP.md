# Scenario

**Feature**: agent format facts block includes exact create title after id and before stage

```
# create stores title in task.json; agent status prints it in the leading facts block
tsk create "add dark mode" -> tsk status --format=agent <id>
# facts: id → title → stage → terminal → topic → dir (exact title text; inbox topic)
```

## Steps

1. Create task with known title `add dark mode` (stage `create`, inbox).
2. Run `tsk status --format=agent <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "add dark mode"
	id := createTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = agentStatusArgs(id)
	return nil
}
```
