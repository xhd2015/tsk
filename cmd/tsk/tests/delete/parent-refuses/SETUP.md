# Scenario

**Feature**: delete parent without `--recursive` refuses and keeps tree

```
create parent + child -> delete 1 -> error nested tasks; nothing removed
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	parentID := addTask(t, req, "parent work", "", nil)
	runTskOK(t, req, "add", "--parent", fmt.Sprintf("%d", parentID), "child work")
	req.TaskID = parentID
	req.Title = "parent work"
	req.Args = []string{"delete", fmt.Sprintf("%d", parentID)}
	return nil
}
```
