# Scenario

**Feature**: empty project notes list prints `0 notes`

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	req.Args = []string{"project", "notes", "--project", "seatalk"}
	return nil
}
```
