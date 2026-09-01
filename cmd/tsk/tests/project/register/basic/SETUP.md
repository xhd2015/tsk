# Scenario

**Feature**: `tsk project register` writes unique name + cwd into projects.json

```
tsk project register --name seatalk --cwd <nongit> -> registered
tsk project list shows seatalk in the NAME column
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	req.Args = []string{"project", "register", "--name", "seatalk", "--cwd", dir}
	return nil
}
```
