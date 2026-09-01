# Scenario

**Feature**: inbox tasks with project group under `@` project node

```
git repo; project add "p"; create "solo"; tree --plain
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/xhd2015/dot-pkgs.git")
	runTskOK(t, req, "project", "add", "p")
	addTask(t, req, "solo", "", nil)
	req.Args = []string{"tree", "--plain"}
	return nil
}
```
