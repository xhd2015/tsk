# Scenario

**Feature**: advancing a parent renames its dir and cascades child index paths

```
create parent + child -> advance parent -> indexes use new parent dirname
```

```go
import (
	"fmt"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	parentID := createTask(t, req, "add dark mode", "", nil)
	runTskOK(t, req, "create", "--parent", fmt.Sprintf("%d", parentID), "child detail")
	req.TaskID = parentID
	req.Title = "add dark mode"
	req.Args = []string{"advance", fmt.Sprintf("%d", parentID)}
	return nil
}
```
