# Scenario

**Feature**: `tsk note add --label session=abc` stores a key=value label

```
create -> note add --label grok --label session=abc --id N text -> list
```

## Steps

1. Create a task.
2. `tsk note add --label grok --label session=abc --id <id> backfill`.
3. `tsk note list --id <id>`.

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "kv note", "", nil)
	runTskOK(t, req, "note", "add", "--label", "grok", "--label", "session=abc", "--id", fmt.Sprintf("%d", id), "backfill")
	req.Args = []string{"note", "list", "--id", fmt.Sprintf("%d", id)}
	return nil
}
```
