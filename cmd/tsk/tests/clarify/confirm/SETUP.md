# Scenario

**Feature**: clarify confirm -y with all questions advances to implementation

```
create -> advance x2 -> clarify add x2 -> clarify confirm -y -> implementation
```

## Steps

1. Create task and advance to `clarification`.
2. Add two questions.
3. Run `tsk clarify confirm <id> -y`.

```go
import (
	"fmt"
)

func Setup(t *testing.T, req *Request) error {
	req.Title = "clarify me"
	id := createTask(t, req, req.Title, "", nil)
	advanceTask(t, req, id, "")
	advanceTask(t, req, id, "")
	runTskOK(t, req, "clarify", "add", fmt.Sprintf("%d", id), "first question?")
	runTskOK(t, req, "clarify", "add", fmt.Sprintf("%d", id), "second question?")
	req.TaskID = id
	req.Args = []string{"clarify", "confirm", fmt.Sprintf("%d", id), "-y"}
	return nil
}
```