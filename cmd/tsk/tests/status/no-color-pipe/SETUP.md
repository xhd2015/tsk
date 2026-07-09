# Scenario

**Feature**: status without --color on non-TTY stdout omits ANSI escapes

```
create -> advance x2 (clarification) -> tsk status <id>  # stdout captured (piped)
```

## Steps

1. Advance task to `clarification`.
2. Run `tsk status <id>` without `--color` (doctest captures stdout via pipe/buffer).

```go
func Setup(t *testing.T, req *Request) error {
	req.Title = "no color pipe"
	id := createTask(t, req, req.Title, "", nil)
	advanceTask(t, req, id, "")
	advanceTask(t, req, id, "")
	req.TaskID = id
	req.Args = statusArgs(id)
	return nil
}
```