# Scenario

**Feature**: delete nested child leaf; parent remains

```
create parent -> create --parent 1 child -> delete 2 -> child gone; parent stays
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	parentID := addTask(t, req, "parent work", "", nil)
	runTskOK(t, req, "add", "--parent", fmt.Sprintf("%d", parentID), "child work")
	childID := maxTaskID(t, req)
	req.TaskID = childID
	req.Title = "child work"
	req.Message = fmt.Sprintf("%d", parentID) // parent id for assert
	req.Args = []string{"delete", fmt.Sprintf("%d", childID)}
	return nil
}
```
