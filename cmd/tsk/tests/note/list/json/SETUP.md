# Scenario

**Feature**: `tsk note list --json` emits a JSON array with labels

```
note add labeled -> list --json
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := createTask(t, req, "json notes", "", nil)
	runTskOK(t, req, "note", "add", "--label", "grok", "--label", "session-id", "--id", fmt.Sprintf("%d", id), "sess-abc")
	req.Args = []string{"note", "list", "--json", "--id", fmt.Sprintf("%d", id)}
	return nil
}
```
