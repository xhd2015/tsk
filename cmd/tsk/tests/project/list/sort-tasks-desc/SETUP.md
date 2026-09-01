# Scenario

**Feature**: `tsk project list` sorts by TASKS descending

```
two projects (2 tasks vs 1) -> higher count first
```

## Steps

1. Add two tasks under origin A, one under origin B.
2. List and assert order.

```go
import "path/filepath"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dirA := filepath.Join(req.WorkRoot, "proj-a")
	dirB := filepath.Join(req.WorkRoot, "proj-b")
	initGitRepo(t, dirA, "https://github.com/example/alpha.git")
	initGitRepo(t, dirB, "https://github.com/example/beta.git")
	projectAddOK(t, req, "a1", "--dir", dirA)
	projectAddOK(t, req, "a2", "--dir", dirA)
	projectAddOK(t, req, "b1", "--dir", dirB)
	req.Args = []string{"project", "list"}
	return nil
}
```
