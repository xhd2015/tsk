# Scenario

**Feature**: `tsk note add --label grok --label session-id` stores both labels

```
create -> note add --label grok --label session-id --id N sess
```

## Steps

1. Create a task.
2. `tsk note add --label grok --label session-id --id <id> sess-abc`.
3. `tsk note list --id <id>`.

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "report", "", nil)
	runTskOK(t, req, "note", "add", "--label", "grok", "--label", "session-id", "--id", fmt.Sprintf("%d", id), "sess-abc")
	req.Args = []string{"note", "list", "--id", fmt.Sprintf("%d", id)}
	return nil
}
```
