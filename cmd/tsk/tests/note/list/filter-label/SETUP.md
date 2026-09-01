# Scenario

**Feature**: `tsk note list --label` AND-filters

```
unlabeled + grok+session-id -> list --label grok --label session-id -> one
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "filter notes", "", nil)
	idStr := fmt.Sprintf("%d", id)
	runTskOK(t, req, "note", "add", "--id", idStr, "plain")
	runTskOK(t, req, "note", "add", "--label", "grok", "--label", "session-id", "--id", idStr, "sess-abc")
	req.Args = []string{"note", "list", "--label", "grok", "--label", "session-id", "--id", idStr}
	return nil
}
```
