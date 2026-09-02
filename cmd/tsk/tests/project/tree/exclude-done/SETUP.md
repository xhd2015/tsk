# Scenario

**Feature**: default hides terminals; `--done` / `--archived` / `--all` filter stages

```
add open + done + archived; default hides terminals; flags filter as documented
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/xhd2015/dot-pkgs.git")
	openID := projectAddOK(t, req, "active-me")
	doneID := projectAddOK(t, req, "finish-me")
	runTskOK(t, req, "done", fmt.Sprintf("%d", doneID))
	archivedID := projectAddOK(t, req, "shelve-me")
	runTskOK(t, req, "archive", fmt.Sprintf("%d", archivedID))
	req.TaskID = doneID
	req.OpenID = openID
	req.ArchivedID = archivedID
	req.Args = []string{"project", "tree", "--plain"}
	return nil
}
```
