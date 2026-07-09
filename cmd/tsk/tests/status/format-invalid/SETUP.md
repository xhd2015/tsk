# Scenario

**Feature**: invalid `--format` value fails with exit 1 and allowed values on stderr

```
tsk create -> tsk status --format=nope <id> -> exit 1; stderr lists allowed formats
```

## Steps

1. Create a task (valid id).
2. Run `tsk status --format=nope <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "format invalid"
	id := createTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = statusArgs(id, "--format=nope")
	return nil
}
```
