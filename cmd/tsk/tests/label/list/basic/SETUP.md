# Scenario

**Feature**: `tsk label list` dedupes task labels and note label keys

```
create --label report + note grok-session-id=… + progress -> list names
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "label list", "", []string{"report"})
	idStr := fmt.Sprintf("%d", id)
	runTskOK(t, req, "note", "add", "--id", idStr, "--label", "grok-session-id=abc-1", "desc one")
	runTskOK(t, req, "note", "add", "--id", idStr, "--label", "grok-session-id=abc-2", "desc two")
	runTskOK(t, req, "progress", "add", "--id", idStr, "--status", "in-progress", "working")
	req.Args = []string{"label", "list"}
	return nil
}
```
