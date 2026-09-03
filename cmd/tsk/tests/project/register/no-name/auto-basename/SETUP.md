# Scenario

**Feature**: omit `--name` → register with basename of dir

```
register --cwd <dir> -> registered <basename>
```

```go
import "path/filepath"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	dir := outsideGitDir(t)
	req.Title = filepath.Base(dir)
	req.Args = []string{"project", "register", "--cwd", dir}
	return nil
}
```
