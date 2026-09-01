# Scenario

**Feature**: update --set-project attaches origin via basename after project add

```
git + project add seed -> add task -> update --set-project dot-pkgs
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/xhd2015/dot-pkgs.git")
	runTskOK(t, req, "project", "add", "seed")
	id := addTask(t, req, "attach me", "", nil)
	req.TaskID = id
	req.Args = []string{"update", fmt.Sprintf("%d", id), "--set-project", "dot-pkgs"}
	return nil
}
```
