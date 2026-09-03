# Scenario

**Feature**: `--dir` with git origin stores notes under shared project id

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	repo := filepath.Join(req.WorkRoot, "repo")
	initGitRepo(t, repo, "https://github.com/xhd2015/notes-demo.git")
	runTskOK(t, req, "project", "notes", "add", "--dir", repo, "from-dir")
	req.Args = []string{"project", "notes", "list", "--dir", repo}
	return nil
}
```
