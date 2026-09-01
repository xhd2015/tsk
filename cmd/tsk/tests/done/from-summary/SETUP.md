# Scenario

**Feature**: done from summary stage

```
pipeline to summary -> tsk done <id> -> *-done-*
```

## Steps

1. Create task and advance through workflow to `summary`.
2. Run `tsk done <id>`.

```go
import (
	"fmt"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Title = "finish line"
	id := addTask(t, req, req.Title, "", nil)
	advanceTo(t, req, id, "summary")
	req.TaskID = id
	req.Args = []string{"done", fmt.Sprintf("%d", id)}
	return nil
}
```