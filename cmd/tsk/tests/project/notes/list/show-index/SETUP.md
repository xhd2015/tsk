# Scenario

**Feature**: `--show-index` prefixes 1-based indices

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	runTskOK(t, req, "project", "notes", "add", "--project", "seatalk", "first")
	runTskOK(t, req, "project", "notes", "add", "--project", "seatalk", "second")
	req.Args = []string{"project", "notes", "list", "--project", "seatalk", "--show-index"}
	return nil
}
```
