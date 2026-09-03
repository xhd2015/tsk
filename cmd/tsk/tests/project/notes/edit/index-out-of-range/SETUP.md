# Scenario

**Feature**: edit with bad index errors

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	runTskOK(t, req, "project", "notes", "add", "--project", "seatalk", "only")
	req.Args = []string{"project", "notes", "edit", "--project", "seatalk", "--index", "9", "x"}
	return nil
}
```
