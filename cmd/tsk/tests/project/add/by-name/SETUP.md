# Scenario

**Feature**: `tsk project add --project NAME` for registered non-git dir stores name only

```
register seatalk; add --project seatalk "x" -> project.name=seatalk
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	req.Title = "x"
	req.Args = []string{"project", "add", "--project", "seatalk", req.Title}
	return nil
}
```
