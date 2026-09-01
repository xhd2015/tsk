# Scenario

**Feature**: `tsk note list --label session` matches `session=…` by key

```
plain + session=abc -> list --label session -> one
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "filter key", "", nil)
	idStr := fmt.Sprintf("%d", id)
	runTskOK(t, req, "note", "add", "--id", idStr, "plain")
	runTskOK(t, req, "note", "add", "--label", "grok", "--label", "session=abc", "--id", idStr, "sess-hit")
	req.Args = []string{"note", "list", "--label", "session", "--id", idStr}
	return nil
}
```
