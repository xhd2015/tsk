# Scenario

**Feature**: `--all` with `--task` still searches all surfaces (no error)

```
note has unique id -> search --task --all <id> -> note hit (all wins)
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "plain title", "", nil)
	runTskOK(t, req, "note", "add", "--id", fmt.Sprintf("%d", id), "unique-all-wins-token")
	req.Args = []string{"search", "--task", "--all", "unique-all-wins-token"}
	return nil
}
```
