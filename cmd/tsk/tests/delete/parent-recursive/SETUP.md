# Scenario

**Feature**: delete parent with `--recursive` removes subtree

```
create parent + child -> delete --recursive 1 -> both gone
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	parentID := createTask(t, req, "parent work", "", nil)
	runTskOK(t, req, "create", "--parent", fmt.Sprintf("%d", parentID), "child work")
	req.TaskID = parentID
	req.Title = "parent work"
	req.Args = []string{"delete", "--recursive", fmt.Sprintf("%d", parentID)}
	return nil
}
```
