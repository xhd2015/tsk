# Scenario

**Feature**: followup from summary writes context file and changes stage

```
pipeline to summary -> tsk followup <id> "please revise scope"
```

## Steps

1. Advance task to `summary`.
2. Run `tsk followup <id> "please revise scope"`.

```go
import (
	"fmt"
)

func Setup(t *testing.T, req *Request) error {
	req.Title = "needs followup"
	req.Message = "please revise scope"
	id := createTask(t, req, req.Title, "", nil)
	advanceTo(t, req, id, "summary")
	req.TaskID = id
	req.Args = []string{"followup", fmt.Sprintf("%d", id), req.Message}
	return nil
}
```