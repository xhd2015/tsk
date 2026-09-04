# Scenario

**Feature**: two notes list oldest first under `notes`

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	projectAddOK(t, req, "one", "--project", "seatalk")
	runTskOK(t, req, "project", "notes", "add", "--project", "seatalk", "first remark")
	runTskOK(t, req, "project", "notes", "add", "--project", "seatalk", "second remark")
	req.Args = []string{"project", "tree", "--project", "seatalk", "--plain"}
	return nil
}
```
