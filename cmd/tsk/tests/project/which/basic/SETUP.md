# Scenario

**Feature**: `tsk project which` prints resolved project identity for cwd

```
tsk project which -> project / name / cwd lines
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/xhd2015/wrk.git")
	req.Args = []string{"project", "which"}
	return nil
}
```
