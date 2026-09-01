# Scenario

**Feature**: `tsk project tree --all` shows a multi-project forest

```
add --dir repo-a; add --dir repo-b; list --all --plain from non-git WorkRoot
```

```go
import "path/filepath"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	repoA := filepath.Join(req.WorkRoot, "repo-a")
	repoB := filepath.Join(req.WorkRoot, "repo-b")
	initGitRepo(t, repoA, "https://github.com/xhd2015/dot-pkgs.git")
	initGitRepo(t, repoB, "https://github.com/xhd2015/wrk.git")
	projectAddOK(t, req, "from-a", "--dir", repoA)
	projectAddOK(t, req, "from-b", "--dir", repoB)
	req.Args = []string{"project", "tree", "--all", "--plain"}
	return nil
}
```
