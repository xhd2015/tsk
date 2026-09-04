# Scenario

**Feature**: no project notes → tree has no `notes` node

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	projectAddOK(t, req, "one", "--project", "seatalk")
	req.Args = []string{"project", "tree", "--project", "seatalk", "--plain"}
	return nil
}
```
