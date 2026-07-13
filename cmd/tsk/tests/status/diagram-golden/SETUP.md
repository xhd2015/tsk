# Scenario

**Feature**: `tsk status --format=diagram` stdout is byte-equal to the committed unicode golden

```
tsk create -> tsk status --format=diagram <id>  # no --color; art is stage-independent without ANSI
```

## Steps

1. Create task at `create` stage (any stage works; Highlight only adds ANSI when colored).
2. Run `tsk status --format=diagram <id>` without `--color`.
3. Compare full stdout to `expected.txt` (exact equality).

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "diagram golden"
	id := createTask(t, req, req.Title, "", nil)
	req.TaskID = id
	req.Args = statusArgs(id, "--format=diagram")
	return nil
}
```
