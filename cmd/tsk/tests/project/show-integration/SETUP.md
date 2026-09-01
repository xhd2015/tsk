# Scenario

**Feature**: `tsk show` prints cwd and project for project tasks

```
project add; show <id> -> cwd: and project: lines
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/xhd2015/dot-pkgs.git")
	id := projectAddOK(t, req, "show-me")
	req.TaskID = id
	req.Args = []string{"show", fmt.Sprintf("%d", id)}
	return nil
}
```
