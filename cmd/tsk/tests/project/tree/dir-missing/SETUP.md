# Scenario

**Feature**: `tsk project tree --dir` errors clearly when PATH does not exist

```
project tree --dir missing-dir --plain -> exit 1, resolve --dir error
```

```go
import "path/filepath"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	missing := filepath.Join(req.WorkRoot, "no-such-dir")
	req.Args = []string{"project", "tree", "--dir", missing, "--plain"}
	return nil
}
```
