# Scenario

**Feature**: project origin path is grey when tree --color is on

```
git repo; project add; tree --color -> short name plain, origin grey
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/xhd2015/dot-pkgs.git")
	runTskOK(t, req, "project", "add", "p")
	req.Args = []string{"tree", "--color"}
	return nil
}
```
