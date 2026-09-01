# Scenario

**Feature**: `--limit 1` prints only the last task note

```
note one; note two -> list --limit 1 -> two only
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := addTask(t, req, "limit notes", "", nil)
	idStr := fmt.Sprintf("%d", id)
	runTskOK(t, req, "note", "add", "--id", idStr, "one")
	runTskOK(t, req, "note", "add", "--id", idStr, "two")
	req.Args = []string{"note", "list", "--limit", "1", "--id", idStr}
	return nil
}
```
