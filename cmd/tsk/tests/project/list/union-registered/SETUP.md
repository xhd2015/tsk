# Scenario

**Feature**: default `tsk project list` includes registered projects with TASKS 0

```
register seatalk; list -> NAME seatalk, TASKS 0
list --auto -> 0; list --registered -> seatalk (no TASKS col)
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	req.Args = []string{"project", "list"}
	return nil
}
```
