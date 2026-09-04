# Scenario

**Feature**: `--dry-run` on a nested parent without `--recursive` still errors

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	parentID := addTask(t, req, "parent work", "", nil)
	runTskOK(t, req, "add", "--parent", fmt.Sprintf("%d", parentID), "child work")
	req.TaskID = parentID
	req.Title = "parent work"
	req.Args = []string{"delete", "--dry-run", fmt.Sprintf("%d", parentID)}
	return nil
}
```
