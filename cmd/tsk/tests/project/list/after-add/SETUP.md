# Scenario

**Feature**: `tsk project list` shows auto-seen projects with task counts

```
project add "one"; project list -> table with origin + TASKS 1
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/xhd2015/dot-pkgs.git")
	projectAddOK(t, req, "one")
	req.Args = []string{"project", "list"}
	return nil
}
```
