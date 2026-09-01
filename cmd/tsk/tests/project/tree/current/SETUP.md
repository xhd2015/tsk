# Scenario

**Feature**: `tsk project tree` shows a tree for the current project

```
project add "one"; project tree --plain -> tree with task leaf
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/xhd2015/dot-pkgs.git")
	projectAddOK(t, req, "one")
	req.Args = []string{"project", "tree", "--plain"}
	return nil
}
```
