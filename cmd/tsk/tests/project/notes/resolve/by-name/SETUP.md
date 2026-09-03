# Scenario

**Feature**: `--name` resolves a registered project

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	runTskOK(t, req, "project", "notes", "add", "--name", "seatalk", "via-name")
	req.Args = []string{"project", "notes", "list", "--name", "seatalk"}
	return nil
}
```
