# Scenario

**Feature**: default list excludes done; `--stage done` shows them

```
add; done --force; list --plain hides; list --stage done --plain shows
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/xhd2015/dot-pkgs.git")
	id := projectAddOK(t, req, "finish-me")
	runTskOK(t, req, "done", "--force", fmt.Sprintf("%d", id))
	req.TaskID = id
	req.Args = []string{"project", "tree", "--plain"}
	return nil
}
```
