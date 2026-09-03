# Scenario

**Feature**: second register with same name and dir prints already up to date

```
register --name seatalk --cwd DIR; register again -> already up to date
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	runTskOK(t, req, "project", "register", "--name", "seatalk", "--cwd", dir)
	req.Args = []string{"project", "register", "--name", "seatalk", "--cwd", dir}
	return nil
}
```
