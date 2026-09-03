# Scenario

**Feature**: `notes add` without text errors

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	req.Args = []string{"project", "notes", "add", "--project", "seatalk"}
	return nil
}
```
