# Scenario

**Feature**: `--json` includes `notes` on the project object

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	projectAddOK(t, req, "one", "--project", "seatalk")
	runTskOK(t, req, "project", "notes", "add", "--project", "seatalk", "dev command")
	req.Args = []string{"project", "tree", "--project", "seatalk", "--json"}
	return nil
}
```
