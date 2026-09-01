# Scenario

**Feature**: `tsk project add` stores normalized origin (not name) when git remote exists

```
git init + origin https://github.com/xhd2015/dot-pkgs.git
tsk project add "leave a remark" -> project.origin set
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/xhd2015/dot-pkgs.git")
	req.Title = "leave a remark"
	req.Args = []string{"project", "add", req.Title}
	return nil
}
```
