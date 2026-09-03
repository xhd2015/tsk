# Scenario

**Feature**: `tsk show` prints task cwd and `project:` as ledger location (not origin)

```
project add; show <id> -> cwd: <abs>; project: <location path>
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
