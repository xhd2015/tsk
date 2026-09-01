# Scenario

**Feature**: list with resolved project but no open tasks still shows the branch

```
git repo, no project tasks -> list --plain -> 0 tasks, 1 project
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/xhd2015/agent-pro.git")
	req.Args = []string{"project", "tree", "--plain"}
	return nil
}
```
