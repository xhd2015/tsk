# Scenario

**Feature**: omit `--name` matches existing row by cwd then already up to date

```
register --name seatalk --cwd DIR; register --cwd DIR (no name) -> up to date
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	req.Args = []string{"project", "register", "--cwd", dir}
	return nil
}
```
