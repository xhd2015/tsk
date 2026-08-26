# Scenario

**Feature**: `tsk search --json` emits a JSON array of hits

```
note session -> search --json id -> JSON array with kind note
```

```go
import "fmt"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	id := createTask(t, req, "json search", "eng/backend", nil)
	runTskOK(t, req, "note", "add", "--label", "grok", "--id", fmt.Sprintf("%d", id), "sess-json-1")
	req.Args = []string{"search", "--json", "sess-json-1"}
	return nil
}
```
