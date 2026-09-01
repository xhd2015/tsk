# Scenario

**Feature**: `tsk project add` outside a git repo errors

```
cd outside-any-repo; tsk project add "x" -> Error, exit 1
```

```go
import (
	"os"
	"path/filepath"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.WorkRoot = outsideGitDir(t)
	req.TskHome = filepath.Join(req.WorkRoot, ".tsk")
	if err := os.MkdirAll(req.TskHome, 0o755); err != nil {
		return err
	}
	req.Title = "x"
	req.Args = []string{"project", "add", req.Title}
	return nil
}
```
