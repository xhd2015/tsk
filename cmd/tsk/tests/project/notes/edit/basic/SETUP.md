# Scenario

**Feature**: edit replaces note text; list shows replacement

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	runTskOK(t, req, "project", "notes", "add", "--project", "seatalk", "original text")
	runTskOK(t, req, "project", "notes", "edit", "--project", "seatalk", "--index", "1", "replaced text")
	req.Args = []string{"project", "notes", "list", "--project", "seatalk"}
	return nil
}
```
