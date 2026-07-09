# Scenario

**Feature**: status --plain renders ASCII box art without ANSI or theme

```
tsk create -> tsk status --plain <id>
```

## Steps

1. Create task at `create` stage.
2. Run `tsk status --plain <id>`.

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "plain ascii"
	id := createTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = statusArgs(id, "--plain")
	return nil
}
```