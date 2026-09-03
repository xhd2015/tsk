# Scenario

**Feature**: delete with bad index errors

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	req.Args = []string{"project", "notes", "delete", "--project", "seatalk", "--index", "1"}
	return nil
}
```
