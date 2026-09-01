# Scenario

**Feature**: update --clear-project removes project from task

```
project add -> update --clear-project
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	initGitRepo(t, req.WorkRoot, "https://github.com/xhd2015/dot-pkgs.git")
	id := 0
	before := maxTaskID(t, req)
	runTskOK(t, req, "project", "add", "seed")
	id = maxTaskID(t, req)
	if id <= before {
		t.Fatalf("expected new project task")
	}
	req.TaskID = id
	req.Args = []string{"update", fmt.Sprintf("%d", id), "--clear-project"}
	return nil
}
```
