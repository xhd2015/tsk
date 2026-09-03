# Scenario

**Feature**: delete removes the selected note

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	runTskOK(t, req, "project", "notes", "add", "--project", "seatalk", "keep")
	runTskOK(t, req, "project", "notes", "add", "--project", "seatalk", "drop")
	runTskOK(t, req, "project", "notes", "delete", "--project", "seatalk", "--index", "2")
	req.Args = []string{"project", "notes", "list", "--project", "seatalk"}
	return nil
}
```
