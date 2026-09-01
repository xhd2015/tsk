# Scenario

**Feature**: advancing a parent updates stage only; nested child paths stay put

```
create parent + child -> advance parent -> parent stage in_process; dirs unchanged
```

```go
import (
	"fmt"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	parentID := addTask(t, req, "add dark mode", "", nil)
	runTskOK(t, req, "add", "--parent", fmt.Sprintf("%d", parentID), "child detail")
	req.TaskID = parentID
	req.Title = "add dark mode"
	req.Args = []string{"advance", fmt.Sprintf("%d", parentID)}
	return nil
}
```
