# Scenario

**Feature**: `tsk note list --label session=abc` is exact key=value match

```
session=abc + session=other -> list --label session=abc -> one
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := createTask(t, req, "filter kv", "", nil)
	idStr := fmt.Sprintf("%d", id)
	runTskOK(t, req, "note", "add", "--label", "session=abc", "--id", idStr, "hit")
	runTskOK(t, req, "note", "add", "--label", "session=other", "--id", idStr, "miss")
	req.Args = []string{"note", "list", "--label", "session=abc", "--id", idStr}
	return nil
}
```
