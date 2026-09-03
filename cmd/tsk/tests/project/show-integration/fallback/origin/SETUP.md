# Scenario

**Feature**: `tsk show` prints `project: <origin>` when no ledger location and no name

```
add task; update --set-project <unknown-origin> -> show project: origin
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "origin fallback", "", nil)
	runTskOK(t, req, "update", fmt.Sprintf("%d", id), "--set-project", "github.com/xhd2015/unknown-proj")
	req.TaskID = id
	req.Args = []string{"show", fmt.Sprintf("%d", id)}
	return nil
}
```
