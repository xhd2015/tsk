# Scenario

**Feature**: `--json` emits a JSON array

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	runTskOK(t, req, "project", "notes", "add", "--project", "seatalk", "hello")
	req.Args = []string{"project", "notes", "list", "--project", "seatalk", "--json"}
	return nil
}
```
