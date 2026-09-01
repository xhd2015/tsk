# Scenario

**Feature**: `--color` colors only the stage box span; left refine rail stays outside the SGR span

```
create -> … -> implementation -> tsk status --color <id>
# mid row is "│          │ implementation │" — rail │ not inside green box SGR
```

## Steps

1. Create task and advance to `implementation` (box mid-row shares the left refine rail).
2. Run `tsk status --color <id>`.
3. Assert green still on the implementation box; leading left-rail `│` is outside the box SGR.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "color box only"
	id := addTask(t, req, req.Title, "", nil)
	advanceTo(t, req, id, "implementation")
	req.TaskID = id
	req.Args = statusArgs(id, "--color")
	return nil
}
```
