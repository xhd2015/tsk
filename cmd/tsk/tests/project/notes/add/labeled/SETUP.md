# Scenario

**Feature**: labeled project note lists with label brackets

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	runTskOK(t, req, "project", "notes", "add", "--project", "seatalk", "--label", "run", "go test ./...")
	req.Args = []string{"project", "notes", "list", "--project", "seatalk"}
	return nil
}
```
