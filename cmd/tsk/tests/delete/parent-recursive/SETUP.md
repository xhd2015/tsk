# Scenario

**Feature**: delete parent with `--recursive` removes subtree

```
create parent + child -> delete --recursive 1 -> both gone
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	parentID := addTask(t, req, "parent work", "", nil)
	runTskOK(t, req, "add", "--parent", fmt.Sprintf("%d", parentID), "child work")
	req.TaskID = parentID
	req.Title = "parent work"
	req.Args = []string{"delete", "--recursive", fmt.Sprintf("%d", parentID)}
	return nil
}
```
