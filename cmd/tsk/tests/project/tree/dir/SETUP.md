# Scenario

**Feature**: `tsk project tree --dir PATH` resolves the project from PATH (not cwd)

```
add --dir repo; tree --dir repo --plain from non-git WorkRoot -> that project's tasks
```

```go
import "path/filepath"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	repo := filepath.Join(req.WorkRoot, "repo")
	initGitRepo(t, repo, "https://github.com/xhd2015/dot-pkgs.git")
	projectAddOK(t, req, "from-dir", "--dir", repo)
	req.Args = []string{"project", "tree", "--dir", repo, "--plain"}
	return nil
}
```
