# Scenario

**Feature**: streaming `project tree` still prints notes under the current project

```
git origin; add; notes add --dir; tree --plain --no-sub-dirs
```

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/xhd2015/dot-pkgs.git")
	projectAddOK(t, req, "one")
	runTskOK(t, req, "project", "notes", "add", "--dir", req.WorkRoot, "from stream")
	req.Args = []string{"project", "tree", "--plain", "--no-sub-dirs"}
	return nil
}
```
